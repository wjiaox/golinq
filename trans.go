package linq

import (
	"errors"
	"fmt"
	"strconv"
)

func datatrans(arr []T) error {
	var err error
	for i, v := range arr {
		switch v.(type) {
		case uint8, uint16, uint32, uint64, int, int8, int32, int64, float32, float64:
			{
				var f float64

				fmt.Sscanf(fmt.Sprintf("%v", v), "%v", &f)
				fmt.Println("test:", fmt.Sprintf("%v", v), f)
				arr[i] = f
				break
			}

		default:
			{
				arr[i] = v
				err = errors.New("not num!")
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
