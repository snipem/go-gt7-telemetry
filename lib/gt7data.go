package gt7

import (
	"encoding/binary"
	"fmt"
	"math"
)

type GTData struct {
	PackageID            int32
	BestLap              int32
	LastLap              int32
	CurrentLap           int16
	CurrentGear          uint8
	SuggestedGear        uint8
	FuelCapacity         float32
	CurrentFuel          float32
	Boost                float32
	TyreDiameterFL       float32
	TyreDiameterFR       float32
	TyreDiameterRL       float32
	TyreDiameterRR       float32
	TypeSpeedFL          float32
	TypeSpeedFR          float32
	TypeSpeedRL          float32
	TyreSpeedRR          float32
	CarSpeed             float32
	TyreSlipRatioFL      float32
	TyreSlipRatioFR      float32
	TyreSlipRatioRL      float32
	TyreSlipRatioRR      float32
	TimeOnTrack          Duration
	TotalLaps            int16
	CurrentPosition      int16
	TotalPositions       int16
	CarID                int32
	Throttle             float32
	RPM                  float32
	RPMRevWarning        uint16
	Brake                float32
	RPMRevLimiter        uint16
	EstimatedTopSpeed    int16
	Clutch               float32
	ClutchEngaged        float32
	RPMAfterClutch       float32
	OilTemp              float32
	WaterTemp            float32
	OilPressure          float32
	RideHeight           float32
	TyreTempFL           float32
	TyreTempFR           float32
	SuspensionFL         float32
	SuspensionFR         float32
	TyreTempRL           float32
	TyreTempRR           float32
	SuspensionRL         float32
	SuspensionRR         float32
	Gear1                float32
	Gear2                float32
	Gear3                float32
	Gear4                float32
	Gear5                float32
	Gear6                float32
	Gear7                float32
	Gear8                float32
	PositionX            float32
	PositionY            float32
	PositionZ            float32
	VelocityX            float32
	VelocityY            float32
	VelocityZ            float32
	RotationPitch        float32
	RotationYaw          float32
	RotationRoll         float32
	AngularVelocityX     float32
	AngularVelocityY     float32
	AngularVelocityZ     float32
	IsPaused             bool
	InRace               bool
	IsLoading            bool
	IsInGear             bool
	CarHasTurbo          bool
	IsRevLimiterFlashing bool
	IsHandbrakeEngaged   bool
	IsLightsOn           bool
	IsLowBeamOn          bool
	IsHighBeamOn         bool
	IsASMEngaged         bool
	IsTCSEngaged         bool
	// Following flags are unknown
	// See https://www.gtplanet.net/forum/threads/gt7-is-compatible-with-motion-rig.410728/page-4#post-13799643
	// for more information
	UnknownFlag12 byte
	UnknownFlag13 byte
	UnknownFlag14 byte
	UnknownFlag15 byte
}

type Duration struct {
	Seconds int
}

func (d Duration) String() string {
	return fmt.Sprintf("%d seconds", d.Seconds)
}

func NewGTData(ddata []byte) GTData {
	data := GTData{}

	if len(ddata) == 0 {
		return data
	}

	data.PackageID = int32(binary.LittleEndian.Uint32(ddata[0x70 : 0x70+4]))
	data.BestLap = int32(binary.LittleEndian.Uint32(ddata[0x78 : 0x78+4]))
	data.LastLap = int32(binary.LittleEndian.Uint32(ddata[0x7C : 0x7C+4]))
	data.CurrentLap = int16(binary.LittleEndian.Uint16(ddata[0x74 : 0x74+2]))
	data.CurrentGear = uint8(ddata[0x90 : 0x90+1][0]) & 0b00001111
	data.SuggestedGear = uint8(ddata[0x90 : 0x90+1][0]) >> 4
	data.FuelCapacity = math.Float32frombits(binary.LittleEndian.Uint32(ddata[0x48 : 0x48+4]))
	data.CurrentFuel = math.Float32frombits(binary.LittleEndian.Uint32(ddata[0x44 : 0x44+4]))
	data.Boost = math.Float32frombits(binary.LittleEndian.Uint32(ddata[0x50:0x50+4])) - 1

	data.TyreDiameterFL = math.Float32frombits(binary.LittleEndian.Uint32(ddata[0xB4 : 0xB4+4]))
	data.TyreDiameterFR = math.Float32frombits(binary.LittleEndian.Uint32(ddata[0xB8 : 0xB8+4]))
	data.TyreDiameterRL = math.Float32frombits(binary.LittleEndian.Uint32(ddata[0xBC : 0xBC+4]))
	data.TyreDiameterRR = math.Float32frombits(binary.LittleEndian.Uint32(ddata[0xC0 : 0xC0+4]))

	data.TypeSpeedFL = float32(math.Abs(float64(3.6 * data.TyreDiameterFL * math.Float32frombits(binary.LittleEndian.Uint32(ddata[0xA4:0xA4+4])))))
	data.TypeSpeedFR = float32(math.Abs(float64(3.6 * data.TyreDiameterFR * math.Float32frombits(binary.LittleEndian.Uint32(ddata[0xA8:0xA8+4])))))
	data.TypeSpeedRL = float32(math.Abs(float64(3.6 * data.TyreDiameterRL * math.Float32frombits(binary.LittleEndian.Uint32(ddata[0xAC:0xAC+4])))))
	data.TyreSpeedRR = float32(math.Abs(float64(3.6 * data.TyreDiameterRR * math.Float32frombits(binary.LittleEndian.Uint32(ddata[0xB0:0xB0+4])))))

	data.CarSpeed = 3.6 * math.Float32frombits(binary.LittleEndian.Uint32(ddata[0x4C:0x4C+4]))

	if data.CarSpeed > 0 {
		data.TyreSlipRatioFL = data.TypeSpeedFL / data.CarSpeed
		data.TyreSlipRatioFR = data.TypeSpeedFR / data.CarSpeed
		data.TyreSlipRatioRL = data.TypeSpeedRL / data.CarSpeed
		data.TyreSlipRatioRR = data.TyreSpeedRR / data.CarSpeed
	}

	data.TimeOnTrack = Duration{
		Seconds: int(binary.LittleEndian.Uint32(ddata[0x80:0x80+4])) / 1000,
	}

	data.TotalLaps = int16(binary.LittleEndian.Uint16(ddata[0x76 : 0x76+2]))
	data.CurrentPosition = int16(binary.LittleEndian.Uint16(ddata[0x84 : 0x84+2]))
	data.TotalPositions = int16(binary.LittleEndian.Uint16(ddata[0x86 : 0x86+2]))
	data.CarID = int32(binary.LittleEndian.Uint32(ddata[0x124 : 0x124+4]))
	data.Throttle = float32(ddata[0x91 : 0x91+1][0]) / 2.55
	data.RPM = math.Float32frombits(binary.LittleEndian.Uint32(ddata[0x3C : 0x3C+4]))
	data.RPMRevWarning = binary.LittleEndian.Uint16(ddata[0x88 : 0x88+2])
	data.Brake = float32(ddata[0x92 : 0x92+1][0]) / 2.55
	data.Boost = math.Float32frombits(binary.LittleEndian.Uint32(ddata[0x50:0x50+4])) - 1
	data.RPMRevLimiter = binary.LittleEndian.Uint16(ddata[0x8A : 0x8A+2])
	data.EstimatedTopSpeed = int16(binary.LittleEndian.Uint16(ddata[0x8C : 0x8C+2]))
	data.Clutch = math.Float32frombits(binary.LittleEndian.Uint32(ddata[0xF4 : 0xF4+4]))
	data.ClutchEngaged = math.Float32frombits(binary.LittleEndian.Uint32(ddata[0xF8 : 0xF8+4]))
	data.RPMAfterClutch = math.Float32frombits(binary.LittleEndian.Uint32(ddata[0xFC : 0xFC+4]))
	data.OilTemp = math.Float32frombits(binary.LittleEndian.Uint32(ddata[0x5C : 0x5C+4]))
	data.WaterTemp = math.Float32frombits(binary.LittleEndian.Uint32(ddata[0x58 : 0x58+4]))
	data.OilPressure = math.Float32frombits(binary.LittleEndian.Uint32(ddata[0x54 : 0x54+4]))
	data.RideHeight = 1000 * math.Float32frombits(binary.LittleEndian.Uint32(ddata[0x38:0x38+4]))
	data.TyreTempFL = math.Float32frombits(binary.LittleEndian.Uint32(ddata[0x60 : 0x60+4]))
	data.TyreTempFR = math.Float32frombits(binary.LittleEndian.Uint32(ddata[0x64 : 0x64+4]))
	data.SuspensionFL = math.Float32frombits(binary.LittleEndian.Uint32(ddata[0xC4 : 0xC4+4]))
	data.SuspensionFR = math.Float32frombits(binary.LittleEndian.Uint32(ddata[0xC8 : 0xC8+4]))
	data.TyreTempRL = math.Float32frombits(binary.LittleEndian.Uint32(ddata[0x68 : 0x68+4]))
	data.TyreTempRR = math.Float32frombits(binary.LittleEndian.Uint32(ddata[0x6C : 0x6C+4]))
	data.SuspensionRL = math.Float32frombits(binary.LittleEndian.Uint32(ddata[0xCC : 0xCC+4]))
	data.SuspensionRR = math.Float32frombits(binary.LittleEndian.Uint32(ddata[0xD0 : 0xD0+4]))
	data.Gear1 = math.Float32frombits(binary.LittleEndian.Uint32(ddata[0x104 : 0x104+4]))
	data.Gear2 = math.Float32frombits(binary.LittleEndian.Uint32(ddata[0x108 : 0x108+4]))
	data.Gear3 = math.Float32frombits(binary.LittleEndian.Uint32(ddata[0x10C : 0x10C+4]))
	data.Gear4 = math.Float32frombits(binary.LittleEndian.Uint32(ddata[0x110 : 0x110+4]))
	data.Gear5 = math.Float32frombits(binary.LittleEndian.Uint32(ddata[0x114 : 0x114+4]))
	data.Gear6 = math.Float32frombits(binary.LittleEndian.Uint32(ddata[0x118 : 0x118+4]))
	data.Gear7 = math.Float32frombits(binary.LittleEndian.Uint32(ddata[0x11C : 0x11C+4]))
	data.Gear8 = math.Float32frombits(binary.LittleEndian.Uint32(ddata[0x120 : 0x120+4]))
	data.PositionX = math.Float32frombits(binary.LittleEndian.Uint32(ddata[0x2C : 0x2C+4]))
	data.PositionY = math.Float32frombits(binary.LittleEndian.Uint32(ddata[0x30 : 0x30+4]))
	data.PositionZ = math.Float32frombits(binary.LittleEndian.Uint32(ddata[0x34 : 0x34+4]))
	data.VelocityX = math.Float32frombits(binary.LittleEndian.Uint32(ddata[0x40 : 0x40+4]))
	data.VelocityY = math.Float32frombits(binary.LittleEndian.Uint32(ddata[0x44 : 0x44+4]))
	data.VelocityZ = math.Float32frombits(binary.LittleEndian.Uint32(ddata[0x48 : 0x48+4]))
	data.RotationPitch = math.Float32frombits(binary.LittleEndian.Uint32(ddata[0x10 : 0x10+4]))
	data.RotationYaw = math.Float32frombits(binary.LittleEndian.Uint32(ddata[0x14 : 0x14+4]))
	data.RotationRoll = math.Float32frombits(binary.LittleEndian.Uint32(ddata[0x18 : 0x18+4]))
	data.AngularVelocityX = math.Float32frombits(binary.LittleEndian.Uint32(ddata[0x1C : 0x1C+4]))
	data.AngularVelocityY = math.Float32frombits(binary.LittleEndian.Uint32(ddata[0x20 : 0x20+4]))
	data.AngularVelocityZ = math.Float32frombits(binary.LittleEndian.Uint32(ddata[0x24 : 0x24+4]))

	flags := ddata[0x8E : 0x8E+2]

	data.InRace = flags[0]&1 == 1
	data.IsPaused = (flags[0]>>1)&1 == 1
	data.IsLoading = (flags[0]>>2)&1 == 1
	data.IsInGear = (flags[0]>>3)&1 == 1
	data.CarHasTurbo = (flags[0]>>4)&1 == 1
	data.IsRevLimiterFlashing = (flags[0]>>5)&1 == 1
	data.IsHandbrakeEngaged = (flags[0]>>6)&1 == 1
	data.IsLightsOn = (flags[0]>>7)&1 == 1

	data.IsLowBeamOn = (flags[1])&1 == 1
	data.IsHighBeamOn = (flags[1]>>1)&1 == 1
	data.IsASMEngaged = (flags[1]>>2)&1 == 1
	data.IsTCSEngaged = (flags[1]>>3)&1 == 1

	data.UnknownFlag12 = (flags[1] >> 4) & 1
	data.UnknownFlag13 = (flags[1] >> 5) & 1
	data.UnknownFlag14 = (flags[1] >> 6) & 1
	data.UnknownFlag15 = (flags[1] >> 7) & 1

	return data
}
