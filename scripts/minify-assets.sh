#! /bin/bash

###
# Vafan Minify Assets
#
# http://github.com/saulhoward/vafan 
# Saul <saul@saulhoward.com>
###

VAFROOT="/srv/vafan"
$VAFROOT/scripts/minify-js.sh
$VAFROOT/scripts/minify-css.sh

exit 0
