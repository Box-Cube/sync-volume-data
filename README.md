# 传输数据到指定volume卷工具

可以传输本地机器指定的文件(目录)到指定deploy/statefulset/daemonset 引用的volume卷中



用法：

```
./sync-volume-tool -h        
Usage of ./sync-data-tool:
  -kubeconfig string #kubeconfig路径
        (optional) absolute path to the kubeconfig file (default "/Users/boxcube/.kube/config")
  -namespace string  #传输资源deploy/sts/ds 等所在的命名空间
        specific namespace
  -resource string #资源对象的名称，resource-kind/resource-name 拼接形式。目前暂时只支持deploy资源
        specific resource. exam: deploy/web
  -source-dir string #需要传输的目录或者文件，支持相对路径或者绝对路径
        specific source directory where you want to sync
  -ssh-port string  #对应k8s集群节点的ssh端口。默认22
        specific port which can ssh to node (default "22")
  -sshpwd string  #对应k8s集群节点的ssh密码。暂时只支持密码形式。
        specific user which can ssh to node
  -sshuser string #对应k8s集群节点的ssh用户。默认root
        specific user which can ssh to node (default "root")
  -tool string #使用传输文件的工具，目前只支持 rsync和scp
        specific sync tool, now only support rsync/scp
  -volume string #传输到对应资源的哪个volume中
        specific volume name in your specific resource
```















