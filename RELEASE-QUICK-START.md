# GoReleaser Quick Start

## Phase 1 Setup Complete! ✅

You now have automated releases with:
- **GitHub Releases** - Pre-built binaries for Linux, macOS, Windows (amd64 & arm64)
- **Homebrew** - Easy installation via `brew install`
- **Checksums** - SHA256 verification for all binaries
- **Automated changelog** - From your git commits

## Before Your First Release

### 1. Create the Homebrew Tap Repository

Visit: https://github.com/new

- **Repository name**: `homebrew-tap`
- **Visibility**: Public
- **Initialize**: Empty (no README, .gitignore, or license)

Click "Create repository"

That's it! GoReleaser will automatically populate it on your first release.

## Creating a Release

### Simple 3-Step Process:

```bash
# 1. Make sure all changes are committed
git add .
git commit -m "feat: your changes here"
git push

# 2. Create and push a version tag
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0

# 3. Watch the magic happen!
# GitHub Actions will automatically:
# - Run tests
# - Build binaries for all platforms
# - Create GitHub release with notes
# - Update Homebrew formula
# - Generate checksums
```

Visit: https://github.com/1broseidon/promptext/actions to watch the release build.

## Testing Locally (Optional)

Before creating a real release, test the build:

```bash
# Build without publishing
goreleaser release --snapshot --clean --skip=publish

# Check the output
ls -lh dist/

# Test the binary
./dist/promptext_linux_amd64_v1/prx --version
```

## After Your First Release

Users can install promptext via:

**Homebrew (macOS/Linux):**
```bash
brew tap 1broseidon/tap
brew install promptext
prx --version
```

**Direct Download:**
Download from: https://github.com/1broseidon/promptext/releases

**Go Install:**
```bash
go install github.com/1broseidon/promptext/cmd/promptext@latest
```

## Version Naming

Use semantic versioning:
- `v1.0.0` - Major release (breaking changes)
- `v1.1.0` - Minor release (new features)
- `v1.1.1` - Patch release (bug fixes)
- `v1.0.0-beta.1` - Pre-release (marked as pre-release automatically)

## Commit Message Tips

For a nice changelog, use conventional commits:

```bash
git commit -m "feat: add new awesome feature"     # Shows in "Features"
git commit -m "fix: correct bug in parser"        # Shows in "Bug Fixes"
git commit -m "docs: update README"               # Shows in "Others"
git commit -m "chore: update dependencies"        # Not in changelog
```

## Troubleshooting

**Error: homebrew-tap repository not found**
- Make sure you created the `homebrew-tap` repository
- Make sure it's public
- Wait a minute and try again

**Error: tests failed**
- Fix the failing tests
- Delete the tag: `git tag -d v1.0.0 && git push origin :refs/tags/v1.0.0`
- Create the tag again after fixes

**Want to test without releasing?**
```bash
goreleaser release --snapshot --clean --skip=publish
```

## Files Created

- `.goreleaser.yaml` - GoReleaser configuration (simple, clean)
- `.github/workflows/release.yml` - GitHub Actions workflow
- `RELEASE.md` - Detailed documentation
- This file - Quick start guide

## What's Next?

Phase 1 is complete! When you're ready for more package managers:

**Phase 2** - Windows package managers (Scoop, Chocolatey)
**Phase 3** - Linux packages (deb, rpm, AUR, Snap)

See `RELEASE.md` for details on expanding to more platforms.

---

**That's it!** Create a tag, push it, and you're done. 🚀
