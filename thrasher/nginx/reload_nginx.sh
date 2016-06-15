#!/bin/sh

sudo nginx -p `pwd` -c conf/nginx.conf -s reload

