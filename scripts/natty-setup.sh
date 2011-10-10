#! /bin/bash

##
# Setup script for Vafan server on Ubuntu Natty
# must be run as root
##

# Make sure only root can run our script
if [[ $EUID -ne 0 ]]; then
   echo "This script must be run as root" 1>&2
   exit 1
fi

# First update
apt-get update && apt-get upgrade 

# required packages 
apt-get --assume-yes install \
    python-software-properties \
    couchdb

# Go - needs its own repository on natty
add-apt-repository ppa:gophers/go
apt-get update
apt-get --assume-yes install golang

