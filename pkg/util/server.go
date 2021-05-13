package util

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"time"

	"github.com/couchbase/couchbase-exporter/pkg/log"
)

const (
	certAndKeyError = "please specify both cert and key arguments"
)

var (
	errCertAndKey = fmt.Errorf(certAndKeyError)
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
		log.Error("Server shutdown failed")
	} else {
		<-s.err
		s.Start(s.server.Addr, s.server.Handler)
	}
}

func (s *Server) RestartWithTLS(cert string, key string) {
	log.Info("Restarting server with TLS")

	if err := s.server.Shutdown(context.TODO()); err != nil {
		log.Error("Server shutdown failed")
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
		log.Error("%s", err)
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
			log.Error("Error when attempting to retrieve cert or key files: %s", err)
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
		log.Error("TLS Load failed %s", err)

		blankConfig := &tls.Config{}

		return blankConfig, err
	}

	return &tls.Config{
		Certificates: []tls.Certificate{keypair},
	}, nil
}

func Serve(address string, handler http.Handler, cert string, key string) {
	server := &Server{}

	server.pickStart(address, handler, cert, key)

	for {
		select {
		case err := <-server.err:
			log.Error("Server failed unexpectedly %s", err)
			server.pickStart(address, handler, cert, key)
		case <-time.After(10 * time.Second):
		}

		log.Debug("Restarting to pick up TLS config if necessary.")

		server.restartWithTLSIfNecessary(cert, key)
	}
}
