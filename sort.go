package linq

import (
	"fmt"
	"sort"
)
import "reflect"
import "encoding/json"
import "strings"

func quicksort(a []T, low int, high int) {
	if low >= high {
		return
	}
	first := low
	last := high
	key := a[first].(float64)
	for first < last {
		for key <= a[last].(float64) && first < last {
			last -= 1
		}
		a[first] = a[last]
		for key >= a[first].(float64) && first < last {
			first += 1
		}
		a[last] = a[first]

		a[first] = key
		quicksort(a[:], low, first-1)
		quicksort(a[:], first+1, high)
	}

}
func descquicksort(a []T, low int, high int) {
	if low >= high {
		return
	}
	first := low
	last := high
	key := a[first].(float64)
	for first < last {
		for key >= a[last].(float64) && first < last {
			last -= 1
		}
		a[first] = a[last]
		for key <= a[first].(float64) && first < last {
			first += 1
		}
		a[last] = a[first]

		a[first] = key
		descquicksort(a[:], low, first-1)
		descquicksort(a[:], first+1, high)
	}

}

func datatrans(arr []T) {
	for i, v := range arr {
		switch reflect.TypeOf(v).Kind() {
		case reflect.Int:
			{
				arr[i] = float64(v.(int))
				break
			}
		case reflect.Int64:
			{
				arr[i] = float64(v.(int64))
				break
			}
		case reflect.Float32:
			{
				arr[i] = float64(v.(float32))
				break
			}
		case reflect.Float64:
			{
				arr[i] = float64(v.(float64))
				break
			}
		default:
			{
				arr[i] = v
			}
		}

	}
}
func maxs(a, b int) int {
	return b&((a-b)>>31) | a&^((a-b)>>31)
}
func minis(a, b int) int {
	return a&((a-b)>>31) | b&^((a-b)>>31)
}
func stringcompare(a, b string) bool {
	var c = false
	if a != "" && b != "" {
		a = strings.ToLower(a)
		b = strings.ToLower(b)
		l := minis(len(a), len(b))
		for i := 0; i < l; i++ {
			if int(a[i]) > int(b[i]) {
				c = true
				break
			} else if int(a[i]) < int(b[i]) {
				c = false
				break
			} else if a[i] == b[i] {
				continue
			}

		}
	} else {
		c = true
	}
	return c

}
func unitiysort(s4c, s4r []string, val interface{}, v T) ([]string, []string) {
	if s4c[0] == "" {
		s4c[0] = val.(string)
		s4r[0] = v.(string)
		return s4c, s4r
	}
	for j := 0; j < len(s4c); j++ {
		if !stringcompare(val.(string), s4c[0]) {
			t := []string{val.(string)}
			t = append(t, s4c...)
			s4c = t
			tr := []string{v.(string)}
			tr = append(tr, s4r...)
			s4r = tr
			break
		} else if stringcompare(val.(string), s4c[len(s4c)-1]) {
			s4c = append(s4c, val.(string))
			s4r = append(s4r, v.(string))
			break
		} else if stringcompare(val.(string), s4c[j]) && !stringcompare(val.(string), s4c[j+1]) {
			var t, t1, t2 []string
			t1 = append(t1, s4c[:j+1]...) //array begin from 0
			t2 = append(t2, s4c[j+1:]...)
			t = append(t, t1...)
			t = append(t, val.(string))
			t = append(t, t2...)
			s4c = t

			var tr, tr1, tr2 []string
			tr1 = append(tr1, s4r[:j+1]...)
			tr2 = append(tr2, s4r[j+1:]...)
			tr = append(tr, tr1...)
			tr = append(tr, v.(string))
			tr = append(tr, tr2...)
			s4r = tr
			break
		} else {
			continue
		}
	}
	return s4c, s4r

}

func ascsort(field string, arr []T) []T {
	t4a := reflect.TypeOf(arr[0]).Kind()
	if field == "" && (t4a == reflect.Int || t4a == reflect.Float64 || t4a == reflect.Float32) {
		datatrans(arr)
		quicksort(arr, 0, len(arr)-1)
	} else if field == "" && t4a == reflect.String {
		t := make([]string, len(arr))
		for i, v := range arr {
			t[i] = v.(string)
		}
		sort.Strings(t)
		for i, v := range t {
			arr[i] = v
		}
	} else if field != "" && t4a == reflect.String {
		var m interface{}
		var s4c = make([]string, 1)
		var s4r = make([]string, 1)

		for _, v := range arr {
			err := json.Unmarshal([]byte(v.(string)), &m)
			if err != nil {
				return nil
			}
			if val, ok := m.(map[string]interface{})[field]; ok {
				var t4v reflect.Kind
				if val != nil {
					t4v = reflect.TypeOf(val).Kind()
				} else {
					continue
				}
				switch t4v {
				case reflect.Slice: //ok
					{
						s := fmt.Sprintf("%v", val)

						if s != "" && s != "null" {
							ss := strings.Split(s[1:len(s)-1], " ")
							sort.Strings(ss)

							m.(map[string]interface{})[field] = ss
							var jv T
							j, _ := json.Marshal(m)
							jv = string(j)
							if ss[0] != "" {
								s4c, s4r = unitiysort(s4c, s4r, ss[0], jv)
							}
						}
						break

					}
				case reflect.String:
					{
						s4c, s4r = unitiysort(s4c, s4r, val.(string), v)
						break

					}
				case reflect.Int,
					reflect.Int64,
					reflect.Float32,
					reflect.Float64:
					{
						s := fmt.Sprintf("%v", val)
						s4c, s4r = unitiysort(s4c, s4r, s, v)
						break

					}

				}
			}
		}
		arr = make([]T, len(s4r))
		for i := 0; i < len(s4r); i++ {
			arr[i] = s4r[i]
		}
		fmt.Println("arr:", arr)

	}
	return arr
}
func descsort(field string, arr []T) []T {
	arr = ascsort(field, arr)
	l := len(arr)
	for i := 0; i < l/2; i++ {
		arr[i], arr[l-i-1] = arr[l-i-1], arr[i]
	}
	return arr
}
