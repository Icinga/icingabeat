# Release Workflow
Before submitting a new release, make sure all relevant pull requests and local branches have been merged to the `master`
branch. All tests must pass before a release is tagged.

## 1. AUTHORS
Update the [AUTHORS] and [.mailmap] file

``` bash
git checkout master
git log --use-mailmap | grep ^Author: | cut -f2- -d' ' | sort | uniq > AUTHORS
git commit -am "Update AUTHORS"
```

## 2. Changelog
Update [CHANGELOG.md] with all relevant information.

## 3. Version
Version numbers are incremented regarding the [SemVer 1.0.0] specification. 
Update the version number in the following files:

* `version.yml`
* `vendor/github.com/elastic/beats/dev-tools/packer/version.yml`

## 4. Build
Build packages:

``` bash
export SNAPSHOT=false
make package
```

Create dashboard zip files:

``` bash
export SNAPSHOT=false
make package-dashboards
```

## 5. Git Tag
Commit all changes to the `master` branch

``` bash
git commit -v -a -m "Release version <VERSION>"
git push
```

Tag the release

``` bash
git tag -m "Version <VERSION>" v<VERSION>
```

Push tags

``` bash
git push --tags
```

[SemVer 1.0.0]: http://semver.org/spec/v1.0.0.html
[CHANGELOG.md]: CHANGELOG.md
[AUTHORS]: AUTHORS
[.mailmap]: .mailmap
