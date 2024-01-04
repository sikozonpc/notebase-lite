package medium

import (
	"bytes"
	"testing"

	"github.com/sikozonpc/notebase/types"
)

func TestBuildInsightsMailTemplate(t *testing.T) {
	{

		t.Run("BuildInsightsMailTemplate should return a list with insights", func(t *testing.T) {
			insights := []*types.DailyInsight{
				{
					Text:        "This is an insight",
					Note:        "This is a note",
					BookAuthors: "John Doe",
					BookTitle:   "Gopher",
				},
			}

			u := &types.User{
				FirstName: "Test",
				LastName:  "Test",
				Email:     "gopher@gopher.xyz",
				ID:        42,
				IsActive:  true,
			}

			html := BuildInsightsMailTemplate("../template", u, insights, "some-random-token")

			if html == "" {
				t.Errorf("BuildInsightsMailTemplate() = %v; want %v", html, "html")
			}

			if !bytes.Contains([]byte(html), []byte("This is an insight")) {
				t.Errorf("BuildInsightsMailTemplate() = %v; want %v", html, "html")
			}

			if !bytes.Contains([]byte(html), []byte("This is a note")) {
				t.Errorf("BuildInsightsMailTemplate() = %v; want %v", html, "html")
			}

			if !bytes.Contains([]byte(html), []byte("John Doe")) {
				t.Errorf("BuildInsightsMailTemplate() = %v; want %v", html, "html")
			}

			if !bytes.Contains([]byte(html), []byte("Gopher")) {
				t.Errorf("BuildInsightsMailTemplate() = %v; want %v", html, "html")
			}
		})
	}
}
