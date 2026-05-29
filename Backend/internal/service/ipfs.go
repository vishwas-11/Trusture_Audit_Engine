package service

type IPFSUploader interface {
	UploadFile(filePath string) (string, error)
}
