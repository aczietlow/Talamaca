package net

import "net/http"

type AuthTransport struct {
	Token     string
	Transport http.RoundTripper
}

var Client *http.Client

func (t *AuthTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	if t.Token != "" {
		request.Header.Set("Authorization", t.Token)
	}
	return t.Transport.RoundTrip(request)
}

func NewClient(token string) {
	Client = &http.Client{
		Transport: &AuthTransport{
			Token:     token,
			Transport: http.DefaultTransport,
		},
	}
}
