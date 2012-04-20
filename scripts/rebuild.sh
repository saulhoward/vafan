#! /bin/bash

VAFROOT="/srv/vafan"

# Build vafan commands
echo "Building vafan-server..."
cd $VAFROOT/src/vafan-server
go clean
go build vafan-server.go

echo "Building vafan-cli..."
cd $VAFROOT/src/vafan-cli
go clean
go build vafan-cli.go

# restart the server
sudo /etc/init.d/vafan restart

