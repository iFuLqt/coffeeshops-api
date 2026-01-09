package entity

import "mime/multipart"

type FileUploadImageEntity struct {
	File *multipart.FileHeader
	Name string
}