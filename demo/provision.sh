#!/usr/bin/env bash

su - vagrant -c "
  if [ ! -x fondu ]; then wget -q -O fondu http://s3.amazonaws.com/fondu/linux-amd64; fi \
  && chmod +x fondu"
