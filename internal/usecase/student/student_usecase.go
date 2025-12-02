package usecase

import (
	"context"
	"net/mail"
	"time"

	"github.com/google/uuid"

	"todo/internal/domain/student"
	studentrepo "todo/internal/repository/student"
)

// UseCase handles all business logic for student management.
// Satisfies scenarios from lines 5-83 of the feature specification.
type UseCase struct {
	repo studentrepo.Repository
}

// NewUseCase creates a new StudentUseCase.
func NewUseCase(repo studentrepo.Repository) *UseCase {
	return &UseCase{
		repo: repo,
	}
}

// CreateStudent creates a new student with validation.
// Source: "我提交新學生資訊" (第 5-10 行)
//
// Given: 系統已初始化
// When: 我提交新學生資訊，包含姓名、學號、電子郵件和班級
// Then: 系統應該成功建立學生記錄，並返回學生 ID
func (uc *UseCase) CreateStudent(ctx context.Context, req *student.CreateStudentRequest) (*student.Student, error) {
	// Validate required fields (第 36-40 行)
	if err := validateCreateRequest(req); err != nil {
		return nil, err
	}

	// Validate email format (第 48-52 行)
	if err := validateEmail(req.Email); err != nil {
		return nil, err
	}

	// Validate grade if provided (第 79-83 行)
	if req.Grade != nil && (*req.Grade < student.MinGrade || *req.Grade > student.MaxGrade) {
		return nil, student.NewInvalidGradeError()
	}

	// Check student number uniqueness (第 42-46 行)
	exists, err := uc.repo.ExistsByStudentNumber(ctx, req.StudentNumber)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, student.NewStudentNumberAlreadyExistsError()
	}

	// Create student entity
	now := time.Now()
	s := &student.Student{
		ID:            uuid.New().String(),
		StudentNumber: req.StudentNumber,
		Name:          req.Name,
		Email:         req.Email,
		Class:         req.Class,
		Grade:         req.Grade,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	// Save to repository
	if err := uc.repo.Save(ctx, s); err != nil {
		return nil, err
	}

	return s, nil
}

// GetStudent retrieves a student by student number.
// Source: "我使用學號查詢學生" (第 12-16 行)
//
// Given: 系統中已存在學號為「2024001」的學生記錄
// When: 我使用學號「2024001」查詢學生
// Then: 系統應該返回該學生的完整資訊
func (uc *UseCase) GetStudent(ctx context.Context, studentNumber string) (*student.Student, error) {
	s, err := uc.repo.FindByStudentNumber(ctx, studentNumber)
	if err != nil {
		// Returns StudentNotFound error (第 54-58 行)
		return nil, err
	}
	return s, nil
}

// GetAllStudents retrieves all students.
// Source: "我請求查詢所有學生" (第 18-22 行)
//
// Given: 系統中已存在 5 筆學生記錄
// When: 我請求查詢所有學生
// Then: 系統應該返回所有 5 筆學生記錄
func (uc *UseCase) GetAllStudents(ctx context.Context) ([]*student.Student, error) {
	students, err := uc.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	// Return empty slice if no students exist
	if students == nil {
		students = make([]*student.Student, 0)
	}

	return students, nil
}

// UpdateStudent updates an existing student with partial update support.
// Source: "我將該學生的電子郵件更新" (第 24-28 行)
//
// Given: 系統中已存在學號為「2024001」的學生記錄
// When: 我將該學生的電子郵件更新為「wang.new@school.edu」
// Then: 系統應該成功更新學生記錄
func (uc *UseCase) UpdateStudent(ctx context.Context, studentNumber string, req *student.UpdateStudentRequest) (*student.Student, error) {
	// Get existing student (第 60-64 行 for not found error)
	existing, err := uc.repo.FindByStudentNumber(ctx, studentNumber)
	if err != nil {
		return nil, err
	}

	// Apply partial updates (第 72-77 行)
	if req.StudentNumber != nil {
		// Check uniqueness if changing student number
		if *req.StudentNumber != existing.StudentNumber {
			exists, err := uc.repo.ExistsByStudentNumber(ctx, *req.StudentNumber)
			if err != nil {
				return nil, err
			}
			if exists {
				return nil, student.NewStudentNumberAlreadyExistsError()
			}
			existing.StudentNumber = *req.StudentNumber
		}
	}

	if req.Name != nil {
		if *req.Name == "" {
			return nil, student.NewMissingRequiredFieldError("Name")
		}
		existing.Name = *req.Name
	}

	if req.Email != nil {
		if err := validateEmail(*req.Email); err != nil {
			return nil, err
		}
		existing.Email = *req.Email
	}

	if req.Class != nil {
		if *req.Class == "" {
			return nil, student.NewMissingRequiredFieldError("Class")
		}
		existing.Class = *req.Class
	}

	if req.Grade != nil {
		if *req.Grade < student.MinGrade || *req.Grade > student.MaxGrade {
			return nil, student.NewInvalidGradeError()
		}
		existing.Grade = req.Grade
	}

	// Update timestamp
	existing.UpdatedAt = time.Now()

	// Save updated student
	if err := uc.repo.Update(ctx, existing); err != nil {
		return nil, err
	}

	return existing, nil
}

// DeleteStudent deletes a student by student number.
// Source: "我請求刪除該學生記錄" (第 30-34 行)
//
// Given: 系統中已存在學號為「2024001」的學生記錄
// When: 我請求刪除該學生記錄
// Then: 系統應該成功刪除該學生
func (uc *UseCase) DeleteStudent(ctx context.Context, studentNumber string) error {
	// Verify student exists before deletion
	_, err := uc.repo.FindByStudentNumber(ctx, studentNumber)
	if err != nil {
		// Returns StudentNotFound error (第 66-70 行)
		return err
	}

	return uc.repo.Delete(ctx, studentNumber)
}

// validateCreateRequest validates required fields in CreateStudentRequest.
// Source: "新增時缺少必填欄位" (第 36-40 行)
func validateCreateRequest(req *student.CreateStudentRequest) error {
	if req.StudentNumber == "" {
		return student.NewMissingRequiredFieldError("StudentNumber")
	}
	if req.Name == "" {
		return student.NewMissingRequiredFieldError("Name")
	}
	if req.Email == "" {
		return student.NewMissingRequiredFieldError("Email")
	}
	if req.Class == "" {
		return student.NewMissingRequiredFieldError("Class")
	}
	return nil
}

// validateEmail validates email format.
// Source: "電子郵件格式驗證" (第 48-52 行)
func validateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return student.NewInvalidEmailError()
	}
	return nil
}
