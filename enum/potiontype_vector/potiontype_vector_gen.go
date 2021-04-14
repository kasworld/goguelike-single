// Code generated by "genenum.exe -typename=PotionType -packagename=potiontype -basedir=enum -vectortype=int"

package potiontype_vector

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"

	"github.com/kasworld/goguelike-single/enum/potiontype"
)

type PotionTypeVector [potiontype.PotionType_Count]int

func (es PotionTypeVector) String() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "PotionTypeVector[")
	for i, v := range es {
		fmt.Fprintf(&buf,
			"%v:%v ",
			potiontype.PotionType(i), v)
	}
	buf.WriteString("]")
	return buf.String()
}
func (es *PotionTypeVector) Dec(e potiontype.PotionType) {
	es[e] -= 1
}
func (es *PotionTypeVector) Inc(e potiontype.PotionType) {
	es[e] += 1
}
func (es *PotionTypeVector) Add(e potiontype.PotionType, v int) {
	es[e] += v
}
func (es *PotionTypeVector) SetIfGt(e potiontype.PotionType, v int) {
	if es[e] < v {
		es[e] = v
	}
}
func (es PotionTypeVector) Get(e potiontype.PotionType) int {
	return es[e]
}

// Iter return true if iter stop, return false if iter all
// fn return true to stop iter
func (es PotionTypeVector) Iter(fn func(i potiontype.PotionType, v int) bool) bool {
	for i, v := range es {
		if fn(potiontype.PotionType(i), v) {
			return true
		}
	}
	return false
}

// VectorAdd add element to element
func (es PotionTypeVector) VectorAdd(arg PotionTypeVector) PotionTypeVector {
	var rtn PotionTypeVector
	for i, v := range es {
		rtn[i] = v + arg[i]
	}
	return rtn
}

// VectorSub sub element to element
func (es PotionTypeVector) VectorSub(arg PotionTypeVector) PotionTypeVector {
	var rtn PotionTypeVector
	for i, v := range es {
		rtn[i] = v - arg[i]
	}
	return rtn
}

func (es *PotionTypeVector) ToWeb(w http.ResponseWriter, r *http.Request) error {
	tplIndex, err := template.New("index").Funcs(IndexFn).Parse(`
		<html>
		<head>
		<title>PotionType statistics</title>
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
	return potiontype.PotionType(i).String()
}

var IndexFn = template.FuncMap{
	"PotionTypeIndex": Index,
}

const (
	HTML_tableheader = `<tr>
		<th>Name</th>
		<th>Value</th>
		</tr>`
	HTML_row = `<tr>
		<td>{{PotionTypeIndex $i}}</td>
		<td>{{$v}}</td>
		</tr>
		`
)
