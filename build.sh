#!/usr/bin/env bash
set -Eeuo pipefail

ROOT="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
cd "$ROOT"

MODE=release
case "$(uname -s)" in
  Linux) OS=linux ;;
  Darwin) OS=darwin ;;
  MINGW*|MSYS*|CYGWIN*) OS=windows ;;
  *) OS=linux ;;
esac
case "$(uname -m)" in
  x86_64|amd64) ARCH=amd64 ;;
  aarch64|arm64) ARCH=arm64 ;;
  *) ARCH=amd64 ;;
esac
RUN_TESTS=false
NO_CACHE=false
VERSION=""
APPIMAGE=true

usage() {
  cat <<'USAGE'
mdmark Docker/BuildKit Builder (desktop only)

Usage: ./build.sh [options]
  -mode MODE       release (default) or debug
  -os OS           linux, windows, darwin, or all (comma-separated)
  -arch ARCH       amd64, arm64, or all (comma-separated)
  -run-tests       Run go test / go vet (best-effort)
  -no-cache        Disable Docker build cache
  -no-appimage     Skip Linux AppImage packaging
  -v, --version V  Application version (default: git tag or dev-SHA)
  -h, --help       Show this help

Examples:
  ./build.sh
  ./build.sh -os linux -arch amd64
  ./build.sh -os linux,windows -arch amd64 -v 1.0.0
  ./build.sh -no-cache -no-appimage

Artifacts are written flat into ./dist/
USAGE
}

while (($#)); do
  case "$1" in
    -mode) MODE="${2:?missing value for -mode}"; shift 2 ;;
    -os) OS="${2:?missing value for -os}"; shift 2 ;;
    -arch) ARCH="${2:?missing value for -arch}"; shift 2 ;;
    -run-tests) RUN_TESTS=true; shift ;;
    -no-cache) NO_CACHE=true; shift ;;
    -no-appimage) APPIMAGE=false; shift ;;
    -v|--version) VERSION="${2:?missing value for version}"; shift 2 ;;
    -h|--help) usage; exit 0 ;;
    *) echo "Unknown option: $1" >&2; usage >&2; exit 2 ;;
  esac
done

[[ "$MODE" == release || "$MODE" == debug ]] || { echo "Invalid -mode: $MODE" >&2; exit 2; }

if [[ -z "$VERSION" ]]; then
  VERSION="$(git describe --tags --exact-match 2>/dev/null || true)"
  if [[ -z "$VERSION" ]]; then
    SHA="$(git rev-parse --short HEAD 2>/dev/null || printf unknown)"
    VERSION="dev-${SHA}"
  fi
fi
BUILD_TIME="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
BUILD_COMMIT="$(git rev-parse HEAD 2>/dev/null || true)"
if [[ -z "$BUILD_COMMIT" || "$BUILD_COMMIT" == "HEAD" ]]; then
  BUILD_COMMIT="unknown"
fi
BUILD_SOURCE="$(git remote get-url origin 2>/dev/null || printf local)"

normalize_list() {
  local value="$1" default_all="$2"
  if [[ "$value" == all ]]; then printf '%s' "$default_all"; else printf '%s' "$value"; fi
}
OS="$(normalize_list "$OS" 'linux,windows,darwin')"
ARCH="$(normalize_list "$ARCH" 'amd64,arm64')"
IFS=',' read -r -a OSS <<< "$OS"
IFS=',' read -r -a ARCHES <<< "$ARCH"
PLATFORMS=""
for os in "${OSS[@]}"; do
  case "$os" in linux|windows|darwin) ;; *) echo "Unsupported OS: $os" >&2; exit 2 ;; esac
  for arch in "${ARCHES[@]}"; do
    case "$arch" in amd64|arm64) ;; *) echo "Unsupported architecture: $arch" >&2; exit 2 ;; esac
    [[ -n "$PLATFORMS" ]] && PLATFORMS+=','
    PLATFORMS+="${os}/${arch}"
  done
done

export DOCKER_BUILDKIT=1
OUT_DIR="$ROOT/.buildout"
rm -rf "$OUT_DIR"
mkdir -p "$OUT_DIR"
DOCKER_CACHE_OPTS=()
[[ "$NO_CACHE" == true ]] && DOCKER_CACHE_OPTS+=(--no-cache)

cleanup() { rm -rf "$OUT_DIR"; }
trap cleanup EXIT

echo "mdmark build: version=$VERSION mode=$MODE platforms=$PLATFORMS tests=$RUN_TESTS appimage=$APPIMAGE"
echo "Building entirely in Docker with BuildKit..."

docker build \
  "${DOCKER_CACHE_OPTS[@]}" \
  --build-arg RUN_TESTS="$RUN_TESTS" \
  --build-arg BUILD_MODE="$MODE" \
  --build-arg APP_VERSION="$VERSION" \
  --build-arg BUILD_TIME="$BUILD_TIME" \
  --build-arg BUILD_COMMIT="$BUILD_COMMIT" \
  --build-arg BUILD_SOURCE="$BUILD_SOURCE" \
  --build-arg PLATFORMS="$PLATFORMS" \
  --build-arg BUILD_APPIMAGE="$APPIMAGE" \
  --target exporter \
  --output "type=local,dest=$OUT_DIR" \
  .

mkdir -p "$ROOT/dist"
while IFS= read -r -d '' file; do
  base="$(basename "$file")"
  dest="$ROOT/dist/$base"
  install -m 0755 "$file" "$dest"
  echo "Built: dist/$base"
done < <(find "$OUT_DIR" -type f -print0)

echo "Done. Artifacts are in ./dist/"
ls -la "$ROOT/dist/" || true
