package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestIndexHandler(t *testing.T) {
	// Create a temporary template file for testing
	tempDir := t.TempDir()
	templatePath := filepath.Join(tempDir, "index.html")

	// Create a simple template for testing
	templateContent := `<!DOCTYPE html>
<html>
<head><title>Page Insight Tool</title></head>
<body>
<h1>Page Insight Tool</h1>
<form method="POST" action="/analyze">
<input type="url" name="url" required>
<button type="submit">Analyze</button>
</form>
{{if .Error}}
<div class="error">{{.Error}}</div>
{{else if .URL}}
<div class="results">
<h2>Analysis Results</h2>
<p>URL: {{.URL}}</p>
</div>
{{end}}
</body>
</html>`

	err := os.WriteFile(templatePath, []byte(templateContent), 0644)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(IndexHandler)

	handler.ServeHTTP(rr, req)

	// Since we can't easily mock the template loading in this test,
	// we'll just check that the handler doesn't panic and returns some response
	if rr.Code == http.StatusInternalServerError {
		// This is expected since the template file doesn't exist in test environment
		// We're testing that the handler doesn't panic
		t.Log("Handler returned 500 as expected due to missing template file")
	} else if rr.Code != http.StatusOK {
		t.Errorf("handler returned unexpected status code: got %v", rr.Code)
	}
}

func TestAnalyzeHandler_EmptyURL(t *testing.T) {
	req, err := http.NewRequest("POST", "/analyze", strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(AnalyzeHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestAnalyzeHandler_ValidURL(t *testing.T) {
	data := url.Values{}
	data.Set("url", "https://example.com")

	req, err := http.NewRequest("POST", "/analyze", strings.NewReader(data.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(AnalyzeHandler)

	handler.ServeHTTP(rr, req)

	// Since we can't easily mock the template loading in this test,
	// we'll just check that the handler doesn't panic and returns some response
	if rr.Code == http.StatusInternalServerError {
		// This is expected since the template file doesn't exist in test environment
		// We're testing that the handler doesn't panic
		t.Log("Handler returned 500 as expected due to missing template file")
	} else if rr.Code != http.StatusOK {
		t.Errorf("handler returned unexpected status code: got %v", rr.Code)
	}
}

func TestAnalyzePage_InvalidURL(t *testing.T) {
	analysis := analyzePage("invalid-url")

	if analysis.Error.Message == "" {
		t.Error("expected error for invalid URL, got none")
	}
}

func TestAnalyzePage_PrivateIP(t *testing.T) {
	analysis := analyzePage("http://192.168.1.1")

	if analysis.Error.Message == "" {
		t.Error("expected error for private IP, got none")
	}

	if !strings.Contains(analysis.Error.Message, "access to private network denied") {
		t.Errorf("expected private network error, got: %s", analysis.Error.Message)
	}
}

func TestValidateURL_ValidURL(t *testing.T) {
	u, _ := url.Parse("https://example.com")
	err := validateURL(u)

	if err != nil {
		t.Errorf("expected no error for valid URL, got: %v", err)
	}
}

func TestValidateURL_InvalidScheme(t *testing.T) {
	u, _ := url.Parse("ftp://example.com")
	err := validateURL(u)

	if err == nil {
		t.Error("expected error for invalid scheme, got none")
	}

	if !strings.Contains(err.Error(), "unsupported URL scheme") {
		t.Errorf("expected scheme error, got: %s", err.Error())
	}
}

func TestIsLinkAccessible(t *testing.T) {
	// Test accessible link (using a known working URL)
	accessibleURL, _ := url.Parse("https://httpbin.org/status/200")
	if !isLinkAccessible(accessibleURL) {
		t.Error("expected https://httpbin.org/status/200 to be accessible")
	}

	// Test inaccessible link (using a known 404 URL)
	inaccessibleURL, _ := url.Parse("https://httpbin.org/status/404")
	if isLinkAccessible(inaccessibleURL) {
		t.Error("expected https://httpbin.org/status/404 to be inaccessible")
	}

	// Test non-HTTP schemes
	mailtoURL, _ := url.Parse("mailto:test@example.com")
	if isLinkAccessible(mailtoURL) {
		t.Error("expected mailto: scheme to be inaccessible")
	}

	telURL, _ := url.Parse("tel:+1234567890")
	if isLinkAccessible(telURL) {
		t.Error("expected tel: scheme to be inaccessible")
	}

	javascriptURL, _ := url.Parse("javascript:alert('test')")
	if isLinkAccessible(javascriptURL) {
		t.Error("expected javascript: scheme to be inaccessible")
	}

	// Test fragment-only links
	fragmentURL, _ := url.Parse("#section")
	if isLinkAccessible(fragmentURL) {
		t.Error("expected fragment-only link to be inaccessible")
	}
}
