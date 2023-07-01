package gt7

import (
	"encoding/binary"
	"fmt"
	"golang.org/x/crypto/salsa20"
	"io"
	"log"
	"net"
	"time"
)

type Session struct {
	SpecialPacketTime int64
	BestLap           int32
	MinBodyHeight     int
	MaxSpeed          int
}
type GT7Communication struct {
	playstationIP          string
	sendPort, receivePort  int
	lastTimeDataReceived   time.Time
	currentLap             Lap
	session                Session
	laps                   []Lap
	LastData               GTData
	alwaysRecordData       bool
	shallRun, shallRestart bool
}

func (gt7c *GT7Communication) SendHB(conn *net.UDPConn) {
	log.Println("Sending heartbeat")
	_, err := conn.WriteToUDP([]byte("A"), &net.UDPAddr{
		IP:   net.ParseIP(gt7c.playstationIP),
		Port: gt7c.sendPort,
	})
	if err != nil {
		log.Fatal(err)
	}
	err = conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		log.Fatal(err)
	}
}

func salsa20Dec(dat []byte) []byte {
	var key [32]byte
	key = [32]byte([]byte("Simulator Interface Packet GT7 ver 0.0"))

	oiv := dat[0x40:0x44]
	iv1 := binary.LittleEndian.Uint32(oiv)
	iv2 := iv1 ^ 0xDEADBEAF
	iv := make([]byte, 8)
	binary.LittleEndian.PutUint32(iv, iv2)
	binary.LittleEndian.PutUint32(iv[4:], iv1)
	ddata := make([]byte, len(dat))
	salsa20.XORKeyStream(ddata, dat, iv, &key)
	magic := binary.LittleEndian.Uint32(ddata[:4])
	if magic != 0x47375330 {
		return nil
	}
	return ddata
}

func (gt7c *GT7Communication) Start() {
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 5055,
	})
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer func(conn *net.UDPConn) {
		err := conn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(conn)

}

type Lap struct {
}

func NewGT7Communication(playstationIP string) *GT7Communication {
	return &GT7Communication{
		playstationIP: playstationIP,
		sendPort:      33739,
		receivePort:   33740,
		LastData:      GTData{},
		shallRun:      true,
		shallRestart:  false,
	}
}

func (gt7c *GT7Communication) Stop() {
	gt7c.shallRun = false
}

func (gt7c *GT7Communication) Run() {

	// Disable Log Output
	log.SetOutput(io.Discard)

	for gt7c.shallRun {
		s := &net.UDPConn{}
		gt7c.shallRestart = false
		addr, err := net.ResolveUDPAddr("udp", ":33740")
		log.Println(addr)
		if err != nil {
			log.Fatal(err)
		}
		s, err = net.ListenUDP("udp", addr)
		if err != nil {
			log.Fatal(err)
		}
		defer s.Close()

		gt7c.SendHB(s)
		if err != nil {
			log.Fatal(err)
		}
		previousLap := -1
		packageID := 0
		packageNr := 0
		for !gt7c.shallRestart && gt7c.shallRun {
			buffer := make([]byte, 4096)
			log.Println("Reading from udp")
			n, _, err := s.ReadFromUDP(buffer)
			if err != nil {
				gt7c.SendHB(s)
				packageNr = 0
				continue
			}
			packageNr++
			log.Println("Package nr: ", packageNr)
			ddata := salsa20Dec(buffer[:n])
			if len(ddata) > 0 && int(binary.LittleEndian.Uint32(ddata[0x70:0x70+4])) > packageID {
				gt7c.lastTimeDataReceived = time.Now()
				packageID = int(binary.LittleEndian.Uint32(ddata[0x70 : 0x70+4]))

				bstlap := int(binary.LittleEndian.Uint32(ddata[0x78 : 0x78+4]))
				lstlap := int(binary.LittleEndian.Uint32(ddata[0x7C : 0x7C+4]))
				curlap := int(binary.LittleEndian.Uint16(ddata[0x74 : 0x74+2]))

				gt7c.LastData = NewGTData(ddata)
				log.Printf("Speed: %d\n", gt7c.LastData.CarSpeed)

				log.Println(packageID, previousLap, bstlap, lstlap, curlap)

				if curlap == 0 {
					gt7c.session.SpecialPacketTime = 0
				}

				if packageNr > 100 {
					gt7c.SendHB(s)
					packageNr = 0
				}
			}
		}
	}
}
