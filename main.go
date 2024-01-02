package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/zd4r/blur-collection-holders/internal/clients/blur"
	"github.com/zd4r/blur-collection-holders/internal/store/address"
	collectionNameCache "github.com/zd4r/blur-collection-holders/internal/store/collection"
)

var (
	ethereumAddressPattern = `0x[a-fA-F0-9]{40}`
	addressToENS           = make(map[string]string)
)

type collection struct {
	ownerships []ownership
}

type ownership struct {
	OwnerAddress string
	NumberOwned  int
}

func main() {
	a := app.New()
	w := a.NewWindow("blur collection ownerships")
	w.SetMaster()
	w.Resize(fyne.NewSize(750, 500))

	// blur clients
	blurClient, err := blur.NewClient()
	if err != nil {
		log.Fatalf("failed to create blur client: %s", err.Error())
	}

	// address stores
	EOAStore := address.NewMap()
	CAStore := address.NewMap()
	collectionNamesCache := collectionNameCache.NewMap()

	// ethereum address regexp
	ethereumAddressRexExp := regexp.MustCompile(ethereumAddressPattern)

	// main page tree
	collections := make(map[string]collection)

	tree := widget.NewTree(
		func(id widget.TreeNodeID) []widget.TreeNodeID {
			switch {
			case id == "":
				return CAStore.GetAll()
			case CAStore.Exists(id):
				children := make([]widget.TreeNodeID, 0, len(collections[id].ownerships))
				for _, own := range collections[id].ownerships {
					children = append(children, fmt.Sprintf("%s:%d", own.OwnerAddress, own.NumberOwned))
				}

				return children
			}

			return []string{}
		},
		func(id widget.TreeNodeID) bool {
			return id == "" || CAStore.Exists(id)
		},
		func(branch bool) fyne.CanvasObject {
			if branch {
				return widget.NewLabel("")
			}
			return container.NewBorder(
				nil, nil,
				widget.NewHyperlink("", nil),
				widget.NewLabel(""),
			)
		},
		func(id widget.TreeNodeID, branch bool, o fyne.CanvasObject) {
			if branch {
				o.(*widget.Label).SetText(collectionNamesCache.Get(id))
				return
			}
			parts := strings.Split(id, ":")
			addr, num := parts[0], parts[1]

			c := o.(*fyne.Container)
			l1 := c.Objects[0].(*widget.Hyperlink)
			if ens, ok := addressToENS[addr]; ok {
				l1.SetText(ens)
			} else {
				l1.SetText(addr)
			}
			_ = l1.SetURLFromString(fmt.Sprintf("https://blur.io/%s", addr))

			l2 := c.Objects[1].(*widget.Label)
			l2.SetText(num)
		})

	// cookie input
	cookie := widget.NewPasswordEntry()
	cookie.SetPlaceHolder("authToken=?; walletAddress=?;")

	// fetch button
	fetchButton := widget.NewButton("fetch", func() {})
	fetchButton.OnTapped = func() {
		fetchButton.Disable()
		defer fetchButton.Enable()
		defer func() {
			fetchButton.Text = fmt.Sprintf("fetch (%s)", time.Now().UTC().Local().Format("2006-01-02 15:04:05"))
		}()

		leaderboard, err := blurClient.GetLeaderboard(cookie.Text)
		if err != nil {
			dialog.ShowInformation(
				"error occurred",
				err.Error(),
				w,
			)
			return
		}

		EOAStore.Clear()

		for _, trader := range leaderboard.Traders {
			if trader.Username != nil {
				addressToENS[trader.WalletAddress] = *trader.Username
			}
			EOAStore.Set(strings.ToLower(trader.WalletAddress))
		}

		for _, addr := range CAStore.GetAll() {
			ownerships, err := blurClient.GetCollectionOwnerships(addr)
			if err != nil {
				dialog.ShowInformation(
					"error occurred",
					err.Error(),
					w,
				)
				return
			}

			collections[addr] = collection{
				ownerships: filterOwnerships(ownerships, EOAStore),
			}

			tree.Refresh()
		}
	}

	// collections addresses input
	collectionAddresses := widget.NewMultiLineEntry()
	collectionAddresses.SetPlaceHolder("collection addresses")

	collectionAddresses.OnChanged = func(s string) {
		CAStore.Clear()

		matches := ethereumAddressRexExp.FindAllString(collectionAddresses.Text, -1)
		for _, addr := range matches {
			CAStore.Set(strings.ToLower(addr))

			if !collectionNamesCache.Exists(addr) {
				collectionName, err := blurClient.GetCollectionNameByAddress(addr)
				if err != nil {
					dialog.ShowInformation(
						"error occurred",
						err.Error(),
						w,
					)
				}

				collectionNamesCache.Set(addr, collectionName)
			}
		}

		tree.Refresh()
	}

	// proxy input
	proxy := widget.NewEntry()
	proxy.SetPlaceHolder("http://user:pass@host:port")
	proxy.OnChanged = func(s string) {
		if err := blurClient.SetProxy(proxy.Text); err != nil {
			dialog.ShowInformation(
				"error occurred",
				err.Error(),
				w,
			)
			return
		}
	}

	// tabs
	tabs := container.NewAppTabs(
		container.NewTabItem(
			"ownerships",
			container.NewBorder(
				nil,
				fetchButton,
				nil, nil,
				tree,
			),
		),
		container.NewTabItem(
			"settings",
			container.NewBorder(
				nil,
				container.NewVBox(cookie, proxy,
					fetchButton,
				),
				nil, nil,
				collectionAddresses,
			),
		),
	)

	w.SetContent(tabs)

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
