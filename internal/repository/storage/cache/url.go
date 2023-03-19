package cache

import (
	"encoding/json"
	"net/url"

	"github.com/google/uuid"
)

type item struct {
	UserID uuid.UUID
	Value  *url.URL
}

func (i item) MarshalJSON() ([]byte, error) {
	type alias struct {
		UserID uuid.UUID `json:"user_id"`
		Value  string    `json:"value"`
	}
	aliasValue := alias{}
	aliasValue.UserID = i.UserID
	aliasValue.Value = i.Value.String()
	return json.Marshal(aliasValue)
}

func (i *item) UnmarshalJSON(data []byte) error {
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

	i.UserID = aliasValue.UserID
	i.Value = parsedValue

	return nil
}
