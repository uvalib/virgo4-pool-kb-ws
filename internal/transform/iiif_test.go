package transform

import (
	"testing"

	"github.com/uvalib/virgo4-pool-kb-ws/internal/provider"
)

func TestResolveManifestURLFromImageBase(t *testing.T) {
	opts := Options{}
	hit := provider.Hit{IIIFID: "uva-lib:999"}
	got := opts.resolveManifestURL(hit)
	want := "https://iiif.lib.virginia.edu/iiif/uva-lib:999/manifest"
	if got != want {
		t.Fatalf("manifest=%q want %q", got, want)
	}
}
