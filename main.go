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

	err := wails.Run(&options.App{
		Title:            "mdmark",
		Width:            startW,
		Height:           startH,
		MinWidth:         900,
		MinHeight:        600,
		WindowStartState: startState,
		BackgroundColour: &options.RGBA{R: 13, G: 17, B: 23, A: 255},
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: app.startup,
		OnDomReady: func(ctx context.Context) {
			app.DomReady(ctx)
			// Wait until GTK has a real GdkWindow before touching geometry.
			time.AfterFunc(500*time.Millisecond, func() { app.RestoreWindow() })
			time.AfterFunc(1500*time.Millisecond, func() { app.RestoreWindow() })
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
			UniqueId: "com.localtools.mdmark",
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
