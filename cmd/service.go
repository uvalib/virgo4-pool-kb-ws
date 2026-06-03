package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/uvalib/virgo4-api/v4api"
	"github.com/uvalib/virgo4-jwt/v4jwt"
	"github.com/uvalib/virgo4-pool-kb-ws/internal/provider"
	"github.com/uvalib/virgo4-pool-kb-ws/internal/provider/bedrock"
	"github.com/uvalib/virgo4-pool-kb-ws/internal/service"
)

// ServiceContext contains shared handler dependencies.
type ServiceContext struct {
	Version            string
	JWTKey             string
	DetailResourceBase string
	ResourceHTTPClient *http.Client
	SearchService      *service.SearchService
}

// InitializeService wires dependencies from configuration.
func InitializeService(version string, cfg *ServiceConfig) *ServiceContext {
	var kbProvider provider.KnowledgeBaseProvider
	switch strings.ToLower(strings.TrimSpace(cfg.ProviderType)) {
	case "mock":
		log.Printf("using mock knowledge base provider")
		kbProvider = &provider.MockProvider{}
	case "", "bedrock":
		var err error
		kbProvider, err = bedrock.New(cfg.AWSRegion, cfg.KnowledgeBaseID)
		if err != nil {
			log.Fatalf("unable to initialize bedrock provider: %v", err)
		}
	default:
		log.Fatalf("unknown provider type %q (use bedrock or mock)", cfg.ProviderType)
	}

	svc := &ServiceContext{
		Version:            version,
		JWTKey:             cfg.JWTKey,
		DetailResourceBase: cfg.DetailResourceBase,
		ResourceHTTPClient: &http.Client{Timeout: 30 * time.Second},
		SearchService: &service.SearchService{
			Provider:         kbProvider,
			DefaultLimit:     cfg.DefaultLimit,
			ScoreThreshold:   cfg.ScoreThreshold,
			IIIFImageBaseURL: cfg.IIIFImageBaseURL,
		},
	}
	return svc
}

func (svc *ServiceContext) ignoreFavicon(c *gin.Context) {}

func (svc *ServiceContext) getVersion(c *gin.Context) {
	build := "unknown"
	files, _ := filepath.Glob("../buildtag.*")
	if len(files) == 1 {
		build = strings.Replace(files[0], "../buildtag.", "", 1)
	}
	c.JSON(http.StatusOK, gin.H{"version": svc.Version, "build": build})
}

func (svc *ServiceContext) healthCheck(c *gin.Context) {
	type hcResp struct {
		Healthy bool   `json:"healthy"`
		Message string `json:"message,omitempty"`
	}
	c.JSON(http.StatusOK, gin.H{"knowledge_base": hcResp{Healthy: true}})
}

func (svc *ServiceContext) identifyHandler(c *gin.Context) {
	resp := v4api.PoolIdentity{Attributes: make([]v4api.PoolAttribute, 0)}
	resp.Name = "Images (Knowledge Base)"
	resp.Description = "Imagery from the UVA Library digital collections knowledge base."
	resp.Mode = "image"
	resp.Attributes = append(resp.Attributes, v4api.PoolAttribute{Name: "facets", Supported: false})
	resp.Attributes = append(resp.Attributes, v4api.PoolAttribute{Name: "sorting", Supported: false})
	c.JSON(http.StatusOK, resp)
}

func getBearerToken(authorization string) (string, error) {
	components := strings.Split(strings.Join(strings.Fields(authorization), " "), " ")
	if len(components) != 2 || components[0] != "Bearer" || components[1] == "" {
		return "", fmt.Errorf("invalid Authorization header: [%s]", authorization)
	}
	return components[1], nil
}

func (svc *ServiceContext) authMiddleware(c *gin.Context) {
	tokenStr, err := getBearerToken(c.Request.Header.Get("Authorization"))
	if err != nil {
		log.Printf("authentication failed: %s", err.Error())
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if tokenStr == "undefined" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	claims, jwtErr := v4jwt.Validate(tokenStr, svc.JWTKey)
	if jwtErr != nil {
		log.Printf("jwt validation failed: %s", jwtErr.Error())
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	c.Set("jwt", tokenStr)
	c.Set("claims", claims)
}
