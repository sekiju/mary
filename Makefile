.PHONY: release changeset changeset-version changeset-status

changeset:
	npx changeset

changeset-version:
	npx changeset version

changeset-status:
	npx changeset status --since=origin/main

release:
	./scripts/release.sh
