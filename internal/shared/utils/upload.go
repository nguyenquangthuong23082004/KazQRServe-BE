package utils

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Cấu hình giới hạn upload
const (
	MaxFileSize    = 2 * 1024 * 1024 // 2MB
	UploadsRootDir = "./uploads"
)

// AllowedExtensions danh sách các định dạng ảnh được phép
var AllowedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".webp": true,
}

// SaveUploadedFile nhận file từ form-data, kiểm tra tính hợp lệ và lưu vào thư mục phân loại.
// Trả về đường dẫn tương đối (ví dụ: "/uploads/products/1712345678_abc.jpg") hoặc lỗi.
func SaveUploadedFile(file *multipart.FileHeader, subFolder string) (string, error) {
	// 1. Kiểm tra kích thước file
	if file.Size > MaxFileSize {
		return "", fmt.Errorf("kích thước file vượt quá giới hạn cho phép (tối đa %d MB)", MaxFileSize/(1024*1024))
	}

	// 2. Kiểm tra phần mở rộng của file
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !AllowedExtensions[ext] {
		return "", errors.New("định dạng file không hợp lệ (chỉ chấp nhận .jpg, .jpeg, .png, .webp)")
	}

	// 3. Tạo thư mục lưu trữ đích (ví dụ: ./uploads/products)
	targetDir := filepath.Join(UploadsRootDir, subFolder)
	if err := os.MkdirAll(targetDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("không thể tạo thư mục lưu trữ: %v", err)
	}

	// 4. Sinh tên file duy nhất sử dụng timestamp nano-giây để tránh trùng tên
	timestamp := time.Now().UnixNano()
	uniqueName := fmt.Sprintf("%d%s", timestamp, ext)
	targetFilePath := filepath.Join(targetDir, uniqueName)

	// 5. Mở file nguồn
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// 6. Tạo file đích để lưu
	out, err := os.Create(targetFilePath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// 7. Sao chép nội dung file
	if _, err = io.Copy(out, src); err != nil {
		return "", err
	}

	// 8. Trả về đường dẫn tương đối sử dụng định dạng slash cho URL tĩnh
	relativePath := fmt.Sprintf("/uploads/%s/%s", subFolder, uniqueName)
	return relativePath, nil
}
