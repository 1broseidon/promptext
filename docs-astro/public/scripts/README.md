# Installation Scripts

These scripts are served via Astro endpoints at:
- `chain.sh/promptext/install` → `install.sh`
- `chain.sh/promptext/uninstall` → `uninstall.sh`
- `chain.sh/promptext/install.ps1` → `install.ps1`

## Keeping Scripts in Sync

**Important:** These scripts are copies from `/scripts` in the main repo.

When updating installation scripts:
1. Update the source scripts in `/scripts/`
2. Copy them here: `cp scripts/*.{sh,ps1} docs-astro/public/scripts/`
3. Commit both locations together

## Endpoints

The Astro endpoints are defined in:
- `src/pages/install.ts`
- `src/pages/uninstall.ts`
- `src/pages/install.ps1.ts`

These endpoints serve the scripts as plain text with proper headers for piping to bash/PowerShell.
