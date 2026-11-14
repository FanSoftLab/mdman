package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// Bench fixtures live in a temp dir; setup cost is excluded from timing.
// Benchmarks measure the public API only, so they compile and compare
// cleanly before and after optimization changes.

func benchWriteFile(b *testing.B, dir, name string, size int) string {
	b.Helper()
	path := filepath.Join(dir, name)
	buf := make([]byte, size)
	for i := range buf {
		buf[i] = byte(i*7 + i>>8)
	}
	if err := os.WriteFile(path, buf, 0644); err != nil {
		b.Fatal(err)
	}
	return path
}

func benchApp(b *testing.B) *App {
	b.Helper()
	return &App{
		configDir:     b.TempDir(),
		settings:      make(map[string]interface{}),
		recentFiles:   []RecentFile{},
		recentFileMax: 20,
		media:         newMediaCache(32, 64<<20),
	}
}

// BenchmarkResolveMediaURL measures the media -> data URL path.
// warm: same path every iteration (cache hits after optimization).
// cold: rotates across 64 distinct files so a bounded cache always misses.
func BenchmarkResolveMediaURL(b *testing.B) {
	a := benchApp(b)
	dir := b.TempDir()
	png := benchWriteFile(b, dir, "img.png", 1<<20)

	b.Run("warm", func(b *testing.B) {
		for b.Loop() {
			if _, err := a.ResolveMediaURL("", png); err != nil {
				b.Fatal(err)
			}
		}
	})

	paths := make([]string, 64)
	for i := range paths {
		paths[i] = benchWriteFile(b, dir, fmt.Sprintf("cold-%02d.png", i), 1<<20)
	}
	i := 0
	b.Run("cold", func(b *testing.B) {
		for b.Loop() {
			if _, err := a.ResolveMediaURL("", paths[i%len(paths)]); err != nil {
				b.Fatal(err)
			}
			i++
		}
	})
}

// BenchmarkReadFile measures the user-facing open-file path
// (resolve + read + AddRecent incl. recent.json persist).
func BenchmarkReadFile(b *testing.B) {
	a := benchApp(b)
	dir := b.TempDir()
	md := benchWriteFile(b, dir, "doc.md", 2<<20)

	for b.Loop() {
		if _, err := a.ReadFile(md); err != nil {
			b.Fatal(err)
		}
	}
}
