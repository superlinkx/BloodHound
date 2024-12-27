package dbmodel

import "github.com/specterops/bloodhound/src/database/types/null"

type IngestTask struct {
	FileName    string
	RequestGUID string
	TaskID      null.Int64
	FileType    FileType

	BigSerial
}

type IngestTasks []IngestTask

type FileType int

const (
	FileTypeInvalid FileType = iota
	FileTypeJson
	FileTypeZip
)
