package content

import (
	"net/url"
	"strings"
)

func normalizeLocale(locale string) string {
	return strings.ToLower(strings.TrimSpace(locale))
}

func isSupportedLocale(locale string) bool {
	switch normalizeLocale(locale) {
	case "en", "th":
		return true
	default:
		return false
	}
}

func isValidURL(value string) bool {
	if strings.TrimSpace(value) == "" {
		return true
	}
	parsed, err := url.ParseRequestURI(value)
	if err != nil {
		return false
	}
	return parsed.Scheme == "http" || parsed.Scheme == "https"
}

func defaultContentBody() ContentBody {
	return ContentBody{
		Technical: []TechnicalItem{},
		Projects:  []ProjectItem{},
		PortfolioInfo: PortfolioInfo{
			OwnerName:    "",
			Title:        "",
			Subtitle:     "",
			About:        "",
			ContactEmail: "",
			ContactPhone: "",
			Location:     "",
			Github:       "",
			Linkedin:     "",
			Instagram:    "",
		},
	}
}
