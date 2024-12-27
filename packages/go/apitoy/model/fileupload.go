package model

import (
	"io"
	"time"

	"github.com/gofrs/uuid"
)

type FileValidator func(src io.Reader, dst io.Writer) error

type IngestTask struct {
	FileName    string
	RequestGUID string
	TaskID      int64
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

type FileUploadJob struct {
	UserID           uuid.UUID
	UserEmailAddress string
	User             User
	Status           JobStatus
	StatusMessage    string
	StartTime        time.Time
	EndTime          time.Time
	LastIngest       time.Time
	TotalFiles       int
	FailedFiles      int

	BigSerial
}

type FileUploadJobs []FileUploadJob

type JobStatus int

const (
	JobStatusInvalid           JobStatus = -1
	JobStatusReady             JobStatus = 0
	JobStatusRunning           JobStatus = 1
	JobStatusComplete          JobStatus = 2
	JobStatusCanceled          JobStatus = 3
	JobStatusTimedOut          JobStatus = 4
	JobStatusFailed            JobStatus = 5
	JobStatusIngesting         JobStatus = 6
	JobStatusAnalyzing         JobStatus = 7
	JobStatusPartiallyComplete JobStatus = 8
)
