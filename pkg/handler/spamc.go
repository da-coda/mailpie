package handler

import (
	"bytes"
	"context"
	"github.com/da-coda/mailpie/pkg/event"
	"github.com/da-coda/mailpie/pkg/instances"
	"github.com/da-coda/mailpie/pkg/store"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/teamwork/spamc"
)

type Spamc struct {
	client spamc.Client
}

func NewSpamc(client spamc.Client, subscribable event.Subscribable) *Spamc {
	s := Spamc{client: client}
	subscribable.Subscribe(store.NewMailStoredEvent, s.Handler)
	return &s
}

func (s Spamc) Handler(dispatcher string, data interface{}) {
	mail, isMail := data.(instances.Mail)
	if !isMail {
		logrus.Error("received non mail data")
		return
	}
	ctx := context.Background()
	err := s.client.Ping(ctx)
	if err != nil {
		logrus.WithError(err).Error("spamd is not pingable")
		for unwrapped := errors.Unwrap(err); unwrapped != nil; {
			logrus.WithError(err).Error("spamd is not pingable")
		}
		return
	}
	mailReader := bytes.NewReader(mail.RawMessage)
	check, err := s.client.Check(ctx, mailReader, nil)
	if err != nil {
		logrus.WithError(err).Error("Error on score")
	}
	logrus.WithField("score", check.Score).WithField("basescore", check.BaseScore).WithField("isSpam", check.IsSpam).Debug("Spamscore received")
	mailReader = bytes.NewReader(mail.RawMessage)
	report, err := s.client.Report(ctx, mailReader, nil)
	if err != nil {
		logrus.WithError(err).Error("Error on score")
	}
	logrus.Debug(report.Report.String())
}
