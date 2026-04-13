package main

import (
	"embed"
	"encoding/json"
	"fmt"
	_ "image/jpeg"
	_ "image/png"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed wails.json
var wailsConfig []byte

func main() {
	app := NewApp()

	var cfg struct {
		Version string `json:"version"`
	}
	_ = json.Unmarshal(wailsConfig, &cfg)

	appMenu := menu.NewMenu()
	appMenu.Append(menu.AppMenu()) // Yosegi menu (About, Quit, etc.)
	windowMenu := appMenu.AddSubmenu("Window")
	windowMenu.AddText("Minimize", keys.CmdOrCtrl("m"), func(_ *menu.CallbackData) {})
	windowMenu.AddText("Close", keys.CmdOrCtrl("w"), func(_ *menu.CallbackData) {})

	err := wails.Run(&options.App{
		Title:     "Yosegi",
		Width:     960,
		Height:    700,
		MinWidth:  800,
		MinHeight: 700,
		Menu:  appMenu,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 255, G: 255, B: 255, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
		Mac: &mac.Options{
			TitleBar:             mac.TitleBarHiddenInset(),
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
			About: &mac.AboutInfo{
				Title:   "Yosegi",
				Message: fmt.Sprintf("Version %s\n\nTurn any photo into a mosaic of your memories", cfg.Version),
			},
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
