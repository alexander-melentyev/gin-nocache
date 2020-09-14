package nocache_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	nocache "github.com/alexander-melentyev/gin-nocache"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

// most duplicated text in test.
const (
	test = "test"
	etag = "ETag"
)

func TestNoCache(t *testing.T) {
	var (
		epoch          = time.Unix(0, 0).Format(time.RFC1123)
		noCacheHeaders = map[string]string{
			"Expires":         epoch,
			"Cache-Control":   "no-cache, no-store, no-transform, must-revalidate, private, max-age=0",
			"Pragma":          "no-cache",
			"X-Accel-Expires": "0",
		}
	)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", test, nil)
	r.Header.Set(etag, test)

	gin.SetMode(gin.ReleaseMode)

	g := gin.New()
	g.Use(nocache.NoCache())
	g.GET(test, func(c *gin.Context) {
		c.JSON(200, gin.H{
			test: test,
		})
	})

	g.ServeHTTP(w, r)

	for k, v := range noCacheHeaders {
		t.Run(k, func(t *testing.T) {
			require.Equal(t, w.Header().Get(k), v)
		})

		t.Run(k, func(t *testing.T) {
			require.Equal(t, r.Header.Get(etag), "")
		})
	}
}
