package main

import (
	"flag"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"golang.org/x/net/html"
	"gopkg.in/mail.v2"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"
)

type MailpieTest struct {
	suite.Suite
}

func (suite *MailpieTest) TestRun() {
	tmp := os.TempDir()
	confFile, err := os.CreateTemp(tmp, "mailpieTestRun")
	if err != nil {
		suite.T().Skip("Unable to create temp config file for test")
	}
	configFileContent := []byte(
		`
networkconfigs:
    smtp:
        host: 0.0.0.0
        port: 1025
    imap:
        host: 0.0.0.0
        port: 1143
    http:
        host: 127.0.0.1
        port: 8000
disable_imap: false
disable_smtp: false
disable_http: false
`)
	_, err = confFile.Write(configFileContent)
	if err != nil {
		suite.T().Skip("Unable to write temp config file for test")
	}
	err = confFile.Close()
	if err != nil {
		suite.T().Skip("Unable to close temp config file after write")
	}
	flags := flag.NewFlagSet("TestMailpieRun", flag.PanicOnError)
	arguments := []string{"-config", confFile.Name(), "-imapPort", "2143"}
	go Run(flags, arguments)
	//give the services some time to start
	time.Sleep(2 * time.Second)

	// Testing start

	// SMTP running
	if suite.True(checkPortOpen("127.0.0.1", 1025), "SMTP not running") {
		m := mail.NewMessage()
		m.SetHeader("From", "alex@example.com")
		m.SetHeader("To", "bob@example.com", "cora@example.com")
		m.SetAddressHeader("Cc", "dan@example.com", "Dan")
		m.SetHeader("Subject", "Hello!")
		m.SetBody("text/html", "Hello <b>Bob</b> and <i>Cora</i>!")

		d := mail.NewDialer("127.0.0.1", 1025, "user", "123456")

		// Send the email to Bob, Cora and Dan.
		err = d.DialAndSend(m)
		suite.Nil(err)
	}

	//IMAP running
	if suite.True(checkPortOpen("127.0.0.1", 2143), "IMAP not running") {
		imapConnect(suite.T(), "127.0.0.1:2143")
	}

	//HTTP running
	if suite.True(checkPortOpen("127.0.0.1", 8000), "HTTP not running") {
		response, err := http.Get("http://localhost:8000")
		suite.Nil(err)
		z := html.NewTokenizer(response.Body)
		for {
			tt := z.Next()
			if tt == html.ErrorToken {
				err := z.Err()
				if err == io.EOF {
					// Not an error, we're done and it's valid!
					break
				}
				suite.FailNow("invalid html recieved", err)
			}
		}
	}
}

func TestMailpie(t *testing.T) {
	if !testing.Short() {
		suite.Run(t, new(MailpieTest))
	}
}

func checkPortOpen(host string, port int) bool {
	timeout := time.Second
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, strconv.Itoa(port)), timeout)
	if err != nil {
		return false
	}
	if conn == nil {
		return false
	}
	_ = conn.Close()
	return true
}

func imapConnect(t *testing.T, address string) {
	// Connect to server
	c, err := client.Dial(address)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Logout()
	if err := c.Login("username", "password"); err != nil {
		log.Fatal(err)
	}
	mbox, err := c.Select("INBOX", false)
	if err != nil {
		assert.FailNow(t, "Error during selecting INBOX", err)
	}
	assert.NotNil(t, mbox)

	seqset := new(imap.SeqSet)
	seqset.AddRange(0, 1)

	messagesChan := make(chan *imap.Message, 10)
	done := make(chan error, 1)
	go func() {
		done <- c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messagesChan)
	}()

	messages := []*imap.Message{}
	for msg := range messagesChan {
		messages = append(messages, msg)
	}
	if !assert.Len(t, messages, 1) {
		assert.FailNow(t, "No messages found")
	}
	assert.Equal(t, "Hello!", messages[0].Envelope.Subject)

	if err := <-done; err != nil {
		assert.FailNow(t, "Error during fetching", err)
	}
}
