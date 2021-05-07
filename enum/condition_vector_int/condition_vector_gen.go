// Code generated by "genenum.exe -typename=Condition -packagename=condition -basedir=enum -flagtype=uint16 -vectortype=int"

package condition_vector_int

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"

	"github.com/kasworld/goguelike-single/enum/condition"
)

type ConditionVector_int [condition.Condition_Count]int

func (es ConditionVector_int) String() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "ConditionVector_int[")
	for i, v := range es {
		fmt.Fprintf(&buf,
			"%v:%v ",
			condition.Condition(i), v)
	}
	buf.WriteString("]")
	return buf.String()
}
func (es *ConditionVector_int) Dec(e condition.Condition) {
	es[e] -= 1
}
func (es *ConditionVector_int) Inc(e condition.Condition) {
	es[e] += 1
}
func (es *ConditionVector_int) Add(e condition.Condition, v int) {
	es[e] += v
}
func (es *ConditionVector_int) SetIfGt(e condition.Condition, v int) {
	if es[e] < v {
		es[e] = v
	}
}
func (es ConditionVector_int) Get(e condition.Condition) int {
	return es[e]
}

// Iter return true if iter stop, return false if iter all
// fn return true to stop iter
func (es ConditionVector_int) Iter(fn func(i condition.Condition, v int) bool) bool {
	for i, v := range es {
		if fn(condition.Condition(i), v) {
			return true
		}
	}
	return false
}

// VectorAdd add element to element
func (es ConditionVector_int) VectorAdd(arg ConditionVector_int) ConditionVector_int {
	var rtn ConditionVector_int
	for i, v := range es {
		rtn[i] = v + arg[i]
	}
	return rtn
}

// VectorSub sub element to element
func (es ConditionVector_int) VectorSub(arg ConditionVector_int) ConditionVector_int {
	var rtn ConditionVector_int
	for i, v := range es {
		rtn[i] = v - arg[i]
	}
	return rtn
}

func (es *ConditionVector_int) ToWeb(w http.ResponseWriter, r *http.Request) error {
	tplIndex, err := template.New("index").Funcs(IndexFn).Parse(`
		<html>
		<head>
		<title>Condition Vector int</title>
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
	return condition.Condition(i).String()
}

var IndexFn = template.FuncMap{
	"ConditionIndex": Index,
}

const (
	HTML_tableheader = `<tr>
		<th>Name</th>
		<th>Value</th>
		</tr>`
	HTML_row = `<tr>
		<td>{{ConditionIndex $i}}</td>
		<td>{{$v}}</td>
		</tr>
		`
)