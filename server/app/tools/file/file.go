package file

import (
	"archive/zip"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

func Save(file *multipart.FileHeader, path string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer func() { _ = src.Close() }()

	dst, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() { _ = dst.Close() }()

	_, err = io.Copy(dst, src)
	return err
}

func Mkdir(path string) error {
	return os.MkdirAll(path, os.ModeDir|0775)
}

func Remove(path string) error {
	return os.RemoveAll(path)
}

func Unzip(src, dst string) error {
	reader, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() { _ = reader.Close() }()

	if err := os.MkdirAll(dst, os.ModeDir|0775); err != nil {
		return err
	}

	for _, file := range reader.File {
		path := filepath.Join(dst, file.Name)
		if file.FileInfo().IsDir() {
			_ = os.MkdirAll(path, file.Mode())
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer func() { _ = fileReader.Close() }()

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer func() { _ = targetFile.Close() }()

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			return err
		}
	}
	return nil
}

func Zip(src, dst string) error {
	files := make([]string, 0)
	err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})

	zipfile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() { _ = zipfile.Close() }()

	writer := zip.NewWriter(zipfile)
	defer func() { _ = writer.Close() }()

	for _, x := range files {
		pathInZip, err := filepath.Rel(src, x)
		if err != nil {
			return err
		}

		fdst, err := writer.Create(pathInZip)
		if err != nil {
			return err
		}

		fsrc, err := os.Open(x)
		defer func() { _ = fsrc.Close() }()
		if err != nil {
			return err
		}

		_, err = io.Copy(fdst, fsrc)
		if err != nil {
			return err
		}
	}
	return err
}
