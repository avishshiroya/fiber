package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Auth struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserId    uuid.UUID `gorm:"type:uuid" json:"userId"`
	Email     string    `gorm:"size:255,not null" json:"email" validate:"lte=255,required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (a *Auth) BeforeCreate(tx *gorm.DB) (err error) {
	// Set UUID if not provided
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}

	// Set timestamps
	a.CreatedAt = time.Now()
	a.UpdatedAt = time.Now()
	return nil
}

// BeforeUpdate hook to update UpdatedAt timestamp
func (a *Auth) BeforeUpdate(tx *gorm.DB) (err error) {
	a.UpdatedAt = time.Now()
	return nil
}
