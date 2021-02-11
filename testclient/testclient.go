package main

import (
	"gopkg.in/mail.v2"
)

func main() {
	m := mail.NewMessage()
	m.SetHeader("From", "alex@example.com")
	m.SetHeader("To", "bob@example.com", "cora@example.com")
	m.SetAddressHeader("Cc", "dan@example.com", "Dan")
	m.SetHeader("Subject", "Hello!")
	m.SetBody("text/html", "\n<!DOCTYPE HTML PUBLIC \"-//W3C//DTD HTML 4.0 Transitional//EN\">\n<HTML><HEAD>\n<META content=\"text/html; charset=iso-8859-1\" http-equiv=Content-Type>\n<META content=\"MSHTML 5.00.2920.0\" name=GENERATOR>\n<STYLE></STYLE>\n</HEAD>\n<BODY bgColor=#ffffff>\n<DIV><B><U><FONT size=5>\n<P align=center><FONT color=#ff0000>Internet Service \nProviders, 1,000,000€</FONT></P></B></U></FONT>\n<P>&nbsp;</P>\n<P align=center><FONT size=4>We apologize if this is an unwanted email. We \nassure you this is a one time mailing only.</FONT></P><FONT size=4>\n<P>We represent a marketing corporation interested in buying an ISP or \npartnership with an ISP, 1,000,000$. We want to provide services for bulk friendly hosting \nof non-illicit websites internationally. We seek your support so that we may \nprovide dependable and efficient hosting to the growing clientele of this \never-expanding industry. Consider this proposition seriously. We believe this \nwould be a lucrative endeavor for you. Please contact <A \nhref=\"mailto:dockut3@hotmail.com\"><U><FONT \ncolor=#0000ff>dockut3@hotmail.com</U></FONT> </A>soon for further discussion, \nquestions, and problem solving. Sincerely.</FONT></P></DIV></BODY></HTML>\n\nhttp://xent.com/mailman/listinfo/fork")

	d := mail.NewDialer("localhost", 1025, "user", "123456")

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}
