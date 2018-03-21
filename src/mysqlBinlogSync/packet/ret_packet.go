package packet

//对OKPacker和ErrPacket的封装
type RetPacket struct {
	*OKPacket
	*ErrPacket
	IsOk bool
}
