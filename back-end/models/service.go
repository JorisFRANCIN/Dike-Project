package models

import (
	"time"
)

type About struct {
	Client struct {
		Host string `json:"host"`
	} `json:"client"`
	Server struct {
		CurrentTime int64     `json:"current_time"`
		Services    []Service `json:"services"`
	} `json:"server"`
}

type Service struct {
	Name          string `json:"name"`
	Description   string `json:"description,omitempty"`
	BackgroundColor string `json:"backgroundColor"`
}

type ServiceAccessToken struct {
    ServiceName   string `json:"service_name"`
    AccessToken   string `json:"access_token"`
    RefreshToken  string `json:"refresh_token"`
    ExpiresAt     time.Time `json:"expires_at"`
}

