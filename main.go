package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"github.com/sqweek/dialog"
	"github.com/tidwall/gjson"
	"golang.design/x/clipboard"
)

var APIKEY string

const (
	url     = "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent?key="
	payload = `{"contents": [{"parts": [{"text": "%s : %s"}]}]}`
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Quigo")
	myWindow.Resize(fyne.NewSize(800, 400))
	load(&config)
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}

	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("AI", theme.ComputerIcon(), mainTab(myWindow, myApp)),
		container.NewTabItemWithIcon("Settings", theme.SettingsIcon(), settingTab()),
	)

	myWindow.SetCloseIntercept(
		func() {
			if unstagedChanges {
				ok := dialog.Message("%s", "Changes not saved. Do you still want to Quit?").
					Title("Quit ?").
					YesNo()

				if ok {
					myApp.Quit()
				}

			} else {
				myApp.Quit()
			}
		},
	)

	tabs.SetTabLocation(container.TabLocationTop)
	myWindow.SetContent(tabs)
	myWindow.ShowAndRun()
}

func handle(value string, prompt string) (respond string, err error) {
	req, err := http.NewRequest(
		"POST",
		url+config.Apikey,
		strings.NewReader(fmt.Sprintf(payload, prompt, value)),
	)
	if err != nil {
		return "", err
	}

	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	generatedText := gjson.Get(string(body), "candidates.0.content.parts.0.text").String()
	return generatedText, nil
}
