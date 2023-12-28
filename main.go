package main

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("blur-collection-holders")
	w.SetMaster()

	collectionAddress := widget.NewEntry()
	collectionAddress.SetPlaceHolder("collection address")

	showHoldersTableButton := widget.NewButton(
		"show",
		func() {
			tableWindow := a.NewWindow(fmt.Sprintf("%s holders", collectionAddress.Text))

			tableWindow.SetContent(
				container.NewBorder(
					widget.NewButton(
						"refresh",
						func() {
							log.Println("refreshing")
						},
					),
					nil, nil, nil,
					getHoldersTableWidget(collectionAddress.Text),
				),
			)

			tableWindow.Show()
		},
	)

	w.SetContent(container.New(
		layout.NewGridLayout(2),
		collectionAddress,
		showHoldersTableButton,
	))

	w.ShowAndRun()
}
