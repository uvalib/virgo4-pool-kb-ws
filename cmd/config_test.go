package main

import (
	"os"
	"testing"
)


func TestApplyEnv(t *testing.T) {
	t.Setenv("V4_POOL_KB_PORT", "7777")
	t.Setenv("V4_POOL_KB_JWT_KEY", "test-key")
	t.Setenv("V4_POOL_KB_AWS_REGION", "test-region")
	t.Setenv("V4_POOL_KB_ID", "KBTEST")
	t.Setenv("V4_POOL_KB_SCORE_THRESHOLD", "0.5")
	t.Setenv("V4_POOL_KB_DEFAULT_LIMIT", "10")
	t.Setenv("V4_POOL_KB_PROVIDER", "bedrock")
	t.Setenv("V4_POOL_KB_IIIF_IMAGE_BASE_URL", "https://iiif.lib.virginia.edu/iiif")
	cfg := &ServiceConfig{}
	applyEnv(cfg)
	if cfg.Port != 7777 || cfg.KnowledgeBaseID != "KBTEST" {
		t.Fatalf("unexpected cfg: %+v", cfg)
	}
	if cfg.JWTKey != "test-key" {
		t.Fatalf("jwt=%q want test-key", cfg.JWTKey)
	}
	if cfg.AWSRegion != "test-region" {
		t.Fatalf("region=%q want test-region", cfg.AWSRegion)
	}
	if cfg.ScoreThreshold != 0.5 {
		t.Fatalf("score threshold=%f want 0.5", cfg.ScoreThreshold)
	}
	if cfg.DefaultLimit != 10 {
		t.Fatalf("default limit=%d want 10", cfg.DefaultLimit)
	}
	if cfg.ProviderType != "bedrock" {
		t.Fatalf("provider=%q want bedrock", cfg.ProviderType)
	}
	if cfg.IIIFImageBaseURL != "https://iiif.lib.virginia.edu/iiif" {
		t.Fatalf("iiif image base url=%q want https://iiif.lib.virginia.edu/iiif", cfg.IIIFImageBaseURL)
	}
}

func TestDetailBaseEnv(t *testing.T) {
	t.Setenv("V4_KB_DETAIL_BASE", "http://localhost:8983/api/resource")
	cfg := &ServiceConfig{}
	applyEnv(cfg)
	if cfg.DetailResourceBase != "http://localhost:8983/api/resource" {
		t.Fatalf("detail base=%q", cfg.DetailResourceBase)
	}
}

func TestFlagOverridesEnv(t *testing.T) {
	t.Setenv("V4_POOL_KB_PORT", "8888")
	t.Setenv("V4_POOL_KB_JWT_KEY", "env-key")
	t.Setenv("V4_POOL_KB_AWS_REGION", "env-region")

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"cmd.test", "-port", "3001", "-jwtkey", "flag-key", "-region", "flag-region"}

	cfg := LoadConfiguration()
	if cfg.Port != 3001 {
		t.Fatalf("port=%d want flag 3001", cfg.Port)
	}
	if cfg.JWTKey != "flag-key" {
		t.Fatalf("jwt=%q want flag-key", cfg.JWTKey)
	}
	if cfg.AWSRegion != "flag-region" {
		t.Fatalf("region=%q want file-region (file below env, flag only set port/jwt)", cfg.AWSRegion)
	}
}
