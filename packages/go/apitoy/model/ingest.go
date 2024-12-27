package model

import (
	"encoding/json"
	"errors"

	"github.com/specterops/bloodhound/mediatypes"
)

var AllowedZipFileUploadTypes = []string{
	mediatypes.ApplicationZip.String(),
	"application/x-zip-compressed", // Not currently available in mediatypes
	"application/zip-compressed",   // Not currently available in mediatypes
}

var (
	ErrMetaTagNotFound     = errors.New("no valid meta tag found")
	ErrDataTagNotFound     = errors.New("no data tag found")
	ErrNoTagFound          = errors.New("no valid meta tag or data tag found")
	ErrInvalidDataTag      = errors.New("invalid data tag found")
	ErrJSONDecoderInternal = errors.New("json decoder internal error")
	ErrInvalidZipFile      = errors.New("failed to find zip file header")
)

type Metadata struct {
	Type    DataType         `json:"type"`
	Methods CollectionMethod `json:"methods"`
	Version int              `json:"version"`
}

type DataType string

const (
	DataTypeSession        DataType = "sessions"
	DataTypeUser           DataType = "users"
	DataTypeGroup          DataType = "groups"
	DataTypeComputer       DataType = "computers"
	DataTypeGPO            DataType = "gpos"
	DataTypeOU             DataType = "ous"
	DataTypeDomain         DataType = "domains"
	DataTypeRemoved        DataType = "deleted"
	DataTypeContainer      DataType = "containers"
	DataTypeLocalGroups    DataType = "localgroups"
	DataTypeAIACA          DataType = "aiacas"
	DataTypeRootCA         DataType = "rootcas"
	DataTypeEnterpriseCA   DataType = "enterprisecas"
	DataTypeNTAuthStore    DataType = "ntauthstores"
	DataTypeCertTemplate   DataType = "certtemplates"
	DataTypeAzure          DataType = "azure"
	DataTypeIssuancePolicy DataType = "issuancepolicies"
)

func AllIngestDataTypes() []DataType {
	return []DataType{
		DataTypeSession,
		DataTypeUser,
		DataTypeGroup,
		DataTypeComputer,
		DataTypeGPO,
		DataTypeOU,
		DataTypeDomain,
		DataTypeRemoved,
		DataTypeContainer,
		DataTypeLocalGroups,
		DataTypeAIACA,
		DataTypeRootCA,
		DataTypeEnterpriseCA,
		DataTypeNTAuthStore,
		DataTypeCertTemplate,
		DataTypeAzure,
		DataTypeIssuancePolicy,
	}
}

func (s DataType) IsValid() bool {
	for _, method := range AllIngestDataTypes() {
		if s == method {
			return true
		}
	}

	return false
}

type CollectionMethod uint64

const (
	CollectionMethodGroup CollectionMethod = 1 << iota
	CollectionMethodLocalAdmin
	CollectionMethodGPOLocalGroup
	CollectionMethodSession
	CollectionMethodLoggedOn
	CollectionMethodTrusts
	CollectionMethodACL
	CollectionMethodContainer
	CollectionMethodRDP
	CollectionMethodObjectProps
	CollectionMethodSessionLoop
	CollectionMethodLoggedOnLoop
	CollectionMethodDCOM
	CollectionMethodSPNTargets
	CollectionMethodPSRemote
	CollectionMethodUserRights
	CollectionMethodCARegistry
	CollectionMethodDCRegistry
	CollectionMethodCertServices
)

const (
	DelimOpenBracket        = json.Delim('{')
	DelimCloseBracket       = json.Delim('}')
	DelimOpenSquareBracket  = json.Delim('[')
	DelimCloseSquareBracket = json.Delim(']')
)
