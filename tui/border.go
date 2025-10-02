package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	highlight = lipgloss.AdaptiveColor{
		Light: "#FF4757",
		Dark:  "#FF4757",
	}
	title = lipgloss.AdaptiveColor{
		Light: "#FF4757",
		Dark:  "#FF4757",
	}
	option = lipgloss.AdaptiveColor{
		Light: "#A4B0BE",
		Dark:  "#A4B0BE",
	}
	bg = lipgloss.AdaptiveColor{
		Light: "#222222",
		Dark:  "#222222",
	}
	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(title)
	paneBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "╰",
		BottomRight: "╯",
	}
	moonSymbols = lipgloss.Border{
		Left:  "☾",
		Right: "☽",
	}
	pane = lipgloss.NewStyle().
		Border(paneBorder, true).BorderForeground(highlight).Padding(0, 2).BorderBackground(bg)
)

type borderPosition int

const (
	TopLeftBorder borderPosition = iota
	TopRightBorder
	BottomLeftBorder
	BottomRightBorder
)

func borderize(content string, embeddedText map[borderPosition]string) string {
	width := lipgloss.Width(content)
	style := lipgloss.NewStyle().Foreground(highlight).Background(bg)

	if embeddedText == nil {
		embeddedText = make(map[borderPosition]string)
	}

	writeBorderedText := func(text string) string {
		if text != "" {
			return fmt.Sprintf("%s%s%s", style.Render(moonSymbols.Left), style.Render(text), style.Render(moonSymbols.Right))
		}
		return text
	}

	buildHorizontalBorder := func(leftText, rightText string) string {
		leftText = writeBorderedText(leftText)
		rightText = writeBorderedText(rightText)

		remainingLength := max(0, width-lipgloss.Width(leftText)-lipgloss.Width(rightText)-2)

		s := leftText + style.Render(strings.Repeat(paneBorder.Top, remainingLength)) + rightText
		s = lipgloss.NewStyle().Inline(true).MaxWidth(width).Render(s)
		return style.Render(paneBorder.TopLeft) + s + style.Render(paneBorder.TopRight)
	}

	return strings.Join([]string{
		buildHorizontalBorder(embeddedText[TopLeftBorder], embeddedText[TopRightBorder]),
		lipgloss.NewStyle().BorderForeground(highlight).Border(paneBorder, false).Render(content),
		buildHorizontalBorder(embeddedText[BottomLeftBorder], embeddedText[BottomLeftBorder]),
	}, "/n")
}
