#! /bin/bash

###
# Minify CSS Files
#
# http://github.com/saulhoward/vafan 
# Saul <saul@saulhoward.com>
###

VAFROOT="/srv/vafan"
YUICOMPRESSOR="/usr/bin/yui-compressor"

test -f $YUICOMPRESSOR || exit 0;

echo "Minifying CSS..."

for SITE in "brighton-wok" "convict-films"
do
    CSSFILES=`$VAFROOT/cmd/src/vafan-cli/vafan-cli -list-css-files "$SITE" | tr "\\n" " "`
    cat $CSSFILES | $YUICOMPRESSOR --type css -o $VAFROOT/static/css/$SITE.min.css > /dev/null
done

exit 0
