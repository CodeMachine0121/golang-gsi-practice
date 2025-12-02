package usecase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"todo/internal/domain/student"
	studentrepo "todo/internal/repository/student"
)

func TestCreateStudent_Success(t *testing.T) {
	// Scenario: 成功新增學生 (第 5-10 行)
	// Given: 系統已初始化
	repo := studentrepo.NewMemoryRepository()
	uc := NewUseCase(repo)

	// When: 我提交新學生資訊
	grade := 1
	req := &student.CreateStudentRequest{
		StudentNumber: "2024001",
		Name:          "王小明",
		Email:         "wang@school.edu",
		Class:         "一年一班",
		Grade:         &grade,
	}

	// Then: 系統應該成功建立學生記錄
	s, err := uc.CreateStudent(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, s)

	// And: 返回的學生 ID 應該不為空
	assert.NotEmpty(t, s.ID)

	// And: 返回的學生資訊應該與提交的資訊相符
	assert.Equal(t, "2024001", s.StudentNumber)
	assert.Equal(t, "王小明", s.Name)
	assert.Equal(t, "wang@school.edu", s.Email)
	assert.Equal(t, "一年一班", s.Class)
	assert.Equal(t, &grade, s.Grade)
}

func TestGetStudent_Success(t *testing.T) {
	// Scenario: 查詢單一學生 (第 12-16 行)
	// Given: 系統中已存在學號為「2024001」、姓名為「王小明」的學生記錄
	repo := studentrepo.NewMemoryRepository()
	uc := NewUseCase(repo)

	grade := 1
	s := &student.Student{
		ID:            "test-id",
		StudentNumber: "2024001",
		Name:          "王小明",
		Email:         "wang@school.edu",
		Class:         "一年一班",
		Grade:         &grade,
	}
	require.NoError(t, repo.Save(context.Background(), s))

	// When: 我使用學號「2024001」查詢學生
	result, err := uc.GetStudent(context.Background(), "2024001")

	// Then: 系統應該返回該學生的完整資訊
	require.NoError(t, err)
	assert.Equal(t, "2024001", result.StudentNumber)
	assert.Equal(t, "王小明", result.Name)

	// And: 返回的資訊應該包含姓名、學號、電子郵件和班級
	assert.Equal(t, "wang@school.edu", result.Email)
	assert.Equal(t, "一年一班", result.Class)
}

func TestGetAllStudents_Success(t *testing.T) {
	// Scenario: 查詢所有學生 (第 18-22 行)
	// Given: 系統中已存在 5 筆學生記錄
	repo := studentrepo.NewMemoryRepository()
	uc := NewUseCase(repo)

	for i := 1; i <= 5; i++ {
		s := &student.Student{
			ID:            "test-id-" + string(rune(i)),
			StudentNumber: "202400" + string(rune(i+'0')),
			Name:          "学生" + string(rune(i+'0')),
			Email:         "student" + string(rune(i+'0')) + "@school.edu",
			Class:         "一年一班",
		}
		require.NoError(t, repo.Save(context.Background(), s))
	}

	// When: 我請求查詢所有學生
	students, err := uc.GetAllStudents(context.Background())

	// Then: 系統應該返回所有 5 筆學生記錄
	require.NoError(t, err)
	assert.Len(t, students, 5)

	// And: 每筆記錄都應該包含學生的完整資訊
	for _, s := range students {
		assert.NotEmpty(t, s.ID)
		assert.NotEmpty(t, s.StudentNumber)
		assert.NotEmpty(t, s.Name)
		assert.NotEmpty(t, s.Email)
		assert.NotEmpty(t, s.Class)
	}
}

func TestUpdateStudent_Success(t *testing.T) {
	// Scenario: 成功更新學生資訊 (第 24-28 行)
	// Given: 系統中已存在學號為「2024001」的學生記錄
	repo := studentrepo.NewMemoryRepository()
	uc := NewUseCase(repo)

	grade := 1
	s := &student.Student{
		ID:            "test-id",
		StudentNumber: "2024001",
		Name:          "王小明",
		Email:         "wang@school.edu",
		Class:         "一年一班",
		Grade:         &grade,
	}
	require.NoError(t, repo.Save(context.Background(), s))

	// When: 我將該學生的電子郵件更新為「wang.new@school.edu」
	newEmail := "wang.new@school.edu"
	updateReq := &student.UpdateStudentRequest{
		Email: &newEmail,
	}
	updated, err := uc.UpdateStudent(context.Background(), "2024001", updateReq)

	// Then: 系統應該成功更新學生記錄
	require.NoError(t, err)

	// And: 查詢該學生時應該返回新的電子郵件地址
	assert.Equal(t, "wang.new@school.edu", updated.Email)
	assert.Equal(t, "王小明", updated.Name) // Other fields unchanged
}

func TestDeleteStudent_Success(t *testing.T) {
	// Scenario: 成功刪除學生 (第 30-34 行)
	// Given: 系統中已存在學號為「2024001」的學生記錄
	repo := studentrepo.NewMemoryRepository()
	uc := NewUseCase(repo)

	s := &student.Student{
		ID:            "test-id",
		StudentNumber: "2024001",
		Name:          "王小明",
		Email:         "wang@school.edu",
		Class:         "一年一班",
	}
	require.NoError(t, repo.Save(context.Background(), s))

	// When: 我請求刪除該學生記錄
	err := uc.DeleteStudent(context.Background(), "2024001")

	// Then: 系統應該成功刪除該學生
	require.NoError(t, err)

	// And: 再次查詢該學號時應該返回「學生不存在」的錯誤
	_, err = uc.GetStudent(context.Background(), "2024001")
	require.Error(t, err)
	var studentErr *student.StudentError
	require.ErrorAs(t, err, &studentErr)
	assert.Equal(t, student.ErrorTypeStudentNotFound, studentErr.Type)
}

func TestCreateStudent_MissingRequiredField(t *testing.T) {
	// Scenario: 新增時缺少必填欄位 (第 36-40 行)
	// Given: 系統已初始化
	repo := studentrepo.NewMemoryRepository()
	uc := NewUseCase(repo)

	// When: 我提交不完整的學生資訊，缺少姓名欄位
	req := &student.CreateStudentRequest{
		StudentNumber: "2024001",
		Name:          "", // Missing name
		Email:         "wang@school.edu",
		Class:         "一年一班",
	}

	// Then: 系統應該拒絕並返回錯誤
	_, err := uc.CreateStudent(context.Background(), req)
	require.Error(t, err)

	var studentErr *student.StudentError
	require.ErrorAs(t, err, &studentErr)
	assert.Equal(t, student.ErrorTypeMissingRequiredField, studentErr.Type)

	// And: 學生記錄不應該被建立
	students, _ := uc.GetAllStudents(context.Background())
	assert.Len(t, students, 0)
}

func TestCreateStudent_StudentNumberAlreadyExists(t *testing.T) {
	// Scenario: 學號必須唯一 (第 42-46 行)
	// Given: 系統中已存在學號為「2024001」的學生記錄
	repo := studentrepo.NewMemoryRepository()
	uc := NewUseCase(repo)

	s := &student.Student{
		ID:            "test-id",
		StudentNumber: "2024001",
		Name:          "王小明",
		Email:         "wang@school.edu",
		Class:         "一年一班",
	}
	require.NoError(t, repo.Save(context.Background(), s))

	// When: 我嘗試新增另一個學號相同「2024001」的學生
	req := &student.CreateStudentRequest{
		StudentNumber: "2024001",
		Name:          "另一個學生",
		Email:         "another@school.edu",
		Class:         "一年一班",
	}

	// Then: 系統應該拒絕並返回錯誤「學號已存在」
	_, err := uc.CreateStudent(context.Background(), req)
	require.Error(t, err)

	var studentErr *student.StudentError
	require.ErrorAs(t, err, &studentErr)
	assert.Equal(t, student.ErrorTypeStudentNumberAlreadyExists, studentErr.Type)

	// And: 新的學生記錄不應該被建立
	students, _ := uc.GetAllStudents(context.Background())
	assert.Len(t, students, 1)
}

func TestCreateStudent_InvalidEmail(t *testing.T) {
	// Scenario: 電子郵件格式驗證 (第 48-52 行)
	// Given: 系統已初始化
	repo := studentrepo.NewMemoryRepository()
	uc := NewUseCase(repo)

	// When: 我提交學生資訊，電子郵件為無效格式「invalid-email」
	req := &student.CreateStudentRequest{
		StudentNumber: "2024001",
		Name:          "王小明",
		Email:         "invalid-email",
		Class:         "一年一班",
	}

	// Then: 系統應該拒絕並返回錯誤
	_, err := uc.CreateStudent(context.Background(), req)
	require.Error(t, err)

	var studentErr *student.StudentError
	require.ErrorAs(t, err, &studentErr)
	assert.Equal(t, student.ErrorTypeInvalidEmail, studentErr.Type)

	// And: 學生記錄不應該被建立
	students, _ := uc.GetAllStudents(context.Background())
	assert.Len(t, students, 0)
}

func TestGetStudent_NotFound(t *testing.T) {
	// Scenario: 查詢不存在的學生 (第 54-58 行)
	// Given: 系統已初始化
	repo := studentrepo.NewMemoryRepository()
	uc := NewUseCase(repo)

	// When: 我使用不存在的學號「9999999」查詢學生
	_, err := uc.GetStudent(context.Background(), "9999999")

	// Then: 系統應該返回錯誤「學生不存在」
	require.Error(t, err)

	var studentErr *student.StudentError
	require.ErrorAs(t, err, &studentErr)
	assert.Equal(t, student.ErrorTypeStudentNotFound, studentErr.Type)
}

func TestUpdateStudent_NotFound(t *testing.T) {
	// Scenario: 更新不存在的學生 (第 60-64 行)
	// Given: 系統已初始化
	repo := studentrepo.NewMemoryRepository()
	uc := NewUseCase(repo)

	// When: 我嘗試更新不存在的學號「9999999」的學生資訊
	newEmail := "new@school.edu"
	updateReq := &student.UpdateStudentRequest{
		Email: &newEmail,
	}
	_, err := uc.UpdateStudent(context.Background(), "9999999", updateReq)

	// Then: 系統應該返回錯誤「學生不存在」
	require.Error(t, err)

	var studentErr *student.StudentError
	require.ErrorAs(t, err, &studentErr)
	assert.Equal(t, student.ErrorTypeStudentNotFound, studentErr.Type)

	// And: 不應該建立新的學生記錄
	students, _ := uc.GetAllStudents(context.Background())
	assert.Len(t, students, 0)
}

func TestDeleteStudent_NotFound(t *testing.T) {
	// Scenario: 刪除不存在的學生 (第 66-70 行)
	// Given: 系統已初始化
	repo := studentrepo.NewMemoryRepository()
	uc := NewUseCase(repo)

	// When: 我嘗試刪除不存在的學號「9999999」的學生
	err := uc.DeleteStudent(context.Background(), "9999999")

	// Then: 系統應該返回錯誤「學生不存在」
	require.Error(t, err)

	var studentErr *student.StudentError
	require.ErrorAs(t, err, &studentErr)
	assert.Equal(t, student.ErrorTypeStudentNotFound, studentErr.Type)
}

func TestUpdateStudent_PartialUpdate(t *testing.T) {
	// Scenario: 部分更新學生資訊 (第 72-77 行)
	// Given: 系統中已存在學號為「2024001」的學生記錄，班級為「一年一班」
	repo := studentrepo.NewMemoryRepository()
	uc := NewUseCase(repo)

	s := &student.Student{
		ID:            "test-id",
		StudentNumber: "2024001",
		Name:          "王小明",
		Email:         "wang@school.edu",
		Class:         "一年一班",
	}
	require.NoError(t, repo.Save(context.Background(), s))

	// When: 我只更新該學生的班級為「一年二班」，不更新其他欄位
	newClass := "一年二班"
	updateReq := &student.UpdateStudentRequest{
		Class: &newClass,
	}
	updated, err := uc.UpdateStudent(context.Background(), "2024001", updateReq)

	// Then: 系統應該成功更新班級欄位
	require.NoError(t, err)
	assert.Equal(t, "一年二班", updated.Class)

	// And: 其他欄位應該保持不變
	assert.Equal(t, "王小明", updated.Name)
	assert.Equal(t, "wang@school.edu", updated.Email)

	// And: 學號應該仍然是「2024001」
	assert.Equal(t, "2024001", updated.StudentNumber)
}

func TestCreateStudent_InvalidGrade(t *testing.T) {
	// Scenario: 驗證年級範圍 (第 79-83 行)
	// Given: 系統已初始化
	repo := studentrepo.NewMemoryRepository()
	uc := NewUseCase(repo)

	// When: 我提交學生資訊，年級為無效值「10」（超出範圍）
	grade := 10
	req := &student.CreateStudentRequest{
		StudentNumber: "2024001",
		Name:          "王小明",
		Email:         "wang@school.edu",
		Class:         "一年一班",
		Grade:         &grade,
	}

	// Then: 系統應該拒絕並返回錯誤
	_, err := uc.CreateStudent(context.Background(), req)
	require.Error(t, err)

	var studentErr *student.StudentError
	require.ErrorAs(t, err, &studentErr)
	assert.Equal(t, student.ErrorTypeInvalidGrade, studentErr.Type)

	// And: 學生記錄不應該被建立
	students, _ := uc.GetAllStudents(context.Background())
	assert.Len(t, students, 0)
}
