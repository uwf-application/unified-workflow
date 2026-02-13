#!/bin/bash

# Unified Workflow SDK Release Script
# This script creates a GitHub release for v1.0.0

set -e

echo "=== Creating Unified Workflow SDK Release v1.0.0 ==="
echo ""

# Configuration
VERSION="v1.0.0"
RELEASE_BRANCH="release/v1.0.0"
TAG_NAME="v1.0.0"
RELEASE_FILE="dist/unified-workflow-sdk-v1.0.0.tar.gz"
RELEASE_NOTES="RELEASE_CHECKLIST.md"

# Check if release file exists
if [ ! -f "$RELEASE_FILE" ]; then
    echo "❌ Release file not found: $RELEASE_FILE"
    echo "   Run ./package_sdk.sh first to create the distribution package"
    exit 1
fi

echo "1. Checking git status..."
if ! git status --porcelain | grep -q "^[^?]"; then
    echo "✅ Git working directory is clean"
else
    echo "⚠️  Git working directory has uncommitted changes"
    echo "   Please commit or stash changes before creating a release"
    exit 1
fi

echo ""

echo "2. Checking current branch..."
CURRENT_BRANCH=$(git branch --show-current)
if [ "$CURRENT_BRANCH" != "$RELEASE_BRANCH" ]; then
    echo "⚠️  Not on release branch ($RELEASE_BRANCH)"
    echo "   Current branch: $CURRENT_BRANCH"
    echo "   Switching to release branch..."
    git checkout "$RELEASE_BRANCH"
fi

echo ""

echo "3. Checking tag..."
if git rev-parse "$TAG_NAME" >/dev/null 2>&1; then
    echo "✅ Tag $TAG_NAME already exists"
else
    echo "❌ Tag $TAG_NAME does not exist"
    echo "   Please create the tag first: git tag -a $TAG_NAME -m 'Release $TAG_NAME'"
    exit 1
fi

echo ""

echo "4. Pushing to remote repository..."
echo "   Pushing branch: $RELEASE_BRANCH"
git push origin "$RELEASE_BRANCH"

echo "   Pushing tag: $TAG_NAME"
git push origin "$TAG_NAME"

echo ""

echo "5. Creating GitHub release..."
echo "   Note: This step requires GitHub CLI (gh) to be installed and authenticated"
echo ""

# Check if GitHub CLI is installed
if command -v gh >/dev/null 2>&1; then
    echo "✅ GitHub CLI is installed"
    
    # Create release using GitHub CLI
    echo "   Creating release $TAG_NAME..."
    
    # Read release notes
    if [ -f "$RELEASE_NOTES" ]; then
        RELEASE_BODY=$(cat "$RELEASE_NOTES" | sed -n '/## 📝 Release Notes/,/## 🔧 Support Information/p' | head -n -1)
    else
        RELEASE_BODY="Release $TAG_NAME - Unified Workflow SDK"
    fi
    
    # Create release
    gh release create "$TAG_NAME" \
        --title "Unified Workflow SDK $TAG_NAME" \
        --notes "$RELEASE_BODY" \
        --target "$RELEASE_BRANCH" \
        "$RELEASE_FILE"
    
    echo "✅ GitHub release created!"
else
    echo "⚠️  GitHub CLI (gh) is not installed"
    echo ""
    echo "Manual steps required:"
    echo "1. Go to: https://github.com/your-org/unified-workflow/releases/new"
    echo "2. Select tag: $TAG_NAME"
    echo "3. Title: Unified Workflow SDK $TAG_NAME"
    echo "4. Description: Copy from $RELEASE_NOTES"
    echo "5. Upload: $RELEASE_FILE"
    echo "6. Publish release"
fi

echo ""

echo "6. Verifying release..."
echo "   Release file: $(ls -lh "$RELEASE_FILE")"
echo "   File size: $(du -h "$RELEASE_FILE" | cut -f1)"
echo "   SHA256 checksum:"
sha256sum "$RELEASE_FILE"

echo ""

echo "7. Creating checksum file..."
CHECKSUM_FILE="dist/checksums.txt"
cd dist
sha256sum "unified-workflow-sdk-v1.0.0.tar.gz" > "checksums.txt"
cd ..
echo "   Created: $CHECKSUM_FILE"
cat "$CHECKSUM_FILE"

echo ""

echo "=== Release Summary ==="
echo ""
echo "✅ Release v1.0.0 created successfully!"
echo ""
echo "Release artifacts:"
echo "  - $RELEASE_FILE"
echo "  - $CHECKSUM_FILE"
echo ""
echo "Git:"
echo "  - Branch: $RELEASE_BRANCH (pushed)"
echo "  - Tag: $TAG_NAME (pushed)"
echo ""
echo "Next steps:"
echo "1. Merge release branch to main:"
echo "   git checkout main"
echo "   git merge $RELEASE_BRANCH"
echo "   git push origin main"
echo ""
echo "2. Update documentation site"
echo "3. Notify stakeholders"
echo "4. Monitor for any issues"
echo ""
echo "Release complete! 🎉"