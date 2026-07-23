---
"mdl": patch
---

Versioning and release notes now generated via [changesets](https://github.com/changesets/changesets) instead of changie. Added a CI workflow (lint, build, vet, test) on every PR, and a multi-arch release workflow (linux/darwin/windows × amd64/arm64) triggered on `v*` tags. Added `docs/extractors.md`, a walkthrough for writing a new site extractor. Update-check no longer blocks downloads on network failure.
