package utils

import (
	"regexp"
	"testing"
	"unicode/utf8"

	"github.com/go-playground/validator/v10"
)

func TestGenerateRequestID(t *testing.T) {
	regex := regexp.MustCompile(`^[a-f0-9]{8}-[a-f0-9]{4}-4[a-f0-9]{3}-[89ab][a-f0-9]{3}-[a-f0-9]{12}$`)

	requestID := GenerateRequestID()
	if !regex.MatchString(requestID) {
		t.Errorf("GenerateRequestID() returned invalid UUIDv4: %s", requestID)
	}
}

func TestLuhnCheck(t *testing.T) {
	type args struct {
		nums string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"0 is invalid", args{"0"}, false},
		{"1 is invalid", args{"1"}, false},
		{"11 is invalid", args{"11"}, false},
		{"70483 is invalid", args{"70483"}, false},
		{"349926205465199 is invalid", args{"349926205465199"}, false},
		{"00 is valid", args{"00"}, true},
		{"18 is valid", args{"18"}, true},
		{"70482 is valid", args{"70482"}, true},
		{"349926205465194 is valid", args{"349926205465194"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LuhnCheck(tt.args.nums); got != tt.want {
				t.Errorf("Checksum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateRandomString(t *testing.T) {
	length := 10
	randomString := GenerateRandomString(length)

	if utf8.RuneCountInString(randomString) == 0 {
		t.Errorf("generated string is empty")
	}
}

func TestHashPassword(t *testing.T) {
	password := "testpassword"
	hashedPassword, err := HashPassword(password)

	if err != nil {
		t.Errorf("error hashing password: %v", err)
	}

	if len(hashedPassword) == 0 {
		t.Errorf("hashed password is empty")
	}
}

func TestComparePassword(t *testing.T) {
	password := "testpassword"
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Error hashing password: %v", err)
	}

	if !ComparePassword(hashedPassword, password) {
		t.Errorf("password comparison failed")
	}

	if ComparePassword(hashedPassword, "wrongpassword") {
		t.Errorf("expected comparison to fail, but it succeeded")
	}
}

type LuhnTestStruct struct {
	CardNumber string `validate:"luhn"`
}

func TestLuhnValidation(t *testing.T) {
	validate := validator.New()

	err := validate.RegisterValidation("luhn", luhnValidation)
	if err != nil {
		t.Fatalf("Failed to register luhn validation: %v", err)
	}

	validNumber := LuhnTestStruct{CardNumber: "79927398713"}
	err = validate.Struct(validNumber)
	if err != nil {
		t.Errorf("expected valid Luhn number, got error: %v", err)
	}

	invalidNumber := LuhnTestStruct{CardNumber: "12345678901"}
	err = validate.Struct(invalidNumber)
	if err == nil {
		t.Errorf("expected invalid Luhn number, but got no error")
	}
}

var luhnValidation validator.Func = func(fl validator.FieldLevel) bool {
	number, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}

	return LuhnCheck(number)
}
