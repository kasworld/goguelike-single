// Code generated by "genprotocol.exe -ver=4ab5c823f1f31871cd0cef53b64b0920dee9d59673878aa337f0c3cb12b8a163 -basedir=protocol_c2t -prefix=c2t -statstype=int"

package c2t_packet

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_const"
	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_error"
	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_idcmd"
	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_idnoti"
)

type FlowType byte // packet flow type

const (
	invalid      FlowType = iota // make uninitalized packet error
	Request                      // Request for request packet (response packet expected)
	Response                     // Response is reply of request packet
	Notification                 // Notification is just send and forget packet
)

var _FlowType_str = map[FlowType]string{
	invalid:      "invalid",
	Request:      "Request",
	Response:     "Response",
	Notification: "Notification",
}

func (e FlowType) String() string {
	if s, exist := _FlowType_str[e]; exist {
		return s
	}
	return fmt.Sprintf("FlowType%d", byte(e))
}

///////////////////////////////////////////////////////////////////////////////

const (
	// HeaderLen fixed size of header
	HeaderLen = 4 + 4 + 2 + 2 + 1 + 1 + 2

	// MaxPacketLen max total packet size byte of raw packet
	MaxPacketLen = HeaderLen + c2t_const.MaxBodyLen
)

func (pk Packet) String() string {
	return fmt.Sprintf("Packet[%v %+v]", pk.Header, pk.Body)
}

// Packet is header + body as object (not byte list)
type Packet struct {
	Header Header
	Body   interface{}
}

func (h Header) String() string {
	switch h.FlowType {
	default:
		return fmt.Sprintf(
			"Header[%v:%v ID:%v Error:%v BodyLen:%v BodyType:%v Fill:%v]",
			h.FlowType, h.Cmd, h.ID, h.ErrorCode, h.bodyLen, h.bodyType, h.Fill)
	case invalid:
		return fmt.Sprintf(
			"Header[%v:%v ID:%v Error:%v BodyLen:%v BodyType:%v Fill:%v]",
			h.FlowType, h.Cmd, h.ID, h.ErrorCode, h.bodyLen, h.bodyType, h.Fill)
	case Request:
		return fmt.Sprintf(
			"Header[%v:%v ID:%v Error:%v BodyLen:%v BodyType:%v Fill:%v]",
			h.FlowType, c2t_idcmd.CommandID(h.Cmd), h.ID, h.ErrorCode, h.bodyLen, h.bodyType, h.Fill)
	case Response:
		return fmt.Sprintf(
			"Header[%v:%v ID:%v Error:%v BodyLen:%v BodyType:%v Fill:%v]",
			h.FlowType, c2t_idcmd.CommandID(h.Cmd), h.ID, h.ErrorCode, h.bodyLen, h.bodyType, h.Fill)
	case Notification:
		return fmt.Sprintf(
			"Header[%v:%v ID:%v Error:%v BodyLen:%v BodyType:%v Fill:%v]",
			h.FlowType, c2t_idnoti.NotiID(h.Cmd), h.ID, h.ErrorCode, h.bodyLen, h.bodyType, h.Fill)
	}
}

// Header is fixed size header of packet
type Header struct {
	bodyLen   uint32              // set at marshal(Packet2Bytes)
	ID        uint32              // sender set, unique id per packet (wrap around reuse)
	Cmd       uint16              // sender set, application demux received packet
	ErrorCode c2t_error.ErrorCode // sender set, Response error
	FlowType  FlowType            // sender set, flow control, Request, Response, Notification
	bodyType  byte                // set at marshal(Packet2Bytes), body compress, marshal type
	Fill      uint16              // sender set, any data
}

// MakeHeaderFromBytes unmarshal header from bytelist
func MakeHeaderFromBytes(buf []byte) Header {
	var h Header
	h.bodyLen = binary.LittleEndian.Uint32(buf[0:4])
	h.ID = binary.LittleEndian.Uint32(buf[4:8])
	h.Cmd = binary.LittleEndian.Uint16(buf[8:10])
	h.ErrorCode = c2t_error.ErrorCode(binary.LittleEndian.Uint16(buf[10:12]))
	h.FlowType = FlowType(buf[12])
	h.bodyType = buf[13]
	h.Fill = binary.LittleEndian.Uint16(buf[14:16])
	return h
}

func (h Header) ToBytesAt(buf []byte) {
	binary.LittleEndian.PutUint32(buf[0:4], h.bodyLen)
	binary.LittleEndian.PutUint32(buf[4:8], h.ID)
	binary.LittleEndian.PutUint16(buf[8:10], h.Cmd)
	binary.LittleEndian.PutUint16(buf[10:12], uint16(h.ErrorCode))
	buf[12] = byte(h.FlowType)
	buf[13] = h.bodyType
	binary.LittleEndian.PutUint16(buf[14:16], h.Fill)
}

// ToByteList marshal header to bytelist
func (h Header) ToByteList() []byte {
	buf := make([]byte, HeaderLen)
	h.ToBytesAt(buf)
	return buf
}

// GetBodyLenFromHeaderBytes return packet body len from bytelist of header
func GetBodyLenFromHeaderBytes(buf []byte) uint32 {
	return binary.LittleEndian.Uint32(buf[0:4])
}

// BodyLen return bodylen field
func (h *Header) BodyLen() uint32 {
	return h.bodyLen
}

// BodyType return bodyType field
func (h *Header) BodyType() byte {
	return h.bodyType
}

///////////////////////////////////////////////////////////////////////////////

func Bytes2HeaderBody(rdata []byte) (Header, []byte, error) {
	if len(rdata) < HeaderLen {
		return Header{}, nil, fmt.Errorf("header not complete")
	}
	header := MakeHeaderFromBytes(rdata)
	if len(rdata) != HeaderLen+int(header.bodyLen) {
		return header, nil, fmt.Errorf("packet not complete")
	}
	return header, rdata[HeaderLen : HeaderLen+int(header.bodyLen)], nil
}

func ReadHeaderBody(conn io.Reader, buffer []byte) (Header, []byte, error) {
	recvLen := 0
	toRead := HeaderLen
	for recvLen < toRead {
		n, err := conn.Read(buffer[recvLen:toRead])
		if err != nil {
			return Header{}, nil, err
		}
		recvLen += n
	}
	header := MakeHeaderFromBytes(buffer)
	toRead += int(header.bodyLen)
	for recvLen < toRead {
		n, err := conn.Read(buffer[recvLen:toRead])
		if err != nil {
			return header, nil, err
		}
		recvLen += n
	}
	return header, buffer[HeaderLen:toRead], nil
}

// Packet2Bytes make packet to bytelist
// marshalBodyFn append marshaled(+compress) body to buffer and return total buffer, bodyType, error
// set Packet.Header.bodyLen, Packet.Header.bodyType
// return bytelist, error
func Packet2Bytes(pk *Packet,
	marshalBodyFn func(interface{}, []byte) ([]byte, byte, error),
	oldbuf []byte,
) ([]byte, error) {
	newbuf, bodytype, err := marshalBodyFn(pk.Body, oldbuf)
	if err != nil {
		return nil, err
	}
	bodyLen := len(newbuf) - HeaderLen
	if bodyLen > c2t_const.MaxBodyLen {
		return nil,
			fmt.Errorf("fail to serialize large packet %v, %v", pk.Header, bodyLen)
	}
	pk.Header.bodyType = bodytype
	pk.Header.bodyLen = uint32(bodyLen)
	pk.Header.ToBytesAt(newbuf)
	return newbuf, nil
}
