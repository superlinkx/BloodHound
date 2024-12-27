package dbmodel

import (
	"errors"
	"time"

	"github.com/gofrs/uuid"
	"github.com/specterops/bloodhound/src/database/types"
	"github.com/specterops/bloodhound/src/database/types/null"
	"github.com/specterops/bloodhound/src/serde"
)

const PermissionURIScheme = "permission"

type Installation struct {
	Unique
}

type Permission struct {
	Authority string
	Name      string

	Serial
}

type Permissions []Permission

type AuthToken struct {
	UserID     uuid.NullUUID `gorm:"type:text"`
	ClientID   uuid.NullUUID `gorm:"type:text"`
	Name       null.String
	Key        string
	HmacMethod string
	LastAccess time.Time

	Unique
}

type AuthTokens []AuthToken

type AuthSecret struct {
	UserID        uuid.UUID
	Digest        string
	DigestMethod  string
	ExpiresAt     time.Time
	TOTPSecret    string
	TOTPActivated bool

	Serial
}

func RoleAssociations() []string {
	return []string{
		"Permissions",
	}
}

type Role struct {
	Name        string
	Description string
	Permissions Permissions `gorm:"many2many:roles_permissions"`

	Serial
}

type Roles []Role

// OIDCProvider contains the data needed to initiate an OIDC secure login flow
type OIDCProvider struct {
	ClientID      string
	Issuer        string
	SSOProviderID int

	Serial
}

func (OIDCProvider) TableName() string {
	return "oidc_providers"
}

const (
	ObjectIDAttributeNameFormat = "urn:oasis:names:tc:SAML:2.0:attrname-format:uri"
	ObjectIDEmail               = "urn:oid:0.9.2342.19200300.100.1.3"
	XMLTypeString               = "xs:string"
	XMLSOAPClaimsEmailAddress   = "http://schemas.xmlsoap.org/ws/2005/05/identity/claims/emailaddress"
)

var (
	ErrSAMLAssertion = errors.New("SAML assertion error")
)

// SAMLRootURIVersion is required for payloads to match the ACS / Callback url configured on IDPs
// While the DB column root_uri_version has a default of 2, it is also hardcoded in the db method CreateSAMLIdentityProvider
type SAMLRootURIVersion int

var (
	SAMLRootURIVersion1 SAMLRootURIVersion = 1
	SAMLRootURIVersion2 SAMLRootURIVersion = 2

	SAMLRootURIVersionMap = map[SAMLRootURIVersion]string{
		SAMLRootURIVersion1: "/api/v1/login/saml",
		SAMLRootURIVersion2: "/api/v2/sso",
	}
)

type SAMLProvider struct {
	Name            string `gorm:"unique;index"`
	DisplayName     string
	IssuerURI       string
	SingleSignOnURI string
	MetadataXML     []byte
	RootURIVersion  SAMLRootURIVersion

	// PrincipalAttributeMapping is an array of OID or XML Namespace element mapping strings that can be used to map a
	// SAML assertion to a user in the database.
	//
	// For example: ["http://schemas.xmlsoap.org/ws/2005/05/identity/claims/emailaddress", "urn:oid:0.9.2342.19200300.100.1.3"]
	PrincipalAttributeMappings []string `gorm:"type:text[];column:ous"`

	// The below values generated values that point a client to SAML related resources hosted on the BloodHound instance
	// and should not be persisted to the database due to the fact that the URLs rely on the Host header that the user is
	// using to communicate to the API
	ServiceProviderIssuerURI     serde.URL `gorm:"-"`
	ServiceProviderInitiationURI serde.URL `gorm:"-"`
	ServiceProviderMetadataURI   serde.URL `gorm:"-"`
	ServiceProviderACSURI        serde.URL `gorm:"-"`

	SSOProviderID null.Int32

	Serial
}

type SAMLProviders []SAMLProvider

func (SAMLProvider) TableName() string {
	return "saml_providers"
}

// SSOProvider is the common representation of an SSO provider that can be used to display high level information about that provider
type SSOProvider struct {
	Type SessionAuthProvider `gorm:"column:type"`
	Name string
	Slug string

	OIDCProvider *OIDCProvider `gorm:"foreignKey:SSOProviderID"`
	SAMLProvider *SAMLProvider `gorm:"foreignKey:SSOProviderID"`

	Serial
}

// Used by gorm to preload / instantiate the user FK'd tables data
func UserAssociations() []string {
	return []string{
		"SSOProvider",
		"SSOProvider.SAMLProvider", // Needed to populate the child provider
		"SSOProvider.OIDCProvider", // Needed to populate the child provider
		"AuthSecret",
		"AuthTokens",
		"Roles.Permissions",
	}
}

type User struct {
	SSOProvider   *SSOProvider
	SSOProviderID null.Int32
	AuthSecret    *AuthSecret `gorm:"constraint:OnDelete:CASCADE;"`
	AuthTokens    AuthTokens  `gorm:"constraint:OnDelete:CASCADE;"`
	Roles         Roles       `gorm:"many2many:users_roles"`
	FirstName     null.String
	LastName      null.String
	EmailAddress  null.String
	PrincipalName string `gorm:"unique;index"`
	LastLogin     time.Time
	IsDisabled    bool

	// EULA Acceptance does not pertain to Bloodhound Community Edition; this flag is used for Bloodhound Enterprise users.
	// This value is automatically set to true for Bloodhound Community Edition in the patchEULAAcceptance and CreateUser functions.
	EULAAccepted bool

	Unique
}

type Users []User

func (s Users) IsSortable(column string) bool {
	switch column {
	case "first_name",
		"last_name",
		"email_address",
		"principal_name",
		"last_login",
		"created_at",
		"updated_at",
		"deleted_at":
		return true
	default:
		return false
	}
}

// Used by gorm to preload / instantiate the user FK'd tables data
func UserSessionAssociations() []string {
	return []string{
		"User.SSOProvider",
		"User.SSOProvider.SAMLProvider", // Needed to populate the child provider
		"User.SSOProvider.OIDCProvider", // Needed to populate the child provider
		"User.AuthSecret",
		"User.AuthTokens",
		"User.Roles.Permissions",
	}
}

type SessionAuthProvider int

const (
	SessionAuthProviderSecret SessionAuthProvider = 0
	SessionAuthProviderSAML   SessionAuthProvider = 1
	SessionAuthProviderOIDC   SessionAuthProvider = 2
)

type SessionFlagKey string

const (
	SessionFlagFedEULAAccepted SessionFlagKey = "fed_eula_accepted" // INFO: The FedEULA is only applicable to select enterprise installations
)

type UserSession struct {
	User             User `gorm:"constraint:OnDelete:CASCADE;"`
	UserID           uuid.UUID
	AuthProviderType SessionAuthProvider
	AuthProviderID   int32 // If SSO Session, this will be the child saml or oidc provider id
	ExpiresAt        time.Time
	Flags            types.JSONBBoolObject

	BigSerial
}
