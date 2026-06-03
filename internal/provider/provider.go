package provider

import "context"

// Hit is a normalized knowledge-base retrieval result.
type Hit struct {
	ID              string
	IIIFID          string
	IIIFImageURL    string
	IIIFManifestURL string
	Title           string
	Collection      string
	Score           float64
	Content         string
}

// KnowledgeBaseProvider retrieves records from a backing knowledge base.
type KnowledgeBaseProvider interface {
	Search(ctx context.Context, query string, limit int, threshold float64) ([]Hit, error)
}
