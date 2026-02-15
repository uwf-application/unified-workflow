# Unified Workflow Release CLI

A comprehensive command-line interface for releasing all Unified Workflow components. This tool combines all existing release scripts (`package_sdk.sh`, `release_all_components.sh`, `create_release.sh`) into a single, unified CLI with automatic version management.

## Overview

The `./release` script provides a unified interface for:
- **SDK releases** - Packaging and distributing the Go SDK
- **Component releases** - Building Docker images and binaries for individual services
- **Full stack releases** - Complete release bundles with all components
- **Automatic versioning** - Semantic version incrementing with git tagging
- **GitHub integration** - Automated release creation with artifacts

## Installation

The release CLI is a standalone bash script. Make it executable:

```bash
chmod +x ./release
```

## Quick Start

### Release all components (default)
```bash
./release
```
- Auto-increments minor version (v1.0.0 → v1.1.0)
- Runs tests, builds all components, creates release bundle
- Creates git tag and GitHub release

### Release only the SDK
```bash
./release sdk
```

### Release specific component
```bash
./release registry
./release executor
./release worker
```

### Dry run (preview)
```bash
./release --dry-run
./release sdk --dry-run
```

## Command Reference

### Basic Usage
```
./release [COMMAND] [COMPONENT] [OPTIONS]
```

### Commands
- `release` - Create a release (default command)
- `clean` - Clean all release artifacts

### Components (for release command)
- `sdk` - Unified Workflow Go SDK
- `registry` - Registry Service
- `executor` - Executor Service  
- `worker` - Worker Service
- `all` - All components (default)

### Options (for release command)
| Option | Description | Default |
|--------|-------------|---------|
| `-h, --help` | Show help message | |
| `-v, --version VERSION` | Use specific version | Auto-increment |
| `-i, --increment TYPE` | Increment type: major, minor, patch | minor |
| `-r, --registry REG` | Docker registry | uwf-application |
| `--no-tests` | Skip running tests | |
| `--no-tag` | Skip creating git tag | |
| `--no-github` | Skip creating GitHub release | |
| `--no-push` | Skip pushing Docker images | |
| `--dry-run` | Show plan without executing | |

## Examples

### 1. Standard Release
```bash
# Release all components with auto-incremented minor version
./release
```

### 2. Major Version Release
```bash
# Release all components with major version increment
./release all -i major
```

### 3. Specific Version Release
```bash
# Release registry service with specific version
./release registry -v v2.0.0
```

### 4. SDK Only Release
```bash
# Release only the SDK with patch version increment
./release sdk -i patch
```

### 5. Custom Registry
```bash
# Release to custom Docker registry
./release -r myregistry.azurecr.io
```

### 6. Skip Certain Steps
```bash
# Release without tests and without pushing images
./release --no-tests --no-push
```

### 7. Clean Release Artifacts
```bash
# Clean all release artifacts (dist/, release/, Docker images, local git tags)
./release clean
```

## Version Management

### Version File
The current version is stored in `VERSION` file (e.g., `v1.0.0`).

### Auto-increment Behavior
- When no version is specified (`-v`), the script auto-increments based on `--increment` type
- The `VERSION` file is automatically updated after successful release
- Version format: `vMAJOR.MINOR.PATCH` (semantic versioning)

### Increment Types
- **major**: `v1.0.0` → `v2.0.0` (breaking changes)
- **minor**: `v1.0.0` → `v1.1.0` (new features, backward compatible)
- **patch**: `v1.0.0` → `v1.0.1` (bug fixes, backward compatible)

## Release Process

### For "all" Components
1. **Prerequisites check** - Verify Docker, Go, git, GitHub CLI
2. **Run tests** - Execute Go tests with 60-second timeout (can be skipped with `--no-tests`)
3. **Package SDK** - Create SDK distribution bundle
4. **Build Docker images** - Build registry, executor, worker images
5. **Build binaries** - Create standalone binaries for each service
6. **Create release bundle** - Package docker-compose and deployment guide
7. **Create binaries archive** - Package all binaries
8. **Push Docker images** - Push to registry (can be skipped with `--no-push`)
9. **Create git tag** - Tag the release (can be skipped with `--no-tag`)
10. **Create GitHub release** - Publish to GitHub (can be skipped with `--no-github`)

### For Individual Components
- **SDK**: Steps 1-3, 9-10
- **Registry/Executor/Worker**: Steps 1-2, 4-5, 8-10

## Release Artifacts

### SDK Release
```
dist/unified-workflow-sdk-v1.1.0.tar.gz
dist/unified-workflow-sdk-v1.1.0.tar.gz.sha256
```

### Component Release (e.g., registry)
```
Docker image: uwf-application/uwf-registry:v1.1.0
Binary: dist/binaries-v1.1.0/registry-service
```

### Full Stack Release
```
dist/unified-workflow-sdk-v1.1.0.tar.gz
unified-workflow-stack-v1.1.0.tar.gz
dist/unified-workflow-binaries-v1.1.0.tar.gz
Docker images:
  - uwf-application/uwf-registry:v1.1.0
  - uwf-application/uwf-executor:v1.1.0
  - uwf-application/uwf-worker:v1.1.0
```

## Integration with Existing Scripts

The release CLI integrates and replaces these existing scripts:

| Old Script | New Equivalent |
|------------|----------------|
| `./package_sdk.sh` | `./release sdk` |
| `./release_all_components.sh` | `./release all` |
| `./create_release.sh` | `./release` (with GitHub integration) |

### Migration from Old Scripts
```bash
# Old: Package SDK v1.0.0
./package_sdk.sh v1.0.0

# New: Release SDK with auto-increment
./release sdk

# Old: Release all components v1.0.0
./release_all_components.sh v1.0.0

# New: Release all with specific version
./release all -v v1.0.0
```

## Configuration

### Environment Variables
The script can be configured via environment variables:

```bash
# Custom default registry
export DEFAULT_REGISTRY="myregistry.azurecr.io"

# Skip confirmation prompts (for CI/CD)
export RELEASE_AUTO_CONFIRM="true"
```

### Customizing Release Behavior
Edit the script to modify:
- Default registry (`DEFAULT_REGISTRY`)
- Component list (`COMPONENTS`)
- Release directory (`RELEASE_DIR`)
- Test behavior (`run_tests` function)

## CI/CD Integration

### GitHub Actions Example
```yaml
name: Release
on:
  push:
    tags:
      - 'v*'

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.25'
      
      - name: Set up Docker
        uses: docker/setup-buildx-action@v2
      
      - name: Login to Docker Registry
        uses: docker/login-action@v2
        with:
          registry: ${{ secrets.REGISTRY_URL }}
          username: ${{ secrets.REGISTRY_USERNAME }}
          password: ${{ secrets.REGISTRY_PASSWORD }}
      
      - name: Run Release
        run: |
          chmod +x ./release
          ./release all \
            -r ${{ secrets.REGISTRY_URL }} \
            --no-tests \
            --no-push
        env:
          RELEASE_AUTO_CONFIRM: "true"
```

### Jenkins Pipeline Example
```groovy
pipeline {
    agent any
    
    environment {
        REGISTRY = 'myregistry.azurecr.io'
    }
    
    stages {
        stage('Release') {
            steps {
                sh '''
                    chmod +x ./release
                    ./release all -r ${REGISTRY} --no-tests
                '''
            }
        }
    }
}
```

## Troubleshooting

### Common Issues

1. **Tests hanging or timing out**
   ```
   Tests timed out after 30 seconds (continuing anyway)
   ```
   Solution: Use `--no-tests` flag to skip tests, or install coreutils:
   ```bash
   brew install coreutils  # macOS
   # Then use gtimeout command
   ```

2. **Docker not installed**
   ```
   ❌ Missing required dependencies: docker
   ```
   Solution: Install Docker and ensure Docker daemon is running.

3. **GitHub CLI not installed**
   ```
   ⚠️  GitHub CLI (gh) not installed (skipping GitHub release)
   ```
   Solution: Install GitHub CLI or use `--no-github` flag.

4. **Permission denied**
   ```
   bash: ./release: Permission denied
   ```
   Solution: Run `chmod +x ./release`

5. **Version format error**
   ```
   ❌ Invalid version format: 1.0.0. Expected format: vX.Y.Z
   ```
   Solution: Use `v` prefix: `v1.0.0`

### Debug Mode
For detailed debugging, run with `set -x`:
```bash
bash -x ./release --dry-run
```

## Best Practices

1. **Always use dry-run first**
   ```bash
   ./release --dry-run
   ```

2. **Version control your releases**
   - The script automatically creates git tags
   - Keep `VERSION` file in version control

3. **Test before release**
   - Run `./release --no-tag --no-github --no-push` for local testing
   - Verify artifacts before pushing to registry

4. **Use semantic versioning**
   - `major`: Breaking changes
   - `minor`: New features
   - `patch`: Bug fixes

5. **Document release notes**
   - Update `RELEASE_CHECKLIST.md` before release
   - The script uses this for GitHub release notes

## Extending the Script

### Adding New Components
1. Add component to `COMPONENTS` array
2. Add handling in `release_component` function
3. Add Dockerfile: `Dockerfile.<component>`
4. Add binary build logic in `build_component_binary`

### Custom Release Steps
Override functions in the script:
- `check_prerequisites` - Add custom dependency checks
- `run_tests` - Custom test execution
- `create_release_bundle` - Custom bundle creation

## Support

- **Documentation**: This file and inline script help (`./release --help`)
- **Issues**: Report issues in the project repository
- **Contributing**: Fork and submit pull requests

## License

Part of the Unified Workflow project. See project LICENSE for details.

---

**Last Updated**: February 12, 2026  
**Version**: 1.0.0  
**Author**: Unified Workflow Team