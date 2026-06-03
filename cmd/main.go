package main

import (
	"fmt"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

const version = "1.0.0"

func main() {
	log.Printf("===> V4 Knowledge Base pool starting <===")
	cfg := LoadConfiguration()
	svc := InitializeService(version, cfg)

	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()
	router := gin.Default()
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	corsCfg := cors.DefaultConfig()
	corsCfg.AllowAllOrigins = true
	corsCfg.AllowCredentials = true
	corsCfg.AddAllowHeaders("Authorization")
	router.Use(cors.New(corsCfg))

	router.GET("/", svc.getVersion)
	router.GET("/favicon.ico", svc.ignoreFavicon)
	router.GET("/version", svc.getVersion)
	router.GET("/healthcheck", svc.healthCheck)
	router.GET("/identify", svc.identifyHandler)

	api := router.Group("/api")
	{
		api.GET("/providers", svc.providersHandler)
		api.POST("/search", svc.authMiddleware, svc.search)
		api.POST("/search/facets", svc.authMiddleware, svc.facets)
		api.GET("/resource/:id", svc.authMiddleware, svc.getResource)
	}

	router.GET("/:id", svc.getResourceByID)

	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("listening on %s (v%s)", addr, version)
	log.Fatal(router.Run(addr))
}
