package schema

import "time"

type LongURL struct {
	Hash      string `bson:"hash"`
	Assembled string `bson:"assembled"`

	Scheme string `bson:"scheme"`
	Host   string `bson:"host"`
	Path   string `bson:"path"`
	Query  string `bson:"query"`
}

type User struct {
	ID    string `bson:"id"`
	Email string `bson:"email"`
}

type ShortURL struct {
	OrganizationID string     `bson:"organization_id"`
	LongURL        LongURL    `bson:"long_url"`
	Scheme         string     `bson:"scheme"`
	Domain         string     `bson:"domain"`
	Alias          string     `bson:"alias"`
	CreatedAt      time.Time  `bson:"created_at"`
	ExpiresAt      *time.Time `bson:"expires_at"`
	CreatedBy      User       `bson:"created_by"`
	Description    string     `bson:"description"`
	Tags           []string   `bson:"tags"`
}

func (s ShortURL) OrganizationKey() ShortURLOrganizationKey {
	return ShortURLOrganizationKey{OrganizationID: s.OrganizationID, Domain: s.Domain, LongURLHash: s.LongURL.Hash}
}

func (s ShortURL) GlobalKey() ShortURLGlobalKey {
	return ShortURLGlobalKey{Domain: s.Domain, Alias: s.Alias}
}

type ShortURLOrganizationKey struct {
	OrganizationID string
	Domain         string
	LongURLHash    string
}

type ShortURLGlobalKey struct {
	Domain string
	Alias  string
}
