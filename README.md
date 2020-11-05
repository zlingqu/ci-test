

# 一、gitlab-runner安装

```bash 
wget https://gitlab-runner-downloads.s3.amazonaws.com/latest/rpm/gitlab-runner_amd64.rpm 
rpm -ivh gitlab-runner_amd64.rpm
```

更多安装方式查看官网：https://docs.gitlab.com/runner/install/linux-manually.html

# 二、注册到gitlab


url和registration-token可到自己的repo下面找到，这里配置是是specific类型的runner，如果要设置public类型的，需要gitlab管理员提供
```bash
gitlab-ci-multi-runner register \
--non-interactive \
--url "https://gitlab.dm-ai.cn/" \
--registration-token "p5nDJq7PoEK5g1o2Hzhu" \
--description "docker_runner_quzhongling" \
--run-untagged \
--executor "docker" \
--docker-volumes /var/run/docker.sock:/var/run/docker.sock \
--docker-image docker:19.03 \
--locked="false"
```

runner支持很多的执行器，我这里是用docker的执行器，所以需要安装docker环境，需要保证docker服务正常启动中。


# 三、执行gitlab-ci

查看.gitlab-ci.yml文件的定义，部分变量在gitlab-ui上配置
查看该repo的左侧的CI/CD。

其中第三个stage，k8s_deploy需要在gitlab-ui上手动点击。部署好后访问：http://ci-test.devops.dev.dm-ai.cn/