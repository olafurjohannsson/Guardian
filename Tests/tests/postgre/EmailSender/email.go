
package EmailSender

import (
    "net/smtp"
    "strconv"
)

type EmailSender interface {
    Send() error
    Create() EmailMessage
    Init() * EmailAuth
}

type EmailAuth struct {
    User string
    Password string
    Hostname string
    Identity string
    Port int
}

type EmailMessage struct {
    To []string
    From string
    Body []byte
}


// Create mail client
func Init(user string, password string, hostname string, port int, identity string) *EmailAuth {
    return &EmailAuth{
        User: user,
        Password: password,
        Hostname: hostname,
        Port: port,
        Identity: identity,
    }
}

// Create email
func Create(to []string, from string, body []byte) EmailMessage {
    return EmailMessage{
        To: to,
        From: from,
        Body: body,
    }
}

// Send email
func Send(email EmailMessage, emailAuth EmailAuth) error {
    auth := smtp.PlainAuth(emailAuth.Identity, emailAuth.User, emailAuth.Password, emailAuth.Hostname)
    
    err := smtp.SendMail(emailAuth.Hostname + ":" + strconv.Itoa(emailAuth.Port), auth, email.From, email.To, email.Body)
    
    if err != nil {
        return err
    }
    
    return nil
}
