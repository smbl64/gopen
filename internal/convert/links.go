package convert

import (
	"fmt"
	"regexp"
)

var regz []*regexp.Regexp = []*regexp.Regexp{
	// Matches this format:
	// git@gitlab.example.org:team/subteam/project.git
	regexp.MustCompile(`^git@(?P<base>.+):(?P<path>.+)\.git$`),
	// Matches these formats:
	// https://gitlab.example.org/team/project
	// https://gitlab.example.org/team/project.git
	regexp.MustCompile(`^https://(?P<base>.+)/(?P<path>.+?)(\.git)?$`),
}

func ConvertGitRemoteToHTTP(remoteURL string) (string, error) {
	for _, re := range regz {
		parts := re.FindStringSubmatch(remoteURL)
		if len(parts) > 0 {
			pathIdx := re.SubexpIndex("path")
			path := parts[pathIdx]

			pathIdx = re.SubexpIndex("base")
			base := parts[pathIdx]
			return fmt.Sprintf("https://%s/%s", base, path), nil
		}
	}

	return "", fmt.Errorf("unknown remote url format: '%s'", remoteURL)
}
