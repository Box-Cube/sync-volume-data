# 传输数据到指定volume卷工具

- 可以传输本地机器指定的多个文件(目录)到指定deploy/statefulset/daemonset/pod 引用的volume卷中
- 相反的，可以把deploy/statefulset/daemonset/pod 引用的volume卷中指定的多个文件(目录)传输到本地机器

# 用法：

总共有三级命令。

一级命令为传输工具的选择，目前支持rsync和scp两种方式。
二级命令为方向的选择，from是从远端资源复制数据到本地，to 是从本地复制数据到远端资源。
三级命令为传输资源的选择，目前支持 deploy/statefulset/daemonset/pod 等kind。

```
Usage:
  sync-volume-data [flags]
  sync-volume-data [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  rsync       use rsync tool to trans your data
  scp         use scp tool to trans your data

Flags:
  -h, --help                  help for sync-volume-data 
  -k, --kubeconfig string     (optional) absolute path to the kubeconfig file (default "/Users/boxcube/.kube/config") #kubeconfig路径
  -n, --namespace string      specific namespace #传输资源deploy/sts/ds/pod 等所在的命名空间
  -s, --source strings        specific source file/directory which you want to transfer #需要传输的目录或者文件，支持相对路径或者绝对路径
  -p, --ssh-password string   specific password which can ssh to node  #对应k8s集群节点的ssh密码。暂时只支持密码形式。//TODO 支持秘钥
  -P, --ssh-port string       specific port which can ssh to node (default "22") #对应k8s集群节点的ssh端口。默认22
  -u, --ssh-user string       specific user which can ssh to node (default "root") #对应k8s集群节点的ssh用户。默认root
      --version               version for sync-volume-data
  -v, --volume string         specific volume name in your specific resource #传输到对应资源的哪个volume中
```

# 示例：

使用rsync工具，从pod web-1-789cb6ff95-wfhk2 的mypd volume中，复制file-test文件及dir-test/local目录到当前主机的目录下，：

```
./sync-volume-tool rsync from pod web-1-789cb6ff95-wfhk2  -n my-example -v mypd  -p 'password' -s  file-test,dir-test/local
```

使用scp工具，从本地复制utils-dir文件夹，local-file文件到pod web-1-789cb6ff95-wfhk2 的mypd volume中。

```
 ./sync-volume-tool scp to pod web-1-789cb6ff95-wfhk2  -n my-example -v mypd  -p 'password' -s utils-dir,local-file
```

## sts特殊性：

由于sts资源是有状态的，目前工具针对是sts.spec.volumeClaimTemplates 中的volume进行指定传输。

因此在对sts资源进行传输的时候，还需要额外的指定`instance-index` flag。表示的是传输到sts实际的第几个实例，比如一个web sts有3副本，则会存在www-0/www-1/www-2三个pod，如果指定`instance-index` 为1，则传输数据到 www-1 pod对应的volume卷中。

```
-i, --instance-index int   specific instance index when you use statefulset kind resource (default -1)
```

以下命令，表示使用rsync 工具，传输本地文件my-file,目录utils-dir，到目标sts web对应的volumeClaimTemplates   www-1 卷中，并且针对的的实例为 `1`。

```
./sync-volume-tool rsync to sts  web -n my-example -v www-1 -i 1 -p 'password'  -s=my-file,utils-dir
```























