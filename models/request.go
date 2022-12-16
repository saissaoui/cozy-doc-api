package models

import (
	"encoding/json"

	"github.com/asaskevich/govalidator"
)

type DocumentRequest struct {
	Database  string            `json:"database" valid:"required"`
	Documents []json.RawMessage `json:"docs" valid:"required"`
}

func (req DocumentRequest) Validate() error {
	if _, err := govalidator.ValidateStruct(req); err != nil {
		return err
	}
	return nil
}
