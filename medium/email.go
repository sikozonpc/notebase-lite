package medium

import (
	"bytes"
	"fmt"
	"log"
	"path/filepath"
	"text/template"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/sikozonpc/notebase/config"
	t "github.com/sikozonpc/notebase/types"
)

var FromName = "Notebase"

type Mailer struct {
	FromEmail string
	Client    *sendgrid.Client
}

func NewMailer(apiKey, fromEmail string) *Mailer {
	client := sendgrid.NewSendClient(apiKey)

	return &Mailer{
		FromEmail: fromEmail,
		Client:    client,
	}
}

func (m *Mailer) SendInsights(u *t.User, insights []*t.DailyInsight, authToken string) error {
	from := mail.NewEmail(FromName, m.FromEmail)
	subject := "Daily Insight(s)"
	userName := fmt.Sprintf("%v %v", u.FirstName, u.LastName)

	if u.Email == "" {
		return fmt.Errorf("user has no email")
	}

	to := mail.NewEmail(userName, u.Email)

	htmlContent := BuildInsightsMailTemplate("template", u, insights, authToken)

	message := mail.NewSingleEmail(from, subject, to, "", htmlContent)
	response, err := m.Client.Send(message)
	if err != nil {
		log.Println(err)
	}

	log.Printf("Email sent to %v with status code %v", u.Email, response.StatusCode)

	return nil
}

func BuildInsightsMailTemplate(templateDir string, u *t.User, insights []*t.DailyInsight, authToken string) string {
	filename := filepath.Join(templateDir, "daily.tmpl")
	tmpl, err := template.ParseFiles(filename)
	if err != nil {
		panic(err)
	}

	payload := struct {
		UnsubscribeURL string
		User           *t.User
		Insights       []*t.DailyInsight
	}{
		UnsubscribeURL: fmt.Sprintf("%s/unsubscribe.html?token=%s", config.Envs.PublicURL, authToken),
		User:           u,
		Insights:       insights,
	}

	var out bytes.Buffer
	err = tmpl.Execute(&out, payload)
	if err != nil {
		panic(err)
	}

	return out.String()
}
