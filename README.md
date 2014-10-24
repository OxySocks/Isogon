IsoGOn
=========

IsoGOn is a web application/API written in Go used for (simple) domotica. It provides a minimal API internet connected devices to register their current state, and in the future to edit this state.
It is currently under development; and is mainly a 'toy' project (for now). It's primary goals are speed and simplicity in a minimal form-factor.

IsoGOn uses the [Martini] go web framework.

Written in  Go, IsoGOn is fast and can be run under most operating systems including Microsoft Windows, Mac OSX and Linux.

Installation
===============

Note: By default the HTTP server starts on port 3000. This can changed by declaring PORT environment variable.

    Install Go
    go get github.com/tools/godep
    git clone github.com/oxysocks/isogon
    cd isogon && godep get ./ && godep go build
    Start IsoGOn: PORT="80" MARTINI_ENV="production" ./isogon

Environment variables
===============

* PORT: Declares the port that should be used.
* MARTINI_ENV - used by Martini to enable production optimizations.

Support
=========

Currently no real support is available. However, pull request are more than welcome. A goal for IsoGOn is to write clean understandable code with proper documentation.

License
----

MIT

[Martini]:https://github.com/go-martini/martini