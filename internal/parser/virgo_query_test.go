package parser

import (
	"strings"
	"testing"
)

func TestValidateAcceptsKeywordQuery(t *testing.T) {
	ok, msg := Validate(`keyword: {cats}`)
	if !ok {
		t.Fatalf("expected valid query, got %q", msg)
	}
}

func TestToKBTextQueryStripsFieldPrefixes(t *testing.T) {
	got := ToKBTextQuery(`keyword: {title: {dogs}}`)
	if got == "" {
		t.Fatal("expected non-empty query")
	}
	if strings.Contains(got, "keyword:") || strings.Contains(got, "title:") {
		t.Fatalf("unexpected prefixes in %q", got)
	}
}
