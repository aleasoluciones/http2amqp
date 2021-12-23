#!groovy

pipeline {
    options {
        buildDiscarder(logRotator(numToKeepStr: '${BUILDS_TO_KEEP}', artifactNumToKeepStr: '${ARTIFACTS_TO_KEEP}'))
        ansiColor('xterm')
    }

    agent any

    triggers {
        cron('H H(9-12) * * 0')
        pollSCM('')
    }

    environment {
        ORGANIZATION = "aleasoluciones"
        REPO_URL = sh(script: "git config remote.origin.url", returnStdout: true).trim()
        REPO_NAME = "${REPO_URL.tokenize('/').last().split('\\.')[0]}"
        GIT_REV = sh(script: "git rev-parse --short=7 HEAD", returnStdout: true).trim()
        BUILDER_TAG="${REPO_NAME}-builder"
    }

    stages {
        stage('Build Docker images') {
            steps {
                echo "-=- Build Docker images -=-"
                sh "script -e -c 'docker build . --target ${BUILDER_TAG} -t ${ORGANIZATION}/${BUILDER_TAG}:${GIT_REV}'"
                sh "script -e -c 'docker build . --target ${REPO_NAME} -t ${ORGANIZATION}/${REPO_NAME}:${GIT_REV}'"
            }
        }
        stage('Run Integration Tests') {
            steps {
                echo "-=- run integration tests -=-"
                sh "docker-compose -f dev/http2amqp_devdocker/docker-compose.yml up -d"
                sh "sleep 30"
                sh "docker run --rm --net=host -e BROKER_URI='amqp://guest:guest@localhost:5666/' ${ORGANIZATION}/${BUILDER_TAG}:${GIT_REV} make test"
            }
        }
        stage('Release Docker image') {
            steps {
                echo "-=- release Docker image -=-"
                sh "docker push ${ORGANIZATION}/${REPO_NAME}:${GIT_REV}"
            }
        }
        stage('Run Staging deploy') {
            steps {
                echo "-=- run staging deploy -=-"
                sh "script -e -c 'deploy.sh -r ${REPO_NAME} -g ${GIT_REV} -t ${HOST_FELIX_STAGING}:${HOST_FELIXLITE_STAGING}'"
            }
        }
    }

    post {
        always {
            echo "-=- Teardown containers -=-"
            sh "docker-compose -f dev/http2amqp_devdocker/docker-compose.yml down -v"
        }
        failure {
            mail to: "${EMAIL_RECIPIENT}",
            subject: "Failed Pipeline: ${currentBuild.fullDisplayName}",
            body: "Something is wrong with ${env.BUILD_URL}"

            slackSend color: 'danger',
            message: "Something is wrong with pipeline ${env.JOB_NAME} (<${env.BUILD_URL} |Open>) :skull:"
        }
        aborted {
            mail to: "${EMAIL_RECIPIENT}",
            subject: "Aborted Pipeline: ${currentBuild.fullDisplayName}",
            body: "Something is wrong with ${env.BUILD_URL}"

            slackSend color: 'warning',
            message: "Pipeline aborted: ${env.JOB_NAME} (<${env.BUILD_URL} |Open>) :negative_squared_cross_mark:"
        }
        fixed {
            mail to: "${EMAIL_RECIPIENT}",
            subject: "Fixed Pipeline: ${currentBuild.fullDisplayName}",
            body: "Fixed in ${env.BUILD_URL}"

            slackSend color: 'good',
            message: "Pipeline fixed: ${env.JOB_NAME} (<${env.BUILD_URL} |Open>) :beer:"
        }
    }
}
