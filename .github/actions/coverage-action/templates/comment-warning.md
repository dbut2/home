Hey, test coverage has dropped below {{.Config.SoftTarget}}%. Good news is that it's still above {{.Config.HardTarget}}% so you can merge this PR, but you might want to track this as a [new issue](https://github.com/{{.Config.Repository}}/issues/new?title=Improve+unit+test+coverage&body={{.IssueBody}}).

![Coverage: {{.Percent | printf "%.1f"}}%](https://img.shields.io/badge/Coverage-{{.Percent | printf "%.1f"}}%25-{{.Color}})

<details><summary>What can I do to improve test coverage?</summary>
<p>

1. :bookmark_tabs:&nbsp;&nbsp;See the [files tab](https://github.com/{{.Config.Repository}}/pull/{{.PullRequest}}/files), functions that have 0% coverage have been annotated with a warning.
2. :beetle:&nbsp;&nbsp;[Create an issue to improve test coverage](https://github.com/{{.Config.Repository}}/issues/new?title=Improve+unit+test+coverage&body={{.IssueBody}}). All the details have already been filled in.
3. :eyes:&nbsp;&nbsp;Check out the [action results](https://github.com/{{.Config.Repository}}/actions/runs/{{.RunID}}).
5. :gear:&nbsp;&nbsp;[Configure](https://github.com/{{.Config.Repository}}/tree/master/.github/workflows) this warning.

</p>
</details>

I'll update this comment if anything changes.
