package auth

import (
	"testing"
)

func TestCreateToken(t *testing.T) {
	token, err := genTokenByTest(Info{
		Uid:            "admin",
		StaffId:        "admin",
		Name:           "admin",
		IsRefreshToken: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(token)
}
