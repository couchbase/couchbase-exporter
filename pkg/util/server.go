package util

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"time"

	"github.com/couchbase/couchbase-exporter/pkg/objects"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	server *http.Server
	err    chan error
}

func (s *Server) Start(address string, handler http.Handler) {
	s.server = &http.Server{
		Addr:    address,
		Handler: handler,
	}

	s.err = make(chan error)

	go func() {
		s.err <- s.server.ListenAndServe()
	}()
}

func (s *Server) StartWithTLS(address string, handler http.Handler, cert string, key string) {
	s.server = &http.Server{
		Addr:    address,
		Handler: handler,
	}

	s.err = make(chan error)

	go func() {
		s.err <- s.server.ListenAndServeTLS(cert, key)
	}()
}

func (s *Server) Restart() {
	log.Info("Restarting server")

	if err := s.server.Shutdown(context.TODO()); err != nil {
		log.Error(err, "Server shutdown failed")
	} else {
		<-s.err
		s.Start(s.server.Addr, s.server.Handler)
	}
}

func (s *Server) RestartWithTLS(cert string, key string) {
	log.Info("Restarting server with TLS")

	if err := s.server.Shutdown(context.TODO()); err != nil {
		log.Error(err, "Server shutdown failed")
	} else {
		<-s.err
		s.StartWithTLS(s.server.Addr, s.server.Handler, cert, key)
	}
}

func useTLS(cert string, key string) (bool, error) {
	if len(cert) == 0 && len(key) == 0 {
		return false, nil
	}

	if len(cert) != 0 && len(key) != 0 {
		return true, nil
	}

	return false, errCertAndKey
}

func (s *Server) pickStart(address string, handler http.Handler, cert string, key string) {
	useTLS, err := useTLS(cert, key)

	if err != nil {
		log.Error(err, "error setting TLS certificates")
		os.Exit(1)
	}

	if useTLS {
		s.StartWithTLS(address, handler, cert, key)
	} else {
		s.Start(address, handler)
	}
}

func (s *Server) restartWithTLSIfNecessary(cert string, key string) {
	if tls, err := useTLS(cert, key); err == nil && tls {
		oldConfig := s.server.TLSConfig
		newConfig, err := createTLSConfig(cert, key)

		if err != nil {
			log.Error(err, "Error when attempting to retrieve cert or key files")

			return
		}

		if reflect.DeepEqual(oldConfig.Certificates, newConfig.Certificates) {
			return
		}

		s.RestartWithTLS(cert, key)
	}
}

func createTLSConfig(cert string, key string) (*tls.Config, error) {
	keypair, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		log.Error(err, "TLS Load failed")

		blankConfig := &tls.Config{} // nolint:gosec

		return blankConfig, fmt.Errorf("TLS Load Failed: %w", err)
	}

	return &tls.Config{ //nolint:gosec
		Certificates: []tls.Certificate{keypair},
	}, nil
}

const waitTimeForTLSCheck = 10

func Serve(exporterConfig *objects.ExporterConfig) {
	handler := AuthHandler{
		ServeMux:      http.NewServeMux(),
		TokenLocation: "",
	}

	if len(exporterConfig.Token) != 0 {
		handler.TokenLocation = exporterConfig.Token
	}

	handler.ServeMux.Handle("/metrics", promhttp.Handler())

	metricsServer := fmt.Sprintf("%v:%v", exporterConfig.ServerAddress, exporterConfig.ServerPort)
	log.Info("starting server", "listeningOn", metricsServer)

	server := &Server{
		server: nil,
		err:    make(chan error),
	}

	server.pickStart(metricsServer, handler, exporterConfig.Certificate, exporterConfig.Key)

	for {
		select {
		case err := <-server.err:
			log.Error(err, "Server failed unexpectedly")
			server.pickStart(metricsServer, handler, exporterConfig.Certificate, exporterConfig.Key)
		case <-time.After(waitTimeForTLSCheck * time.Second):
		}

		log.Info("Restarting to pick up TLS config if necessary.")

		server.restartWithTLSIfNecessary(exporterConfig.Certificate, exporterConfig.Key)
	}
}
