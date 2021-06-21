package main

//go:generate swagger generate spec -o ./swagger/swagger.yml
import (
	"embed"
	"flag"
	"fmt"
	"github.com/da-coda/mailpie/pkg/api"
	_ "github.com/da-coda/mailpie/pkg/api"
	"github.com/da-coda/mailpie/pkg/config"
	"github.com/da-coda/mailpie/pkg/event"
	"github.com/da-coda/mailpie/pkg/handler"
	"github.com/da-coda/mailpie/pkg/handler/imap"
	"github.com/da-coda/mailpie/pkg/store"
	"github.com/emersion/go-imap/server"
	"github.com/gorilla/mux"
	"github.com/mhale/smtpd"
	"github.com/sirupsen/logrus"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

type errorOrigin string

const (
	SMTP errorOrigin = "smtp"
	SPA  errorOrigin = "spa"
	IMAP errorOrigin = "imap"
	API  errorOrigin = "api"
)

type errorState struct {
	err    error
	origin errorOrigin
}

func main() {
	Run(flag.CommandLine, os.Args[1:])
}

//Run main entry point for mailpie. Loads the config, setup of globalMailStore and globalMessageQueue, starts all services
func Run(flags *flag.FlagSet, arguments []string) {
	err := config.Load(flags, arguments)
	if err != nil {
		logrus.WithError(err).Fatal("Error during configuration setup")
	}
	logrus.SetLevel(config.GetConfig().LogrusLevel)
	conf := config.GetConfig()
	globalMessageQueue := event.CreateOrGet()
	globalMailStore := store.CreateMailStore(*globalMessageQueue)

	errorChannel := make(chan errorState)
	if !conf.DisableHTTP {
		go serveSPA(errorChannel)
	}

	if !conf.DisableSMTP {
		smtpHandler := handler.CreateSmtpHandler(*globalMailStore)
		go serveSMTP(errorChannel, smtpHandler)
	}

	if !conf.DisableIMAP {
		go serveIMAP(errorChannel)
	}

	go serveApi(errorChannel, globalMailStore)

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

//serveSMTP Setup SMTP-Server and run ListenAndServe. If some error occurs during service runtime, the error gets send to Run
//via the errorChannel. Needs an SMTP handler which handles incoming mails
func serveSMTP(errorChannel chan errorState, smtpHandler handler.SmtpHandler) {
	addr := config.GetConfig().NetworkConfigs.SMTP.Host + ":" + strconv.Itoa(config.GetConfig().NetworkConfigs.SMTP.Port)
	srv := &smtpd.Server{
		Addr:         addr,
		Handler:      smtpHandler.Handle,
		Appname:      "Mailpie",
		Hostname:     "localhost",
		AuthRequired: false,
		//currently no auth is needed and implemented, so always return true on login
		AuthHandler: func(remoteAddr net.Addr, mechanism string, username []byte, password []byte, shared []byte) (bool, error) {
			return true, nil
		},
		LogWrite: func(remoteIP, verb, line string) {
			logrus.WithField("ip", remoteIP).WithField("verb", verb).Debug(line)
		},
		AuthMechs: map[string]bool{"PLAIN": true, "LOGIN": true, "CRAM-MD5": false},
	}
	logrus.WithField("Address", addr).Info("Starting SMTP server")
	//run the server. In best case, this will never stop. If there is some error, send it to Run via errorChannel
	err := srv.ListenAndServe()
	if err != nil {
		errorChannel <- errorState{err: err, origin: SMTP}
	}
}

//embed the index html and the dist directory(introduced in go 1.16)
//go:embed "dist/index.html"
var indexHtml string

//go:embed "dist"
var dist embed.FS

//serveSPA serve the MailPie Single-Page-Application
func serveSPA(errorChannel chan errorState) {
	router := mux.NewRouter()
	spa := handler.NewSpaHandler(dist, indexHtml)
	router.PathPrefix("/").Handler(spa).Methods("GET")

	srv := &http.Server{
		Handler:      router,
		Addr:         config.GetConfig().NetworkConfigs.HTTP.Host + ":" + strconv.Itoa(config.GetConfig().NetworkConfigs.HTTP.Port),
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

//serveIMAP runs the IMAP server
func serveIMAP(errorChannel chan errorState) {
	be := imap.NewBackend()
	s := server.New(be)
	imapLogger := logrus.StandardLogger()
	s.Debug = imapLogger.Writer()
	s.Addr = config.GetConfig().NetworkConfigs.IMAP.Host + ":" + strconv.Itoa(config.GetConfig().NetworkConfigs.IMAP.Port)
	s.AllowInsecureAuth = true
	logrus.WithField("Address", s.Addr).Info("Starting IMAP server")
	err := s.ListenAndServe()
	if err != nil {
		errorChannel <- errorState{err: err, origin: IMAP}
	}
}

func serveApi(errorChannel chan errorState, mailStore *store.MailStore) {
	baseRouter := api.NewBaseRouter(mailStore)
	srv := &http.Server{
		Handler:      baseRouter.Mux,
		Addr:         config.GetConfig().NetworkConfigs.API.Host + ":" + strconv.Itoa(config.GetConfig().NetworkConfigs.API.Port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logrus.WithField("Address", fmt.Sprintf("http://%s", srv.Addr)).Info("Starting API server")
	//should run forever unless an error occurs
	err := srv.ListenAndServe()
	if err != nil {
		errorChannel <- errorState{err: err, origin: API}
	}
}

func (state errorState) String() string {
	return fmt.Sprintf("error at %s: %s", state.origin, state.err.Error())
}
