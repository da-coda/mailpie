package main

import (
	"log"
	"strings"

	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
)

func main() {
	// Set up authentication information.
	auth := sasl.NewPlainClient("", "", "")

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	to := []string{"recipient@example.net"}
	msg := strings.NewReader(`
		I
	Gonna
		Test
This
Shit!
`)
	err := smtp.SendMail("localhost:1025", auth, "sender@example.org", to, msg)
	if err != nil {
		log.Fatal(err)
	}
}