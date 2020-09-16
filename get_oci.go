package getter

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	auth "github.com/deislabs/oras/pkg/auth/docker"
	"github.com/deislabs/oras/pkg/content"
	"github.com/deislabs/oras/pkg/oras"
)

// OCIGetter is responsible for handling OCI repositories
type OCIGetter struct {
	getter
}

// ClientMode returns the client mode directory
func (g *OCIGetter) ClientMode(u *url.URL) (ClientMode, error) {
	return ClientModeDir, nil
}

// Get gets the repository as the specified url
func (g *OCIGetter) Get(path string, u *url.URL) error {
	ctx := g.Context()

	if !pathContainsTag(u.Path) {
		u.Path = u.Path + ":latest"
	}

	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return fmt.Errorf("make policy directory: %w", err)
	}

	cli, err := auth.NewClient()
	if err != nil {
		return fmt.Errorf("new auth client: %w", err)
	}

	resolver, err := cli.Resolver(ctx, http.DefaultClient, false)
	if err != nil {
		return fmt.Errorf("new resolver: %w", err)
	}

	fileStore := content.NewFileStore(path)
	defer fileStore.Close()

	repository := u.Host + u.Path
	_, _, err = oras.Pull(ctx, resolver, repository, fileStore)
	if err != nil {
		return fmt.Errorf("pulling policy: %w", err)
	}

	return nil
}

// GetFile is currently a NOOP
func (g *OCIGetter) GetFile(dst string, u *url.URL) error {
	return nil
}