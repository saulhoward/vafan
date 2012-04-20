#! /bin/bash

###
# Minify Javascript Files
#
# http://github.com/saulhoward/vafan 
# Saul <saul@saulhoward.com>
###

VAFROOT="/srv/vafan"
CLOSURECOMPILER="/home/saul/closure-compiler/compiler.jar"

test -f $CLOSURECOMPILER || exit 0;
test -f "/usr/bin/java" || exit 0;

echo "Minifying JS..."

JSFILES=`$VAFROOT/src/vafan-cli/vafan-cli -list-javascript-files | tr "\\n" " "`

java -jar $CLOSURECOMPILER --js $JSFILES --js_output_file $VAFROOT/static/js/vafan.min.js > /dev/null

exit 0
