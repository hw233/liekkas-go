package token

import "testing"

func TestGenerateBearerToken(t *testing.T) {
	token, err := GenerateBearerToken(nil, "overlord")
	t.Log(token)
	t.Log(err)

}
