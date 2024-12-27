package app

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"io"

	"github.com/specterops/bloodhound/log"
	"github.com/specterops/bloodhound/packages/go/apitoy/model"
)

// IngestFile ingests a given file in the form of an io.ReadCloser. It also requires a valid context, requestID, jobID, and fileType.
func (s BHApp) IngestFile(ctx context.Context, requestID string, jobID int, fileType model.FileType, content io.ReadCloser) error {
	var validationStrategy model.FileValidator

	switch fileType {
	case model.FileTypeJson:
		validationStrategy = writeAndValidateJSON
	case model.FileTypeZip:
		validationStrategy = writeAndValidateZip
	default:
		return model.ErrInvalidFile
	}

	if fileUploadJob, err := s.dbAdapter.GetFileUploadJobByID(ctx, jobID); err != nil {
		return err
	} else if tempFile, err := s.fileAdapter.SaveIngestFile(content, validationStrategy); err != nil {
		return err
	} else if _, err := s.dbAdapter.CreateIngestTask(ctx, tempFile, fileType, requestID, jobID); err != nil {
		return err
	} else if err := s.dbAdapter.TouchFileUploadJobLastIngest(ctx, fileUploadJob); err != nil {
		return err
	} else {
		return nil
	}
}

var zipMagicBytes = []byte{0x50, 0x4b, 0x03, 0x04}

// validateMetaTag ensures that the correct tags are present in a json file for data model.
// If readToEnd is set to true, the stream will read to the end of the file (needed for TeeReader)
func validateMetaTag(reader io.Reader, readToEnd bool) (model.Metadata, error) {
	var (
		depth            = 0
		decoder          = json.NewDecoder(reader)
		dataTagFound     = false
		dataTagValidated = false
		metaTagFound     = false
		meta             model.Metadata
	)

	for {
		if token, err := decoder.Token(); err != nil {
			if errors.Is(err, io.EOF) {
				if !metaTagFound && !dataTagFound {
					return model.Metadata{}, model.ErrNoTagFound
				} else if !dataTagFound {
					return model.Metadata{}, model.ErrDataTagNotFound
				} else {
					return model.Metadata{}, model.ErrMetaTagNotFound
				}
			} else {
				return model.Metadata{}, model.ErrInvalidJSONFile
			}
		} else {
			//Validate that our data tag is actually opening correctly
			if dataTagFound && !dataTagValidated {
				if typed, ok := token.(json.Delim); ok && typed == model.DelimOpenSquareBracket {
					dataTagValidated = true
				} else {
					dataTagFound = false
				}
			}
			switch typed := token.(type) {
			case json.Delim:
				switch typed {
				case model.DelimCloseBracket, model.DelimCloseSquareBracket:
					depth--
				case model.DelimOpenBracket, model.DelimOpenSquareBracket:
					depth++
				}
			case string:
				if !metaTagFound && depth == 1 && typed == "meta" {
					if err := decoder.Decode(&meta); err != nil {
						log.Warnf("Found invalid metatag, skipping")
					} else if meta.Type.IsValid() {
						metaTagFound = true
					}
				}

				if !dataTagFound && depth == 1 && typed == "data" {
					dataTagFound = true
				}
			}
		}

		if dataTagValidated && metaTagFound {
			break
		}
	}

	if readToEnd {
		if _, err := io.Copy(io.Discard, reader); err != nil {
			return model.Metadata{}, err
		}
	}

	return meta, nil
}

func validateZipFile(reader io.Reader) error {
	bytes := make([]byte, 4)
	if readBytes, err := reader.Read(bytes); err != nil {
		return err
	} else if readBytes < 4 {
		return model.ErrInvalidZipFile
	} else {
		for i := 0; i < 4; i++ {
			if bytes[i] != zipMagicBytes[i] {
				return model.ErrInvalidZipFile
			}
		}

		_, err := io.Copy(io.Discard, reader)

		return err
	}
}

func writeAndValidateZip(src io.Reader, dst io.Writer) error {
	tr := io.TeeReader(src, dst)
	return validateZipFile(tr)
}

const (
	UTF8BOM1 = 0xef
	UTF8BOM2 = 0xbb
	UTF8BMO3 = 0xbf
)

func writeAndValidateJSON(src io.Reader, dst io.Writer) error {
	tr := io.TeeReader(src, dst)
	bufReader := bufio.NewReader(tr)
	if b, err := bufReader.Peek(3); err != nil {
		return err
	} else {
		if b[0] == UTF8BOM1 && b[1] == UTF8BOM2 && b[2] == UTF8BMO3 {
			if _, err := bufReader.Discard(3); err != nil {
				return err
			}
		}
	}
	_, err := validateMetaTag(bufReader, true)
	return err
}
