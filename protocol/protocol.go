package protocol

import "net"

type Color struct {
	R, G, B uint8
}

type Light struct {
	Active   Color
	Inactive Color
}

const (
	MessageBegin = 0x0B
	MessageEnd   = 0x0A

	MessageTypePing                    = 0x0C
	MessageTypePong                    = 0x0D
	MessageTypeSensorsStatus           = 0x0E
	MessageTypeSensorsRequest          = 0x0F
	MessageTypeLightsStatus            = 0x10
	MessageTypeLightsSet               = 0x11
	MessageTypeLightsRequest           = 0x12
	MessageTypeLightsBrightnessSet     = 0x13
	MessageTypeLightsBrightnessRequest = 0x14
)

func FormatMessage(payload []byte) []byte {
	length := len(payload) - 1
	// resulting message should be "start of message" + "length of message" + "payload" + "end of message"
	return append([]byte{MessageBegin, byte(length)}, append(payload, MessageEnd)...)
}

func SendMessage(conn net.Conn, payload []byte) error {
	_, err := conn.Write(FormatMessage(payload))
	return err
}

func Ping() []byte {
	return FormatMessage([]byte{MessageTypePing})
}

func RequestSensorsStatus() []byte {
	return FormatMessage([]byte{MessageTypeSensorsRequest})
}

func RequestLightsStatus() []byte {
	return FormatMessage([]byte{MessageTypeLightsRequest})
}

func SetBrightness(brightness uint8) []byte {
	return FormatMessage([]byte{MessageTypeLightsBrightnessSet, brightness})
}

func SetLights(lights map[int]Light) []byte {
	payload := []byte{MessageTypeLightsSet}
	for i, light := range lights {
		payload = append(payload, byte(i))
		payload = append(payload, light.Active.R, light.Active.G, light.Active.B)
		payload = append(payload, light.Inactive.R, light.Inactive.G, light.Inactive.B)
	}
	return FormatMessage(payload)
}
