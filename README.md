<p align="center">
  <img width="100px" src="https://github.com/da-coda/mailpie/blob/main/readme/magpie.png?raw=true" alt="Mailpie logo"/>
</p>
<h1 align="center">MailPie</h1>

**The Mailcatcher for testing and developing**

[![Maintenance](https://img.shields.io/badge/Maintained%3F-yes-green.svg)](https://github.com/da-coda/mailpie/graphs/commit-activity)
[![made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](http://golang.org)
[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/da-coda/mailpie.svg)](https://github.com/da-coda/mailpie)
[![Go Report Card](https://goreportcard.com/badge/github.com/da-coda/mailpie)](https://goreportcard.com/report/github.com/da-coda/mailpie)
[![GitHub license](https://img.shields.io/github/license/da-coda/mailpie.svg)](https://github.com/da-coda/mailpie/blob/master/LICENSE)
[![GitHub release](https://img.shields.io/github/release/da-coda/mailpie.svg)](https://GitHub.com/da-coda/mailpie/releases/)
[![GitHub stars](https://img.shields.io/github/stars/da-coda/mailpie.svg?style=social&label=Star&maxAge=2592000)](https://GitHub.com/da-coda/mailpie/stargazers/)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat-square)](http://makeapullrequest.com)

![Current](https://github.com/da-coda/mailpie/workflows/Current/badge.svg?branch=main)
![Release](https://github.com/da-coda/mailpie/workflows/Release/badge.svg?branch=release)
![Develop](https://github.com/da-coda/mailpie/workflows/Develop/badge.svg?branch=develop)
### Build with
This project was made possible thanks to the amazing work of other people
* [go-imap](https://github.com/emersion/go-imap)
* [mux](https://github.com/gorilla/mux)
* [smtpd](https://github.com/mhale/smtpd)
* [sse](https://github.com/r3labs/sse)

## Why MailPie?
MailPie aims to satisfy your needs in development and testing environments regarding mails.
With multiple ways to view your mails you are able to test and debug in dev and test environments
without any actual SMTP servers.

### Features
#### Implemented
Currently, MailPie only supports simple SMTP and IMAP. You are able to add MailPie as an SMTP server for your project
and as an IMAP server for you mail client. Any mail sent via SMTP to MailPie will be visible in your mail client.

#### Planned
- Webinterface with Vue 3 communicating over Server-Send-Events and REST Api with the backend
- REST-API that can be used in test suites for mail testing
- Codeception(PHP) Module for testing with the REST-API
- Advanced SMTP and IMAP handling
- Maybe supporting usage as a proxy mail server for mail logging?
- Implement [spamassassin](https://github.com/Teamwork/spamc) support

## How to use MailPie?
### SMTP - settings in your project
1. Execute the MailPie binary or build the binary yourself
2. Add MailPie as you SMTP Server to your project
    - *Port*: 1025
    - *Username*: any
    - *Password*: any

Due to the fact that MailPie only acts as a mail catcher in Dev and Test environments there
is no credentials check needed.

### IMAP & SMTP - settings in your mail client
Add MailPie as a new account to your mail client
- Email address, Username and Password can be anything

![Mail Settings in Thunderbird 68][mail-settings]

Currently, you will need to trigger the mail receiving manually in your mail client or set the auto update to a short 
duration.

[mail-settings]: readme/mail_settings.png

## Contact
Daniel Müller - Twitter: [@da_coda_](https://twitter.com/da_coda_) - LinkedIn: [daniel96mueller](https://www.linkedin.com/in/daniel96mueller/) - E-Mail: [contact@daniel-mueller.de](mailto:contact@daniel-mueller.de)