package file

import (
	"fmt"
	"io"
	"os"

	"github.com/specterops/bloodhound/log"
	"github.com/specterops/bloodhound/packages/go/apitoy/model"
)

func (s Adapter) SaveIngestFile(body io.ReadCloser, validationStrategy model.FileValidator) (string, error) {
	tempFile, err := os.CreateTemp(s.tempDir, s.ingestFilePrefix)
	if err != nil {
		return "", fmt.Errorf("creating ingest file: %w: %v", model.ErrGeneralApplicationFailure, err)
	}

	if err := validationStrategy(body, tempFile); err != nil {
		if err := tempFile.Close(); err != nil {
			log.Errorf("Error closing temp file %s with failed validation: %v", tempFile.Name(), err)
		} else if err := os.Remove(tempFile.Name()); err != nil {
			log.Errorf("Error deleting temp file %s: %v", tempFile.Name(), err)
		}
		return tempFile.Name(), fmt.Errorf("saving ingest file: %w: %v", model.ErrInvalidFile, err)
	} else {
		if err := tempFile.Close(); err != nil {
			log.Errorf("Error closing temp file with successful validation %s: %v", tempFile.Name(), err)
		}
		return tempFile.Name(), nil
	}
}
