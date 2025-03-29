package main

import (
	"fmt"
	"math"
	"os"
	"strconv"

	"github.com/computerdane/color"
	tui "github.com/computerdane/flextui"
	"github.com/computerdane/flextui/components"
	"github.com/computerdane/hll-arty-calc/lib"
	"github.com/eiannone/keyboard"
)

type Theme struct {
	Angle       *color.Color
	Dist        *color.Color
	Offset      *color.Color
	CountryName *color.Color
}

type CountryTheme struct {
	Primary   *color.Color
	Secondary *color.Color
}

var Theme_Default = Theme{
	Angle:       color.RGB(255, 165, 0),
	Dist:        color.RGB(205, 186, 155),
	Offset:      color.RGB(155, 136, 105),
	CountryName: color.RGB(241, 234, 216),
}

var CountryTheme_Germany = CountryTheme{
	Primary:   color.RGB(255, 60, 60),
	Secondary: color.RGB(142, 152, 157),
}

var CountryTheme_Usa = CountryTheme{
	Primary:   color.RGB(110, 149, 255),
	Secondary: color.RGB(255, 60, 60),
}

var CountryTheme_Britain = CountryTheme{
	Primary:   color.RGB(255, 60, 60),
	Secondary: color.RGB(110, 149, 225),
}

var CountryTheme_Russia = CountryTheme{
	Primary:   color.RGB(255, 217, 0),
	Secondary: color.RGB(255, 60, 60),
}

var theme = &Theme_Default
var countryTheme = &CountryTheme_Germany
var team = &lib.Germany
var dist float64

var savedDistances = make(map[float64]struct{})

func angleCell(dist float64, offset float64) *tui.Component {
	angle := lib.GetAngle(team, dist+offset)

	cell := tui.NewComponent()
	cell.SetIsVertical(true)

	cellRows := tui.NewComponent()
	cellRows.SetIsVertical(true)
	cellRows.SetLength(3)
	cell.AddChild(cellRows)

	angleRow := tui.NewComponent()
	cellRows.AddChild(angleRow)

	angleContent := tui.NewComponent()
	angleContent.SetContent(fmt.Sprintf(" %d ", int(math.Round(angle))))
	angleContent.SetLength(len(*angleContent.Content()))
	angleRow.AddChild(tui.NewComponent())
	angleRow.AddChild(angleContent)
	angleRow.AddChild(tui.NewComponent())

	offsetRow := tui.NewComponent()
	cellRows.AddChild(offsetRow)

	offsetContent := tui.NewComponent()
	if offset > 0 {
		offsetContent.SetContent(fmt.Sprintf(" +%dm ", int(math.Round(offset))))
	} else if offset < 0 {
		offsetContent.SetContent(fmt.Sprintf(" %dm ", int(math.Round(offset))))
	} else {
		offsetContent.SetContent(fmt.Sprintf(" %dm ", int(math.Round(dist))))
	}
	offsetContent.SetLength(len(*offsetContent.Content()))
	offsetRow.AddChild(tui.NewComponent())
	offsetRow.AddChild(offsetContent)
	offsetRow.AddChild(tui.NewComponent())

	cell.SetLength(max(len(*angleContent.Content()), len(*offsetContent.Content())))

	if offset == 0 {
		angleContent.SetColorFunc(theme.Angle.Clone().Add(color.Bold).Sprint)
		offsetContent.SetColorFunc(theme.Dist.Clone().Add(color.Bold).Sprint)
	} else if math.Abs(offset) > 40 {
		angleContent.SetColorFunc(theme.Angle.Clone().Add(color.Faint).Sprint)
		offsetContent.SetColorFunc(theme.Offset.Clone().Add(color.Faint).Sprint)
	} else {
		angleContent.SetColorFunc(theme.Angle.Sprint)
		offsetContent.SetColorFunc(theme.Offset.Sprint)
	}

	return cell
}

func angleRow(dist float64) *tui.Component {
	row := tui.NewComponent()
	row.SetLength(3)

	cells := []*tui.Component{
		angleCell(dist, -60),
		angleCell(dist, -40),
		angleCell(dist, -20),
		angleCell(dist, 0),
		angleCell(dist, 20),
		angleCell(dist, 40),
		angleCell(dist, 60),
	}

	rowContent := tui.NewComponent()
	row.AddChild(tui.NewComponent())
	row.AddChild(rowContent)
	row.AddChild(tui.NewComponent())

	length := 0
	for _, cell := range cells {
		rowContent.AddChild(cell)
		length += cell.Length()
	}

	rowContent.SetLength(length)

	return row
}

func main() {
	keyboard.Open()
	defer keyboard.Close()

	tui.HideCursor()
	defer tui.ShowCursor()
	defer tui.Clear()

	tui.HandleShellSignals()

	tui.Screen.SetIsVertical(true)

	borders := components.NewBorders()
	// borders.SetBorderSymbols(components.BordersSymbols_Double)
	borders.SetTitle(" Arty Calculator ")
	borders.Inner.SetIsVertical(true)
	borders.SetTitleColorFunc(theme.Angle.Clone().Add(color.Bold).Sprint)
	tui.Screen.AddChild(borders.Outer)

	anglesArea := tui.NewComponent()
	anglesArea.SetIsVertical(true)
	borders.Inner.AddChild(anglesArea)

	anglesAreaSpacer := tui.NewComponent()
	anglesAreaSpacer.SetLength(1)
	anglesArea.AddChild(anglesAreaSpacer)

	angles := tui.NewComponent()
	angles.SetIsVertical(true)
	anglesArea.AddChild(angles)

	inputArea := tui.NewComponent()
	inputArea.SetLength(3)
	borders.Inner.AddChild(inputArea)

	inputBorders := components.NewBorders()
	inputBorders.Outer.SetLength(23)
	inputBorders.SetTitle(" Distance (meters) ")
	inputArea.AddChild(tui.NewComponent())
	inputArea.AddChild(inputBorders.Outer)

	input := components.NewInput()
	tui.CursorOwner = input.Outer
	inputBorders.Inner.AddChild(input.Outer)

	inputBorders.SetColorFunc(theme.Dist.Sprint)
	inputBorders.SetTitleColorFunc(theme.Dist.Sprint)
	inputBorders.Inner.SetColorFunc(theme.Dist.Sprint)
	input.SetColorFunc(theme.Dist.Sprint)

	angleBorders := components.NewBorders()
	angleBorders.Outer.SetLength(19)
	angleBorders.SetTitle(" Angle (mils) ")
	angleBorders.SetColorFunc(theme.Angle.Sprint)
	angleBorders.SetTitleColorFunc(theme.Angle.Sprint)
	angleBorders.Inner.SetColorFunc(theme.Angle.Sprint)
	inputArea.AddChild(angleBorders.Outer)
	inputArea.AddChild(tui.NewComponent())

	teamMenuArea := tui.NewComponent()
	teamMenuArea.SetLength(1)
	borders.Inner.AddChild(teamMenuArea)

	teamMenu := components.NewMenu([]string{
		" [g] Germany ",
		" [u] USA ",
		" [b] Britain ",
		" [r] Russia ",
	})
	teamMenu.SetIsVertical(false)
	teamMenu.SetColorFunc(theme.CountryName.Clone().Add(color.Faint).Sprint)
	teamMenu.SetSelectedColorFunc(theme.CountryName.Clone().Add(color.Bold).Sprint)
	teamMenu.AddSelection(0)
	teamMenuArea.AddChild(tui.NewComponent())
	teamMenuArea.AddChild(teamMenu.Outer)
	teamMenuArea.AddChild(tui.NewComponent())

	footer := tui.NewComponent()
	footer.SetLength(1)
	footer.SetContent(" Quit: [q] or [Ctrl-c].  Save a Distance: [Enter].  Clear Distances: [Esc]. ")
	tui.Screen.AddChild(footer)

	applyCountryTheme := func() {
		borders.SetColorFunc(countryTheme.Secondary.Sprint)
		footer.SetColorFunc(countryTheme.Primary.Sprint)
	}
	applyCountryTheme()

	tui.Screen.UpdateLayout()
	tui.Screen.Render()

	tui.ShowCursor()
	input.UpdateCursorPos()

	clearAngles := func() {
		savedDistances = make(map[float64]struct{})
		angles.RemoveAllChildren()
		angles.UpdateLayout()
		go angles.Render()
	}

	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to get key: %s", err)
			break
		}

		if char == 'q' || key == keyboard.KeyCtrlC {
			break
		}

		if key == keyboard.KeyEnter {
			if dist != 0 {
				if _, exists := savedDistances[dist]; !exists {
					savedDistances[dist] = struct{}{}
					angles.AddChild(angleRow(dist))
					angles.UpdateLayout()
					go angles.Render()
				}
				dist = 0
			}

			input.SetContent("")
			input.UpdateCursorPos()
			go input.Outer.Render()

			angleBorders.Inner.SetContent("")
			go angleBorders.Inner.Render()

			continue
		}

		if key == keyboard.KeyEsc {
			clearAngles()

			continue
		}

		var content string

		if char == 'g' || char == 'u' || char == 'b' || char == 'r' {
			clearAngles()
			teamMenu.RemoveAllSelections()
			switch char {
			case 'g':
				countryTheme = &CountryTheme_Germany
				team = &lib.Germany
				teamMenu.AddSelection(0)
			case 'u':
				countryTheme = &CountryTheme_Usa
				team = &lib.Usa
				teamMenu.AddSelection(1)
			case 'b':
				countryTheme = &CountryTheme_Britain
				team = &lib.Britain
				teamMenu.AddSelection(2)
			case 'r':
				countryTheme = &CountryTheme_Russia
				team = &lib.Russia
				teamMenu.AddSelection(3)
			}
			applyCountryTheme()
			go borders.Outer.Render()
			go footer.Render()
			go teamMenu.Outer.Render()
			go teamMenu.RenderChanges()
			goto UpdateAngle
		}

		if '0' <= char && char <= '9' {
			content = input.Content()
			content += string(char)
		} else if key == keyboard.KeyBackspace || key == keyboard.KeyBackspace2 {
			content = input.Content()
			if len(content) > 0 {
				content = content[:len(content)-1]
			} else {
				dist = 0
			}
		} else {
			continue
		}

		input.SetContent(content)
		input.UpdateCursorPos()
		go input.Outer.Render()

	UpdateAngle:
		content = input.Content()
		if len(content) > 0 {
			dist, err = strconv.ParseFloat(content, 64)
			if err != nil {
				angleBorders.Inner.SetContent("error")
				continue
			}

			angle := lib.GetAngle(team, dist)

			angleBorders.Inner.SetContent(fmt.Sprintf("%d", int(math.Round(angle))))
		} else {
			angleBorders.Inner.SetContent("")
		}

		go angleBorders.Inner.Render()
	}
}
