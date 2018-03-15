distgo
======
distgo is a gödel plugin that runs, builds, distributes and publishes products in a Go project based on declarative
project configuration.

Plugin Tasks
------------
distgo provides the following tasks:

* `artifacts`: prints the artifacts (build, dist or Docker) for the specified products.
* `build`: builds the executables for the specified products.
* `clean`: removes the outputs (build, dist and Docker) generated for the specified products.
* `dist`: creates the distribution outputs for the specified products.
* `docker`: creates the Docker images for the specified products.
* `products`: prints all of the products for the project.
* `project-version`: prints the version of the project.
* `publish`: publishes the distribution artifacts for the specified products.
* `run`: runs the build output for the specified product.

Assets
------
distgo assets are executables that provide specific functionality for distgo. Assets can provide distribution actions
for the "dist" task, Docker build actions for the "docker" task and publish actions for the "publish" task. Refer to
the `assetapi`, `dister`, `dockerbuilder` and `publisher` packages for more information.

Writing an asset
----------------
distgo provides helper APIs to facilitate writing new assets. More detailed instructions for writing assets are
forthcoming. In the meantime, the most effective way to write an asset is to examine the implementation of an existing
asset.

Core concepts
-------------
distgo operates on a single project. A project is a logical unit that groups code, and is typically a GitHub repository.
A single project may contain one or more products, where a product is a Go "main" package. Users specify the products in
a project using a YAML configuration file.

A product may specify other products as dependencies of the product. Dependencies cannot form a cycle.

A product typically has a build configuration. The build configuration specifies things such as the target
OS/architectures, build flags that should be used when building the project, etc. A product can only have a single build
configuration, and the number of outputs for the build is equivalent to the number of target OS/architectures specified.

A product can have one or more dist configurations. A dist configuration specifies a distribution type and then creates
one or more files as output when run. When a dist task for a product is run, the output of the product's (and any of its
dependencies') build task is available.

A product can have one or more docker configurations. A Docker configuration specifies a Docker builder type and its
parameters, along with an input Docker context directory and Dockerfile. When a Docker builder for a product is run, the
output of the product's (and any of its dependencies') dist task is available.

The "publish" operation for a product specifies one or more of a product's distributions. distgo ships with some
built-in publish operations and supports customizing publish operations using assets. A publish operation takes the dist
outputs for a product dist configuration and publishes them.
