package dashboard

import (
	"github.com/gouniverse/uncdn"
)

func uncdnThemeStyle(themeName string) string {
	css := map[string]func() string{
		THEME_DEFAULT:   uncdn.BootstrapCss523,
		THEME_CERULEAN:  uncdn.BootstrapCeruleanCss523,
		THEME_COSMO:     uncdn.BootstrapCosmoCss523,
		THEME_CYBORG:    uncdn.BootstrapCyborgCss523,
		THEME_DARKLY:    uncdn.BootstrapDarklyCss523,
		THEME_FLATLY:    uncdn.BootstrapFlatlyCss523,
		THEME_JOURNAL:   uncdn.BootstrapJournalCss523,
		THEME_LITERA:    uncdn.BootstrapLiteraCss523,
		THEME_LUMEN:     uncdn.BootstrapLumenCss523,
		THEME_LUX:       uncdn.BootstrapLuxCss523,
		THEME_MATERIA:   uncdn.BootstrapMateriaCss523,
		THEME_MINTY:     uncdn.BootstrapMintyCss523,
		THEME_MORPH:     uncdn.BootstrapMorphCss523,
		THEME_PULSE:     uncdn.BootstrapPulseCss523,
		THEME_QUARTZ:    uncdn.BootstrapQuartzCss523,
		THEME_SANDSTONE: uncdn.BootstrapSandstoneCss523,
		THEME_SIMPLEX:   uncdn.BootstrapSimplexCss523,
		THEME_SKETCHY:   uncdn.BootstrapSketchyCss523,
		THEME_SLATE:     uncdn.BootstrapSlateCss523,
		THEME_SOLAR:     uncdn.BootstrapSolarCss523,
		THEME_SPACELAB:  uncdn.BootstrapSpacelabCss523,
		THEME_SUPERHERO: uncdn.BootstrapSuperheroCss523,
		THEME_UNITED:    uncdn.BootstrapUnitedCss523,
		THEME_VAPOR:     uncdn.BootstrapVaporCss523,
		THEME_YETI:      uncdn.BootstrapYetiCss523,
		THEME_ZEPHYR:    uncdn.BootstrapZephyrCss523,
	}
	if style, ok := css[themeName]; ok {
		return style()
	} else {
		return css[THEME_DEFAULT]()
	}
}
