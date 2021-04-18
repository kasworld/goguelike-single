// Code generated by "genenum.exe -typename=TurnAction -packagename=turnaction -basedir=enum -vectortype=int"

package turnaction_vector

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"

	"github.com/kasworld/goguelike-single/enum/turnaction"
)

type TurnActionVector [turnaction.TurnAction_Count]int

func (es TurnActionVector) String() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "TurnActionVector[")
	for i, v := range es {
		fmt.Fprintf(&buf,
			"%v:%v ",
			turnaction.TurnAction(i), v)
	}
	buf.WriteString("]")
	return buf.String()
}
func (es *TurnActionVector) Dec(e turnaction.TurnAction) {
	es[e] -= 1
}
func (es *TurnActionVector) Inc(e turnaction.TurnAction) {
	es[e] += 1
}
func (es *TurnActionVector) Add(e turnaction.TurnAction, v int) {
	es[e] += v
}
func (es *TurnActionVector) SetIfGt(e turnaction.TurnAction, v int) {
	if es[e] < v {
		es[e] = v
	}
}
func (es TurnActionVector) Get(e turnaction.TurnAction) int {
	return es[e]
}

// Iter return true if iter stop, return false if iter all
// fn return true to stop iter
func (es TurnActionVector) Iter(fn func(i turnaction.TurnAction, v int) bool) bool {
	for i, v := range es {
		if fn(turnaction.TurnAction(i), v) {
			return true
		}
	}
	return false
}

// VectorAdd add element to element
func (es TurnActionVector) VectorAdd(arg TurnActionVector) TurnActionVector {
	var rtn TurnActionVector
	for i, v := range es {
		rtn[i] = v + arg[i]
	}
	return rtn
}

// VectorSub sub element to element
func (es TurnActionVector) VectorSub(arg TurnActionVector) TurnActionVector {
	var rtn TurnActionVector
	for i, v := range es {
		rtn[i] = v - arg[i]
	}
	return rtn
}

func (es *TurnActionVector) ToWeb(w http.ResponseWriter, r *http.Request) error {
	tplIndex, err := template.New("index").Funcs(IndexFn).Parse(`
		<html>
		<head>
		<title>TurnAction statistics</title>
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
	return turnaction.TurnAction(i).String()
}

var IndexFn = template.FuncMap{
	"TurnActionIndex": Index,
}

const (
	HTML_tableheader = `<tr>
		<th>Name</th>
		<th>Value</th>
		</tr>`
	HTML_row = `<tr>
		<td>{{TurnActionIndex $i}}</td>
		<td>{{$v}}</td>
		</tr>
		`
)
