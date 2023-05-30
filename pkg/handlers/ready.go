// Copyright 2017 Igor Dolzhikov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package handlers

import (
	"github.com/takama/bit"
	"net/http"
	"sync/atomic"
	// Alternative of the Bit router with the same Router interface
	// "github.com/takama/k8sapp/pkg/router/httprouter"
)

// Ready returns "OK" if service is ready to serve traffic
func (h *Handler) Ready(c bit.Control) {
	// TODO: possible use cases:
	// load data from a database, a message broker, any external services, etc

	c.Code(http.StatusOK)
	c.Body(http.StatusText(http.StatusOK))
}

// Base handler implements middleware logic
func (h *Handler) ReadyWaiting(isReady *atomic.Value) func(bit.Control) {
	return func(c bit.Control) {
		if isReady == nil || !isReady.Load().(bool) {
			c.Code(http.StatusServiceUnavailable)
			c.Body(http.StatusText(http.StatusServiceUnavailable))
			return
		}
		c.Code(http.StatusOK)
		c.Body(http.StatusText(http.StatusOK))
	}
}
