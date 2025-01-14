// Copyright 2017 Igor Dolzhikov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package service

import (
	"net/http"
	"sync/atomic"
	"time"

	"github.com/takama/bit"
	// Alternative of the Bit router with the same Router interface
	// "github.com/takama/k8sapp/pkg/router/httprouter"
	"github.com/vic999/k8sapp/pkg/config"
	"github.com/vic999/k8sapp/pkg/handlers"
	"github.com/vic999/k8sapp/pkg/logger"
	stdlog "github.com/vic999/k8sapp/pkg/logger/standard"
	"github.com/vic999/k8sapp/pkg/version"
)

// Setup configures the service
func Setup(cfg *config.Config) (r bit.Router, log logger.Logger, err error) {
	// Setup logger
	log = stdlog.New(&logger.Config{
		Level: cfg.LogLevel,
		Time:  true,
		UTC:   true,
	})

	log.Info("Version:", version.RELEASE)
	log.Warnf("%s log level is used", logger.LevelDebug.String())
	log.Infof("Service %s listened on %s:%d", config.SERVICENAME, cfg.LocalHost, cfg.LocalPort)
	isReady := &atomic.Value{}
	isReady.Store(false)
	go func() {
		log.Infof("Readyz probe is negative by default... waiting 30 sec")
		time.Sleep(30 * time.Second)
		isReady.Store(true)
		log.Infof("Readyz probe is positive.")
	}()

	// Define handlers
	h := handlers.New(log, cfg)

	// Register new router
	r = bit.NewRouter()

	// Response for undefined methods
	r.SetupNotFoundHandler(h.Base(notFound))

	// Configure router
	r.SetupMiddleware(h.Base)
	r.GET("/", h.Root)
	r.GET("/healthz", h.Health)
	r.GET("/readyz", h.ReadyWaiting(isReady))
	r.GET("/info", h.Info)

	return
}

// Response for undefined methods
func notFound(c bit.Control) {
	c.Code(http.StatusNotFound)
	c.Body("Method not found for " + c.Request().URL.Path)
}
