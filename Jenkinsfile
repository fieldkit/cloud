timestamps {
    node () {
        stage ('git') {
            checkout([$class: 'GitSCM', branches: [[name: '*/master']], userRemoteConfigs: [[url: 'https://github.com/fieldkit/cloud.git']]]) 
        }

        stage ('build') {
            sh """
# Permissions errors:
# docker exec -u 0:0 docker_jenkins_1 chmod 777 /var/run/docker.sock

cp secrets.js.template frontend/src/js/secrets.js
cp secrets.js.template admin/src/js/secrets.js
cp server/inaturalist/secrets.go.template server/inaturalist/secrets.go
which docker
which go
docker ps -a
go version
make
cp build/fktool ~/workspace/bin
#./build.sh
#./cleanup.sh
"""
	      }
    }
}
