package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/uvalib/virgo4-api/v4api"
)

// isNamespacedResourceID reports whether path looks like a Virgo item id (e.g. uva-lib:123, tsb:59492).
func isNamespacedResourceID(id string) bool {
	id = strings.TrimSpace(id)
	if id == "" || strings.Contains(id, "/") {
		return false
	}
	prefix, suffix, ok := strings.Cut(id, ":")
	return ok && prefix != "" && suffix != ""
}

type providerDetails struct {
	Provider    string `json:"provider"`
	Label       string `json:"label,omitempty"`
	HomepageURL string `json:"homepage_url,omitempty"`
	LogoURL     string `json:"logo_url,omitempty"`
}

type poolProviders struct {
	Providers []providerDetails `json:"providers"`
}

func (svc *ServiceContext) providersHandler(c *gin.Context) {
	p := poolProviders{Providers: []providerDetails{}}
	c.JSON(http.StatusOK, p)
}

func (svc *ServiceContext) search(c *gin.Context) {
	var req v4api.SearchRequest
	if err := c.BindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, "invalid request")
		return
	}
	resp, status, err := svc.SearchService.Search(c.Request.Context(), &req)
	if err != nil {
		if status == http.StatusNotImplemented {
			c.String(status, err.Error())
			return
		}
		c.String(status, err.Error())
		return
	}
	c.JSON(status, resp)
}

func (svc *ServiceContext) facets(c *gin.Context) {
	c.JSON(http.StatusOK, svc.SearchService.Facets())
}

// getResource proxies GET /api/resource/:id to the Solr images pool.
func (svc *ServiceContext) getResource(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	if !isNamespacedResourceID(id) {
		c.String(http.StatusBadRequest, "missing resource id")
		return
	}
	svc.proxyResourceDetail(c, id)
}

// getResourceByID proxies GET /{namespace:id} (pool root identifier URLs) to the Solr images pool.
func (svc *ServiceContext) getResourceByID(c *gin.Context) {
	id := strings.TrimSpace(strings.TrimPrefix(c.Request.URL.Path, "/"))
	if !isNamespacedResourceID(id) {
		c.String(http.StatusNotFound, "not found")
		return
	}
	svc.proxyResourceDetail(c, id)
}
