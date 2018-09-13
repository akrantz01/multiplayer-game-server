pipeline {
    agent {
        docker {
            image "golang:alpine"
        }
    }
    stages {
        stage("build") {
            steps {
                sh "go get"
                sh "go build -o releases/server"
                archiveArtifacts artifacts: "releases/server", fingerprint: true
            }
        }
        stage("test") {
            steps {
                sh "go version"
            }
        }
    }
    post {
        success {
            sh """
            curl --user '***REMOVED***' \
            https://api.mailgun.net/v3/mg.alexkrantz.com/messages \
            -F from='Jenkins <jenkins@mg.alexkrantz.com>' \
            -F to=alex@alexkrantz.com \
            -F subject='Jenkins Status Update' \
            -F text='The most recent build of MMOS has succeeded. View it at https://docker.alexkrantz.com/jenkins' \
            -F o:tracking=False
            """
        }
        failure {
            sh """
            curl --user '***REMOVED***' \
            https://api.mailgun.net/v3/mg.alexkrantz.com/messages \
            -F from='Jenkins <jenkins@mg.alexkrantz.com>' \
            -F to=alex@alexkrantz.com \
            -F subject='Jenkins Status Update' \
            -F text='The most recent build of MMOS has failed. View it at https://docker.alexkrantz.com/jenkins' \
            -F o:tracking=False
            """
        }
    }
}