package main

import (
	"net/http"
)

type Client struct {
	HTTPClient *http.Client
}

var DefaultClient = &Client{
	HTTPClient: http.DefaultClient,
}
