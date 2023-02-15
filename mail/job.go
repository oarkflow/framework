package mail

import "fmt"

type SendMailJob struct {
}

// Signature The name and signature of the job.
func (r *SendMailJob) Signature() string {
	return "goravel_send_mail_job"
}

// Handle Execute the job.
func (r *SendMailJob) Handle(args ...any) error {
	msg := Mail{
		To:          args[4].([]string),
		From:        fmt.Sprintf("%s<%s>", args[3].(string), args[2].(string)),
		Subject:     args[0].(string),
		Body:        args[1].(string),
		Bcc:         args[6].([]string),
		Cc:          args[5].([]string),
		AttachFiles: args[7].([]Attachment),
	}
	return SendMail(msg)
}
