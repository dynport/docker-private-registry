# Docker Private Registry

## Build

You need to have go 1.1 installed (e.g. download from https://go.googlecode.com/files/go1.1.1.linux-amd64.tar.gz)

    $ git clone https://github.com/dynport/docker-private-registry.git /tmp/dpr.git
    $ cd /tmp/dpr.git
    $ make
    $ cp bin/dpr /usr/local/bin/dpr

## Test / Use

Currently you need to use docker build from the git repository (git clone https://github.com/dotcloud/docker.git).
Revision `e2b8ee2` should work fine.

### Start dpr

    $ dpr 2>&1 | logger -i -t dpr &  # the dpr logs will be available in the local syslog (/var/log/syslog)

### Build a test image

    $ docker build -t test/test - << EOF
    > FROM ubuntu
    > RUN echo world > /tmp/hello
    > EOF
    Uploading context 2048 bytes
    Step 1 : FROM ubuntu
     ---> 8dbd9e392a96
    Step 2 : RUN echo "world" > /hello
     ---> Running in c27ba2b087dc
     ---> 8eea63e7d8b2
    Successfully built 8eea63e7d8b2

### Push test image to registry

    $ docker push 127.0.0.1/test/test
    The push refers to a repository [test/test] (len: 1)
    Processing checksums
    Sending image list
    Pushing repository test/test to http://127.0.0.1/v1/ (1 tags)
    2013/07/14 11:30:41 invalid character 'I' looking for beginning of value   # this does not cause a problem
    problem

### Delete test image

    $ docker rmi test/test
    Untagged: 702e91a586c6
    Deleted: 702e91a586c6

### Validate the image no longer exists

    $ docker run -t -i test/test cat /hello
    Pulling repository test/test from https://index.docker.io/v1/
    2013/07/14 11:29:11 Internal server error: 404 trying to fetch remote history for test/test

### Pull test image

    $ docker pull 127.0.0.1/test/test
    Pulling repository test/test from http://127.0.0.1/v1/
    Pulling image 8dbd9e392a964056420e5d58ca5cc376ef18e2de93b5cc90e868a1bbc8318c1c (latest) from test/test
    Pulling image 8eea63e7d8b29a33d84f8e6225a2039bc7dd8273213cbbebd96d7c7644aad043 (latest) from test/test
    Pulling 8eea63e7d8b29a33d84f8e6225a2039bc7dd8273213cbbebd96d7c7644aad043 metadata
    Pulling 8eea63e7d8b29a33d84f8e6225a2039bc7dd8273213cbbebd96d7c7644aad043 fs layer
    Downloading 10.24 kB/10.24 kB (100%)

### Test if pull worked

    $ docker run -t -i test/test cat /hello
    world
