package integration_test

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProxyIntegration(t *testing.T) {
	directResp, err := http.Get("http://localhost:8000")
	require.NoError(t, err)
	defer directResp.Body.Close()

	directBody, err := io.ReadAll(directResp.Body)
	require.NoError(t, err)

	proxyResp, err := http.Get("http://localhost:8080")
	require.NoError(t, err)
	defer proxyResp.Body.Close()

	proxyBody, err := io.ReadAll(proxyResp.Body)
	require.NoError(t, err)

	assert.Equal(t, directBody, proxyBody, "Response from direct request and proxy chain should be the same")
}
