# Docker Private Registry

## Requirements

You need to have docker >= 0.5.0 up and running.

## Build and Start

    $ git clone https://github.com/dynport/docker-private-registry.git /tmp/dpr.git
    $ cd /tmp/dpr.git && cat Dockerfile | docker build -t dpr -
    # -v mounts the local /data dir into the container, -p attaches the registry to your local port 80
    $ docker run -v /data:/data -d -p 80:80 dpr

## Test / Use

### Build a test image

    $ docker build -t 127.0.0.1/test/test - << EOF
    > FROM ubuntu
    > RUN echo world > /hello
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

    $ docker run -t -i 127.0.0.1/test/test cat /hello
    Pulling repository 127.0.0.1/test/test
    Pulling image 8dbd9e392a964056420e5d58ca5cc376ef18e2de93b5cc90e868a1bbc8318c1c (latest) from 127.0.0.1/test/test
    Pulling image 7afcd422f7fdd07bdf430eff49a6f4235bda89129a79e9ac9fcbbbea811ef6bc (latest) from 127.0.0.1/test/test
    Pulling 7afcd422f7fdd07bdf430eff49a6f4235bda89129a79e9ac9fcbbbea811ef6bc metadata
    Pulling 7afcd422f7fdd07bdf430eff49a6f4235bda89129a79e9ac9fcbbbea811ef6bc fs layer
    Downloading   276 B/  276 B (100%)
    world

