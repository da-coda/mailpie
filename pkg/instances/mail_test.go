package instances

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	gomail "net/mail"
	"testing"
)

var mail = []byte(`Received: from localhost (localhost [127.0.0.1])
        by localhost (Mailpie) with SMTP
        for <bob@example.com>; Wed, 27 Jan 2021 17:00:48 +0100 (CET)
MIME-Version: 1.0
Date: Wed, 27 Jan 2021 17:00:48 +0100
From: alex@example.com
To: bob@example.com, cora@example.com
Cc: "Dan" <dan@example.com>
Subject: Hello!
Content-Type: text/html; charset=UTF-8
Content-Transfer-Encoding: quoted-printable

Hello <b>Bob</b> and <i>Cora</i>!
`)

type MailUnitTestSuite struct {
	suite.Suite
}

func (suite *MailUnitTestSuite) TestParseMail() {
	parsed, err := ParseMail(mail)
	assert.Nil(suite.T(), err, "No error expected")
	parseRawDirectly, err := gomail.ReadMessage(bytes.NewReader(parsed.RawMessage))
	assert.Nil(suite.T(), err, "No error expected")
	bodyFromParseMail, err := ioutil.ReadAll(parsed.Body)
	assert.Nil(suite.T(), err, "No error expected")
	bodyFromTestParse, err := ioutil.ReadAll(parseRawDirectly.Body)
	assert.Nil(suite.T(), err, "No error expected")
	assert.Equal(suite.T(), bodyFromParseMail, bodyFromTestParse, "Body from ParseMail and direct parsing are not equal")
}

func (suite *MailUnitTestSuite) TestRead_ValidMail() {
	parsed, err := ParseMail(mail)
	assert.Nil(suite.T(), err, "No error expected")
	readParsed, err := ioutil.ReadAll(parsed)
	assert.Nil(suite.T(), err, "No error expected")
	assert.Equal(suite.T(), mail, readParsed)
}

func (suite *MailUnitTestSuite) TestRead_InvalidMail() {
	invalidMail := []byte("I am not a Mail")
	parsed, err := ParseMail(invalidMail)
	assert.Nil(suite.T(), parsed, "parsed should be nil on error")
	assert.Error(suite.T(), err, "missing error with invalid mail")
}

func (suite *MailUnitTestSuite) TestLen() {
	parsed, err := ParseMail(mail)
	assert.Nil(suite.T(), err, "No error expected")
	assert.Equal(suite.T(), len(mail), parsed.Len(), "Length of parsed mail should be the Same as length of original mail")
}

func TestMailUnitTestSuite(t *testing.T) {
	suite.Run(t, new(MailUnitTestSuite))
}
