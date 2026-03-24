package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bitly/go-simplejson"
	flags "github.com/jessevdk/go-flags"
)

func TestBuildJSONMarkdown(t *testing.T) {
	// デフォルトは markdown 有効
	b, err := buildJSON("hello", Options{})
	if err != nil {
		t.Fatal(err)
	}
	js, _ := simplejson.NewJson(b)
	if mrkdwn, _ := js.Get("mrkdwn").Bool(); !mrkdwn {
		t.Error("expected mrkdwn to be true by default")
	}

	// DisableMarkdown を指定すると無効になる
	b, err = buildJSON("hello", Options{DisableMarkdown: true})
	if err != nil {
		t.Fatal(err)
	}
	js, _ = simplejson.NewJson(b)
	if mrkdwn, _ := js.Get("mrkdwn").Bool(); mrkdwn {
		t.Error("expected mrkdwn to be false when DisableMarkdown is set")
	}
}

func TestParseField(t *testing.T) {
	v := parseField("name|value|true")
	if v.Title != "name" || v.Value != "value" || v.Short != true {
		t.Error("parse failed")
	}

	v = parseField("name|value|0")
	if v.Title != "name" || v.Value != "value" || v.Short != false {
		t.Error("parse failed")
	}

	v = parseField("name|value")
	if v.Title != "name" || v.Value != "value" || v.Short != false {
		t.Error("parse failed for title|value format")
	}

	v = parseField("name")
	if v.Title != "name" || v.Value != "" || v.Short != false {
		t.Error("parse failed for title-only format")
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

func TestPostMessageSuccess(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	err := postMessage(ts.URL, []byte(`{"text":"hello"}`))
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestPostMessageClientError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer ts.Close()

	err := postMessage(ts.URL, []byte(`{"text":"hello"}`))
	if err == nil {
		t.Error("expected error for 4xx response, got nil")
	}
	if !strings.Contains(err.Error(), "400") {
		t.Errorf("expected error to contain status code, got: %v", err)
	}
}

func TestPostMessageServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	err := postMessage(ts.URL, []byte(`{"text":"hello"}`))
	if err == nil {
		t.Error("expected error for 5xx response, got nil")
	}
}

func TestPostMessageNetworkError(t *testing.T) {
	err := postMessage("http://127.0.0.1:1", []byte(`{"text":"hello"}`))
	if err == nil {
		t.Error("expected error for unreachable server, got nil")
	}
}

func TestBuildJSON1(t *testing.T) {
	args := []string{}
	text := "test"
	var opts Options
	_, err := flags.ParseArgs(&opts, args)
	if err != nil {
		t.Error(err)
	}
	b, err := buildJSON(text, opts)
	if err != nil {
		t.Fatal(err)
	}
	js, err := simplejson.NewJson(b)
	if err != nil {
		t.Error(err)
	}
	js_text, err := js.Get("text").String()
	if err != nil {
		t.Error(err)
	}
	if text != js_text {
		t.Errorf("text in JSON is unexpected: %s != %s", text, js_text)
	}
}

func TestBuildJSON2(t *testing.T) {
	text := "test"
	title := "test title"
	args := []string{
		"--attachment",
		"--title", title,
	}
	var opts Options
	_, err := flags.ParseArgs(&opts, args)
	if err != nil {
		t.Error(err)
	}
	b, err := buildJSON(text, opts)
	if err != nil {
		t.Fatal(err)
	}
	js, err := simplejson.NewJson(b)
	if err != nil {
		t.Error(err)
	}

	js_text, err := js.Get("attachments").GetIndex(0).Get("text").String()
	if err != nil {
		t.Error(err)
	}
	if text != js_text {
		t.Errorf("text in JSON is unexpected: %s != %s", text, js_text)
	}

	js_title, err := js.Get("attachments").GetIndex(0).Get("title").String()
	if err != nil {
		t.Error(err)
	}
	if title != js_title {
		t.Errorf("title in JSON is unexpected: %s != %s", title, js_title)
	}
}

func TestBuildJSON3(t *testing.T) {
	channel := "#mychannel"
	text := "test"
	title := "test title"
	color := "daner"
	args := []string{
		"--channel", channel,
		"--attachment",
		"--title", title,
		"--color", color,
		"--field", "Environment|production|true",
		"--field", "Service|test|1",
	}
	var opts Options
	_, err := flags.ParseArgs(&opts, args)
	if err != nil {
		t.Error(err)
	}
	b, err := buildJSON(text, opts)
	if err != nil {
		t.Fatal(err)
	}
	js, err := simplejson.NewJson(b)
	if err != nil {
		t.Error(err)
	}

	js_channel, err := js.Get("channel").String()
	if err != nil {
		t.Error(err)
	}
	if channel != js_channel {
		t.Errorf("channel in JSON is unexpected: %s != %s", channel, js_channel)
	}

	js_text, err := js.Get("attachments").GetIndex(0).Get("text").String()
	if err != nil {
		t.Error(err)
	}
	if text != js_text {
		t.Errorf("text in JSON is unexpected: %s != %s", text, js_text)
	}

	js_title, err := js.Get("attachments").GetIndex(0).Get("title").String()
	if err != nil {
		t.Error(err)
	}
	if title != js_title {
		t.Errorf("title in JSON is unexpected: %s != %s", title, js_title)
	}

	js_color, err := js.Get("attachments").GetIndex(0).Get("color").String()
	if err != nil {
		t.Error(err)
	}
	if color != js_color {
		t.Errorf("color in JSON is unexpected: %s != %s", color, js_color)
	}

	js_f_title1, err := js.Get("attachments").GetIndex(0).Get("fields").GetIndex(0).Get("title").String()
	if err != nil {
		t.Error(err)
	}
	if js_f_title1 != "Environment" {
		t.Errorf("field[0]['title'] in JSON is unexpected: %s != %s", "Environment", js_f_title1)
	}

	js_f_value1, err := js.Get("attachments").GetIndex(0).Get("fields").GetIndex(0).Get("value").String()
	if err != nil {
		t.Error(err)
	}
	if js_f_value1 != "production" {
		t.Errorf("field[0]['value'] in JSON is unexpected: %s != %s", "production", js_f_value1)
	}

	js_f_short1, err := js.Get("attachments").GetIndex(0).Get("fields").GetIndex(0).Get("short").Bool()
	if err != nil {
		t.Error(err)
	}
	if true != js_f_short1 {
		t.Errorf("field[0]['short'] in JSON is unexpected: %s != %v", "true", js_f_short1)
	}
}
