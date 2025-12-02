# Student CRUD API

> 本項目使用 **[GSI Protocol](https://github.com/CodeMachine0121/GSI-Protocol)** 框架進行開發

## 項目概述

這是一個學校系統的學生管理 API，提供學生資訊的建立、查看、更新和刪除功能。作為學校系統管理員，可以透過此 API 有效地管理學生記錄。

## 技術棧

- **後端框架**: Go + Gin
- **測試框架**: Testify
- **專案架構**: 用例驅動的領域驅動設計 (DDD)
- **開發方法論**: GSI Protocol (Specification-Architecture-Implementation-Verification)

## 專案結構

```
todo-list/
├── docs/                    # 文檔
│   └── features/           # 功能特性文檔
│       └── student_crud_api/
│           └── architecture.md
├── features/               # Gherkin 行為規格
│   └── student_crud_api.feature
├── internal/               # 內部實現
│   ├── domain/            # 領域層 (實體、值物件)
│   │   └── student/
│   ├── handler/           # 應用層 (控制器)
│   │   └── student/
│   ├── repository/        # 基礎設施層 (資料存取)
│   │   └── student/
│   └── usecase/           # 業務邏輯層 (用例)
│       └── student/
├── go.mod                  # Go 模組定義
└── go.sum                  # 依賴鎖定檔
```

## 主要功能

### API 端點

| 方法 | 端點 | 功能 |
|------|------|------|
| POST | `/students` | 新增學生 |
| GET | `/students/:id` | 查詢單一學生 |
| GET | `/students` | 查詢所有學生 |
| PATCH | `/students/:id` | 更新學生資訊 |
| DELETE | `/students/:id` | 刪除學生 |

### 學生資訊結構

- **姓名** (name): 必填
- **學號** (student_id): 必填、唯一
- **電子郵件** (email): 必填、格式驗證
- **班級** (class): 可選
- **年級** (grade): 可選 (1-6 之間)

## 開發流程

本項目遵循 GSI Protocol 的完整開發流程：

1. **Specification (規格)** - 使用 Gherkin 定義功能需求
2. **Architecture (架構)** - 設計系統架構和資料模型
3. **Implementation (實現)** - 根據規格實現程式碼
4. **Verification (驗證)** - 執行測試驗證功能正確性

## 快速開始

### 安裝依賴

```bash
go mod download
```

### 執行測試

```bash
go test ./...
```

### 運行應用

```bash
go run main.go
```

## 驗證規則

- ✓ 學號必須唯一
- ✓ 電子郵件必須符合有效格式
- ✓ 年級必須在 1-6 之間
- ✓ 姓名為必填欄位
- ✓ 支援部分更新 (PATCH)

## 錯誤處理

API 返回標準化的錯誤回應：

- `400 Bad Request` - 請求資料驗證失敗
- `404 Not Found` - 學生不存在
- `409 Conflict` - 學號已存在
- `500 Internal Server Error` - 伺服器錯誤

## 開發參考

- **GSI Protocol**: https://github.com/CodeMachine0121/GSI-Protocol
- **Gin 文檔**: https://gin-gonic.com/
- **Testify 文檔**: https://github.com/stretchr/testify

## License

MIT
