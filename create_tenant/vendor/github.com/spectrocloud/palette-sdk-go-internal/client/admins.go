package client

import (
	"fmt"

	"github.com/spectrocloud/hapi/models"

	systemC "github.com/spectrocloud/hapi/system/client/v1"
)

func (h *V1Client) GetSystemAdmins() (*models.V1SystemAdmins, error) {
	params := systemC.V1SystemAdminsListParams{}

	resp, err := h.SystemC.V1SystemAdminsList(&params)
	if err != nil {
		return nil, err
	}
	return resp.Payload, err
}

func (h *V1Client) GetSystemAdminByEmail(email string) (*models.V1SystemAdmin, error) {
	admins, err := h.GetSystemAdmins()
	if err != nil {
		return nil, err
	}

	if admins != nil {
		for _, admin := range admins.Items {
			if admin != nil && admin.Spec != nil {
				if admin.Spec.EmailID != nil && *admin.Spec.EmailID == email {
					return admin, nil
				}
			}
		}
	}
	return nil, fmt.Errorf("admin with email '%s' not found", email)
}

func (h *V1Client) DeleteSystemAdmin(uid string) error {
	params := systemC.NewV1SystemAdminsUIDDeleteParams().WithUID(uid)

	_, err := h.SystemC.V1SystemAdminsUIDDelete(params)

	return err
}
