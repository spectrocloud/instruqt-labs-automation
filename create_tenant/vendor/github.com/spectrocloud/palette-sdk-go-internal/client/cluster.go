package client

import (
	"github.com/spectrocloud/hapi/models"
	clusterC "github.com/spectrocloud/hapi/spectrocluster/client/v1"
)

func (h *V1Client) UpdatePauseAgentUpgradeSettingCluster(upgradeSetting *models.V1ClusterUpgradeSettingsEntity, clusterUID string, context string) error {
	params := clusterC.NewV1SpectroClustersUIDUpgradeSettingsParamsWithContext(h.ctx).
		WithUID(clusterUID).
		WithBody(upgradeSetting)
	_, err := h.ClusterC.V1SpectroClustersUIDUpgradeSettings(params)
	if err != nil {
		return err
	}
	return nil
}

func (h *V1Client) UpdatePauseAgentUpgradeSettingContext(upgradeSetting *models.V1ClusterUpgradeSettingsEntity, context string) error {
	params := clusterC.NewV1SpectroClustersUpgradeSettingsParamsWithContext(h.ctx).
		WithBody(upgradeSetting)
	_, err := h.ClusterC.V1SpectroClustersUpgradeSettings(params)
	return err
}

func (h *V1Client) GetPauseAgentUpgradeSettingContext(context string) (string, error) {
	params := clusterC.NewV1SpectroClustersUpgradeSettingsGetParamsWithContext(h.ctx)
	resp, err := h.ClusterC.V1SpectroClustersUpgradeSettingsGet(params)
	if err != nil {
		return "", err
	}
	return resp.Payload.SpectroComponents, nil
}
