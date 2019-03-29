package middleware

import (
	"fmt"
	"testing"
)

func TestJwt(t *testing.T) {
	token, err := CreateToken(CustomClaims{
		ID: "5c98aeaafb7d8c4e73a1f712",
	})
	fmt.Println(token, err)
	entry, code := ParseToken(token)
	fmt.Printf("%+v, %d", entry, code)
}
