BEATNAME=icingabeat
BEAT_DIR=github.com/icinga/icingabeat
BEAT_DESCRIPTION=Icingabeat ships Icinga 2 events and states to Elasticsearch or Logstash.
BEAT_VENDOR=Icinga
BEAT_DOC_URL=https://github.com/Icinga/icingabeat
SYSTEM_TESTS=false
TEST_ENVIRONMENT=false
ES_BEATS?=./vendor/github.com/elastic/beats
GOPACKAGES=$(shell glide novendor)
PREFIX?=.

#TARGETS="linux/amd64 linux/386 windows/amd64 windows/386 darwin/amd64"
#PACKAGES=${BEATNAME}/deb ${BEATNAME}/rpm ${BEATNAME}/darwin ${BEATNAME}/win ${BEATNAME}/bin
#SNAPSHOT=false

# Path to the libbeat Makefile
-include $(ES_BEATS)/libbeat/scripts/Makefile

# Initial beat setup
.PHONY: setup
setup: copy-vendor
	make update

# Copy beats into vendor directory
.PHONY: copy-vendor
copy-vendor:
	mkdir -p vendor/github.com/elastic/
	cp -R ${GOPATH}/src/github.com/elastic/beats vendor/github.com/elastic/
	rm -rf vendor/github.com/elastic/beats/.git

.PHONY: git-init
git-init:
	git init
	git add README.md CONTRIBUTING.md
	git commit -m "Initial commit"
	git add LICENSE
	git commit -m "Add the LICENSE"
	git add .gitignore
	git commit -m "Add git settings"
	git add .
	git reset -- .travis.yml
	git commit -m "Add icingabeat"
	git add .travis.yml
	git commit -m "Add Travis CI"

# This is called by the beats packer before building starts
.PHONY: before-build
before-build:

# Collects all dependencies and then calls update
.PHONY: collect
collect:
