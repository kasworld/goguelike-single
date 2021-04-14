// Code generated by "genenum.exe -typename=FactionType -packagename=factiontype -basedir=enum -vectortype=int"

package factiontype_vector

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"

	"github.com/kasworld/goguelike-single/enum/factiontype"
)

type FactionTypeVector [factiontype.FactionType_Count]int

func (es FactionTypeVector) String() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "FactionTypeVector[")
	for i, v := range es {
		fmt.Fprintf(&buf,
			"%v:%v ",
			factiontype.FactionType(i), v)
	}
	buf.WriteString("]")
	return buf.String()
}
func (es *FactionTypeVector) Dec(e factiontype.FactionType) {
	es[e] -= 1
}
func (es *FactionTypeVector) Inc(e factiontype.FactionType) {
	es[e] += 1
}
func (es *FactionTypeVector) Add(e factiontype.FactionType, v int) {
	es[e] += v
}
func (es *FactionTypeVector) SetIfGt(e factiontype.FactionType, v int) {
	if es[e] < v {
		es[e] = v
	}
}
func (es FactionTypeVector) Get(e factiontype.FactionType) int {
	return es[e]
}

// Iter return true if iter stop, return false if iter all
// fn return true to stop iter
func (es FactionTypeVector) Iter(fn func(i factiontype.FactionType, v int) bool) bool {
	for i, v := range es {
		if fn(factiontype.FactionType(i), v) {
			return true
		}
	}
	return false
}

// VectorAdd add element to element
func (es FactionTypeVector) VectorAdd(arg FactionTypeVector) FactionTypeVector {
	var rtn FactionTypeVector
	for i, v := range es {
		rtn[i] = v + arg[i]
	}
	return rtn
}

// VectorSub sub element to element
func (es FactionTypeVector) VectorSub(arg FactionTypeVector) FactionTypeVector {
	var rtn FactionTypeVector
	for i, v := range es {
		rtn[i] = v - arg[i]
	}
	return rtn
}

func (es *FactionTypeVector) ToWeb(w http.ResponseWriter, r *http.Request) error {
	tplIndex, err := template.New("index").Funcs(IndexFn).Parse(`
		<html>
		<head>
		<title>FactionType statistics</title>
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
	return factiontype.FactionType(i).String()
}

var IndexFn = template.FuncMap{
	"FactionTypeIndex": Index,
}

const (
	HTML_tableheader = `<tr>
		<th>Name</th>
		<th>Value</th>
		</tr>`
	HTML_row = `<tr>
		<td>{{FactionTypeIndex $i}}</td>
		<td>{{$v}}</td>
		</tr>
		`
)