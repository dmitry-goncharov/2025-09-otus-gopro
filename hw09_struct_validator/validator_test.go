package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			User{
				ID:     "123456789012345678901234567890123456",
				Name:   "some name",
				Age:    20,
				Email:  "example@example.ru",
				Role:   "admin",
				Phones: []string{"79111234567"},
				meta:   nil,
			},
			nil,
		},
		{
			User{
				ID:     "12345678901234567890123456789012345",
				Name:   "some name",
				Age:    17,
				Email:  "example@example.ru",
				Role:   "user",
				Phones: []string{"79111234567"},
				meta:   nil,
			},
			ValidationErrors{
				ValidationError{
					Field: "ID",
					Err:   fmt.Errorf("%s should be length %d", "12345678901234567890123456789012345", 36),
				},
				ValidationError{
					Field: "Age",
					Err:   fmt.Errorf("%d should not be less than %d", 17, 18),
				},
				ValidationError{
					Field: "Role",
					Err:   fmt.Errorf("%s should be in %s", "user", []string{"admin", "stuff"}),
				},
			},
		},
		{
			App{
				Version: "1.0.0",
			},
			nil,
		},
		{
			App{
				Version: "1.0",
			},
			ValidationErrors{
				ValidationError{
					Field: "Version",
					Err:   fmt.Errorf("%s should be length %d", "1.0", 5),
				},
			},
		},
		{
			Token{
				Header:    []byte{0, 1, 2, 3, 4},
				Payload:   []byte{5, 6, 7, 8, 9},
				Signature: []byte{2, 3, 4, 5, 6},
			},
			nil,
		},
		{
			Response{
				Code: 200,
				Body: "some body",
			},
			nil,
		},
		{
			Response{
				Code: 201,
				Body: "some body",
			},
			ValidationErrors{
				ValidationError{
					Field: "Code",
					Err:   fmt.Errorf("%d should be in %s", 201, []string{"200", "404", "500"}),
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)

			require.Equal(t, tt.expectedErr, err)
		})
	}
}
