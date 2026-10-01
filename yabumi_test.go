package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	flags "github.com/jessevdk/go-flags"
)

var discardLogger = log.New(io.Discard, "", 0)

// closedServerURL は接続が即座に拒否される URL を返す
func closedServerURL(t *testing.T) string {
	t.Helper()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	u := ts.URL
	ts.Close()
	return u
}

func TestBuildJSON(t *testing.T) {
	cases := []struct {
		name     string
		args     []string
		text     string
		expected SlackMessage
	}{
		{
			name:     "simple text",
			text:     "hello",
			expected: SlackMessage{Text: "hello", Markdown: true},
		},
		{
			name:     "channel",
			args:     []string{"--channel", "#mychannel"},
			text:     "hello",
			expected: SlackMessage{Text: "hello", Channel: "#mychannel", Markdown: true},
		},
		{
			name:     "disable markdown",
			args:     []string{"--disable-markdown"},
			text:     "hello",
			expected: SlackMessage{Text: "hello", Markdown: false},
		},
		{
			name: "attachment with title",
			args: []string{"--attachment", "--title", "test title"},
			text: "hello",
			expected: SlackMessage{
				Markdown: true,
				Attachments: []Attachment{
					{Fallback: "hello", Text: "hello", Title: "test title"},
				},
			},
		},
		{
			name: "attachment with all options",
			args: []string{
				"--channel", "#mychannel",
				"--attachment",
				"--title", "title",
				"--title-link", "https://example.com/title",
				"--color", "danger",
				"--pretext", "pretext",
				"--author-name", "author",
				"--author-link", "https://example.com/author",
				"--author-icon", "https://example.com/author.png",
				"--image-url", "https://example.com/image.png",
				"--thumb-url", "https://example.com/thumb.png",
				"--footer", "footer",
				"--footer-icon", "https://example.com/footer.png",
				"--field", "Environment|production|true",
				"--field", "Service|test|0",
				"--field", "Note",
			},
			text: "hello",
			expected: SlackMessage{
				Channel:  "#mychannel",
				Markdown: true,
				Attachments: []Attachment{
					{
						Fallback:   "hello",
						Text:       "hello",
						Title:      "title",
						TitleLink:  "https://example.com/title",
						Color:      "danger",
						Pretext:    "pretext",
						AuthorName: "author",
						AuthorLink: "https://example.com/author",
						AuthorIcon: "https://example.com/author.png",
						ImageUrl:   "https://example.com/image.png",
						ThumbUrl:   "https://example.com/thumb.png",
						Footer:     "footer",
						FooterIcon: "https://example.com/footer.png",
						Fields: []Field{
							{Title: "Environment", Value: "production", Short: true},
							{Title: "Service", Value: "test", Short: false},
							{Title: "Note"},
						},
					},
				},
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var opts Options
			if _, err := flags.ParseArgs(&opts, c.args); err != nil {
				t.Fatal(err)
			}
			b, err := buildJSON(c.text, opts)
			if err != nil {
				t.Fatal(err)
			}
			var got SlackMessage
			if err := json.Unmarshal(b, &got); err != nil {
				t.Fatalf("failed to parse JSON: %v (%s)", err, b)
			}
			if !reflect.DeepEqual(got, c.expected) {
				t.Errorf("unexpected message:\n got: %+v\nwant: %+v", got, c.expected)
			}
		})
	}
}

func TestParseField(t *testing.T) {
	cases := []struct {
		input    string
		expected Field
	}{
		{"name|value|true", Field{Title: "name", Value: "value", Short: true}},
		{"name|value|0", Field{Title: "name", Value: "value", Short: false}},
		{"name|value", Field{Title: "name", Value: "value"}},
		{"name", Field{Title: "name"}},
		{"name|a|true|b", Field{Title: "name", Value: "a", Short: true}},
	}
	for _, c := range cases {
		if got := parseField(c.input); got != c.expected {
			t.Errorf("parseField(%q) = %+v, want %+v", c.input, got, c.expected)
		}
	}
}

func TestParseBool(t *testing.T) {
	cases := []struct {
		input    string
		expected bool
	}{
		{"1", true},
		{"true", true},
		{"TRUE", true},
		{"yes", true},
		{"0", false},
		{"false", false},
		{"FALSE", false},
		{"", false},
	}
	for _, c := range cases {
		if got := parseBool(c.input); got != c.expected {
			t.Errorf("parseBool(%q) = %v, want %v", c.input, got, c.expected)
		}
	}
}

func TestPostMessageStatus(t *testing.T) {
	cases := []struct {
		name       string
		status     int
		wantErr    bool
		permanent  bool // permanentError (リトライしない) であるべきか
		retryAfter bool // retryAfterError であるべきか
	}{
		{name: "200 OK", status: http.StatusOK},
		{name: "204 No Content", status: http.StatusNoContent},
		{name: "400 Bad Request", status: http.StatusBadRequest, wantErr: true, permanent: true},
		{name: "404 Not Found", status: http.StatusNotFound, wantErr: true, permanent: true},
		{name: "429 Too Many Requests", status: http.StatusTooManyRequests, wantErr: true, retryAfter: true},
		{name: "500 Internal Server Error", status: http.StatusInternalServerError, wantErr: true},
		{name: "503 Service Unavailable", status: http.StatusServiceUnavailable, wantErr: true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var gotContentType string
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotContentType = r.Header.Get("Content-Type")
				w.WriteHeader(c.status)
			}))
			defer ts.Close()

			err := postMessage(ts.URL, []byte(`{"text":"hello"}`))
			if gotContentType != "application/json" {
				t.Errorf("Content-Type = %q, want application/json", gotContentType)
			}
			if !c.wantErr {
				if err != nil {
					t.Errorf("expected no error, got: %v", err)
				}
				return
			}
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !strings.Contains(err.Error(), strconv.Itoa(c.status)) {
				t.Errorf("expected error to contain status code, got: %v", err)
			}
			var pe *permanentError
			if got := errors.As(err, &pe); got != c.permanent {
				t.Errorf("permanentError = %v, want %v", got, c.permanent)
			}
			var ra *retryAfterError
			if got := errors.As(err, &ra); got != c.retryAfter {
				t.Errorf("retryAfterError = %v, want %v", got, c.retryAfter)
			}
		})
	}
}

func TestPostMessageNetworkError(t *testing.T) {
	err := postMessage(closedServerURL(t), []byte(`{"text":"hello"}`))
	if err == nil {
		t.Error("expected error for unreachable server, got nil")
	}
}

func TestSendWithRetry(t *testing.T) {
	cases := []struct {
		name          string
		statuses      []int // 各試行で返すステータス (足りない分は最後の値)
		wantErr       bool
		expectedCalls int
	}{
		{name: "success", statuses: []int{200}, expectedCalls: 1},
		{name: "succeeds on second attempt", statuses: []int{500, 200}, expectedCalls: 2},
		{name: "succeeds on last attempt", statuses: []int{500, 503, 200}, expectedCalls: 3},
		{name: "all fail", statuses: []int{500}, wantErr: true, expectedCalls: 3},
		{name: "no retry on 4xx", statuses: []int{400}, wantErr: true, expectedCalls: 1},
		{name: "4xx after 5xx", statuses: []int{500, 400}, wantErr: true, expectedCalls: 2},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			called := 0
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(c.statuses[min(called, len(c.statuses)-1)])
				called++
			}))
			defer ts.Close()

			err := sendWithRetry(discardLogger, ts.URL, []byte(`{}`), 3, 0)
			if (err != nil) != c.wantErr {
				t.Errorf("error = %v, wantErr %v", err, c.wantErr)
			}
			if called != c.expectedCalls {
				t.Errorf("expected %d calls, got %d", c.expectedCalls, called)
			}
		})
	}
}

func TestSendWithRetryBackoff(t *testing.T) {
	var waits []time.Duration
	orig := sleep
	sleep = func(d time.Duration) { waits = append(waits, d) }
	defer func() { sleep = orig }()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	if err := sendWithRetry(discardLogger, ts.URL, []byte(`{}`), 4, time.Second); err == nil {
		t.Error("expected error when all attempts fail, got nil")
	}
	expected := []time.Duration{time.Second, 2 * time.Second, 4 * time.Second}
	if !reflect.DeepEqual(waits, expected) {
		t.Errorf("waits = %v, want %v", waits, expected)
	}
}

func TestPostMessageErrorDoesNotContainURL(t *testing.T) {
	const secret = "T000/B000/secret-token"
	cases := []struct {
		name string
		url  string
	}{
		{"network error", closedServerURL(t) + "/services/" + secret},
		{"invalid url", "http://[::1/services/" + secret},
	}
	for _, c := range cases {
		err := postMessage(c.url, []byte(`{"text":"hello"}`))
		if err == nil {
			t.Errorf("%s: expected error, got nil", c.name)
			continue
		}
		if strings.Contains(err.Error(), secret) {
			t.Errorf("%s: error message contains webhook url: %v", c.name, err)
		}
	}
}

func TestParseRetryAfter(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	cases := []struct {
		input    string
		expected time.Duration
	}{
		{"", 0},
		{"5", 5 * time.Second},
		{"0", 0},
		{"-1", 0},
		{"invalid", 0},
		{now.Add(10 * time.Second).Format(http.TimeFormat), 10 * time.Second},
		{now.Add(-10 * time.Second).Format(http.TimeFormat), 0},
	}
	for _, c := range cases {
		if got := parseRetryAfter(c.input, now); got != c.expected {
			t.Errorf("parseRetryAfter(%q) = %v, want %v", c.input, got, c.expected)
		}
	}
}

func TestSendWithRetryOn429(t *testing.T) {
	cases := []struct {
		name       string
		retryAfter string
		expected   time.Duration
	}{
		{"without Retry-After", "", time.Second},
		{"with Retry-After", "5", 5 * time.Second},
		{"Retry-After shorter than backoff", "0", time.Second},
		{"Retry-After exceeds max", "3600", maxRetryAfter},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var waits []time.Duration
			orig := sleep
			sleep = func(d time.Duration) { waits = append(waits, d) }
			defer func() { sleep = orig }()

			called := 0
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				called++
				if called < 2 {
					if c.retryAfter != "" {
						w.Header().Set("Retry-After", c.retryAfter)
					}
					w.WriteHeader(http.StatusTooManyRequests)
				} else {
					w.WriteHeader(http.StatusOK)
				}
			}))
			defer ts.Close()

			if err := sendWithRetry(discardLogger, ts.URL, []byte(`{}`), 3, time.Second); err != nil {
				t.Errorf("expected no error, got: %v", err)
			}
			if called != 2 {
				t.Errorf("expected 2 calls, got %d", called)
			}
			if len(waits) != 1 || waits[0] != c.expected {
				t.Errorf("expected wait [%v], got %v", c.expected, waits)
			}
		})
	}
}

func TestPostMessageErrorContainsResponseBody(t *testing.T) {
	long := strings.Repeat("a", maxErrorBodySize+100)
	cases := []struct {
		name     string
		status   int
		body     string
		expected string
	}{
		{"slack error", http.StatusBadRequest, "invalid_payload", "unexpected response status: 400 Bad Request: invalid_payload"},
		{"rate limited", http.StatusTooManyRequests, "rate_limited", "unexpected response status: 429 Too Many Requests: rate_limited"},
		{"server error", http.StatusInternalServerError, "internal_error", "unexpected response status: 500 Internal Server Error: internal_error"},
		{"empty body", http.StatusBadRequest, "", "unexpected response status: 400 Bad Request"},
		{"multi-line body", http.StatusBadGateway, "<html>\n  <body>Bad Gateway</body>\n</html>\n", "unexpected response status: 502 Bad Gateway: <html> <body>Bad Gateway</body> </html>"},
		{"long body", http.StatusBadRequest, long, "unexpected response status: 400 Bad Request: " + long[:maxErrorBodySize] + "..."},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(c.status)
				_, _ = w.Write([]byte(c.body))
			}))
			defer ts.Close()

			err := postMessage(ts.URL, []byte(`{"text":"hello"}`))
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if err.Error() != c.expected {
				t.Errorf("unexpected error message:\n got: %v\nwant: %v", err, c.expected)
			}
		})
	}
}

func TestWebhookURL(t *testing.T) {
	cases := []struct {
		name     string
		arg      string
		env      string
		expected string
	}{
		{"argument only", "https://example.com/arg", "", "https://example.com/arg"},
		{"env only", "", "https://example.com/env", "https://example.com/env"},
		{"argument takes precedence", "https://example.com/arg", "https://example.com/env", "https://example.com/arg"},
		{"neither", "", "", ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Setenv(webhookURLEnv, c.env)
			var opts Options
			opts.Args.Url = c.arg
			if got := webhookURL(opts); got != c.expected {
				t.Errorf("webhookURL() = %q, want %q", got, c.expected)
			}
		})
	}
}

func TestRun(t *testing.T) {
	var received []byte
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		received, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()
	failTS := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("invalid_payload"))
	}))
	defer failTS.Close()

	cases := []struct {
		name         string
		args         []string
		stdin        string
		env          string
		expectedCode int
		stdout       string // stdout に含まれるべき文字列
		stderr       string // stderr に含まれるべき文字列
		sentText     string // 送信された JSON の text (空なら送信を検証しない)
	}{
		{name: "help", args: []string{"--help"}, expectedCode: 0, stdout: "Usage:"},
		{name: "unknown flag", args: []string{"--bogus"}, expectedCode: 1, stderr: "unknown flag `bogus'"},
		{name: "version", args: []string{"--version"}, expectedCode: 0, stdout: "version: dev"},
		{name: "debug reads stdin", args: []string{"-D"}, stdin: "hello\n", expectedCode: 0, stdout: `"text": "hello"`},
		{name: "debug with message", args: []string{"-D", "-m", "from flag"}, stdin: "from stdin", expectedCode: 0, stdout: `"text": "from flag"`},
		{name: "no url", args: []string{"-m", "hi"}, expectedCode: 1, stderr: "set SLACK_WEBHOOK_URL"},
		{name: "send via argument", args: []string{"-m", "hi", ts.URL}, expectedCode: 0, sentText: "hi"},
		{name: "send via env", args: []string{}, stdin: "from env\n", env: ts.URL, expectedCode: 0, sentText: "from env"},
		{name: "send fails", args: []string{"-m", "hi", failTS.URL}, expectedCode: 1, stderr: "failed to post message: unexpected response status: 400 Bad Request: invalid_payload"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Setenv(webhookURLEnv, c.env)
			received = nil
			var stdout, stderr bytes.Buffer

			code := run(c.args, strings.NewReader(c.stdin), &stdout, &stderr)
			if code != c.expectedCode {
				t.Errorf("exit code = %d, want %d (stderr: %s)", code, c.expectedCode, stderr.String())
			}
			if !strings.Contains(stdout.String(), c.stdout) {
				t.Errorf("stdout does not contain %q: %s", c.stdout, stdout.String())
			}
			if !strings.Contains(stderr.String(), c.stderr) {
				t.Errorf("stderr does not contain %q: %s", c.stderr, stderr.String())
			}
			if c.sentText != "" {
				var m SlackMessage
				if err := json.Unmarshal(received, &m); err != nil {
					t.Fatalf("failed to parse sent JSON: %v (%s)", err, received)
				}
				if m.Text != c.sentText {
					t.Errorf("sent text = %q, want %q", m.Text, c.sentText)
				}
			}
		})
	}
}
