package main

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

func saveImage(fileHeader *multipart.FileHeader) (string, error) {

	src, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	err = os.MkdirAll("uploads", 0755)
	if err != nil {
		return "", err
	}

	filename := uuid.New().String() + filepath.Ext(fileHeader.Filename)

	dstPath := filepath.Join("uploads", filename)

	dst, err := os.Create(dstPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return "", err
	}

	return "/uploads/" + filename, nil
}
