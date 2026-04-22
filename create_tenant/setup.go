package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/sethvargo/go-password/password"

	// Internal SDK for tenant operations
	internalClient "github.com/spectrocloud/palette-sdk-go-internal/client"
	"github.com/spectrocloud/hapi/models"

	// Public SDK for edge token operations (after tenant login)
	"github.com/spectrocloud/palette-sdk-go/api/client/version1"
	publicModels "github.com/spectrocloud/palette-sdk-go/api/models"
	publicClient "github.com/spectrocloud/palette-sdk-go/client"
)

// TenantData holds the tenant setup results
type TenantData struct {
	TenantUID     string
	TenantName    string
	AdminEmail    string
	AdminPassword string
	EdgeToken     string
	APIKey        string
	PaletteHost   string
}

// TenantSetup creates a new tenant with the Instruqt participant ID
func TenantSetup(sysAdminUsername, sysAdminPassword string) (*TenantData, error) {
	SANDBOX_ID := os.Getenv("INSTRUQT_PARTICIPANT_ID")
	host := os.Getenv("PALETTE_HOST")

	if host == "" {
		return nil, fmt.Errorf("PALETTE_HOST environment variable is required")
	}
	if SANDBOX_ID == "" {
		return nil, fmt.Errorf("INSTRUQT_PARTICIPANT_ID environment variable is required")
	}

	tenantName := fmt.Sprintf("instruqt-%s", SANDBOX_ID)
	adminEmail := fmt.Sprintf("instruqt+%s@spectrocloud.com", SANDBOX_ID)

	// Generate random password with L3@rN- prefix
	randomPart, err := password.Generate(16, 4, 4, false, false)
	if err != nil {
		return nil, fmt.Errorf("failed to generate password: %w", err)
	}
	adminPassword := "L3@rN-" + randomPart

	log.Printf("Creating tenant: %s with admin: %s", tenantName, adminEmail)

	// Step 1: Create system-scoped client using internal SDK
	// Internal SDK expects just hostname, not full URL with scheme
	hostOnly := extractHost(host)
	sysClient := internalClient.New(
		internalClient.WithHubbleURI(hostOnly),
		internalClient.WithScopeSystem(sysAdminUsername, sysAdminPassword),
	)

	// Step 2: Create tenant
	tenantEntity := &models.V1TenantEntity{
		Metadata: &models.V1ObjectMeta{},
		Spec: &models.V1TenantSpecEntity{
			OrgName:   tenantName,
			FirstName: "Instruqt",
			LastName:  "Lab",
			EmailID:   adminEmail,
		},
	}

	tenantUID, err := sysClient.CreateTenant(tenantEntity)
	if err != nil {
		return nil, fmt.Errorf("failed to create tenant: %w", err)
	}
	log.Printf("Tenant created with UID: %s", tenantUID)

	// Step 3: Get password token for activation
	// The tenant UID is also the admin user UID for the first user
	passwordToken, err := sysClient.GetPasswordToken(tenantUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get password token: %w", err)
	}
	log.Printf("Got password token for activation")

	// Step 4: Activate user with password
	err = sysClient.ActivateUser(adminPassword, passwordToken)
	if err != nil {
		return nil, fmt.Errorf("failed to activate user: %w", err)
	}
	log.Printf("User activated successfully")

	// Step 5: Login as tenant admin using public SDK and create edge token
	// Public SDK also expects just hostname, not full URL with scheme
	tenantClient := publicClient.New(
		publicClient.WithPaletteURI(hostOnly),
		publicClient.WithUsername(adminEmail),
		publicClient.WithPassword(adminPassword),
	)
	publicClient.WithScopeTenant()(tenantClient)

	// Get the Default project UID
	defaultProjectUID, err := tenantClient.GetProjectUID("Default")
	if err != nil {
		return nil, fmt.Errorf("failed to get Default project UID: %w", err)
	}
	log.Printf("Found Default project UID: %s", defaultProjectUID)

	// Create edge registration token with Default project
	edgeEntity := &publicModels.V1EdgeTokenEntity{
		Metadata: &publicModels.V1ObjectMeta{
			Name: fmt.Sprintf("instruqt-%s", SANDBOX_ID),
		},
		Spec: &publicModels.V1EdgeTokenSpecEntity{
			DefaultProjectUID: defaultProjectUID,
			Expiry:            publicModels.V1Time(time.Now().Add(7 * 24 * time.Hour)),
		},
	}

	edgeTokenParams := version1.NewV1EdgeTokensCreateParams().WithBody(edgeEntity)
	edgeTokenResp, err := tenantClient.Client.V1EdgeTokensCreate(edgeTokenParams)
	if err != nil {
		return nil, fmt.Errorf("failed to create edge token: %w", err)
	}

	// Get the actual token value
	edgeTokenGetParams := version1.NewV1EdgeTokensUIDGetParams().WithUID(*edgeTokenResp.Payload.UID)
	edgeTokenDetails, err := tenantClient.Client.V1EdgeTokensUIDGet(edgeTokenGetParams)
	if err != nil {
		return nil, fmt.Errorf("failed to get edge token details: %w", err)
	}

	edgeToken := edgeTokenDetails.Payload.Spec.Token
	log.Printf("Edge registration token created")

	// Step 6: Create API key for tenant
	apiKeyName := fmt.Sprintf("instruqt-%s", SANDBOX_ID)
	apiKey, err := tenantClient.CreateAPIKey(apiKeyName, nil, 7*24*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("failed to create API key: %w", err)
	}
	log.Printf("API key created")

	result := &TenantData{
		TenantUID:     tenantUID,
		TenantName:    tenantName,
		AdminEmail:    adminEmail,
		AdminPassword: adminPassword,
		EdgeToken:     edgeToken,
		APIKey:        apiKey,
		PaletteHost:   host,
	}

	// Set environment variables for downstream use
	os.Setenv("PALETTE_TENANT_UID", tenantUID)
	os.Setenv("PALETTE_TENANT_NAME", tenantName)
	os.Setenv("PALETTE_ADMIN_EMAIL", adminEmail)
	os.Setenv("PALETTE_ADMIN_PASSWORD", adminPassword)
	os.Setenv("PALETTE_EDGE_TOKEN", edgeToken)
	os.Setenv("PALETTE_API_KEY", apiKey)

	// Generate HTML output if template exists
	if err := generateHTML(result); err != nil {
		log.Printf("Warning: failed to generate HTML: %v", err)
	}

	// Generate JSON file for other containers to fetch via HTTP
	if err := generateJSON(result); err != nil {
		log.Printf("Warning: failed to generate JSON: %v", err)
	}

	log.Printf("Tenant setup complete!")
	log.Printf("  Tenant: %s", tenantName)
	log.Printf("  Admin Email: %s", adminEmail)
	log.Printf("  Admin Password: %s", adminPassword)
	log.Printf("  Edge Token: %s", edgeToken)
	log.Printf("  API Key: %s", apiKey)

	return result, nil
}
