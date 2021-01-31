package instances

import (
	"bytes"
	"io"
	gomail "net/mail"
)

type Mail struct {
	gomail.Message
	RawMessage []byte
	readIndex  int64
	Flags      []string
}

func (m *Mail) Read(p []byte) (n int, err error) {
	if m.readIndex >= int64(len(m.RawMessage)) {
		err = io.EOF
		m.readIndex = 0
		return
	}

	n = copy(p, m.RawMessage[m.readIndex:])
	m.readIndex += int64(n)
	return
}

func (m *Mail) Len() int {
	return len(m.RawMessage)
}

func (m Mail) Get32BitTimestamp() uint32 {
	time, _ := m.Header.Date()
	return uint32(time.Unix())
}

func ParseMail(data []byte) (*Mail, error) {
	parsedMail, err := gomail.ReadMessage(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	return &Mail{Message: *parsedMail, RawMessage: data}, nil
}
