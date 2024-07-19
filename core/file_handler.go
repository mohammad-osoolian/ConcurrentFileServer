package core

import (
	"ConcurrentFileServer/utils"
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type FileHandler interface {
	UploadFile(ctx context.Context, file []byte, mimeType string) (string, error)
	DownloadFile(ctx context.Context, fileID string) ([]byte, string, error)
}

type FileHandlerImpl struct {
	RootDir   string
	FileCount int64
}

func NewFileHandlerImpl() FileHandler {
	return &FileHandlerImpl{RootDir: "../files/", FileCount: 0}
}

func (f *FileHandlerImpl) UploadFile(ctx context.Context, file []byte, mimeType string) (string, error) {
	//TODO implement me
	f.FileCount++ //wanted to use for generating access_hash but changed my mind. However it seemed useful so i kept it.

	access_hash := fmt.Sprintf("%.19d", rand.Int63())
	md5hash := md5.Sum(file)
	file_name := hex.EncodeToString(md5hash[:])
	file_id := fmt.Sprintf("%s-%s", access_hash, file_name)
	extention := utils.GetExtensionByMimeType(mimeType)
	path := fmt.Sprintf("%s%s.%s", f.RootDir, file_id, extention)
	if ctx.Err() != nil {
		return "", errors.New("error: request cancelled")
	}
	err := ConcurrentWrite(path, file, 1024*1024)
	return file_id, err
}

func (f *FileHandlerImpl) DownloadFile(ctx context.Context, fileID string) ([]byte, string, error) {
	//TODO implement me
	var extention string
	filepath.Walk(f.RootDir, func(path string, info fs.FileInfo, err error) error {
		if (!info.IsDir()) && (strings.Split(info.Name(), ".")[0] == fileID) {
			extention = strings.Split(info.Name(), ".")[1]
			return nil
		}
		return nil
	})
	if extention == "" {
		return nil, "", errors.New("error: file id not found")
	}

	mime_type := utils.GetMimeTypeByExtension(extention)
	path := fmt.Sprintf("%s%s.%s", f.RootDir, fileID, extention)
	if ctx.Err() != nil {
		return nil, "", errors.New("error: request cancelled")
	}
	file, err := ConcurrentRead(path, 1024*1024)
	return *file, mime_type, err
}

func ConcurrentRead(filepath string, chunkSize int) (*[]byte, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	fileInfo, _ := file.Stat()
	fileSize := fileInfo.Size()

	data := make([]byte, fileSize)
	var wg sync.WaitGroup
	for offset := int64(0); offset < fileSize; offset += int64(chunkSize) {
		wg.Add(1)
		if offset+int64(chunkSize) > fileSize {
			chunkSize = int(fileSize - offset)
		}
		go func() {
			file.ReadAt(data[offset:offset+int64(chunkSize)], offset)
			wg.Done()
		}()
	}
	wg.Wait()
	return &data, nil
}

func ConcurrentWrite(filepath string, data []byte, chunkSize int) error {
	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	fileSize := int64(len(data))
	var wg sync.WaitGroup
	for offset := int64(0); offset < fileSize; offset += int64(chunkSize) {
		wg.Add(1)
		if offset+int64(chunkSize) > fileSize {
			chunkSize = int(fileSize - offset)
		}
		go func() {
			file.WriteAt(data[offset:offset+int64(chunkSize)], offset)
			wg.Done()
		}()
	}
	wg.Wait()
	return nil
}

func SequentialWrite(filepath string, data []byte) error {
	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func SequentialRead(filepath string) (*[]byte, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return &data, nil
}
