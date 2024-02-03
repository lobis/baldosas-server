package protocol

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

func formatMessage(payload []byte) []byte {
	length := len(payload) - 1
	// resulting message should be "start of message" + "length of message" + "payload" + "end of message"
	return append([]byte{MessageBegin, byte(length)}, append(payload, MessageEnd)...)
}

func Ping() []byte {
	return formatMessage([]byte{MessageTypePing})
}
