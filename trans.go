package linq

import (
	"errors"
	"reflect"
	"strconv"
)

func datatrans(arr []T) error {
	var err error
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
				err = errors.New("not int!")
			}
		}

	}
	return err
}

func interface2float(arr []string) ([]float64, error) {
	var iarr []float64
	for _, v := range arr {
		i, err := strconv.Atoi(v)
		if err != nil {
			f, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return nil, err
			} else {
				iarr = append(iarr, f)
			}

		} else {
			iarr = append(iarr, float64(i))
		}
	}
	return iarr, nil

}
