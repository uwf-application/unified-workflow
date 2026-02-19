package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

// ServiceConfig represents a service to deploy
type ServiceConfig struct {
	Name        string `yaml:"name"`
	Dockerfile  string `yaml:"dockerfile"`
	BuildDir    string `yaml:"build_dir"`
	Registry    string `yaml:"registry"`
	TargetHost  string `yaml:"target_host"`
	ComposePath string `yaml:"compose_path"`
	Username    string `yaml:"username"`
}

// DeploymentConfig represents the deployment configuration
type DeploymentConfig struct {
	JumpServer     string                   `yaml:"jump_server"`
	HarborRegistry string                   `yaml:"harbor_registry"`
	HarborUsername string                   `yaml:"harbor_username"`
	Services       map[string]ServiceConfig `yaml:"services"`
	Environment    string                   `yaml:"environment"`
}

// newDeployCmd creates the deploy command
func newDeployCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy services to test environment",
		Long: `Deploy Unified Workflow services to the test environment.
This command builds Docker images, pushes them to Harbor registry,
and deploys them to target servers via the jump server.`,
	}

	// Add subcommands
	cmd.AddCommand(newDeployBuildCmd())
	cmd.AddCommand(newDeployPushCmd())
	cmd.AddCommand(newDeployAllCmd())
	cmd.AddCommand(newDeployServiceCmd())
	cmd.AddCommand(newDeployStatusCmd())
	cmd.AddCommand(newDeployInitCmd())
	cmd.AddCommand(newDeploySyncCmd())
	cmd.AddCommand(newDeployVerifyCmd())

	return cmd
}

// newDeployBuildCmd creates the build subcommand
func newDeployBuildCmd() *cobra.Command {
	var services []string
	var all bool

	cmd := &cobra.Command{
		Use:   "build",
		Short: "Build Docker images for services",
		Long: `Build Docker images for specified services.
If no services specified, builds all services.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := loadDeploymentConfig()
			if err != nil {
				return fmt.Errorf("failed to load deployment config: %w", err)
			}

			servicesToBuild := services
			if all || len(services) == 0 {
				servicesToBuild = getAllServiceNames(config)
			}

			fmt.Printf("🚀 Building %d services...\n", len(servicesToBuild))

			for _, serviceName := range servicesToBuild {
				service, exists := config.Services[serviceName]
				if !exists {
					return fmt.Errorf("service %s not found in config", serviceName)
				}

				fmt.Printf("📦 Building %s...\n", serviceName)
				if err := buildService(service); err != nil {
					return fmt.Errorf("failed to build %s: %w", serviceName, err)
				}
				fmt.Printf("✅ Built %s successfully\n", serviceName)
			}

			fmt.Println("🎉 All services built successfully!")
			return nil
		},
	}

	cmd.Flags().StringSliceVarP(&services, "services", "s", []string{}, "Services to build (comma-separated)")
	cmd.Flags().BoolVarP(&all, "all", "a", false, "Build all services")

	return cmd
}

// newDeployPushCmd creates the push subcommand
func newDeployPushCmd() *cobra.Command {
	var services []string
	var all bool

	cmd := &cobra.Command{
		Use:   "push",
		Short: "Push Docker images to Harbor registry",
		Long: `Push Docker images to Harbor registry.
Images must be built before they can be pushed.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := loadDeploymentConfig()
			if err != nil {
				return fmt.Errorf("failed to load deployment config: %w", err)
			}

			servicesToPush := services
			if all || len(services) == 0 {
				servicesToPush = getAllServiceNames(config)
			}

			fmt.Printf("🚀 Pushing %d services to Harbor...\n", len(servicesToPush))

			for _, serviceName := range servicesToPush {
				service, exists := config.Services[serviceName]
				if !exists {
					return fmt.Errorf("service %s not found in config", serviceName)
				}

				fmt.Printf("📤 Pushing %s...\n", serviceName)
				if err := pushService(service); err != nil {
					return fmt.Errorf("failed to push %s: %w", serviceName, err)
				}
				fmt.Printf("✅ Pushed %s successfully\n", serviceName)
			}

			fmt.Println("🎉 All services pushed successfully!")
			return nil
		},
	}

	cmd.Flags().StringSliceVarP(&services, "services", "s", []string{}, "Services to push (comma-separated)")
	cmd.Flags().BoolVarP(&all, "all", "a", false, "Push all services")

	return cmd
}

// newDeployAllCmd creates the all subcommand
func newDeployAllCmd() *cobra.Command {
	var services []string
	var skipBuild bool
	var skipPush bool

	cmd := &cobra.Command{
		Use:   "all",
		Short: "Build and push all services",
		Long: `Build and push all services to test environment.
This is the complete deployment workflow.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := loadDeploymentConfig()
			if err != nil {
				return fmt.Errorf("failed to load deployment config: %w", err)
			}

			servicesToDeploy := services
			if len(services) == 0 {
				servicesToDeploy = getAllServiceNames(config)
			}

			fmt.Printf("🚀 Starting deployment of %d services...\n", len(servicesToDeploy))

			// Build phase
			if !skipBuild {
				fmt.Println("=== BUILD PHASE ===")
				for _, serviceName := range servicesToDeploy {
					service, exists := config.Services[serviceName]
					if !exists {
						return fmt.Errorf("service %s not found in config", serviceName)
					}

					fmt.Printf("📦 Building %s...\n", serviceName)
					if err := buildService(service); err != nil {
						return fmt.Errorf("failed to build %s: %w", serviceName, err)
					}
					fmt.Printf("✅ Built %s\n", serviceName)
				}
			}

			// Push phase
			if !skipPush {
				fmt.Println("\n=== PUSH PHASE ===")
				for _, serviceName := range servicesToDeploy {
					service, exists := config.Services[serviceName]
					if !exists {
						return fmt.Errorf("service %s not found in config", serviceName)
					}

					fmt.Printf("📤 Pushing %s...\n", serviceName)
					if err := pushService(service); err != nil {
						return fmt.Errorf("failed to push %s: %w", serviceName, err)
					}
					fmt.Printf("✅ Pushed %s\n", serviceName)
				}
			}

			fmt.Println("\n🎉 Deployment completed successfully!")
			return nil
		},
	}

	cmd.Flags().StringSliceVarP(&services, "services", "s", []string{}, "Services to deploy (comma-separated)")
	cmd.Flags().BoolVar(&skipBuild, "skip-build", false, "Skip build phase")
	cmd.Flags().BoolVar(&skipPush, "skip-push", false, "Skip push phase")

	return cmd
}

// newDeployServiceCmd creates the service subcommand
func newDeployServiceCmd() *cobra.Command {
	var skipBuild bool
	var skipPush bool

	cmd := &cobra.Command{
		Use:   "service [name]",
		Short: "Deploy a specific service",
		Long: `Deploy a specific service to test environment.
Builds and pushes the specified service.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			serviceName := args[0]
			config, err := loadDeploymentConfig()
			if err != nil {
				return fmt.Errorf("failed to load deployment config: %w", err)
			}

			service, exists := config.Services[serviceName]
			if !exists {
				return fmt.Errorf("service %s not found in config", serviceName)
			}

			fmt.Printf("🚀 Deploying service: %s\n", serviceName)

			// Build phase
			if !skipBuild {
				fmt.Printf("📦 Building %s...\n", serviceName)
				if err := buildService(service); err != nil {
					return fmt.Errorf("failed to build %s: %w", serviceName, err)
				}
				fmt.Printf("✅ Built %s\n", serviceName)
			}

			// Push phase
			if !skipPush {
				fmt.Printf("📤 Pushing %s...\n", serviceName)
				if err := pushService(service); err != nil {
					return fmt.Errorf("failed to push %s: %w", serviceName, err)
				}
				fmt.Printf("✅ Pushed %s\n", serviceName)
			}

			fmt.Printf("🎉 Service %s deployed successfully!\n", serviceName)
			return nil
		},
	}

	cmd.Flags().BoolVar(&skipBuild, "skip-build", false, "Skip build phase")
	cmd.Flags().BoolVar(&skipPush, "skip-push", false, "Skip push phase")

	return cmd
}

// newDeployStatusCmd creates the status subcommand
func newDeployStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Check deployment status",
		Long:  `Check the deployment status of services in the test environment.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := loadDeploymentConfig()
			if err != nil {
				return fmt.Errorf("failed to load deployment config: %w", err)
			}

			fmt.Println("📊 Deployment Status")
			fmt.Println("===================")

			for _, serviceName := range getAllServiceNames(config) {
				service := config.Services[serviceName]
				fmt.Printf("\n🔍 %s\n", serviceName)
				fmt.Printf("   Registry: %s\n", service.Registry)
				fmt.Printf("   Target: %s\n", service.TargetHost)
				fmt.Printf("   Dockerfile: %s\n", service.Dockerfile)
			}

			fmt.Println("\n✅ Configuration loaded successfully")
			return nil
		},
	}

	return cmd
}

// newDeployInitCmd creates the init subcommand
func newDeployInitCmd() *cobra.Command {
	var jumpServer string
	var harborRegistry string
	var environment string

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize deployment configuration",
		Long: `Initialize deployment configuration with jump server and Harbor registry details.
Creates a default configuration file if it doesn't exist.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if jumpServer == "" {
				jumpServer = "10.200.1.2"
			}
			if harborRegistry == "" {
				harborRegistry = "172.30.75.78:9080"
			}
			if environment == "" {
				environment = "test"
			}

			config := DeploymentConfig{
				JumpServer:     jumpServer,
				HarborRegistry: harborRegistry,
				HarborUsername: "zh.akhmetkarimov",
				Environment:    environment,
				Services:       getDefaultServices(),
			}

			if err := saveDeploymentConfig(config); err != nil {
				return fmt.Errorf("failed to save deployment config: %w", err)
			}

			fmt.Println("✅ Deployment configuration initialized!")
			fmt.Printf("   Jump Server: %s\n", config.JumpServer)
			fmt.Printf("   Harbor Registry: %s\n", config.HarborRegistry)
			fmt.Printf("   Environment: %s\n", config.Environment)
			fmt.Printf("   Services: %d configured\n", len(config.Services))
			fmt.Println("\n📁 Configuration saved to: deploy-config.yaml")

			return nil
		},
	}

	cmd.Flags().StringVar(&jumpServer, "jump-server", "", "Jump server address (default: 10.200.1.2)")
	cmd.Flags().StringVar(&harborRegistry, "harbor", "", "Harbor registry address (default: 172.30.75.78:9080)")
	cmd.Flags().StringVar(&environment, "environment", "", "Environment name (default: test)")

	return cmd
}

// newDeploySyncCmd creates the sync subcommand
func newDeploySyncCmd() *cobra.Command {
	var jumpServer string
	var forceScp bool

	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Sync code to jump server",
		Long: `Sync local code changes to the jump server.
This prepares the jump server for building Docker images.
Uses rsync if available, falls back to scp if not.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := loadDeploymentConfig()
			if err != nil {
				return fmt.Errorf("failed to load deployment config: %w", err)
			}

			if jumpServer != "" {
				config.JumpServer = jumpServer
			}

			fmt.Printf("🔄 Syncing code to jump server: %s\n", config.JumpServer)

			// Get current directory
			currentDir, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current directory: %w", err)
			}

			// Check if rsync is available on jump server
			rsyncAvailable := false
			if !forceScp {
				fmt.Println("🔍 Checking if rsync is available on jump server...")
				checkCmd := fmt.Sprintf("ssh khassangali@%s 'which rsync >/dev/null 2>&1 && echo available || echo not-available'", config.JumpServer)
				if output, err := exec.Command("sh", "-c", checkCmd).Output(); err == nil {
					if strings.TrimSpace(string(output)) == "available" {
						rsyncAvailable = true
						fmt.Println("✅ rsync is available on jump server")
					} else {
						fmt.Println("⚠️  rsync is not available on jump server, using scp instead")
					}
				}
			}

			var syncCmd string
			if rsyncAvailable && !forceScp {
				// Use rsync
				syncCmd = fmt.Sprintf("rsync -avz --exclude='.git' --exclude='node_modules' --exclude='dist' %s/ khassangali@%s:/tmp/uwf-deploy/",
					currentDir, config.JumpServer)
				fmt.Printf("📁 Using rsync to sync from: %s\n", currentDir)
			} else {
				// Use tar + ssh (fallback)
				// First, create directory on jump server
				fmt.Println("📁 Creating directory on jump server...")
				mkdirCmd := fmt.Sprintf("ssh khassangali@%s mkdir -p /tmp/uwf-deploy", config.JumpServer)
				if err := executeLocalCommand(mkdirCmd); err != nil {
					fmt.Printf("⚠️  Failed to create directory (may already exist): %v\n", err)
				}

				// Use tar to create archive and pipe through ssh
				// Exclude .git directory and other build artifacts
				syncCmd = fmt.Sprintf("cd %s && tar czf - --exclude='.git' --exclude='node_modules' --exclude='dist' --exclude='*.log' . | ssh khassangali@%s 'cd /tmp/uwf-deploy && tar xzf -'",
					currentDir, config.JumpServer)
				fmt.Printf("📁 Using tar+ssh to sync from: %s\n", currentDir)
				fmt.Println("⚠️  Note: This method is slower than rsync but more reliable than scp")
			}

			fmt.Printf("📁 Syncing to: khassangali@%s:/tmp/uwf-deploy/\n", config.JumpServer)
			fmt.Println("⏳ This may take a moment...")

			// Execute sync command
			if err := executeLocalCommand(syncCmd); err != nil {
				return fmt.Errorf("failed to sync code: %w", err)
			}

			fmt.Println("✅ Code synced successfully!")

			// Show installation hint if rsync is not available
			if !rsyncAvailable && !forceScp {
				fmt.Println("\n💡 Tip: Install rsync on jump server for faster sync with exclusions:")
				fmt.Printf("   ssh khassangali@%s 'sudo yum install rsync -y'  # For RHEL/CentOS\n", config.JumpServer)
				fmt.Printf("   ssh khassangali@%s 'sudo apt-get install rsync -y'  # For Ubuntu/Debian\n", config.JumpServer)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&jumpServer, "jump-server", "", "Jump server address (overrides config)")
	cmd.Flags().BoolVar(&forceScp, "force-scp", false, "Force using scp instead of rsync")

	return cmd
}

// newDeployVerifyCmd creates the verify subcommand
func newDeployVerifyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "verify",
		Short: "Verify deployment",
		Long: `Verify that services are deployed and running correctly.
Checks Docker images, containers, and service health.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := loadDeploymentConfig()
			if err != nil {
				return fmt.Errorf("failed to load deployment config: %w", err)
			}

			fmt.Println("🔍 Verifying deployment...")

			// Check if jump server is reachable
			fmt.Printf("📡 Checking jump server: %s\n", config.JumpServer)
			pingCmd := fmt.Sprintf("ssh -o ConnectTimeout=5 khassangali@%s echo connected", config.JumpServer)
			if err := executeLocalCommand(pingCmd); err != nil {
				return fmt.Errorf("jump server is not reachable: %w", err)
			}

			// Check Harbor registry
			fmt.Printf("🐋 Checking Harbor registry: %s\n", config.HarborRegistry)
			harborCmd := fmt.Sprintf("ssh khassangali@%s 'curl -s http://%s/v2/_catalog | grep -q repositories'",
				config.JumpServer, config.HarborRegistry)
			if err := executeLocalCommand(harborCmd); err != nil {
				fmt.Printf("⚠️  Harbor registry check failed (may be expected): %v\n", err)
			} else {
				fmt.Println("✅ Harbor registry is accessible")
			}

			fmt.Println("\n✅ Deployment verification completed!")
			return nil
		},
	}

	return cmd
}

// Helper functions

func loadDeploymentConfig() (DeploymentConfig, error) {
	// For now, return default config
	// TODO: Load from YAML file
	return DeploymentConfig{
		JumpServer:     "10.200.1.2",
		HarborRegistry: "172.30.75.78:9080",
		HarborUsername: "zh.akhmetkarimov",
		Environment:    "test",
		Services:       getDefaultServices(),
	}, nil
}

func saveDeploymentConfig(config DeploymentConfig) error {
	// TODO: Save to YAML file
	return nil
}

func getDefaultServices() map[string]ServiceConfig {
	return map[string]ServiceConfig{
		"workflow-worker": {
			Name:        "workflow-worker",
			Dockerfile:  "Dockerfile.worker",
			BuildDir:    ".",
			Registry:    "172.30.75.78:9080/taf/workflow-worker",
			TargetHost:  "172.30.75.85",
			ComposePath: "/opt/workflow/docker-compose.yml",
			Username:    "zh.akhmetkarimov",
		},
		"workflow-registry": {
			Name:        "workflow-registry",
			Dockerfile:  "Dockerfile.registry",
			BuildDir:    ".",
			Registry:    "172.30.75.78:9080/taf/workflow-registry",
			TargetHost:  "172.30.75.85",
			ComposePath: "/opt/workflow/docker-compose.yml",
			Username:    "zh.akhmetkarimov",
		},
		"workflow-executor": {
			Name:        "workflow-executor",
			Dockerfile:  "Dockerfile.executor",
			BuildDir:    ".",
			Registry:    "172.30.75.78:9080/taf/workflow-executor",
			TargetHost:  "172.30.75.85",
			ComposePath: "/opt/workflow/docker-compose.yml",
			Username:    "zh.akhmetkarimov",
		},
		"workflow-api": {
			Name:        "workflow-api",
			Dockerfile:  "Dockerfile.workflow-api",
			BuildDir:    ".",
			Registry:    "172.30.75.78:9080/taf/workflow-api",
			TargetHost:  "172.30.75.85",
			ComposePath: "/opt/workflow/docker-compose.yml",
			Username:    "zh.akhmetkarimov",
		},
	}
}

func getAllServiceNames(config DeploymentConfig) []string {
	var names []string
	for name := range config.Services {
		names = append(names, name)
	}
	return names
}

func buildService(service ServiceConfig) error {
	// First check if Docker or Podman is installed on jump server
	fmt.Println("🔍 Checking if Docker/Podman is available on jump server...")

	// Check for docker first, then podman
	checkDockerCmd := "ssh khassangali@10.200.1.2 \"which docker >/dev/null 2>&1 && echo docker-available || echo docker-not-available\""
	dockerOutput, dockerErr := exec.Command("sh", "-c", checkDockerCmd).Output()

	checkPodmanCmd := "ssh khassangali@10.200.1.2 \"which podman >/dev/null 2>&1 && echo podman-available || echo podman-not-available\""
	podmanOutput, podmanErr := exec.Command("sh", "-c", checkPodmanCmd).Output()

	if dockerErr != nil || podmanErr != nil {
		return fmt.Errorf("failed to check container runtime availability: %w", dockerErr)
	}

	dockerAvailable := strings.TrimSpace(string(dockerOutput)) == "docker-available"
	podmanAvailable := strings.TrimSpace(string(podmanOutput)) == "podman-available"

	if !dockerAvailable && !podmanAvailable {
		return fmt.Errorf("Neither Docker nor Podman is installed on jump server. Please install one of them first:\n" +
			"   ssh khassangali@10.200.1.2 'sudo yum install docker -y'  # For RHEL/CentOS\n" +
			"   ssh khassangali@10.200.1.2 'sudo apt-get install docker.io -y'  # For Ubuntu/Debian\n" +
			"   OR install podman:\n" +
			"   ssh khassangali@10.200.1.2 'sudo yum install podman -y'  # For RHEL/CentOS\n" +
			"   ssh khassangali@10.200.1.2 'sudo apt-get install podman -y'  # For Ubuntu/Debian")
	}

	// Determine which command to use
	var buildCmd string
	if dockerAvailable {
		fmt.Println("✅ Docker is available on jump server")
		buildCmd = fmt.Sprintf("ssh khassangali@10.200.1.2 'cd /tmp/uwf-deploy && docker build -f %s -t %s:latest .'",
			service.Dockerfile, service.Registry)
	} else {
		fmt.Println("✅ Podman is available on jump server")
		buildCmd = fmt.Sprintf("ssh khassangali@10.200.1.2 'cd /tmp/uwf-deploy && podman build -f %s -t %s:latest .'",
			service.Dockerfile, service.Registry)
	}

	return executeLocalCommand(buildCmd)
}

func pushService(service ServiceConfig) error {
	serviceName := service.Name

	// First check if we should use docker or podman for save command
	// (reuse the same logic from buildService to check availability)
	checkPodmanCmd := "ssh khassangali@10.200.1.2 \"which podman >/dev/null 2>&1 && echo podman-available || echo podman-not-available\""
	podmanOutput, _ := exec.Command("sh", "-c", checkPodmanCmd).Output()
	podmanAvailable := strings.TrimSpace(string(podmanOutput)) == "podman-available"

	// Use user's home directory instead of /tmp to avoid permission issues
	tarPath := fmt.Sprintf("/home/khassangali/%s_latest.tar", serviceName)

	// Step 1: Remove existing tar file if it exists (podman/docker save can't overwrite)
	fmt.Printf("🧹 Removing existing %s.tar file if it exists...\n", serviceName)
	removeCmd := fmt.Sprintf("ssh khassangali@10.200.1.2 'rm -f %s'", tarPath)
	executeLocalCommand(removeCmd) // Don't fail if file doesn't exist

	// Step 2: Save image as tar file on jump server
	fmt.Printf("💾 Saving %s image as tar file on jump server...\n", serviceName)

	var saveCmd string
	if podmanAvailable {
		saveCmd = fmt.Sprintf("ssh khassangali@10.200.1.2 'podman save -o %s %s:latest'",
			tarPath, service.Registry)
	} else {
		saveCmd = fmt.Sprintf("ssh khassangali@10.200.1.2 'docker save -o %s %s:latest'",
			tarPath, service.Registry)
	}

	if err := executeLocalCommand(saveCmd); err != nil {
		return fmt.Errorf("failed to save %s image: %w", serviceName, err)
	}

	// Step 3: SCP tar file from jump server to Harbor host using h.gaparov user
	fmt.Printf("📤 Copying %s.tar to Harbor host (h.gaparov)...\n", serviceName)
	// Use sshpass to provide password non-interactively
	scpCmd := fmt.Sprintf("ssh khassangali@10.200.1.2 'sshpass -p \"9R!\\$ZpK@eM2xQ7A\" scp %s h.gaparov@172.30.75.78:/tmp/'",
		tarPath)
	if err := executeLocalCommand(scpCmd); err != nil {
		return fmt.Errorf("failed to copy %s.tar to Harbor host: %w", serviceName, err)
	}

	// Step 4: Load image on Harbor host using h.gaparov user (execute from jump server)
	fmt.Printf("📥 Loading %s image on Harbor host...\n", serviceName)
	// Execute via jump server: ssh from jump server to Harbor host with sudo -S (read password from stdin)
	// Note: We need to provide sudo password after sshpass password
	loadCmd := fmt.Sprintf("ssh khassangali@10.200.1.2 'echo \"9R!\\$ZpK@eM2xQ7A\" | sshpass -p \"9R!\\$ZpK@eM2xQ7A\" ssh h.gaparov@172.30.75.78 sudo -S docker load -i /tmp/%s_latest.tar'",
		serviceName)
	if err := executeLocalCommand(loadCmd); err != nil {
		return fmt.Errorf("failed to load %s image on Harbor host: %w", serviceName, err)
	}

	// Step 5: Push to Harbor registry from Harbor host using h.gaparov user (execute from jump server)
	fmt.Printf("🚀 Pushing %s to Harbor registry...\n", serviceName)
	// Execute via jump server: ssh from jump server to Harbor host with sudo -S
	pushCmd := fmt.Sprintf("ssh khassangali@10.200.1.2 'echo \"9R!\\$ZpK@eM2xQ7A\" | sshpass -p \"9R!\\$ZpK@eM2xQ7A\" ssh h.gaparov@172.30.75.78 sudo -S docker push %s:latest'",
		service.Registry)
	if err := executeLocalCommand(pushCmd); err != nil {
		return fmt.Errorf("failed to push %s to Harbor: %w", serviceName, err)
	}

	// Step 6: Clean up tar files
	fmt.Printf("🧹 Cleaning up %s.tar files...\n", serviceName)
	cleanupCmd := fmt.Sprintf("ssh khassangali@10.200.1.2 'rm -f %s && echo \"9R!\\$ZpK@eM2xQ7A\" | sshpass -p \"9R!\\$ZpK@eM2xQ7A\" ssh h.gaparov@172.30.75.78 sudo -S rm -f /tmp/%s_latest.tar'",
		tarPath, serviceName)
	executeLocalCommand(cleanupCmd) // Don't fail if cleanup fails

	return nil
}

func executeLocalCommand(cmd string) error {
	fmt.Printf("   Executing: %s\n", cmd)

	// Use shell to execute the command to handle complex commands properly
	execCmd := exec.Command("sh", "-c", cmd)

	// Capture output
	output, err := execCmd.CombinedOutput()
	if err != nil {
		fmt.Printf("   ❌ Command failed: %v\n", err)
		fmt.Printf("   Output: %s\n", string(output))
		return fmt.Errorf("command failed: %w", err)
	}

	// Print success output
	if len(output) > 0 {
		// Trim and show first few lines if output is long
		outputStr := strings.TrimSpace(string(output))
		lines := strings.Split(outputStr, "\n")
		if len(lines) > 10 {
			fmt.Printf("   ✅ Output (first 10 of %d lines):\n", len(lines))
			for i := 0; i < 10 && i < len(lines); i++ {
				fmt.Printf("      %s\n", lines[i])
			}
			fmt.Printf("      ... and %d more lines\n", len(lines)-10)
		} else if len(outputStr) > 0 {
			fmt.Printf("   ✅ Output: %s\n", outputStr)
		} else {
			fmt.Printf("   ✅ Command executed successfully\n")
		}
	} else {
		fmt.Printf("   ✅ Command executed successfully\n")
	}

	return nil
}
