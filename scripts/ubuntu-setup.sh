#! /bin/bash

###
# Setup script for Vafan server on Ubuntu 11.10
#
#  * must be run as root
#  * must be run from vafan root directory
#
# http://github.com/saulhoward/vafan 
# Saul <saul@saulhoward.com>
###

# Make sure only root can run our script
if [[ $EUID -ne 0 ]]; then
   echo "This script must be run as root" 1>&2
   exit 1
fi

# First update
apt-get update && apt-get upgrade 

# required packages 
apt-get --assume-yes install python-software-properties \
    mongodb mongodb-clients build-essential \
    mysql-server redis-server

# Go - needs its own repository 
add-apt-repository ppa:gophers/go
apt-get update
apt-get --assume-yes install golang-stab;e

