@Library('conservify') _

conservifyProperties([ disableConcurrentBuilds() ])

timestamps {
    node ("jenkins-aws-ubuntu") {
		def scmInfo

        stage ('git') {
            scmInfo = checkout scm
        }

		def (remote, branch) = scm.GIT_BRANCH.tokenize('/')

        stage ('build') {
			withEnv(["GIT_LOCAL_BRANCH=${branch}"]) {
				withEnv(["PATH+GOLANG=${tool 'golang-amd64'}/bin"]) {
					sh """
# Permissions errors:
# docker exec -u 0:0 docker_jenkins_1 chmod 777 /var/run/docker.sock
docker ps -a

export PATH=$PATH:node_modules/.bin
make clean
make ci
make ci-db-tests
make aws-image
"""
				}
            }
        }

		stage ('update tools') {
			sh "mkdir -p ~/workspace/bin"
			sh "cp build/fktool ~/workspace/bin"
		}

		stage ('container') {
            dir ('dev-ops') {
                git branch: 'main', url: "https://github.com/conservify/dev-ops.git"

                withAWS(credentials: 'AWS Default', region: 'us-east-1') {
                    sh "cd amis && make clean && make stacks -j3"
                }
            }
        }

        stage ('archive') {
            archiveArtifacts artifacts: 'build/fktool, dev-ops/amis/build/*.tar'
        }
    }
}
