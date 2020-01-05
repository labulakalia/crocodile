package jwt

import "testing"

var (
	token string
)

func Test(t *testing.T) {
	t.Run("generate token", TestGenerateToken)
	t.Run("parse token", TestParseToken)
}

func TestGenerateToken(t *testing.T) {
	var (
		err error
	)
	token, err = GenerateToken("121212121")
	if err != nil {
		t.Errorf("GenerateToken failed: %v", err)
	}
	t.Logf("Token: %s\n", token)
}

func TestParseToken(t *testing.T) {
	calims, err := ParseToken(token)
	if err != nil {
		t.Errorf("ParseToken failed: %v", err)
	}
	t.Logf("User Id: %s", calims.Id)
}
