package medium

import (
	"bytes"
	"fmt"
	"log"
	"text/template"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
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

func (m *Mailer) SendInsights(u *t.User, insights []*t.DailyInsight) error {
	from := mail.NewEmail(FromName, m.FromEmail)
	subject := "Daily Insight(s)"
	userName := fmt.Sprintf("%v %v", u.FirstName, u.LastName)

	if u.Email == "" {
		return fmt.Errorf("user has no email")
	}

	to := mail.NewEmail(userName, u.Email)

	htmlContent := buildInsightsMailTemplate(u, insights)

	message := mail.NewSingleEmail(from, subject, to, "", htmlContent)
	response, err := m.Client.Send(message)
	if err != nil {
		log.Println(err)
	}

	log.Printf("Email sent to %v with status code %v", u.Email, response.StatusCode)

	return nil
}

func buildInsightsMailTemplate(u *t.User, insights []*t.DailyInsight) string {
	tmpl, err := template.ParseFiles("daily.tmpl")
	if err != nil {
		panic(err)
	}

	payload := struct {
		User     *t.User
		Insights []*t.DailyInsight
	}{
		User:     u,
		Insights: insights,
	}

	var out bytes.Buffer
	err = tmpl.Execute(&out, payload)
	if err != nil {
		panic(err)
	}

	return out.String()
}

