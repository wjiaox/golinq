# golinq
simple linq for test

#How to...
    type student2 struct {
	     Name   string
	     Grades []float64 //also can []int
    }

	ss := []student2{
		{"John", []float64{78, 84, 83, 89, 75}},
		{"Tony", []float64{60, 76.4, 83, 71, 77}},
		{"Lily", []float64{97, 84, 83, 89, 87}},
		{"Jenny", []float64{88, 80, 83, 88.5, 77}},
		{"Tom", []float64{90, 96.3, 98, 99, 89}},
	}
	var si = make([]T, len(ss))
	for i, v := range ss {
		si[i] = v
	}
	var val = From(si...).Where("Grades", func(s T) (bool, error) {
		//must be s.(float64)
		return s.(float64) > 86, nil
	}).AverageByField("Grades")
