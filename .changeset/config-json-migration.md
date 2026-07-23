---
"mdl": major
---

Config format changed from HCL to JSON. Config now lives in an OS-appropriate global directory (e.g. `~/.config/mdl/config.json` on Linux) instead of a local `config.hcl`; old `config.hcl` files are not auto-imported, and `example.config.hcl` has been removed in favor of an auto-generated default config plus a JSON Schema (`schema/config.schema.json`).

Dropped the `koanf`+HCL dependencies. Config is now split into `FileConfig` + `RuntimeFlags`. `--config` flag default is now empty (resolves to the OS default) instead of `config.hcl`.

Added:
- `$schema` field in config for editor autocomplete/validation
- `version` field with an automatic migration framework for future config schema changes
- `MDL_CONFIG` env var support for overriding the config path
- `check_updates` config option to check GitHub for new releases

Also fixed the default config path on Windows, which used `%AppData%/Roaming` instead of `~/.config/mdl`.
