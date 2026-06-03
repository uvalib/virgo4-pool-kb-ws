package transform

import (
	"testing"

	"github.com/uvalib/virgo4-pool-kb-ws/internal/provider"
)

func TestHitsToPoolResultBuildsGroups(t *testing.T) {
	hits := []provider.Hit{{ID: "img-1", Title: "Test Image", Score: 0.9}}
	resp := HitsToPoolResult(hits, 0, 5, Options{})
	if len(resp.Groups) != 1 {
		t.Fatalf("expected 1 group, got %d", len(resp.Groups))
	}
	if resp.Groups[0].Value != "img-1" {
		t.Fatalf("unexpected group value %q", resp.Groups[0].Value)
	}
}

func TestRecordFieldsIncludesIIIFImageURL(t *testing.T) {
	hit := provider.Hit{
		ID:     "uva-lib:123",
		IIIFID: "uva-lib:123",
		Title:  "Sample",
	}
	fields := recordFields(hit, Options{})
	var imageURL string
	for _, f := range fields {
		if f.Name == "iiif_image_url" {
			imageURL = f.Value
			if f.Type != "iiif-image-url" {
				t.Fatalf("type=%q want iiif-image-url", f.Type)
			}
		}
	}
	want := "https://iiif.lib.virginia.edu/iiif/uva-lib:123"
	if imageURL != want {
		t.Fatalf("iiif_image_url=%q want %q", imageURL, want)
	}
}

func TestRecordFieldsPrefersMetadataImageURL(t *testing.T) {
	hit := provider.Hit{
		ID:           "uva-lib:123",
		IIIFID:       "uva-lib:123",
		IIIFImageURL: "https://example.test/iiif/uva-lib:123",
	}
	fields := recordFields(hit, Options{})
	for _, f := range fields {
		if f.Name == "iiif_image_url" && f.Value == "https://example.test/iiif/uva-lib:123" {
			return
		}
	}
	t.Fatal("expected metadata iiif_image_url to be used")
}
