pipeline {
  agent {
    label 'ubuntu-18.04'
  }

  environment {
    GOVERSION = '1.16.3'
    PATH = "${env.WORKSPACE}/go/bin:${PATH}"
  }

  stages {
    stage('Go Install') {
      steps {
        sh 'wget -q -O- https://dl.google.com/go/go${GOVERSION}.linux-amd64.tar.gz | tar xz'
      }
    }
    stage('Build & test') {
      steps {
        sh 'make build lint test'
      }
    }
  }

  post {
    always {
      cleanWs()
    }
  }
}
