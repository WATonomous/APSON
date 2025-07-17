package plantops

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Announcement represents a parsed PlantOps service interruption announcement.
type Announcement struct {
	Title string
	Link  string
}

// normalizeBuilding tries to map common building names to codes (e.g., 'Engineering 2' -> 'E2').
func normalizeBuilding(s string) string {
	s = strings.ToUpper(strings.TrimSpace(s))
	replacements := map[string]string{
		"ENGINEERING 2":              "E2",
		"ENGINEERING 3":              "E3",
		"CARL A POLLOCK HALL":        "CPH",
		"CARL A.POLLOCK HALL":        "CPH",
		"CARL POLLOCK HALL":          "CPH",
		"DOUGLAS WRIGHT ENGINEERING": "DWE",
		"PHYSICS":                    "PHY",
		"SOUTH CAMPUS HALL":          "SCH",
		// ...
	}
	if v, ok := replacements[s]; ok {
		return v
	}
	return s
}

// containsBuilding checks if any building code or normalized name is present in the text.
func containsBuilding(text string, buildings []string) bool {
	text = strings.ToUpper(text)

	// Regex to extract 'ENGINEERING N' and 'ENGINEERING N & M' patterns
	engPattern := regexp.MustCompile(`ENGINEERING [0-9](?: ?[&AND]+ ?[0-9]+)*`)
	matches := engPattern.FindAllString(text, -1)
	for _, match := range matches {
		// e.g., match = 'ENGINEERING 2 & 3' or 'ENGINEERING 2 AND 3'
		nums := regexp.MustCompile(`[0-9]+`).FindAllString(match, -1)
		for _, n := range nums {
			bldg := normalizeBuilding("ENGINEERING " + n)
			for _, b := range buildings {
				if bldg == strings.ToUpper(b) {
					return true
				}
			}
		}
	}

	// Fallback: Split by comma, ' and ', ' & '
	parts := splitBuildings(text)
	for _, part := range parts {
		p := strings.TrimSpace(part)
		pNorm := normalizeBuilding(p)
		for _, b := range buildings {
			if p == strings.ToUpper(b) || pNorm == strings.ToUpper(b) {
				return true
			}
			if strings.Contains(pNorm, strings.ToUpper(b)) || strings.Contains(p, strings.ToUpper(b)) {
				return true
			}
		}
	}
	return false
}

// isNumber returns true if the string is a number.
func isNumber(s string) bool {
	if len(s) == 0 {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// splitBuildings splits a string by ',', ' and ', and ' & ' delimiters.
func splitBuildings(text string) []string {
	// Replace ' and ' and ' & ' with ','
	replacer := strings.NewReplacer(" and ", ",", " & ", ",")
	t := replacer.Replace(text)
	return strings.Split(t, ",")
}

// containsKeyword checks if any relevant keyword is present in the text.
func containsKeyword(text string) bool {
	text = strings.ToUpper(text)
	return strings.Contains(text, "ELECTRICAL") || strings.Contains(text, "POWER")
}

// FetchAndParse fetches the PlantOps Service Interruptions page and parses relevant announcements.
// buildings is a list of building codes/names to match (e.g., ["CPH", "E2"]).
// It fetches each notice.php page to check for relevance.
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
		link, _ := s.Attr("href")
		fullLink := link
		if !strings.HasPrefix(link, "http") {
			fullLink = "https://plantops.uwaterloo.ca/service-interruptions/" + link
		}

		// Fetch the full notice page
		nResp, err := http.Get(fullLink)
		if err != nil || nResp.StatusCode != 200 {
			if nResp != nil {
				nResp.Body.Close()
			}
			return // skip this notice if fetch fails
		}
		defer nResp.Body.Close()
		bodyBytes, err := io.ReadAll(nResp.Body)
		if err != nil {
			return
		}
		bodyStr := string(bodyBytes)
		nDoc, err := goquery.NewDocumentFromReader(strings.NewReader(bodyStr))
		if err != nil {
			return
		}
		nTitle := ""
		nDoc.Find("h3").EachWithBreak(func(i int, s *goquery.Selection) bool {
			nTitle = strings.TrimSpace(s.Text())
			return false // only first h3
		})
		relevant := IsRelevantAnnouncement(nTitle, bodyStr, buildings)
		if relevant {
			results = append(results, Announcement{Title: nTitle, Link: link})
		}
	})

	return results, nil
}

// IsRelevantAnnouncement checks if an announcement is relevant by looking for buildings and keywords in the title and description.
func IsRelevantAnnouncement(title, fullText string, buildings []string) bool {
	if containsKeyword(title) && containsBuilding(title, buildings) {
		return true
	}
	if fullText != "" {
		if containsKeyword(fullText) && containsBuilding(fullText, buildings) {
			return true
		}
		// Try to extract 'Where is this happening?' and 'What is happening?' sections
		sections := extractSections(fullText)
		for k, v := range sections {
			if (k == "where" || k == "what") && containsKeyword(v) && containsBuilding(v, buildings) {
				return true
			}
			if (k == "where" || k == "what") && (containsKeyword(title) || containsKeyword(v)) && containsBuilding(v, buildings) {
				return true
			}
		}
	}
	return false
}

// extractSections tries to extract key sections from the HTML/text for more robust matching.
func extractSections(html string) map[string]string {
	sections := make(map[string]string)
	// Very basic regex-based extraction for demo purposes
	whereRe := regexp.MustCompile(`(?i)<div class="section-header">Where is this happening\?</div>\s*<p class="section-text">([^<]+)</p>`)
	whatRe := regexp.MustCompile(`(?i)<div class="section-header">What is happening\?</div>\s*<p class="section-text">([^<]+)</p>`)
	if m := whereRe.FindStringSubmatch(html); len(m) > 1 {
		sections["where"] = m[1]
	}
	if m := whatRe.FindStringSubmatch(html); len(m) > 1 {
		sections["what"] = m[1]
	}
	return sections
}
