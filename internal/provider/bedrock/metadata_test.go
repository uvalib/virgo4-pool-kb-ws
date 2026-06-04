package bedrock

import "testing"

type memberString struct {
	Value string
}

func TestMetadataValueToStringPlainString(t *testing.T) {
	got := metadataValueToString("uva-lib:123")
	if got != "uva-lib:123" {
		t.Fatalf("got %q", got)
	}
}

func TestMetadataValueToStringMemberStruct(t *testing.T) {
	got := metadataValueToString(&memberString{Value: "uva-lib:329370"})
	if got != "uva-lib:329370" {
		t.Fatalf("got %q", got)
	}
}

func TestMetadataValueToStringWrappedStructPrint(t *testing.T) {
	got := metadataValueToString("&{uva-lib:329370}")
	if got != "uva-lib:329370" {
		t.Fatalf("got %q", got)
	}
}

func TestMetadataValueToStringValueFieldFormat(t *testing.T) {
	got := metadataValueToString(`&{Value: "My Title"}`)
	if got != "My Title" {
		t.Fatalf("got %q", got)
	}
}

func TestExtractMetadataStringFromMap(t *testing.T) {
	meta := map[string]interface{}{
		"id":         &memberString{Value: "uva-lib:123"},
		"iiif_id":    &memberString{Value: "uva-lib:999"},
		"title":      &memberString{Value: "Sample Title"},
		"collection": &memberString{Value: "Sample Collection"},
		"subject":    &memberString{Value: "Sample Subject"},
		"notes":      &memberString{Value: "Sample Notes"},
		"location":   &memberString{Value: "Charlottesville, VA"},
	}
	cases := map[string]string{
		"id":         "uva-lib:123",
		"iiif_id":    "uva-lib:999",
		"title":      "Sample Title",
		"collection": "Sample Collection",
		"subject":    "Sample Subject",
		"notes":      "Sample Notes",
		"location":   "Charlottesville, VA",
	}
	for key, want := range cases {
		if got := extractMetadataString(meta, key); got != want {
			t.Fatalf("%s=%q want %q", key, got, want)
		}
	}
}
