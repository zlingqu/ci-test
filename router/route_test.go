package router

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var r *gin.Engine

func init() {
	r = SetupRouter()
}

func TestGetHostname(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/getHostname", nil)

	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	hostname, _ := os.Hostname()
	assert.Equal(t, hostname, w.Body.String())

}

func TestPing(t *testing.T) {

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/ping", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}
