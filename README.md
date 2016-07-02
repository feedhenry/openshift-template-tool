# OpenShift template tool

This tool is meant to facilitate working with OpenShift templates.

Initially, the tool is able to combine multiple templates into a single unified
template.

It can be used in a workflow in which you maintain several small templates for
individual components and then automatically generate an all-in-one template to
describe a more complex application.

It can also be used as a way to pretty-print templates in a given API version.

See the next sections for details on how to build, test and use the tool.

## Building

You will need Go 1.5 or later. If using Go 1.5, set the environment variable
`GO15VENDOREXPERIMENT=1`. Later versions don't need that environment variable
and have vendoring support enabled by default.

### Using `go get`

```
go get github.com/feedhenry/openshift-template-tool
```

### Alternative

Clone this project into a path defined in your `GOPATH` environment variable.

From the project root where this `README.md` is located, build using the
standard `go` tool:

```
cd $GOPATH/src/github.com/feedhenry/openshift-template-tool
go install
```

## Testing

Run tests using the standard `go` tool:

```
go test ./template/
```

## Using

It's suggested to add `$GOPATH/bin` to your `PATH` environment variable, to make
it easier to call the tool by it's name.

That can be accomplished by adding this to your `~/.bashrc`:

```
# Include every ./bin directory from GOPATHs into PATH
export PATH=${GOPATH//://bin:}/bin:$PATH
```

Then, the tool is available as a regular command:

```
openshift-template-tool --help

openshift-template-tool merge base-template.json \
                              component-1.json \
                              component-2.json \
                              ... \
                              > merged-template.json
```

The output of `openshift-template-tool merge` is a template object that inherits
the metadata from the base template, the first command line argument. The lists
of objects and parameters from each of the subsequent templates are appended to
that of the base template and duplicates are removed.

## Release procedure

1) Modify VERSION file and commit all changes

```
vi VERSION
git commit -a -m"Version bump"
git push upstream master
```

2) Execute build script

`./scripts/release`

3) Create release in github and attach executables from dist folder.

https://github.com/feedhenry/openshift-template-tool/releases/new

## Updating vendored dependencies

This tool depends heavily on
[openshift/origin](https://github.com/openshift/origin). Origin itself depends
on several other packages. Some of those packages have patches that don't exist
upstream, code that lives only vendored within the Origin code base.

None of the more popular vendoring tools in the Go ecosystem can handle
re-vendoring the packages vendored in Origin. We tried using
[godep](https://github.com/tools/godep) and
[glide](https://github.com/Masterminds/glide), solved error after error, but
arrived at no satisfying result.

In essence, what we need is a way to copy all transitive dependencies into a
`vendor/` subdirectory. With that in mind, we have a script to automate the
process.

To update the vendored dependencies:

0. Start from a clean state. Use `git status` to ensure you have a clean
   working tree.
1. If you want to base the dependencies on a different version of Origin, edit
   the `copy-dependencies` script and set the value of `ORIGIN_TAG` to point to
   the appropriate tag/branch/commit.
2. Run `./copy-dependencies`.
3. Use `git status` to verify what changes were made.
4. Make sure you can still build the code with `go install`.
5. Once you are happy with the changes, run `git add copy-dependencies vendor`
   to add the changes to the index, then `git commit`.

## License

The OpenShift template tool is licensed under the [Apache License, Version 2.0](http://www.apache.org/licenses/).
