package url

import (
	"encoding/json"
	"net/url"
)

type (
	Request interface {
		URL() *url.URL
	}

	Response interface {
		Result() *url.URL
	}

	response struct {
		result *url.URL
	}

	request struct {
		url *url.URL
	}
)

func (req *request) URL() *url.URL {
	return req.url
}

func (req *request) UnmarshalJSON(data []byte) error {
	raw := struct {
		URL string `json:"url"`
	}{}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	u, err := url.ParseRequestURI(raw.URL)
	if err != nil {
		return ErrParseURL
	}

	if u.Host == "" || u.Scheme == "" {
		return ErrParseURL
	}

	req.url = u

	return nil
}

func (res *response) Result() *url.URL {
	return res.result
}

func (res *response) MarshalJSON() ([]byte, error) {
	type alias struct {
		Result string `json:"result"`
	}
	var empty url.URL
	if *res.result == empty {
		return nil, ErrEmptyURL
	}
	aliasValue := new(alias)
	aliasValue.Result = res.result.String()
	return json.Marshal(aliasValue)
}

func NewURLRequest(data *url.URL) *request {
	return &request{
		url: data,
	}
}

func NewURLResponse(result *url.URL) *response {
	return &response{
		result: result,
	}
}
