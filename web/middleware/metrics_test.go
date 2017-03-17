package middleware

import (
	"testing"
)

func TestNormalizePath(t *testing.T) {
	if in := ""; normalizePathForMetric(in) != in {
		t.Fatalf("Unexpected return for %s", in)
	}
	if in := "/"; normalizePathForMetric(in) != in {
		t.Fatalf("Unexpected return for %s", in)
	}
	if in := "/0/1"; normalizePathForMetric(in) != "/{id}/{id}" {
		t.Fatalf("Unexpected return for %s", in)
	}
	if in := "//foo/09/"; normalizePathForMetric(in) != "//foo/{id}/" {
		t.Fatalf("Unexpected return for %s", in)
	}
}
