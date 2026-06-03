package service

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/uvalib/virgo4-api/v4api"
	"github.com/uvalib/virgo4-pool-kb-ws/internal/parser"
	"github.com/uvalib/virgo4-pool-kb-ws/internal/provider"
	"github.com/uvalib/virgo4-pool-kb-ws/internal/transform"
)

// SearchService orchestrates validation, query translation, retrieval, and mapping.
type SearchService struct {
	Provider         provider.KnowledgeBaseProvider
	DefaultLimit     int
	ScoreThreshold   float64
	IIIFImageBaseURL string
}

func (s *SearchService) transformOpts() transform.Options {
	return transform.Options{IIIFImageBaseURL: s.IIIFImageBaseURL}
}

// Search runs a pool search for the given Virgo request.
func (s *SearchService) Search(ctx context.Context, req *v4api.SearchRequest) (*v4api.PoolResult, int, error) {
	start := time.Now()
	if req == nil {
		return nil, http.StatusBadRequest, fmt.Errorf("missing request")
	}

	valid, errMsg := parser.Validate(req.Query)
	if !valid {
		return nil, http.StatusBadRequest, fmt.Errorf("malformed search: %s", errMsg)
	}

	if hasUnsupportedFilters(req) || strings.Contains(req.Query, "filter:") {
		empty := transform.HitsToPoolResult(nil, req.Pagination.Start, 0, s.transformOpts())
		empty.StatusCode = http.StatusOK
		return empty, http.StatusOK, nil
	}

	for _, token := range []string{"date:", "identifier:", "journal_title:", "fulltext:", "series:"} {
		if strings.Contains(req.Query, token) {
			return nil, http.StatusNotImplemented, fmt.Errorf("%s queries are not supported", strings.TrimSuffix(token, ":"))
		}
	}

	limit := s.DefaultLimit
	if req.Pagination.Rows > 0 {
		limit = req.Pagination.Rows
	}
	kbQuery := parser.ToKBTextQuery(req.Query)
	hits, err := s.Provider.Search(ctx, kbQuery, limit, s.ScoreThreshold)
	if err != nil {
		resp := &v4api.PoolResult{StatusCode: http.StatusBadGateway, StatusMessage: err.Error(), Confidence: "low"}
		return resp, http.StatusBadGateway, nil
	}

	elapsed := time.Since(start).Milliseconds()
	resp := transform.HitsToPoolResult(hits, req.Pagination.Start, elapsed, s.transformOpts())
	return resp, http.StatusOK, nil
}

// Facets returns an empty facet list (KB pool does not support facets).
func (s *SearchService) Facets() map[string]any {
	out := make(map[string]any)
	out["facets"] = make([]v4api.Facet, 0)
	return out
}

// Resource loads a single record by id using a targeted KB query.
// NOT USED. Solr images pool handles this for now.
func (s *SearchService) Resource(ctx context.Context, id string) ([]v4api.RecordField, int, error) {
	hits, err := s.Provider.Search(ctx, id, 10, s.ScoreThreshold)
	if err != nil {
		return nil, http.StatusBadGateway, err
	}
	for _, hit := range hits {
		if hit.ID == id || hit.IIIFID == id {
			return transform.HitToDetailFields(hit, s.transformOpts()), http.StatusOK, nil
		}
	}
	return nil, http.StatusNotFound, fmt.Errorf("resource not found")
}

func hasUnsupportedFilters(req *v4api.SearchRequest) bool {
	if len(req.Filters) > 1 {
		return true
	}
	if len(req.Filters) == 1 {
		return len(req.Filters[0].Facets) > 0
	}
	return false
}
