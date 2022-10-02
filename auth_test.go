package ft_test

import (
	"os"
	"testing"

	"github.com/goexl/ft"
)

func TestPublicKey(t *testing.T) {
	client, _ := ft.New()
	req := new(ft.PublicKeyReq)
	req.AppId = os.Getenv(`APP_ID`)
	if rsp, err := client.PublicKey(req); nil != err {
		t.Fail()
	} else if `` == rsp.Key {
		t.Fail()
	}
}
