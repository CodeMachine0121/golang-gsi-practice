# 學生資訊管理 - 架構設計

> 來源：features/student_crud_api.feature
> 建立日期：2025-12-02

## 1. 專案上下文

- **程式語言**：Go 1.25.4
- **框架**：Gin Web Framework
- **測試框架**：Testify
- **架構模式**：DDD（Domain-Driven Design）- Use Case 為主
- **命名慣例**：PascalCase（型別/結構體）、camelCase（方法/變數）

## 2. 功能概述

學生資訊管理系統提供完整的 CRUD 操作，允許學校管理員創建、查詢、更新和刪除學生記錄。系統需驗證學生資訊的完整性、唯一性及格式正確性，確保資料品質。

## 3. 資料模型

### 3.1 列舉/常數

#### GradeLevel（年級）
來源：「年級必須在 1-6 之間」(第 82 行)

| 值 | 說明 |
|---|---|
| 1-6 | 有效年級範圍（國小 1-6 年級） |

### 3.2 核心實體

#### Student（學生）
來源：「我提交新學生資訊，包含姓名、學號、電子郵件和班級」(第 7 行)

| 欄位 | 型別 | 必填 | 說明 |
|---|---|---|---|
| ID | string | ✅ | 學生系統唯一識別碼（UUID） |
| StudentNumber | string | ✅ | 學號，唯一值 |
| Name | string | ✅ | 學生姓名 |
| Email | string | ✅ | 電子郵件地址，需驗證格式 |
| Class | string | ✅ | 班級名稱（如「一年一班」） |
| Grade | int | - | 年級（1-6，可選） |
| CreatedAt | time.Time | ✅ | 建立時間 |
| UpdatedAt | time.Time | ✅ | 最後更新時間 |

#### CreateStudentRequest（建立學生請求）
來源：「我提交新學生資訊」(第 7 行)

| 欄位 | 型別 | 必填 | 說明 |
|---|---|---|---|
| StudentNumber | string | ✅ | 學號 |
| Name | string | ✅ | 姓名 |
| Email | string | ✅ | 電子郵件 |
| Class | string | ✅ | 班級 |
| Grade | int | - | 年級（可選） |

#### UpdateStudentRequest（更新學生請求）
來源：「我將該學生的電子郵件更新」(第 26 行)

| 欄位 | 型別 | 必填 | 說明 |
|---|---|---|---|
| StudentNumber | string | - | 學號（部分更新） |
| Name | string | - | 姓名（部分更新） |
| Email | string | - | 電子郵件（部分更新） |
| Class | string | - | 班級（部分更新） |
| Grade | int | - | 年級（部分更新） |

### 3.3 錯誤/異常

#### StudentError
來源：各驗證場景（第 36-83 行）

| 錯誤型別 | 狀態碼 | 說明 |
|---|---|---|
| MissingRequiredField | 400 | 缺少必填欄位 |
| InvalidEmail | 400 | 電子郵件格式無效 |
| InvalidGrade | 400 | 年級超出範圍（1-6） |
| StudentNumberAlreadyExists | 409 | 學號已存在（重複） |
| StudentNotFound | 404 | 學生不存在 |

## 4. 服務介面

### 4.1 StudentUseCase
職責：協調學生資訊管理的業務邏輯，包括驗證、建立、查詢、更新和刪除操作

#### CreateStudent()
來源：「我提交新學生資訊」(第 5-10 行)

**簽名：** `CreateStudent(ctx context.Context, req *CreateStudentRequest) (*Student, error)`

**參數：**
| 參數 | 型別 | 說明 |
|---|---|---|
| ctx | context.Context | 請求上下文 |
| req | *CreateStudentRequest | 建立學生的請求資訊 |

**回傳：** `(*Student, error)` - 建立的學生物件或錯誤

**業務規則：**
1. 驗證必填欄位（姓名、學號、電子郵件、班級）
2. 驗證電子郵件格式有效性
3. 驗證學號唯一性（不能重複）
4. 驗證年級在 1-6 範圍內（若提供）
5. 生成唯一的學生 ID
6. 記錄建立和更新時間戳

---

#### GetStudent()
來源：「我使用學號查詢學生」(第 12-16 行)

**簽名：** `GetStudent(ctx context.Context, studentNumber string) (*Student, error)`

**參數：**
| 參數 | 型別 | 說明 |
|---|---|---|
| ctx | context.Context | 請求上下文 |
| studentNumber | string | 要查詢的學號 |

**回傳：** `(*Student, error)` - 學生物件或錯誤

**業務規則：**
1. 根據學號查詢學生
2. 若學生不存在，返回 StudentNotFound 錯誤（404）
3. 返回完整的學生資訊

---

#### GetAllStudents()
來源：「我請求查詢所有學生」(第 18-22 行)

**簽名：** `GetAllStudents(ctx context.Context) ([]*Student, error)`

**參數：**
| 參數 | 型別 | 說明 |
|---|---|---|
| ctx | context.Context | 請求上下文 |

**回傳：** `([]*Student, error)` - 所有學生物件的切片或錯誤

**業務規則：**
1. 查詢系統中所有學生記錄
2. 返回完整的學生資訊列表
3. 若無學生，返回空切片

---

#### UpdateStudent()
來源：「我將該學生的電子郵件更新」(第 24-28 行)

**簽名：** `UpdateStudent(ctx context.Context, studentNumber string, req *UpdateStudentRequest) (*Student, error)`

**參數：**
| 參數 | 型別 | 說明 |
|---|---|---|
| ctx | context.Context | 請求上下文 |
| studentNumber | string | 要更新的學號 |
| req | *UpdateStudentRequest | 包含要更新的欄位 |

**回傳：** `(*Student, error)` - 更新後的學生物件或錯誤

**業務規則：**
1. 驗證學生存在（根據學號）
2. 支援部分更新（只更新提供的欄位）
3. 驗證更新的欄位（電子郵件格式、年級範圍）
4. 若更新學號，驗證新學號唯一性
5. 更新 UpdatedAt 時間戳
6. 返回更新後的完整學生資訊

---

#### DeleteStudent()
來源：「我請求刪除該學生記錄」(第 30-34 行)

**簽名：** `DeleteStudent(ctx context.Context, studentNumber string) error`

**參數：**
| 參數 | 型別 | 說明 |
|---|---|---|
| ctx | context.Context | 請求上下文 |
| studentNumber | string | 要刪除的學號 |

**回傳：** `error` - 成功為 nil，失敗返回錯誤

**業務規則：**
1. 驗證學生存在（根據學號）
2. 若不存在，返回 StudentNotFound 錯誤（404）
3. 刪除該學生記錄
4. 刪除成功返回 nil

---

### 4.2 StudentRepository
職責：學生資料持久化層，管理資料的儲存和查詢

#### Save()
來源：「系統應該成功建立學生記錄」(第 8 行)

**簽名：** `Save(ctx context.Context, student *Student) error`

**參數：**
| 參數 | 型別 | 說明 |
|---|---|---|
| ctx | context.Context | 請求上下文 |
| student | *Student | 要保存的學生物件 |

**回傳：** `error` - 成功為 nil

**業務規則：**
1. 插入新的學生記錄到資料庫
2. 驗證學號唯一性（資料庫層面）

---

#### FindByStudentNumber()
來源：「我使用學號查詢學生」(第 14 行)

**簽名：** `FindByStudentNumber(ctx context.Context, studentNumber string) (*Student, error)`

**參數：**
| 參數 | 型別 | 說明 |
|---|---|---|
| ctx | context.Context | 請求上下文 |
| studentNumber | string | 要查詢的學號 |

**回傳：** `(*Student, error)` - 學生物件或錯誤

**業務規則：**
1. 根據學號查詢單一學生
2. 若不存在返回 nil

---

#### FindAll()
來源：「我請求查詢所有學生」(第 20 行)

**簽名：** `FindAll(ctx context.Context) ([]*Student, error)`

**參數：**
| 參數 | 型別 | 說明 |
|---|---|---|
| ctx | context.Context | 請求上下文 |

**回傳：** `([]*Student, error)` - 學生物件切片或錯誤

**業務規則：**
1. 查詢所有學生記錄
2. 返回完整列表，若無記錄返回空切片

---

#### Update()
來源：「系統應該成功更新學生記錄」(第 27 行)

**簽名：** `Update(ctx context.Context, student *Student) error`

**參數：**
| 參數 | 型別 | 說明 |
|---|---|---|
| ctx | context.Context | 請求上下文 |
| student | *Student | 要更新的學生物件 |

**回傳：** `error` - 成功為 nil

**業務規則：**
1. 更新現有的學生記錄
2. 返回錯誤若學生不存在

---

#### Delete()
來源：「系統應該成功刪除該學生」(第 33 line)

**簽名：** `Delete(ctx context.Context, studentNumber string) error`

**參數：**
| 參數 | 型別 | 說明 |
|---|---|---|
| ctx | context.Context | 請求上下文 |
| studentNumber | string | 要刪除的學號 |

**回傳：** `error` - 成功為 nil

**業務規則：**
1. 刪除指定學號的學生記錄
2. 返回錯誤若學生不存在

---

### 4.3 StudentHandler（Gin HTTP Handler）
職責：處理 HTTP 請求，轉換為 Use Case 調用，返回 HTTP 響應

#### CreateStudent()
來源：第 5-10 行

**簽名：** `CreateStudent(c *gin.Context)`

**路由：** `POST /api/students`

**請求體：** `CreateStudentRequest`

**響應：**
- 成功 (201): `Student` 物件
- 失敗 (400): 驗證錯誤
- 失敗 (409): 學號已存在

---

#### GetStudent()
來源：第 12-16 行

**簽名：** `GetStudent(c *gin.Context)`

**路由：** `GET /api/students/:studentNumber`

**參數：** URL 路徑參數 `studentNumber`

**響應：**
- 成功 (200): `Student` 物件
- 失敗 (404): 學生不存在

---

#### GetAllStudents()
來源：第 18-22 行

**簽名：** `GetAllStudents(c *gin.Context)`

**路由：** `GET /api/students`

**響應：**
- 成功 (200): `[]*Student` 陣列
- 成功 (200): 空陣列（若無記錄）

---

#### UpdateStudent()
來源：第 24-28 行

**簽名：** `UpdateStudent(c *gin.Context)`

**路由：** `PUT /api/students/:studentNumber`

**參數：** URL 路徑參數 `studentNumber`

**請求體：** `UpdateStudentRequest`

**響應：**
- 成功 (200): 更新後的 `Student` 物件
- 失敗 (400): 驗證錯誤
- 失敗 (404): 學生不存在

---

#### DeleteStudent()
來源：第 30-34 行

**簽名：** `DeleteStudent(c *gin.Context)`

**路由：** `DELETE /api/students/:studentNumber`

**參數：** URL 路徑參數 `studentNumber`

**響應：**
- 成功 (204): 無內容
- 失敗 (404): 學生不存在

---

## 5. 架構決策

### DDD 使用案例（Use Case）為主的架構

**為何選擇此架構模式：**
- **業務邏輯清晰隔離**：Use Case 層集中所有業務規則（驗證、唯一性檢查等），使業務邏輯不依賴框架
- **可測試性強**：每個 Use Case 可獨立測試，無需依賴 HTTP 框架或資料庫
- **易於維護和擴展**：新增驗證規則或業務邏輯只需修改 Use Case，不影響 Handler 或 Repository
- **框架無關**：業務核心不依賴 Gin，可輕鬆遷移至其他框架

**資料模型設計理由：**
- **Student 實體**：代表業務核心領域，包含所有學生相關資訊
- **CreateStudentRequest / UpdateStudentRequest**：隔離 HTTP 層與 Use Case 層，使用專定的 Request 物件而非直接使用 Student 實體
- **錯誤型別分類**：使用語義明確的錯誤型別（StudentNotFound、StudentNumberAlreadyExists）便於 Handler 映射到正確的 HTTP 狀態碼

**整合方式：**
```
HTTP Request
    ↓
Gin Handler
    ↓
StudentUseCase
    ↓
StudentRepository
    ↓
Database
```

- **Handler** 驗證 HTTP 請求格式，呼叫 Use Case
- **Use Case** 執行所有業務邏輯和驗證
- **Repository** 處理資料持久化和查詢

## 6. 情境對應

| 情境 | 行數 | 核心實體 | Use Case 方法 |
|---|---|---|---|
| 成功新增學生 | 5-10 | Student, CreateStudentRequest | CreateStudent() |
| 查詢單一學生 | 12-16 | Student | GetStudent() |
| 查詢所有學生 | 18-22 | Student | GetAllStudents() |
| 成功更新學生資訊 | 24-28 | Student, UpdateStudentRequest | UpdateStudent() |
| 成功刪除學生 | 30-34 | Student | DeleteStudent() |
| 新增時缺少必填欄位 | 36-40 | CreateStudentRequest | CreateStudent() (驗證) |
| 學號必須唯一 | 42-46 | Student | CreateStudent() (驗證) |
| 電子郵件格式驗證 | 48-52 | CreateStudentRequest | CreateStudent() (驗證) |
| 查詢不存在的學生 | 54-58 | Student | GetStudent() (錯誤) |
| 更新不存在的學生 | 60-64 | Student | UpdateStudent() (錯誤) |
| 刪除不存在的學生 | 66-70 | Student | DeleteStudent() (錯誤) |
| 部分更新學生資訊 | 72-77 | Student, UpdateStudentRequest | UpdateStudent() (部分) |
| 驗證年級範圍 | 79-83 | Student | CreateStudent() (驗證) |

## 7. 檔案結構

```
internal/
├── domain/
│   └── student/
│       ├── student.go                # Student 實體定義
│       ├── request.go                # CreateStudentRequest, UpdateStudentRequest
│       └── error.go                  # StudentError 定義
├── usecase/
│   └── student/
│       └── student_usecase.go        # StudentUseCase 實現
├── repository/
│   └── student/
│       ├── interface.go              # StudentRepository 介面
│       └── memory.go                 # 記憶體實現（測試用）
│       └── database.go               # 資料庫實現
└── handler/
    └── student/
        └── student_handler.go        # Gin HTTP Handler

pkg/
└── logger/                           # 日誌工具

tests/
└── student/
    ├── usecase_test.go              # Use Case 單元測試
    ├── handler_test.go              # Handler 集成測試
    └── repository_test.go           # Repository 單元測試
```

## 8. 技術決策補充

### 驗證策略
- **欄位層級驗證**（Use Case）：格式驗證（電子郵件、年級範圍）
- **業務層級驗證**（Use Case）：唯一性檢查（學號）、必填檢查
- **資料庫層級驗證**：唯一約束（防止競態條件）

### 錯誤處理
- Use Case 層拋出具體的 Domain Error（StudentNotFound、InvalidEmail 等）
- Handler 捕獲 Domain Error 並映射到 HTTP 狀態碼
- 統一的錯誤回應格式：`{ "error": "error_message", "code": "error_code" }`

### 時間戳管理
- 所有學生記錄自動記錄 CreatedAt 和 UpdatedAt
- CreatedAt 建立時設定，之後不可更改
- UpdatedAt 每次更新時自動更新

### 部分更新支援
- UpdateStudentRequest 所有欄位為指標（*string、*int 等）
- 只更新非 nil 的欄位，保留其他欄位不變
