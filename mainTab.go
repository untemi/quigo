package main

import (
	"fmt"
	"log"
	"unicode/utf8"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/sqweek/dialog"
	"golang.design/x/clipboard"
	"golang.org/x/exp/maps"
)

var combo *widget.Select

func mainTab(window fyne.Window, app fyne.App) *fyne.Container {
	confirmBtn := widget.NewButtonWithIcon("", theme.MailSendIcon(), nil)
	confirmBtn.Disable()

	combo = widget.NewSelect(maps.Keys(config.Prompts), nil)

	input := widget.NewEntry()
	input.SetPlaceHolder("Enter text...")
	input.MultiLine = true
	input.Wrapping = fyne.TextWrapWord

	aitext := widget.NewRichText()
	aitext.ParseMarkdown(`# Quigo
  1. Choose your option
  2. Write your prompt
  3. Go !!`)
	aitext.Wrapping = fyne.TextWrapWord
	aitext.Scroll = container.ScrollVerticalOnly

	loading := widget.NewProgressBarInfinite()
	loading.Hidden = true

	copyText := widget.NewButtonWithIcon("", theme.ContentCopyIcon(), func() {
		clipboard.Write(clipboard.FmtText, []byte(aitext.String()))
	})
	copyText.Disable()

	btnValide := func(_ string) {
		if combo.Selected != "" && utf8.RuneCountInString(input.Text) > 2 &&
			utf8.RuneCountInString(config.Apikey) > 30 {
			confirmBtn.Enable()
		} else {
			confirmBtn.Disable()
		}
	}
	btnAction := func() {
		loading.Hidden, aitext.Hidden = false, true
		input.Disable()
		confirmBtn.Disable()

		value := input.Text
		respond, err := handle(value, config.Prompts[combo.Selected].Text)

		if err != nil {
			log.Println("ERROR : ", err)

			if window.Canvas().Focused() != nil {
				dialog.Message("%s", err).Title("Error").Error()
			} else {
				app.SendNotification(&fyne.Notification{Content: fmt.Sprintf("%s", err), Title: "Quigo"})
			}
		} else {
			aitext.ParseMarkdown(respond)
			copyText.Enable()
		}

		loading.Hidden, aitext.Hidden = true, false
		input.Enable()
		confirmBtn.Enable()
	}

	combo.OnChanged, input.OnChanged = btnValide, btnValide
	confirmBtn.OnTapped = func() { go btnAction() }

	main := container.NewBorder(
		container.NewVBox(
			widget.NewSeparator(),
			combo,
			container.NewBorder(nil, nil, nil, confirmBtn, input),
		),
		copyText,
		nil,
		nil,
		container.NewGridWithColumns(1, loading, aitext),
	)

	return main
}
