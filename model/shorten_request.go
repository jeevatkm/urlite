package model

type ShortenRequest struct {
	LongUrl    string `json:"long_url"`
	Domain     string `json:"domain"`
	CustomName string `json:"custom_name"`
	Secure     bool   `json:"secure"`
	Password   string `json:"password"`
}
