package resin

import (
	"os"
	"testing"
)

func TestAuth(t *testing.T) {
	e := os.Getenv("RESIN_EMAIL")
	p := os.Getenv("RESIN_PASSWORD")

	if e == "" || p == "" {
		t.Skip("Skipping integration test, env vars not defined")
	} else {
		token, err := Signup("https://api.resinstaging.io", e, p)
		if err != nil {
			t.Errorf("Error signing up: %s", err)
		} else if token == "" {
			t.Error("Empty token after signup")
		} else if id1, err := GetUserId(token); err != nil {
			t.Error("Invalid token after signup")
		} else {
			token2, err := Login("https://api.resinstaging.io", e, p)
			if err != nil {
				t.Errorf("Error logging in: %s", err)
			} else if token == "" {
				t.Error("Empty token after login")
			} else if id2, err := GetUserId(token2); err != nil {
				t.Error("Invalid token after login")
			} else if id1 != id2 {
				t.Error("User ids don't match on login and signup")
			}
		}
	}

}
