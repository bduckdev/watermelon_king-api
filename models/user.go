package models

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	CreatedAt *time.Time     `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-"`
	Username  string         `json:"username" gorm:"unique" validate:"required"`
	Email     string         `json:"email" gorm:"unique" validate:"required,email"`
	Score     int            `json:"score"`
	Password  string         `json:"password" validate:"required"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()

	return
}

var validate *validator.Validate

func ValidateNewUser(user *User) error {
	validate = validator.New(validator.WithRequiredStructEnabled())

	errs := validate.Struct(user)
	if errs != nil {
		return errs
	}
	return errs
}
