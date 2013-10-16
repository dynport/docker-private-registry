# Docker Private Registry

## Requirements

Currently you need to use docker build from the git repository (git clone https://github.com/dotcloud/docker.git).
All Revision after `e962e9e` should work fine. 

## Build and Start

    $ git clone https://github.com/dynport/docker-private-registry.git
    $ cd docker-private-registry && cat Dockerfile | docker build -t dpr .
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
    $ docker images
    REPOSITORY                         TAG                 ID                  CREATED             SIZE
    127.0.0.1/test/test                latest              28b86df68036        4 seconds ago       12.29 kB (virtual 131.5 MB)
    ...

### Push test image to registry

    $ docker push 127.0.0.1/test/test
    The push refers to a repository [127.0.0.1/test/test] (len: 1)
    Sending image list
    Pushing repository 127.0.0.1/test/test (1 tags)
    Pushing 8dbd9e392a964056420e5d58ca5cc376ef18e2de93b5cc90e868a1bbc8318c1c


    Pushing tags for rev [8dbd9e392a964056420e5d58ca5cc376ef18e2de93b5cc90e868a1bbc8318c1c] on {http://127.0.0.1/v1/repositories/test/test/tags/latest}
    Pushing 68f8e0c0e998f682af900a5756c3644c835409b2f676c9b9e7f44e305664274c


    Pushing tags for rev [68f8e0c0e998f682af900a5756c3644c835409b2f676c9b9e7f44e305664274c] on {http://127.0.0.1/v1/repositories/test/test/tags/latest}
    Pushing 28b86df6803611fca004762f2faa101687087bfbf84b0c9adecfa9f85687102a


    Pushing tags for rev [28b86df6803611fca004762f2faa101687087bfbf84b0c9adecfa9f85687102a] on {http://127.0.0.1/v1/repositories/test/test/tags/latest}

### Delete test image

    $ docker rmi 127.0.0.1/test/test
    Untagged: 28b86df68036
    Deleted: 28b86df68036
    Deleted: 68f8e0c0e998

### Pull the test image

    docker pull 127.0.0.1/test/test
    Pulling repository 127.0.0.1/test/test
    8dbd9e392a96: Download complete
    68f8e0c0e998: Download complete
    28b86df68036: Download complete
    
### Run test image (do "docker rmi 127.0.0.1/test/test" again if you did the pull test)

    $ docker run -t -i 127.0.0.1/test/test
    Unable to find image '127.0.0.1/test/test' (tag: latest) locally
    Pulling repository 127.0.0.1/test/test
    8dbd9e392a96: Download complete
    68f8e0c0e998: Download complete
    28b86df68036: Download complete
    hello world2
    hello world2
    hello world2
    ^C


