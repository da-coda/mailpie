package main

import (
	"fmt"
	"github.com/emersion/go-imap/backend/memory"
	"github.com/emersion/go-imap/server"
	"github.com/gorilla/mux"
	"github.com/mhale/smtpd"
	"github.com/r3labs/sse"
	"log"
	"mailmon-go/pkg/handler"
	"mailmon-go/pkg/handler/api"
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

type errorState struct {
	err    error
	origin errorOrigin
}

func main() {
	errorChannel := make(chan errorState)
	go serveSPA(errorChannel)
	go serveSSE(errorChannel)
	go serveSMTP(errorChannel)
	go serveAPI(errorChannel)
	go serveIMAP(errorChannel)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signals
		fmt.Println("Received SIGTERM")
		os.Exit(0)
	}()

	var errorState errorState
	for {
		errorState = <-errorChannel
		fmt.Println(errorState)
	}
}

func serveSMTP(errorChannel chan errorState) {
	addr := "127.0.0.1:1025"
	srv := &smtpd.Server{
		Addr:         addr,
		Handler:      handler.SmtpHandler,
		Appname:      "Mailmon",
		Hostname:     "localhost",
		AuthRequired: false,
		AuthHandler: func(remoteAddr net.Addr, mechanism string, username []byte, password []byte, shared []byte) (bool, error) {
			return true, nil
		},
		AuthMechs: map[string]bool{"PLAIN": true, "LOGIN": true},
	}
	log.Println("Starting smtp server at: ", addr)
	err := srv.ListenAndServe()
	if err != nil {
		errorChannel <- errorState{err: err, origin: SMTP}
	}
}

func serveSPA(errorChannel chan errorState) {
	router := mux.NewRouter()
	spa := handler.SpaHandler{StaticPath: "dist", IndexPath: "dist/index.html"}
	router.PathPrefix("/").Handler(spa).Methods("GET")

	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("Starting SPA server at: ", srv.Addr)
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
		Addr:         "127.0.0.1:8001",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Println("Starting API server at: ", srv.Addr)
	err := srv.ListenAndServe()
	if err != nil {
		errorChannel <- errorState{err: err, origin: API}
	}
}

func serveSSE(errorChannel chan errorState) {
	sseServer := sse.New()
	sseServer.CreateStream("messages")
	sseHandler := handler.NewOrGetSSEHandler(sseServer)
	router := http.NewServeMux()
	router.Handle("/events", sseHandler)

	log.Println("Starting SSE server at: 127.0.0.1:8002")
	err := http.ListenAndServe("127.0.0.1:8002", router)
	if err != nil {
		errorChannel <- errorState{err: err, origin: SSE}
	}
}

func serveIMAP(errorChannel chan errorState) {
	// Create a memory backend
	be := memory.New()

	s := server.New(be)
	s.Addr = ":1143"
	s.AllowInsecureAuth = true
	log.Println("Starting IMAP server at 127.0.0.1:1143")
	err := s.ListenAndServe()
	if err != nil {
		errorChannel <- errorState{err: err, origin: IMAP}
	}
}

func (state errorState) String() string {
	return fmt.Sprintf("error at %s: %s", state.origin, state.err.Error())
}
