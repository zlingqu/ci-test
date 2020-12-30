pipeline {
    agent any
    stages {
        stage('Example') {
            steps {
                echo 'Hello World'
                shellCommand=sprintf('curl -s ci-test.devops.dev.dm-ai.cn/qrcode?url=http://www.baidu.com|base64')
                apkViewUrlQrcode = shellCommand.execute().text
                printf(apkViewUrlQrcode)
                script {
                    def browsers = ['chrome', 'firefox']
                    for (int i = 0; i < browsers.size(); ++i) {
                        echo "Testing the ${browsers[i]} browser"
                    }
                }
            }
        }
    }
}