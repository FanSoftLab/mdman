package main

import (
	"context"
	"embed"
	"flag"
	"log"
	"time"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:frontend/dist
var assets embed.FS

var (
	version     = "dev"
	buildTime   = "unknown"
	buildCommit = "unknown"
	buildSource = "local"
)

func main() {
	flag.Parse()
	fileArgs := flag.Args()

	app := NewApp()
	app.SetPendingFiles(fileArgs)

	ws := app.GetWindowState()
	startW, startH := ws.Width, ws.Height
	if startW < 900 {
		startW = 1280
	}
	if startH < 600 {
		startH = 800
	}
	startState := options.Normal
	if ws.Maximised {
		startState = options.Maximised
	}

	// Match the window background to the saved theme to avoid a flash on start.
	bg := &options.RGBA{R: 13, G: 17, B: 23, A: 255} // dark #0d1117
	if app.GetStartupTheme() == "light" {
		bg = &options.RGBA{R: 255, G: 255, B: 255, A: 255}
	}

	err := wails.Run(&options.App{
		Title:            "mdman",
		Width:            startW,
		Height:           startH,
		MinWidth:         900,
		MinHeight:        600,
		WindowStartState: startState,
		BackgroundColour: bg,
		// Keep the window hidden until the frontend is ready, so the GTK
		// default (light) background never flashes before the first frame.
		StartHidden: true,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: app.startup,
		OnDomReady: func(ctx context.Context) {
			app.DomReady(ctx)
			// Wait until GTK has a real GdkWindow before touching geometry,
			// then show the window with the correct size.
			time.AfterFunc(500*time.Millisecond, func() {
				app.RestoreWindow()
				runtime.WindowShow(ctx)
			})
			time.AfterFunc(1500*time.Millisecond, func() { app.RestoreWindow() })
			// Safety: never leave the window hidden if load is slow.
			time.AfterFunc(3*time.Second, func() { runtime.WindowShow(ctx) })
		},
		OnShutdown:    app.Shutdown,
		OnBeforeClose: app.BeforeClose,
		Bind: []interface{}{
			app,
		},
		Linux: &linux.Options{
			// OnDemand avoids black webview on some NVIDIA/Wayland setups
			WebviewGpuPolicy: linux.WebviewGpuPolicyOnDemand,
		},
		SingleInstanceLock: &options.SingleInstanceLock{
			UniqueId: "com.localtools.mdman",
			OnSecondInstanceLaunch: func(secondInstanceData options.SecondInstanceData) {
				ctx := app.GetContext()
				if ctx == nil {
					return
				}
				runtime.WindowUnminimise(ctx)
				runtime.Show(ctx)
				if len(secondInstanceData.Args) > 0 {
					runtime.EventsEmit(ctx, "open-files", secondInstanceData.Args)
				}
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}
