package main

import (
	"fmt"
	"math"
	"os"
	"strconv"

	tui "github.com/computerdane/flextui"
	"github.com/computerdane/flextui/components"
	"github.com/computerdane/hll-arty-calc/lib"
	"github.com/eiannone/keyboard"
	"github.com/fatih/color"
)

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
	borders.SetColorFunc(color.New(color.FgRed).SprintFunc())
	borders.SetTitle(" Hell Let Loose - Artillery Calculator ")
	borders.SetTitleColorFunc(color.New(color.Bold).Add(color.FgGreen).SprintFunc())
	borders.Inner.SetIsVertical(true)
	tui.Screen.AddChild(borders.Outer)

	angles := tui.NewComponent()
	angles.SetContent("Angles will go here")
	borders.Inner.AddChild(angles)

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

	angleBorders := components.NewBorders()
	angleBorders.Outer.SetLength(19)
	angleBorders.SetTitle(" Angle (mils) ")
	colorBlue := color.New(color.FgBlue).SprintFunc()
	angleBorders.SetColorFunc(colorBlue)
	angleBorders.SetTitleColorFunc(colorBlue)
	angleBorders.Inner.SetColorFunc(colorBlue)
	inputArea.AddChild(angleBorders.Outer)
	inputArea.AddChild(tui.NewComponent())

	teamMenuArea := tui.NewComponent()
	teamMenuArea.SetLength(1)
	borders.Inner.AddChild(teamMenuArea)

	team := lib.Germany

	teamMenu := components.NewMenu([]string{
		" [g] Germany ",
		" [u] USA ",
		" [b] Britain ",
		" [r] Russia ",
	})
	teamMenu.SetIsVertical(false)
	teamMenu.SetColorFunc(color.New(color.FgYellow).SprintFunc())
	teamMenu.SetSelectedColorFunc(color.New(color.Bold).Add(color.FgMagenta).SprintFunc())
	teamMenu.AddSelection(0)
	teamMenuArea.AddChild(tui.NewComponent())
	teamMenuArea.AddChild(teamMenu.Outer)
	teamMenuArea.AddChild(tui.NewComponent())

	footer := tui.NewComponent()
	footer.SetLength(1)
	footer.SetContent(" Press [q] or [Ctrl-c] to quit. Press [Enter] to save a distance. ")
	tui.Screen.AddChild(footer)

	tui.Screen.UpdateLayout()
	tui.Screen.Render()

	tui.ShowCursor()
	input.UpdateCursorPos()

	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to get key: %s", err)
			break
		}

		if char == 'q' || key == keyboard.KeyCtrlC {
			break
		}

		var content string

		if char == 'g' || char == 'u' || char == 'b' || char == 'r' {
			teamMenu.RemoveAllSelections()
			switch char {
			case 'g':
				team = lib.Germany
				teamMenu.AddSelection(0)
			case 'u':
				team = lib.Usa
				teamMenu.AddSelection(1)
			case 'b':
				team = lib.Britain
				teamMenu.AddSelection(2)
			case 'r':
				team = lib.Russia
				teamMenu.AddSelection(3)
			}
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
			dist, err := strconv.ParseFloat(content, 64)
			if err != nil {
				angleBorders.Inner.SetContent("error")
				continue
			}

			angle := lib.GetAngle(&team, dist)

			angleBorders.Inner.SetContent(fmt.Sprintf("%d", int(math.Round(angle))))
		} else {
			angleBorders.Inner.SetContent("")
		}

		go angleBorders.Inner.Render()
	}
}
