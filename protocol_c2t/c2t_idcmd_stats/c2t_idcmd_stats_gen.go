// Code generated by "genprotocol.exe -ver=4ab5c823f1f31871cd0cef53b64b0920dee9d59673878aa337f0c3cb12b8a163 -basedir=protocol_c2t -prefix=c2t -statstype=int"

package c2t_idcmd_stats

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"

	"github.com/kasworld/goguelike-single/protocol_c2t/c2t_idcmd"
)

type CommandIDStat [c2t_idcmd.CommandID_Count]int

func (es CommandIDStat) String() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "CommandIDStats[")
	for i, v := range es {
		fmt.Fprintf(&buf,
			"%v:%v ",
			c2t_idcmd.CommandID(i), v)
	}
	buf.WriteString("]")
	return buf.String()
}
func (es *CommandIDStat) Inc(e c2t_idcmd.CommandID) {
	es[e] += 1
}
func (es *CommandIDStat) Add(e c2t_idcmd.CommandID, v int) {
	es[e] += v
}
func (es *CommandIDStat) SetIfGt(e c2t_idcmd.CommandID, v int) {
	if es[e] < v {
		es[e] = v
	}
}
func (es CommandIDStat) Get(e c2t_idcmd.CommandID) int {
	return es[e]
}

// Iter return true if iter stop, return false if iter all
// fn return true to stop iter
func (es CommandIDStat) Iter(fn func(i c2t_idcmd.CommandID, v int) bool) bool {
	for i, v := range es {
		if fn(c2t_idcmd.CommandID(i), v) {
			return true
		}
	}
	return false
}

// VectorAdd add element to element
func (es CommandIDStat) VectorAdd(arg CommandIDStat) CommandIDStat {
	var rtn CommandIDStat
	for i, v := range es {
		rtn[i] = v + arg[i]
	}
	return rtn
}

// VectorSub sub element to element
func (es CommandIDStat) VectorSub(arg CommandIDStat) CommandIDStat {
	var rtn CommandIDStat
	for i, v := range es {
		rtn[i] = v - arg[i]
	}
	return rtn
}

func (es CommandIDStat) ToWeb(w http.ResponseWriter, r *http.Request) error {
	tplIndex, err := template.New("index").Funcs(IndexFn).Parse(`
		<html>
		<head>
		<title>CommandID statistics</title>
		</head>
		<body>
		<table border=1 style="border-collapse:collapse;">` +
		HTML_tableheader +
		`{{range $i, $v := .}}` +
		HTML_row +
		`{{end}}` +
		HTML_tableheader +
		`</table>
	
		<br/>
		</body>
		</html>
		`)
	if err != nil {
		return err
	}
	if err := tplIndex.Execute(w, es); err != nil {
		return err
	}
	return nil
}

func Index(i int) string {
	return c2t_idcmd.CommandID(i).String()
}

var IndexFn = template.FuncMap{
	"CommandIDIndex": Index,
}

const (
	HTML_tableheader = `<tr>
		<th>Name</th>
		<th>Value</th>
		</tr>`
	HTML_row = `<tr>
		<td>{{CommandIDIndex $i}}</td>
		<td>{{$v}}</td>
		</tr>
		`
)
