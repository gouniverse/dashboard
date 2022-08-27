package dashboard

import (
	"strings"

	"github.com/gouniverse/hb"
)

func scripts(scriptsArray []string) string {
	requiredScripts := []string{
		// templates.Jquery360(),
		// templates.BootstrapJs513(),
		// templates.VueJs3(),
		// templates.WebJs260(),
		// templates.SharedJs(),
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
