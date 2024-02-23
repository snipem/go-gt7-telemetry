package gt7

import (
	"encoding/binary"
	"fmt"
	"golang.org/x/crypto/salsa20"
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

func (gt7c *GT7Communication) SendHB(conn *net.UDPConn) error {
	_, err := conn.WriteToUDP([]byte("A"), &net.UDPAddr{
		IP:   net.ParseIP(gt7c.playstationIP),
		Port: gt7c.sendPort,
	})
	if err != nil {
		return fmt.Errorf("error sending heart beat: %v", err)
	}
	err = conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		return fmt.Errorf("error setting read deadline: %v", err)
	}
	return nil
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

func (gt7c *GT7Communication) Start() error {
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 5055,
	})
	if err != nil {
		return fmt.Errorf("error starting connection: %v", err)
	}
	defer conn.Close()

	return nil
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

func (gt7c *GT7Communication) Run() error {

	for gt7c.shallRun {
		s := &net.UDPConn{}
		gt7c.shallRestart = false
		addr, err := net.ResolveUDPAddr("udp", ":33740")
		if err != nil {
			return fmt.Errorf("error resolving address: %v", err)
		}
		s, err = net.ListenUDP("udp", addr)
		if err != nil {
			return fmt.Errorf("error listening on udp %s: %v", addr, err)
		}
		defer s.Close()

		gt7c.SendHB(s)
		if err != nil {
			return fmt.Errorf("error sending heart beat: %v", addr)
		}
		packageID := 0
		packageNr := 0
		for !gt7c.shallRestart && gt7c.shallRun {
			buffer := make([]byte, 4096)
			n, _, err := s.ReadFromUDP(buffer)
			if err != nil {
				gt7c.SendHB(s)
				packageNr = 0
				continue
			}
			packageNr++
			ddata := salsa20Dec(buffer[:n])
			if len(ddata) > 0 && int(binary.LittleEndian.Uint32(ddata[0x70:0x70+4])) > packageID {
				gt7c.lastTimeDataReceived = time.Now()
				packageID = int(binary.LittleEndian.Uint32(ddata[0x70 : 0x70+4]))

				curlap := int(binary.LittleEndian.Uint16(ddata[0x74 : 0x74+2]))

				gt7c.LastData = NewGTData(ddata)

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
	return nil
}
