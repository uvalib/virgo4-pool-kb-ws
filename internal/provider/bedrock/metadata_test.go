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
		"iiif_id": &memberString{Value: "uva-lib:999"},
		"title":   &memberString{Value: "Sample Title"},
	}
	if got := extractMetadataString(meta, "iiif_id"); got != "uva-lib:999" {
		t.Fatalf("iiif_id=%q", got)
	}
	if got := extractMetadataString(meta, "title"); got != "Sample Title" {
		t.Fatalf("title=%q", got)
	}
}
