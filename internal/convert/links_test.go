package convert

import "testing"

func TestConvertGitLink(t *testing.T) {
	testData := []struct {
		Name      string
		Input     string
		Expected  string
		MustError bool
	}{
		{
			Name:     "ssh format",
			Input:    "git@gitlab.example.org:team/subteam/project.git",
			Expected: "https://gitlab.example.org/team/subteam/project",
		},
		{
			Name:     "https format",
			Input:    "https://gitlab.example.org/team/project.git",
			Expected: "https://gitlab.example.org/team/project",
		},
		{
			Name:      "unknown address",
			Input:     "https://www.google.com",
			MustError: true,
		},
	}

	for _, testcase := range testData {
		t.Run(testcase.Name, func(t *testing.T) {
			got, err := ConvertGitRemoteToHTTP(testcase.Input)
			if !testcase.MustError && err != nil {
				t.Fatalf("failed for input '%s': %v", testcase.Input, err)
			} else if testcase.MustError && err == nil {
				t.Fatalf("failed for input '%s': expected an error but got nil", testcase.Input)
			}

			if got != testcase.Expected {
				t.Fatalf("input '%s'\nwanted '%s'\ngot '%s'", testcase.Input, testcase.Expected, got)
			}
		})
	}
}
