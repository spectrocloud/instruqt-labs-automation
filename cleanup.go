package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spectrocloud/palette-sdk-go/api/client/version1"
	"github.com/spectrocloud/palette-sdk-go/client"
)

func Cleanup(SuperApiKey string) {
	SANDBOX_ID := os.Getenv("INSTRUQT_PARTICIPANT_ID")
	host := os.Getenv("PALETTE_HOST")

	apiKey := SuperApiKey
	projectUid := ""

	scope := "tenant"

	if host == "" || apiKey == "" {
		log.Printf("You must specify the PALETTE_HOST and PALETTE_API_KEY environment variables.")
		os.Exit(1)
	}
	if projectUid != "" {
		scope = "project"
	}

	// Initialize a Palette client
	fmt.Println(scope)

	pc := client.New(
		client.WithPaletteURI(host),
		client.WithAPIKey(apiKey),
	)
	if projectUid != "" {
		client.WithScopeProject(projectUid)(pc)
	} else {
		client.WithScopeTenant()(pc)
	}

	projectUid, err := pc.GetProjectUID(fmt.Sprintf("instruqt-%s", SANDBOX_ID))
	if err != nil {
		log.Printf("Failed to Get Project by UID: %s\n", fmt.Sprintf("instruqt-%s", SANDBOX_ID))
	}

	userId, err := pc.GetUserByEmail(fmt.Sprintf("instruqt+%s@spectrocloud.com", SANDBOX_ID))
	if err != nil {
		log.Printf("Failed to Get User by Email: %s\n", fmt.Sprintf("instruqt+%s@spectrocloud.com", SANDBOX_ID))
	}

	if err = pc.DeleteAPIKeyByName(fmt.Sprintf("instruqt-%s-api-key", SANDBOX_ID)); err != nil {
		log.Printf("Failed to Delete API Key: %s\n", fmt.Sprintf("instruqt-%s-api-key", SANDBOX_ID))
	}

	if err = pc.DeleteUser(userId.Metadata.UID); err != nil {
		log.Printf("Failed to Delete User: %s\n", userId.Metadata.UID)
	}

	client.WithScopeProject(projectUid)(pc)

	// Delete Project
	edgeClusters, err := pc.GetClusterGroupSummaries()

	for _, cluster := range edgeClusters {
		err = pc.ForceDeleteCluster(cluster.Metadata.UID, true)
		if err != nil {
			log.Printf("Failed to Force Delete Cluster: %s\n", cluster.Metadata.UID)
		}
	}

	// Delete Edge Devices
	edgeDevices, err := pc.ListEdgeHosts()
	for _, edgeDevice := range edgeDevices {
		err = pc.DeleteAppliance(edgeDevice.Metadata.UID)
		if err != nil {
			log.Printf("Failed to Delete Appliance: %s\n", edgeDevice.Metadata.UID)
		}
	}

	// Delete Registration Tokens
	params := &version1.V1EdgeTokensListParams{}
	edgeTokens, err := pc.Client.V1EdgeTokensList(params)

	for _, token := range edgeTokens.Payload.Items {
		if token.Spec.DefaultProject.UID == projectUid {
			err = pc.DeleteRegistrationToken(token.Metadata.UID)
			if err != nil {
				log.Printf("Failed to Delete Registration Token: %s\n", token.Metadata.UID)
			}
		}
	}

	// Delete Project Last
	if err = pc.DeleteProject(projectUid); err != nil {
		log.Printf("Failed to Delete Project: %s\n", projectUid)
	}

	fmt.Println("Deleted objects.")
}
