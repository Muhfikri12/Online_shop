package cmd

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ApiServer(logger *zap.Logger, port string, name string, engine *gin.Engine) error {
	if engine == nil {
		return errors.New("engine is nil")
	}

	// Validate port
	if _, err := strconv.Atoi(port); err != nil {
		return fmt.Errorf("invalid port: %v", err)
	}

	srv := &http.Server{Addr: ":" + port, Handler: engine}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("can't run service", zap.Error(err))
		}
	}()
	logger.Info(name + " initiated at port " + port)

	// gracefully shutdown ------------------------------------------------------------------------
	quit := make(chan os.Signal, 2)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	logger.Info("Shutdown " + name + " service")

	cts, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(cts); err != nil {
		logger.Error("can't shutdown "+name+" service", zap.Error(err))
	}

	logger.Info(name + " service exiting")

	logger.Info("Running cleanup tasks...")

	return nil
}
