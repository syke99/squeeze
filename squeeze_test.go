package squeeze

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type Person struct {
	Name string `json:"Name"`
	Age  int    `json:"Age"`
	Body Body   `json:"Body"`
}

type Body struct {
	Weight int `json:"Weight"`
	Height int `json:"Height"`
}

var Kevin = Person{
	Name: "Kevin",
	Age:  33,
	Body: Body{
		Weight: 210,
		Height: 6,
	},
}

var Dave = Person{
	Name: "Dave",
	Age:  42,
	Body: Body{
		Weight: 198,
		Height: 7,
	},
}

var Mary = Person{
	Name: "Mary",
	Age:  20,
	Body: Body{
		Weight: 155,
		Height: 5,
	},
}

func Name(left, right Person) Result {
	if left.Name == "Kevin" {
		return Left
	} else if right.Name == "Dave" {
		return Right
	}

	return Left
}

func Age(left, right Person) Result {
	if left.Age <= 35 {
		return Left
	}

	if right.Age > 35 && right.Age < 50 {
		return Right
	}

	return Left
}

func Weight(left, right Person) Result {
	if left.Body.Weight < right.Body.Weight {
		return Left
	} else if left.Body.Weight > right.Body.Weight {
		return Right
	}

	return Left
}

func Height(left, right Person) Result {
	if left.Body.Height == 10 {
		return Left
	} else if right.Body.Height == 10 {
		return Right
	}

	return Left
}

func TestSqueeze(t *testing.T) {
	rules := Rules[Person]{
		"Name":   Name,
		"Age":    Age,
		"Weight": Weight,
		"Height": Height,
	}

	res, err := Squeeze([]Person{Kevin, Dave, Mary}, rules)
	assert.NoError(t, err)

	assert.Equal(t, "Kevin", res.Name)
	assert.Equal(t, 33, res.Age)
	assert.Equal(t, 6, res.Body.Height)
	assert.Equal(t, 155, res.Body.Weight)
}

func TestSqueezeTNotStruct(t *testing.T) {
	_, err := Squeeze([]string{""}, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "expected T to be a struct type")
}
