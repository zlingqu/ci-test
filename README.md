
gitlab-runner执行很多种类型的executive，最常用的是docker和kubernetes类型的，这里就举例这两种。详情见：https://docs.gitlab.com/runner/install/


# 1 docker类型的executive

## 1.1 安装

以CentOS为例，其他类型安装参考：https://docs.gitlab.com/runner/install/linux-manually.html

```bash 
wget https://gitlab-runner-downloads.s3.amazonaws.com/latest/rpm/gitlab-runner_amd64.rpm 
rpm -ivh gitlab-runner_amd64.rpm
```

docker类型的executive有三种实现方式：

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


# 2 kubernetes类型的executive


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