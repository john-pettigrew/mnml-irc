package main

import ui "github.com/gizak/termui"

var textInput *ui.Par
var messagesList *ui.List

func initializeView() error {
	err := ui.Init()
	if err != nil {
		return err
	}

	messagesList = ui.NewList()

	messagesList.Height = ((ui.TermHeight() * 9) / 10) - 3
	messagesList.Overflow = "wrap"
	messagesList.Border = false
	messagesList.PaddingTop = 3
	messagesList.PaddingBottom = 0
	messagesList.PaddingLeft = 3
	messagesList.PaddingRight = 3

	textInput = ui.NewPar(inputPrefix + buffer.Contents)
	textInput.Height = ui.TermHeight() - messagesList.Height + 3
	textInput.Border = false
	textInput.PaddingTop = 0
	textInput.PaddingBottom = 0
	textInput.PaddingLeft = 3
	textInput.PaddingRight = 3

	ui.Body.AddRows(
		ui.NewRow(
			ui.NewCol(12, 0, messagesList),
		),
		ui.NewRow(
			ui.NewCol(12, 0, textInput),
		),
	)
	ui.Body.Align()
	ui.Render(ui.Body)

	ui.Handle("/sys/wnd/resize", func(e ui.Event) {

		messagesList.Height = ((ui.TermHeight() * 9) / 10) - 3
		textInput.Height = ui.TermHeight() - messagesList.Height + 3
		ui.Body.Align()
		ui.Render(ui.Body)
	})

	return nil
}

func renderScreen() {

	textInput.Text = inputPrefix + buffer.Contents
	channel := serverList.CurrentChannel()
	channel.SetMessageItems(messagesList.InnerHeight(), messagesList.Width)

	messagesList.Items = channel.VisibleMessages

	ui.Render(ui.Body)
}
