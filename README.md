# Page Insight Tool

A web application that accepts a URL input, analyzes the contents of the corresponding web page, and displays specific information about it.

## 🚀 Features

- **Web Form Interface**: Clean, modern web form for URL input
- **HTML Analysis**: Extracts HTML version, page title, and heading structure
- **Link Analysis**: Counts internal vs external links and inaccessible links
- **Security Analysis**: Detects login forms and provides security insights
- **Error Handling**: Graceful handling of network errors, malformed URLs, and security violations
- **SSRF Protection**: Blocks access to private networks and internal IPs
- **Responsive Design**: Modern, mobile-friendly interface

## 📋 Technical Requirements

- **Backend**: Go 1.21+
- **Frontend**: Server-side rendered HTML with Go templates
- **Dependencies**: 
  - `github.com/PuerkitoBio/goquery` - HTML parsing
  - `github.com/gorilla/mux` - HTTP routing
  - `gopkg.in/yaml.v2` - Configuration management

## 🛠️ Build and Run Instructions

### Prerequisites

- Go 1.21 or higher
- Git

### Installation

1. **Clone the repository**:
   ```bash
   git clone <repository-url>
   cd page-insight-tool
   ```

2. **Install dependencies**:
   ```bash
   go mod tidy
   ```

3. **Build the application**:
   ```bash
   make build
   ```

4. **Run the application**:
   ```bash
   make run
   ```

   Or run directly with Go:
   ```bash
   go run app/cmd/page-insight-tool.go
   ```

5. **Access the application**:
   Open your browser and navigate to `http://localhost:8080`

### Development

For development with hot reloading:
```bash
make dev
```

## 🏗️ Project Structure

Following this architecture pattern:

```
page-insight-tool/
├── app/
│   ├── cmd/
│   │   └── page-insight-tool.go    # Application entry point
│   ├── config/
│   │   ├── config.go               # Configuration management
│   │   └── page-insight-tool.yaml  # YAML configuration file
│   ├── handlers/
│   │   └── analyze.go              # HTTP handlers and analysis logic
│   ├── helper/
│   │   └── network.go              # helper function for networking checks
│   ├── router/
│   │   └── router.go               # HTTP routing setup
│   ├── static/
│   │   └── style.css               # CSS styles for the interface
│   ├── templates/
│   │   └── index.html              # HTML template for the web interface
│   └── .bin/                       # Build artifacts (gitignored)
├── go.mod                          # Go module dependencies
├── go.sum                          # Dependency checksums
├── Makefile                        # Build and development commands
├── Dockerfile                      # Container configuration
├── .gitignore                      # Git ignore patterns
└── README.md                       # This file
```

## 🔧 Configuration

The application supports multiple configuration methods:

### Environment Variables
- `PORT`: Server port (default: 8080)
- `DEBUG`: Enable debug logging (default: false)


### YAML Configuration
Create a configuration file following the pattern in `app/config/page-insight-tool.yaml`:

```yaml
Local:
  Host: localhost
  Port: "8080"
```

### Command Line Options
- `--config`: Path to configuration file
- `--debug`: Enable debug logging

## 🔒 Security Features

### SSRF Protection
- **Private IP Blocking**: Prevents access to private network ranges (10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16)
- **Localhost Protection**: Blocks access to 127.0.0.0/8 and ::1
- **Link-local Protection**: Blocks 169.254.0.0/16 range
- **Scheme Validation**: Only allows HTTP and HTTPS protocols

### Input Validation
- **URL Format Validation**: Ensures proper URL structure
- **Hostname Resolution**: Validates hostname can be resolved
- **Timeout Protection**: Configurable timeout for all requests

## 🎯 Assumptions and Design Decisions

### Security Assumptions
- **Network Access**: Only HTTP/HTTPS schemes are allowed
- **Private Networks**: Access to private IP ranges is blocked for security
- **Timeouts**: Configurable timeout enforced on all HTTP requests
- **User Agent**: Custom user agent to avoid being blocked by websites

### Technical Decisions
- **Server-Side Rendering**: Uses Go templates for simplicity and SEO
- **No JavaScript**: Pure server-side implementation for reliability
- **Responsive Design**: Mobile-first CSS approach
- **Error Handling**: User-friendly error messages for all failure scenarios
- **Modular Architecture**: Clean separation of concerns

### Analysis Limitations
- **JavaScript Content**: Cannot analyze dynamically generated content
- **Authentication**: Cannot access password-protected pages
- **Rate Limiting**: May be blocked by sites with strict rate limiting
- **Redirects**: Follows redirects but may not capture all redirect chains

## 🧪 Testing

Run the test suite:
```bash
make test
```
Test app running :
```bash
make test-app
```

## 📈 Future Improvements
