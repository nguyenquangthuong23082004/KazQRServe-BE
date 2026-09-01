package utils

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/skip2/go-qrcode"
)

// GenerateQRCodeImage sinh file ảnh QR Code PNG từ chuỗi content và lưu vào ./uploads/qrcodes/
// Trả về đường dẫn tương đối (ví dụ: "/uploads/qrcodes/table_e8f9a2b1-c4d3.png")
func GenerateQRCodeImage(content string, fileName string) (string, error) {
	subFolder := "qrcodes"
	targetDir := filepath.Join(UploadsRootDir, subFolder)
	if err := os.MkdirAll(targetDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("không thể tạo thư mục lưu mã QR: %v", err)
	}

	fullName := fmt.Sprintf("%s.png", fileName)
	targetFilePath := filepath.Join(targetDir, fullName)

	// Sinh file ảnh PNG QR với kích thước 256x256 pixel, mức sửa lỗi Medium
	if err := qrcode.WriteFile(content, qrcode.Medium, 256, targetFilePath); err != nil {
		return "", fmt.Errorf("lỗi sinh ảnh QR Code: %v", err)
	}

	relativePath := fmt.Sprintf("/uploads/%s/%s", subFolder, fullName)
	return relativePath, nil
}
