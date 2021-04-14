// Code generated by "genenum.exe -typename=ScrollType -packagename=scrolltype -basedir=enum -vectortype=int"

package scrolltype_vector

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"

	"github.com/kasworld/goguelike-single/enum/scrolltype"
)

type ScrollTypeVector [scrolltype.ScrollType_Count]int

func (es ScrollTypeVector) String() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "ScrollTypeVector[")
	for i, v := range es {
		fmt.Fprintf(&buf,
			"%v:%v ",
			scrolltype.ScrollType(i), v)
	}
	buf.WriteString("]")
	return buf.String()
}
func (es *ScrollTypeVector) Dec(e scrolltype.ScrollType) {
	es[e] -= 1
}
func (es *ScrollTypeVector) Inc(e scrolltype.ScrollType) {
	es[e] += 1
}
func (es *ScrollTypeVector) Add(e scrolltype.ScrollType, v int) {
	es[e] += v
}
func (es *ScrollTypeVector) SetIfGt(e scrolltype.ScrollType, v int) {
	if es[e] < v {
		es[e] = v
	}
}
func (es ScrollTypeVector) Get(e scrolltype.ScrollType) int {
	return es[e]
}

// Iter return true if iter stop, return false if iter all
// fn return true to stop iter
func (es ScrollTypeVector) Iter(fn func(i scrolltype.ScrollType, v int) bool) bool {
	for i, v := range es {
		if fn(scrolltype.ScrollType(i), v) {
			return true
		}
	}
	return false
}

// VectorAdd add element to element
func (es ScrollTypeVector) VectorAdd(arg ScrollTypeVector) ScrollTypeVector {
	var rtn ScrollTypeVector
	for i, v := range es {
		rtn[i] = v + arg[i]
	}
	return rtn
}

// VectorSub sub element to element
func (es ScrollTypeVector) VectorSub(arg ScrollTypeVector) ScrollTypeVector {
	var rtn ScrollTypeVector
	for i, v := range es {
		rtn[i] = v - arg[i]
	}
	return rtn
}

func (es *ScrollTypeVector) ToWeb(w http.ResponseWriter, r *http.Request) error {
	tplIndex, err := template.New("index").Funcs(IndexFn).Parse(`
		<html>
		<head>
		<title>ScrollType statistics</title>
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
	return scrolltype.ScrollType(i).String()
}

var IndexFn = template.FuncMap{
	"ScrollTypeIndex": Index,
}

const (
	HTML_tableheader = `<tr>
		<th>Name</th>
		<th>Value</th>
		</tr>`
	HTML_row = `<tr>
		<td>{{ScrollTypeIndex $i}}</td>
		<td>{{$v}}</td>
		</tr>
		`
)