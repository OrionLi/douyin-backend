package utils

import "testing"

func TestJWT(t *testing.T) {
	token, err := GenerateToken(1, "test", 1)
	if err != nil {
		t.Error(err)
	}
	t.Log(token)
}
