package screens

import (
	"gophkeeper/internal/keeper/tui/styles"
	"strings"
)

func RenderContent(header, content string) string {
	var b strings.Builder

	b.WriteString(styles.HeaderStyle.Render(header))
	b.WriteString("\n\n")
	b.WriteString(content)

	return styles.ContentPaddedStyle.Render(b.String())
}
