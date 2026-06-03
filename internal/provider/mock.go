package provider

import "context"

// MockProvider is an in-memory provider for tests.
type MockProvider struct {
	Hits []Hit
	Err  error
}

func (m *MockProvider) Search(ctx context.Context, query string, limit int, threshold float64) ([]Hit, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	if limit <= 0 || limit > len(m.Hits) {
		limit = len(m.Hits)
	}
	return m.Hits[:limit], nil
}
