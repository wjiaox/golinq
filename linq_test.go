package linq

import (
	"encoding/json"
	"reflect"
	"testing"
)

func Test_linq(t *testing.T) {
	a := []T{1, 2, 9.01, 32, 3, 4.2, 4, 5, 6, 7, 8}
	var val = From(a...).Where("", func(s T) (bool, error) {
		return s.(float64) > 4, nil
	}).OrderBy("", "asc")

	t.Log("test:", val)
	a = []T{1, 2, 9.3, 32, 3, 4, 5, 6, 4.2, 8}
	val = From(a...).OrderBy("", "asc")

	t.Log("test:", val)
	a = []T{"Tom", "Alice", "Jenny", "John", "JoJo", "Tonny", "Allen", "Ocean"}
	val = From(a...).OrderBy("", "asc")

	t.Log("test:", val)
}

type student struct {
	name  string
	grade int
}

func Test_Select(t *testing.T) {

	var ss = []student{
		{"John", 80},
		{"Tony", 90},
		{"Lily", 88},
		{"Jenny", 87},
		{"Jenny", 98},
		{"Tom", 78},
	}
	var si = make([]T, len(ss))
	for i, v := range ss {
		si[i] = v
	}
	var val = From(si...).Where("grade", func(s T) (bool, error) {
		return s.(float64) > 86, nil
	})
	//	var val = From(ss).Where("", func(s T) (bool, error) {
	//		return s.(student).grade > 86, nil
	//	}).Select(func(s T) (T, error) {
	//		return s.(student).name, nil
	//	}).Result()
	t.Log("test jval:", val.Jval)
	t.Log("test values:", val.Values)

	val = From(si...).Where("grade", func(s T) (bool, error) {
		return s.(float64) > 86, nil
	}).OrderBy("name", "asc")
	t.Log("test jval:", val.Jval)
	t.Log("test values:", val.Values)

	val = From(si...).Where("grade", func(s T) (bool, error) {
		return s.(float64) > 86, nil
	}).OrderBy("grade", "asc")
	t.Log("test jval:", val.Jval)
	t.Log("test values:", val.Values)

}

type student2 struct {
	Name   string
	Grades []int
}

func Test_SelectMany(t *testing.T) {
	ss := []student2{
		{"John", []int{78, 84, 83, 89, 75}},
		{"Tony", []int{60, 76, 83, 71, 77}},
		{"Lily", []int{97, 84, 83, 89, 87}},
		{"Jenny", []int{88, 80, 83, 88, 77}},
		{"Tom", []int{90, 96, 98, 99, 89}},
	}
	var si = make([]T, len(ss))
	for i, v := range ss {
		si[i] = v
	}
	var val = From(si...).Where("Grades", func(s T) (bool, error) {
		return s.(float64) > 86, nil
	})

	//var val = From(ss)
	t.Log("test jval:", val.Jval)
	t.Log("test values:", val.Values)
	val = From(si...).Where("Name", func(s T) (bool, error) {
		return s.(string) == "Tony", nil
	})

	t.Log("test jval:", val.Jval)
	t.Log("test values:", val.Values)

	val = From(si...).Where("Grades", func(s T) (bool, error) {
		return s.(float64) > 86, nil
	}).OrderBy("Grades", "asc")
	t.Log("test jval:", val.Jval)
	t.Log("test values:", val.Values)

}

func Test_mytest(t *testing.T) {
	var val = student2{
		"John", []int{78, 84, 83, 89, 75},
	}
	if reflect.TypeOf(val).Kind() == reflect.Struct {
		t.Log("ok")
	}
	json.Marshal(val)
	v := reflect.ValueOf(val)

	v2 := v.FieldByName("Grades")

	t.Log(v2, v.Field(1), reflect.TypeOf(v).Field(0).Name)
}
