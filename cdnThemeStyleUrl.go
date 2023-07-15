package dashboard

import (
	"github.com/gouniverse/cdn"
)

func cdnThemeStyleUrl(themeName string) string {
	css := map[string]func() string{
		THEME_DEFAULT:   cdn.BootstrapCss_5_3_0,
		THEME_CERULEAN:  cdn.BootstrapCeruleanCss_5_3_0,
		THEME_COSMO:     cdn.BootstrapCosmoCss_5_3_0,
		THEME_CYBORG:    cdn.BootstrapCyborgCss_5_3_0,
		THEME_DARKLY:    cdn.BootstrapDarklyCss_5_3_0,
		THEME_FLATLY:    cdn.BootstrapFlatlyCss_5_3_0,
		THEME_JOURNAL:   cdn.BootstrapJournalCss_5_3_0,
		THEME_LITERA:    cdn.BootstrapLiteraCss_5_3_0,
		THEME_LUMEN:     cdn.BootstrapLumenCss_5_3_0,
		THEME_LUX:       cdn.BootstrapLuxCss_5_3_0,
		THEME_MATERIA:   cdn.BootstrapMateriaCss_5_3_0,
		THEME_MINTY:     cdn.BootstrapMintyCss_5_3_0,
		THEME_MORPH:     cdn.BootstrapMorphCss_5_3_0,
		THEME_PULSE:     cdn.BootstrapPulseCss_5_3_0,
		THEME_QUARTZ:    cdn.BootstrapQuartzCss_5_3_0,
		THEME_SANDSTONE: cdn.BootstrapSandstoneCss_5_3_0,
		THEME_SIMPLEX:   cdn.BootstrapSimplexCss_5_3_0,
		THEME_SKETCHY:   cdn.BootstrapSketchyCss_5_3_0,
		THEME_SLATE:     cdn.BootstrapSlateCss_5_3_0,
		THEME_SOLAR:     cdn.BootstrapSolarCss_5_3_0,
		THEME_SPACELAB:  cdn.BootstrapSpacelabCss_5_3_0,
		THEME_SUPERHERO: cdn.BootstrapSuperheroCss_5_3_0,
		THEME_UNITED:    cdn.BootstrapUnitedCss_5_3_0,
		THEME_VAPOR:     cdn.BootstrapVaporCss_5_3_0,
		THEME_YETI:      cdn.BootstrapYetiCss_5_3_0,
		THEME_ZEPHYR:    cdn.BootstrapZephyrCss_5_3_0,
	}
	if style, ok := css[themeName]; ok {
		return style()
	} else {
		return css[THEME_DEFAULT]()
	}
}
