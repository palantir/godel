package fail

fail

/*
This is a non-compiling file that has been added to explicitly ensure that CI fails.
It also contains the command that caused the failure and its output.
Remove this file if debugging locally.

./godelw verify failed after updating godel plugins and assets

Command that caused error:
./godelw exec -- go fix ./. ./docsgenerator ./docsgenerator/cmd ./docsgenerator/generator ./framework/artifactresolver ./framework/builtintasks ./framework/builtintasks/githooks ./framework/builtintasks/githubwiki ./framework/builtintasks/idea ./framework/builtintasks/installupdate ./framework/builtintasks/installupdate/layout ./framework/builtintasks/packages ./framework/godel ./framework/godel/config ./framework/godel/config/internal/v0 ./framework/godellauncher ./framework/godellauncher/defaulttasks ./framework/internal/pathsinternal ./framework/internal/pluginsinternal ./framework/pluginapi ./framework/pluginapi/v2/pluginapi ./framework/pluginapitester ./framework/plugins ./framework/verifyorder ./godelgetter ./godelinit ./godelinit/cmd ./integration_test ./pkg/dirchecksum ./pkg/osarch ./pkg/versionedconfig

Output:
stat /repo/pkg/products/v2/docsgenerator: directory not found
stat /repo/pkg/products/v2/docsgenerator/cmd: directory not found
stat /repo/pkg/products/v2/docsgenerator/generator: directory not found
stat /repo/pkg/products/v2/framework/artifactresolver: directory not found
stat /repo/pkg/products/v2/framework/builtintasks: directory not found
stat /repo/pkg/products/v2/framework/builtintasks/githooks: directory not found
stat /repo/pkg/products/v2/framework/builtintasks/githubwiki: directory not found
stat /repo/pkg/products/v2/framework/builtintasks/idea: directory not found
stat /repo/pkg/products/v2/framework/builtintasks/installupdate: directory not found
stat /repo/pkg/products/v2/framework/builtintasks/installupdate/layout: directory not found
stat /repo/pkg/products/v2/framework/builtintasks/packages: directory not found
stat /repo/pkg/products/v2/framework/godel: directory not found
stat /repo/pkg/products/v2/framework/godel/config: directory not found
stat /repo/pkg/products/v2/framework/godel/config/internal/v0: directory not found
stat /repo/pkg/products/v2/framework/godellauncher: directory not found
stat /repo/pkg/products/v2/framework/godellauncher/defaulttasks: directory not found
stat /repo/pkg/products/v2/framework/internal/pathsinternal: directory not found
stat /repo/pkg/products/v2/framework/internal/pluginsinternal: directory not found
stat /repo/pkg/products/v2/framework/pluginapi: directory not found
stat /repo/pkg/products/v2/framework/pluginapi/v2/pluginapi: directory not found
stat /repo/pkg/products/v2/framework/pluginapitester: directory not found
stat /repo/pkg/products/v2/framework/plugins: directory not found
stat /repo/pkg/products/v2/framework/verifyorder: directory not found
stat /repo/pkg/products/v2/godelgetter: directory not found
stat /repo/pkg/products/v2/godelinit: directory not found
stat /repo/pkg/products/v2/godelinit/cmd: directory not found
stat /repo/pkg/products/v2/integration_test: directory not found
stat /repo/pkg/products/v2/pkg/dirchecksum: directory not found
stat /repo/pkg/products/v2/pkg/osarch: directory not found
stat /repo/pkg/products/v2/pkg/versionedconfig: directory not found
Error: exit status 1

*/
