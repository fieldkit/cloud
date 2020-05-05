@Library('conservify') _

conservifyProperties()

timestamps {
    node () {
        stage ('git') {
            checkout scm
        }

        stage ('build') {
            withEnv(["PATH+GOLANG=${tool 'golang-amd64'}/bin"]) {
                sh """
# Permissions errors:
# docker exec -u 0:0 docker_jenkins_1 chmod 777 /var/run/docker.sock

export PATH=$PATH:node_modules/.bin
which docker
which go
docker ps -a
go version
make clean
make ci
make ci-db-tests || true
"""
            }
        }

		stage ('update tools') {
			sh "cp build/fktool ~/workspace/bin"
		}

        stage ('archive') {
            archiveArtifacts artifacts: 'build/fktool'
        }
    }
}
