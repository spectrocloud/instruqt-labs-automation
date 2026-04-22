package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/url"
	"os"
	"strings"

	internalClient "github.com/spectrocloud/palette-sdk-go-internal/client"
	"github.com/spectrocloud/hapi/models"
	userC "github.com/spectrocloud/hapi/user/client/v1"
)

// extractHost extracts just the hostname from a URL (removes scheme)
// e.g., "https://training.spectrocloud.com" -> "training.spectrocloud.com"
func extractHost(rawURL string) string {
	// If it doesn't have a scheme, return as-is
	if !strings.Contains(rawURL, "://") {
		return rawURL
	}
	parsed, err := url.Parse(rawURL)
	if err != nil {
		// Fallback: try to strip scheme manually
		if _, after, found := strings.Cut(rawURL, "://"); found {
			return after
		}
		return rawURL
	}
	return parsed.Host
}

// findTenantByName searches for a tenant by name using server-side filtering
func findTenantByName(sysClient *internalClient.V1Client, name string) (*models.V1Tenant, error) {
	filter := fmt.Sprintf("spec.orgName=%s", name)
	params := userC.NewV1TenantsListParams().WithFilters(&filter)
	resp, err := sysClient.UserC.V1TenantsList(params)
	if err != nil {
		return nil, fmt.Errorf("failed to list tenants: %w", err)
	}

	if len(resp.Payload.Items) == 0 {
		return nil, fmt.Errorf("tenant with name '%s' not found", name)
	}

	return resp.Payload.Items[0], nil
}

// generateJSON creates a JSON file with credentials for other containers to fetch
func generateJSON(data *TenantData) error {
	// Check for assets/ (local) or current dir (container)
	outputPath := "assets/credentials.json"
	if _, err := os.Stat("assets"); os.IsNotExist(err) {
		outputPath = "credentials.json"
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create credentials.json: %w", err)
	}
	defer file.Close()

	hostOnly := extractHost(data.PaletteHost)
	tenantURL := fmt.Sprintf("https://%s.%s", data.TenantName, hostOnly)

	creds := map[string]string{
		"palette_host":   data.PaletteHost,
		"tenant_name":    data.TenantName,
		"tenant_uid":     data.TenantUID,
		"tenant_url":     tenantURL,
		"admin_email":    data.AdminEmail,
		"admin_password": data.AdminPassword,
		"edge_token":     data.EdgeToken,
		"api_key":        data.APIKey,
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(creds)
}

// generateHTML creates an HTML file from a template with tenant data
func generateHTML(data *TenantData) error {
	// Check for template in assets/ (local) or current dir (container)
	tmplPath := "assets/index.tmpl.html"
	outputPath := "assets/index.html"

	if _, err := os.Stat(tmplPath); os.IsNotExist(err) {
		// Try current directory (container environment)
		tmplPath = "index.tmpl.html"
		outputPath = "index.html"
		if _, err := os.Stat(tmplPath); os.IsNotExist(err) {
			// Template doesn't exist, skip HTML generation
			return nil
		}
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create index.html: %w", err)
	}
	defer file.Close()

	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	// Check if registration token should be shown (default: true)
	showRegToken := os.Getenv("SHOW_REGISTRATION_TOKEN") != "false"

	// Build tenant URL: https://<tenant_name>.<host>
	hostOnly := extractHost(data.PaletteHost)
	tenantURL := fmt.Sprintf("https://%s.%s", data.TenantName, hostOnly)

	// Convert to the expected template structure
	pageData := struct {
		EmailID               string
		Password              string
		PaletteHost           string
		RegistrationToken     string
		ShowRegistrationToken bool
		TenantName            string
		TenantURL             string
		Theme                 string
		APIKey                string
	}{
		EmailID:               data.AdminEmail,
		Password:              data.AdminPassword,
		PaletteHost:           data.PaletteHost,
		RegistrationToken:     data.EdgeToken,
		ShowRegistrationToken: showRegToken,
		TenantName:            data.TenantName,
		TenantURL:             tenantURL,
		Theme:                 os.Getenv("INSTRUQT_THEME"),
		APIKey:                data.APIKey,
	}

	return tmpl.Execute(file, pageData)
}
