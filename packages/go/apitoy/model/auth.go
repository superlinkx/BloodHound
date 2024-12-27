package model

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/specterops/bloodhound/src/serde"
)

// OIDCProvider contains the data needed to initiate an OIDC secure login flow
type OIDCProvider struct {
	ClientID      string
	Issuer        string
	SSOProviderID int

	Serial
}

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
	Name            string
	DisplayName     string
	IssuerURI       string
	SingleSignOnURI string
	MetadataXML     []byte
	RootURIVersion  SAMLRootURIVersion

	// PrincipalAttributeMapping is an array of OID or XML Namespace element mapping strings that can be used to map a
	// SAML assertion to a user in the database.
	//
	// For example: ["http://schemas.xmlsoap.org/ws/2005/05/identity/claims/emailaddress", "urn:oid:0.9.2342.19200300.100.1.3"]
	PrincipalAttributeMappings []string

	// The below values generated values that point a client to SAML related resources hosted on the BloodHound instance
	// and should not be persisted to the database due to the fact that the URLs rely on the Host header that the user is
	// using to communicate to the API
	ServiceProviderIssuerURI     serde.URL
	ServiceProviderInitiationURI serde.URL
	ServiceProviderMetadataURI   serde.URL
	ServiceProviderACSURI        serde.URL

	SSOProviderID int32

	Serial
}

// SSOProvider is the common representation of an SSO provider that can be used to display high level information about that provider
type SSOProvider struct {
	Type SessionAuthProvider
	Name string
	Slug string

	OIDCProvider *OIDCProvider
	SAMLProvider *SAMLProvider

	Serial
}

type Permission struct {
	Authority string
	Name      string

	Serial
}

type Permissions []Permission

type AuthToken struct {
	UserID     uuid.NullUUID
	ClientID   uuid.NullUUID
	Name       string
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

type Role struct {
	Name        string
	Description string
	Permissions Permissions

	Serial
}

type Roles []Role

type User struct {
	SSOProvider   *SSOProvider
	SSOProviderID int32
	AuthSecret    *AuthSecret
	AuthTokens    AuthTokens
	Roles         Roles
	FirstName     string
	LastName      string
	EmailAddress  string
	PrincipalName string
	LastLogin     time.Time
	IsDisabled    bool

	// EULA Acceptance does not pertain to Bloodhound Community Edition; this flag is used for Bloodhound Enterprise users.
	// This value is automatically set to true for Bloodhound Community Edition in the patchEULAAcceptance and CreateUser functions.
	EULAAccepted bool

	Unique
}

type Users []User

type SessionAuthProvider int

const (
	SessionAuthProviderSecret SessionAuthProvider = 0
	SessionAuthProviderSAML   SessionAuthProvider = 1
	SessionAuthProviderOIDC   SessionAuthProvider = 2
)
