package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	flags "github.com/jessevdk/go-flags"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
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
	f.Value = p[1]
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

func postMessage(url string, json []byte) error {
	req, err := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer(json),
	)
	if err != nil {
		log.Println(err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: time.Duration(3) * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return err
	}
	if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		log.Fatal(resp.Status)
	} else if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("Response status code: %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	return nil
}

func buildJSON(text string, opts Options) []byte {
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

	if opts.DisableMarkdown {
		m.Markdown = false
	} else {
		m.Markdown = true
	}

	// b, err := json.Marshal(m)
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	return b
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
		bytes, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			panic(err)
		}
		text = strings.TrimRight(string(bytes), "\n")
	} else {
		text = opts.Message
	}

	b := buildJSON(text, opts)

	if opts.Debug {
		fmt.Println(string(b))
	} else {
		if opts.Args.Url == "" {
			log.Fatal("the required argument `Url` was not provided")
		}
		for i := 0; i < 3; i++ {
			err := postMessage(opts.Args.Url, b)
			if err == nil {
				break
			}
		}
	}
}
