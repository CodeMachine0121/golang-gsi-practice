package student

import "fmt"

// ErrorType represents different types of student domain errors.
// Source: 各驗證場景（第 36-83 行）
type ErrorType string

const (
	// ErrorTypeMissingRequiredField indicates a required field is missing.
	// Source: "姓名為必填欄位" (第 39 行)
	ErrorTypeMissingRequiredField ErrorType = "MISSING_REQUIRED_FIELD"

	// ErrorTypeInvalidEmail indicates the email format is invalid.
	// Source: "無效的電子郵件格式" (第 51 行)
	ErrorTypeInvalidEmail ErrorType = "INVALID_EMAIL"

	// ErrorTypeInvalidGrade indicates the grade is outside valid range.
	// Source: "年級必須在 1-6 之間" (第 82 行)
	ErrorTypeInvalidGrade ErrorType = "INVALID_GRADE"

	// ErrorTypeStudentNumberAlreadyExists indicates student number is duplicate.
	// Source: "學號已存在" (第 45 行)
	ErrorTypeStudentNumberAlreadyExists ErrorType = "STUDENT_NUMBER_ALREADY_EXISTS"

	// ErrorTypeStudentNotFound indicates student does not exist.
	// Source: "學生不存在" (第 57 行)
	ErrorTypeStudentNotFound ErrorType = "STUDENT_NOT_FOUND"
)

// StudentError represents a domain error in student operations.
type StudentError struct {
	Type    ErrorType
	Message string
	Field   string // For field-specific errors
}

// Error implements the error interface.
func (e *StudentError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("[%s] %s: %s", e.Type, e.Field, e.Message)
	}
	return fmt.Sprintf("[%s] %s", e.Type, e.Message)
}

// NewMissingRequiredFieldError creates a new missing required field error.
func NewMissingRequiredFieldError(field string) *StudentError {
	return &StudentError{
		Type:    ErrorTypeMissingRequiredField,
		Message: fmt.Sprintf("%s為必填欄位", field),
		Field:   field,
	}
}

// NewInvalidEmailError creates a new invalid email error.
func NewInvalidEmailError() *StudentError {
	return &StudentError{
		Type:    ErrorTypeInvalidEmail,
		Message: "無效的電子郵件格式",
		Field:   "email",
	}
}

// NewInvalidGradeError creates a new invalid grade error.
func NewInvalidGradeError() *StudentError {
	return &StudentError{
		Type:    ErrorTypeInvalidGrade,
		Message: fmt.Sprintf("年級必須在 %d-%d 之間", MinGrade, MaxGrade),
		Field:   "grade",
	}
}

// NewStudentNumberAlreadyExistsError creates a new duplicate student number error.
func NewStudentNumberAlreadyExistsError() *StudentError {
	return &StudentError{
		Type:    ErrorTypeStudentNumberAlreadyExists,
		Message: "學號已存在",
		Field:   "student_number",
	}
}

// NewStudentNotFoundError creates a new student not found error.
func NewStudentNotFoundError() *StudentError {
	return &StudentError{
		Type:    ErrorTypeStudentNotFound,
		Message: "學生不存在",
	}
}
