package cache

import (
	"encoding/json"
	"net/url"

	"github.com/google/uuid"
)

type storageURL struct {
	UserID  uuid.UUID
	Value   url.URL
	Deleted bool
}

func (s storageURL) MarshalJSON() ([]byte, error) {
	type alias struct {
		UserID uuid.UUID `json:"user_id"`
		Value  string    `json:"value"`
	}
	aliasValue := alias{}
	aliasValue.UserID = s.UserID
	aliasValue.Value = s.Value.String()
	return json.Marshal(aliasValue)
}

func (s *storageURL) UnmarshalJSON(data []byte) error {
	type alias struct {
		UserID uuid.UUID `json:"user_id"`
		Value  string    `json:"value"`
	}

	aliasValue := alias{}
	if err := json.Unmarshal(data, &aliasValue); err != nil {
		return err
	}

	parsedValue, err := url.ParseRequestURI(aliasValue.Value)
	if err != nil {
		return err
	}

	s.UserID = aliasValue.UserID
	s.Value = *parsedValue

	return nil
}
