package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/anzx/fabric-actions/coverage/pkg/pr"
	"github.com/anzx/fabric-actions/coverage/pkg/render"
	"github.com/mbndr/figlet4go"
	"github.com/sethvargo/go-envconfig"
)

var cfg config

func init() {
	if err := envconfig.Process(context.Background(), &cfg); err != nil {
		log.Fatal(err)
	}
}

func main() {
	fmt.Printf("::debug::config: %+v\n", cfg)
	reader := bufio.NewReader(os.Stdin)
	for {
		rowStr, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		row := parseRow(rowStr)

		// if row.file is empty, ignore, go to the next row
		if row.file == "" {
			continue
		}

		// "total" is the last row and the percent is the average of all the function-level coverage
		if row.file == "total" {
			finish(row)
			// we're done, stop processing the file
			break
		}

		// if the function is not covered by tests at all, write a warning
		if row.percent == 0 {
			fmt.Printf("::warning file=%s,line=%s::`%s` not covered by tests\n", row.file, row.line, row.funcName)
		}
	}
}

func finish(row *row /* your boat */) {
	fmt.Printf("::set-output name=coverage::%.1f%%\n", row.percent)
	bigOutput(row.percent)

	ctx := context.Background()
	renderer := render.NewRenderer("/coverage/templates/", cfg)
	client := pr.NewClient(ctx, cfg.GitHubToken, cfg, renderer)

	if row.percent >= cfg.SoftTarget {
		// The check has passed and we're in the 'success' range. If we've previously commented
		// we'll update the comment with a nice message, otherwise do nothing.
		if cfg.EnableWarningComment {
			client.CommentSuccess(ctx, row.percent)
		}
	}
	if row.percent < cfg.SoftTarget {
		fmt.Printf("::warning::Coverage %.1f%% below soft target %.1f%%\n", row.percent, cfg.SoftTarget)

		// The check has passed, but we're in the 'warning' range so if commenting is enabled,
		// post a comment to notify the user that coverage could be improved.
		if cfg.EnableWarningComment && row.percent >= cfg.HardTarget {
			client.CommentWarning(ctx, row.percent)
		}
	}
	if row.percent < cfg.HardTarget {
		fmt.Printf("::error::Coverage %.1f%% below hard target %.1f%%\n", row.percent, cfg.HardTarget)

		// The check has failed, and in the case where I've commented previously, delete the comment.
		// This is a deliberate design decision to reduce comment noise, to increase the likelyhood
		// that the warning will not be ignored. The only time a comment will occur on a Pull Request
		// is the following:
		//
		//   1. The coverage is in the 'warning' range
		//   2. The coverage was in the 'warning' range and is now in the 'success' range
		//
		// That means there will be no comment for 'fails', they are blocked from merging anyway so
		// any comment will just be unnecessary noise.
		if cfg.EnableWarningComment {
			client.DeleteCoverageCommentIfExists(ctx)
		}

		// Below the the hard target, so a non-zero exit code fails the step
		os.Exit(1)
	}
}

func bigOutput(percent float64) {
	ascii := figlet4go.NewAsciiRender()
	ascii.LoadFont("/coverage/fonts/")
	options := figlet4go.NewRenderOptions()
	options.FontName = "colossal"
	renderStr, _ := ascii.RenderOpts(fmt.Sprintf("%.1f%%", percent), options)
	fmt.Println("\n\n")
	fmt.Print(renderStr)
}

type row struct {
	file     string
	line     string
	funcName string
	percent  float64
}

// matches the output from "go tool cover -func=cover.out"
var rx = regexp.MustCompile("(.*?)(:\\d+)?:\\s+(.*?)\\s+(\\d?\\d?\\d.\\d)%")

func parseRow(rowStr string) *row {
	r := &row{}
	m := rx.FindStringSubmatch(rowStr)
	if len(m) == 0 {
		return r
	}
	r.file = strings.ReplaceAll(m[1], "github.com/"+cfg.Repository, "")
	if len(m[2]) > 0 {
		r.line = m[2][1:]
	}
	r.funcName = m[3]
	pct, err := strconv.ParseFloat(m[4], 64)
	if err != nil {
		panic(err)
	}
	r.percent = pct
	return r
}
