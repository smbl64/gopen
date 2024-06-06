package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/alexflint/go-arg"
	"github.com/smbl64/gopen/internal/convert"
	"github.com/smbl64/gopen/internal/log"
)

var opts struct {
	Verbose bool   `arg:"-v,--verbose" help:"Verbose output"`
	Target  string `arg:"positional,required" help:"The file to open in web browser"`
}

var logger *log.Logger

func die(err error) {
	logger.Infof("fatal: %v\n", err)
	os.Exit(1)
}

func main() {
	arg.MustParse(&opts)

	logger = log.NewLogger(opts.Verbose)

	root, relative, err := findRepoAndRelativePath(opts.Target)
	if err != nil {
		die(err)
	}

	remoteURL, err := readGitRemoteURL(root)
	if err != nil {
		die(err)
	}

	remoteURL, err = convert.ConvertGitRemoteToHTTP(remoteURL)
	if err != nil {
		die(err)
	}

	branch, err := getDefaultBranchName(root)
	if err != nil {
		die(err)
	}

	url := fmt.Sprintf("%s/-/blob/%s/%s", remoteURL, branch, relative)
	logger.Debugf("Root     : %s\n", root)
	logger.Debugf("Repo     : %s\n", remoteURL)
	logger.Debugf("Branch   : %s\n", branch)
	logger.Debugf("Relative : %s\n", relative)
	logger.Debugf("URL      : %s\n", url)

	err = openBrowser(url)
	if err != nil {
		die(err)
	}
}

func findRepoAndRelativePath(target string) (string, string, error) {
	target, err := filepath.Abs(target)
	if err != nil {
		return "", "", err
	}

	root, err := findGitRoot(target)
	if err != nil {
		return "", "", err
	}

	relative, err := filepath.Rel(root, target)
	if err != nil {
		return "", "", err
	}

	return root, relative, nil

}

func findGitRoot(filename string) (string, error) {
	separator := string(os.PathSeparator)
	parts := strings.Split(filename, separator)

	// Add a slash for non-Windows OS
	if runtime.GOOS != "windows" {
		parts = append([]string{separator}, parts...)
	}

	// Go upwards within the folder hierarchy and see if we
	// can find the CODEOWNERS file.
	for i := len(parts); i > 1; i-- {
		current := parts[0:i]

		root := filepath.Join(current...)
		found := hasDir(root, ".git")
		if found {
			return root, nil
		}
	}

	return "", errors.New("No git repository found")
}

func hasDir(root string, dirName string) bool {
	dir := filepath.Join(root, dirName)
	stat, err := os.Stat(dir)

	if err != nil {
		return false
	}
	return stat.IsDir()
}

func readGitRemoteURL(repoDir string) (string, error) {
	cmd := exec.Command("git", "-C", repoDir, "remote", "get-url", "origin")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git: %w", err)
	}

	url := string(output)
	return strings.Trim(url, "\r\n"), nil
}

func getDefaultBranchName(repoDir string) (string, error) {
	cmd := exec.Command("git", "-C", repoDir, "symbolic-ref", "refs/remotes/origin/HEAD")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git: %w", err)
	}

	outStr := string(output)
	outStr = strings.Trim(outStr, "\r\n")
	outStr, _ = strings.CutPrefix(outStr, "refs/remotes/origin/")

	return outStr, nil
}

func openBrowser(url string) error {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = errors.New("unsupported platform")
	}

	return err
}
