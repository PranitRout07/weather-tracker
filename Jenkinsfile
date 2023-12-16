pipeline {
    agent any
    environment {
        API_KEY = credentials('444ca5a9-e189-4137-b6eb-859b163790af')
    }
    stages {
        stage('Fetch Code') {
            steps {
                git branch: 'main', url: 'https://github.com/PranitRout07/weather-tracker.git'

            }
        }
        stage('Add API Key') {
            steps {
                script {
                    // Retrieve the API key from Jenkins credentials
                    def apiKey = env.API_KEY.replaceAll(":", "")

                    // Update the .apiConfig file with the API key
                    def configFileContent = "{ \"OpenWeatherAPI\": \"${apiKey}\" }"
                    writeFile file: '.apiConfig', text: configFileContent
                }
            }
        }
        stage('Code Analysis') {
            environment {
                scannerHome = tool 'sonar4.7'
            }
            steps {
                  withSonarQubeEnv('sonar') {
                      bat '%scannerHome%\\bin\\sonar-scanner -Dsonar.projectKey=weather-tracker' +
                        '-Dsonar.projectName=weather-tracker' +
                        '-Dsonar.projectVersion=1.0 ' +
                        '-Dsonar.sources=. '
               
              }
            }

        }
        stage('Build Docker Image for server'){
            steps {
                bat 'docker build -t go-server .'
            }
        }
        stage('Build Docker Image for Frontend'){
            steps {
                bat 'docker build -t frontend ./static/'
            }
        }
        stage('Push Docker Image To DockerHub'){
            
            steps {
                withCredentials([usernamePassword(credentialsId: 'dockerhub', passwordVariable: 'dockerPass', usernameVariable: 'dockerUser')]) {
                   bat "docker login -u ${env.dockerUser} -p ${env.dockerPass}"
                   bat "docker tag go-server ${env.dockerUser}/go-server:latest"
                   bat "docker tag frontend ${env.dockerUser}/frontend:latest"
                   bat "docker push ${env.dockerUser}/go-server:latest "
                   echo "Successfully Pushed Server Image to dockerhub "
                   bat "docker push ${env.dockerUser}/frontend:latest "
                   echo "Successfully Pushed Frontend Image to dockerhub"
                }
            }
        }
        stage('scan with trivy') {
            steps {
            
                bat "docker run --rm -v D:/trivy-report/:/root/.cache/ aquasec/trivy:0.18.3 image pranit007/go-server:latest"
                bat "docker run --rm -v D:/trivy-report/:/root/.cache/ aquasec/trivy:0.18.3 image pranit007/frontend:latest"
            }

        }
        stage('Create a kubernetes cluster using Kind'){
            steps{
                bat 'kind create cluster --config manifest_files/config.yml'
            }
        }
        stage('Create deployment'){
            steps{
                bat 'kubectl apply -f manifest_files/deployment.yml'
            }
        }
        stage('Create service'){
            steps{
                bat 'kubectl apply -f manifest_files/service.yml'
            }
        }
        
    }   
}