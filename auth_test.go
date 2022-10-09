package ft_test

import (
	"os"
	"testing"

	"github.com/goexl/ft"
)

func TestPublicKey(t *testing.T) {
	client, _ := ft.New()
	id := os.Getenv(`APP_ID`)
	key := os.Getenv(`APP_KEY`)
	secret := os.Getenv(`APP_SECRET`)
	if _key, err := client.PublicKey(ft.App(id, key, secret)); nil != err {
		t.Fail()
	} else if `` == _key {
		t.Fail()
	}
}

func TestToken(t *testing.T) {
	client, _ := ft.New()
	id := os.Getenv(`APP_ID`)
	key := os.Getenv(`APP_KEY`)
	secret := os.Getenv(`APP_SECRET`)
	if token, err := client.Token(ft.App(id, key, secret)); nil != err {
		t.Fail()
	} else if `` == token {
		t.Fail()
	}
}
