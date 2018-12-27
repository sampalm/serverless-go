package chatsess

import (
	"fmt"
	"testing"
)

func TestBSalt(t *testing.T) {
	//BCrypt Password
	b := NewPasswordBcrypt("something")

	if !CheckPasswordBcrypt("something", b) {
		t.Errorf("bcrypt want true got %t", CheckPasswordBcrypt("something", b))
	}
	if CheckPasswordBcrypt("helloworld", b) {
		t.Errorf("bcrypt: want false got %t", CheckPasswordBcrypt("helloworld", b))
	}
}

func TestBCrypt(t *testing.T) {
	// BSalt Password
	p := NewPassword("something")

	if !CheckPassword("something", p) {
		t.Errorf("bsalt want true got %t", CheckPassword("something", p))
	}
	if CheckPassword("helloworld", p) {
		t.Errorf("bsalt: want false got %t", CheckPassword("helloworld", p))
	}
	if !CheckPassword("devsam", "413e09d0d7e2cec716e2_dec76cd496e158d8fd1ac5cd1c1df6eac7dc48ec6406c265442a36850f8ef5bc") {
		t.Errorf("bsalt: want true got %t", CheckPassword("devsam", "413e09d0d7e2cec716e2_dec76cd496e158d8fd1ac5cd1c1df6eac7dc48ec6406c265442a36850f8ef5bc"))
	}
}

func ExampleBSalt() {
	p := NewPassword("something")
	fmt.Printf("bsalt: 'something' matches? %t", CheckPassword("something", p))
	// Output:
	// bsalt: 'something' matches? true
}
func ExampleBCrypt() {
	b := NewPasswordBcrypt("something")
	fmt.Printf("bcrypt: 'something' matches? %t", CheckPasswordBcrypt("something", b))
	// Output:
	// bcrypt: 'something' matches? true
}
func BenchmarkBSalt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewPassword("something")
	}
}

func BenchmarkBCrypt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewPasswordBcrypt("something")
	}
}
