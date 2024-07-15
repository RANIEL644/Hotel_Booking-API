package models

import (
	"time"
)

type User struct {
	// ID            int       `json:"id"`
	Username      string    `json:"username,omitempty"`
	Email         string    `json:"email,omitempty"`
	Password      string    `json:"password"`
	Phone_Number  string    `json:"phone_number,omitempty"`
	Token         string    `json:"token,omitempty"`
	User_Type     string    `json:"user_type,omitempty"`
	Refresh_token string    `json:"refresh_token,omitempty"`
	Created_at    time.Time `json:"created_at,omitempty"`
	Updated_at    time.Time `json:"updated_at,omitempty"`
	User_id       string    `json:"user_id,omitempty"`
}
