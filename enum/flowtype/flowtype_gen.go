// Code generated by "genenum.exe -typename=FlowType -packagename=flowtype -basedir=enum"

package flowtype

import "fmt"

type FlowType uint8

const (
	Request      FlowType = iota // Request for request packet (response packet expected)
	Response                     // Response is reply of request packet
	Notification                 // Notification is just send and forget packet
	//

	FlowType_Count int = iota
)

var _FlowType2string = [FlowType_Count][2]string{
	Request:      {"Request", "Request for request packet (response packet expected)"},
	Response:     {"Response", "Response is reply of request packet"},
	Notification: {"Notification", "Notification is just send and forget packet"},
}

func (e FlowType) String() string {
	if e >= 0 && e < FlowType(FlowType_Count) {
		return _FlowType2string[e][0]
	}
	return fmt.Sprintf("FlowType%d", uint8(e))
}

func (e FlowType) CommentString() string {
	if e >= 0 && e < FlowType(FlowType_Count) {
		return _FlowType2string[e][1]
	}
	return ""
}

var _string2FlowType = map[string]FlowType{
	"Request":      Request,
	"Response":     Response,
	"Notification": Notification,
}

func String2FlowType(s string) (FlowType, bool) {
	v, b := _string2FlowType[s]
	return v, b
}
