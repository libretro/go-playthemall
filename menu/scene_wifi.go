package menu

import (
	"github.com/libretro/ludo/ludos"
)

type sceneWiFi struct {
	entry
}

func buildWiFi() Scene {
	var list sceneWiFi
	list.label = "WiFi Menu"

	list.children = append(list.children, entry{
		label: "Looking for networks",
		icon:  "reload",
	})

	list.segueMount()

	go func() {
		networks := ludos.ScanNetworks()

		if len(networks) > 0 {
			list.children = []entry{}
			for _, network := range networks {
				network := network
				list.children = append(list.children, entry{
					label:       network,
					icon:        "menu_network",
					stringValue: func() string { return ludos.NetworkStatus(network) },
					callbackOK: func() {
						list.segueNext()
						menu.stack = append(menu.stack, buildKeyboard("Passpharse for "+network, func() {
							go ludos.ConnectNetwork(network)
						}))
					},
				})
				list.segueMount()
				fastForwardTweens()
			}
		} else {
			list.children[0].label = "No network found"
			list.children[0].icon = "menu_close"
		}
	}()

	return &list
}

func (s *sceneWiFi) Entry() *entry {
	return &s.entry
}

func (s *sceneWiFi) segueMount() {
	genericSegueMount(&s.entry)
}

func (s *sceneWiFi) segueNext() {
	genericSegueNext(&s.entry)
}

func (s *sceneWiFi) segueBack() {
	genericAnimate(&s.entry)
}

func (s *sceneWiFi) update(dt float32) {
	genericInput(&s.entry, dt)
}

func (s *sceneWiFi) render() {
	genericRender(&s.entry)
}

func (s *sceneWiFi) drawHintBar() {
	genericDrawHintBar()
}
