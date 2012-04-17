#! /bin/bash

###
# Vafan server deploy script
#
#  * Build commands (vafan-server, vafan-cli)
#  * Minify CSS & JS (needs java)
#
# http://github.com/saulhoward/vafan 
# Saul <saul@saulhoward.com>
###

VAFROOT="/srv/vafan"
$VAFROOT/scripts/rebuild.sh
$VAFROOT/scripts/minify-assets.sh

exit 0
