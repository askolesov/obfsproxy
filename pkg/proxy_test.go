package pkg

import (
	"bytes"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/askolesov/obfsproxy/pkg/codec"
	"github.com/stretchr/testify/require"
)

func TestProxyChain(t *testing.T) {
	// Create a mock HTTP server that mirrors input
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.Copy(w, r.Body)
	}))
	defer mockServer.Close()

	mockServerAddr := mockServer.Listener.Addr().String()

	// Create codec for proxies
	key := []byte("test-key")
	xorer, err := codec.NewXorer(key)
	require.NoError(t, err)

	injector, err := codec.NewInjector(456, 20)
	require.NoError(t, err)

	chain, err := codec.NewChain([]codec.Codec{
		codec.NewInverter(),
		xorer,
		injector,
	})
	require.NoError(t, err)

	// Start proxy 1 (server)
	proxyServer := NewProxy("localhost:50502", mockServerAddr, true, chain)
	go func() {
		_ = proxyServer.Start()
	}()

	// Start proxy 2 (client)
	proxyClient := NewProxy("localhost:50503", proxyServer.ListenAddr, false, chain)
	go func() {
		_ = proxyClient.Start()
	}()

	// Generate test data (10MB)
	dataSize := 10 * 1024 * 1024 // 10MB
	testData := make([]byte, dataSize)
	for i := range testData {
		testData[i] = byte(i * 345876 % 256)
	}

	// Send data through proxy chain
	resp, err := http.Post("http://"+proxyClient.ListenAddr, "application/octet-stream", bytes.NewReader(testData))
	require.NoError(t, err)
	defer resp.Body.Close()

	// Read response
	result, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	// Compare input and output
	require.Equal(t, testData, result, "Response data should match input data")
}

func TestProxy_Proxy(t *testing.T) {
	// Create a proxy instance
	p := &Proxy{}

	// Create pipe connections
	in1, out1 := net.Pipe()
	defer in1.Close()
	defer out1.Close()
	in2, out2 := net.Pipe()
	defer in2.Close()
	defer out2.Close()

	// Test data
	testData := []byte("Hello, World!")
	expectedData := make([]byte, len(testData))
	for i, b := range testData {
		expectedData[i] = ^b // Invert bytes
	}

	codec := codec.NewInverter()

	// Run proxy in a goroutine
	go p.proxy(out1, in2, codec.NewEncoder())

	go func() {
		// Write test data to server
		_, err := in1.Write(testData)
		require.NoError(t, err, "Failed to write to server")
		in1.Close()
	}()

	// Read from client
	result := make([]byte, len(testData))
	_, err := io.ReadFull(out2, result)
	require.NoError(t, err, "Failed to read from client")

	// Compare result with expected data
	require.Equal(t, expectedData, result, "Proxy data mismatch")
}
