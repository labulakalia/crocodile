package notify

import "testing"

import "net/http"

func Test_JSONPost(t *testing.T) {
	_, err := JSONPost("http://webhook.test",nil,http.DefaultClient)
	if err != nil {
		t.Fatal(err)
	}
}
