package main

import (
	"testing"
	"os"
)

func TestLoadEnvironmentVariables(t *testing.T) {
	ENV_KEYS := []string{
		"RELATIVE_PATH_GCP_KEYFILE",
		"FIRESTORE_TEST_DATA_COLLECTION_ID",
		"FIRESTORE_DATA_COLLECTION_ID",
	}

	for _, key := range ENV_KEYS {
		val := os.Getenv(key)
		if val == "" {
			t.Fatalf(`Can't load environment variable %q`, key)
		}
	}

}