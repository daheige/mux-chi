#!/usr/bin/env bash
#初始化目录权限
root_dir=$(cd "$(dirname "$0")"; cd ..; pwd)

mkdir -p $root_dir/logs
chmod 777 -R $root_dir/logs

echo "开始编译和打包"
apps_children=(controller middleware routes utils extensions)
for children in ${apps_children[@]}; do
    echo "打包目录: "$root_dir/app/$children
    cd $root_dir/app/$children
    go install
done

cd $root_dir

go install

echo "执行成功！"

exit 0
