package transform

import (
	"testing"

	"github.com/uvalib/virgo4-api/v4api"
	"github.com/uvalib/virgo4-pool-kb-ws/internal/provider"
)

func fullTestHit() provider.Hit {
	return provider.Hit{
		ID:         "uva-lib:123",
		IIIFID:     "uva-lib:123",
		Title:      "Sample Title",
		Collection: "Sample Collection",
		Subject:    "Sample Subject",
		Notes:      "Sample Notes",
		Location:   "Charlottesville, VA",
		Content:    "Sample summary content",
	}
}

func TestHitsToPoolResultBuildsGroups(t *testing.T) {
	hits := []provider.Hit{fullTestHit()}
	resp := HitsToPoolResult(hits, 0, 5, Options{})
	if len(resp.Groups) != 1 {
		t.Fatalf("expected 1 group, got %d", len(resp.Groups))
	}
	if resp.Groups[0].Value != "uva-lib:123" {
		t.Fatalf("unexpected group value %q", resp.Groups[0].Value)
	}
	if len(resp.Groups[0].Records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(resp.Groups[0].Records))
	}
	assertRecordFields(t, resp.Groups[0].Records[0].Fields)
}

func TestHitToDetailFieldsIncludesAllFields(t *testing.T) {
	assertRecordFields(t, HitToDetailFields(fullTestHit(), Options{}))
}

func TestRecordFieldsIncludesAllFields(t *testing.T) {
	assertRecordFields(t, recordFields(fullTestHit(), Options{}))
}

func assertRecordFields(t *testing.T, fields []v4api.RecordField) {
	t.Helper()

	want := map[string]struct {
		value string
		typ   string
		label string
	}{
		"id":                 {value: "uva-lib:123", typ: "identifier", label: "Identifier"},
		"title":              {value: "Sample Title", typ: "title", label: "Title"},
		"digital_collection": {value: "Sample Collection", typ: "collection", label: "Digital Collection"},
		"subject":            {value: "Sample Subject", typ: "subject", label: "Subject"},
		"notes":              {value: "Sample Notes", typ: "notes", label: "Notes"},
		"location":           {value: "Charlottesville, VA", typ: "location", label: "Location"},
		"iiif_image_url":     {value: "https://iiif.lib.virginia.edu/iiif/uva-lib:123", typ: "iiif-image-url", label: "IIIF Image"},
		"iiif_manifest_url":  {value: "https://iiif.lib.virginia.edu/iiif/uva-lib:123/manifest", typ: "iiif-manifest-url", label: "IIIF Manifest"},
		"summary":            {value: "Sample summary content", typ: "summary", label: "Summary"},
	}

	got := make(map[string]v4api.RecordField, len(fields))
	for _, f := range fields {
		got[f.Name] = f
	}

	for name, w := range want {
		f, ok := got[name]
		if !ok {
			t.Fatalf("missing field %q", name)
		}
		if f.Value != w.value {
			t.Fatalf("%q value=%q want %q", name, f.Value, w.value)
		}
		if f.Type != w.typ {
			t.Fatalf("%q type=%q want %q", name, f.Type, w.typ)
		}
		if f.Label != w.label {
			t.Fatalf("%q label=%q want %q", name, f.Label, w.label)
		}
	}
	if len(fields) != len(want) {
		t.Fatalf("got %d fields want %d: %+v", len(fields), len(want), fields)
	}
}
