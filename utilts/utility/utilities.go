package utility

import (
	"io"
	"mime/multipart"
	"os"
)

func SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	// ทำการ copy ไฟล์จาก src ไปยัง out
	if _, err = io.Copy(out, src); err != nil {
		return err
	}

	return nil
}
