# Hướng dẫn Tự động tạo bảng (Migration) & Quy tắc tổ chức Thư mục Model

Tài liệu này hướng dẫn cách thiết lập tự động tạo bảng dữ liệu khi khởi chạy ứng dụng Go (sử dụng GORM) và giải thích chi tiết tư duy thiết kế cấu trúc thư mục chứa các Model/Entity của bạn.

---

## 1. Thiết lập Tự động tạo bảng (AutoMigrate) khi khởi chạy App

Để ứng dụng tự động đồng bộ cấu trúc từ code Go (Struct) vào PostgreSQL khi khởi chạy, chúng ta thực hiện theo quy trình 3 bước sau:

### Bước 1: Định nghĩa Model Struct
Các struct biểu diễn cấu trúc bảng dữ liệu sẽ được đặt trong thư mục của từng module tương ứng (Ví dụ: `User` thuộc module `user`).
* File: `internal/modules/user/entity/user.go`

```go
package entity

import (
	"time"

	"gorm.io/gorm"
)

// User đại diện cho bảng "users" trong Database
type User struct {
	ID        uint           `gorm:"primaryKey;autoIncrement"`
	Username  string         `gorm:"type:varchar(100);unique;not null"`
	Email     string         `gorm:"type:varchar(150);unique;not null"`
	Password  string         `gorm:"type:varchar(255);not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"` // Hỗ trợ Soft Delete (Xóa tạm thời)
}
```

### Bước 2: Tạo hàm điều phối Migration
Tạo một file quản lý migration chung để đăng ký tất cả các bảng dữ liệu sẽ được đồng bộ.
* File: [migration.go](file:///home/uwal/Desktop/KazQRServe-BE/internal/shared/database/migration.go)

```go
package database

import (
	"log"

	"github.com/nguyenquangthuong23082004/KazQRServe-BE/internal/modules/user/entity"
	"gorm.io/gorm"
)

// AutoMigrate nhận vào kết nối DB và tự động đồng bộ hóa các struct thành bảng dữ liệu
func AutoMigrate(db *gorm.DB) error {
	log.Println("Starting database auto-migration...")

	// Khai báo tất cả các model muốn tạo bảng ở đây
	err := db.AutoMigrate(
		&entity.User{},
		// &entity.QRCode{}, // Đăng ký thêm các bảng khác khi phát triển
	)

	if err != nil {
		return err
	}

	log.Println("Database migration completed successfully!")
	return nil
}
```

### Bước 3: Kích hoạt Migration trong main.go
Gọi hàm migration ngay sau khi kết nối cơ sở dữ liệu thành công tại entrypoint của ứng dụng.
* File: [main.go](file:///home/uwal/Desktop/KazQRServe-BE/cmd/api/main.go)

```go
	// Khởi tạo kết nối databases
	db, err := database.InitDb()
	if err != nil {
		log.Fatalln("Error initializing database: ", err)
	}

	// Tự động tạo hoặc nâng cấp cấu trúc bảng
	if err := database.AutoMigrate(db); err != nil {
		log.Fatalln("Error running auto migration: ", err)
	}

	// Truyền DB vừa kết nối vào App
	application := app.NewApp(db)
```

---

## 2. Tại sao lại chia thư mục dạng `/modules/user/entity`?

Phương pháp tổ chức thư mục này tuân theo triết lý **Modular Design** kết hợp với **Clean Architecture**, mang lại nhiều lợi ích thực tế:

* **Tính đóng gói cao (Encapsulation):** Tất cả các thành phần thuộc một nghiệp vụ cụ thể (như entity, repository, service, handler của `User`) sẽ nằm trọn vẹn trong thư mục `/modules/user`. Việc phát triển, sửa lỗi hay nâng cấp chức năng liên quan đến user sẽ được khoanh vùng tại một chỗ.
* **Tránh lỗi Circular Dependency (Vòng lặp import trong Go):** Go là ngôn ngữ nghiêm ngặt không cho phép Package A import Package B và ngược lại Package B cũng import Package A. Bằng cách thiết kế theo từng Module độc lập, các luồng import sẽ đi từ ngoài vào trong một chiều nhất quán, giảm thiểu lỗi biên dịch.
* **Định nghĩa rõ vai trò của Entity:** `entity` chỉ là các cấu trúc dữ liệu mô tả cơ sở dữ liệu và quy luật nghiệp vụ lõi của module đó.

---

## 3. Vai trò của thư mục `shared/model` (hoặc `shared`)

Thư mục `shared` dùng để chứa các cấu trúc dữ liệu không thuộc về bất kỳ một nghiệp vụ/module cụ thể nào, mà được dùng chung trên toàn hệ thống.

* **Base Model:** Dùng để tái sử dụng các thuộc tính chung của database.
  * *Ví dụ:* Tạo một `BaseModel` chứa các trường audit logs mặc định:
    ```go
    package model

    import "time"

    type BaseModel struct {
        ID        uint      `gorm:"primaryKey"`
        CreatedAt time.Time
        UpdatedAt time.Time
    }
    ```
  * Sau đó ở các Model Entity, ta chỉ cần nhúng `BaseModel` của `shared` vào:
    ```go
    type User struct {
        model.BaseModel // Kế thừa ID, CreatedAt, UpdateAt
        Username string
    }
    ```
* **Common DTOs (Data Transfer Objects):** Các cấu trúc dùng để giao tiếp dữ liệu giữa Client và Server mà hầu như API nào cũng cần.
  * *Ví dụ:* Cấu trúc phân trang chung (`PaginationRequest`, `PaginationResponse`), Cấu trúc phản hồi API chuẩn (`StandardResponse`).

---

## 4. Tóm tắt nguyên tắc phân chia Model

Để giữ hệ thống gọn gàng, hãy phân loại Model theo bảng so sánh sau:

| Tiêu chí | Đặt tại `modules/{name}/entity` | Đặt tại `shared` |
| :--- | :--- | :--- |
| **Phạm vi** | Chỉ thuộc về nghiệp vụ cụ thể đó. | Dùng chung cho toàn hệ thống. |
| **Ví dụ** | `User`, `QRCode`, `Transaction`, `History`. | `BaseModel`, `StandardResponse`, `Pagination`. |
| **Liên kết** | Có thể liên kết với các entity của module khác qua khoá ngoại. | Hoàn toàn độc lập, không tham chiếu trực tiếp đến bất kỳ thực thể nghiệp vụ cụ thể nào. |
| **Khuyên dùng** | **Lựa chọn mặc định** cho mọi bảng dữ liệu chính của ứng dụng. | **Hạn chế tối đa**, chỉ chứa các cấu trúc dữ liệu nền tảng. |
