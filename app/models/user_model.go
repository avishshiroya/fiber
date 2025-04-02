package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Username  string    `gorm:"size:255;not null" json:"username" validate:"required,lte=255"`
	Email     string    `gorm:"size:255;unique;not null" json:"email" validate:"required,email,lte=255"`
	Password  string    `gorm:"size:255;not null" json:"password" validate:"required,lte=255"`
}

// BeforeCreate hook to set ID and timestamps
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	// Set UUID if not provided
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}

	// Set timestamps
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	return nil
}

// BeforeUpdate hook to update UpdatedAt timestamp
func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	u.UpdatedAt = time.Now()
	return nil
}
