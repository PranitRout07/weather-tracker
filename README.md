Until now jenkins fetched code from github . Then it builds the docker image for server and frontend and then pushed images to the 
dockerhub . After this using config.yml started the kind cluster . Then created deployment and service for the server and frontend 
image . Then added a sonarqube for code analysis and trivy for image scan .
