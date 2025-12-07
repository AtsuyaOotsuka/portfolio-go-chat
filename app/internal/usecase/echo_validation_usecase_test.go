package usecase

import (
	"fmt"
	"testing"

	"github.com/go-playground/validator/v10"
)

func TestValidate(t *testing.T) {
	validate := validator.New()
	cv := &CustomValidator{Validator: validate}

	type TestStruct struct {
		Name  string `validate:"required"`
		Email string `validate:"required,email"`
	}

	validInput := TestStruct{
		Name:  "John Doe",
		Email: "john.doe@example.com",
	}

	invalidInput := TestStruct{
		Name:  "",
		Email: "invalid-email",
	}

	tests := []struct {
		input    TestStruct
		expected bool
	}{
		{input: validInput, expected: true},
		{input: invalidInput, expected: false},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Testing input: %+v", test.input), func(t *testing.T) {
			err := cv.Validate(test.input)
			if (err == nil) != test.expected {
				t.Errorf("Expected validation result for %+v to be %v, got %v",
					test.input, test.expected, err == nil)
			}
		})
	}
}
