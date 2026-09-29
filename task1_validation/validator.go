// Package validation містить завдання 1 домашньої роботи: валідатор форми
// реєстрації, що повертає власний (custom) тип помилки з переліком УСІХ
// невалідних полів, а не лише першого.
package validation

import (
	"strings"
	"unicode/utf8"
)

// RegistrationForm — вхідні дані форми реєстрації, які потрібно перевірити.
type RegistrationForm struct {
	Email    string
	Password string
	Age      int
}

// ValidationError — власний тип помилки, що переносить структуровані дані:
// список назв усіх полів, які не пройшли валідацію.
type ValidationError struct {
	Fields []string
}

// Error реалізує інтерфейс error.
func (e *ValidationError) Error() string {
	return "registration invalid: fields " + strings.Join(e.Fields, ", ")
}

// ValidateRegistration перевіряє форму реєстрації та повертає
// *ValidationError з переліком ВСІХ невалідних полів, якщо форма невалідна.
// Якщо форма валідна, повертає nil.
//
// Правила валідації:
//   - Email не може бути порожнім
//   - Password не може бути порожнім і має містити щонайменше 8 символів
//   - Age має бути в межах [0, 150]
func ValidateRegistration(f RegistrationForm) error {
	var invalidFields []string

	// Check email
	if f.Email == "" {
		invalidFields = append(invalidFields, "email")
	}

	// Check password
	if utf8.RuneCountInString(f.Password) < 8 {
		invalidFields = append(invalidFields, "password")
	}

	// Check age
	if f.Age < 0 || f.Age > 150 {
		invalidFields = append(invalidFields, "age")
	}

	if len(invalidFields) > 0 {
		return &ValidationError{Fields: invalidFields}
	}

	return nil
}
