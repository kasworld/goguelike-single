// Code generated by "genenum.exe -typename=TowerAchieve -packagename=towerachieve -basedir=enum -vectortype=float64"

package towerachieve_vector

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"

	"github.com/kasworld/goguelike-single/enum/towerachieve"
)

type TowerAchieveVector [towerachieve.TowerAchieve_Count]float64

func (es TowerAchieveVector) String() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "TowerAchieveVector[")
	for i, v := range es {
		fmt.Fprintf(&buf,
			"%v:%v ",
			towerachieve.TowerAchieve(i), v)
	}
	buf.WriteString("]")
	return buf.String()
}
func (es *TowerAchieveVector) Dec(e towerachieve.TowerAchieve) {
	es[e] -= 1
}
func (es *TowerAchieveVector) Inc(e towerachieve.TowerAchieve) {
	es[e] += 1
}
func (es *TowerAchieveVector) Add(e towerachieve.TowerAchieve, v float64) {
	es[e] += v
}
func (es *TowerAchieveVector) SetIfGt(e towerachieve.TowerAchieve, v float64) {
	if es[e] < v {
		es[e] = v
	}
}
func (es TowerAchieveVector) Get(e towerachieve.TowerAchieve) float64 {
	return es[e]
}

// Iter return true if iter stop, return false if iter all
// fn return true to stop iter
func (es TowerAchieveVector) Iter(fn func(i towerachieve.TowerAchieve, v float64) bool) bool {
	for i, v := range es {
		if fn(towerachieve.TowerAchieve(i), v) {
			return true
		}
	}
	return false
}

// VectorAdd add element to element
func (es TowerAchieveVector) VectorAdd(arg TowerAchieveVector) TowerAchieveVector {
	var rtn TowerAchieveVector
	for i, v := range es {
		rtn[i] = v + arg[i]
	}
	return rtn
}

// VectorSub sub element to element
func (es TowerAchieveVector) VectorSub(arg TowerAchieveVector) TowerAchieveVector {
	var rtn TowerAchieveVector
	for i, v := range es {
		rtn[i] = v - arg[i]
	}
	return rtn
}

func (es *TowerAchieveVector) ToWeb(w http.ResponseWriter, r *http.Request) error {
	tplIndex, err := template.New("index").Funcs(IndexFn).Parse(`
		<html>
		<head>
		<title>TowerAchieve statistics</title>
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
	return towerachieve.TowerAchieve(i).String()
}

var IndexFn = template.FuncMap{
	"TowerAchieveIndex": Index,
}

const (
	HTML_tableheader = `<tr>
		<th>Name</th>
		<th>Value</th>
		</tr>`
	HTML_row = `<tr>
		<td>{{TowerAchieveIndex $i}}</td>
		<td>{{$v}}</td>
		</tr>
		`
)
