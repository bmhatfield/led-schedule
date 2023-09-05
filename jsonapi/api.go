package jsonapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func Get[T any](u url.URL) (*T, error) {
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid status code: %d", resp.StatusCode)
	}

	t := new(T)
	if err := json.NewDecoder(resp.Body).Decode(t); err != nil {
		return nil, err
	}

	return t, nil
}
