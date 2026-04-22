package client

import (
	"fmt"

	authC "github.com/spectrocloud/hapi/auth/client/v1"
	"github.com/spectrocloud/hapi/models"
	"github.com/spectrocloud/palette-sdk-go/client/herr"
)

func (h *V1Client) GetOrgLoginBanner(orgName string) (*models.V1LoginBannerSettings, error) {
	params := authC.NewV1AuthOrgLoginBannerGetParamsWithContext(h.ctx).WithOrgName(orgName)
	resp, err := h.AuthC.V1AuthOrgLoginBannerGet(params)
	if err != nil || resp == nil {
		if herr.IsNotFound(err) {
			return nil, fmt.Errorf("invalid Organization: %s", orgName)
		}
		return nil, err
	} else if !resp.Payload.IsEnabled {
		return nil, nil
	}
	return resp.Payload, nil
}

func (h *V1Client) GetSystemLoginBanner() (*models.V1LoginBannerSettings, error) {
	params := authC.NewV1AuthSystemLoginBannerGetParamsWithContext(h.ctx)
	resp, err := h.AuthC.V1AuthSystemLoginBannerGet(params)
	if err != nil || resp == nil {
		return nil, err
	} else if !resp.Payload.IsEnabled {
		return nil, nil
	}
	return resp.Payload, nil
}

func (h *V1Client) SwitchOrganization(orgName string) (string, error) {
	params := authC.NewV1AuthOrgSwitchParamsWithContext(h.ctx).
		WithOrgName(orgName)
	resp, err := h.AuthC.V1AuthOrgSwitch(params)
	if err != nil || resp == nil {
		return "", err
	}
	return resp.Payload.Authorization, nil
}
