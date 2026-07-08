// Package vapid provides helpers for generating WebPush VAPID key pairs.
package vapid

import (
	"fmt"

	webpush "github.com/SherClockHolmes/webpush-go"
)

// Keys holds a base64url-encoded VAPID key pair used for WebPush.
type Keys struct {
	Public  string
	Private string
}

// GenerateKeys creates a new VAPID key pair suitable for WebPush.
func GenerateKeys() (Keys, error) {
	private, public, err := webpush.GenerateVAPIDKeys()
	if err != nil {
		return Keys{}, fmt.Errorf("generate VAPID keys: %w", err)
	}
	return Keys{Public: public, Private: private}, nil
}

// EnvLines renders the key pair as .env assignments for the API and frontend.
// The public key is emitted for both the backend (VAPID_PUBLIC_KEY) and the
// frontend (NEXT_PUBLIC_VAPID_PUBLIC_KEY), which must share the same value.
func (k Keys) EnvLines() string {
	return fmt.Sprintf(
		"VAPID_PUBLIC_KEY=%s\nVAPID_PRIVATE_KEY=%s\nNEXT_PUBLIC_VAPID_PUBLIC_KEY=%s\n",
		k.Public, k.Private, k.Public,
	)
}
