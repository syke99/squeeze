package squeeze

import (
	"encoding/json"
	"strings"

	"github.com/doublerebel/bellows"
)

// Result signals to Squeeze which struct's field to
// set the value from
type Result int

const (
	// Left signals to Squeeze that it should set the
	// resulting T's field with the corresponding field
	// from the left struct passed to a Rule
	Left Result = iota
	// Right signals to Squeeze that it should set the
	// resulting T's field with the corresponding field
	// from the right struct passed to a Rule
	Right
)

func marshalIntermediary[T any](field string, parent, intermediary map[string]interface{}, dst *T) error {
	newField := parent[field]

	intermediary[field] = newField

	intermediary = bellows.Expand(intermediary)

	intermediaryBytes, err := json.Marshal(intermediary)
	if err != nil {
		return err
	}

	err = json.Unmarshal(intermediaryBytes, dst)
	if err != nil {
		return err
	}

	return nil
}

// Rule is a generic function type used to define rules for
// which T to use a field from
type Rule[T any] func(left, right T) Result

// Rules maps a Rule to a given T's field. It accomplishes
// this by matching defined json struct tags. These two strings
// MUST equal
type Rules[T any] map[string]Rule[T]

func squeeze[T any](left, right T, rules Rules[T]) (*T, error) {
	newT := new(T)

	leftMap := bellows.Flatten(left)
	rightMap := bellows.Flatten(right)

	intermediaryLeftMap := bellows.Flatten(newT)
	intermediaryRightMap := bellows.Flatten(newT)

	intermediaryLeft := *new(T)
	intermediaryRight := *new(T)

	newTMap := bellows.Flatten(newT)

	for name, field := range newTMap {
		switch field.(type) {
		// if the field is a nested struct type, then skip the field
		case map[string]interface{}:
			continue
		default:
			err := marshalIntermediary(name, leftMap, intermediaryLeftMap, &intermediaryLeft)
			if err != nil {
				return nil, err
			}

			err = marshalIntermediary(name, rightMap, intermediaryRightMap, &intermediaryRight)
			if err != nil {
				return nil, err
			}

			splitName := strings.Split(name, ".")

			var rule Rule[T]

			switch len(splitName) {
			case 0:
				continue
			case 1:
				rule = rules[splitName[0]]
			default:
				rule = rules[splitName[len(splitName)-1]]
			}

			res := rule(intermediaryLeft, intermediaryRight)

			switch res {
			case Left:
				f := bellows.Flatten(intermediaryLeft)[name]
				newTMap[name] = f
			case Right:
				f := bellows.Flatten(intermediaryRight)[name]
				newTMap[name] = f
			}
		}
	}

	expandedNewTMap := bellows.Expand(newTMap)

	newTBytes, err := json.Marshal(expandedNewTMap)
	if err != nil {
		return newT, err
	}

	err = json.Unmarshal(newTBytes, newT)
	if err != nil {
		return newT, err
	}

	return newT, nil
}

// Squeeze loops through a slice of T (structs), comparing
// the current and previous T by applying the given rules
// to them and returns the resulting, flattened T. If T's
// type is  not a struct, it will immediately return an
// error. If structs has a length of zero, it will return
// a T with all fields zet to their zero-values. If structs
// has a length of one, it will immediately return the T
// stored at index 0
func Squeeze[T any](structs []T, rules Rules[T]) (T, error) {
	res, err := validateStruct(*new(T))

	if len(structs) == 0 {
		return *new(T), nil
	} else if len(structs) == 1 {
		return structs[0], nil
	} else {
		res = structs[0]

		for i, t := range structs {
			if i == 0 {
				continue
			}

			left := res
			right := t

			r, err := squeeze(left, right, rules)
			if err != nil {
				if r == nil {
					return *new(T), err
				} else {
					return res, err
				}
			}

			res = *r
		}
	}

	return res, err
}
