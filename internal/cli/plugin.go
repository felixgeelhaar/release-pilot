package cli

import (
	"fmt"
	"strings"

	"github.com/felixgeelhaar/release-pilot/internal/plugin/manager"
	"github.com/spf13/cobra"
)

var pluginCmd = &cobra.Command{
	Use:   "plugin",
	Short: "Manage ReleasePilot plugins",
	Long: `Manage plugins for ReleasePilot.

Plugins extend ReleasePilot's functionality for version control systems,
package managers, notification services, and more.

Examples:
  # List available plugins
  release-pilot plugin list --available

  # Install a plugin
  release-pilot plugin install github

  # Configure a plugin interactively
  release-pilot plugin configure github

  # Update a plugin
  release-pilot plugin update github

  # Get plugin information
  release-pilot plugin info github`,
}

var pluginListCmd = &cobra.Command{
	Use:   "list",
	Short: "List plugins",
	Long: `List installed plugins or available plugins from the registry.

By default, shows installed plugins. Use --available to show all plugins
from the registry.`,
	RunE: runPluginList,
}

var pluginInstallCmd = &cobra.Command{
	Use:   "install <name>",
	Short: "Install a plugin",
	Long: `Install a plugin from the registry.

Downloads the plugin binary for your platform and makes it available
for use. Plugins must be enabled after installation with 'plugin enable'.`,
	Args: cobra.ExactArgs(1),
	RunE: runPluginInstall,
}

var pluginUninstallCmd = &cobra.Command{
	Use:     "uninstall <name>",
	Aliases: []string{"remove"},
	Short:   "Uninstall a plugin",
	Long:    `Remove an installed plugin and its associated files.`,
	Args:    cobra.ExactArgs(1),
	RunE:    runPluginUninstall,
}

var pluginEnableCmd = &cobra.Command{
	Use:   "enable <name>",
	Short: "Enable a plugin",
	Long:  `Enable an installed plugin to use it in releases.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runPluginEnable,
}

var pluginDisableCmd = &cobra.Command{
	Use:   "disable <name>",
	Short: "Disable a plugin",
	Long:  `Disable a plugin without uninstalling it.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runPluginDisable,
}

var (
	pluginListAvailable bool
	pluginListRefresh   bool
)

func init() {
	// Add plugin command to root
	rootCmd.AddCommand(pluginCmd)

	// Add subcommands to plugin
	pluginCmd.AddCommand(pluginListCmd)
	pluginCmd.AddCommand(pluginInstallCmd)
	pluginCmd.AddCommand(pluginUninstallCmd)
	pluginCmd.AddCommand(pluginEnableCmd)
	pluginCmd.AddCommand(pluginDisableCmd)

	// Flags for plugin list
	pluginListCmd.Flags().BoolVarP(&pluginListAvailable, "available", "a", false, "Show all available plugins from registry")
	pluginListCmd.Flags().BoolVarP(&pluginListRefresh, "refresh", "r", false, "Force refresh registry cache")
}

func runPluginList(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	mgr, err := manager.NewManager()
	if err != nil {
		return fmt.Errorf("failed to create plugin manager: %w", err)
	}

	var entries []manager.PluginListEntry
	if pluginListAvailable {
		entries, err = mgr.ListAvailable(ctx, pluginListRefresh)
		if err != nil {
			return fmt.Errorf("failed to list available plugins: %w", err)
		}
	} else {
		entries, err = mgr.ListInstalled(ctx)
		if err != nil {
			return fmt.Errorf("failed to list installed plugins: %w", err)
		}
	}

	// Display plugins
	if len(entries) == 0 {
		if pluginListAvailable {
			fmt.Println("No plugins available in registry.")
		} else {
			fmt.Println("No plugins installed.")
			fmt.Println()
			fmt.Println("Use 'release-pilot plugin list --available' to see available plugins.")
			fmt.Println("Use 'release-pilot plugin install <name>' to install a plugin.")
		}
		return nil
	}

	if pluginListAvailable {
		displayAvailablePlugins(entries)
	} else {
		displayInstalledPlugins(entries)
	}

	return nil
}

func displayInstalledPlugins(entries []manager.PluginListEntry) {
	fmt.Println("Installed Plugins:")
	fmt.Println()

	for _, entry := range entries {
		// Status icon
		var statusIcon, statusText string
		switch entry.Status {
		case manager.StatusEnabled:
			statusIcon = "✓"
			statusText = "enabled"
		case manager.StatusInstalled:
			statusIcon = "✗"
			statusText = "disabled"
		case manager.StatusUpdateAvailable:
			statusIcon = "⚠"
			statusText = "update available"
		default:
			statusIcon = " "
			statusText = "unknown"
		}

		version := entry.Info.Version
		if entry.Installed != nil {
			version = entry.Installed.Version
		}

		fmt.Printf("  %s %-15s (%-8s)  %s  %s\n",
			statusIcon,
			entry.Info.Name,
			version,
			formatStatus(statusText),
			entry.Info.Description,
		)
	}

	fmt.Println()
	fmt.Println("Use 'release-pilot plugin info <name>' for more details.")
}

func displayAvailablePlugins(entries []manager.PluginListEntry) {
	fmt.Println("Available Plugins:")
	fmt.Println()

	// Group by category
	categories := make(map[string][]manager.PluginListEntry)
	for _, entry := range entries {
		category := entry.Info.Category
		if category == "" {
			category = "other"
		}
		categories[category] = append(categories[category], entry)
	}

	// Display by category
	categoryNames := []string{"vcs", "notification", "package_manager", "project_management", "container", "other"}
	categoryTitles := map[string]string{
		"vcs":                "Version Control:",
		"notification":       "Notifications:",
		"package_manager":    "Package Managers:",
		"project_management": "Project Management:",
		"container":          "Containers:",
		"other":              "Other:",
	}

	for _, category := range categoryNames {
		plugins, ok := categories[category]
		if !ok || len(plugins) == 0 {
			continue
		}

		fmt.Println(categoryTitles[category])
		for _, entry := range plugins {
			// Status indicator
			var status string
			switch entry.Status {
			case manager.StatusEnabled:
				status = "✓ installed"
			case manager.StatusInstalled:
				status = "✓ installed"
			case manager.StatusUpdateAvailable:
				status = "⚠ update"
			default:
				status = ""
			}

			fmt.Printf("  %-12s %-8s  %-12s %s\n",
				entry.Info.Name,
				entry.Info.Version,
				status,
				entry.Info.Description,
			)
		}
		fmt.Println()
	}

	fmt.Println("Use 'release-pilot plugin install <name>' to install a plugin.")
	fmt.Println("Use 'release-pilot plugin info <name>' for more details.")
}

func formatStatus(status string) string {
	statusColors := map[string]string{
		"enabled":          "enabled",
		"disabled":         "disabled",
		"update available": "update",
	}

	if colored, ok := statusColors[status]; ok {
		return colored
	}
	return status
}

func getCategoryTitle(category string) string {
	switch category {
	case "vcs":
		return "Version Control"
	case "notification":
		return "Notifications"
	case "package_manager":
		return "Package Managers"
	case "project_management":
		return "Project Management"
	case "container":
		return "Containers"
	default:
		return strings.Title(category)
	}
}

func runPluginInstall(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	pluginName := args[0]

	mgr, err := manager.NewManager()
	if err != nil {
		return fmt.Errorf("failed to create plugin manager: %w", err)
	}

	fmt.Printf("Installing plugin %q...\n", pluginName)

	if err := mgr.Install(ctx, pluginName); err != nil {
		return fmt.Errorf("failed to install plugin: %w", err)
	}

	fmt.Println()
	printSuccess(fmt.Sprintf("Plugin %q installed successfully", pluginName))
	fmt.Println()
	fmt.Println("To use this plugin:")
	fmt.Printf("  1. Enable it: release-pilot plugin enable %s\n", pluginName)
	fmt.Printf("  2. Configure it in release.config.yaml\n")

	return nil
}

func runPluginUninstall(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	pluginName := args[0]

	mgr, err := manager.NewManager()
	if err != nil {
		return fmt.Errorf("failed to create plugin manager: %w", err)
	}

	fmt.Printf("Uninstalling plugin %q...\n", pluginName)

	if err := mgr.Uninstall(ctx, pluginName); err != nil {
		return fmt.Errorf("failed to uninstall plugin: %w", err)
	}

	printSuccess(fmt.Sprintf("Plugin %q uninstalled successfully", pluginName))

	return nil
}

func runPluginEnable(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	pluginName := args[0]

	mgr, err := manager.NewManager()
	if err != nil {
		return fmt.Errorf("failed to create plugin manager: %w", err)
	}

	if err := mgr.Enable(ctx, pluginName); err != nil {
		return fmt.Errorf("failed to enable plugin: %w", err)
	}

	printSuccess(fmt.Sprintf("Plugin %q enabled", pluginName))
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("  1. Configure the plugin in release.config.yaml")
	fmt.Println("  2. Run release-pilot commands to use the plugin")

	return nil
}

func runPluginDisable(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	pluginName := args[0]

	mgr, err := manager.NewManager()
	if err != nil {
		return fmt.Errorf("failed to create plugin manager: %w", err)
	}

	if err := mgr.Disable(ctx, pluginName); err != nil {
		return fmt.Errorf("failed to disable plugin: %w", err)
	}

	printSuccess(fmt.Sprintf("Plugin %q disabled", pluginName))

	return nil
}
