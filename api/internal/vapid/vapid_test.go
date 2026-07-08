package vapid_test

import (
	"strings"
	"testing"

	"github.com/akaitigo/kenketsu-plus/api/internal/vapid"
)

func TestGenerateKeys(t *testing.T) {
	t.Parallel()

	keys, err := vapid.GenerateKeys()
	if err != nil {
		t.Fatalf("GenerateKeys failed: %v", err)
	}
	if keys.Public == "" || keys.Private == "" {
		t.Fatalf("expected non-empty keys, got public=%q private=%q", keys.Public, keys.Private)
	}
	if keys.Public == keys.Private {
		t.Error("public and private keys must differ")
	}
}

func TestGenerateKeys_Unique(t *testing.T) {
	t.Parallel()

	first, err := vapid.GenerateKeys()
	if err != nil {
		t.Fatalf("first GenerateKeys failed: %v", err)
	}
	second, err := vapid.GenerateKeys()
	if err != nil {
		t.Fatalf("second GenerateKeys failed: %v", err)
	}
	if first.Private == second.Private || first.Public == second.Public {
		t.Error("expected distinct key pairs across calls")
	}
}

func TestEnvLines(t *testing.T) {
	t.Parallel()

	keys := vapid.Keys{Public: "pub-key", Private: "priv-key"}
	out := keys.EnvLines()

	for _, want := range []string{
		"VAPID_PUBLIC_KEY=pub-key\n",
		"VAPID_PRIVATE_KEY=priv-key\n",
		"NEXT_PUBLIC_VAPID_PUBLIC_KEY=pub-key\n",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("EnvLines output missing %q; got:\n%s", want, out)
		}
	}
}
