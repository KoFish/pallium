pallium
=======

Pallium is an experimental golang homeserver implementation for the [Matrix.org](matrix.org) standard.

Matrix is a new open standard for interoperable Instant Messaging and VoIP, providing pragmatic
HTTP APIs and open source reference implementations for creating and running your own real-time
communication infrastructure. 

To get up and running:

    # get the latest Go from https://golang.org/dl/
    # and define a GOPATH for your Go workspace if you don't have one already:
    export GOPATH=~/go
    mkdir -p $GOPATH

    # grab the latest pallium with all its dependencies:
    go get github.com/KoFish/pallium

    # set up a default config
    cd $GOPATH/github.com/KoFish/pallium
    cp config.json.dist config.json

    # edit the hostname param in config.json to specify how the server
    # should refer to itself and expect to be accessed from the internet
    # (use localhost for local experimentation)

    # set the server running. This will create a local sqlite db for storage
    # and start listening for traffic.
    $GOPATH/bin/pallium

To use the server, select a client from [matrix.org](matrix.org).  For instance, to use
the webclient hosted at https://matrix.org/beta against your new server, just specify
the URL of your pallium server (e.g. http://localhost:8008) as the "Home Server" parameter
on the login and registration pages rather than http://matrix.org.  You do not need to
enter the captcha, as pallium does not support captcha-based registration yet.

*Currently pallium does not support the full Matrix API set, so the webclient may well not
work correctly*

Alternatively, to run your own webclient:

    git clone http://github.com/matrix-org/matrix-angular-sdk

...and follow the instructions in syweb/webclient/README

For more information on Matrix, please see [Matrix.org](matrix.org), the
[Matrix Specification](github.com/matrix-org/matrix-doc/tree/master/specification) or
[Synapse](github.com/matrix-org/synapse) - the Python reference Matrix home server implementation.
