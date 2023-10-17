package model

import (
	"time"

	_ "github.com/asaskevich/govalidator"
)

type Base struct {
	ID string `json:"id" valid:"uuid"`
	CreatedAt time.Time `json:"created_at" valid:"-"`
	UpdatedAt time.Time `json:"updated_at" valid:"-"`
}