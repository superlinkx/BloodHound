package dbmodel

import (
	"github.com/specterops/bloodhound/packages/go/apitoy/model"
	"github.com/specterops/bloodhound/src/database/types/null"
)

type IngestTask struct {
	FileName    string
	RequestGUID string
	TaskID      null.Int64
	FileType    FileType

	BigSerial
}

func (s IngestTask) ConvertToAppModel() model.IngestTask {
	return model.IngestTask{
		FileName:    s.FileName,
		RequestGUID: s.RequestGUID,
		TaskID:      s.TaskID.Int64,
		FileType:    model.FileType(s.FileType),
		BigSerial:   s.BigSerial.ConvertToAppModel(),
	}
}

type IngestTasks []IngestTask

type FileType int

const (
	FileTypeInvalid FileType = iota
	FileTypeJson
	FileTypeZip
)
