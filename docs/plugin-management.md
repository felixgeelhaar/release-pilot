# Plugin Management System

## Problem Statement

Current plugin UX is poor:
- Users must manually build plugin binaries
- Manual copy to `~/.release-pilot/plugins/`
- Manual config file editing
- No easy way to discover available plugins
- No update mechanism

**User friction:**
```bash
# Current (BAD UX)
cd plugins/github
go build -o ~/.release-pilot/plugins/github .
# Manually edit release.config.yaml
vim release.config.yaml
```

## Solution: Plugin Management CLI

Inspired by:
- `gh extension install owner/repo`
- `kubectl krew install plugin`
- `brew install package`

### User Experience (Target)

```bash
# Discover available plugins
release-pilot plugin list --available

# Install a plugin
release-pilot plugin install github

# Enable a plugin (adds to config)
release-pilot plugin enable github --owner=felixgeelhaar --repo=release-pilot

# Configure a plugin interactively
release-pilot plugin configure github

# Update plugins
release-pilot plugin update github
release-pilot plugin update --all

# Disable a plugin (keeps binary, removes from config)
release-pilot plugin disable github

# Uninstall a plugin (removes binary and config)
release-pilot plugin uninstall github

# List installed plugins
release-pilot plugin list

# Show plugin info
release-pilot plugin info github
```

## Architecture

### 1. Plugin Registry

**Central registry of official plugins:**

```yaml
# registry.yaml (hosted at github.com/felixgeelhaar/release-pilot-plugins)
plugins:
  - name: github
    description: Create GitHub releases and upload assets
    repository: felixgeelhaar/release-pilot
    path: plugins/github
    version: v1.0.0
    category: vcs
    hooks:
      - pre_publish
      - post_publish
    config_schema:
      owner: {type: string, required: true, description: "Repository owner"}
      repo: {type: string, required: true, description: "Repository name"}
      token: {type: string, env: GITHUB_TOKEN, description: "GitHub token"}
      draft: {type: boolean, default: false, description: "Create as draft"}

  - name: gitlab
    description: Create GitLab releases
    repository: felixgeelhaar/release-pilot
    path: plugins/gitlab
    version: v1.0.0
    category: vcs

  - name: slack
    description: Send release notifications to Slack
    repository: felixgeelhaar/release-pilot
    path: plugins/slack
    version: v1.0.0
    category: notification

  - name: npm
    description: Publish packages to npm registry
    repository: felixgeelhaar/release-pilot
    path: plugins/npm
    version: v1.0.0
    category: package_manager

  - name: homebrew
    description: Publish to Homebrew tap
    repository: felixgeelhaar/release-pilot
    path: plugins/homebrew
    version: v1.0.0
    category: package_manager

  - name: docker
    description: Build and push Docker images
    repository: felixgeelhaar/release-pilot
    path: plugins/docker
    version: v1.0.0
    category: container

  - name: jira
    description: Create and link Jira release versions
    repository: felixgeelhaar/release-pilot
    path: plugins/jira
    version: v1.0.0
    category: project_management
```

### 2. Plugin Manager Service

```go
package plugin

type Manager struct {
    registry   *Registry
    installer  *Installer
    config     *ConfigManager
    pluginDir  string // ~/.release-pilot/plugins
}

// List available plugins from registry
func (m *Manager) ListAvailable(ctx context.Context) ([]PluginInfo, error)

// List installed plugins
func (m *Manager) ListInstalled(ctx context.Context) ([]InstalledPlugin, error)

// Install plugin binary
func (m *Manager) Install(ctx context.Context, name string, version string) error

// Uninstall plugin binary
func (m *Manager) Uninstall(ctx context.Context, name string) error

// Update plugin to latest version
func (m *Manager) Update(ctx context.Context, name string) error

// Enable plugin in config
func (m *Manager) Enable(ctx context.Context, name string, config map[string]any) error

// Disable plugin in config
func (m *Manager) Disable(ctx context.Context, name string) error

// Configure plugin interactively
func (m *Manager) Configure(ctx context.Context, name string) error

// Get plugin info
func (m *Manager) Info(ctx context.Context, name string) (*PluginInfo, error)
```

### 3. Plugin Installer

```go
package plugin

type Installer struct {
    registry  *Registry
    httpClient *http.Client
}

// Download and install plugin binary
func (i *Installer) Install(ctx context.Context, pluginInfo PluginInfo, destDir string) error {
    // 1. Download plugin binary from GitHub releases
    // 2. Verify checksum
    // 3. Extract to plugin directory
    // 4. Make executable
}

// Build plugin from source (fallback if no binary available)
func (i *Installer) BuildFromSource(ctx context.Context, pluginInfo PluginInfo) error {
    // 1. Clone repository
    // 2. Build plugin
    // 3. Move to plugin directory
}
```

### 4. Interactive Configuration

```go
package plugin

type ConfigWizard struct {
    schema ConfigSchema
    ui     *InteractiveUI
}

// Run interactive configuration wizard
func (w *ConfigWizard) Run(ctx context.Context) (map[string]any, error) {
    // Prompt user for each config field based on schema
    // - String inputs
    // - Boolean yes/no
    // - Select from options
    // - Validate inputs
    // - Show defaults and descriptions
}
```

### 5. Plugin Directory Structure

```
~/.release-pilot/
├── plugins/
│   ├── github              # Plugin binary
│   ├── gitlab              # Plugin binary
│   ├── slack               # Plugin binary
│   ├── npm                 # Plugin binary
│   ├── homebrew            # Plugin binary
│   └── manifest.yaml       # Installed plugins metadata
├── config/
│   └── plugins.yaml        # Plugin configurations (alternative to release.config.yaml)
└── cache/
    └── registry.yaml       # Cached plugin registry
```

**manifest.yaml:**
```yaml
installed:
  - name: github
    version: v1.0.0
    installed_at: 2025-12-11T18:00:00Z
    binary_path: /Users/user/.release-pilot/plugins/github
    checksum: sha256:abc123...

  - name: slack
    version: v1.0.0
    installed_at: 2025-12-11T18:01:00Z
    binary_path: /Users/user/.release-pilot/plugins/slack
    checksum: sha256:def456...
```

## CLI Commands

### `release-pilot plugin list`

Lists installed plugins and their status.

```bash
$ release-pilot plugin list

Installed Plugins:
  github (v1.0.0)  ✓ enabled    Create GitHub releases
  slack (v1.0.0)   ✗ disabled   Send Slack notifications
  npm (v0.9.0)     ⚠ update     Publish npm packages

Use 'release-pilot plugin list --available' to see all available plugins.
```

### `release-pilot plugin list --available`

Shows all plugins from the registry.

```bash
$ release-pilot plugin list --available

Available Plugins:

Version Control:
  github    v1.0.0  ✓ installed  Create GitHub releases
  gitlab    v1.0.0               Create GitLab releases

Notifications:
  slack     v1.0.0  ✓ installed  Send Slack notifications
  discord   v1.0.0               Send Discord notifications

Package Managers:
  npm       v1.0.0  ✓ installed  Publish npm packages
  homebrew  v1.0.0               Publish to Homebrew
  docker    v1.0.0               Build and push Docker images

Project Management:
  jira      v1.0.0               Create Jira release versions
  linear    v1.0.0               Create Linear releases

Use 'release-pilot plugin install <name>' to install a plugin.
```

### `release-pilot plugin install <name>`

Downloads and installs a plugin binary.

```bash
$ release-pilot plugin install github

Installing github plugin...
  ✓ Downloading from registry
  ✓ Verifying checksum
  ✓ Installing to ~/.release-pilot/plugins/github
  ✓ Making executable

GitHub plugin installed successfully!

Next steps:
  1. Enable the plugin: release-pilot plugin enable github
  2. Or configure interactively: release-pilot plugin configure github
```

### `release-pilot plugin enable <name>`

Enables a plugin by adding it to the config.

```bash
$ release-pilot plugin enable github --owner=felixgeelhaar --repo=release-pilot

Enabling github plugin...
  ✓ Added to release.config.yaml

Configuration:
  owner: felixgeelhaar
  repo: release-pilot
  token: (from GITHUB_TOKEN env)

GitHub plugin enabled!

Use 'release-pilot plugin configure github' to modify settings.
```

### `release-pilot plugin configure <name>`

Interactive configuration wizard.

```bash
$ release-pilot plugin configure github

Configuring GitHub plugin...

? Repository owner: felixgeelhaar
? Repository name: release-pilot
? GitHub token: (leave empty to use GITHUB_TOKEN env)
? Create releases as drafts? (y/N) n
? Mark releases as prereleases? (y/N) n
? Generate release notes from GitHub? (y/N) n

Configuration saved to release.config.yaml

Test your configuration:
  release-pilot plan --dry-run
```

### `release-pilot plugin info <name>`

Shows detailed plugin information.

```bash
$ release-pilot plugin info github

GitHub Plugin (v1.0.0)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Description:
  Create GitHub releases and upload release assets

Status:        ✓ Installed and enabled
Version:       v1.0.0 (latest)
Category:      Version Control
Repository:    felixgeelhaar/release-pilot
Documentation: https://github.com/felixgeelhaar/release-pilot/blob/main/plugins/github/README.md

Hooks:
  • pre_publish  - Validate GitHub configuration
  • post_publish - Create GitHub release and upload assets

Configuration:
  owner    (string, required)  Repository owner
  repo     (string, required)  Repository name
  token    (string, from env)  GitHub token (uses GITHUB_TOKEN)
  draft    (boolean, default: false)  Create as draft
  prerelease (boolean, default: false)  Mark as prerelease
  assets   (array, optional)   Glob patterns for assets to upload

Example Configuration:
  plugins:
    - name: github
      enabled: true
      config:
        owner: felixgeelhaar
        repo: release-pilot
        draft: false
        assets:
          - "dist/*.tar.gz"
          - "dist/checksums.txt"

Commands:
  Install:     release-pilot plugin install github
  Enable:      release-pilot plugin enable github
  Configure:   release-pilot plugin configure github
  Update:      release-pilot plugin update github
  Disable:     release-pilot plugin disable github
  Uninstall:   release-pilot plugin uninstall github
```

### `release-pilot plugin update <name>`

Updates plugin to the latest version.

```bash
$ release-pilot plugin update github

Checking for updates...
  Current version: v1.0.0
  Latest version:  v1.1.0

Updating github plugin...
  ✓ Downloading v1.1.0
  ✓ Verifying checksum
  ✓ Replacing binary
  ✓ Updated manifest

GitHub plugin updated to v1.1.0!

Changelog:
  - Added support for discussion categories
  - Improved error messages
  - Fixed asset upload race condition
```

### `release-pilot plugin disable <name>`

Disables plugin (keeps binary, removes from config).

```bash
$ release-pilot plugin disable github

Disabling github plugin...
  ✓ Removed from release.config.yaml
  ℹ Binary kept in ~/.release-pilot/plugins/github

GitHub plugin disabled.

To re-enable: release-pilot plugin enable github
To uninstall: release-pilot plugin uninstall github
```

### `release-pilot plugin uninstall <name>`

Completely removes plugin.

```bash
$ release-pilot plugin uninstall github

Uninstalling github plugin...
  ⚠ This will remove the plugin binary and configuration

? Are you sure? (y/N) y

  ✓ Removed from release.config.yaml
  ✓ Removed binary from ~/.release-pilot/plugins/github
  ✓ Updated manifest

GitHub plugin uninstalled.

To reinstall: release-pilot plugin install github
```

## Implementation Plan

### Phase 1: Core Infrastructure (Week 1)
- [ ] Create plugin registry schema
- [ ] Implement Registry service (fetch/cache/parse)
- [ ] Implement Manager service structure
- [ ] Add plugin CLI commands structure

### Phase 2: Installation (Week 2)
- [ ] Implement Installer service
- [ ] Download from GitHub releases
- [ ] Checksum verification
- [ ] Binary extraction and permissions
- [ ] Manifest management

### Phase 3: Configuration (Week 3)
- [ ] Implement ConfigWizard
- [ ] Interactive prompts with validation
- [ ] Schema-driven configuration
- [ ] Config file management

### Phase 4: Enable/Disable (Week 4)
- [ ] Enable plugin (add to config)
- [ ] Disable plugin (remove from config)
- [ ] Update plugin config
- [ ] Validate plugin configuration

### Phase 5: Updates & Info (Week 5)
- [ ] Version checking
- [ ] Plugin updates
- [ ] Changelog display
- [ ] Plugin info command

### Phase 6: Polish & Documentation (Week 6)
- [ ] Beautiful CLI output
- [ ] Error handling and recovery
- [ ] User documentation
- [ ] Video tutorials

## Configuration Integration

**Option 1: Extend release.config.yaml**
```yaml
plugins:
  - name: github
    enabled: true
    path: ~/.release-pilot/plugins/github  # Auto-managed
    config:
      owner: felixgeelhaar
      repo: release-pilot
```

**Option 2: Separate plugins.yaml**
```yaml
# ~/.release-pilot/config/plugins.yaml
github:
  enabled: true
  config:
    owner: felixgeelhaar
    repo: release-pilot

slack:
  enabled: false
  config:
    webhook: ${SLACK_WEBHOOK}
```

**Recommendation:** Option 1 (extend release.config.yaml) for consistency.

## Binary Distribution

### Option 1: GitHub Releases (Recommended)
Release plugin binaries alongside ReleasePilot releases:
```
release-pilot-v1.1.0/
├── release-pilot_Darwin_arm64.tar.gz
├── release-pilot_Linux_amd64.tar.gz
├── plugins/
│   ├── github_Darwin_arm64.tar.gz
│   ├── github_Linux_amd64.tar.gz
│   ├── slack_Darwin_arm64.tar.gz
│   └── ...
└── checksums.txt
```

### Option 2: Separate Repository
Create `release-pilot-plugins` repository with versioned plugin releases.

**Recommendation:** Option 1 for simplicity, Option 2 for independent versioning.

## Security Considerations

1. **Checksum Verification:** Always verify downloaded binaries
2. **Signature Verification:** Optional cosign/GPG signature checking
3. **Registry Security:** Use GitHub as trusted source
4. **Plugin Permissions:** Document what each plugin can access
5. **Config Validation:** Validate user inputs against schema
6. **Secret Handling:** Never log or display secrets

## User Documentation

### Quick Start Guide
```markdown
# Getting Started with Plugins

## 1. Discover Available Plugins
release-pilot plugin list --available

## 2. Install a Plugin
release-pilot plugin install github

## 3. Configure the Plugin
release-pilot plugin configure github

## 4. Test Your Release
release-pilot plan --dry-run

That's it! Your plugin is ready to use.
```

## Success Metrics

- [ ] Users can install a plugin in < 30 seconds
- [ ] No manual file editing required
- [ ] Clear error messages for troubleshooting
- [ ] 90% of users successfully install first plugin
- [ ] Plugin discovery is intuitive

## Future Enhancements

- Third-party plugin support (community plugins)
- Plugin templates for creating custom plugins
- Plugin testing framework
- Plugin marketplace/directory website
- Auto-update notifications
- Plugin dependency resolution
- Plugin sandboxing for security
