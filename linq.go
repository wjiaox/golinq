package linq

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type T interface{}
type Query struct {
	Values []T
	Jval   []string // try to use json to save data
	Err    error
}

func getarray(t ...T) ([]string, []T, error) {
	var arr []T
	var jval []string

	for _, v := range t {
		v4v := reflect.ValueOf(v)
		t4v := reflect.TypeOf(v)
		temp := make(map[string]interface{})
		if t4v.Kind() == reflect.Struct {
			for j := 0; j < t4v.NumField(); j++ {
				k4v := v4v.Field(j).Type().Kind()
				if k4v == reflect.String {
					temp[t4v.Field(j).Name] = v4v.Field(j).String()
				} else if k4v == reflect.Int {
					temp[t4v.Field(j).Name] = v4v.Field(j).Int()
				} else if k4v == reflect.Slice {
					temp[t4v.Field(j).Name] = v4v.Field(j).Interface()
				}

			}
			j, err := json.Marshal(temp)
			if err != nil {
				return nil, nil, err
			}
			jval = append(jval, string(j))
		} else {
			arr = append(arr, v)
		}

	}
	fmt.Print("get jval:", jval)

	return jval, arr, nil
}
func From(t ...T) *Query {
	if t == nil {
		return nil
	}
	jin, in, err := getarray(t...)
	if err != nil {
		return &Query{Err: err}
	}

	return &Query{Values: in, Jval: jin}

}

func (q *Query) Where(field string, f func(s T) (bool, error)) *Query {
	var vals []T
	if field == "" {
		var temp []T
		temp = append(temp, q.Values...)
		datatrans(temp)
		for _, v := range temp {
			if ok, _ := f(v); ok {
				vals = append(vals, v)
			}
		}
	} else {
		for _, v := range q.Jval {
			var m interface{}
			err := json.Unmarshal([]byte(v), &m)
			if err == nil {
				jf := m.(map[string]interface{})[field]
				if reflect.TypeOf(jf).Kind() == reflect.Slice {
					var ms []interface{}
					for _, jv := range jf.([]interface{}) {
						if ok, _ := f(jv); ok {
							ms = append(ms, jv)
						}
					}
					m.(map[string]interface{})[field] = ms
					jbyte, _ := json.Marshal(m)
					vals = append(vals, string(jbyte))

				} else {
					if ok, _ := f(jf); ok {
						vals = append(vals, v)
					}

				}

			} else {
				q.Err = err
			}

		}
	}
	q.Values = vals

	return q
}

//sort  method 1.ASC 2.DESC
func (q *Query) OrderBy(field, method string) *Query {

	switch strings.ToUpper(method) {
	case "ASC":
		{
			q.Values = ascsort(field, q.Values)
			break
		}
	case "DESC":
		{
			q.Values = descsort(field, q.Values)
			break
		}
	default:
		{
			q.Err = errors.New("can not match method!")
			break
		}
	}
	return q
}

//group
func (q *Query) GroupBy() *Query {
	return q
}

func (q *Query) Select(f func(s T) (T, error)) *Query {
	var vals []T
	for _, v := range q.Values {
		sel, _ := f(v)
		vals = append(vals, sel)
	}
	q.Values = vals
	return q

}

func (q *Query) Result() interface{} {
	return q.Values
}

func init() {

}
