#!/bin/sh
DIR=$(cd "$(dirname "$0")" && pwd)
${DIR}/catalina.sh health
