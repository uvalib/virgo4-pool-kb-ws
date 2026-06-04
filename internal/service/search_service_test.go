package service

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/uvalib/virgo4-api/v4api"
	"github.com/uvalib/virgo4-pool-kb-ws/internal/provider"
)

func TestSearchWithMockProvider(t *testing.T) {
	svc := &SearchService{
		Provider: &provider.MockProvider{
			Hits: []provider.Hit{
				{ID: "img-1", Title: "First"},
				{ID: "img-2", Title: "Second"},
			},
		},
		DefaultLimit: 20,
	}

	resp, status, err := svc.Search(context.Background(), &v4api.SearchRequest{
		Query:      `keyword: {test}`,
		Pagination: v4api.Pagination{Rows: 1},
	})
	if err != nil {
		t.Fatal(err)
	}
	if status != http.StatusOK {
		t.Fatalf("status=%d", status)
	}
	if len(resp.Groups) != 1 {
		t.Fatalf("expected 1 group (limit 1), got %d", len(resp.Groups))
	}
	if resp.Groups[0].Value != "img-1" {
		t.Fatalf("unexpected group id %q", resp.Groups[0].Value)
	}
}

func TestSearchProviderError(t *testing.T) {
	svc := &SearchService{
		Provider:     &provider.MockProvider{Err: errors.New("kb down")},
		DefaultLimit: 20,
	}
	resp, status, err := svc.Search(context.Background(), &v4api.SearchRequest{Query: `keyword: {x}`})
	if err != nil {
		t.Fatal(err)
	}
	if status != http.StatusBadGateway {
		t.Fatalf("status=%d", status)
	}
	if resp.StatusCode != http.StatusBadGateway {
		t.Fatalf("pool status=%d", resp.StatusCode)
	}
}

func TestResourceExactMatch(t *testing.T) {
	hit := provider.Hit{
		ID:         "target-id",
		IIIFID:     "target-id",
		Title:      "Right",
		Collection: "Sample Collection",
		Subject:    "Sample Subject",
		Notes:      "Sample Notes",
		Location:   "Charlottesville, VA",
		Content:    "Sample summary content",
	}
	svc := &SearchService{
		Provider: &provider.MockProvider{
			Hits: []provider.Hit{
				{ID: "other", Title: "Wrong"},
				hit,
			},
		},
		DefaultLimit: 20,
	}
	fields, status, err := svc.Resource(context.Background(), "target-id")
	if err != nil {
		t.Fatal(err)
	}
	if status != http.StatusOK {
		t.Fatalf("status=%d", status)
	}

	want := map[string]string{
		"id":                 "target-id",
		"title":              "Right",
		"digital_collection": "Sample Collection",
		"subject":            "Sample Subject",
		"notes":              "Sample Notes",
		"location":           "Charlottesville, VA",
		"iiif_image_url":     "https://iiif.lib.virginia.edu/iiif/target-id",
		"iiif_manifest_url":  "https://iiif.lib.virginia.edu/iiif/target-id/manifest",
		"summary":            "Sample summary content",
	}
	got := make(map[string]string, len(fields))
	for _, f := range fields {
		got[f.Name] = f.Value
	}
	for name, value := range want {
		if got[name] != value {
			t.Fatalf("%q=%q want %q", name, got[name], value)
		}
	}
	if len(fields) != len(want) {
		t.Fatalf("got %d fields want %d: %+v", len(fields), len(want), fields)
	}
}

func TestResourceNoExactMatchReturnsNotFound(t *testing.T) {
	svc := &SearchService{
		Provider: &provider.MockProvider{
			Hits: []provider.Hit{{ID: "only-other", Title: "Nope"}},
		},
		DefaultLimit: 20,
	}
	_, status, err := svc.Resource(context.Background(), "missing-id")
	if err == nil {
		t.Fatal("expected error")
	}
	if status != http.StatusNotFound {
		t.Fatalf("status=%d want 404", status)
	}
}
