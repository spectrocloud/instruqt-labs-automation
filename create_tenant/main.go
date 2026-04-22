package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// Get sysadmin credentials from environment or arguments
	sysAdminUsername := os.Getenv("PALETTE_SYSADMIN_USERNAME")
	sysAdminPassword := os.Getenv("PALETTE_SYSADMIN_PASSWORD")

	// Allow override via command line arguments for setup/cleanup
	if len(os.Args) >= 4 && os.Args[1] != "delete-uid" && os.Args[1] != "webserver" {
		sysAdminUsername = os.Args[2]
		sysAdminPassword = os.Args[3]
	}

	switch os.Args[1] {
	case "setup":
		if sysAdminUsername == "" || sysAdminPassword == "" {
			fmt.Println("Error: Sysadmin credentials required")
			fmt.Println("Set PALETTE_SYSADMIN_USERNAME and PALETTE_SYSADMIN_PASSWORD environment variables")
			fmt.Println("Or pass them as arguments: setup setup <username> <password>")
			os.Exit(1)
		}
		result, err := TenantSetup(sysAdminUsername, sysAdminPassword)
		if err != nil {
			fmt.Printf("Setup failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("\nTenant Setup Complete!\n")
		fmt.Printf("=====================\n")
		fmt.Printf("Tenant UID:     %s\n", result.TenantUID)
		fmt.Printf("Tenant Name:    %s\n", result.TenantName)
		fmt.Printf("Admin Email:    %s\n", result.AdminEmail)
		fmt.Printf("Admin Password: %s\n", result.AdminPassword)
		fmt.Printf("Edge Token:     %s\n", result.EdgeToken)
		fmt.Printf("Palette Host:   %s\n", result.PaletteHost)

	case "cleanup":
		if sysAdminUsername == "" || sysAdminPassword == "" {
			fmt.Println("Error: Sysadmin credentials required")
			fmt.Println("Set PALETTE_SYSADMIN_USERNAME and PALETTE_SYSADMIN_PASSWORD environment variables")
			fmt.Println("Or pass them as arguments: setup cleanup <username> <password>")
			os.Exit(1)
		}
		err := TenantCleanup(sysAdminUsername, sysAdminPassword)
		if err != nil {
			fmt.Printf("Cleanup failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Tenant cleanup complete!")

	case "delete-uid":
		if len(os.Args) < 3 {
			fmt.Println("Usage: setup delete-uid <tenant_uid> [username] [password]")
			os.Exit(1)
		}
		tenantUID := os.Args[2]
		// Shift args for credentials
		if len(os.Args) >= 5 {
			sysAdminUsername = os.Args[3]
			sysAdminPassword = os.Args[4]
		}
		if sysAdminUsername == "" || sysAdminPassword == "" {
			fmt.Println("Error: Sysadmin credentials required")
			os.Exit(1)
		}
		err := TenantDeleteByUID(sysAdminUsername, sysAdminPassword, tenantUID)
		if err != nil {
			fmt.Printf("Delete failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Tenant deleted!")

	case "webserver":
		Webserver()

	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: setup <command> [arguments]")
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println("  setup              Create a new tenant for the Instruqt participant")
	fmt.Println("  cleanup            Delete the tenant and all its resources")
	fmt.Println("  delete-uid <uid>   Delete a tenant directly by UID")
	fmt.Println("  webserver          Start the HTTP server to serve static files")
	fmt.Println("")
	fmt.Println("Environment variables:")
	fmt.Println("  PALETTE_HOST               Palette API host (required for setup/cleanup)")
	fmt.Println("  INSTRUQT_PARTICIPANT_ID    Instruqt sandbox ID (required for setup/cleanup)")
	fmt.Println("  PALETTE_SYSADMIN_USERNAME  Sysadmin username")
	fmt.Println("  PALETTE_SYSADMIN_PASSWORD  Sysadmin password")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  setup setup")
	fmt.Println("  setup setup admin mypassword")
	fmt.Println("  setup cleanup")
	fmt.Println("  setup delete-uid 698b2882ddf964e58042deba")
	fmt.Println("  setup webserver")
}
