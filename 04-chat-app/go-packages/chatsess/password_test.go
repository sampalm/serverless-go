package chatsess

import "testing"

func TestPass(t *testing.T) {
	p := NewPassword("something")

	if !CheckPassword("something", p) {
		t.Errorf("something no match")
	}

	if CheckPassword("helloworld", p) {
		t.Errorf("helloworld matches")
	}

	if !CheckPassword("devsam", "413e09d0d7e2cec716e2_dec76cd496e158d8fd1ac5cd1c1df6eac7dc48ec6406c265442a36850f8ef5bc") {
		t.Errorf("something no match")
	}
}
