package i18n

import (
	"embed"
	"log"

	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"
	toml "github.com/pelletier/go-toml/v2"
	"golang.org/x/text/language"
)

//go:embed locales/*.toml
var localeFS embed.FS

var bundle *goi18n.Bundle

func Init() {
	bundle = goi18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	files, err := localeFS.ReadDir("locales")
	if err != nil {
		log.Fatalf("i18n: read locales: %v", err)
	}
	for _, f := range files {
		if _, err := bundle.LoadMessageFileFS(localeFS, "locales/"+f.Name()); err != nil {
			log.Fatalf("i18n: load %s: %v", f.Name(), err)
		}
	}
}

// T translates messageID for the given locale, falling back to English.
func T(locale, messageID string) string {
	if bundle == nil {
		return messageID
	}
	localizer := goi18n.NewLocalizer(bundle, locale, "en")
	msg, err := localizer.Localize(&goi18n.LocalizeConfig{MessageID: messageID})
	if err != nil {
		return messageID
	}
	return msg
}
