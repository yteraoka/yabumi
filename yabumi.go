package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	flags "github.com/jessevdk/go-flags"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

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
	ImageUrl        string   `long:"image-url" description:"image url (attachment)"`
	ThumbUrl        string   `long:"thumb-url" description:"thumbnail image url (attachment)"`
	Footer          string   `long:"footer" description:"footer text (attachment)"`
	FooterIcon      string   `long:"footer-icon" description:"footer icon url (attachment)"`
	Message         string   `short:"m" long:"message" description:"pass message instead of read stdin"`
	Fields          []string `short:"f" long:"field" description:"\"title|value|short\" (attachment)"`
	DisableMarkdown bool     `short:"M" long:"disable-markdown" description:"disable markdown processing"`
	Debug           bool     `short:"D" long:"debug" description:"enable debug mode. do not send request, show json only"`
	Version         bool     `short:"v" long:"version" description:"show version"`
	Args            struct {
		Url string `description:"slack webhook endpoint url"`
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
	ImageUrl   string  `json:"image_url,omitempty"` // URL
	ThumbUrl   string  `json:"thumb_url,omitempty"` // URL
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
	var result bool
	switch strings.ToLower(s) {
	case "0":
		result = false
	case "false":
		result = false
	case "":
		result = false
	default:
		result = true
	}
	return result
}

// permanentError はリトライすべきでないエラーを表す
type permanentError struct {
	err error
}

func (e *permanentError) Error() string { return e.err.Error() }
func (e *permanentError) Unwrap() error { return e.err }

func sendWithRetry(url string, body []byte, retries int, baseWait time.Duration) error {
	var lastErr error
	for i := 0; i < retries; i++ {
		if i > 0 {
			wait := baseWait * (1 << (i - 1))
			log.Printf("waiting %v before retry...", wait)
			time.Sleep(wait)
		}
		lastErr = postMessage(url, body)
		if lastErr == nil {
			return nil
		}
		log.Printf("attempt %d failed: %v", i+1, lastErr)
		var pe *permanentError
		if errors.As(lastErr, &pe) {
			return lastErr
		}
	}
	return lastErr
}

func postMessage(url string, json []byte) error {
	req, err := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer(json),
	)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: time.Duration(3) * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck
	if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		return &permanentError{err: fmt.Errorf("unexpected response status: %s", resp.Status)}
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected response status: %s", resp.Status)
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
		a.Color = opts.Color
		if len(opts.Fields) > 0 {
			for _, field := range opts.Fields {
				a.Fields = append(a.Fields, parseField(field))
			}
		}
		a.AuthorName = opts.AuthorName
		a.AuthorLink = opts.AuthorLink
		a.AuthorIcon = opts.AuthorIcon
		a.ImageUrl = opts.ImageUrl
		a.ThumbUrl = opts.ThumbUrl
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
	var opts Options
	var text string
	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}

	if opts.Version {
		fmt.Println("yabumi (Post message to slack)")
		fmt.Println("version:", version)
		fmt.Println("commit:", commit)
		fmt.Println("build date:", date)
		os.Exit(0)
	}

	if opts.Message == "" {
		bytes, err := io.ReadAll(os.Stdin)
		if err != nil {
			log.Fatal(err)
		}
		text = strings.TrimRight(string(bytes), "\n")
	} else {
		text = opts.Message
	}

	b, err := buildJSON(text, opts)
	if err != nil {
		log.Fatal(err)
	}

	if opts.Debug {
		fmt.Println(string(b))
	} else {
		if opts.Args.Url == "" {
			log.Fatal("the required argument `Url` was not provided")
		}
		if err := sendWithRetry(opts.Args.Url, b, 3, time.Second); err != nil {
			log.Fatal("all attempts failed")
		}
	}
}
