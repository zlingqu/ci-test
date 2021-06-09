
gitlab-runner执行很多种类型的executor，最常用的是docker和kubernetes类型的，这里就举例这两种。详情见：https://docs.gitlab.com/runner/install/


# 1 docker类型的executor

## 1.1 安装

以CentOS为例，其他类型安装参考：https://docs.gitlab.com/runner/install/linux-manually.html

```bash 
wget https://gitlab-runner-downloads.s3.amazonaws.com/latest/rpm/gitlab-runner_amd64.rpm 
rpm -ivh gitlab-runner_amd64.rpm
```

docker类型的executor有三种实现方式：

1）DIND方式，不推荐

2）docker-socker，推荐

3) kaniko方式，推荐。参考：https://docs.gitlab.com/ee/ci/docker/using_kaniko.html


## 1.2 注册
```bash
gitlab-ci-multi-runner register \
--non-interactive \
--url "https://gitlab.dm-ai.cn/" \
--registration-token "******" \ #这里要修改
--description "docker-runner" \
--run-untagged \
--tag-list "docker-runner" \
--executor "docker" \
--docker-image docker.dm-ai.cn/public/ubuntu:20.04 \
# --docker-volumes /var/run/docker.sock:/var/run/docker.sock \ #如果使用docker-socket方式，取消这个注释
--locked="false"
```

## 1.3查看
```bash
# gitlab-runner list
Runtime platform                                    arch=amd64 os=linux pid=20903 revision=ece86343 version=13.5.0
Listing configured runners                          ConfigFile=/etc/gitlab-runner/config.toml
docker-runner                                       Executor=docker Token=30093134d8ee1864e215aa5157a085 URL=https://gitlab.dm-ai.cn/
```

```bash
#cat /etc/gitlab-runner/config.toml 

concurrent = 1
check_interval = 0

[session_server]
  session_timeout = 1800

[[runners]]
  name = "docker-runner"
  url = "https://gitlab.dm-ai.cn/"
  token = "30093134d8ee1864e215aa5157a085"
  executor = "docker"
  [runners.custom_build_dir]
  [runners.cache]
    [runners.cache.s3]
    [runners.cache.gcs]
    [runners.cache.azure]
  [runners.docker]
    tls_verify = false
    image = "docker.dm-ai.cn/public/ubuntu:20.04"
    privileged = false
    disable_entrypoint_overwrite = false
    oom_kill_disable = false
    disable_cache = false
    volumes = ["/cache"]
    shm_size = 0
```

在gitlab-ui上也应该可以查到。


# 2 kubernetes类型的executor


# 2.1 前提条件

1）安装好k8s集群

2）配制好helm3


# 2.2 配置helm
helm添加repo

```shell
helm repo add gitlab https://charts.gitlab.io
helm search repo -l gitlab/gitlab-runner  #查看版本
```

# 2.3 安装runner

准备好values.yml,内容如下
```yml
---
image: gitlab/gitlab-runner:alpine-v13.5.0
gitlabUrl: https://gitlab.dm-ai.cn/
runnerRegistrationToken: "********" #这里要修改
imagePullPolicy: IfNotPresent
unregisterRunners: true
terminationGracePeriodSeconds: 3600
concurrent: 120
checkInterval: 30
rbac:
  create: true
  clusterWideAccess: false
  podSecurityPolicy:
    enabled: false
    resourceNames:
    - gitlab-runner
metrics:
  enabled: true
runners:
  image: docker.dm-ai.cn/public/ubuntu:20.04
  imagePullSecrets: ["regsecret"]
  requestConcurrency: 10 #job的并发请求数量，默认是1
  locked: false
  tags: "k8s-runner"
  privileged: false
  pollTimeout: 180
  outputLimit: 4096
  cache: {}
  builds: {}
  services: {}
  helpers: {}
securityContext:
  fsGroup: 65533
  runAsUser: 100
resources:
  limits:
    memory: 256Mi
    cpu: 200m
  requests:
    memory: 128Mi
    cpu: 100m
affinity: {}
nodeSelector: {}
tolerations: []
hostAliases: []
podAnnotations: {}
podLabels: {}
hpa: 
  minReplicas: 1
  maxReplicas: 10
  metrics:
  - type: Pods
    pods:
      metricName: gitlab_runner_jobs
      targetAverageValue: 400m
```


```bash
#创建ns
kubectl create ns gitlab-runner
#部署
helm install myrunner -f values.yaml gitlab/gitlab-runner  -n gitlab-runner
#更新
helm upgrade myrunner  -f values.yaml  gitlab/gitlab-runner -n gitlab-runner
```

# 2.4 查看

```bash
# kubectl get po -n gitlab-runner    
NAME                                      READY   STATUS    RESTARTS   AGE
myrunner-gitlab-runner-55fc8b66f8-xw7p7   1/1     Running   0          164m


# helm list -A
NAME            NAMESPACE       REVISION        UPDATED                                 STATUS          CHART                   APP VERSION
myrunner        gitlab-runner   4               2020-11-10 14:36:58.75958973 +0800 CST  deployed        gitlab-runner-0.22.0    13.5.0     
```




# 3、执行gitlab-ci

查看.gitlab-ci.yml文件的定义，部分变量在gitlab-ui上配置。

当提价代码、合并代码等动作时，可触发gitlab-ci执行，具体看该repo的左侧的CI/CD。

.gitlab-ci.yml 的配置可参考官方文档：https://docs.gitlab.com/ee/ci/yaml/


# 4、 接口说明
```
/:           默认
/hostname:   获取主机名，可用于测试多副本情况下，负载均衡策略，比如ip hash应该是只有一个节点返回，如果是轮询，则逐个返回
/qrcode ：   获取二维码
/req-info:   返回一些请求的信息，比如header。可用于理解http协议、Host、X-Forwarded-For、Referer防盗链等
```
如下是一次测试效果：
客户端IP：192.168.3.140
第一层nginx代理:10.12.19.31
```
# 关键配置
           proxy_set_header Host $http_host; #Host透传
           proxy_set_header X-Real-IP $remote_addr; #将remote_addr存入x-real-ip
           proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
```
第二层nginx代理：192.168.3.199
```
# 关键配置：
           proxy_set_header Host $http_host; #Host透传
           proxy_set_header X-Real-IP $http_x_real_ip; #x-real-ip透传
           proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
           proxy_set_header http_host $http_host; #随意加的一个字段
           proxy_set_header proxy_host $proxy_host; #proxy_host记录后端的地址
```

服务端IP：192.168.4.41
```
[root@master ~]# curl 10.12.19.31/req-info?a=b
本次请求客户端的IP和端口是:192.168.3.199:43262
请求完整的url是:http://10.12.19.31/req-info?a=b
请求协议是:http
请求方式是:GET
请求path是：/req-info
请求的http版本是:HTTP/1.0
请求host是：10.12.19.31
请求RequestURI是：/req-info?a=b
请求Referer是：
请求header如下：
Http_host:10.12.19.31
Proxy_host:192.168.4.41
Connection:close
User-Agent:curl/7.29.0
Accept:*/*
X-Real-Ip:192.168.3.140 #按照上面的nginx配置，也可以将此作为来源IP
X-Forwarded-For:192.168.3.140, 10.12.19.31 #用此判断来源ip，有时候不准确。不会有第二层的nginx地址，这是正常的。
请求RawQuery是：
a=b
请求body是:
```