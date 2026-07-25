---
"mdl": major
---

Restructure `config.json`: the `application` and `output` sections are merged into a single `settings` section with camelCase keys (`checkForUpdates`, `autoInstallUpdates`, `maxParallelPageFetches`, `maxParallelChaptersDownload`, `outputDirectory`). Output image format conversion (`png`/`jpeg`/`avif`/`webp`) and `cleanOnStart` are no longer supported — pages are saved as downloaded. Existing config files are not migrated automatically.

If `config.json` fails to parse, `mdl` now shows a TUI prompt offering to overwrite it with defaults instead of exiting silently.
