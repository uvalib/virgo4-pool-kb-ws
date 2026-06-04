package parser

import (
	"strings"
	"log"
	"github.com/uvalib/virgo4-parser/v4parser"
)

// Validate checks that a Virgo query string is syntactically valid.
func Validate(virgoQuery string) (bool, string) {
	return v4parser.Validate(virgoQuery)
}

// ToKBTextQuery translates a Virgo query into plain text for semantic KB retrieval.
func ToKBTextQuery(virgoQuery string) string {
	q := strings.TrimSpace(virgoQuery)
	log.Printf("virgo query=%q", q)

	replacements := []struct{ old, new string }{
		{"{", ""},
		{"}", ""},
		{"keyword:", ""},
		{"title:", ""},
		{"author:", ""},
		{"subject:", ""},
		{"published:", ""},
		{"fulltext:", ""},
		{"series:", ""},
		{"identifier:", ""},
		{"filter:", ""},
		{"date:", ""},
	
	}
	for _, r := range replacements {
		q = strings.ReplaceAll(q, r.old, r.new)
	}

	q = strings.TrimSpace(q)
	if q == "" || q == "()" {
		return "*"
	}
	return q
}
