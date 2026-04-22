package util

import "os"

func IsGithubActionTestEnv() bool {
	return os.Getenv("TEST_ENV") == "github_action"
}
