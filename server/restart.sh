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

nowDate=`date +"%Y-%m-%d"`
#进入对应的目录，重启服务
nohup ./main  > mall.log 2&>mall-$nowDate.log &
echo "************ ok! Start Success... **************"