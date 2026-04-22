package client

import (
	"fmt"

	"github.com/go-openapi/strfmt"
	"github.com/spectrocloud/gomi/pkg/logger"
	authC "github.com/spectrocloud/hapi/auth/client/v1"
	"github.com/spectrocloud/hapi/models"
	userC "github.com/spectrocloud/hapi/user/client/v1"
)

// CreateTenant creates a new tenant in the system.
func (h *V1Client) CreateTenant(body *models.V1TenantEntity) (string, error) {
	params := userC.NewV1TenantsCreateParams().WithBody(body)
	resp, err := h.UserC.V1TenantsCreate(params)
	if err != nil {
		return "", fmt.Errorf("failed to create tenant: %w", err)
	}
	return *resp.Payload.UID, nil
}

// GetTenants gets all tenants within the system.
func (h *V1Client) GetTenants() (*models.V1Tenants, error) {
	params := userC.NewV1TenantsListParams()
	resp, err := h.UserC.V1TenantsList(params)
	if err != nil {
		return nil, fmt.Errorf("failed to list tenants: %w", err)
	}

	return resp.Payload, nil
}

// GetTenantByName gets a tenant by organization name.
func (h *V1Client) GetTenantByName(name string) (*models.V1Tenant, error) {
	tenants, err := h.GetTenants()
	if err != nil {
		return nil, err
	}

	for _, tenant := range tenants.Items {
		if tenant.Spec.OrgName == name {
			return tenant, nil
		}
	}

	return nil, fmt.Errorf("tenant with name '%s' not found", name)
}

// DeleteTenant deletes a tenant in the system by uid.
func (h *V1Client) DeleteTenant(uid string) error {
	params := userC.NewV1TenantsDeleteParams().WithTenantUID(uid)
	_, err := h.UserC.V1TenantsDelete(params)
	if err != nil {
		return fmt.Errorf("failed to delete tenant: %w", err)
	}
	return nil
}

// DeleteTenant deletes a tenant in the system by organization name.
func (h *V1Client) DeleteTenantByName(name string) error {
	tenants, err := h.GetTenants()
	if err != nil {
		return err
	}

	for _, tenant := range tenants.Items {
		if tenant != nil && tenant.Spec != nil {
			if tenant.Spec.OrgName == name {
				return h.DeleteTenant(tenant.Metadata.UID)
			}
		}
	}
	return fmt.Errorf("tenant with name '%s' not found", name)
}

// CleanUpTenant performs cleanup operations on a tenant in the system.
// When force is false, it attempts to delete all resources within the tenant.
// When force is true, it performs a hard delete of all content from the database
// without proper resource cleanup.
func (h *V1Client) CleanUpTenant(uid string, force bool) error {
	params := userC.NewV1TenantsCleanUpParams().WithTenantUID(uid)
	params.SetForceDelete(&force)

	_, err := h.UserC.V1TenantsCleanUp(params)
	if err != nil {
		return fmt.Errorf("failed to cleanup tenant: %w", err)
	}

	return nil
}

// ActivateUser activates a user with the provided password and password token.
func (h *V1Client) ActivateUser(password, passwordToken string) error {
	var body authC.V1PasswordActivateBody
	pass := strfmt.Password(password)
	body.Password = &pass
	params := authC.NewV1PasswordActivateParams().WithPasswordToken(passwordToken).WithBody(body)
	_, err := h.AuthC.V1PasswordActivate(params)
	if err != nil {
		return fmt.Errorf("failed to activate tenant: %w", err)
	}
	logger.Info("User activated")
	return nil
}

// GetPasswordToken retrieves the password token for a given user ID (UID).
func (h *V1Client) GetPasswordToken(uid string) (string, error) {
	params := userC.NewV1UsersPasswordTokenParams().WithUID(uid)
	resp, err := h.UserC.V1UsersPasswordToken(params)
	if err != nil {
		return "", fmt.Errorf("failed to get password token: %w", err)
	}
	if resp.Payload == nil || resp.Payload.PasswordToken == "" {
		return "", fmt.Errorf("received empty password token for UID: %s", uid)
	}
	return resp.Payload.PasswordToken, nil
}

// SysAdminLogin performs a login operation for a sysadmin user.
func (h *V1Client) SysAdminLogin(username, password string) (string, error) {
	httpTransport := h.baseTransport()
	authClient := authC.New(httpTransport, strfmt.Default)
	params := &authC.V1SysLoginParams{
		Body: &models.V1SysLogin{
			Username: username,
			Password: strfmt.Password(password),
		},
	}
	resp, err := authClient.V1SysLogin(params)
	if err != nil {
		return "", err
	}
	h.jwt = resp.Payload.Authorization
	return h.jwt, nil
}
