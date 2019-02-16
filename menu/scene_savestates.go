package menu

import (
	"path/filepath"
	"sort"
	"strings"

	ntf "github.com/libretro/ludo/notifications"
	"github.com/libretro/ludo/savestates"
	"github.com/libretro/ludo/settings"
	"github.com/libretro/ludo/state"
	"github.com/libretro/ludo/utils"
	"github.com/libretro/ludo/video"
)

type screenSavestates struct {
	entry
}

func buildSavestates() Scene {
	var list screenSavestates
	list.label = "Savestates"

	list.children = append(list.children, entry{
		label: "Save State",
		icon:  "savestate",
		callbackOK: func() {
			name := utils.DatedName(state.Global.GamePath)
			vid.TakeScreenshot(name)
			err := savestates.Save(name)
			if err != nil {
				ntf.DisplayAndLog(ntf.Error, "Menu", err.Error())
			} else {
				menu.stack[len(menu.stack)-1] = buildSavestates()
				fastForwardTweens()
				ntf.DisplayAndLog(ntf.Success, "Menu", "State saved.")
			}
		},
	})

	gameName := utils.FileName(state.Global.GamePath)
	paths, _ := filepath.Glob(settings.Current.SavestatesDirectory + "/" + gameName + "@*.state")
	sort.Sort(sort.Reverse(sort.StringSlice(paths)))
	for _, path := range paths {
		path := path
		date := strings.Replace(utils.FileName(path), gameName+"@", "", 1)
		list.children = append(list.children, entry{
			label: "Load " + date,
			icon:  "loadstate",
			path:  path,
			callbackOK: func() {
				err := savestates.Load(path)
				if err != nil {
					ntf.DisplayAndLog(ntf.Error, "Menu", err.Error())
				} else {
					state.Global.MenuActive = false
					ntf.DisplayAndLog(ntf.Success, "Menu", "State loaded.")
				}
			},
		})
	}

	list.segueMount()

	return &list
}

func (s *screenSavestates) Entry() *entry {
	return &s.entry
}

func (s *screenSavestates) segueMount() {
	genericSegueMount(&s.entry)
}

func (s *screenSavestates) segueNext() {
	genericSegueNext(&s.entry)
}

func (s *screenSavestates) segueBack() {
	genericAnimate(&s.entry)
}

func (s *screenSavestates) update(dt float32) {
	genericInput(&s.entry, dt)
}

// Override rendering
func (s *screenSavestates) render() {
	list := &s.entry

	_, h := vid.Window.GetFramebufferSize()

	for i, e := range list.children {
		if e.yp < -0.1 || e.yp > 1.1 {
			continue
		}

		fontOffset := 64 * 0.7 * menu.ratio * 0.3

		color := video.Color{R: 0, G: 0, B: 0, A: e.iconAlpha}
		if state.Global.CoreRunning {
			color = video.Color{R: 1, G: 1, B: 1, A: e.iconAlpha}
		}

		if e.labelAlpha > 0 {
			drawSavestateThumbnail(
				list, i,
				filepath.Join(settings.Current.ScreenshotsDirectory, utils.FileName(e.path)+".png"),
				680*menu.ratio-85*e.scale*menu.ratio,
				float32(h)*e.yp-14*menu.ratio-64*e.scale*menu.ratio+fontOffset,
				170*menu.ratio, 128*menu.ratio,
				e.scale, video.Color{R: 1, G: 1, B: 1, A: e.iconAlpha},
			)
			vid.DrawBorder(
				680*menu.ratio-85*e.scale*menu.ratio,
				float32(h)*e.yp-14*menu.ratio-64*e.scale*menu.ratio+fontOffset,
				170*menu.ratio*e.scale, 128*menu.ratio*e.scale, 0.02/e.scale,
				video.Color{R: color.R, G: color.G, B: color.B, A: e.iconAlpha})
			if i == 0 {
				vid.DrawImage(menu.icons["savestate"],
					680*menu.ratio-64*e.scale*menu.ratio,
					float32(h)*e.yp-14*menu.ratio-64*e.scale*menu.ratio+fontOffset,
					128*menu.ratio, 128*menu.ratio,
					e.scale, video.Color{R: 1, G: 1, B: 1, A: e.iconAlpha})
			}

			vid.Font.SetColor(color.R, color.G, color.B, e.labelAlpha)
			vid.Font.Printf(
				840*menu.ratio,
				float32(h)*e.yp+fontOffset,
				0.6*menu.ratio, e.label)
		}
	}
}

func (s *screenSavestates) drawHintBar() {
	genericDrawHintBar()
}
