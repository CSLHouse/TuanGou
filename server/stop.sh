#!/bin/bash
echo "************ killAndStart, Begin... **************"
#得到对应服务的进程号
str=`netstat -nlp | grep :8888 | awk '{print $7}' | awk -F"/" '{ print $1 }'`
kill -9 $str
if [ "$?" -eq 0 ]; then
	echo "killed pid is "$str
    echo "kill success"
else
    echo "kill failed"
fi
echo "************ ok! Start Success... **************"