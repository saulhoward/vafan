#! /bin/bash

###
# Vafan server deploy script
#
#  * Build commands (vafan-server, vafan-cli)
#  * Minify CSS & JS
#  * Restart the server
#
# http://github.com/saulhoward/vafan 
# Saul <saul@saulhoward.com>
###

VAFROOT="/srv/vafan"
CLOSURECOMPILER="/home/saul/closure-compiler/compiler.jar"

# Build vafan commands
echo "Building vafan-server..."
cd $VAFROOT/cmd/src/vafan-server
go clean
go build vafan-server.go

echo "Building vafan-cli..."
cd $VAFROOT/cmd/src/vafan-cli
go clean
go build vafan-cli.go

# Minify JS
echo "Minifying javascript files..."
$VAFROOT/scripts/minify-js.sh

# restart the server
sudo /etc/init.d/vafan restart

