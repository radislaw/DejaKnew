package auth

import "testing"

func TestHashAndCheckPassword(t *testing.T) {
	hash, err := HashPassword("s3cret")
	if err != nil {
		t.Fatalf("hash failed: %v", err)
	}
	if hash == "s3cret" {
		t.Fatal("hash must not equal plaintext")
	}
	if !CheckPassword(hash, "s3cret") {
		t.Fatal("correct password rejected")
	}
	if CheckPassword(hash, "wrong") {
		t.Fatal("wrong password accepted")
	}
}
