# Contributing to mdman

Thank you for contributing! This document outlines the contribution process.

## Code of Conduct
Be respectful, inclusive, and constructive. Follow the [Contributor Covenant](https://www.contributor-covenant.org/version/2/1/code_of_conduct/).

## How to Contribute

### Reporting Issues
- Search existing issues first.
- Use issue templates.
- Provide: OS, version, steps to reproduce, expected vs actual behavior.

### Suggesting Features
- Open a feature request issue.
- Describe use case and expected behavior.
- Discuss before implementing large features.

### Pull Requests
1. Fork the repository.
2. Create a feature branch: `git checkout -b feat/my-feature`
3. Make changes.
4. Run tests and linters.
5. Update documentation (see below).
6. Submit PR with description and checklist.

## Documentation Requirements
**All user-facing changes require documentation updates.** See `.github/PULL_REQUEST_TEMPLATE.md`.

### Documentation Standards
We follow ISO/IEC/IEEE 26514 and ASD-STE100 (Simplified Technical English adaptation):

- **Task topics**: Goal, Prerequisites, Numbered steps (≤20 words/step), Expected result, Recovery
- **Concept topics**: Purpose, What it is, Key behaviors, Limitations
- **Reference topics**: Exact identifiers, values, defaults, examples
- **Troubleshooting**: Symptom, Conditions, Verified causes, Diagnostics, Actions
- **Style**: Active voice, concrete terms, glossary consistency, UI labels match `en.ts`

### Documentation Structure
```
docs/
├── getting-started/
├── concepts/
├── tasks/
├── reference/
├── troubleshooting/
├── appendix/
│   ├── glossary.md
│   ├── licenses.md
│   └── changelog.md
└── templates/
```

### Glossary
Add new terms to `docs/appendix/glossary.md`. Use exact spelling from `en.ts`.

### UI Strings
All user-facing text in `frontend/src/locales/en.ts` and `uk.ts`. Update both.

## Development Setup

### Prerequisites
- Go 1.27+
- Node.js 20+
- Wails CLI v2.17.1
- Linux: `libwebkit2gtk-4.1-dev`, `libgtk-3-dev`, `libayatana-appindicator3-dev`

### Commands
```bash
# Frontend dev
cd frontend && npm install && cd ..
wails dev

# Frontend build
cd frontend && npm run build && cd ..

# Go build
wails build

# Docker build
./build.sh

# Tests
go test ./...
cd frontend && npm run check
```

## Code Style

### Go
- Follow `gofmt` and `golangci-lint`.
- Run `go vet ./...` before commit.
- Use `golangci-lint` defaults.

### TypeScript/Svelte
- Run `npm run check` in `frontend/`.
- Follow existing patterns in codebase.

### Commits
- Conventional commits: `feat:`, `fix:`, `docs:`, `refactor:`, `test:`, `chore:`
- One logical change per commit.

## Testing
- Add tests for new functionality.
- Run `go test ./...` and frontend type check.
- Manual testing for UI changes.

## Release Process
1. Maintainer tags release: `git tag v1.0.0`
2. CI builds artifacts via `./build.sh -v 1.0.0`
3. GitHub Release created with artifacts.

## Questions?
Open a discussion or issue.