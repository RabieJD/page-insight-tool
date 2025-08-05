package handlers

import (
	"context"
	"fmt"
	"github.com/rabie/page-insight-tool/app/helper"
	"html/template"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	DefaultTimeout  = 30 * time.Second
	UserAgent       = "Page-Insight-Tool/1.0"
	templatePath    = "app/templates/index.html"
	maxLinksToCheck = 500
	maxWorkers      = 10
)

// PageAnalysis holds the result of analyzing a web page
type PageAnalysis struct {
	URL               string
	Title             string
	HTMLVersion       string
	HeadingsCount     map[string]int
	InternalLinks     int
	ExternalLinks     int
	InaccessibleLinks int
	HasLoginForm      bool
	Error             string
}

// IndexHandler renders the form
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	_ = tmpl.Execute(w, nil)
}

// AnalyzeHandler handles the form submission
func AnalyzeHandler(w http.ResponseWriter, r *http.Request) {
	urlStr := r.FormValue("url")
	if urlStr == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	result := analyzePage(urlStr)

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	_ = tmpl.Execute(w, result)
}

// analyzePage orchestrates the full analysis
func analyzePage(urlStr string) PageAnalysis {
	result := PageAnalysis{
		URL:           urlStr,
		HeadingsCount: make(map[string]int),
	}

	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		result.Error = "Invalid URL"
		return result
	}

	if err := validateURL(parsedURL); err != nil {
		result.Error = err.Error()
		return result
	}

	doc, err := fetchPage(parsedURL.String())
	if err != nil {
		result.Error = fmt.Sprintf("Failed to fetch page: %v", err)
		return result
	}

	result.Title = extractTitle(doc)
	result.HTMLVersion = detectHTMLVersion(doc)
	result.HeadingsCount = countHeadings(doc)
	result.InternalLinks, result.ExternalLinks, result.InaccessibleLinks = countLinks(doc, parsedURL)
	result.HasLoginForm = detectLoginForm(doc)

	return result
}

// validateURL ensures the URL is safe to access
func validateURL(u *url.URL) error {
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("unsupported URL scheme")
	}

	host := u.Hostname()
	if host == "" {
		return fmt.Errorf("invalid hostname")
	}

	ips, err := net.LookupIP(host)
	if err != nil {
		return fmt.Errorf("failed to resolve host")
	}

	for _, ip := range ips {
		if helper.IsPrivateIP(ip) {
			return fmt.Errorf("access to private network denied")
		}
	}
	return nil
}

// fetchPage retrieves and parses the remote page
func fetchPage(urlStr string) (*goquery.Document, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", UserAgent)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}

	return goquery.NewDocumentFromReader(resp.Body)
}

// extractTitle gets the page <title>
func extractTitle(doc *goquery.Document) string {
	return strings.TrimSpace(doc.Find("title").First().Text())
}

// detectHTMLVersion infers the HTML version
func detectHTMLVersion(doc *goquery.Document) string {
	if doc.Find("article, aside, footer, header, main, nav, section").Length() > 0 {
		return "HTML5"
	}
	if doc.Find("html[xmlns]").Length() > 0 {
		return "XHTML"
	}
	return "HTML5"
}

// countHeadings returns a map of H1–H6 counts
func countHeadings(doc *goquery.Document) map[string]int {
	counts := make(map[string]int)
	for i := 1; i <= 6; i++ {
		h := fmt.Sprintf("h%d", i)
		counts[h] = doc.Find(h).Length()
	}
	return counts
}

// countLinks separates internal/external/inaccessible links
func countLinks(doc *goquery.Document, base *url.URL) (internal, external, inaccessible int) {
	var links []*url.URL

	doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		raw, ok := s.Attr("href")
		if !ok {
			return
		}

		ref, err := url.Parse(raw)
		if err != nil {
			return
		}

		if !ref.IsAbs() {
			ref = base.ResolveReference(ref)
		}
		links = append(links, ref)
	})

	if len(links) > maxLinksToCheck {
		links = links[:maxLinksToCheck]
	}

	results := checkLinksConcurrently(links)

	for _, res := range results {
		if res.isAccessible {
			if res.link.Hostname() == base.Hostname() {
				internal++
			} else {
				external++
			}
		} else {
			inaccessible++
		}
	}
	return
}

// detectLoginForm finds login forms based on common patterns
func detectLoginForm(doc *goquery.Document) bool {
	found := false
	doc.Find("form").EachWithBreak(func(i int, form *goquery.Selection) bool {
		hasPwd := form.Find("input[type='password']").Length() > 0
		hasUser := form.Find("input[type='text'], input[type='email'], input[name*='user'], input[id*='user']").Length() > 0
		formText := strings.ToLower(form.Text())

		if hasPwd && (hasUser || strings.Contains(formText, "login") || strings.Contains(formText, "sign in")) {
			found = true
			return false
		}
		return true
	})
	return found
}

type linkResult struct {
	link         *url.URL
	isAccessible bool
}

func checkLinksConcurrently(links []*url.URL) []linkResult {
	results := make([]linkResult, len(links))
	jobs := make(chan int, len(links))
	done := make(chan linkResult, len(links))

	for i := 0; i < maxWorkers; i++ {
		go func() {
			for idx := range jobs {
				link := links[idx]
				ok := isLinkAccessible(link)
				done <- linkResult{link: link, isAccessible: ok}
			}
		}()
	}

	for i := range links {
		jobs <- i
	}
	close(jobs)

	for i := 0; i < len(links); i++ {
		result := <-done
		results[i] = result
	}

	return results
}

func isLinkAccessible(link *url.URL) bool {
	if link.Scheme == "mailto" || link.Scheme == "tel" || link.Scheme == "javascript" {
		return false
	}
	if link.Scheme == "" && link.Host == "" && strings.HasPrefix(link.String(), "#") {
		return false
	}

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest(http.MethodHead, link.String(), nil)
	if err != nil {
		return false
	}
	req.Header.Set("User-Agent", UserAgent)

	res, err := client.Do(req)
	return err == nil && res.StatusCode < 400
}
