// Command genvapid generates a WebPush VAPID key pair and prints it as .env
// assignments, so developers do not need to run an external tool during setup (#24).
//
// Usage:
//
//	go run ./cmd/genvapid   (or: make generate-vapid)
package main

import (
	"fmt"
	"os"

	"github.com/akaitigo/kenketsu-plus/api/internal/vapid"
)

func main() {
	keys, err := vapid.GenerateKeys()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to generate VAPID keys:", err)
		os.Exit(1)
	}
	fmt.Print(keys.EnvLines())
}
