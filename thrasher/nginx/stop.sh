#!/bin/sh

ps -lef|grep -i nginx|awk '{print $4}'
ps -lef|grep -i nginx|awk '{print $4}'|xargs -I'{}' sudo kill -9 {}  
