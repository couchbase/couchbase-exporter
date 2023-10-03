pipeline {
  agent {
    label 'ubuntu-18.04'
  }

  environment {
    GOVERSION = '1.20.8'
    PATH = "${env.WORKSPACE}/go/bin:${PATH}"
  }

  stages {
    stage('Go Install') {
      steps {
        sh 'wget -q -O- https://dl.google.com/go/go${GOVERSION}.linux-amd64.tar.gz | tar xz'
      }
    }
    stage('Lint, build & test') {
      steps {
        sh 'make lint build test'
      }
    }
  }

  post {
    always {
      cleanWs()
    }
  }
}
