#!/usr/bin/env bash
# tag-release.sh — read the version from package.json, create a git tag,
# and push it to origin to trigger the release workflow.
set -euo pipefail

# ── Read version from package.json ──────────────────────────────────────────

VERSION=$(node -e "process.stdout.write(require('./package.json').version)")
TAG="v${VERSION}"

echo "Release tag: ${TAG}"

# ── Pre-flight checks ────────────────────────────────────────────────────────

# Require a clean working tree for package.json
if ! git diff --quiet -- package.json; then
  echo "Error: package.json has uncommitted changes. Commit or stash first."
  exit 1
fi

# Refuse to re-tag
if git rev-parse "${TAG}" >/dev/null 2>&1; then
  echo "Error: tag ${TAG} already exists locally."
  exit 1
fi
if git ls-remote --tags origin "${TAG}" | grep -q "${TAG}"; then
  echo "Error: tag ${TAG} already exists on origin."
  exit 1
fi

# Require local branch to be in sync with origin
LOCAL=$(git rev-parse HEAD)
REMOTE=$(git rev-parse "@{u}" 2>/dev/null || echo "")
if [ -n "${REMOTE}" ] && [ "${LOCAL}" != "${REMOTE}" ]; then
  echo "Error: local branch is not in sync with origin. Push or pull first."
  exit 1
fi

# ── Tag and push ─────────────────────────────────────────────────────────────

git tag "${TAG}"
git push origin "${TAG}"

echo "Pushed ${TAG} — GitHub Actions release workflow triggered."
