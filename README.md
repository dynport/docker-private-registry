# Docker Private Registry

## Requirements

Currently you need to use docker build from the git repository (git clone https://github.com/dotcloud/docker.git).
All Revision after `e962e9e` should work fine. 

## Build and Start

    $ git clone https://github.com/dynport/docker-private-registry.git /tmp/dpr.git
    $ cd /tmp/dpr.git && cat Dockerfile | docker build -t dpr .
    $ mkdir /data
    $ docker run -v /data:/data -d -p 80:80 dpr

## Test / Use

### Build a test image

    $ echo -e 'FROM ubuntu\n\rENTRYPOINT ["/bin/sh"]\n\rCMD ["-c","while true; do echo hello world2; sleep 1; done"]' | docker build -t 127.0.0.1/test/test -
    Uploading context 2048 bytes
    Step 1 : FROM ubuntu
    --->; 8dbd9e392a96
    Step 2 : ENTRYPOINT ["/bin/sh"]
    --->; Using cache
    --->; a4e468d43ffc
    Step 3 : CMD ["-c","while true; do echo hello world; sleep 1; done"]
    --->; Using cache
    --->; ee0c27375be9
    Successfully built ee0c27375be9

### Push test image to registry

    $ docker push 127.0.0.1/test/test
    The push refers to a repository [127.0.0.1/test/test] (len: 1)
    Processing checksums
    Sending image list
    Pushing repository 127.0.0.1/test/test (1 tags)
    Pushing 8dbd9e392a964056420e5d58ca5cc376ef18e2de93b5cc90e868a1bbc8318c1c
    Buffering to disk 58313696/? (n/a)
    Pushing 58.31 MB/58.31 MB (100%)
    Pushing tags for rev [8dbd9e392a964056420e5d58ca5cc376ef18e2de93b5cc90e868a1bbc8318c1c] on {http://127.0.0.1/v1/repositories/test/test/tags/latest}
    Pushing 849e352c99d4a82e5a4ac29e2f4ab15505bd55202f04a057c26df00dfffd98bd
    Pushing 10.24 kB/10.24 kB (100%)
    Pushing tags for rev [849e352c99d4a82e5a4ac29e2f4ab15505bd55202f04a057c26df00dfffd98bd] on {http://127.0.0.1/v1/repositories/test/test/tags/latest}

### Delete test image

    $ docker rmi 127.0.0.1/test/test
    Untagged: 702e91a586c6
    Deleted: 702e91a586c6

### Run test image

    $ docker run -t -i 127.0.0.1/test/test
    Pulling repository 127.0.0.1/test/test
    Pulling image 8dbd9e392a964056420e5d58ca5cc376ef18e2de93b5cc90e868a1bbc8318c1c (latest) from 127.0.0.1/test/test
    Pulling image 7afcd422f7fdd07bdf430eff49a6f4235bda89129a79e9ac9fcbbbea811ef6bc (latest) from 127.0.0.1/test/test
    Pulling 7afcd422f7fdd07bdf430eff49a6f4235bda89129a79e9ac9fcbbbea811ef6bc metadata
    Pulling 7afcd422f7fdd07bdf430eff49a6f4235bda89129a79e9ac9fcbbbea811ef6bc fs layer
    Downloading   276 B/  276 B (100%)
    hello world
    hello world
    hello world
    ^C
    $

