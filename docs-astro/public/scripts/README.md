# Installation Scripts

These scripts are served directly as static files at:
- `chain.sh/promptext/scripts/install.sh`
- `chain.sh/promptext/scripts/uninstall.sh`
- `chain.sh/promptext/scripts/install.ps1`

## Keeping Scripts in Sync

**Important:** These scripts are copies from `/scripts` in the main repo.

When updating installation scripts:
1. Update the source scripts in `/scripts/`
2. Copy them here: `cp scripts/*.{sh,ps1} docs-astro/public/scripts/`
3. Commit both locations together

## Publishing

These files are published directly from:
- `public/scripts/install.sh`
- `public/scripts/uninstall.sh`
- `public/scripts/install.ps1`
