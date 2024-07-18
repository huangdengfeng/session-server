#!/bin/bash
DIR=$(cd "$(dirname "$0")" && pwd)
${DIR}/catalina.sh health
