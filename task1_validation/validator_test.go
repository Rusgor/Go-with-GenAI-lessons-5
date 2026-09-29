package validation

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestValidateRegistration_AllValid(t *testing.T) {
	form := RegistrationForm{Email: "user@example.com", Password: "supersecret", Age: 25}

	err := ValidateRegistration(form)
	if err != nil {
		t.Errorf("expected nil error for a valid form, got: %v", err)
	}
}

func TestValidateRegistration_EmptyEmail(t *testing.T) {
	form := RegistrationForm{Email: "", Password: "supersecret", Age: 25}

	err := ValidateRegistration(form)
	if err == nil {
		t.Fatal("expected an error for empty email, got nil")
	}

	var valErr *ValidationError
	if !errors.As(err, &valErr) {
		t.Fatalf("expected *ValidationError, got %T", err)
	}
	if !containsField(valErr.Fields, "email") {
		t.Errorf("expected Fields to contain %q, got %v", "email", valErr.Fields)
	}
}

func TestValidateRegistration_ShortPassword(t *testing.T) {
	form := RegistrationForm{Email: "user@example.com", Password: "short", Age: 25}

	err := ValidateRegistration(form)
	if err == nil {
		t.Fatal("expected an error for too-short password, got nil")
	}

	var valErr *ValidationError
	if !errors.As(err, &valErr) {
		t.Fatalf("expected *ValidationError, got %T", err)
	}
	if !containsField(valErr.Fields, "password") {
		t.Errorf("expected Fields to contain %q, got %v", "password", valErr.Fields)
	}
}

func TestValidateRegistration_PasswordCountsRunes(t *testing.T) {
	form := RegistrationForm{Email: "user@example.com", Password: "парольки", Age: 25}

	if err := ValidateRegistration(form); err != nil {
		t.Errorf("expected an 8-rune password to be valid, got: %v", err)
	}
}

func TestValidateRegistration_InvalidAge(t *testing.T) {
	cases := []struct {
		name string
		age  int
	}{
		{"negative", -1},
		{"too old", 151},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			form := RegistrationForm{Email: "user@example.com", Password: "supersecret", Age: tc.age}

			err := ValidateRegistration(form)
			if err == nil {
				t.Fatal("expected an error for invalid age, got nil")
			}

			var valErr *ValidationError
			if !errors.As(err, &valErr) {
				t.Fatalf("expected *ValidationError, got %T", err)
			}
			if !containsField(valErr.Fields, "age") {
				t.Errorf("expected Fields to contain %q, got %v", "age", valErr.Fields)
			}
		})
	}
}

// Ключовий тест завдання: усі невалідні поля мають потрапити у ВІДПОВІДЬ
// одночасно — валідація не має зупинятись на першій знайденій помилці.
func TestValidateRegistration_ReportsAllInvalidFieldsAtOnce(t *testing.T) {
	form := RegistrationForm{Email: "", Password: "", Age: -5}

	err := ValidateRegistration(form)
	if err == nil {
		t.Fatal("expected an error for a form with multiple invalid fields, got nil")
	}

	var valErr *ValidationError
	if !errors.As(err, &valErr) {
		t.Fatalf("expected *ValidationError, got %T", err)
	}

	wantFields := []string{"email", "password", "age"}
	if !reflect.DeepEqual(valErr.Fields, wantFields) {
		t.Errorf("ValidationError.Fields = %v, want %v", valErr.Fields, wantFields)
	}
}

func TestValidateRegistration_ErrorMessageIsReadable(t *testing.T) {
	form := RegistrationForm{Email: "", Password: "supersecret", Age: 25}

	err := ValidateRegistration(form)
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if err.Error() == "" {
		t.Error("expected Error() to return a non-empty, readable message")
	}
}

func TestValidateRegistration_AgeBoundariesAreValid(t *testing.T) {
	for _, age := range []int{0, 150} {
		t.Run(fmt.Sprintf("age=%d", age), func(t *testing.T) {
			form := RegistrationForm{Email: "user@example.com", Password: "supersecret", Age: age}
			if err := ValidateRegistration(form); err != nil {
				t.Errorf("expected age %d to be valid, got: %v", age, err)
			}
		})
	}
}

func containsField(fields []string, target string) bool {
	for _, f := range fields {
		if f == target {
			return true
		}
	}
	return false
}
