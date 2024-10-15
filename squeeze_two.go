package squeeze

import (
	"errors"
	"iter"
	"reflect"
)

func validateStruct[T any](v T) (T, error) {
	if reflect.ValueOf(v).Type().Kind() != reflect.Struct {
		return *new(T), errors.New("expected T to be a struct type")
	}
	return v, nil
}

// SqueezeIter works just like Squeeze, but instead of taking in a slice of structs,
// it takes in an iter.Seq that yields the type func () (T, error)
func SqueezeIter[T any](structs iter.Seq[func() (T, error)], rules Rules[T]) (T, error) {
	res, err := validateStruct(*new(T))
	if err != nil {
		return res, err
	}

	i := 0

	for str := range structs {
		generated, er := str()
		if er != nil {
			return res, er
		}

		if i == 0 {
			res = generated
			i++
			continue
		}

		squeezed, er := squeeze(res, generated, rules)
		if er != nil {
			return res, er
		}

		res = *squeezed
	}

	return res, nil
}
