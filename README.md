# MailPie
**The Mailcatcher for testing and developing**

### Build with
This project was made possible thanks to the amazing work of other people
* [parsemail](https://github.com/DusanKasan/parsemail)
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
Currently, MailPie only supports simple SMTP and IMAP. You are able to add MailPie as an SMTP server für your project
and as an IMAP server for you mail client. Any mail sent via SMTP to MailPie will be visible in your mail client.

#### Planned
- Dockerize MailPie and put it on Dockerhub for simple usage
- Webinterface with Vue 3 and Vuetify communicating over Server-Send-Events with the backend
- REST-API that can be used in test suites for mail testing
- Switch from hardcoded stuff like ports to a way to use configs and CLI flags; Add debug mode
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