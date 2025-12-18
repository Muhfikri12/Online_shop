package cmd

import (
	"net/http"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestApiServer(t *testing.T) {
	t.Run("Success_Server_Start_And_Shutdown", func(t *testing.T) {
		// Setup
		logger := zap.NewNop()
		gin.SetMode(gin.TestMode)
		app := gin.New()
		port := "8081"
		name := "test-service"

		// Start server in a goroutine
		go ApiServer(logger, port, name, app)

		// Give the server time to start
		time.Sleep(100 * time.Millisecond)

		// Test if server is running by making a request
		resp, err := http.Get("http://localhost:" + port + "/")
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		if resp != nil {
			resp.Body.Close()
		}

		// Trigger graceful shutdown
		process, err := os.FindProcess(os.Getpid())
		assert.NoError(t, err)
		assert.NoError(t, process.Signal(syscall.SIGTERM))

		// Give the server time to shutdown
		time.Sleep(100 * time.Millisecond)

		// Verify server has shut down by attempting another request
		_, err = http.Get("http://localhost:" + port + "/")
		assert.Error(t, err)
	})

	t.Run("Nil_Engine", func(t *testing.T) {
		// Setup
		logger := zap.NewNop()
		port := "8082"
		name := "test-service"

		// Create a channel to capture the error
		errChan := make(chan error, 1)

		// Start server with nil engine
		go func() {
			err := ApiServer(logger, port, name, nil)
			errChan <- err
		}()

		// Give the server time to attempt to start
		time.Sleep(100 * time.Millisecond)

		// Verify that the server failed to start
		select {
		case err := <-errChan:
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "engine is nil")
		default:
			// If no error was received, the test should fail
			t.Fatal("Expected server to fail with nil engine error")
		}

		// Cleanup - ensure server is stopped
		process, _ := os.FindProcess(os.Getpid())
		process.Signal(syscall.SIGTERM)
	})
}
