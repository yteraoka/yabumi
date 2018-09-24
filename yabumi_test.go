package main

import (
	flags "github.com/jessevdk/go-flags"
	"github.com/bitly/go-simplejson"
	"testing"
)

func TestParseField(t *testing.T) {
	v := parseField("name|value|true")
	if v.Title != "name" || v.Value != "value" || v.Short != true {
		t.Error("parse failed")
	}

	v = parseField("name|value|0")
	if v.Title != "name" || v.Value != "value" || v.Short != false {
		t.Error("parse failed")
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
	b := buildJSON(text, opts)
	js, err := simplejson.NewJson(b)
	js_text, err := js.Get("text").String()
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
	b := buildJSON(text, opts)
	js, err := simplejson.NewJson(b)

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
	b := buildJSON(text, opts)
	js, err := simplejson.NewJson(b)

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
	if "Environment" != js_f_title1 {
		t.Errorf("field[0]['title'] in JSON is unexpected: %s != %s", "Environment", js_f_title1)
	}

	js_f_value1, err := js.Get("attachments").GetIndex(0).Get("fields").GetIndex(0).Get("value").String()
	if err != nil {
		t.Error(err)
	}
	if "production" != js_f_value1 {
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
