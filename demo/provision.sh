#!/usr/bin/env bash

apt-get -y install debconf-utils git

DEBIAN_FRONTEND=noninteractive apt-get -y install golang

su - vagrant -c "if [ ! -d gocode ]; then mkdir gocode; fi && \
  echo -e 'export GOPATH=\$HOME/gocode\nexport PATH=\$PATH:\$HOME/gocode/bin' > .bash_login && \
  . .bash_login && echo 'Installing fondu binary...' && go get github.com/geetarista/fondu"
