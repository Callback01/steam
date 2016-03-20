package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

const (
	apiKeyURL = "https://steamcommunity.com/dev/apikey"

	accessDeniedPattern = "<h2>Access Denied</h2>"
)

var (
	keyRegExp = regexp.MustCompile("<p>Key: ([0-9A-F]+)</p>")

	ErrAccessDenied = errors.New("access is denied")
	ErrKeyNotFound  = errors.New("key not found")
)

func (community *Community) getWebAPIKey() (string, error) {
	req, err := http.NewRequest(http.MethodGet, apiKeyURL, nil)
	if err != nil {
		return "", err
	}

	resp, err := community.client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if m, err := regexp.Match(accessDeniedPattern, body); err != nil {
		return "", err
	} else if m {
		return "", ErrAccessDenied
	}

	submatch := keyRegExp.FindStringSubmatch(string(body))
	if len(submatch) == 0 {
		return "", ErrKeyNotFound
	}

	community.apiKey = submatch[1]
	return submatch[1], nil
}
