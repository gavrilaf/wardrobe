package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo"

	"github.com/gavrilaf/wardrobe/pkg/utils/httpx"
	"github.com/gavrilaf/wardrobe/pkg/utils/log"
)

var HealthCheck = "/healthz"

func Measure(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()

		// Do nothing for health-check
		if req.RequestURI == HealthCheck {
			return next(c)
		}

		ctx := newContext(c, req)

		req = req.WithContext(ctx)
		c.SetRequest(req)

		log.FromContext(ctx).Infof("request received")

		start := time.Now()

		err := next(c)

		if err != nil {
			c.Error(err)
		}

		latency := time.Since(start)
		status := c.Response().Status

		logger := log.FromContext(ctx).With("status", status, "latency", latency.Milliseconds())

		if httpx.IsHttpStatusOk(status) {
			logger.Info("request completed")
		} else {
			if err == nil {
				err = fmt.Errorf("http error %d", status)
			}
			log.WithError(logger, err).Error("request completed with error")
		}

		return nil
	}
}

// TODO: configurable log fields
func newContext(c echo.Context, req *http.Request) context.Context {
	ctx := req.Context()

	ctx, _ = log.UpdateContext(ctx, "method", req.Method,
		"uri", req.RequestURI,
		"path", c.Path())

	return ctx
}
