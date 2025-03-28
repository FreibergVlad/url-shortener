package schema

import (
	"time"
)

type Organization struct {
	ID        string
	Name      string
	Slug      string
	CreatedAt time.Time
	CreatedBy string
}

type ShortOrganization struct {
	ID   string
	Slug string
}

type User struct {
	ID           string
	PasswordHash string
	Email        string
	FullName     string
	RoleID       string
	CreatedAt    time.Time
}

type OrganizationMembership struct {
	Organization ShortOrganization
	RoleID       string
	CreatedAt    time.Time
}

type OrganizationMemberships []*OrganizationMembership

func (m OrganizationMemberships) OrganizationMembership(id string) *OrganizationMembership {
	for _, membership := range m {
		if membership.Organization.ID == id {
			return membership
		}
	}
	return nil
}

func (m OrganizationMemberships) HasOrganization(id string) bool {
	return m.OrganizationMembership(id) != nil
}

type Invitation struct {
	ID             string
	Status         string
	OrganizationID string
	Email          string
	RoleID         string
	CreatedAt      time.Time
	ExpiresAt      time.Time
	CreatedBy      string
}
