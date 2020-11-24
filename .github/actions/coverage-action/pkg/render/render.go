package render

import (
	"bytes"
	"html/template"
	"net/url"

	"github.com/mazznoer/colorgrad"
)

// CommentMetadata is used to identify if I've already commented.
// This could be used to store state between runs, but for now it's static.
const CommentMetadata = "<!-- coverage-warning-comment"

type commentInfo struct {
	Config      RequiredTemplateVars
	Percent     float64
	Color       string
	PullRequest int
	IssueBody   string
	RunID       string
	SonarID     string
}

type RequiredTemplateVars interface {
	RepoOwner() string
	RepoName() string
	PullRequestNumber() int
	RunID() string
}

type Renderer struct {
	cfg         RequiredTemplateVars
	templateDir string
}

func NewRenderer(templateDir string, templateVars RequiredTemplateVars) *Renderer {
	return &Renderer{
		cfg:         templateVars,
		templateDir: templateDir,
	}
}

func (r *Renderer) ImprovedComment(percent float64) string {
	return r.renderComment(percent, "comment-improved.md")
}

func (r *Renderer) WarningComment(percent float64) string {
	return r.renderComment(percent, "comment-warning.md")
}

func (r *Renderer) renderComment(percent float64, templateFile string) string {
	info := &commentInfo{
		Config:      r.cfg,
		Percent:     percent,
		Color:       colorForPercent(percent),
		PullRequest: r.cfg.PullRequestNumber(),
		RunID:       r.cfg.RunID(),
		SonarID:     r.cfg.RepoName(),
		IssueBody:   url.QueryEscape(r.renderIssueBody(percent)),
	}

	var b bytes.Buffer
	t := template.Must(template.ParseFiles(r.templateDir + templateFile))
	t.Execute(&b, info)
	return CommentMetadata + " -->\n" + b.String()
}

func (r *Renderer) renderIssueBody(percent float64) string {
	info := &commentInfo{
		Config:      r.cfg,
		Percent:     percent,
		Color:       colorForPercent(percent),
		PullRequest: r.cfg.PullRequestNumber(),
		RunID:       r.cfg.RunID(),
		SonarID:     r.cfg.RepoName(),
	}

	var b bytes.Buffer
	t := template.Must(template.ParseFiles(r.templateDir + "issue.md"))
	t.Execute(&b, info)
	return b.String()
}

func colorForPercent(percent float64) string {
	grad, err := colorgrad.NewGradient().
		HexColors(
			// Using turbo color scale from https://observablehq.com/@d3/color-schemes
			// "#900c00", "#ba2208", "#f65f18", "#ffa423", "#dedd32", "#95fb51", // "#4df884"
			// darken 25%
			"#900c00", "#ba2208", "#f65f18", "#ffa423", "#dedd32", "#64f305", // "#4df884"
		).
		Domain(0, (100 / 5), (100/5)*2, (100/5)*3, (100/5)*4, 100).
		Build()
	if err != nil {
		panic(err)
	}
	return grad.At(percent).Hex()[1:]
}
