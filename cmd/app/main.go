package main

import (
	"log"
	"regexp"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/zd4r/blur-collection-holders/internal/clients/blur"
	"github.com/zd4r/blur-collection-holders/internal/store/address"
)

var (
	ethereumAddressPattern = `0x[a-fA-F0-9]{40}`
)

type ownership struct {
	OwnerAddress string
	NumberOwned  int
}

func main() {
	a := app.New()
	w := a.NewWindow("blur collection ownerships")
	w.SetMaster()
	w.Resize(fyne.NewSize(500, 0))

	blurClient, err := blur.NewClient()
	if err != nil {
		log.Fatalf("failed to create blur client: %s", err.Error())
	}

	addressStore := address.NewMap()

	// collection input
	var (
		ownershipsFiltered []ownership
	)

	collectionAddress := widget.NewEntry()
	collectionAddress.SetPlaceHolder("collection address")

	showHoldersTableButton := widget.NewButton(
		"show",
		func() {
			collectionName, err := blurClient.GetCollectionNameByAddress(collectionAddress.Text)
			if err != nil {
				dialog.ShowInformation(
					"error occurred",
					err.Error(),
					w,
				)
				return
			}

			tableWindow := a.NewWindow(collectionName)
			tableWindow.Resize(fyne.NewSize(500, 250))

			ownerships, err := blurClient.GetCollectionOwnerships(collectionAddress.Text)
			if err != nil {
				dialog.ShowInformation(
					"error occurred",
					err.Error(),
					w,
				)
				return
			}
			ownershipsFiltered = filterOwnerships(ownerships, addressStore)

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
							ownershipsFiltered = filterOwnerships(ownerships, addressStore)
						},
					),
					nil, nil, nil,
					widget.NewList(
						func() int {
							return len(ownershipsFiltered)
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
							l1.SetText(ownershipsFiltered[id].OwnerAddress)

							l2 := c.Objects[1].(*widget.Label)
							l2.SetText(strconv.Itoa(ownershipsFiltered[id].NumberOwned))
						}),
				),
			)

			tableWindow.Show()
		},
	)

	// tracked addresses input
	ethereumAddressRexExp := regexp.MustCompile(ethereumAddressPattern)

	multiLineEntry := widget.NewMultiLineEntry()
	multiLineEntry.SetPlaceHolder("tracked addresses")

	multiLineEntry.OnChanged = func(s string) {
		matches := ethereumAddressRexExp.FindAllString(multiLineEntry.Text, -1)
		for _, addr := range matches {
			addressStore.Set(addr)
		}
	}

	// page layout
	pageLayout := container.NewBorder(
		container.New(
			layout.NewGridLayout(2),
			collectionAddress,
			showHoldersTableButton,
		), nil, nil, nil, multiLineEntry,
	)

	w.SetContent(pageLayout)

	w.ShowAndRun()
}

func filterOwnerships(in blur.CollectionOwnerships, addressStore *address.Map) []ownership {
	out := make([]ownership, 0, len(in.Ownerships))

	for _, own := range in.Ownerships {
		if addressStore.Exists(own.OwnerAddress) {
			out = append(out, ownership{OwnerAddress: own.OwnerAddress, NumberOwned: own.NumberOwned})
		}
	}

	return out
}
