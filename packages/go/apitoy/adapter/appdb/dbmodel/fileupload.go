package dbmodel

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/specterops/bloodhound/src/database/types/null"
)

type FileUploadJob struct {
	UserID           uuid.UUID
	UserEmailAddress null.String
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

type DomainCollectionResult struct {
	JobID             int64
	DomainName        string
	Success           bool
	Message           string
	UserCount         int
	GroupCount        int
	ComputerCount     int
	GPOCount          int
	OUCount           int
	ContainerCount    int
	AIACACount        int `gorm:"column:aiaca_count"`
	RootCACount       int `gorm:"column:rootca_count"`
	EnterpriseCACount int `gorm:"column:enterpriseca_count"`
	NTAuthStoreCount  int `gorm:"column:ntauthstore_count"`
	CertTemplateCount int `gorm:"column:certtemplate_count"`
	DeletedCount      int

	BigSerial
}
