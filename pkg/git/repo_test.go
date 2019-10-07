package git_test

import (
	"os"
	"testing"

	"deploy-wizard/pkg/git"
)

const (
	envRepoURLVar  = "KRUISE_GIT_REPO"
	envUsernameVar = "KRUISE_GIT_USERNAME"
	envPasswordVar = "KRUISE_GIT_PASSWORD"
)

func TestClone(t *testing.T) {
	configMissingMessage := "skipping repo tests because KRUISE_GIT_* variables are not set. Configure KRUISE_GIT_REPO, KRUISE_GIT_USERNAME, and KRUISE_GIT_PASSWORD"
	testRepo := os.Getenv(envRepoURLVar)
	if testRepo == "" {
		t.Skip(configMissingMessage)
	}
	username := os.Getenv(envUsernameVar)
	if username == "" {
		t.Skip(configMissingMessage)
	}
	password := os.Getenv(envPasswordVar)
	if password == "" {
		t.Skip(configMissingMessage)
	}

	repo := git.NewRepo(testRepo, "deploy", "HEAD", &git.RepoCreds{
		username,
		password,
	}, false)

	err := repo.Clone()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	repo.AddFile("test.txt", "Heyo world\n")

	err = repo.Commit("testing repo with test.txt")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	err = repo.Push()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	err = repo.Log()
	if err != nil {
		t.Error(err)
	}
}
