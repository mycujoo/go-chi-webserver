package middleware

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/coocood/freecache"
)

type CacheHandler struct {
	Storage *freecache.Cache
}

func CreateCacheHandler(cacheSize int) *CacheHandler {
	if cacheSize == 0 {
		cacheSize = 100 * 1024 * 1024 // default cache size
	}

	return &CacheHandler{
		Storage: freecache.NewCache(cacheSize),
	}
}

func (c *CacheHandler) Cache(ttl int) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sw := responseWrapper{ResponseWriter: w}

			next.ServeHTTP(&sw, r)

			// Cache only successful requests
			if sw.statusCode == 200 {
				cacheKey := generateCacheKey(r.RequestURI)
				obj := map[string]interface{}{
					"body": string(sw.body),
					"ttl":  strconv.Itoa(ttl),
				}

				cached, err := json.Marshal(obj)

				if err != nil {
					log.Printf("Error parsing cache body: %s", err.Error())
				}

				if err := c.Storage.Set(cacheKey, cached, ttl); err != nil {
					log.Printf("Error when creating cache for '%s'. %s.", r.RequestURI, err)
				}
			}
		})
	}
}

func (c *CacheHandler) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cacheKey := generateCacheKey(r.RequestURI)
		cacheControlHeader := r.Header.Get("Cache-control")

		if !strings.Contains(cacheControlHeader, "no-cache") {
			obj, expireAt, err := c.Storage.GetWithExpiration(cacheKey)
			cached := map[string]interface{}{}

			json.Unmarshal(obj, &cached)

			if err == nil {
				now := time.Now().Unix()
				ttl, _ := strconv.Atoi(cached["ttl"].(string))
				age := ttl - (int(expireAt) - int(now))

				w.Header().Set("Age", strconv.Itoa(age))
				w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", ttl))
				w.Header().Set("Status", strconv.Itoa(http.StatusOK))
				w.WriteHeader(http.StatusOK)

				w.Write([]byte(cached["body"].(string)))

				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func generateCacheKey(uri string) []byte {
	return []byte(fmt.Sprintf("cache:%x", md5.Sum([]byte(uri))))
}

type responseWrapper struct {
	http.ResponseWriter
	body       []byte
	statusCode int
}

func (w *responseWrapper) Write(b []byte) (int, error) {
	w.body = b
	return w.ResponseWriter.Write(b)
}

func (w *responseWrapper) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
