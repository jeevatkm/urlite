package model

import (
	"strings"
)

type UrliteRequest struct {
	LongUrl    string `json:"long_url"`
	Domain     string `json:"domain"`
	CustomName string `json:"custom_name"`
}

func (s *UrliteRequest) IsValid() bool {
	s.LongUrl = strings.TrimSpace(s.LongUrl)
	s.Domain = strings.TrimSpace(s.Domain)
	s.CustomName = strings.TrimSpace(s.CustomName)

	if len(s.LongUrl) > 0 && len(s.Domain) > 0 {
		return true
	}

	return false
}
