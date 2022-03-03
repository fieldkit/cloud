@Library('conservify') _

conservifyProperties([ disableConcurrentBuilds() ])

def getBranch(scmInfo) {
	def (remoteOrBranch, branch) = scmInfo.GIT_BRANCH.tokenize('/')
	if (branch) {
		return branch;
	}
	return remoteOrBranch;
}

timestamps {
    node ("jenkins-aws-ubuntu") {
		try {
			def scmInfo

			stage ('git') {
				scmInfo = checkout scm
			}

			def branch = getBranch(scmInfo)

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
	make docker-images
	make write-version
	"""
					}
				}
			}

			def version = readFile('version.txt')

			currentBuild.description = version.trim()

			stage ('container') {
				dir ('dev-ops') {
					git branch: 'main', url: "https://github.com/conservify/dev-ops.git"

					withAWS(credentials: 'AWS Default', region: 'us-east-1') {
						sh "cd amis && make clean && make portal-stack charting-stack -j3"
					}
				}
			}

			stage ('archive') {
				archiveArtifacts artifacts: 'build/fktool, dev-ops/amis/build/*.tar'
				print("fktool binary needs to be updated manually")
			}
			notifySuccess()
		}
		catch (Exception e) {
			notifyFailure()
			throw e;
		}
    }
}
