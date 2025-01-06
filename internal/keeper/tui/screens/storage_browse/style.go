package storagebrowse

import (
	"gophkeeper/internal/keeper/tui/styles"

	"github.com/charmbracelet/lipgloss"
)

var (
	screenStyle = styles.Regular.PaddingLeft(2)
	tableStyle  = styles.Border.BorderForeground(lipgloss.Color("240"))

	tableSelectedStyle = styles.Regular.
				Foreground(lipgloss.Color("229")).
				Background(lipgloss.Color("57")).
				Bold(false)

	tableHeaderStyle = styles.Padded.
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("240")).
				BorderBottom(true).
				Bold(false)
)
