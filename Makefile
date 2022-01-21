.PHONY: buildwithvars installwithvars gorelease

static:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags='-extldflags "-f no-PIC -static"' -tags 'osusergo netgo static_build'

buildwithvars:
	rm -rf ./adise1941; go build -v -ldflags "-X main.commitHash=$$(git rev-parse --short HEAD) -X main.commitDate=$$(git log -1 --format=%ci | awk '{ print $$1 }')"

installwithvars:
	rm -rf ./adise1941; go install -v -ldflags "-X main.commitHash=$$(git rev-parse --short HEAD) -X main.commitDate=$$(git log -1 --format=%ci | awk '{ print $$1 }')"

gorelease:
	go install -v github.com/goreleaser/goreleaser@latest
	goreleaser --snapshot --skip-publish --rm-dist
