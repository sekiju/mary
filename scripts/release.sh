#!/usr/bin/env bash
set -euo pipefail

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

# 4. Run changeset version (consumes .changeset/*.md, updates CHANGELOG.md, bumps version)
echo "=== Running changeset version ==="
npx changeset version

# 5. Extract new version from package.json
NEW_VERSION="v$(node -p "require('./package.json').version")"
VERSION_NO_V="${NEW_VERSION#v}"

echo "New version: $NEW_VERSION"

# 6. Update version in cmd/cli/main.go
if [[ "$OSTYPE" == "darwin"* ]]; then
	sed -i '' "s/var version = \".*\"/var version = \"$VERSION_NO_V\"/" cmd/cli/main.go
else
	sed -i "s/var version = \".*\"/var version = \"$VERSION_NO_V\"/" cmd/cli/main.go
fi

# 7. Stage all changeset artifacts
git add CHANGELOG.md .changeset/ cmd/cli/main.go

# 8. Commit, tag, push
git commit -m "Release $NEW_VERSION"
npx changeset tag
git push && git push --tags

echo "=== Release $NEW_VERSION complete ==="
