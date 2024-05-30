[private]
default:
  @just --list

test:
  go test -v ./...
