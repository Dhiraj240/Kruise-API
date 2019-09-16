package git_test

import (
	"os"
	"testing"

	"deploy-wizard/pkg/git"
)

const testRepo = "https://github.com/ryane/sampleapp.git"

const (
	envUsernameVar = "KRUISE_STASH_USERNAME"
	envPasswordVar = "KRUISE_STASH_PASSWORD"
)

func TestClone(t *testing.T) {
	username := os.Getenv(envUsernameVar)
	if username == "" {
		t.Errorf("set a valid username for the test repo in a an environment variable called %s", envUsernameVar)
		t.FailNow()
	}
	password := os.Getenv(envPasswordVar)
	if password == "" {
		t.Errorf("set a valid password for the test repo in a an environment variable called %s", envPasswordVar)
		t.FailNow()
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
