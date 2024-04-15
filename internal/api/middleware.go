package api

import (
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/uptrace/bunrouter"
)

func (s *Server) middlewareLog(next bunrouter.HandlerFunc) bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		defer func() {
			s.logger.Logger.WithFields(logrus.Fields{
				"request_time":   time.Now().Format("2006-01-02 15:04:05.000000"),
				"request_ip":     strings.TrimPrefix(strings.Split(req.RemoteAddr, ":")[1], "["),
				"request_method": req.Method,
				"request_url":    req.URL,
			}).Debug("http request served")
		}()

		return next(w, req)
	}
}
