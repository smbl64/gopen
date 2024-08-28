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
	Verbose         bool   `arg:"-v,--verbose" help:"Verbose output"`
	NoSymlinkFollow bool   `arg:"--no-follow" help:"Do not follow symbolic links"`
	NoBrowserOpen   bool   `arg:"-n,--no-open" help:"Do not open web browser"`
	Target          string `arg:"positional,required" help:"The file to open in web browser"`
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
	logger.Debugf("Found git root folder: %s\n", root)
	logger.Debugf("Relative file path: %s\n", relative)

	remoteURL, err := readGitRemoteURL(root)
	if err != nil {
		die(err)
	}
	logger.Debugf("Found git remote URL: %s\n", remoteURL)

	remoteURL, err = convert.ConvertGitRemoteToHTTP(remoteURL)
	if err != nil {
		die(err)
	}
	logger.Debugf("Remote URL as https: %s\n", remoteURL)

	branch, err := getDefaultBranchName(root)
	if err != nil {
		die(err)
	}
	logger.Debugf("Default git branch name: %s\n", branch)

	url := makeFinalURL(remoteURL, branch, relative)
	logger.Debugf("")
	logger.Debugf("Destination URL: %s\n", url)

	if !opts.NoBrowserOpen {
		err = openBrowser(url)
		if err != nil {
			die(err)
		}
	}
}

func makeFinalURL(remoteURL, branch, relativeFileName string) string {
	if strings.Contains(remoteURL, "github") {
		return fmt.Sprintf("%s/blob/%s/%s", remoteURL, branch, relativeFileName)
	} else {
		// Gitlab style
		return fmt.Sprintf("%s/-/blob/%s/%s", remoteURL, branch, relativeFileName)
	}
}

func findRepoAndRelativePath(target string) (string, string, error) {
	var err error
	if !opts.NoSymlinkFollow {
		logger.Debugf("Will follow symbolic links if any")
		target, err = filepath.EvalSymlinks(target)
		if err != nil {
			return "", "", err
		}
	}

	target, err = filepath.Abs(target)
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
	var result string
	cmd := exec.Command("git", "-C", repoDir, "symbolic-ref", "refs/remotes/origin/HEAD")
	output, err := cmd.CombinedOutput()
	if err == nil {
		result = string(output)
		result = strings.Trim(result, "\r\n")
		result, _ = strings.CutPrefix(result, "refs/remotes/origin/")
		return result, nil
	}

	// Try a few common cases if the 'git' command above failed
	candidates := []string{"master", "main", "trunk"}
	for _, c := range candidates {
		_, err := os.Stat(fmt.Sprintf(".git/refs/remotes/origin/%s", c))
		if err != nil {
			continue
		}

		return c, nil
	}

	return "", errors.New("cannot determine the branch name")
}

func openBrowser(url string) error {
	var err error

	switch runtime.GOOS {
	case "linux":
	case "freebsd":
	case "openbsd":
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
