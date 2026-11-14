// Package main is the entry point for mdman, a privacy-first desktop Markdown
// editor built with Wails v2 (Go + system webview).
//
// The Go backend (app.go) provides:
//   - Workspace and file management
//   - Settings and window state persistence
//   - Media file resolution and caching
//   - Single-instance enforcement
//
// The frontend (Svelte 5 + CodeMirror 6) handles all UI, editing, and preview.
//
// Build variables (injected via -ldflags at build time):
//   - version: semantic version string (e.g., "1.0.0")
//   - buildTime: ISO 8601 UTC timestamp
//   - buildCommit: full git commit SHA
//   - buildSource: repository URL or "local"
//
// Build with: `wails build` or `./build.sh` (cross-platform Docker builds).
// Run dev: `wails dev` (live reload).
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

// Build metadata injected at link time via -ldflags (see build.sh and Dockerfile).
// Default values are used when running via `go run` or `wails dev` without ldflags.
var (
	version     = "dev"     // Semantic version (e.g., "1.0.0")
	buildTime   = "unknown" // ISO 8601 UTC build timestamp
	buildCommit = "unknown" // Full git commit SHA
	buildSource = "local"   // Repository URL or "local"
)

// main initializes the Wails application, configures the window, and starts
// the event loop. It handles:
//
//  1. CLI file arguments (passed to App.SetPendingFiles for frontend consumption)
//  2. Window geometry restoration from persisted state (with sensible defaults)
//  3. Startup theme background color to prevent flash before webview loads
//  4. Wails lifecycle hooks (OnStartup, OnDomReady, OnShutdown, OnBeforeClose)
//  5. Single-instance lock with second-instance file opening
//  6. Linux-specific WebView GPU policy (OnDemand) for NVIDIA/Wayland compatibility
//
// The window is initially hidden (StartHidden) and shown after the frontend
// signals readiness via DomReady, with multiple fallback timers to ensure
// the window never stays hidden indefinitely.
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
			// then show the window with the correct size. The runtime calls
			// may panic if the window is gone (fast close) — never take down
			// the process from a timer callback.
			safe := func(fn func()) {
				defer func() {
					if r := recover(); r != nil {
						log.Printf("window timer recovered: %v", r)
					}
				}()
				fn()
			}
			time.AfterFunc(500*time.Millisecond, func() {
				safe(func() {
					app.RestoreWindow()
					runtime.WindowShow(ctx)
				})
			})
			time.AfterFunc(1500*time.Millisecond, func() {
				safe(app.RestoreWindow)
			})
			// Safety: never leave the window hidden if load is slow.
			time.AfterFunc(3*time.Second, func() {
				safe(func() { runtime.WindowShow(ctx) })
			})
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
			UniqueId: "com.fansoftlab.mdman",
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
