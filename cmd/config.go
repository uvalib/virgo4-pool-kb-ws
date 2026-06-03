package main

import (
	"flag"
	"log"
	"os"
	"strconv"
)

// ServiceConfig holds runtime configuration for the KB pool service.
type ServiceConfig struct {
	Port               int
	JWTKey             string
	AWSRegion          string
	KnowledgeBaseID    string
	ScoreThreshold     float64
	DefaultLimit       int
	ProviderType       string
	IIIFImageBaseURL   string
	DetailResourceBase string `yaml:"detail_resource_base"`
}

// LoadConfiguration applies precedence: env < flags.
func LoadConfiguration() *ServiceConfig {
	cfg := &ServiceConfig{}
	applyEnv(cfg)

	flag.IntVar(&cfg.Port, "port", cfg.Port, "Service port")
	flag.StringVar(&cfg.JWTKey, "jwtkey", cfg.JWTKey, "JWT signature key")
	flag.StringVar(&cfg.AWSRegion, "region", cfg.AWSRegion, "AWS region for Bedrock")
	flag.StringVar(&cfg.KnowledgeBaseID, "kbid", cfg.KnowledgeBaseID, "Bedrock knowledge base id")
	flag.Float64Var(&cfg.ScoreThreshold, "threshold", cfg.ScoreThreshold, "Minimum retrieval score")
	flag.IntVar(&cfg.DefaultLimit, "limit", cfg.DefaultLimit, "Default page size for retrieval")
	flag.StringVar(&cfg.ProviderType, "provider", cfg.ProviderType, "Knowledge base provider: bedrock or mock")
	flag.StringVar(&cfg.IIIFImageBaseURL, "iiifbase", cfg.IIIFImageBaseURL, "IIIF image API base URL (used when KB metadata has iiif_id only)")
	flag.StringVar(&cfg.DetailResourceBase, "detailresource", cfg.DetailResourceBase, "Solr images pool /api/resource URL base (appends /{id})")
	flag.String("config", "", "YAML config file (also V4_POOL_KB_CONFIG); path read before flags")

	flag.Parse()

	if cfg.JWTKey == "" {
		log.Fatal("jwtkey is required (flag -jwtkey, env V4_POOL_KB_JWT_KEY)")
	}

	log.Printf("detail resource base=%s", cfg.DetailResourceBase)
	log.Printf("[CONFIG] port=%d region=%s kb=%s threshold=%.3f limit=%d provider=%s", cfg.Port, cfg.AWSRegion, cfg.KnowledgeBaseID, cfg.ScoreThreshold, cfg.DefaultLimit, cfg.ProviderType)
	return cfg
}

func applyEnv(cfg *ServiceConfig) {
	if v := os.Getenv("V4_POOL_KB_PORT"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.Port = n
		}
	}
	if v := os.Getenv("V4_POOL_KB_JWT_KEY"); v != "" {
		cfg.JWTKey = v
	}
	if v := os.Getenv("V4_POOL_KB_AWS_REGION"); v != "" {
		cfg.AWSRegion = v
	}
	if v := os.Getenv("V4_POOL_KB_ID"); v != "" {
		cfg.KnowledgeBaseID = v
	}
	if v := os.Getenv("V4_POOL_KB_SCORE_THRESHOLD"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			cfg.ScoreThreshold = f
		}
	}
	if v := os.Getenv("V4_POOL_KB_PROVIDER"); v != "" {
		cfg.ProviderType = v
	}
	if v := os.Getenv("V4_POOL_KB_DEFAULT_LIMIT"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.DefaultLimit = n
		}
	}
	if v := os.Getenv("V4_POOL_KB_IIIF_IMAGE_BASE_URL"); v != "" {
		cfg.IIIFImageBaseURL = v
	}
	if v := os.Getenv("V4_KB_DETAIL_BASE"); v != "" {
		cfg.DetailResourceBase = v
	}
}
