Test coverage has dropped below {{.Config.SoftTarget}}% and is currently at {{.Percent}}%.

![Coverage: {{.Percent | printf "%.1f"}}%](https://img.shields.io/badge/Coverage-{{.Percent | printf "%.1f"}}%25-{{.Color}})

Related: #{{.PullRequest}}
