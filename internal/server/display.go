package server

import (
	"html/template"
	"net/http"

	"github.com/dbut2/home/pkg/log"
	"github.com/dbut2/home/pkg/pages"
)

func PageDisplay(p *pages.Page, w http.ResponseWriter) error {
	t, err := template.New("Page").Parse(page)
	if err != nil {
		log.Error(err)
		return err
	}
	err = t.Execute(w, p)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func PageDisplayEdit(p *pages.Page, w http.ResponseWriter) error {
	t, err := template.New("EditPage").Parse(editPage)
	if err != nil {
		log.Error(err)
		return err
	}
	err = t.Execute(w, p)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

const (
	page = `
<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8" />
<title>{{ .Title }}</title>
<link rel="stylesheet" href="/static/css/style.css" />
<link rel="stylesheet" href="/static/css/md.css" />
</head>
<body>
<div id="container" class="markdown-body">
<div><div>{{ .ParseContent }}</div></div>
</div>
</body>
</html>`

	editPage = `
<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8" />
<title>{{ .Title }}</title>
<link rel="stylesheet" href="/static/css/style.css" />
<link rel="stylesheet" href="/static/css/md.css" />
<link rel="stylesheet" href="/static/css/edit.css" />
<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bootstrap@4.5.3/dist/css/bootstrap.min.css" integrity="sha384-TX8t27EcRE3e/ihU7zmQxVncDAy5uIKz4rEkgIXeMed4M0jlfIDPvg6uqKI2xXr2" crossorigin="anonymous">
</head>
<body>
<div id="container" class="markdown-body">
<div>
<div>
<form action="{{ .Name }}" method="POST">
<div class="float-left mb-3">
<a class="btn btn-primary" href="/edit">Back</a>
</div>
<div class="float-right mb-3">
<label class="form-check-label">Visible:</label>
<input type="checkbox" name="visible"{{ if .Visible }} checked{{ end }}>
<input class="btn btn-primary" type="submit" value="{{ if .Name }}Update{{ else }}Post{{ end }}" />
</div>
<div class="input-group mb-3">
<div class="input-group-prepend">
<span class="input-group-text">dylanbutler.net/</span>
</div>
<input class="form-control" type="text" name="name" value="{{ .Name }}" placeholder="link" />
</div>
<div class="mb-3">
<input class="form-control" type="text" name="title" value="{{ .Title }}" placeholder="Title" />
</div>
<div class="mb-3">
<textarea class="form-control" name="content" placeholder="Content">{{ .Content }}</textarea>
</div>
</form>
</div>
</div>
<div><div>{{ .ParseContent }}</div></div>
</div>
</body>
</html>`
)
