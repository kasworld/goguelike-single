// Code generated by "genprotocol.exe -ver=4a12840e44ca35eea470dac2384a47d9415e454793d88901bc0b6da0c240e5cd -basedir=protocol_c2t -prefix=c2t -statstype=int"

package c2t_statnoti

import (
	"fmt"
	"net/http"
	"sync"
	"text/template"

	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_idnoti"
	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_packet"
)

func (ns *StatNotification) String() string {
	return fmt.Sprintf("StatNotification[%v]", len(ns))
}

type StatNotification [c2t_idnoti.NotiID_Count]StatRow

func New() *StatNotification {
	ns := new(StatNotification)
	for i := 0; i < c2t_idnoti.NotiID_Count; i++ {
		ns[i].Name = c2t_idnoti.NotiID(i).String()
	}
	return ns
}
func (ns *StatNotification) Add(hd c2t_packet.Header) {
	if int(hd.Cmd) >= c2t_idnoti.NotiID_Count {
		return
	}
	ns[hd.Cmd].add(hd)
}
func (ns *StatNotification) ToWeb(w http.ResponseWriter, r *http.Request) error {
	tplIndex, err := template.New("index").Parse(`
	<html><head><title>Notification packet stat Info</title></head><body>
	<table border=1 style="border-collapse:collapse;">` +
		HTML_tableheader +
		`{{range $i, $v := .}}` +
		HTML_row +
		`{{end}}` +
		HTML_tableheader +
		`</table><br/>
	</body></html>`)
	if err != nil {
		return err
	}
	if err := tplIndex.Execute(w, ns); err != nil {
		return err
	}
	return nil
}

const (
	HTML_tableheader = `<tr>
	<th>Name</th>
	<th>Count</th>
	<th>Total Byte</th>
	<th>Max Byte</th>
	<th>Avg Byte</th>
	</tr>`
	HTML_row = `<tr>
	<td>{{$v.Name}}</td>
	<td>{{$v.Count }}</td>
	<td>{{$v.TotalByte }}</td>
	<td>{{$v.MaxByte }}</td>
	<td>{{printf "%10.3f" $v.Avg }}</td>
	</tr>
	`
)

type StatRow struct {
	mutex     sync.Mutex
	Name      string
	Count     int
	TotalByte int
	MaxByte   int
}

func (ps *StatRow) add(hd c2t_packet.Header) {
	ps.mutex.Lock()
	ps.Count++
	n := int(hd.BodyLen()) + c2t_packet.HeaderLen
	ps.TotalByte += n
	if n > ps.MaxByte {
		ps.MaxByte = n
	}
	ps.mutex.Unlock()
}
func (ps *StatRow) Avg() float64 {
	return float64(ps.TotalByte) / float64(ps.Count)
}
