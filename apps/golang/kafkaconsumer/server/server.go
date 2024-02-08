package server

import (
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/logger"
	"github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/mysql"
)

type Server struct {
	logger *logger.Logger
	MySql  *mysql.MySqlDatabase
}

func New(
	log *logger.Logger,
	db *mysql.MySqlDatabase,
) *Server {
	return &Server{
		logger: log,
		MySql:  db,
	}
}

// Liveness
func (s *Server) Livez(
	w http.ResponseWriter,
	r *http.Request,
) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// Readiness
func (s *Server) Readyz(
	w http.ResponseWriter,
	r *http.Request,
) {
	err := s.MySql.Instance.Ping()
	if err != nil {
		s.logger.Log(logrus.ErrorLevel, r.Context(), "", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Not OK"))
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}
