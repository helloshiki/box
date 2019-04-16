package commute

import (
	"github.com/go-gomail/gomail"
)

//type unencryptedAuth struct {
//	smtp.Auth
//}
//
//func (a unencryptedAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
//	s := *server
//	s.TLS = true
//	return a.Auth.Start(&s)
//}

type smtpInfo struct {
	host     string
	port     int
	username string
	password string
}

var smtpConf *smtpInfo

type EmailOpt struct {
	FromName  string
	ToAddress string
	ToName    string
	Subject   string
	HtmlBody  string
}

func SendRegisterEmail(opt *EmailOpt) error {
	if smtpConf == nil {
		panic("init smtp conf first!")
	}

	m := gomail.NewMessage()
	m.SetAddressHeader("From", smtpConf.username, opt.FromName)
	m.SetHeader("To", m.FormatAddress(opt.ToAddress, opt.ToName))
	m.SetHeader("Subject", opt.Subject)
	m.SetBody("text/html", opt.HtmlBody)

	//d := gomail.NewDialer("smtp.qq.com", 587, "329365307@qq.com", "vmqlkjdpxfpucafj")
	// d.TLSConfig = &tls.Config{ InsecureSkipVerify: true, }
	d := gomail.NewDialer(smtpConf.host, smtpConf.port, smtpConf.username, smtpConf.password)
	return d.DialAndSend(m)
}

func SetupEmailConfig(host string, port int, username, password string) {
	smtpConf = &smtpInfo{
		host:     host,
		port:     port,
		username: username,
		password: password,
	}
}
