package models

import "time"

type SupportRequested struct {
	ClientName  string `json:"client_name"`
	Timestamp   time.Time
	IsSupported bool `json:"is_supported"`
}
