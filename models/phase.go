package models

import (
	"encoding/json"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
	"time"
)

type Phase struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	User    User        `belongs_to:"user"`
	UserID  uuid.UUID   `json:"user_id" db:"user_id"`
	Data      string    `json:"data" db:"data"`
}

// String is not required by pop and may be deleted
func (p Phase) String() string {
	jp, _ := json.Marshal(p)
	return string(jp)
}

// Phases is not required by pop and may be deleted
type Phases []Phase

// String is not required by pop and may be deleted
func (p Phases) String() string {
	jp, _ := json.Marshal(p)
	return string(jp)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (p *Phase) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: p.Data, Name: "Data"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (p *Phase) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (p *Phase) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// AddPhase example
type AddPhase struct {
	Data string `json:"data" example:'{"user_id": "11111"}'`
}

// UpdatePhase example
type UpdatePhase struct {
	Name string `json:"name" example:"Phase name"`
}
