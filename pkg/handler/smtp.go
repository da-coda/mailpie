package handler

import (
	"fmt"
	"mailmon-go/pkg/store"
	"net"
)

var SmtpHandler = func(remoteAddr net.Addr, from string, to []string, data []byte) {
	err, _ := store.AddMail(data)
	fmt.Println(err)
}
