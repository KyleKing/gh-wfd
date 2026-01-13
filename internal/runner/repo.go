package runner

import (
	"fmt"

	"github.com/cli/go-gh/v2/pkg/repository"
)

// DetectRepo returns the current repository in "owner/repo" format.
func DetectRepo() (string, error) {
	repo, err := repository.Current()
	if err != nil {
		return "", fmt.Errorf("failed to detect repository: %w", err)
	}

	return fmt.Sprintf("%s/%s", repo.Owner, repo.Name), nil
}
