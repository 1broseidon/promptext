# Homebrew Tap Setup - Fix Permission Issue

## The Problem

The v0.5.2 release **mostly worked**! ✅
- Built all binaries successfully
- Created GitHub release
- Uploaded all assets

But it failed at the last step with:
```
403 Resource not accessible by integration
```

This is because the default `GITHUB_TOKEN` can only push to the `promptext` repository, not to the separate `homebrew-tap` repository.

## The Solution (2 minutes)

### Step 1: Create a Personal Access Token (PAT)

1. Visit: https://github.com/settings/tokens/new
2. **Note**: "GoReleaser Homebrew Tap"
3. **Expiration**: 90 days (or longer if preferred)
4. **Scopes**: Check **ONLY** `repo` (Full control of private repositories)
   - This gives access to push to your `homebrew-tap` repository
5. Click "Generate token"
6. **COPY THE TOKEN** - you won't see it again!

### Step 2: Add Token to Repository Secrets

1. Visit: https://github.com/1broseidon/promptext/settings/secrets/actions
2. Click "New repository secret"
3. **Name**: `HOMEBREW_TAP_TOKEN`
4. **Secret**: Paste the token you copied
5. Click "Add secret"

### Step 3: Push the Updated Workflow

```bash
git add .github/workflows/release.yml
git commit -m "fix: use PAT for Homebrew tap updates"
git push
```

## Test the Fix

After setting up the token, create a new patch release to test:

```bash
git tag -a v0.5.3 -m "Release v0.5.3: Fix Homebrew tap automation"
git push origin v0.5.3
```

Watch at: https://github.com/1broseidon/promptext/actions

This time it should complete successfully and update the Homebrew tap!

## Verify Success

After the workflow completes, check:

1. **GitHub Release**: https://github.com/1broseidon/promptext/releases/tag/v0.5.3
2. **Homebrew Formula**: https://github.com/1broseidon/homebrew-tap/blob/main/Casks/promptext.rb
3. **Test installation**:
   ```bash
   brew tap 1broseidon/tap
   brew install promptext
   prx --version
   ```

## Alternative: Skip Homebrew Tap (Temporary)

If you want to skip the Homebrew tap for now, you can temporarily disable it in `.goreleaser.yaml` by commenting out the `homebrew_casks` section.

Users can still install via:
- Direct download from GitHub releases
- `go install github.com/1broseidon/promptext/cmd/promptext@latest`

## Why This Is Necessary

GitHub's default `GITHUB_TOKEN` has limited permissions for security reasons. It can only write to the repository where the workflow is running. To push to a different repository (`homebrew-tap`), you need a Personal Access Token with broader permissions.

This is a one-time setup - once the token is added as a secret, all future releases will work automatically.
