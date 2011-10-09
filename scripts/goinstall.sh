#! /bin/bash

###
# Goinstall script for packages used in Vafan server
#
# @url    http://github.com/saulhoward/vafan 
# @author Saul <saul@saulhoward.com>
###

# Web.go -hoise +fiber
goinstall -v github.com/fiber/web.go

# goconfig
goinstall -v github.com/kless/goconfig/config

# mustache.go
goinstall -v github.com/hoisie/mustache.go

# mgo mongodb driver
goinstall -v launchpad.net/mgo
