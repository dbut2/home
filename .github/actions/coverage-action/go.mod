module github.com/anzx/fabric-actions/coverage

go 1.13

require (
	github.com/google/go-github v17.0.0+incompatible
	github.com/google/go-querystring v1.0.0 // indirect
	github.com/mazznoer/colorgrad v0.2.0
	github.com/mbndr/figlet4go v0.0.0-20190224160619-d6cef5b186ea
	github.com/sethvargo/go-envconfig v0.2.2
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
)

replace github.com/anzx/fabric-actions/coverage/pkg => ./pkg
