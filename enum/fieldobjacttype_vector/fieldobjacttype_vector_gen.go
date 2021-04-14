// Code generated by "genenum.exe -typename=FieldObjActType -packagename=fieldobjacttype -basedir=enum -vectortype=int"

package fieldobjacttype_vector

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"

	"github.com/kasworld/goguelike-single/enum/fieldobjacttype"
)

type FieldObjActTypeVector [fieldobjacttype.FieldObjActType_Count]int

func (es FieldObjActTypeVector) String() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "FieldObjActTypeVector[")
	for i, v := range es {
		fmt.Fprintf(&buf,
			"%v:%v ",
			fieldobjacttype.FieldObjActType(i), v)
	}
	buf.WriteString("]")
	return buf.String()
}
func (es *FieldObjActTypeVector) Dec(e fieldobjacttype.FieldObjActType) {
	es[e] -= 1
}
func (es *FieldObjActTypeVector) Inc(e fieldobjacttype.FieldObjActType) {
	es[e] += 1
}
func (es *FieldObjActTypeVector) Add(e fieldobjacttype.FieldObjActType, v int) {
	es[e] += v
}
func (es *FieldObjActTypeVector) SetIfGt(e fieldobjacttype.FieldObjActType, v int) {
	if es[e] < v {
		es[e] = v
	}
}
func (es FieldObjActTypeVector) Get(e fieldobjacttype.FieldObjActType) int {
	return es[e]
}

// Iter return true if iter stop, return false if iter all
// fn return true to stop iter
func (es FieldObjActTypeVector) Iter(fn func(i fieldobjacttype.FieldObjActType, v int) bool) bool {
	for i, v := range es {
		if fn(fieldobjacttype.FieldObjActType(i), v) {
			return true
		}
	}
	return false
}

// VectorAdd add element to element
func (es FieldObjActTypeVector) VectorAdd(arg FieldObjActTypeVector) FieldObjActTypeVector {
	var rtn FieldObjActTypeVector
	for i, v := range es {
		rtn[i] = v + arg[i]
	}
	return rtn
}

// VectorSub sub element to element
func (es FieldObjActTypeVector) VectorSub(arg FieldObjActTypeVector) FieldObjActTypeVector {
	var rtn FieldObjActTypeVector
	for i, v := range es {
		rtn[i] = v - arg[i]
	}
	return rtn
}

func (es *FieldObjActTypeVector) ToWeb(w http.ResponseWriter, r *http.Request) error {
	tplIndex, err := template.New("index").Funcs(IndexFn).Parse(`
		<html>
		<head>
		<title>FieldObjActType statistics</title>
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
	return fieldobjacttype.FieldObjActType(i).String()
}

var IndexFn = template.FuncMap{
	"FieldObjActTypeIndex": Index,
}

const (
	HTML_tableheader = `<tr>
		<th>Name</th>
		<th>Value</th>
		</tr>`
	HTML_row = `<tr>
		<td>{{FieldObjActTypeIndex $i}}</td>
		<td>{{$v}}</td>
		</tr>
		`
)