# gopen 
Open Git-tracked Files in Your Web Browser

## Introduction

`gopen` is an open-source utility designed to simplify the process of opening files tracked by Git directly from your terminal or command line interface. Whether you're working with a repository hosted on
GitHub, GitLab, Bitbucket, or any other Git service that supports web-based access, `gopen` ensures you can quickly and easily navigate to the relevant file in your browser.

## Features

- **Cross-platform support**: Works seamlessly across Windows, macOS, and Linux.
- **Auto-detection of repository URL**: gopen automatically detects whether a local file path is part of a Git-managed project and fetches its corresponding remote URL.
- **Support for major Git platforms**: GitHub, GitLab, Bitbucket, and more!

## Installation

Before installing `gopen`, please ensure you have Go installed on your system. You can download it from the official Go website: [https://go.dev/doc/install](https://go.dev/doc/install).

### Installing gopen via Go

You can use the Go to install directly from GitHub:

```sh
go install github.com/smbl64/gopen@latest
```

## Usage

To open a file or folder in your web browser using `gopen`, simply run the following command, replacing `/path/to/file` with the actual path of the file you want to view:

```sh
gopen /path/to/file
```

If `gopen` detects that this file is part of a Git-managed project, it will open the corresponding URL in your configured browser. If not found or no repository detected, it will return an error.

### Example

To open a file named `main.go` located in your current directory:

```sh
gopen main.go
```

To open current folder:
```sh
gopen .
```

## Contributing

We welcome contributions! If you find any issues or have suggestions for improvements, please feel free to open an issue or submit a pull request.

## License

gopen is released under the MIT License. See the LICENSE file for more details.

