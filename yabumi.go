// yabumi は Slack の Incoming Webhook にメッセージを投稿するコマンドラインツール
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	flags "github.com/jessevdk/go-flags"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

const (
	httpTimeout   = 3 * time.Second
	retryCount    = 3
	retryBaseWait = time.Second
	// Retry-After で指定された待ち時間の上限
	maxRetryAfter = 30 * time.Second
	// エラーメッセージに含めるレスポンスボディの最大バイト数
	maxErrorBodySize = 512
)

var httpClient = &http.Client{Timeout: httpTimeout}

// webhookURLEnv は Webhook URL を指定する環境変数名
const webhookURLEnv = "SLACK_WEBHOOK_URL"

// webhookURL は引数で指定された Webhook URL を返す。未指定の場合は環境変数から取得する
func webhookURL(opts Options) string {
	if opts.Args.URL != "" {
		return opts.Args.URL
	}
	return os.Getenv(webhookURLEnv)
}

// sleep はテストで差し替えられるようにする
var sleep = time.Sleep

type Options struct {
	Channel         string   `short:"C" long:"channel" description:"slack channel to post"`
	UseAttach       bool     `short:"a" long:"attachment" description:"use attachment"`
	Title           string   `short:"t" long:"title" description:"Title text (attachment)"`
	TitleLink       string   `long:"title-link" description:"Title link url (attachment)"`
	Color           string   `short:"c" long:"color" description:"Color code or 'good', 'warning', 'danger' (attachment)"`
	PreText         string   `short:"p" long:"pretext" description:"optional text that appears above the message attachment block (attachment)"`
	AuthorName      string   `long:"author-name" description:"author_name (attachment)"`
	AuthorLink      string   `long:"author-link" description:"author_link (attachment)"`
	AuthorIcon      string   `long:"author-icon" description:"author_icon (attachment)"`
	ImageURL        string   `long:"image-url" description:"image url (attachment)"`
	ThumbURL        string   `long:"thumb-url" description:"thumbnail image url (attachment)"`
	Footer          string   `long:"footer" description:"footer text (attachment)"`
	FooterIcon      string   `long:"footer-icon" description:"footer icon url (attachment)"`
	Message         string   `short:"m" long:"message" description:"pass message instead of read stdin"`
	Fields          []string `short:"f" long:"field" description:"\"title|value|short\" (attachment)"`
	DisableMarkdown bool     `short:"M" long:"disable-markdown" description:"disable markdown processing"`
	Debug           bool     `short:"D" long:"debug" description:"enable debug mode. do not send request, show json only"`
	Version         bool     `short:"v" long:"version" description:"show version"`
	Args            struct {
		URL string `positional-arg-name:"Url" description:"slack webhook endpoint url (default: $SLACK_WEBHOOK_URL)"`
	} `positional-args:"yes"`
}

type SlackMessage struct {
	Text        string       `json:"text,omitempty"`
	Channel     string       `json:"channel,omitempty"`
	Markdown    bool         `json:"mrkdwn"`
	Attachments []Attachment `json:"attachments,omitempty"`
}

// https://api.slack.com/docs/message-attachments
type Attachment struct {
	Fallback   string  `json:"fallback,omitempty"`
	Color      string  `json:"color,omitempty"`
	Pretext    string  `json:"pretext,omitempty"`
	AuthorName string  `json:"author_name,omitempty"`
	AuthorLink string  `json:"author_link,omitempty"` // URL
	AuthorIcon string  `json:"author_icon,omitempty"` // URL
	Title      string  `json:"title,omitempty"`
	TitleLink  string  `json:"title_link,omitempty"` // URL
	Text       string  `json:"text,omitempty"`
	Fields     []Field `json:"fields,omitempty"`
	ImageURL   string  `json:"image_url,omitempty"` // URL
	ThumbURL   string  `json:"thumb_url,omitempty"` // URL
	Footer     string  `json:"footer,omitempty"`
	FooterIcon string  `json:"footer_icon,omitempty"` // URL
}

type Field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short,omitempty"`
}

func parseField(s string) Field {
	p := strings.SplitN(s, "|", 3)
	var f Field
	f.Title = p[0]
	if len(p) > 1 {
		f.Value = p[1]
	}
	if len(p) > 2 {
		f.Short = parseBool(p[2])
	}
	return f
}

func parseBool(s string) bool {
	switch strings.ToLower(s) {
	case "0", "false", "":
		return false
	default:
		return true
	}
}

// permanentError はリトライすべきでないエラーを表す
type permanentError struct {
	err error
}

func (e *permanentError) Error() string { return e.err.Error() }
func (e *permanentError) Unwrap() error { return e.err }

// retryAfterError はサーバーから待ち時間が指定されたリトライ可能なエラーを表す
type retryAfterError struct {
	err   error
	after time.Duration
}

func (e *retryAfterError) Error() string { return e.err.Error() }
func (e *retryAfterError) Unwrap() error { return e.err }

// parseRetryAfter は Retry-After ヘッダ (秒数または HTTP-date) を解釈する
func parseRetryAfter(v string, now time.Time) time.Duration {
	if v == "" {
		return 0
	}
	if sec, err := strconv.Atoi(v); err == nil {
		if sec < 0 {
			return 0
		}
		return time.Duration(sec) * time.Second
	}
	if t, err := http.ParseTime(v); err == nil {
		return max(t.Sub(now), 0)
	}
	return 0
}

func sendWithRetry(logger *log.Logger, endpoint string, body []byte, retries int, baseWait time.Duration) error {
	var lastErr error
	for i := range retries {
		if i > 0 {
			wait := baseWait * (1 << (i - 1))
			var ra *retryAfterError
			if errors.As(lastErr, &ra) {
				wait = max(wait, min(ra.after, maxRetryAfter))
			}
			logger.Printf("waiting %v before retry...", wait)
			sleep(wait)
		}
		lastErr = postMessage(endpoint, body)
		if lastErr == nil {
			return nil
		}
		logger.Printf("attempt %d failed: %v", i+1, lastErr)
		var pe *permanentError
		if errors.As(lastErr, &pe) {
			return lastErr
		}
	}
	return lastErr
}

// redactURL は err に含まれる Webhook URL (秘密情報) を伏せる
func redactURL(err error) error {
	var ue *url.Error
	if errors.As(err, &ue) {
		ue.URL = "<redacted>"
	}
	return err
}

// statusError はレスポンスステータスとボディ (Slack のエラー理由) を含むエラーを返す
func statusError(resp *http.Response) error {
	b, _ := io.ReadAll(io.LimitReader(resp.Body, maxErrorBodySize+1))
	truncated := len(b) > maxErrorBodySize
	if truncated {
		b = b[:maxErrorBodySize]
	}
	// HTML 等の複数行のボディでもログが 1 行に収まるよう空白をまとめる
	body := strings.Join(strings.Fields(strings.ToValidUTF8(string(b), "")), " ")
	if body == "" {
		return fmt.Errorf("unexpected response status: %s", resp.Status)
	}
	if truncated {
		body += "..."
	}
	return fmt.Errorf("unexpected response status: %s: %s", resp.Status, body)
}

func postMessage(endpoint string, json []byte) error {
	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		endpoint,
		bytes.NewBuffer(json),
	)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", redactURL(err))
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", redactURL(err))
	}
	defer resp.Body.Close() //nolint:errcheck
	if resp.StatusCode == http.StatusTooManyRequests {
		return &retryAfterError{
			err:   statusError(resp),
			after: parseRetryAfter(resp.Header.Get("Retry-After"), time.Now()),
		}
	}
	if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		return &permanentError{err: statusError(resp)}
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return statusError(resp)
	}

	return nil
}

func buildJSON(text string, opts Options) ([]byte, error) {
	var m SlackMessage
	m.Channel = opts.Channel

	if opts.UseAttach {
		var a Attachment
		a.Fallback = text
		a.Text = text
		a.Title = opts.Title
		a.TitleLink = opts.TitleLink
		a.Color = opts.Color
		a.Pretext = opts.PreText
		if len(opts.Fields) > 0 {
			for _, field := range opts.Fields {
				a.Fields = append(a.Fields, parseField(field))
			}
		}
		a.AuthorName = opts.AuthorName
		a.AuthorLink = opts.AuthorLink
		a.AuthorIcon = opts.AuthorIcon
		a.ImageURL = opts.ImageURL
		a.ThumbURL = opts.ThumbURL
		a.Footer = opts.Footer
		a.FooterIcon = opts.FooterIcon
		m.Attachments = append(m.Attachments, a)
	} else {
		m.Text = text
	}

	m.Markdown = !opts.DisableMarkdown

	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return b, nil
}

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

// run は CLI 本体。テストできるよう入出力を引数で受け取り、終了コードを返す
func run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	logger := log.New(stderr, "", log.LstdFlags)

	var opts Options
	parser := flags.NewParser(&opts, flags.HelpFlag|flags.PassDoubleDash)
	if _, err := parser.ParseArgs(args); err != nil {
		if flags.WroteHelp(err) {
			fmt.Fprintln(stdout, err) //nolint:errcheck
			return 0
		}
		fmt.Fprintln(stderr, err) //nolint:errcheck
		return 1
	}

	if opts.Version {
		fmt.Fprintln(stdout, "yabumi (Post message to slack)") //nolint:errcheck
		fmt.Fprintln(stdout, "version:", version)              //nolint:errcheck
		fmt.Fprintln(stdout, "commit:", commit)                //nolint:errcheck
		fmt.Fprintln(stdout, "build date:", date)              //nolint:errcheck
		return 0
	}

	var text string
	if opts.Message == "" {
		b, err := io.ReadAll(stdin)
		if err != nil {
			logger.Print(err)
			return 1
		}
		text = strings.TrimRight(string(b), "\n")
	} else {
		text = opts.Message
	}

	b, err := buildJSON(text, opts)
	if err != nil {
		logger.Print(err)
		return 1
	}

	if opts.Debug {
		fmt.Fprintln(stdout, string(b)) //nolint:errcheck
		return 0
	}

	endpoint := webhookURL(opts)
	if endpoint == "" {
		logger.Printf("webhook url is not specified: pass it as an argument or set %s", webhookURLEnv)
		return 1
	}
	if err := sendWithRetry(logger, endpoint, b, retryCount, retryBaseWait); err != nil {
		logger.Printf("failed to post message: %v", err)
		return 1
	}
	return 0
}
