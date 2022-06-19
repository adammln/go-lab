
package main

import (
	"testing"

	"github.com/joho/godotenv"
)

func TestDotEnvLoader(t *testing.T) {
	err := godotenv.Load()
	if (err != nil) {
		t.Fatal(err)
	}
}