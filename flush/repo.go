package flush

import "os/exec"

type RepoWrapper interface {
	GetDiff() (string, error)
	Commit(string) (string, error)
}

type GitWrapper struct{}

func InitRepoWrapper() RepoWrapper {
	return &GitWrapper{}
}

func (w *GitWrapper) GetDiff() (string, error) {
	cmd := exec.Command("git", "diff", "--cached")
	res, err := cmd.Output()
	if err != nil {
		logger.Error("Failed to read 'git diff --cached': ", string(res))
	}
	return string(res), err
}

func (w *GitWrapper) Commit(message string) (string, error) {
	cmd := exec.Command("git", "commit", "-m", message)
	res, err := cmd.Output()
	if err != nil {
		logger.Error("Failed to commit to local repository: ", string(res))
	}
	return string(res), err
}
