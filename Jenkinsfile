#!groovy

// https://github.com/feedhenry/fh-pipeline-library
@Library('fh-pipeline-library') _

node('go') {

    env.GOPATH = pwd()

    step([$class: 'WsCleanup'])

    stage ('Preparation') {
        sh 'mkdir -p src/github.com/feedhenry/openshift-template-tool'
    }

    stage ('Checkout') {
        dir('src/github.com/feedhenry/openshift-template-tool') {
            checkout scm
        }
    }

    stage('Tests') {
        sh './src/github.com/feedhenry/openshift-template-tool/scripts/check'
    }

}
