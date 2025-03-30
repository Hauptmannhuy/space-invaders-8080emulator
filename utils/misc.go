package misc

import "log"

func RegToString(code uint8) string {
	switch code & 0b111 {
	case 0x07:
		return "A"
	case 0x00:
		return "B"
	case 0x01:
		return "C"
	case 0x02:
		return "D"
	case 0x03:
		return "E"
	case 0x04:
		return "H"
	case 0x05:
		return "L"

	default:
		log.Fatal("Invalid register code")
		return ""
	}
}

func RegPairToString(code uint8) string {
	switch code & 0b11 {
	case 0x00:
		return "BC"
	case 0x01:
		return "DE"
	case 0x02:
		return "HL"
	case 0x03:
		return "SP"
	default:
		log.Fatal("Invalid register pair code")
		return ""
	}
}

func Make16bit(hi, lo uint8) uint16 {
	return (uint16(hi) << 8) | uint16(lo)
}
