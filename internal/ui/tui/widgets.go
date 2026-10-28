package tui

import (
	"fmt"
	"image/color"
	"math"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/josephburgess/gust/internal/ui/styles"
)

var sparkChars = []rune("▁▂▃▄▅▆▇█")

// Bar returns a width-cell gauge of v against max.
func Bar(v, max float64, width int) string {
	filled := 0
	if max > 0 {
		filled = clamp(int(math.Round(v/max*float64(width))), 0, width)
	}
	return strings.Repeat("▰", filled) + strings.Repeat("▱", width-filled)
}

// Spark returns one block char per value, scaled to the min and max of vals.
func Spark(vals []float64) string {
	if len(vals) == 0 {
		return ""
	}
	lo, hi := minMax(vals)
	var sb strings.Builder
	for _, v := range vals {
		idx := len(sparkChars) / 2
		if hi > lo {
			idx = int(math.Round((v - lo) / (hi - lo) * float64(len(sparkChars)-1)))
		}
		sb.WriteRune(sparkChars[idx])
	}
	return sb.String()
}

func HBlocks(n int) string {
	return strings.Repeat("█", max(n, 0))
}

// RangeBar places lo..hi on a track that spans weekMin..weekMax. The fill is at least 1 cell.
func RangeBar(lo, hi, weekMin, weekMax float64, width int) (pre, fill, post string) {
	start, end := 0, width
	if span := weekMax - weekMin; span > 0 {
		start = clamp(int(math.Round((lo-weekMin)/span*float64(width))), 0, width)
		end = clamp(int(math.Round((hi-weekMin)/span*float64(width))), 0, width)
	}
	if end-start < 1 {
		if start >= width {
			start = width - 1
		}
		end = start + 1
	}
	return strings.Repeat("·", start), strings.Repeat("━", end-start), strings.Repeat("·", width-end)
}

func minMax(vals []float64) (float64, float64) {
	lo, hi := vals[0], vals[0]
	for _, v := range vals[1:] {
		lo = math.Min(lo, v)
		hi = math.Max(hi, v)
	}
	return lo, hi
}

func clamp(v, lo, hi int) int {
	return min(max(v, lo), hi)
}

var bigGlyphs = map[rune][5]string{
	'0': {"█▀▀▀█", "█   █", "█   █", "█   █", "▀▀▀▀▀"},
	'1': {" ▀█  ", "  █  ", "  █  ", "  █  ", "▀▀▀▀▀"},
	'2': {"▀▀▀▀█", "    █", "█▀▀▀▀", "█    ", "▀▀▀▀▀"},
	'3': {"▀▀▀▀█", "    █", " ▀▀▀█", "    █", "▀▀▀▀▀"},
	'4': {"█   █", "█   █", "▀▀▀▀█", "    █", "    ▀"},
	'5': {"█▀▀▀▀", "█    ", "▀▀▀▀█", "    █", "▀▀▀▀▀"},
	'6': {"█▀▀▀▀", "█    ", "█▀▀▀█", "█   █", "▀▀▀▀▀"},
	'7': {"▀▀▀▀█", "    █", "    █", "    █", "    ▀"},
	'8': {"█▀▀▀█", "█   █", "█▀▀▀█", "█   █", "▀▀▀▀▀"},
	'9': {"█▀▀▀█", "█   █", "▀▀▀▀█", "    █", "▀▀▀▀▀"},
	'-': {"   ", "   ", "▀▀▀", "   ", "   "},
	'°': {"█▀█", "▀▀▀", "   ", "   ", "   "},
}

// BigTemp renders t followed by a degree sign in a 5-row block font.
func BigTemp(t int) []string {
	rows := make([]string, 5)
	for i, r := range fmt.Sprintf("%d°", t) {
		for row := range rows {
			if i > 0 {
				rows[row] += " "
			}
			rows[row] += bigGlyphs[r][row]
		}
	}
	return rows
}

type iconPart struct {
	text  string
	color color.Color
}

const iconWidth = 13

var (
	cloudTop = []iconPart{{"     .-.     ", styles.Subtle}}
	cloudMid = []iconPart{{"    (   ).   ", styles.Subtle}}
	cloudLow = []iconPart{{"   (___(__)  ", styles.Subtle}}

	icons = map[string][][]iconPart{
		"storm": {cloudTop, cloudMid, cloudLow,
			{{"    ‚'", styles.Foam}, {"/_", styles.Gold}, {"‚'", styles.Foam}},
			{{"     ‚' ", styles.Foam}, {"/", styles.Gold}, {" ‚'", styles.Foam}}},
		"drizzle": {cloudTop, cloudMid, cloudLow,
			{{"     ' ' '", styles.Foam}},
			{{"    ' ' '", styles.Foam}}},
		"rain": {cloudTop, cloudMid, cloudLow,
			{{"    ‚'‚'‚'", styles.Foam}},
			{{"    ‚'‚'‚'", styles.Foam}}},
		"snow": {cloudTop, cloudMid, cloudLow,
			{{"    *  *  *", styles.Foam}},
			{{"   *  *  *", styles.Foam}}},
		"fog": {{},
			{{" _ - _ - _ - ", styles.Subtle}},
			{{"  _ - _ - _  ", styles.Subtle}},
			{{" _ - _ - _ - ", styles.Subtle}},
			{}},
		"clear": {
			{{"    \\   /    ", styles.Gold}},
			{{"     .-.     ", styles.Gold}},
			{{"  ― (   ) ―  ", styles.Gold}},
			{{"     `-'     ", styles.Gold}},
			{{"    /   \\    ", styles.Gold}}},
		"clear night": {
			{{"     _.-.  ", styles.Gold}, {"* ", styles.Subtle}},
			{{"   ,'  /     ", styles.Gold}},
			{{"  (   (    ", styles.Gold}, {"* ", styles.Subtle}},
			{{"   `.  \\     ", styles.Gold}},
			{{"     `-'  ", styles.Gold}, {".  ", styles.Subtle}}},
		"cloudy": {{},
			{{"     .--.    ", styles.Subtle}},
			{{"  .-(    ).  ", styles.Subtle}},
			{{" (___.__)__) ", styles.Subtle}},
			{}},
	}
)

// Icon returns 5 styled lines of weather art, each iconWidth cells wide.
// It uses the same condition ID ranges as models.GetWeatherEmoji.
func Icon(id int, night bool) []string {
	var name string
	switch {
	case id >= 200 && id <= 232:
		name = "storm"
	case id >= 300 && id <= 321:
		name = "drizzle"
	case id >= 500 && id <= 531:
		name = "rain"
	case id >= 600 && id <= 622:
		name = "snow"
	case id >= 700 && id <= 781:
		name = "fog"
	case id == 800 && night:
		name = "clear night"
	case id == 800:
		name = "clear"
	default:
		name = "cloudy"
	}

	lines := make([]string, len(icons[name]))
	for i, parts := range icons[name] {
		var sb strings.Builder
		for _, p := range parts {
			sb.WriteString(lipgloss.NewStyle().Foreground(p.color).Render(p.text))
		}
		line := sb.String()
		lines[i] = line + strings.Repeat(" ", max(iconWidth-ansi.StringWidth(line), 0))
	}
	return lines
}

// Panel boxes body under an upper-case title. w is the full width, border included.
func Panel(title, body string, w int, focused bool) string {
	accent, titleColor := styles.HighlightMed, styles.Subtle
	if focused {
		accent, titleColor = styles.Iris, styles.Iris
	}
	content := body
	if title != "" {
		content = lipgloss.NewStyle().Foreground(titleColor).Render(strings.ToUpper(title)) + "\n" + body
	}
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(accent).
		Padding(0, 1).
		Width(w).
		Render(content)
}

func Wordmark() []string {
	return []string{
		"█▀▀▀ █  █ █▀▀▀ ▀█▀",
		"█ ▀█ █  █ ▀▀▀█  █ ",
		"▀▀▀▀ ▀▀▀▀ ▀▀▀▀  ▀ ",
	}
}

// StepRail renders the setup progress, e.g. "✓ city ─── ● preferences ─── ○ sign in".
func StepRail(steps []string, current int) string {
	connector := lipgloss.NewStyle().Foreground(styles.HighlightMed).Render(" ─── ")
	parts := make([]string, len(steps))
	for i, step := range steps {
		switch {
		case i < current:
			parts[i] = lipgloss.NewStyle().Foreground(styles.Foam).Render("✓") + " " +
				lipgloss.NewStyle().Foreground(styles.Subtle).Render(step)
		case i == current:
			parts[i] = lipgloss.NewStyle().Foreground(styles.Rose).Bold(true).Render("● " + step)
		default:
			parts[i] = lipgloss.NewStyle().Foreground(styles.Subtle).Render("○ " + step)
		}
	}
	return strings.Join(parts, connector)
}
