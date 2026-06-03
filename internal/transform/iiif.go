package transform

import (
	"strings"

	"github.com/uvalib/virgo4-pool-kb-ws/internal/provider"
)

const defaultIIIFImageBaseURL = "https://iiif.lib.virginia.edu/iiif"

// Options controls Virgo field normalization for KB hits.
type Options struct {
	IIIFImageBaseURL string
}

func (o Options) imageBaseURL() string {
	if strings.TrimSpace(o.IIIFImageBaseURL) != "" {
		return strings.TrimRight(strings.TrimSpace(o.IIIFImageBaseURL), "/")
	}
	return defaultIIIFImageBaseURL
}

func (o Options) resolveImageURL(hit provider.Hit) string {
	if url := strings.TrimSpace(hit.IIIFImageURL); url != "" {
		return strings.TrimRight(url, "/")
	}
	if hit.IIIFID == "" {
		return ""
	}
	return o.imageBaseURL() + "/" + strings.TrimPrefix(hit.IIIFID, "/")
}

func (o Options) resolveManifestURL(hit provider.Hit) string {
	if url := strings.TrimSpace(hit.IIIFManifestURL); url != "" {
		return strings.TrimRight(url, "/")
	}
	if imageURL := o.resolveImageURL(hit); imageURL != "" {
		return imageURL + "/manifest"
	}
	return ""
}
