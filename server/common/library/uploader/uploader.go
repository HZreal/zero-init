package uploader

import (
	"github.com/zeromicro/go-zero/core/logx"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	fileTool "overall/common/utils/file"
	TimeTools "overall/common/utils/time"
	"overall/common/xerr"
	"path/filepath"
	"time"
)

const (
	// uploadFailed       = "file upload failed"
	fileRetrieveFailed = "failed to retrieve file err: %v"
	fileSizeLimit      = "%v file size %v exceeds the limit %v"
	fileTypeNotAllow   = "%v file type %v is not allow"
	fileTypeForbidden  = "%v file type %v is forbidden"
	fileTypeInvalid    = "%v file type %v is invalid, real type is %v"
	filePathInvalid    = "%v file path %v is invalid, err: %v"
	fileCreateFailed   = "%v file create err: %v"
	fileCopyError      = "%v file copy err: %v"

	MaxSize = 50 << 20 // 50MB

	// 规约 表单文件字段为 file
	// FormKey = "file"
)

var (
	UploadFailedError             = xerr.NewErrCode(xerr.UploadFailed)
	UploadFileTooLargeError       = xerr.NewErrCode(xerr.UploadFileTooLarge)
	UploadFileTypeNotAllowedError = xerr.NewErrCode(xerr.UploadFileTypeNotAllowed)

	AllowAllType = map[string]string{
		".xls":  "application/vnd.ms-excel",
		".xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		".txt":  "text/plain",
		".csv":  "text/csv",
		".png":  "image/png",
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".gif":  "image/gif",
	}
)

type Uploader struct {
	Request     *http.Request
	AllowType   []string
	FormKey     []string
	TmpStoreDir string
	Prefix      string
	Now         time.Time
}

type InitUploaderReq struct {
	Request     *http.Request
	AllowType   []string
	FormKeyList []string
	TmpStoreDir string
	Prefix      string
	Now         time.Time
}

type FileMetadata struct {
	FileName     string
	TmpStorePath string
	FileType     string
}

func NewUploader(initUploaderReq *InitUploaderReq) *Uploader {
	return &Uploader{
		Request:     initUploaderReq.Request,
		AllowType:   initUploaderReq.AllowType,
		FormKey:     initUploaderReq.FormKeyList,
		TmpStoreDir: initUploaderReq.TmpStoreDir,
		Prefix:      initUploaderReq.Prefix,
		Now:         initUploaderReq.Now,
	}
}

func (u *Uploader) UploadHandle() ([]*FileMetadata, error) {
	metadata := make([]*FileMetadata, 0)

	// 必须要先调用一下 ParseMultipartForm 或者 request.FormFile(formKey) 否则 Request.MultipartForm 空指针异常
	err := u.Request.ParseMultipartForm(MaxSize)
	if err != nil {
		return nil, err
	}

	for _, formKey := range u.FormKey {
		// 一个文件类型字段包含多个文件
		fileHeaders, ok := u.Request.MultipartForm.File[formKey]
		if !ok {
			return nil, UploadFailedError
		}
		for _, handler := range fileHeaders {
			file, err := handler.Open()
			if err != nil {
				return nil, err
			}
			defer file.Close()

			//
			err = u.CheckFile(handler)
			if err != nil {
				return nil, err
			}

			//
			err = u.CheckTmpStoreDir()
			if err != nil {
				logx.Errorf(filePathInvalid, handler.Filename, u.TmpStoreDir, err)
				return nil, UploadFailedError
			}

			//
			now := TimeTools.Format(TimeTools.TIME_FORMAT_YYYYMMDDHHMMSS, u.Now)
			ext := filepath.Ext(handler.Filename)
			newFileName := fileTool.GetRandName(u.Prefix, now) + ext
			tmpStorePath := filepath.Join(u.TmpStoreDir, newFileName)

			//

			//
			err = u.WriteFile(file, tmpStorePath)
			if err != nil {
				return nil, err
			}

			metadata = append(metadata, &FileMetadata{
				FileName:     newFileName,
				FileType:     ext,
				TmpStorePath: tmpStorePath,
			})

		}

	}

	return metadata, nil
}

func (u *Uploader) CheckFile(handler *multipart.FileHeader) error {
	fileName := handler.Filename
	ext := filepath.Ext(fileName)
	fileType := ext

	if handler.Size > MaxSize {
		logx.Errorf(fileSizeLimit, fileName, handler.Size, MaxSize)
		return UploadFileTooLargeError
	}

	if len(u.AllowType) > 0 {
		flag := checkAllowType(u.AllowType, fileType)
		if !flag {
			logx.Errorf(fileTypeNotAllow, fileName, fileType)
			return UploadFileTypeNotAllowedError
		}

		for _, v := range u.AllowType {
			if _, ok := AllowAllType[v]; !ok {
				logx.Errorf(fileTypeForbidden, fileName, v)
				return UploadFileTypeNotAllowedError
			}
		}
	}

	mineType := handler.Header.Get("Content-Type")
	if AllowAllType[fileType] != mineType {
		logx.Errorf(fileTypeInvalid, fileName, AllowAllType[fileType], mineType)
		return UploadFailedError
	}

	return nil
}

func (u *Uploader) CheckTmpStoreDir() error {
	err := fileTool.CheckFileDir(u.TmpStoreDir)
	if err != nil {
		return err
	}
	return nil

}

func (u *Uploader) WriteFile(file multipart.File, tmpStorePath string) error {
	// 写入临时文件
	out, err := os.Create(tmpStorePath)
	if err != nil {
		logx.Errorf(fileCreateFailed, tmpStorePath, err)
		return UploadFailedError
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		logx.Errorf(fileCopyError, tmpStorePath, err)
		return UploadFailedError
	}
	return nil
}

func checkAllowType(items []string, item string) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}
