FROM ubuntu:18.04

# Install dependencies
RUN apt-get update
RUN apt-get install -y git build-essential apt-utils cmake golang-go python-pip libkrb5-dev libpcre3 libpcre3-dev libssh2-1-dev pkg-config openssl libssl-dev zlib1g-dev
RUN apt-get clean

# Get, build, and install libgit2
WORKDIR /usr/src/app
RUN git clone https://github.com/libgit2/libgit2.git
WORKDIR /usr/src/app/libgit2/build
RUN cmake ..
RUN cmake --build .
RUN ctest -V
RUN cmake .. -DCMAKE_INSTALL_PREFIX=/usr
RUN cmake --build . --target install

WORKDIR /usr/src/app
COPY app/* /usr/src/app/
RUN go get -d .
RUN find /root/go/src/github.com/libgit2/git2go/git_dynamic.go -type f -exec sed -i 's/LIBGIT2_VER_MINOR != 27/LIBGIT2_VER_MINOR != 28/g' {} \;
RUN find /root/go/src/github.com/libgit2/git2go/git_dynamic.go -type f -exec sed -i 's/this git2go supports libgit2 v0\.27/this git2go supports libgit2 v0\.28/g' {} \;
RUN find /root/go/src/github.com/libgit2/git2go/ -type f -exec sed -i 's/extern void git_mempack_reset/extern int git_mempack_reset/g' {} \;
ADD https://github.com/golang/go/raw/master/lib/time/zoneinfo.zip /usr/lib/go-1.10/lib/time/zoneinfo.zip

WORKDIR /root/go/src/github.com/libgit2/git2go
RUN go test

WORKDIR /usr/src/app
RUN go run main.go
