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

func TestNoCache(t *testing.T) {
	// most duplicated text in test.
	const (
		test = "test"
		etag = "ETag"
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

	type noCacheHeaders struct {
		header string
		value  string
	}

	for _, tst := range [...]noCacheHeaders{
		{
			header: "Expires",
			value:  time.Unix(0, 0).Format(time.RFC1123),
		},
		{
			header: "Cache-Control",
			value:  "no-cache, no-store, no-transform, must-revalidate, private, max-age=0",
		},
		{
			header: "Pragma",
			value:  "no-cache",
		},
		{
			header: "X-Accel-Expires",
			value:  "0",
		},
	} {
		tst := tst

		t.Run(tst.header, func(t *testing.T) {
			require.Equal(t, w.Header().Get(tst.header), tst.value)
		})

		t.Run(tst.header, func(t *testing.T) {
			require.Equal(t, r.Header.Get(etag), "")
		})
	}
}
