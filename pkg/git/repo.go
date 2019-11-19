package git

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-billy.v4/memfs"
	gitclient "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	gittransport "gopkg.in/src-d/go-git.v4/plumbing/transport/client"
	githttp "gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

var (
	// ErrRepoIsNotCloned is returned when an operation depends on a cloned repo
	ErrRepoIsNotCloned = errors.New("repo is not cloned")
)

const (
	authorName  = "kruise-deploy-wizard"
	authorEmail = "kruise@mastercard.com"
)

// Repo wraps a git repository
type Repo struct {
	repoURL string
	creds   *RepoCreds
	prefix  string
	ref     string
	fs      billy.Filesystem
	r       *gitclient.Repository
	files   map[string]string
}

// RepoCreds contains https basic authentication for a git repo
type RepoCreds struct {
	Username string
	Password string
}

// NewRepo creates a new instance of a repo
func NewRepo(repoURL, prefix, ref string, creds *RepoCreds, insecureSkipVerify bool) *Repo {
	customClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: insecureSkipVerify},
		},
	}

	gittransport.InstallProtocol("https", githttp.NewClient(customClient))

	return &Repo{
		repoURL: repoURL,
		creds:   creds,
		prefix:  prefix,
		ref:     ref,
		fs:      memfs.New(),
		files:   map[string]string{},
	}
}

// Clone clones the repository
func (r *Repo) Clone() error {
	gitrepo, err := gitclient.Clone(memory.NewStorage(), r.fs, &gitclient.CloneOptions{
		URL: r.repoURL,
		Auth: &githttp.BasicAuth{
			Username: r.creds.Username,
			Password: r.creds.Password,
		},
	})
	if err != nil {
		return err
	}
	r.r = gitrepo

	return nil
}

// AddFile adds a file to the Repo
func (r *Repo) AddFile(fileName string, content string) {
	r.files[fileName] = content
}

// AddDeploySpec adds a deploy spec to the repo
func (r *Repo) AddDeploySpec(fileName string, content string) error {
	if r.r == nil {
		return ErrRepoIsNotCloned
	}
	wt, err := r.r.Worktree()
	if err != nil {
		return err
	}

	err = r.fs.MkdirAll(filepath.Dir(fileName), 0755)
	if err != nil {
		return err
	}

	f, err := r.fs.Create(fileName)
	if err != nil {
		return err
	}

	_, err = io.WriteString(f, content)
	if err != nil {
		return err
	}

	err = f.Close()
	if err != nil {
		return err
	}

	if _, err := wt.Add(fileName); err != nil {
		return err
	}

	return nil
}

// Commit commits the current state of the repo
func (r *Repo) Commit(msg string) error {
	if r.r == nil {
		return ErrRepoIsNotCloned
	}
	wt, err := r.r.Worktree()
	if err != nil {
		return err
	}

	// write all files to the in-memory filesystem
	for name, content := range r.files {
		filename := filepath.Join(strings.TrimPrefix(r.prefix, "/"), name)

		err = r.fs.MkdirAll(filepath.Dir(filename), 0755)
		if err != nil {
			return err
		}

		f, err := r.fs.Create(filename)
		if err != nil {
			return err
		}

		_, err = io.WriteString(f, content)
		if err != nil {
			return err
		}

		err = f.Close()
		if err != nil {
			return err
		}

		if _, err := wt.Add(filename); err != nil {
			return err
		}
	}

	// commit the added files
	_, err = wt.Commit(
		msg,
		&gitclient.CommitOptions{
			Author: &object.Signature{
				Name: authorName, Email: authorEmail, When: time.Now(),
			},
		})
	if err != nil {
		return err
	}

	return nil
}

// Push pushes the repository
func (r *Repo) Push() error {
	if err := r.r.Push(&gitclient.PushOptions{
		Auth: &githttp.BasicAuth{
			Username: r.creds.Username,
			Password: r.creds.Password,
		},
	}); err != nil {
		return err
	}
	return nil
}

// Log outputs log entries of the repo TODO: remove me? just useful in testing
func (r *Repo) Log() error {
	ref, err := r.r.Head()
	if err != nil {
		return err
	}

	iter, err := r.r.Log(&gitclient.LogOptions{From: ref.Hash()})
	if err != nil {
		return err
	}

	// ... just iterates over the commits, printing it
	return iter.ForEach(func(c *object.Commit) error {
		fmt.Println(c)
		return nil
	})

	return nil
}
