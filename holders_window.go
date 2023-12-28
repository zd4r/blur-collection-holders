package main

import (
	"log"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func getHoldersTableWidget(collectionAddress string) *widget.List {
	//var data [][]string

	ownerships, err := getCollectionHolders(collectionAddress)
	if err != nil {
		log.Fatal(err) // TODO: fix
	}

	//for idx, ownership := range ownerships.Ownerships {
	//	data = append(data, make([]string, 2))
	//	data[idx][0], data[idx][1] = ownership.OwnerAddress, strconv.Itoa(ownership.NumberOwned)
	//}

	table := widget.NewList(
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
		})

	return table
}
