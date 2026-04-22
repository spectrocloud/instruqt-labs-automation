package client

import (
	"errors"
	"github.com/spectrocloud/hapi/apiutil/transport"
	"github.com/spectrocloud/hapi/models"
	clusterC "github.com/spectrocloud/hapi/spectrocluster/client/v1"
)

// GetServiceImagesCLI ..
func (h *V1Client) GetServiceImagesCLI(cloudType, version, clusterUID, edgeHostUID string, fips bool) ([]*models.V1ServiceImage, error) {
	params := clusterC.NewV1ServicesImagesGetParamsWithContext(h.ctx).
		WithCloudTypes(&cloudType).WithVersion(&version).WithIsFipsEnabled(&fips)

	if clusterUID != "" {
		params.SetClusterUID(&clusterUID)
	}

	if edgeHostUID != "" {
		params.SetEdgeHostUID(&edgeHostUID)
	}

	resp, err := h.ClusterC.V1ServicesImagesGet(params)
	var e *transport.TransportError
	if errors.As(err, &e) && e.HttpCode == 404 {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return resp.Payload.ServiceImages, nil
}
