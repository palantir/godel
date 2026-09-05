<p align="right">
<a href="https://autorelease.general.dmz.palantir.tech/palantir/witchcraft-go-error"><img src="https://img.shields.io/badge/Perform%20an-Autorelease-success.svg" alt="Autorelease"></a>
</p>

witchcraft-go-error
===================
[![](https://godoc.org/github.com/palantir/witchcraft-go-error?status.svg)](http://godoc.org/github.com/palantir/witchcraft-go-error)

`witchcraft-error-go` defines the `werror` package, which provides an implementation of the `error` interface that
stores safe and unsafe parameters and has the ability to specify another error as a cause.

Associating structured safe and unsafe parameters with an error allows other infrastructure such as logging to make
decisions about what parameters should or should not be extricated.

TODO:
* Provide example usage and output in README

License
-------
This project is made available under the [Apache 2.0 License](http://www.apache.org/licenses/LICENSE-2.0).
