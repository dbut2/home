package pr

import (
	"context"
	"strings"

	"github.com/anzx/fabric-actions/coverage/pkg/render"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type Client struct {
	client      *github.Client
	repoOwner   string
	repoName    string
	pullRequest int
	renderer    *render.Renderer
}

type Config interface {
	RepoOwner() string
	RepoName() string
	PullRequestNumber() int
}

func NewClient(ctx context.Context, token string, config Config, renderer *render.Renderer) *Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return &Client{
		client:      client,
		repoOwner:   config.RepoOwner(),
		repoName:    config.RepoName(),
		pullRequest: config.PullRequestNumber(),
		renderer:    renderer,
	}
}

// update existing 'warning' comment with a success message
func (c *Client) CommentSuccess(ctx context.Context, percent float64) {
	if commentID, exists := c.commentExists(ctx); exists {
		comment := c.renderer.ImprovedComment(percent)
		c.updateComment(ctx, commentID, comment)
	}
}

func (c *Client) CommentWarning(ctx context.Context, percent float64) {
	comment := c.renderer.WarningComment(percent)
	if commentID, exists := c.commentExists(ctx); exists {
		c.updateComment(ctx, commentID, comment)
	} else {
		c.postComment(ctx, comment)
	}
}

func (c *Client) DeleteCoverageCommentIfExists(ctx context.Context) {
	if commentID, exists := c.commentExists(ctx); exists {
		c.client.Issues.DeleteComment(ctx, c.repoOwner, c.repoName, commentID)
	}
}

func (c *Client) commentExists(ctx context.Context) (int64, bool) {
	comments, _, err := c.client.Issues.ListComments(ctx, c.repoOwner, c.repoName, c.pullRequest, &github.IssueListCommentsOptions{})

	if err != nil {
		panic(err)
	}
	if len(comments) < 1 {
		return 0, false
	}
	for _, comment := range comments {
		if strings.HasPrefix(*comment.Body, render.CommentMetadata) {
			return *comment.ID, true
		}
	}
	return 0, false
}

func (c *Client) updateComment(ctx context.Context, commentID int64, commentStr string) {
	comment := &github.IssueComment{
		Body: &commentStr,
	}
	_, _, err := c.client.Issues.EditComment(ctx, c.repoOwner, c.repoName, commentID, comment)
	if err != nil {
		panic(err)
	}
}

func (c *Client) postComment(ctx context.Context, commentStr string) {
	comment := &github.IssueComment{
		Body: &commentStr,
	}
	_, _, err := c.client.Issues.CreateComment(ctx, c.repoOwner, c.repoName, c.pullRequest, comment)
	if err != nil {
		panic(err)
	}
}
