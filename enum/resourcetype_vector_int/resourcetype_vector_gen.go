// Code generated by "genenum.exe -typename=ResourceType -packagename=resourcetype -basedir=enum -vectortype=int"

package resourcetype_vector_int

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"

	"github.com/kasworld/goguelike-single/enum/resourcetype"
)

type ResourceTypeVector_int [resourcetype.ResourceType_Count]int

func (es ResourceTypeVector_int) String() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "ResourceTypeVector_int[")
	for i, v := range es {
		fmt.Fprintf(&buf,
			"%v:%v ",
			resourcetype.ResourceType(i), v)
	}
	buf.WriteString("]")
	return buf.String()
}
func (es *ResourceTypeVector_int) Dec(e resourcetype.ResourceType) {
	es[e] -= 1
}
func (es *ResourceTypeVector_int) Inc(e resourcetype.ResourceType) {
	es[e] += 1
}
func (es *ResourceTypeVector_int) Add(e resourcetype.ResourceType, v int) {
	es[e] += v
}
func (es *ResourceTypeVector_int) SetIfGt(e resourcetype.ResourceType, v int) {
	if es[e] < v {
		es[e] = v
	}
}
func (es ResourceTypeVector_int) Get(e resourcetype.ResourceType) int {
	return es[e]
}

// Iter return true if iter stop, return false if iter all
// fn return true to stop iter
func (es ResourceTypeVector_int) Iter(fn func(i resourcetype.ResourceType, v int) bool) bool {
	for i, v := range es {
		if fn(resourcetype.ResourceType(i), v) {
			return true
		}
	}
	return false
}

// VectorAdd add element to element
func (es ResourceTypeVector_int) VectorAdd(arg ResourceTypeVector_int) ResourceTypeVector_int {
	var rtn ResourceTypeVector_int
	for i, v := range es {
		rtn[i] = v + arg[i]
	}
	return rtn
}

// VectorSub sub element to element
func (es ResourceTypeVector_int) VectorSub(arg ResourceTypeVector_int) ResourceTypeVector_int {
	var rtn ResourceTypeVector_int
	for i, v := range es {
		rtn[i] = v - arg[i]
	}
	return rtn
}

func (es *ResourceTypeVector_int) ToWeb(w http.ResponseWriter, r *http.Request) error {
	tplIndex, err := template.New("index").Funcs(IndexFn).Parse(`
		<html>
		<head>
		<title>ResourceType Vector int</title>
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
	return resourcetype.ResourceType(i).String()
}

var IndexFn = template.FuncMap{
	"ResourceTypeIndex": Index,
}

const (
	HTML_tableheader = `<tr>
		<th>Name</th>
		<th>Value</th>
		</tr>`
	HTML_row = `<tr>
		<td>{{ResourceTypeIndex $i}}</td>
		<td>{{$v}}</td>
		</tr>
		`
)