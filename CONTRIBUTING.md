# Contributing

Before starting your work on this module, you should
[fork the project](https://help.github.com/articles/fork-a-repo/) to your
GitHub account. This allows you to freely experiment with your changes. When
your changes are complete, submit a
[pull request](https://help.github.com/articles/using-pull-requests/). All pull
requests will be reviewed and merged if they suit some general guidelines:

* Changes are located in a topic branch
* For new functionality, proper tests are written
* Changes should not solve certain problems on special environments

## Issue Tracker
If you think your changes need further discussion, you found a bug or have a
feature request, create an issue on [GitHub](https://github.com/Icinga/icingabeat/issues). The issue number should
be used in your branch name as reference to a certain bug or feature discussion.

## Branches
Choosing a proper name for a branch helps us identify its purpose and possibly
find an associated bug or feature. Generally a branch name should include a
topic such as `bug` or `feature` followed by a description and an issue number
if applicable. Branches should have only changes relevant to a specific issue.

```
git checkout -b bug/field-type-missing-1234
git checkout -b feature/new-setting-for-foo-1235
```
