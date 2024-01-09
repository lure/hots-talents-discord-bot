package github

import (
	"net/http"
	"time"
)

var githubClient http.Client = http.Client{
	Timeout: time.Second * 100,
}
