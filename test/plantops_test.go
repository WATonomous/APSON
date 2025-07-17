package test

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/WATonomous/APSON/internal/plantops"
)

func TestFetchAndParseNoticesOnline(t *testing.T) {
	cases := []struct {
		Link           string
		ExpectRelevant bool
	}{
		{"https://plantops.uwaterloo.ca/service-interruptions/notice.php?ID=1454", true},
		{"https://plantops.uwaterloo.ca/service-interruptions/notice.php?ID=1855", true},
		{"https://plantops.uwaterloo.ca/service-interruptions/notice.php?ID=2478", true},
		{"https://plantops.uwaterloo.ca/service-interruptions/notice.php?ID=2655", true},
		{"https://plantops.uwaterloo.ca/service-interruptions/notice.php?ID=2818", true},
	}

	buildings := []string{"CPH", "E2"}
	passCount := 0
	failCount := 0

	for _, c := range cases {
		resp, err := http.Get(c.Link)
		if err != nil {
			t.Errorf("Failed to fetch %s: %v", c.Link, err)
			failCount++
			continue
		}
		if resp.StatusCode != 200 {
			t.Errorf("Non-200 status for %s: %d", c.Link, resp.StatusCode)
			resp.Body.Close()
			failCount++
			continue
		}
		bodyBytes, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			t.Errorf("Failed to read body for %s: %v", c.Link, err)
			failCount++
			continue
		}
		bodyStr := string(bodyBytes)
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(bodyStr))
		if err != nil {
			t.Errorf("Failed to parse %s: %v", c.Link, err)
			failCount++
			continue
		}

		title := ""
		doc.Find("h3").EachWithBreak(func(i int, s *goquery.Selection) bool {
			title = strings.TrimSpace(s.Text())
			return false // only first h3
		})

		relevant := plantops.IsRelevantAnnouncement(title, bodyStr, buildings)
		if c.ExpectRelevant && !relevant {
			t.Errorf("%s: expected relevant announcement, got none", c.Link)
			failCount++
		} else if c.ExpectRelevant && relevant {
			t.Logf("PASS: %s: found relevant announcement: %q", c.Link, title)
			passCount++
		} else if !c.ExpectRelevant && !relevant {
			t.Logf("PASS: %s: correctly found 0 relevant announcements", c.Link)
			passCount++
		} else if !c.ExpectRelevant && relevant {
			t.Errorf("%s: expected no relevant announcement, but found one", c.Link)
			failCount++
		}
	}
	t.Logf("Test summary: %d passed, %d failed", passCount, failCount)
}
