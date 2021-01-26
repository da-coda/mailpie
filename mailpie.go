package main

import (
	"embed"
	"fmt"
	"github.com/emersion/go-imap/server"
	"github.com/gorilla/mux"
	"github.com/mhale/smtpd"
	"github.com/r3labs/sse"
	"github.com/sirupsen/logrus"
	"mailpie/pkg/event"
	"mailpie/pkg/handler"
	"mailpie/pkg/handler/api"
	"mailpie/pkg/handler/imap"
	"mailpie/pkg/store"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type errorOrigin string

const (
	SMTP errorOrigin = "smtp"
	SPA  errorOrigin = "spa"
	API  errorOrigin = "api"
	SSE  errorOrigin = "sse"
	IMAP errorOrigin = "imap"
)

const listenOnAddress = "0.0.0.0"

type errorState struct {
	err    error
	origin errorOrigin
}

func main() {
	logrus.SetLevel(logrus.DebugLevel)

	globalMessageQueue := event.CreateOrGet()
	globalMailStore := store.CreateMailStore(*globalMessageQueue)

	errorChannel := make(chan errorState)
	go serveSPA(errorChannel)
	go serveSSE(errorChannel)

	smtpHandler := handler.CreateSmtpHandler(*globalMailStore)
	go serveSMTP(errorChannel, smtpHandler)
	go serveAPI(errorChannel)
	go serveIMAP(errorChannel)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signals
		fmt.Print("\r")
		logrus.Info("Received SIGTERM")
		os.Exit(0)
	}()

	var errorState errorState
	for {
		errorState = <-errorChannel
		logrus.WithError(errorState.err).WithField("Origin", errorState.origin).Error("Service received unexpected error")
	}
}

func serveSMTP(errorChannel chan errorState, smtpHandler handler.SmtpHandler) {
	addr := listenOnAddress + ":1025"
	srv := &smtpd.Server{
		Addr:         addr,
		Handler:      smtpHandler.Handle,
		Appname:      "Mailpie",
		Hostname:     "localhost",
		AuthRequired: false,
		AuthHandler: func(remoteAddr net.Addr, mechanism string, username []byte, password []byte, shared []byte) (bool, error) {
			return true, nil
		},
		LogWrite: func(remoteIP, verb, line string) {
			logrus.WithField("ip", remoteIP).WithField("verb", verb).Debug(line)
		},
		AuthMechs: map[string]bool{"PLAIN": true, "LOGIN": true, "CRAM-MD5": false},
	}
	logrus.WithField("Address", addr).Info("Starting SMTP server")
	err := srv.ListenAndServe()
	if err != nil {
		errorChannel <- errorState{err: err, origin: SMTP}
	}
}

func serveSPA(errorChannel chan errorState) {
	//go:embed "dist/index.html"
	var indexHtml string

	//go:embed "dist"
	var dist embed.FS
	router := mux.NewRouter()
	spa := handler.SpaHandler{
		Dist:  dist,
		Index: indexHtml,
	}
	router.PathPrefix("/").Handler(spa).Methods("GET")

	srv := &http.Server{
		Handler:      router,
		Addr:         listenOnAddress + ":8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logrus.WithField("Address", fmt.Sprintf("http://%s", srv.Addr)).Info("Starting SPA server")
	//should run forever unless an error occurs
	err := srv.ListenAndServe()
	if err != nil {
		errorChannel <- errorState{err: err, origin: SPA}
	}
}

func serveAPI(errorChannel chan errorState) {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/").Subrouter()
	api.RegisterApiRoutes(subrouter)
	srv := &http.Server{
		Handler:      router,
		Addr:         listenOnAddress + ":8001",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	logrus.WithField("Address", srv.Addr).Info("Starting API server")
	err := srv.ListenAndServe()
	if err != nil {
		errorChannel <- errorState{err: err, origin: API}
	}
}

func serveSSE(errorChannel chan errorState) {
	address := listenOnAddress + ":8002"
	sseServer := sse.New()
	sseServer.CreateStream("messages")
	sseHandler := handler.NewOrGetSSEHandler(sseServer)
	router := http.NewServeMux()
	router.Handle("/events", sseHandler)

	logrus.WithField("Address", address).Info("Starting SSE server")
	err := http.ListenAndServe(address, router)
	if err != nil {
		errorChannel <- errorState{err: err, origin: SSE}
	}
}

func serveIMAP(errorChannel chan errorState) {
	be := imap.NewBackend()
	s := server.New(be)
	imapLogger := logrus.New()
	s.Debug = imapLogger.Writer()
	s.Addr = listenOnAddress + ":1143"
	s.AllowInsecureAuth = true
	logrus.WithField("Address", s.Addr).Info("Starting IMAP server")
	err := s.ListenAndServe()
	if err != nil {
		errorChannel <- errorState{err: err, origin: IMAP}
	}
}

func (state errorState) String() string {
	return fmt.Sprintf("error at %s: %s", state.origin, state.err.Error())
}
