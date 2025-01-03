// Provides lipgloss styles and render helpers for TUI app
package styles

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	/* Common */

	Regular = lipgloss.NewStyle()
	Bold    = Regular.Bold(true)
	Padded  = Regular.Padding(0, 1)

	Border         = Regular.Border(lipgloss.RoundedBorder())
	ThickBorder    = Regular.Border(lipgloss.ThickBorder()).BorderForeground(lipgloss.AdaptiveColor{Dark: string(Violet), Light: string(Violet)})
	ActiveBorder   = Border.BorderForeground(lipgloss.AdaptiveColor{Dark: string(Violet), Light: string(Violet)})
	InactiveBorder = Border.BorderForeground(lipgloss.AdaptiveColor{Dark: string(White), Light: string(White)})

	Focused = Regular.Foreground(lipgloss.AdaptiveColor{Dark: "205", Light: "205"})
	Blurred = Regular.Foreground(lipgloss.AdaptiveColor{Dark: "240", Light: "240"})

	Highlighted = Regular.Foreground(Purple)

	/* App styles */

	HeaderStyle = Bold.Foreground(lipgloss.Color("#FF79C6"))

	HelpKeyStyle  = Bold.Foreground(lipgloss.AdaptiveColor{Dark: "ff", Light: ""}).Margin(0, 1, 0, 0)
	HelpDescStyle = Regular.Foreground(lipgloss.AdaptiveColor{Dark: "248", Light: "246"})
)
