package main

import (
	"fmt"
	"log"
	"os"

	// Internal SDK for tenant operations
	internalClient "github.com/spectrocloud/palette-sdk-go-internal/client"
)

// TenantCleanup deletes the tenant and all its resources
func TenantCleanup(sysAdminUsername, sysAdminPassword string) error {
	SANDBOX_ID := os.Getenv("INSTRUQT_PARTICIPANT_ID")
	host := os.Getenv("PALETTE_HOST")

	if host == "" {
		return fmt.Errorf("PALETTE_HOST environment variable is required")
	}
	if SANDBOX_ID == "" {
		return fmt.Errorf("INSTRUQT_PARTICIPANT_ID environment variable is required")
	}

	tenantName := fmt.Sprintf("instruqt-%s", SANDBOX_ID)

	log.Printf("Cleaning up tenant: %s", tenantName)

	// Create system-scoped client using internal SDK
	// Internal SDK expects just hostname, not full URL with scheme
	hostOnly := extractHost(host)
	sysClient := internalClient.New(
		internalClient.WithHubbleURI(hostOnly),
		internalClient.WithScopeSystem(sysAdminUsername, sysAdminPassword),
	)

	// Find tenant by name (with pagination)
	tenant, err := findTenantByName(sysClient, tenantName)
	if err != nil {
		log.Printf("Tenant not found, may already be deleted: %v", err)
		return nil
	}
	log.Printf("Found tenant: %s (UID: %s)", tenant.Spec.OrgName, tenant.Metadata.UID)

	// Use CleanUpTenant with force=true for hard delete
	err = sysClient.CleanUpTenant(tenant.Metadata.UID, true)
	if err != nil {
		log.Printf("Warning: cleanup tenant failed (continuing with delete): %v", err)
	}

	// Delete the tenant
	err = sysClient.DeleteTenant(tenant.Metadata.UID)
	if err != nil {
		return fmt.Errorf("failed to delete tenant: %w", err)
	}

	log.Printf("Tenant %s deleted successfully", tenantName)
	return nil
}

// TenantDeleteByUID deletes a tenant directly by UID
func TenantDeleteByUID(sysAdminUsername, sysAdminPassword, tenantUID string) error {
	host := os.Getenv("PALETTE_HOST")

	if host == "" {
		return fmt.Errorf("PALETTE_HOST environment variable is required")
	}

	log.Printf("Deleting tenant by UID: %s", tenantUID)

	// Create system-scoped client using internal SDK
	hostOnly := extractHost(host)
	sysClient := internalClient.New(
		internalClient.WithHubbleURI(hostOnly),
		internalClient.WithScopeSystem(sysAdminUsername, sysAdminPassword),
	)

	// Try cleanup first
	err := sysClient.CleanUpTenant(tenantUID, true)
	if err != nil {
		log.Printf("Warning: cleanup tenant failed (continuing with delete): %v", err)
	}

	// Delete the tenant
	err = sysClient.DeleteTenant(tenantUID)
	if err != nil {
		return fmt.Errorf("failed to delete tenant: %w", err)
	}

	log.Printf("Tenant %s deleted successfully", tenantUID)
	return nil
}
