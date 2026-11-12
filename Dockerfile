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
    if [ "$BUILD_APPIMAGE" != "true" ]; then \
      echo "Skipping AppImage"; \
      exit 0; \
    fi; \
    ICON="build/appicon.png"; \
    if [ ! -f "$ICON" ]; then \
      mkdir -p build; \
      printf '\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR\x00\x00\x00\x01\x00\x00\x00\x01\x08\x02\x00\x00\x00\x90wS\xde\x00\x00\x00\x0cIDATx\x9cc\xf8\x0f\x00\x00\x01\x01\x00\x05\x18\xd8N\x00\x00\x00\x00IEND\xaeB`\x82' > "$ICON"; \
    fi; \
    if [ ! -f build/AppRun ]; then \
      printf '%s\n' '#!/bin/sh' 'HERE="$(dirname "$(readlink -f "$0")")"' 'exec "$HERE/usr/bin/mdmark" "$@"' > build/AppRun; \
      chmod +x build/AppRun; \
    fi; \
    if [ ! -f build/mdmark.desktop ]; then \
      printf '%s\n' \
        '[Desktop Entry]' \
        'Type=Application' \
        'Name=mdmark' \
        'GenericName=Markdown Editor' \
        'Comment=Privacy-first local Markdown editor (desktop SvelteMark)' \
        'Exec=mdmark %F' \
        'Icon=mdmark' \
        'Terminal=false' \
        'Categories=Office;TextEditor;Utility;' \
        'MimeType=text/markdown;text/x-markdown;text/plain;' \
        'StartupWMClass=mdmark' > build/mdmark.desktop; \
    fi; \
    for platform in $(printf '%s' "$PLATFORMS" | tr ',' ' '); do \
      target_os="${platform%/*}"; target_arch="${platform#*/}"; \
      [ "$target_os" = "linux" ] || continue; \
      BIN="/app/build/bin/mdmark-linux-${target_arch}"; \
      if [ ! -f "$BIN" ]; then \
        echo "No binary for AppImage: $BIN"; continue; \
      fi; \
      APPDIR="/tmp/mdmark-${target_arch}.AppDir"; \
      rm -rf "$APPDIR"; \
      mkdir -p "${APPDIR}/usr/bin"; \
      cp "$BIN" "${APPDIR}/usr/bin/mdmark"; \
      chmod +x "${APPDIR}/usr/bin/mdmark"; \
      cp "$ICON" "${APPDIR}/mdmark.png"; \
      ln -sf mdmark.png "${APPDIR}/.DirIcon"; \
      cp build/mdmark.desktop "${APPDIR}/mdmark.desktop"; \
      cp build/AppRun "${APPDIR}/AppRun"; \
      chmod +x "${APPDIR}/AppRun"; \
      export NO_STRIP=1 APPIMAGE_EXTRACT_AND_RUN=1; \
      (cd /tmp && /usr/local/bin/linuxdeploy --appdir "$APPDIR" --output appimage) || \
        (cd /tmp && /usr/local/bin/linuxdeploy --appdir "$APPDIR" --output appimage --appimage-extract-and-run) || true; \
      AI="$(find /tmp -maxdepth 1 -name 'mdmark*.AppImage' -type f | head -1)"; \
      if [ -n "$AI" ]; then \
        mv "$AI" "/app/build/bin/mdmark-linux-${target_arch}.AppImage"; \
        chmod +x "/app/build/bin/mdmark-linux-${target_arch}.AppImage"; \
        echo "AppImage: /app/build/bin/mdmark-linux-${target_arch}.AppImage"; \
      else \
        echo "WARNING: AppImage not produced for $target_arch"; \
      fi; \
      rm -rf "$APPDIR"; \
    done

FROM scratch AS exporter
COPY --from=builder /app/build/bin /build/bin
