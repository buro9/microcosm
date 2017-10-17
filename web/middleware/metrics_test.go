package middleware

import (
	"testing"
)

func TestNormalizePath(t *testing.T) {
	if testPath := ""; normalizePathForMetric(testPath) != testPath {
		t.Fatalf("Unexpected return for %s", testPath)
	}
	if testPath := "/"; normalizePathForMetric(testPath) != testPath {
		t.Fatalf("Unexpected return for %s", testPath)
	}
	if testPath := "/0/1"; normalizePathForMetric(testPath) != "/{id}/{id}" {
		t.Fatalf("Unexpected return for %s", testPath)
	}
	if testPath := "//foo/09/"; normalizePathForMetric(testPath) != "//foo/{id}/" {
		t.Fatalf("Unexpected return for %s", testPath)
	}
}
