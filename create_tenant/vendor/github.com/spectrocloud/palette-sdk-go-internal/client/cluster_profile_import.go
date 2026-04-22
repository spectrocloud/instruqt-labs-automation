package client

import (
	"errors"

	"github.com/spectrocloud/hapi/apiutil/transport"
	"github.com/spectrocloud/hapi/models"
	clusterC "github.com/spectrocloud/hapi/spectrocluster/client/v1"
)

func (h *V1Client) ClusterProfileExport(uid string) (*models.V1ClusterProfile, error) {
	params := clusterC.NewV1ClusterProfilesGetParamsWithContext(h.ctx).WithUID(uid)
	resp, err := h.ClusterC.V1ClusterProfilesGet(params)
	var e *transport.TransportError
	if errors.As(err, &e) && e.HttpCode == 404 {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return resp.Payload, nil
}
