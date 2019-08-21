package main

import (
	"crypto/sha256"
	"fmt"
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

type Response struct {
	Sha      string `json:"sha"`
	Filename string `json:"filename"`
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
		pluginPath := filepath.Join(path, plugin)
		fstat, err := os.Stat(pluginPath)
		if err != nil {
			c.JSON(http.StatusNotFound, err.Error())
			return
		}
		record, ok := cache.Load(pluginPath)
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
		cache.Store(pluginPath, CacheRecord{lastUpdate: fstat.ModTime(), sha: sha})
		log.Printf("%x", sha)
		c.JSON(http.StatusOK, Response{Sha: fmt.Sprintf("%x", sha), Filename: plugin})
	})
	r.Run(":8080") // listen and serve on 0.0.0.0:8080
}
