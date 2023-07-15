package dashboard

import (
	"net/http"
	"strings"

	"github.com/gouniverse/responses"
	"github.com/gouniverse/uncdn"
	"github.com/samber/lo"
)

func UncdnHandler(w http.ResponseWriter, r *http.Request) {
	c := &cdnController{}
	required, extension := c.findRequiredAndExtension(r)

	if len(required) < 1 {
		responses.HTMLResponse(w, r, "Nothing requested")
		return
	}
	if extension == "" {
		responses.HTMLResponse(w, r, "No extension provided")
		return
	}
	if !lo.Contains([]string{"css", "js"}, extension) {
		w.Write([]byte("Extension " + extension + " not supported"))
		return
	}

	if extension == "js" {
		responses.GzipJSResponse(w, r, c.compileJS(w, required))
		return
	}

	if extension == "css" {
		responses.GzipCSSResponse(w, r, c.compileCSS(w, required))
		return
	}

	responses.HTMLResponse(w, r, "Extension "+extension+" not found")
}

type cdnController struct{}

func (c cdnController) compileJS(w http.ResponseWriter, required []string) string {
	js := []string{}
	lo.ForEach(required, func(item string, index int) {
		js = append(js, findJS(item))
	})
	return strings.Join(js, "\n;\n")
}

func (c cdnController) compileCSS(w http.ResponseWriter, required []string) string {
	css := []string{}
	lo.ForEach(required, func(item string, index int) {
		css = append(css, c.findCss(item))
	})
	return strings.Join(css, "\n\n")
}

func findJS(required string) string {
	if required == "jq360" {
		return uncdn.Jquery360()
	}

	if required == "bs523" {
		return uncdn.BootstrapJs523()
	}

	if required == "vue3" {
		return uncdn.VueJs3()
	}

	if required == "web260" {
		return uncdn.WebJs260()
	}

	if required == "ntf" {
		return uncdn.NotifyJs()
	}

	if required == "swal" {
		return uncdn.Sweetalert2_11432()
	}

	return ""
}

func (c cdnController) findCss(required string) string {
	css := map[string]func() string{
		"bs523":                  uncdn.BootstrapCss523,
		"bs523" + THEME_CERULEAN: uncdn.BootstrapCeruleanCss523,
		"bs523" + THEME_COSMO:    uncdn.BootstrapCosmoCss523,
		"bs523cyborg":            uncdn.BootstrapCyborgCss523,
		"bs523darkly":            uncdn.BootstrapDarklyCss523,
		"bs523flatly":            uncdn.BootstrapFlatlyCss523,
		"bs523journal":           uncdn.BootstrapJournalCss523,
		"bs523litera":            uncdn.BootstrapLiteraCss523,
		"bs523lumen":             uncdn.BootstrapLumenCss523,
		"bs523lux":               uncdn.BootstrapLuxCss523,
		"bs523materia":           uncdn.BootstrapMateriaCss523,
		"bs523minty":             uncdn.BootstrapMintyCss523,
		"bs523morph":             uncdn.BootstrapMorphCss523,
		"bs523pulse":             uncdn.BootstrapPulseCss523,
		"bs523quartz":            uncdn.BootstrapQuartzCss523,
		"bs523sandstone":         uncdn.BootstrapSandstoneCss523,
		"bs523simplex":           uncdn.BootstrapSimplexCss523,
		"bs523sketchy":           uncdn.BootstrapSketchyCss523,
		"bs523slate":             uncdn.BootstrapSlateCss523,
		"bs523solar":             uncdn.BootstrapSolarCss523,
		"bs523spacelab":          uncdn.BootstrapSpacelabCss523,
		"bs523superhero":         uncdn.BootstrapSuperheroCss523,
		"bs523united":            uncdn.BootstrapUnitedCss523,
		"bs523vapor":             uncdn.BootstrapVaporCss523,
		"bs523yeti":              uncdn.BootstrapYetiCss523,
		"bs523" + THEME_ZEPHYR:   uncdn.BootstrapZephyrCss523,
	}
	if style, ok := css[required]; ok {
		return style()
	}
	return ""
}

func (c cdnController) findRequiredAndExtension(req *http.Request) (required []string, extension string) {
	uriParts := strings.Split(strings.Trim(req.RequestURI, "/"), "/")
	name := lo.Ternary(len(uriParts) > 1, uriParts[len(uriParts)-1], "")
	if name == "" {
		return []string{}, ""
	}

	nameParts := strings.Split(name, ".")

	if len(nameParts) < 2 {
		return []string{}, ""
	}

	requiredArray := strings.Split(nameParts[0], "-")

	return requiredArray, nameParts[1]
}
