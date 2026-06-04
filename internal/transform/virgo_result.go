package transform

import (
	"net/http"

	"github.com/uvalib/virgo4-api/v4api"
	"github.com/uvalib/virgo4-pool-kb-ws/internal/provider"
)

// HitsToPoolResult maps KB hits into the pool search response.
func HitsToPoolResult(hits []provider.Hit, start int, elapsedMS int64, opts Options) *v4api.PoolResult {
	resp := &v4api.PoolResult{
		ElapsedMS:  elapsedMS,
		Confidence: "low",
		StatusCode: http.StatusOK,
		Groups:     make([]v4api.Group, 0, len(hits)),
		Pagination: v4api.Pagination{Start: start, Total: len(hits), Rows: len(hits)},
	}

	for _, hit := range hits {
		id := hit.ID
		if id == "" {
			id = hit.IIIFID
		}
		if id == "" {
			continue
		}
		rec := v4api.Record{Fields: recordFields(hit, opts)}
		group := v4api.Group{Value: id, Count: 1, Records: []v4api.Record{rec}}
		resp.Groups = append(resp.Groups, group)
	}

	if len(resp.Groups) > 0 {
		resp.Confidence = "medium"
	}
	return resp
}

// HitToDetailFields returns detail fields for a single resource.
func HitToDetailFields(hit provider.Hit, opts Options) []v4api.RecordField {
	return recordFields(hit, opts)
}

func recordFields(hit provider.Hit, opts Options) []v4api.RecordField {
	fields := make([]v4api.RecordField, 0, 8)
	id := hit.ID
	fields = append(fields, v4api.RecordField{Name: "id", Type: "identifier", Label: "Identifier", Value: id, CitationPart: "id"})
	if hit.Title != "" {
		fields = append(fields, v4api.RecordField{Name: "title", Type: "title", Label: "Title", Value: hit.Title, CitationPart: "title"})
	}
	if hit.Collection != "" {
		fields = append(fields, v4api.RecordField{Name: "digital_collection", Type: "collection", Label: "Digital Collection", Value: hit.Collection})
	}
	// Omit iiif_id field to avoid duplicate identifiers in the search results
	//if hit.IIIFID != "" {
	//	fields = append(fields, v4api.RecordField{Name: "iiif_id", Type: "identifier", Label: "IIIF ID", Value: hit.IIIFID, Visibility: "detailed"})
	//}
	if imageURL := opts.resolveImageURL(hit); imageURL != "" {
		fields = append(fields, v4api.RecordField{Name: "iiif_image_url", Type: "iiif-image-url", Label: "IIIF Image", Value: imageURL})
	}
	if manifestURL := opts.resolveManifestURL(hit); manifestURL != "" {
		fields = append(fields, v4api.RecordField{Name: "iiif_manifest_url", Type: "iiif-manifest-url", Label: "IIIF Manifest", Value: manifestURL})
	}
	if hit.Content != "" {
		fields = append(fields, v4api.RecordField{Name: "summary", Type: "summary", Label: "Summary", Value: hit.Content, CitationPart: "abstract"})
	}
	return fields
}
