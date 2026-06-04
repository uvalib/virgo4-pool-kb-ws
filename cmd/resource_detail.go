package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func detailResourceURL(base, id string) string {
	base = strings.TrimRight(strings.TrimSpace(base), "/")
	return base + "/" + strings.TrimSpace(id)
}

func requestJWT(c *gin.Context) string {
	if v, ok := c.Get("jwt"); ok {
		if token, ok := v.(string); ok && token != "" && token != "undefined" {
			return token
		}
	}
	token, err := getBearerToken(c.GetHeader("Authorization"))
	if err != nil {
		return ""
	}
	return token
}

func (svc *ServiceContext) fetchDetailResource(ctx context.Context, id, jwt string) (int, []byte, string, error) {
	if svc.DetailResourceBase == "" {
		return http.StatusInternalServerError, nil, "", fmt.Errorf("detail resource base not configured (set VIRGO_DETAIL_SOURCE)")
	}

	tgt := detailResourceURL(svc.DetailResourceBase, id)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, tgt, nil)
	if err != nil {
		return http.StatusInternalServerError, nil, "", err
	}
	req.Header.Set("Authorization", "Bearer "+jwt)

	client := svc.ResourceHTTPClient
	if client == nil {
		client = http.DefaultClient
	}

	log.Printf("resource %s fetching from %s", id, tgt)
	resp, err := client.Do(req)
	if err != nil {
		return http.StatusBadGateway, nil, "", fmt.Errorf("upstream request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return http.StatusBadGateway, nil, "", fmt.Errorf("reading upstream response: %w", err)
	}
	return resp.StatusCode, body, resp.Header.Get("Content-Type"), nil
}

func (svc *ServiceContext) proxyResourceDetail(c *gin.Context, id string) {
	jwt := requestJWT(c)
	if jwt == "" {
		c.String(http.StatusUnauthorized, "authentication required")
		return
	}

	status, body, contentType, err := svc.fetchDetailResource(c.Request.Context(), id, jwt)
	if err != nil {
		if status == 0 {
			status = http.StatusBadGateway
		}
		c.String(status, err.Error())
		return
	}

	if contentType == "" {
		contentType = "application/json"
	}
	c.Data(status, contentType, body)
}
