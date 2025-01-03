// Provides lipgloss styles and render helpers for TUI app
package styles

import (
	"github.com/charmbracelet/lipgloss"
	tint "github.com/lrstanley/bubbletint"
)

var (
	BordersHeight = 2

	ErrorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff1744"))

	NoStyle      = lipgloss.NewStyle()
	FocusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	BlurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	NewSecretStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#388E3C"))
	UpdatedSecretStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFEB3B"))

	HelpColor lipgloss.TerminalColor

	// Style for currently selected block - menu or content
	ActiveBlock lipgloss.Style
)

// ============================================================

const (
	Black           = lipgloss.Color("#000000")
	DarkRed         = lipgloss.Color("#FF0000")
	Red             = lipgloss.Color("#FF5353")
	Purple          = lipgloss.Color("135")
	Orange          = lipgloss.Color("214")
	BurntOrange     = lipgloss.Color("214")
	Yellow          = lipgloss.Color("#DBBD70")
	Green           = lipgloss.Color("34")
	Turquoise       = lipgloss.Color("86")
	DarkGreen       = lipgloss.Color("#325451")
	LightGreen      = lipgloss.Color("47")
	GreenBlue       = lipgloss.Color("#00A095")
	DeepBlue        = lipgloss.Color("39")
	LightBlue       = lipgloss.Color("81")
	LightishBlue    = lipgloss.Color("75")
	Blue            = lipgloss.Color("63")
	Violet          = lipgloss.Color("13")
	Grey            = lipgloss.Color("#737373")
	LightGrey       = lipgloss.Color("245")
	LighterGrey     = lipgloss.Color("250")
	EvenLighterGrey = lipgloss.Color("253")
	DarkGrey        = lipgloss.Color("#606362")
	White           = lipgloss.Color("#ffffff")
	OffWhite        = lipgloss.Color("#a8a7a5")
	HotPink         = lipgloss.Color("200")
)

var (
	DebugLogLevel = Blue
	InfoLogLevel  = lipgloss.AdaptiveColor{Dark: string(Turquoise), Light: string(Green)}
	ErrorLogLevel = Red
	WarnLogLevel  = Yellow

	LogRecordAttributeKey = lipgloss.AdaptiveColor{Dark: string(LightGrey), Light: string(LightGrey)}

	HelpKey = lipgloss.AdaptiveColor{
		Dark:  "ff",
		Light: "",
	}
	HelpDesc = lipgloss.AdaptiveColor{
		Dark:  "248",
		Light: "246",
	}

	InactivePreviewBorder = lipgloss.AdaptiveColor{
		Dark:  "244",
		Light: "250",
	}

	CurrentBackground            = Grey
	CurrentForeground            = White
	SelectedBackground           = lipgloss.Color("110")
	SelectedForeground           = Black
	CurrentAndSelectedBackground = lipgloss.Color("117")
	CurrentAndSelectedForeground = Black

	TitleColor = lipgloss.AdaptiveColor{
		Dark:  "",
		Light: "",
	}

	GroupReportBackgroundColor = EvenLighterGrey
	TaskSummaryBackgroundColor = EvenLighterGrey

	ScrollPercentageBackground = lipgloss.AdaptiveColor{
		Dark:  string(DarkGrey),
		Light: string(EvenLighterGrey),
	}
)

var (
	Regular = lipgloss.NewStyle()
	Bold    = Regular.Bold(true)
	Padded  = Regular.Padding(0, 1)

	Border      = Regular.Border(lipgloss.NormalBorder())
	ThickBorder = Regular.Border(lipgloss.ThickBorder()).BorderForeground(Violet)

	ModuleStyle = Regular.Foreground(lipgloss.AdaptiveColor{
		Dark:  string(LightishBlue),
		Light: "27",
	})

	WorkspaceStyle = Regular.Foreground(Purple)
)

func init() {
	tint.NewDefaultRegistry()
	tint.SetTint(tint.TintDracula)
	tint.SetTintID("dracula")

	HelpColor = tint.BrightBlack()

	ActiveBlock = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(tint.Purple())
}
