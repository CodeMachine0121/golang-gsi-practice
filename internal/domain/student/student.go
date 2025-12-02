package student

import "time"

// Student represents a student entity in the system.
// Source: "我提交新學生資訊，包含姓名、學號、電子郵件和班級" (第 7 行)
type Student struct {
	ID            string    `json:"id"`
	StudentNumber string    `json:"student_number"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	Class         string    `json:"class"`
	Grade         *int      `json:"grade,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// CreateStudentRequest represents the request for creating a student.
// Source: "我提交新學生資訊" (第 7 行)
type CreateStudentRequest struct {
	StudentNumber string `json:"student_number"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	Class         string `json:"class"`
	Grade         *int   `json:"grade,omitempty"`
}

// UpdateStudentRequest represents the request for updating a student.
// Supports partial updates where only provided fields are updated.
// Source: "我將該學生的電子郵件更新" (第 26 行)
type UpdateStudentRequest struct {
	StudentNumber *string `json:"student_number,omitempty"`
	Name          *string `json:"name,omitempty"`
	Email         *string `json:"email,omitempty"`
	Class         *string `json:"class,omitempty"`
	Grade         *int    `json:"grade,omitempty"`
}

// MinGrade and MaxGrade define the valid range for student grade.
// Source: "年級必須在 1-6 之間" (第 82 行)
const (
	MinGrade = 1
	MaxGrade = 6
)
