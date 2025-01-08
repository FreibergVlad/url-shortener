package db

import "time"

type LongURL struct {
	Hash      string `bson:"hash"`
	Assembled string `bson:"assembled"`

	Scheme string `bson:"scheme"`
	Host   string `bson:"host"`
	Path   string `bson:"path"`
	Query  string `bson:"query"`
}

type ShortURL struct {
	OrganizationID string     `bson:"organization_id"`
	LongURL        LongURL    `bson:"long_url"`
	Scheme         string     `bson:"scheme"`
	Domain         string     `bson:"domain"`
	Alias          string     `bson:"alias"`
	Status         string     `bson:"status"`
	CreatedAt      time.Time  `bson:"created_at"`
	ExpiresAt      *time.Time `bson:"expires_at"`
	DeletedAt      *time.Time `bson:"deleted_at"`
	CreatedBy      string     `bson:"created_by"`
	Description    *string    `bson:"description"`
	Tags           []string   `bson:"tags"`
}
