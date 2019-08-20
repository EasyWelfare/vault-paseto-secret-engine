package main

import (
	"crypto/sha256"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type CacheRecord struct {
	lastUpdate time.Time
	sha        []byte
}

func main() {
	path := os.Getenv("PLUGIN_PATH")
	if path == "" {
		log.Fatalf("PLUGIN_PATH env var not set")
	}
	cache := sync.Map{}
	r := gin.Default()
	r.Use(cors.Default())
	r.GET("/sha/:plugin", func(c *gin.Context) {
		plugin := c.Param("plugin")
		plugin_path := filepath.Join(path, plugin)
		fstat, err := os.Stat(plugin_path)
		if err != nil {
			c.JSON(http.StatusNotFound, err.Error())
			return
		}
		record, ok := cache.Load(plugin_path)
		if ok {
			record, okcast := record.(CacheRecord)
			if !okcast {
				log.Fatalf("not able to cast %v of type %T to CacheRecord", record, record)
			}
			if !record.lastUpdate.Before(fstat.ModTime()) {
				log.Printf("%x", record.sha)
				c.JSON(http.StatusOK, string(record.sha))
				return
			}
		}
		f, err := os.Open(filepath.Join(path, plugin))
		defer f.Close()
		h := sha256.New()
		if err != nil {
			c.JSON(http.StatusNotFound, err.Error())
			return
		}
		if _, err := io.Copy(h, f); err != nil {
			log.Fatal(err)
		}
		sha := h.Sum(nil)
		cache.Store(plugin_path, CacheRecord{lastUpdate: fstat.ModTime(), sha: sha})
		log.Printf("%x", sha)
		c.JSON(http.StatusOK, string(sha))
	})
	r.Run(":8080") // listen and serve on 0.0.0.0:8080
}
