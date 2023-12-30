package main

import (
	"fmt"
	"log"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/zd4r/blur-collection-holders/internal/clients/blur"
)

func main() {
	a := app.New()
	w := a.NewWindow("blur collection ownerships")
	w.SetMaster()
	w.Resize(fyne.NewSize(500, 0))

	collectionAddress := widget.NewEntry()
	collectionAddress.SetPlaceHolder("collection address")

	var (
		ownerships blur.CollectionOwnerships
		err        error
	)

	blurClient, err := blur.NewClient()
	if err != nil {
		log.Fatalf("failed to create blur client: %s", err.Error())
	}

	showHoldersTableButton := widget.NewButton(
		"show",
		func() {
			tableWindow := a.NewWindow(fmt.Sprintf("%s holders", collectionAddress.Text))
			tableWindow.Resize(fyne.NewSize(500, 250))

			ownerships, err = blurClient.GetCollectionOwnerships(collectionAddress.Text)
			if err != nil {
				if err != nil {
					dialog.ShowInformation(
						"error occurred",
						err.Error(),
						w,
					)
					return
				}
			}

			tableWindow.SetContent(
				container.NewBorder(
					widget.NewButton(
						"refresh",
						func() {
							ownerships, err = blurClient.GetCollectionOwnerships(collectionAddress.Text)
							if err != nil {
								dialog.ShowInformation(
									"error occurred",
									err.Error(),
									w,
								)
								return
							}
						},
					),
					nil, nil, nil,
					widget.NewList(
						func() int {
							return len(ownerships.Ownerships)
						},
						func() fyne.CanvasObject {
							return container.NewBorder(
								nil, nil,
								widget.NewLabel(""),
								widget.NewLabel(""),
							)
						},
						func(id widget.ListItemID, o fyne.CanvasObject) {
							c := o.(*fyne.Container)
							l1 := c.Objects[0].(*widget.Label)
							l1.SetText(ownerships.Ownerships[id].OwnerAddress)

							l2 := c.Objects[1].(*widget.Label)
							l2.SetText(strconv.Itoa(ownerships.Ownerships[id].NumberOwned))
						}),
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
