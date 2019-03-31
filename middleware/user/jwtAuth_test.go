package middleware

import (
	"fmt"
	"testing"
	"time"
)

func TestJwt(t *testing.T) {
	token, err := CreateToken(CustomClaims{
		ID: "5c98aeaafb7d8c4e73a1f712",
	})
	fmt.Println(token, err)
	time.Sleep(time.Second * 3)
	entry, code := ParseToken(token)
	if code > 0 {
		t.Error("error")
	}
	fmt.Printf("%+v, %d", entry, code)
}
