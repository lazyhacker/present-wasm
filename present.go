// Copyright 2012 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build js

package main

import (
	"bytes"
	"html/template"
	"strings"
	"syscall/js"

	"golang.org/x/tools/present"
	"lazyhackergo.com/browser"
)

func main() {

	save := js.NewEventCallback(false, false, false, cbSave)
	defer save.Close()

	present := js.NewEventCallback(false, false, false, cbPresent)
	defer present.Close()

	window := browser.GetWindow()

	getPresentation()

	window.Document.GetElementById("saveButton").AddEventListener(browser.EventClick, save)
	window.Document.GetElementById("presentButton").AddEventListener(browser.EventClick, present)

	window.Document.GetElementById("saveButton").SetProperty("disabled", false)
	window.Document.GetElementById("presentButton").SetProperty("disabled", false)
	keepalive()
}

func getPresentation() {
	window := browser.GetWindow()
	storage := window.LocalStorage

	v := storage.GetItem("preso")

	m := window.Document.GetElementById("markdown")
	m.SetValue(v)
}

func cbSave(e js.Value) {

	window := browser.GetWindow()
	storage := window.LocalStorage

	v := window.Document.GetElementById("markdown")

	storage.SetItem("preso", v.Value())
}

func cbPresent(e js.Value) {

	present.PlayEnabled = false
	present.NotesEnabled = false

	// Initialize the slide template.
	tmpl := present.Template()
	tmpl = tmpl.Funcs(template.FuncMap{"playable": playable})
	tmpl.Parse(action_template)
	tmpl.Parse(slide_template)

	window := browser.GetWindow()
	storage := window.LocalStorage

	slide := storage.GetItem("preso")

	if len(slide) == 0 {
		storage.SetItem("preso", default_slides)
		slide = storage.GetItem("preso")
	}
	// convert to a io.Reader
	reader := strings.NewReader(slide)

	// create an io.Writer for renderer to output to
	out := new(bytes.Buffer)

	doc, _ := present.Parse(reader, "root", 0)
	err := doc.Render(out, tmpl)
	if err != nil {
		println("error rendering: " + err.Error())
	}

	w := browser.NewWindow("")
	w.Document.Open()
	w.Document.Write(out.String())
	w.Document.Close()
}

func keepalive() {
	select {}
}

const default_slides = `
Title of document
Subtitle of document
15:04 2 Jan 2006
Tags: foo, bar, baz

Author Name
Job title, Company
joe@example.com
http://url/
@twitter_name
Some Text

* Title of slide or section (must have asterisk)

Some Text
`

const slide_template = `
{/* This is the slide template. It defines how presentations are formatted. */}

{{define "root"}}
<!DOCTYPE html>
<html>
  <head>
    <title>{{.Title}}</title>
    <meta charset='utf-8'>
    <script>
      var notesEnabled = {{.NotesEnabled}};
    </script>
    <script src='static/slides.js'></script>

    {{if .NotesEnabled}}
    <script>
      var sections = {{.Sections}};
      var titleNotes = {{.TitleNotes}}
    </script>
    <script src='/static/notes.js'></script>
    {{end}}

    <script>
      // Initialize Google Analytics tracking code on production site only.
      if (window["location"] && window["location"]["hostname"] == "talks.golang.org") {
        var _gaq = _gaq || [];
        _gaq.push(["_setAccount", "UA-11222381-6"]);
        _gaq.push(["b._setAccount", "UA-49880327-6"]);
        window.trackPageview = function() {
          _gaq.push(["_trackPageview", location.pathname+location.hash]);
          _gaq.push(["b._trackPageview", location.pathname+location.hash]);
        };
        window.trackPageview();
        window.trackEvent = function(category, action, opt_label, opt_value, opt_noninteraction) {
          _gaq.push(["_trackEvent", category, action, opt_label, opt_value, opt_noninteraction]);
          _gaq.push(["b._trackEvent", category, action, opt_label, opt_value, opt_noninteraction]);
        };
      }
    </script>
  </head>

  <body style='display: none'>

    <section class='slides layout-widescreen'>

      <article>
        <h1>{{.Title}}</h1>
        {{with .Subtitle}}<h3>{{.}}</h3>{{end}}
        {{if not .Time.IsZero}}<h3>{{.Time.Format "2 January 2006"}}</h3>{{end}}
        {{range .Authors}}
          <div class="presenter">
            {{range .TextElem}}{{elem $.Template .}}{{end}}
          </div>
        {{end}}
      </article>

  {{range $i, $s := .Sections}}
  <!-- start of slide {{$s.Number}} -->
      <article>
      {{if $s.Elem}}
        <h3>{{$s.Title}}</h3>
        {{range $s.Elem}}{{elem $.Template .}}{{end}}
      {{else}}
        <h2>{{$s.Title}}</h2>
      {{end}}
      </article>
  <!-- end of slide {{$i}} -->
  {{end}}{{/* of Slide block */}}

      <article>
        <h3>Thank you</h3>
        {{range .Authors}}
          <div class="presenter">
            {{range .Elem}}{{elem $.Template .}}{{end}}
          </div>
        {{end}}
      </article>

    </section>

    <div id="help">
      Use the left and right arrow keys or click the left and right
      edges of the page to navigate between slides.<br>
      (Press 'H' or navigate to hide this message.)
    </div>

    {{if .PlayEnabled}}
    <script src='/play.js'></script>
    {{end}}

    <script>
      (function() {
        // Load Google Analytics tracking code on production site only.
        if (window["location"] && window["location"]["hostname"] == "talks.golang.org") {
          var ga = document.createElement("script"); ga.type = "text/javascript"; ga.async = true;
          ga.src = ("https:" == document.location.protocol ? "https://ssl" : "http://www") + ".google-analytics.com/ga.js";
          var s = document.getElementsByTagName("script")[0]; s.parentNode.insertBefore(ga, s);
        }
      })();
    </script>
  </body>
</html>
{{end}}

{{define "newline"}}
<br>
{{end}}
`

const action_template = `
{/*
This is the action template.
It determines how the formatting actions are rendered.
*/}

{{define "section"}}
  <h{{len .Number}} id="TOC_{{.FormattedNumber}}">{{.FormattedNumber}} {{.Title}}</h{{len .Number}}>
  {{range .Elem}}{{elem $.Template .}}{{end}}
{{end}}

{{define "list"}}
  <ul>
  {{range .Bullet}}
    <li>{{style .}}</li>
  {{end}}
  </ul>
{{end}}

{{define "text"}}
  {{if .Pre}}
  <div class="code"><pre>{{range .Lines}}{{.}}{{end}}</pre></div>
  {{else}}
  <p>
    {{range $i, $l := .Lines}}{{if $i}}{{template "newline"}}
    {{end}}{{style $l}}{{end}}
  </p>
  {{end}}
{{end}}

{{define "code"}}
  <div class="code{{if playable .}} playground{{end}}" {{if .Edit}}contenteditable="true" spellcheck="false"{{end}}>{{.Text}}</div>
{{end}}

{{define "image"}}
<div class="image">
  <img src="{{.URL}}"{{with .Height}} height="{{.}}"{{end}}{{with .Width}} width="{{.}}"{{end}}>
</div>
{{end}}

{{define "video"}}
<div class="video">
  <video {{with .Height}} height="{{.}}"{{end}}{{with .Width}} width="{{.}}"{{end}} controls>
    <source src="{{.URL}}" type="{{.SourceType}}">
  </video>
</div>
{{end}}

{{define "background"}}
<div class="background">
  <img src="{{.URL}}">
</div>
{{end}}

{{define "iframe"}}
<iframe src="{{.URL}}"{{with .Height}} height="{{.}}"{{end}}{{with .Width}} width="{{.}}"{{end}}></iframe>
{{end}}

{{define "link"}}<p class="link"><a href="{{.URL}}" target="_blank">{{style .Label}}</a></p>{{end}}

{{define "html"}}{{.HTML}}{{end}}

{{define "caption"}}<figcaption>{{style .Text}}</figcaption>{{end}}
`
