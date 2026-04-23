package chat

import "strings"

func normalizeID(id string) string {
	return strings.TrimSpace(id)
}
