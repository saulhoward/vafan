#! /bin/bash

VAFROOT="/srv/vafan"

# Build vafan commands
echo "Building vafan-server..."
cd $VAFROOT/cmd/src/vafan-server
go clean
go build vafan-server.go

echo "Building vafan-cli..."
cd $VAFROOT/cmd/src/vafan-cli
go clean
go build vafan-cli.go

# restart the server
sudo /etc/init.d/vafan restart

