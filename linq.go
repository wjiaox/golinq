package linq

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
)

type T interface{}

//if Values!=nil type of datastruct is slice else if Jval!=nil type of datastruct is struct
type Query struct {
	Values []T
	Jval   []string // try to use json to save data
	Err    error
	Kind   reflect.Kind
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
	//fmt.Print("get jval:", jval)

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
	var kind reflect.Kind
	if jin != nil {
		kind = reflect.Struct
	} else if in != nil {
		kind = reflect.Slice
	}
	return &Query{Values: in, Jval: jin, Kind: kind}
}

func (q *Query) Where(field string, f func(s T) (bool, error)) *Query {
	var vals []T
	if field == "" {
		var temp []T
		temp = append(temp, q.Values...)
		//--->for condition xx > xx;Compatible with kinds of int or float
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

//method select Query.Select("")
func (q *Query) Select(f func(s T) (T, error)) *Query {
	var vals []T
	for _, v := range q.Values {
		sel, _ := f(v)
		vals = append(vals, sel)
	}
	q.Values = vals
	return q
}

func (q *Query) SelectField(field string) {
	var vals []T
	var m interface{}
	for _, v := range q.Values {
		err := json.Unmarshal([]byte(v.(string)), &m)
		if err != nil {
			q.Err = err
			return
		} else {
			if val, ok := m.(map[string]interface{})[field]; ok {
				vals = append(vals, val)
			}
		}

	}
	q.Values = vals

}
func (q *Query) Average() float64 {
	var average float64
	var arr []T
	var l int
	arr = append(arr, q.Values...)
	l = len(arr)
	if q.Kind == reflect.Slice {
		err := datatrans(arr)
		if err == nil {
			for _, v := range arr {
				average += v.(float64)
			}
		} else {
			q.Err = err
			return 0
		}
	}
	return average / float64(l)

}

func (q *Query) AverageByField(field string) *Query {
	var vals []T
	var m interface{}
	if q.Values == nil && q.Jval != nil {
		for _, v := range q.Jval {
			vals = append(vals, v)
		}
		q.Values = vals
	}

	vals = make([]T, 0)
	if q.Kind == reflect.Struct {
		for _, v := range q.Values {
			err := json.Unmarshal([]byte(v.(string)), &m)
			if err != nil {
				q.Err = err
				return nil
			}
			//json decode the datastruct that may decode the type of all kinds of num to float64
			//so we can conver it to val.(float64) for convenient to sum
			if val, ok := m.(map[string]interface{})[field]; ok {
				switch val.(type) {
				case []interface{}:
					{
						var sum float64
						if err == nil {
							for _, vi := range val.([]interface{}) {
								sum += vi.(float64)
							}
							average := sum / float64(len(val.([]interface{})))
							m.(map[string]interface{})[field] = average
							j, _ := json.Marshal(m)

							vals = append(vals, string(j))

						} else {
							q.Err = err
							return nil
						}
						break
					}
				case int, int32, int64,
					float32, float64:
					{
						vals = append(vals, val)
						break
					}
				default:
					{
						break
					}
				}
			} else {
				q.Err = err
				return nil
			}
		}
	}

	if len(q.Values) == len(vals) && reflect.TypeOf(vals[0]).Kind() == reflect.Float64 {
		var sum float64
		l := len(vals)
		for _, v := range vals {
			sum += v.(float64)
		}
		vals = make([]T, 0)
		vals = append(vals, sum/float64(l))
	}
	q.Values = vals
	return q
}

func (q *Query) Empty() bool {
	if q.Values == nil && q.Jval == nil {
		return true
	}
	return false
}

func (q *Query) Reverse() []T {
	l := len(q.Values)
	for i := 0; i < l/2; i++ {
		q.Values[i], q.Values[l-1-i] = q.Values[l-1-i], q.Values[i]
	}
	return q.Values
}

func (q *Query) Result() interface{} {
	return q.Values
}

func init() {

}
