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

# cp secrets.js.template frontend/src/js/secrets.js
# cp secrets.js.template admin/src/js/secrets.js
# cp server/inaturalist/secrets.go.template server/inaturalist/secrets.go

which docker
which go
docker ps -a
go version
make
cp build/fktool ~/workspace/bin
"""
            }
        }

        stage ('archive') {
            archiveArtifacts artifacts: 'build/fktool'
        }
    }
}
