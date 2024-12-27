package appdb

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/specterops/bloodhound/packages/go/apitoy/adapter/appdb/dbmodel"
	"github.com/specterops/bloodhound/packages/go/apitoy/model"
	"github.com/specterops/bloodhound/src/database"
	"github.com/specterops/bloodhound/src/database/types/null"
	"gorm.io/gorm"
)

func (s Adapter) GetFileUploadJobByID(ctx context.Context, jobID int) (model.FileUploadJob, error) {
	if job, err := getFileUploadJob(ctx, s.db, int64(jobID)); errors.Is(err, database.ErrNotFound) {
		return model.FileUploadJob{}, fmt.Errorf("get file upload job by id: %w: %v", model.ErrNotFound, err)
	} else if err != nil {
		return model.FileUploadJob{}, fmt.Errorf("get file upload job by id: %w: %v", model.ErrGenericDatabaseFailure, err)
	} else {
		return convertFileUploadJob(job), nil
	}
}

func (s Adapter) CreateIngestTask(ctx context.Context, filename string, fileType model.FileType, requestID string, jobID int) (model.IngestTask, error) {
	newIngestTask := dbmodel.IngestTask{
		FileName:    filename,
		RequestGUID: requestID,
		TaskID:      null.Int64From(int64(jobID)),
		FileType:    dbmodel.FileType(fileType),
	}

	if task, err := createIngestTask(ctx, s.db, newIngestTask); err != nil {
		return task, fmt.Errorf("create ingest task: %w: %v", model.ErrGenericDatabaseFailure, err)
	} else {
		return task, nil
	}
}

func (s Adapter) TouchFileUploadJobLastIngest(ctx context.Context, fileUploadJob model.FileUploadJob) error {
	fileUploadJob.LastIngest = time.Now().UTC()
	if err := updateFileUploadJob(ctx, s.db, fileUploadJob); err != nil {
		return fmt.Errorf("touch last ingest: %w: %v", model.ErrGenericDatabaseFailure, err)
	} else {
		return nil
	}
}

func getFileUploadJob(ctx context.Context, db *gorm.DB, jobID int64) (dbmodel.FileUploadJob, error) {
	var job dbmodel.FileUploadJob
	if result := db.Preload("User").WithContext(ctx).First(&job, jobID); result.Error != nil {
		return job, checkError(result)
	} else {
		return job, nil
	}
}

func createIngestTask(ctx context.Context, db *gorm.DB, ingestTask dbmodel.IngestTask) (dbmodel.IngestTask, error) {
	result := db.WithContext(ctx).Create(&ingestTask)

	return ingestTask, checkError(result)
}

func updateFileUploadJob(ctx context.Context, db *gorm.DB, fileUploadJob dbmodel.FileUploadJob) error {
	result := db.WithContext(ctx).Save(&fileUploadJob)
	return checkError(result)
}

func convertFileUploadJob(job dbmodel.FileUploadJob) model.FileUploadJob {
	return model.FileUploadJob{}
}
