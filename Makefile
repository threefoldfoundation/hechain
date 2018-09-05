all: install

daemonpkgs = ./cmd/hechaind
clientpkgs = ./cmd/hechainc
testpkgs = ./pkg/types
pkgs = $(daemonpkgs) $(clientpkgs) ./pkg/config $(testpkgs)

version = $(shell git describe --abbrev=0)
commit = $(shell git rev-parse --short HEAD)
ifeq ($(commit), $(shell git rev-list -n 1 $(version) | cut -c1-7))
	fullversion = $(version)
	fullversionpath = \/releases\/tag\/$(version)
else
	fullversion = $(version)-$(commit)
	fullversionpath = \/tree\/$(commit)
endif

dockerVersion = $(shell git describe --abbrev=0 | cut -d 'v' -f 2)
dockerVersionEdge = edge

configpkg = github.com/threefoldfoundation/hechain/pkg/config
ldflagsversion = -X $(configpkg).rawVersion=$(fullversion)

stdoutput = $(GOPATH)/bin
daemonbin = $(stdoutput)/hechaind
clientbin = $(stdoutput)/hechainc

install:
	go build -race -tags='debug profile' -ldflags '$(ldflagsversion)' -o $(daemonbin) $(daemonpkgs)
	go build -race -tags='debug profile' -ldflags '$(ldflagsversion)' -o $(clientbin) $(clientpkgs)

install-std:
	go build -ldflags '$(ldflagsversion) -s -w' -o $(daemonbin) $(daemonpkgs)
	go build -ldflags '$(ldflagsversion) -s -w' -o $(clientbin) $(clientpkgs)

test:
	go test -race -v -tags='debug testing' -timeout=60s $(testpkgs)

test-coverage:
	go test -race -v -tags='debug testing' -timeout=60s \
		-coverpkg=all -coverprofile=coverage.out -covermode=atomic $(testpkgs)

test-coverage-web: test-coverage
	go tool cover -html=coverage.out

# xc builds and packages release binaries
# for all windows, linux and mac, 64-bit only,
# using the standard Golang toolchain.
xc:
	docker build -t hechainbuilder -f DockerBuilder .
	docker run --rm -v $(shell pwd):/go/src/github.com/threefoldfoundation/hechain hechainbuilder

docker-minimal: xc
	docker build -t hechain/hechain:$(dockerVersion) -f DockerfileMinimal --build-arg binaries_location=release/hechain-$(version)-linux-amd64/cmd .

# Release images builds and packages release binaries, and uses the linux based binary to create a minimal docker
release-images: get_hub_jwt docker-minimal
	docker push hechain/hechain:$(dockerVersion)
	# also create a latest
	docker tag hechain/hechain:$(dockerVersion) hechain/hechain
	docker push hechain/hechain:latest
	curl -b "active-user=hechain; caddyoauth=$(HUB_JWT)" -X POST --data "image=hechain/hechain:$(dockerVersion)" "https://hub.gig.tech/api/flist/me/docker"
	# symlink the latest flist
	curl -b "active-user=hechain; caddyoauth=$(HUB_JWT)" -X GET "https://hub.gig.tech/api/flist/me/hechain-hechain-$(dockerVersion).flist/link/hechain-hechain-latest.flist"
	# Merge the flist with ubuntu and nmap flist, so we have a tty file etc...
	curl -b "active-user=hechain; caddyoauth=$(HUB_JWT)" -X POST --data "[\"gig-official-apps/ubuntu1604.flist\", \"hechain/hechain-hechain-$(dockerVersion).flist\", \"gig-official-apps/nmap.flist\"]" "https://hub.gig.tech/api/flist/me/merge/hechain-16.04-hechain-$(dockerVersion).flist"
	# And also link in a latest
	curl -b "active-user=hechain; caddyoauth=$(HUB_JWT)" -X GET "https://hub.gig.tech/api/flist/me/ubuntu-16.04-hechain-$(dockerVersion).flist/link/ubuntu-16.04-hechain-latest.flist"

xc-edge:
	docker build -t hechainbuilderedge -f DockerBuilderEdge .
	docker run --rm -v $(shell pwd):/go/src/github.com/threefoldfoundation/hechain hechainbuilderedge

docker-minimal-edge: xc-edge
	docker build -t hechain/hechain:$(dockerVersionEdge) -f DockerfileMinimal --build-arg binaries_location=release/hechain-$(dockerVersionEdge)-linux-amd64/cmd .

# Release images builds and packages release binaries, and uses the linux based binary to create a minimal docker
release-images-edge: get_hub_jwt docker-minimal-edge
	docker push hechain/hechain:$(dockerVersionEdge)
	curl -b "active-user=hechain; caddyoauth=$(HUB_JWT)" -X POST --data "image=hechain/hechain:$(dockerVersionEdge)" "https://hub.gig.tech/api/flist/me/docker"
	# Merge the flist with ubuntu and nmap flist, so we have a tty file etc...
	curl -b "active-user=hechain; caddyoauth=$(HUB_JWT)" -X POST --data "[\"gig-official-apps/ubuntu1604.flist\", \"hechain/hechain-hechain-$(dockerVersionEdge).flist\", \"gig-official-apps/nmap.flist\"]" "https://hub.gig.tech/api/flist/me/merge/ubuntu-16.04-hechain-$(dockerVersionEdge).flist"

release-dir:
	[ -d release ] || mkdir release

get_hub_jwt: check-HUB_APP_ID check-HUB_APP_SECRET
	$(eval HUB_JWT = $(shell curl -X POST "https://itsyou.online/v1/oauth/access_token?response_type=id_token&grant_type=client_credentials&client_id=$(HUB_APP_ID)&client_secret=$(HUB_APP_SECRET)&scope=user:memberof:hechain"))

check-%:
	@ if [ "${${*}}" = "" ]; then \
		echo "Required env var $* not present"; \
		exit 1; \
	fi

ineffassign:
	ineffassign $(pkgs)

.PHONY: all install xc release-images get_hub_jwt check-% ineffassign