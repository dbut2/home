package pages

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPage_Display(t *testing.T) {
	p := Page{Title: "title"}

	require.NoError(t, p.Display(nil))
}
