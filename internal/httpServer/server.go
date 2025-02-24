package httpServer

import (
	"crypto/tls"
	"fmt"
	"net/http"
)

type HttpServer struct {
	address     string
	enableHttps bool
	mux         *http.ServeMux
	controllers map[string]Controller
}

func NewServer(address string, enableHttps bool) *HttpServer {
	return &HttpServer{
		address:     address,
		mux:         http.NewServeMux(),
		controllers: make(map[string]Controller),
		enableHttps: enableHttps,
	}
}

func (s *HttpServer) AddController(ctrl Controller) error {
	if _, ok := s.controllers[ctrl.GetUrl()]; ok {
		return fmt.Errorf("presentation with url \"%s\" already added", ctrl.GetUrl())
	}
	s.controllers[ctrl.GetUrl()] = ctrl
	return nil
}

func (s *HttpServer) Launch() <-chan error {
	errChan := make(chan error, 1)
	go func() {
		defer close(errChan)
		var err error
		defer func() { errChan <- err }()
		s.registerControllers()
		server := http.Server{
			Addr:    s.address,
			Handler: s.mux,
		}
		if err = s.attachCertificate(&server); err != nil {
			errChan <- err
			return
		}
		fmt.Println("starting http server on", s.address)
		if err = server.ListenAndServe(); err != nil {
			errChan <- err
			return
		}
	}()
	return errChan
}

func (s *HttpServer) registerControllers() {
	for _, ctrl := range s.controllers {
		s.mux.HandleFunc(ctrl.GetUrl(), func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
				}
			}()
			for _, midWare := range ctrl.GetOrderedMiddlewares() {
				if !midWare.ProcessBeforeHandler(w, r) {
					return
				}
			}
			ctrl.Handle(w, r)
			for _, midWare := range ctrl.GetOrderedMiddlewares() {
				midWare.ProcessAfterHandler(w, r)
			}
		})
	}
}

func (s *HttpServer) attachCertificate(server *http.Server) error {
	if !s.enableHttps {
		return nil
	}
	certFilePath := "" //path to the cert file
	keyFilePath := ""  //path to the key file
	certificate, err := tls.LoadX509KeyPair(certFilePath, keyFilePath)
	if err != nil {
		return err
	}
	server.TLSConfig = &tls.Config{
		Certificates: []tls.Certificate{certificate},
	}
	return nil
}
