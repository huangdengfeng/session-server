package validator

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStructNil(t *testing.T) {
	type User struct {
		FirstName string `json:"fname"`
		LastName  string `json:"lname"`
		Age       uint8  `validate:"gte=0,lte=10"`
		Email     string `json:"e-mail" validate:"required,email"`
	}
	var u *User
	err := Validator.Struct(u)
	assert.Error(t, err)
	fmt.Println(err)
}
func TestStruct(t *testing.T) {
	type User struct {
		FirstName string `json:"fname"`
		LastName  string `json:"lname"`
		Age       uint8  `validate:"gte=0,lte=10"`
		Email     string `json:"e-mail" validate:"required,email"`
	}
	user := User{
		Age:   100,
		Email: "1",
	}
	err := Validator.Struct(user)
	assert.Error(t, err)
	fmt.Println(err)
}

func TestVarNit(t *testing.T) {
	var s *string
	err := Validator.Var(s, "gte=5,lte=10")
	assert.Error(t, err)
	fmt.Println(err)
}

func TestVar(t *testing.T) {
	s := "aaa"
	err := Validator.Var(s, "gte=5,lte=10")
	assert.Error(t, err)
	fmt.Println(err)
}
