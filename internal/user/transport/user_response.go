package transport

import "time"

type UserResponse struct {
	ID        uint64     `json:"id,omitempty"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}
