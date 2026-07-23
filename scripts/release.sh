#!/usr/bin/env bash
set -euo pipefail

BUMP="${1:-auto}"

case "$BUMP" in
patch | minor | major | auto) ;;
*)
	echo "Usage: release.sh [patch|minor|major|auto]"
	echo "  Default: auto (detect from changeset kinds)"
	exit 1
	;;
esac

# 1. Assert clean working tree
if [[ -n "$(git status --porcelain)" ]]; then
	echo "Error: working tree is not clean. Commit or stash changes first."
	exit 1
fi

# 2. Assert on main branch
BRANCH="$(git branch --show-current)"
if [[ "$BRANCH" != "main" ]]; then
	echo "Error: must be on 'main' branch. Current branch: '$BRANCH'"
	exit 1
fi

# 3. Full gate
echo "=== Running build, vet, and tests ==="
go build ./...
go vet ./...
go test ./...

# 4. Run changie batch
echo "=== Running changie batch ($BUMP) ==="
NEW_VERSION="$(changie batch "$BUMP")"
if [[ -z "$NEW_VERSION" ]]; then
	echo "Error: changie batch failed to produce a version"
	exit 1
fi

# Strip 'v' prefix if present for sed substitution
VERSION_NO_V="${NEW_VERSION#v}"

echo "New version: $NEW_VERSION"

# 5. Update version in cmd/cli/main.go
if [[ "$OSTYPE" == "darwin"* ]]; then
	sed -i '' "s/var version = \".*\"/var version = \"$VERSION_NO_V\"/" cmd/cli/main.go
else
	sed -i "s/var version = \".*\"/var version = \"$VERSION_NO_V\"/" cmd/cli/main.go
fi

# 6. Commit, tag, push
git add CHANGELOG.md .changes/ cmd/cli/main.go
git commit -m "Release $NEW_VERSION"
git tag "$NEW_VERSION"
git push && git push --tags

echo "=== Release $NEW_VERSION complete ==="
