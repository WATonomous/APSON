package plantops

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Announcement represents a parsed PlantOps service interruption announcement.
type Announcement struct {
	Title string
	Link  string
}

// FetchAndParse fetches the PlantOps Service Interruptions page and parses relevant announcements.
// buildings is a list of building codes/names to match (e.g., ["CPH", "E2"]).
func FetchAndParse(buildings []string) ([]Announcement, error) {
	url := "https://plantops.uwaterloo.ca/service-interruptions/"
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch PlantOps page: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	var results []Announcement
	doc.Find("a.w3-leftbar").Each(func(i int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Text())
		link, _ := s.Attr("href")
		for _, b := range buildings {
			if strings.Contains(strings.ToUpper(title), strings.ToUpper(b)) {
				results = append(results, Announcement{
					Title: title,
					Link:  link,
				})
				break
			}
		}
	})

	return results, nil
}
