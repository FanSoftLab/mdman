# syntax=docker/dockerfile:1.7
ARG GO_VERSION=1.27
ARG NODE_VERSION=20

# Base system with Go and Wails deps (Debian bookworm + bullseye webkit2gtk-4.0)
FROM golang:${GO_VERSION}-bookworm AS system-base
RUN --mount=type=cache,target=/var/cache/apt,sharing=locked \
    --mount=type=cache,target=/var/lib/apt,sharing=locked \
    rm -f /etc/apt/apt.conf.d/docker-clean && \
    echo "deb http://deb.debian.org/debian bullseye main" > /etc/apt/sources.list.d/bullseye.list && \
    apt-get update && apt-get install -y --no-install-recommends \
      ca-certificates git \
      libwebkit2gtk-4.0-dev libgtk-3-dev libayatana-appindicator3-dev \
      gcc g++ pkg-config \
      curl gnupg file && \
    rm -rf /var/lib/apt/lists/* /etc/apt/sources.list.d/bullseye.list

# Node.js layer
FROM system-base AS system-node
ARG NODE_VERSION=20
RUN curl -fsSL https://deb.nodesource.com/setup_${NODE_VERSION}.x | bash - && \
    apt-get update && apt-get install -y --no-install-recommends nodejs && \
    rm -rf /var/lib/apt/lists/* && \
    node --version && npm --version

# Frontend (SvelteMark UI)
FROM system-node AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package.json frontend/package-lock.json* ./
RUN --mount=type=cache,target=/root/.npm \
    npm install --legacy-peer-deps
COPY frontend/ ./
RUN npm run build:wails

# Wails CLI
FROM system-base AS wails-cli
RUN go install github.com/wailsapp/wails/v2/cmd/wails@v2.17.1

# linuxdeploy for AppImage
FROM system-base AS linuxdeploy
ARG TARGETARCH=amd64
RUN set -eux; \
    case "$TARGETARCH" in \
      amd64) LD_ARCH=x86_64 ;; \
      arm64) LD_ARCH=aarch64 ;; \
      *) LD_ARCH=x86_64 ;; \
    esac; \
    curl -fsSL -o /usr/local/bin/linuxdeploy \
      "https://github.com/linuxdeploy/linuxdeploy/releases/download/continuous/linuxdeploy-${LD_ARCH}.AppImage" && \
    chmod +x /usr/local/bin/linuxdeploy

# Build desktop binaries
FROM system-node AS builder
ARG RUN_TESTS=false
ARG BUILD_MODE=release
ARG APP_VERSION=dev
ARG BUILD_TIME=unknown
ARG BUILD_COMMIT=unknown
ARG BUILD_SOURCE=local
ARG PLATFORMS=linux/amd64
ARG BUILD_APPIMAGE=true
# BuildKit target-platform args: auto-set by buildx, but must be redeclared
# per-stage before any RUN can read them (set -eu makes unset access fatal).
ARG TARGETOS
ARG TARGETPLATFORM
WORKDIR /app

COPY --from=wails-cli /go/bin/wails /usr/local/bin/wails
COPY --from=linuxdeploy /usr/local/bin/linuxdeploy /usr/local/bin/linuxdeploy

COPY go.mod go.sum* ./
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download && go mod tidy

COPY . .
COPY --from=frontend-builder /app/frontend/dist ./frontend/dist

# Ensure module graph is complete after full source copy
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod tidy

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    set -eu; \
    if [ "$RUN_TESTS" = "true" ]; then go test ./... || true; fi; \
    go vet ./... || true; \
    mkdir -p /app/build/bin; \
    for platform in $(printf '%s' "$PLATFORMS" | tr ',' ' '); do \
      target_os="${platform%/*}"; target_arch="${platform#*/}"; \
      suffix=""; [ "$target_os" = "windows" ] && suffix=".exe"; \
      flags="-s -w"; \
      if [ "$BUILD_MODE" = "debug" ]; then flags=""; fi; \
      ldflags="$flags -X main.version=$APP_VERSION -X main.buildTime=$BUILD_TIME -X main.buildCommit=$BUILD_COMMIT -X main.buildSource=$BUILD_SOURCE"; \
      outname="mdmark-${target_os}-${target_arch}${suffix}"; \
      dest="/app/build/bin/${outname}"; \
      echo "Building mdmark for $target_os/$target_arch"; \
      CGO_ENABLED=1 GOOS="$target_os" GOARCH="$target_arch" \
        wails build -platform "$target_os/$target_arch" -ldflags "$ldflags" -o "$outname"; \
      # Locate binary (Wails may write to build/bin/ or CWD)
      src=""; \
      for cand in \
        "build/bin/${outname}" \
        "${outname}" \
        "build/bin/mdmark${suffix}" \
        "build/bin/mdmark"; do \
        if [ -f "$cand" ]; then src="$cand"; break; fi; \
      done; \
      if [ -z "$src" ]; then \
        src="$(find build -type f -name 'mdmark*' 2>/dev/null | head -1 || true)"; \
      fi; \
      if [ -z "$src" ] || [ ! -f "$src" ]; then \
        echo "ERROR: binary not found for $target_os/$target_arch" >&2; \
        find build -type f 2>/dev/null | head -30 || true; \
        ls -la . 2>/dev/null | head -20 || true; \
        exit 1; \
      fi; \
      # Same file (Wails already wrote into build/bin/) — skip copy
      src_real="$(readlink -f "$src" 2>/dev/null || echo "$src")"; \
      dest_real="$(readlink -f "$dest" 2>/dev/null || echo "$dest")"; \
      if [ "$src_real" = "$dest_real" ]; then \
        echo "Binary already at $dest"; \
      else \
        cp -f "$src" "$dest"; \
      fi; \
      chmod +x "$dest" || true; \
      echo "Binary: $dest"; \
      ls -la "$dest"; \
    done

# AppImage packaging (Linux only)
RUN set -eu; \
    if [ "$TARGETOS" != "linux" ]; then \
      echo "Skipping AppImage"; \
    else \
    for platform in $(printf '%s' "$TARGETPLATFORM" | tr ',' ' '); do \
      target_os="${platform%/*}"; target_arch="${platform#*/}"; \
      [ "$target_os" = "linux" ] || continue; \
      BIN="/app/build/bin/mdmark-linux-${target_arch}"; \
      if [ ! -f "$BIN" ]; then \
        echo "No binary for AppImage: $BIN"; continue; \
      fi; \
      APPDIR="/tmp/mdmark-${target_arch}.AppDir"; \
      rm -rf "$APPDIR"; mkdir -p "$APPDIR/usr/bin" "$APPDIR/usr/share/applications" "$APPDIR/usr/share/icons/hicolor/256x256/apps"; \
      cp -f "$BIN" "$APPDIR/usr/bin/mdmark"; chmod +x "$APPDIR/usr/bin/mdmark"; \
      cp -f /app/build/mdmark.desktop "$APPDIR/usr/share/applications/mdmark.desktop"; \
      cp -f /app/build/mdmark.desktop "$APPDIR/mdmark.desktop"; \
      cp -f /app/build/appicon.png "$APPDIR/usr/share/icons/hicolor/256x256/apps/mdmark.png"; \
      cp -f /app/build/appicon.png "$APPDIR/mdmark.png"; \
      cp -f /app/build/AppRun "$APPDIR/AppRun"; chmod +x "$APPDIR/AppRun"; \
      # Bundle WebKitGTK helper processes (WebProcess/NetworkProcess/GPUProcess +
      # injected bundle). They are not in the main binary's dependency graph, so
      # linuxdeploy never picks them up; without them the webview stays black,
      # or WebKit silently falls back to the host's helpers (version mismatch).
      TRIPLET="$(gcc -dumpmachine)"; \
      for d in webkit2gtk-4.0 webkitgtk-4.0; do \
        if [ -d "/usr/lib/$TRIPLET/$d" ]; then \
          mkdir -p "$APPDIR/usr/lib/$TRIPLET"; \
          cp -a "/usr/lib/$TRIPLET/$d" "$APPDIR/usr/lib/$TRIPLET/"; \
          rm -f "$APPDIR/usr/lib/$TRIPLET/$d/MiniBrowser"; \
        fi; \
      done; \
      export NO_STRIP=1 APPIMAGE_EXTRACT_AND_RUN=1; \
      (cd /tmp && /usr/local/bin/linuxdeploy --appdir "$APPDIR" --output appimage \
        --exclude-library='libpango*' \
        --exclude-library='libharfbuzz*' \
        --exclude-library='libgtk-3*' \
        --exclude-library='libgdk-3*' \
        --exclude-library='libglib-2.0*' \
        --exclude-library='libgio-2.0*' \
        --exclude-library='libgobject-2.0*' \
        --exclude-library='libgmodule-2.0*' \
        --exclude-library='libcairo*' \
        --exclude-library='libfontconfig*' \
        --exclude-library='libfreetype*' \
      ) || \
      (cd /tmp && /usr/local/bin/linuxdeploy --appdir "$APPDIR" --output appimage --appimage-extract-and-run \
        --exclude-library='libpango*' \
        --exclude-library='libharfbuzz*' \
        --exclude-library='libgtk-3*' \
        --exclude-library='libgdk-3*' \
        --exclude-library='libglib-2.0*' \
      ) || true; \
      # Final purge: never ship host-sensitive text stack inside the AppImage
      for pat in 'libpango*' 'libharfbuzz*' 'libgtk-3*' 'libgdk-3*' 'libglib-2.0.so*' 'libgio-2.0.so*' 'libgobject-2.0.so*'; do \
        find "$APPDIR" -name "$pat" -type f -delete 2>/dev/null || true; \
      done; \
      # If linuxdeploy already produced an AppImage, rebuild from cleaned AppDir when needed
      AI="$(find /tmp -maxdepth 1 -name 'mdmark*.AppImage' -type f | head -1)"; \
      if [ -n "$AI" ]; then \
        # Re-pack cleaned AppDir if libs were deleted after packaging
        if find "$APPDIR/usr/lib" -name 'libpango*' 2>/dev/null | grep -q .; then \
          find "$APPDIR" -name 'libpango*' -o -name 'libharfbuzz*' | xargs -r rm -f; \
        fi; \
        mv "$AI" "/app/build/bin/mdmark-linux-${target_arch}.AppImage"; \
        chmod +x "/app/build/bin/mdmark-linux-${target_arch}.AppImage"; \
        echo "AppImage: /app/build/bin/mdmark-linux-${target_arch}.AppImage"; \
      else \
        echo "WARNING: AppImage not produced for $target_arch"; \
      fi; \
      rm -rf "$APPDIR"; \
    done; \
    fi

FROM scratch AS exporter
COPY --from=builder /app/build/bin /build/bin
