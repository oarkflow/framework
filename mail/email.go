package mail

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/oarkflow/frame/pkg/common/bytebufferpool"
	"github.com/oarkflow/frame/pkg/common/utils"
	"github.com/oarkflow/frame/server/render"
	"github.com/oarkflow/log/fqdn"

	"github.com/oarkflow/framework/contracts/mail"
	queue2 "github.com/oarkflow/framework/contracts/queue"
	"github.com/oarkflow/framework/facades"

	"github.com/oarkflow/log"
	sMail "github.com/xhit/go-simple-mail/v2"
)

var maxBigInt = big.NewInt(math.MaxInt64)

type Config struct {
	Host         string `json:"host" yaml:"host" env:"MAIL_HOST"`
	Username     string `json:"username" yaml:"username" env:"MAIL_USERNAME"`
	Password     string `json:"password" yaml:"password" env:"MAIL_PASSWORD"`
	Encryption   string `json:"encryption" yaml:"encryption" env:"MAIL_ENCRYPTION"`
	FromAddress  string `json:"from_address" yaml:"from_address" env:"MAIL_FROM_ADDRESS"`
	AwsAccessKey string `json:"aws_access_key"`
	AwsSecretKey string `json:"aws_secret_key"`
	FromName     string `json:"from_name" yaml:"from_name" env:"MAIL_FROM_NAME"`
	EmailLayout  string `json:"layout" yaml:"layout" env:"MAIL_LAYOUT"`
	Port         int    `json:"port" yaml:"port" env:"MAIL_PORT"`
	Charset      string `json:"charset"`
	Region       string `json:"region"`
}

func GetMailConfig() Config {
	host := facades.Config.GetString("mail.host")
	port := facades.Config.GetInt("mail.port")
	username := facades.Config.GetString("mail.username")
	password := facades.Config.GetString("mail.password")
	encryption := facades.Config.GetString("mail.encryption")
	fromName := facades.Config.GetString("mail.from.name")
	fromAddress := facades.Config.GetString("mail.from.address")
	emailLayout := facades.Config.GetString("mail.layout")
	charset := facades.Config.GetString("mail.charset")
	region := facades.Config.GetString("mail.aws_region")
	awsAccessKey := facades.Config.GetString("mail.aws_access_key")
	awsSecretKey := facades.Config.GetString("mail.aws_secret_key")
	return Config{
		Host:         host,
		Port:         port,
		Username:     username,
		Password:     password,
		Encryption:   encryption,
		FromName:     fromName,
		FromAddress:  fromAddress,
		EmailLayout:  emailLayout,
		Charset:      charset,
		Region:       region,
		AwsAccessKey: awsAccessKey,
		AwsSecretKey: awsSecretKey,
	}
}

type Mailer struct {
	*sMail.SMTPServer
	*sMail.SMTPClient
	*render.HtmlEngine
	Config Config
}

type Attachment struct {
	Data     []byte
	File     string
	FileName string
	MimeType string
}

type Mail struct {
	To          []string     `json:"to,omitempty"`
	From        string       `json:"from,omitempty"`
	Subject     string       `json:"subject,omitempty"`
	Body        string       `json:"body,omitempty"`
	Bcc         []string     `json:"bcc,omitempty"`
	Cc          []string     `json:"cc,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
	AttachFiles []Attachment `json:"attach_files"`
	engine      *render.HtmlEngine
}

var DefaultMailer *Mailer

func Default(cfg Config, templateEngine *render.HtmlEngine) {
	DefaultMailer = New(cfg, templateEngine)
}

func New(cfg Config, templateEngine *render.HtmlEngine) *Mailer {
	m := &Mailer{Config: cfg}
	m.HtmlEngine = templateEngine
	if cfg.AwsAccessKey != "" && cfg.AwsSecretKey != "" {
		if cfg.Region == "" {
			cfg.Region = "us-east-1"
		}
		if cfg.Charset == "" {
			cfg.Charset = "UTF-8"
		}
		return m
	}
	m.SMTPServer = sMail.NewSMTPClient()
	m.SMTPServer.Host = cfg.Host
	m.SMTPServer.Port = cfg.Port
	m.SMTPServer.Username = cfg.Username
	m.SMTPServer.Password = cfg.Password
	if cfg.Encryption == "tls" {
		m.SMTPServer.Encryption = sMail.EncryptionSTARTTLS
	} else {
		m.SMTPServer.Encryption = sMail.EncryptionSSL
	}
	// Variable to keep alive connection
	m.SMTPServer.KeepAlive = false
	// Timeout for connect to SMTP Server
	m.SMTPServer.ConnectTimeout = 10 * time.Second
	// Timeout for send the data and wait respond
	m.SMTPServer.SendTimeout = 10 * time.Second
	return m
}

func (m *Mailer) sendSes(msg Mail) error {
	// Create a new session in the us-west-2 region.
	// Replace us-west-2 with the AWS Region you're using for Amazon SES.
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(m.Config.Region),
		Credentials: credentials.NewStaticCredentials(m.Config.AwsAccessKey, m.Config.AwsSecretKey, ""),
	})
	if err != nil {
		return err
	}
	// Create an SES session.
	svc := ses.New(sess)
	var recipient, ccRecipient, bccRecipient []*string
	for _, to := range msg.To {
		recipient = append(recipient, aws.String(to))
	}
	for _, bcc := range msg.Bcc {
		bccRecipient = append(bccRecipient, aws.String(bcc))
	}
	for _, cc := range msg.Cc {
		ccRecipient = append(ccRecipient, aws.String(cc))
	}
	// Assemble the email.
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses:  ccRecipient,
			BccAddresses: bccRecipient,
			ToAddresses:  recipient,
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(m.Config.Charset),
					Data:    aws.String(msg.Body),
				},
				Text: &ses.Content{
					Charset: aws.String(m.Config.Charset),
					Data:    aws.String(""),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(m.Config.Charset),
				Data:    aws.String(msg.Subject),
			},
		},
		Source: aws.String(msg.From),
		// Uncomment to use a configuration set
		// ConfigurationSetName: aws.String(ConfigurationSet),
	}

	// Attempt to send the email.
	_, err = svc.SendEmail(input)
	return err
}

func (m *Mailer) Send(msg Mail) error {
	if m.SMTPServer == nil && m.Config.AwsAccessKey != "" && m.Config.AwsSecretKey != "" {
		return m.sendSes(msg)
	}
	var err error
	m.SMTPClient, err = m.SMTPServer.Connect()
	if err != nil {
		fmt.Println("Error on connection: " + err.Error())
		return err
	}
	defer m.SMTPClient.Close()
	// New email simple html with inline and CC
	email := sMail.NewMSG()
	if msg.From == "" {
		msg.From = fmt.Sprintf("%s<%s>", m.Config.FromName, m.Config.FromAddress)
	}
	email.SetFrom(msg.From).AddTo(msg.To...).SetSubject(msg.Subject)
	if len(msg.Cc) > 0 { //nolint:wsl
		email.AddCc(msg.Cc...)
	}
	if len(msg.Bcc) > 0 { //nolint:wsl
		email.AddBcc(msg.Bcc...)
	}
	// txt, _ := html2text.FromString(body, html2text.Options{PrettyTables: false})
	// email.AddAlternative(sMail.TextPlain, txt)
	email.SetBody(sMail.TextHTML, msg.Body) //nolint:wsl
	for _, attachment := range msg.Attachments {
		email.AddAttachmentData(attachment.Data, attachment.File, attachment.MimeType)
	}
	for _, attachment := range msg.AttachFiles {
		email.AddAttachment(attachment.File, attachment.FileName)
	}

	// Call Send and pass the client
	err = email.Send(m.SMTPClient)
	if err != nil {
		return err
	} else {
		log.Info().Msg("Email Sent to " + strings.Join(msg.To, ", "))
	}
	return nil
}

func (m *Mailer) Queue(msg Mail, queue *mail.Queue) error {
	job := facades.Queue.Job(&SendMailJob{}, []queue2.Arg{
		{Value: msg.Subject, Type: "string"},
		{Value: msg.Body, Type: "string"},
		{Value: msg.From, Type: "string"},
		{Value: msg.To, Type: "[]string"},
		{Value: msg.Cc, Type: "[]string"},
		{Value: msg.Bcc, Type: "[]string"},
		{Value: msg.Attachments, Type: "[]Attachment"},
		{Value: msg.AttachFiles, Type: "[]Attachment"},
	})
	if queue != nil {
		if queue.Connection != "" {
			job.OnConnection(queue.Connection)
		}
		if queue.Queue != "" {
			job.OnQueue(queue.Queue)
		}
	}

	return job.Dispatch()
}

func (m *Mailer) View(view string, body utils.H) *Body {
	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)
	if err := m.Render(buf, view, body, m.Config.EmailLayout); err != nil {
		panic(err)
	}
	bodyContent := &Body{Content: buf.String(), mailer: m}
	return bodyContent
}

func (m *Mailer) Html(view string, body utils.H) string {
	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)
	if err := m.Render(buf, view, body, DefaultMailer.Config.EmailLayout); err != nil {
		panic(err)
	}
	return buf.String()
}

func View(view string, body utils.H) *Body {
	bodyContent := &Body{Content: Html(view, body)}
	return bodyContent
}

func Html(view string, body utils.H) string {
	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)
	if err := DefaultMailer.Render(buf, view, body, DefaultMailer.Config.EmailLayout); err != nil {
		panic(err)
	}
	return buf.String()
}

func Send(msg Mail) error {
	return DefaultMailer.Send(msg)
}

func Queue(msg Mail, queue *mail.Queue) error {
	return DefaultMailer.Queue(msg, queue)
}

type Body struct {
	Content string
	mailer  *Mailer
}

func SendMail(msg Mail) error {
	mailer := New(GetMailConfig(), nil)
	return mailer.Send(msg)
}

func (t *Body) Send(msg Mail) error {
	msg.Body = t.Content
	if t.mailer != nil {
		return t.mailer.Send(msg)
	}
	if DefaultMailer == nil {
		return errors.New("No mailer configured")
	}
	return DefaultMailer.Send(msg)
}

func (t *Body) Queue(msg Mail, queue *mail.Queue) error {
	msg.Body = t.Content
	if t.mailer != nil {
		return t.mailer.Queue(msg, queue)
	}
	if DefaultMailer == nil {
		return errors.New("No mailer configured")
	}
	return DefaultMailer.Queue(msg, queue)
}

func generateMessageID() (string, error) {
	t := time.Now().UnixNano()
	pid := os.Getpid()
	rint, err := rand.Int(rand.Reader, maxBigInt)
	if err != nil {
		return "", err
	}
	h, err := fqdn.Hostname()
	// If we can't get the hostname, we'll use localhost
	if err != nil {
		h = "localhost.localdomain"
	}
	msgid := fmt.Sprintf("<%d.%d.%d@%s>", t, pid, rint, h)
	return msgid, nil
}
