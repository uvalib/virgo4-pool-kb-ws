package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/uvalib/virgo4-api/v4api"
	"github.com/uvalib/virgo4-pool-kb-ws/internal/provider"
	"github.com/uvalib/virgo4-pool-kb-ws/internal/service"
)

func TestSearchHandlerReturnsPoolResult(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &ServiceContext{
		Version: "test",
		JWTKey:  "secret",
		SearchService: &service.SearchService{
			Provider:     &provider.MockProvider{Hits: []provider.Hit{{ID: "1", Title: "A"}}},
			DefaultLimit: 20,
		},
	}

	r := gin.New()
	r.POST("/api/search", svc.search)

	body, _ := json.Marshal(v4api.SearchRequest{Query: `keyword: {cat}`})
	req := httptest.NewRequest(http.MethodPost, "/api/search", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
	}
	var resp v4api.PoolResult
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if len(resp.Groups) != 1 {
		t.Fatalf("expected 1 group, got %d", len(resp.Groups))
	}
}

func TestGetResourceProxiesSolrPool(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var gotAuth string
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		if r.URL.Path != "/api/resource/uva-lib:2154074" {
			t.Errorf("path=%q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"fields":[{"name":"title","value":"Test"}]}`))
	}))
	defer upstream.Close()

	svc := &ServiceContext{
		DetailResourceBase: upstream.URL + "/api/resource",
		ResourceHTTPClient: upstream.Client(),
	}
	r := gin.New()
	r.GET("/api/resource/:id", svc.getResource)

	req := httptest.NewRequest(http.MethodGet, "/api/resource/uva-lib:2154074", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
	}
	if gotAuth != "Bearer test-token" {
		t.Fatalf("Authorization=%q", gotAuth)
	}
	if !bytes.Contains(w.Body.Bytes(), []byte(`"fields"`)) {
		t.Fatalf("body=%s", w.Body.String())
	}
}

func TestUvaLibRootPathProxiesSolrPool(t *testing.T) {
	gin.SetMode(gin.TestMode)
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/resource/uva-lib:329370" {
			t.Errorf("path=%q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"fields":[]}`))
	}))
	defer upstream.Close()

	svc := &ServiceContext{
		DetailResourceBase: upstream.URL + "/api/resource",
		ResourceHTTPClient: upstream.Client(),
	}
	r := gin.New()
	r.GET("/:id", svc.getResourceByID)

	req := httptest.NewRequest(http.MethodGet, "/uva-lib:329370", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
	}
}

func TestRootPathProxiesNamespacedID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/resource/tsb:59492" {
			t.Errorf("path=%q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"fields":[]}`))
	}))
	defer upstream.Close()

	svc := &ServiceContext{
		DetailResourceBase: upstream.URL + "/api/resource",
		ResourceHTTPClient: upstream.Client(),
	}
	r := gin.New()
	r.GET("/:id", svc.getResourceByID)

	req := httptest.NewRequest(http.MethodGet, "/tsb:59492", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
	}
}

func TestIsNamespacedResourceID(t *testing.T) {
	cases := map[string]bool{
		"uva-lib:329370": true,
		"tsb:59492":      true,
		"tsm:1503507":    true,
		"version":        false,
		"":               false,
		":59492":         false,
		"tsb:":           false,
		"no-colon":       false,
	}
	for id, want := range cases {
		if got := isNamespacedResourceID(id); got != want {
			t.Fatalf("isNamespacedResourceID(%q)=%v want %v", id, got, want)
		}
	}
}

func TestRootPathReturnsNotFoundForNonResourcePaths(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &ServiceContext{DetailResourceBase: "https://example.com/api/resource"}
	r := gin.New()
	r.GET("/:id", svc.getResourceByID)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/version", nil))
	if w.Code != http.StatusNotFound {
		t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
	}
}

func TestDetailResourceURLTrimsBaseSlash(t *testing.T) {
	got := detailResourceURL("https://example.com/items/", "id-1")
	want := "https://example.com/items/id-1"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestHealthcheck(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := &ServiceContext{Version: "test"}
	r := gin.New()
	r.GET("/healthcheck", svc.healthCheck)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/healthcheck", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("unexpected status %d", w.Code)
	}
}
