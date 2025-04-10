package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spectrocloud/palette-sdk-go/client"
)

func Cleanup() {
	SANDBOX_ID := os.Getenv("INSTRUQT_PARTICIPANT_ID")
	host := os.Getenv("PALETTE_HOST")
	apiKey := os.Getenv("PALETTE_API_KEY")
	projectUid := os.Getenv("PALETTE_PROJECT_UID")

	host = "console.vertex.spectrocloud.com"
	apiKey = "NTk4NGE3MjU4MTM0ZjQ1NTE5YzgwMDMwMmM0MTlmMWQ="
	projectUid = ""

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

	if err = pc.DeleteProject(projectUid); err != nil {
		log.Printf("Failed to Delete Project: %s\n", projectUid)
	}

	if err = pc.DeleteUser(userId.Metadata.UID); err != nil {
		log.Printf("Failed to Delete User: %s\n", userId.Metadata.UID)
	}

	fmt.Println("Deleted objects.")
}
