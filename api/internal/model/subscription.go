package model

import "time"

type PushSubscription struct {
	CreatedAt time.Time `json:"createdAt"`
	ID        string    `json:"id"`
	Endpoint  string    `json:"endpoint"`
	P256dh    string    `json:"p256dh"`
	Auth      string    `json:"auth"`
}

func (s *PushSubscription) Validate() error {
	if s.Endpoint == "" {
		return ErrFieldRequired("endpoint")
	}
	if s.P256dh == "" {
		return ErrFieldRequired("p256dh")
	}
	if s.Auth == "" {
		return ErrFieldRequired("auth")
	}
	return nil
}
