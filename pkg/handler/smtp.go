package handler

import (
	"fmt"
	"mailpie/pkg/store"
	"net"
)

var SmtpHandler = func(remoteAddr net.Addr, from string, to []string, data []byte) {
	_, err := store.AddMail(data)
	fmt.Println(err)
}
