package repository

import (
	"context"

	"todo/internal/domain/student"
)

// Repository defines the interface for student data persistence.
type Repository interface {
	// Save saves a new student record.
	// Source: "系統應該成功建立學生記錄" (第 8 行)
	Save(ctx context.Context, s *student.Student) error

	// FindByStudentNumber retrieves a student by student number.
	// Source: "我使用學號查詢學生" (第 14 行)
	FindByStudentNumber(ctx context.Context, studentNumber string) (*student.Student, error)

	// FindAll retrieves all student records.
	// Source: "我請求查詢所有學生" (第 20 行)
	FindAll(ctx context.Context) ([]*student.Student, error)

	// Update updates an existing student record.
	// Source: "系統應該成功更新學生記錄" (第 27 行)
	Update(ctx context.Context, s *student.Student) error

	// Delete deletes a student record by student number.
	// Source: "系統應該成功刪除該學生" (第 33 行)
	Delete(ctx context.Context, studentNumber string) error

	// ExistsByStudentNumber checks if a student number exists.
	// Used for uniqueness validation.
	ExistsByStudentNumber(ctx context.Context, studentNumber string) (bool, error)
}
