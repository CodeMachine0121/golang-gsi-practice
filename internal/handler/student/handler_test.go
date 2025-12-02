package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"todo/internal/domain/student"
	studentrepo "todo/internal/repository/student"
	studentusecase "todo/internal/usecase/student"
)

func setupTestHandler() *Handler {
	gin.SetMode(gin.TestMode)
	repo := studentrepo.NewMemoryRepository()
	uc := studentusecase.NewUseCase(repo)
	return NewHandler(uc)
}

func TestCreateStudent_Success(t *testing.T) {
	// Scenario: 成功新增學生 (第 5-10 行)
	handler := setupTestHandler()
	router := gin.New()
	RegisterRoutes(router, handler)

	// When: 我提交新學生資訊
	grade := 1
	payload := student.CreateStudentRequest{
		StudentNumber: "2024001",
		Name:          "王小明",
		Email:         "wang@school.edu",
		Class:         "一年一班",
		Grade:         &grade,
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/api/students", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then: 系統應該成功建立學生記錄
	assert.Equal(t, http.StatusCreated, w.Code)

	var result student.Student
	json.Unmarshal(w.Body.Bytes(), &result)

	// And: 返回的學生 ID 應該不為空
	assert.NotEmpty(t, result.ID)

	// And: 返回的學生資訊應該與提交的資訊相符
	assert.Equal(t, "2024001", result.StudentNumber)
	assert.Equal(t, "王小明", result.Name)
	assert.Equal(t, "wang@school.edu", result.Email)
	assert.Equal(t, "一年一班", result.Class)
}

func TestGetStudent_Success(t *testing.T) {
	// Scenario: 查詢單一學生 (第 12-16 行)
	handler := setupTestHandler()
	router := gin.New()
	RegisterRoutes(router, handler)

	// Given: 先建立一個學生
	grade := 1
	createPayload := student.CreateStudentRequest{
		StudentNumber: "2024001",
		Name:          "王小明",
		Email:         "wang@school.edu",
		Class:         "一年一班",
		Grade:         &grade,
	}
	body, _ := json.Marshal(createPayload)
	createReq, _ := http.NewRequest("POST", "/api/students", bytes.NewBuffer(body))
	createReq.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, createReq)
	require.Equal(t, http.StatusCreated, w.Code)

	// When: 我使用學號「2024001」查詢學生
	getReq, _ := http.NewRequest("GET", "/api/students/2024001", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, getReq)

	// Then: 系統應該返回該學生的完整資訊
	assert.Equal(t, http.StatusOK, w.Code)

	var result student.Student
	json.Unmarshal(w.Body.Bytes(), &result)

	// And: 返回的資訊應該包含姓名、學號、電子郵件和班級
	assert.Equal(t, "2024001", result.StudentNumber)
	assert.Equal(t, "王小明", result.Name)
	assert.Equal(t, "wang@school.edu", result.Email)
	assert.Equal(t, "一年一班", result.Class)
}

func TestGetAllStudents_Success(t *testing.T) {
	// Scenario: 查詢所有學生 (第 18-22 行)
	handler := setupTestHandler()
	router := gin.New()
	RegisterRoutes(router, handler)

	// Given: 建立 5 筆學生記錄
	for i := 1; i <= 5; i++ {
		payload := student.CreateStudentRequest{
			StudentNumber: "202400" + string(rune(i+'0')),
			Name:          "学生" + string(rune(i+'0')),
			Email:         "student" + string(rune(i+'0')) + "@school.edu",
			Class:         "一年一班",
		}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", "/api/students", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}

	// When: 我請求查詢所有學生
	getReq, _ := http.NewRequest("GET", "/api/students", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, getReq)

	// Then: 系統應該返回所有 5 筆學生記錄
	assert.Equal(t, http.StatusOK, w.Code)

	var results []*student.Student
	json.Unmarshal(w.Body.Bytes(), &results)
	assert.Len(t, results, 5)
}

func TestUpdateStudent_Success(t *testing.T) {
	// Scenario: 成功更新學生資訊 (第 24-28 行)
	handler := setupTestHandler()
	router := gin.New()
	RegisterRoutes(router, handler)

	// Given: 先建立一個學生
	createPayload := student.CreateStudentRequest{
		StudentNumber: "2024001",
		Name:          "王小明",
		Email:         "wang@school.edu",
		Class:         "一年一班",
	}
	body, _ := json.Marshal(createPayload)
	createReq, _ := http.NewRequest("POST", "/api/students", bytes.NewBuffer(body))
	createReq.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, createReq)
	require.Equal(t, http.StatusCreated, w.Code)

	// When: 我將該學生的電子郵件更新為「wang.new@school.edu」
	updatePayload := student.UpdateStudentRequest{
		Email: strPtr("wang.new@school.edu"),
	}
	updateBody, _ := json.Marshal(updatePayload)
	updateReq, _ := http.NewRequest("PUT", "/api/students/2024001", bytes.NewBuffer(updateBody))
	updateReq.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, updateReq)

	// Then: 系統應該成功更新學生記錄
	assert.Equal(t, http.StatusOK, w.Code)

	var result student.Student
	json.Unmarshal(w.Body.Bytes(), &result)

	// And: 查詢該學生時應該返回新的電子郵件地址
	assert.Equal(t, "wang.new@school.edu", result.Email)
}

func TestDeleteStudent_Success(t *testing.T) {
	// Scenario: 成功刪除學生 (第 30-34 行)
	handler := setupTestHandler()
	router := gin.New()
	RegisterRoutes(router, handler)

	// Given: 先建立一個學生
	createPayload := student.CreateStudentRequest{
		StudentNumber: "2024001",
		Name:          "王小明",
		Email:         "wang@school.edu",
		Class:         "一年一班",
	}
	body, _ := json.Marshal(createPayload)
	createReq, _ := http.NewRequest("POST", "/api/students", bytes.NewBuffer(body))
	createReq.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, createReq)
	require.Equal(t, http.StatusCreated, w.Code)

	// When: 我請求刪除該學生記錄
	deleteReq, _ := http.NewRequest("DELETE", "/api/students/2024001", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, deleteReq)

	// Then: 系統應該成功刪除該學生
	assert.Equal(t, http.StatusNoContent, w.Code)

	// And: 再次查詢該學號時應該返回「學生不存在」的錯誤
	getReq, _ := http.NewRequest("GET", "/api/students/2024001", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, getReq)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCreateStudent_MissingRequiredField(t *testing.T) {
	// Scenario: 新增時缺少必填欄位 (第 36-40 行)
	handler := setupTestHandler()
	router := gin.New()
	RegisterRoutes(router, handler)

	// When: 我提交不完整的學生資訊，缺少姓名欄位
	payload := student.CreateStudentRequest{
		StudentNumber: "2024001",
		Name:          "", // Missing name
		Email:         "wang@school.edu",
		Class:         "一年一班",
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/api/students", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then: 系統應該拒絕並返回錯誤
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var errorResp ErrorResponse
	json.Unmarshal(w.Body.Bytes(), &errorResp)
	assert.Equal(t, student.ErrorTypeMissingRequiredField, student.ErrorType(errorResp.Code))
}

func TestCreateStudent_StudentNumberAlreadyExists(t *testing.T) {
	// Scenario: 學號必須唯一 (第 42-46 行)
	handler := setupTestHandler()
	router := gin.New()
	RegisterRoutes(router, handler)

	// Given: 先建立一個學號為「2024001」的學生
	createPayload := student.CreateStudentRequest{
		StudentNumber: "2024001",
		Name:          "王小明",
		Email:         "wang@school.edu",
		Class:         "一年一班",
	}
	body, _ := json.Marshal(createPayload)
	createReq, _ := http.NewRequest("POST", "/api/students", bytes.NewBuffer(body))
	createReq.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, createReq)

	// When: 我嘗試新增另一個學號相同「2024001」的學生
	duplicatePayload := student.CreateStudentRequest{
		StudentNumber: "2024001",
		Name:          "另一個學生",
		Email:         "another@school.edu",
		Class:         "一年一班",
	}

	dupBody, _ := json.Marshal(duplicatePayload)
	dupReq, _ := http.NewRequest("POST", "/api/students", bytes.NewBuffer(dupBody))
	dupReq.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	router.ServeHTTP(w, dupReq)

	// Then: 系統應該拒絕並返回錯誤「學號已存在」
	assert.Equal(t, http.StatusConflict, w.Code)

	var errorResp ErrorResponse
	json.Unmarshal(w.Body.Bytes(), &errorResp)
	assert.Equal(t, student.ErrorTypeStudentNumberAlreadyExists, student.ErrorType(errorResp.Code))
}

func TestCreateStudent_InvalidEmail(t *testing.T) {
	// Scenario: 電子郵件格式驗證 (第 48-52 行)
	handler := setupTestHandler()
	router := gin.New()
	RegisterRoutes(router, handler)

	// When: 我提交學生資訊，電子郵件為無效格式「invalid-email」
	payload := student.CreateStudentRequest{
		StudentNumber: "2024001",
		Name:          "王小明",
		Email:         "invalid-email",
		Class:         "一年一班",
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/api/students", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then: 系統應該拒絕並返回錯誤
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var errorResp ErrorResponse
	json.Unmarshal(w.Body.Bytes(), &errorResp)
	assert.Equal(t, student.ErrorTypeInvalidEmail, student.ErrorType(errorResp.Code))
}

func TestGetStudent_NotFound(t *testing.T) {
	// Scenario: 查詢不存在的學生 (第 54-58 行)
	handler := setupTestHandler()
	router := gin.New()
	RegisterRoutes(router, handler)

	// When: 我使用不存在的學號「9999999」查詢學生
	req, _ := http.NewRequest("GET", "/api/students/9999999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then: 系統應該返回錯誤「學生不存在」
	// And HTTP 狀態碼應該是 404
	assert.Equal(t, http.StatusNotFound, w.Code)

	var errorResp ErrorResponse
	json.Unmarshal(w.Body.Bytes(), &errorResp)
	assert.Equal(t, student.ErrorTypeStudentNotFound, student.ErrorType(errorResp.Code))
}

func TestCreateStudent_InvalidGrade(t *testing.T) {
	// Scenario: 驗證年級範圍 (第 79-83 行)
	handler := setupTestHandler()
	router := gin.New()
	RegisterRoutes(router, handler)

	// When: 我提交學生資訊，年級為無效值「10」（超出範圍）
	grade := 10
	payload := student.CreateStudentRequest{
		StudentNumber: "2024001",
		Name:          "王小明",
		Email:         "wang@school.edu",
		Class:         "一年一班",
		Grade:         &grade,
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/api/students", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then: 系統應該拒絕並返回錯誤
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var errorResp ErrorResponse
	json.Unmarshal(w.Body.Bytes(), &errorResp)
	assert.Equal(t, student.ErrorTypeInvalidGrade, student.ErrorType(errorResp.Code))
}

// Helper function for pointer to string
func strPtr(s string) *string {
	return &s
}
