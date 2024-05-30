package convert

import (
	"errors"
	"fmt"
	"regexp"
)

const basePath = "https://gitlab.example.org"

var ErrUnknownFormat = errors.New("unknown url format")

var regz []*regexp.Regexp = []*regexp.Regexp{
	// Matches this format:
	// git@gitlab.example.org:team/subteam/project.git
	regexp.MustCompile(`^git@gitlab.example.org:(?P<path>.+)\.git$`),
	// Matches these formats:
	// https://gitlab.example.org/team/project
	// https://gitlab.example.org/team/project.git
	regexp.MustCompile(`^https://gitlab.example.org/(?P<path>.+?)(\.git)?$`),
}

func ConvertGitRemoteToHTTP(remoteURL string) (string, error) {
	for _, re := range regz {
		parts := re.FindStringSubmatch(remoteURL)
		if len(parts) > 0 {
			pathIdx := re.SubexpIndex("path")
			path := parts[pathIdx]
			return fmt.Sprintf("%s/%s", basePath, path), nil
		}
	}

	return "", ErrUnknownFormat
}
