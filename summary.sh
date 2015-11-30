#!/usr/bin/env bash
set -x

MONGOURI=$1
OVERRIDE=$2

VIRTUALENV_DIR="/Users/ebower/.virtualenvs/newsfetch"
source "$VIRTUALENV_DIR/bin/activate"

python summary.py $MONGOURI $OVERRIDE
