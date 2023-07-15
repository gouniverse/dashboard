package dashboard

import (
	"strings"

	"github.com/gouniverse/hb"
	"github.com/gouniverse/uncdn"
)

func scripts(scriptsArray []string) string {
	requiredScripts := []string{
		uncdn.BootstrapJs523(),
	}

	scriptsArray = append(requiredScripts, scriptsArray...)
	scripts := ""
	for _, script := range scriptsArray {
		if strings.HasPrefix(script, "http") || strings.HasPrefix(script, "//") {
			script = hb.NewScriptURL(script).ToHTML()
		} else if !strings.HasPrefix(script, "<script") {
			script = hb.NewScript(script).ToHTML()
		}
		scripts += script
	}
	return scripts
}
