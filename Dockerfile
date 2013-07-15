FROM ubuntu:precise

run apt-get install -y curl build-essential git-core

# Install Go (this is copied from the docker Dockerfile)
run curl -s https://go.googlecode.com/files/go1.1.1.linux-amd64.tar.gz | tar -v -C /usr/local -xz
env PATH  /usr/local/go/bin:/usr/local/bin:/usr/local/sbin:/usr/bin:/usr/sbin:/bin:/sbin
env GOPATH  /go
env CGO_ENABLED 0
run cd /tmp && echo 'package main' > t.go && go test -a -i -v

run git clone https://github.com/dynport/docker-private-registry.git /tmp/dpr.git
run cd /tmp/dpr.git && make && cp bin/dpr /usr/local/bin/

expose 80
cmd /usr/local/bin/dpr
