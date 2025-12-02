package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"todo/internal/domain/student"
	studentusecase "todo/internal/usecase/student"
)

// Handler handles HTTP requests for student management.
type Handler struct {
	useCase *studentusecase.UseCase
}

// NewHandler creates a new student HTTP handler.
func NewHandler(useCase *studentusecase.UseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

// ErrorResponse represents a standard error response.
type ErrorResponse struct {
	Error string `json:"error"`
	Code  string `json:"code,omitempty"`
}

// CreateStudent handles POST /api/students
// Source: "我提交新學生資訊" (第 5-10 行)
//
// When: 我提交新學生資訊，包含姓名、學號、電子郵件和班級
// Then: 系統應該成功建立學生記錄
func (h *Handler) CreateStudent(c *gin.Context) {
	var req student.CreateStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request format",
			Code:  "INVALID_REQUEST",
		})
		return
	}

	s, err := h.useCase.CreateStudent(c.Request.Context(), &req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, s)
}

// GetStudent handles GET /api/students/:studentNumber
// Source: "我使用學號查詢學生" (第 12-16 行)
//
// When: 我使用學號「2024001」查詢學生
// Then: 系統應該返回該學生的完整資訊
func (h *Handler) GetStudent(c *gin.Context) {
	studentNumber := c.Param("studentNumber")

	s, err := h.useCase.GetStudent(c.Request.Context(), studentNumber)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, s)
}

// GetAllStudents handles GET /api/students
// Source: "我請求查詢所有學生" (第 18-22 行)
//
// When: 我請求查詢所有學生
// Then: 系統應該返回所有學生記錄
func (h *Handler) GetAllStudents(c *gin.Context) {
	students, err := h.useCase.GetAllStudents(c.Request.Context())
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, students)
}

// UpdateStudent handles PUT /api/students/:studentNumber
// Source: "我將該學生的電子郵件更新" (第 24-28 行)
//
// When: 我將該學生的電子郵件更新為「wang.new@school.edu」
// Then: 系統應該成功更新學生記錄
func (h *Handler) UpdateStudent(c *gin.Context) {
	studentNumber := c.Param("studentNumber")

	var req student.UpdateStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request format",
			Code:  "INVALID_REQUEST",
		})
		return
	}

	s, err := h.useCase.UpdateStudent(c.Request.Context(), studentNumber, &req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, s)
}

// DeleteStudent handles DELETE /api/students/:studentNumber
// Source: "我請求刪除該學生記錄" (第 30-34 行)
//
// When: 我請求刪除該學生記錄
// Then: 系統應該成功刪除該學生
func (h *Handler) DeleteStudent(c *gin.Context) {
	studentNumber := c.Param("studentNumber")

	err := h.useCase.DeleteStudent(c.Request.Context(), studentNumber)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// handleError maps domain errors to HTTP responses.
func (h *Handler) handleError(c *gin.Context, err error) {
	var studentErr *student.StudentError
	if errors.As(err, &studentErr) {
		switch studentErr.Type {
		case student.ErrorTypeMissingRequiredField:
			// Source: "姓名為必填欄位" (第 39 行)
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: studentErr.Message,
				Code:  string(studentErr.Type),
			})
		case student.ErrorTypeInvalidEmail:
			// Source: "無效的電子郵件格式" (第 51 行)
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: studentErr.Message,
				Code:  string(studentErr.Type),
			})
		case student.ErrorTypeInvalidGrade:
			// Source: "年級必須在 1-6 之間" (第 82 行)
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: studentErr.Message,
				Code:  string(studentErr.Type),
			})
		case student.ErrorTypeStudentNumberAlreadyExists:
			// Source: "學號已存在" (第 45 行)
			c.JSON(http.StatusConflict, ErrorResponse{
				Error: studentErr.Message,
				Code:  string(studentErr.Type),
			})
		case student.ErrorTypeStudentNotFound:
			// Source: "學生不存在" (第 57 行)
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: studentErr.Message,
				Code:  string(studentErr.Type),
			})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "Internal server error",
				Code:  "INTERNAL_ERROR",
			})
		}
		return
	}

	// Unknown error
	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Error: "Internal server error",
		Code:  "INTERNAL_ERROR",
	})
}

// RegisterRoutes registers all student routes to the router.
func RegisterRoutes(router *gin.Engine, handler *Handler) {
	group := router.Group("/api/students")
	{
		group.POST("", handler.CreateStudent)
		group.GET("", handler.GetAllStudents)
		group.GET("/:studentNumber", handler.GetStudent)
		group.PUT("/:studentNumber", handler.UpdateStudent)
		group.DELETE("/:studentNumber", handler.DeleteStudent)
	}
}
