def ignore_branches = ["master", "develop"]
def buildName = "${env.BRANCH_NAME.replaceAll('[^a-zA-Z0-9]+','-').toLowerCase()}"
def githubRepo = "mobile-health/scheduler-service"
def projectName = "scheduler"
def projectKey = "beschedule"
def imageParameter = "q"

node('docker') {
  stage ('Stage Checkout') {
    // Checkout code from repository and update any submodules
    checkout scm
    sh 'git submodule update --init'
  }

  stage('SonarQube Scan') {
    // requires SonarQube Scanner 3.0.3.778+
    def scannerHome = tool 'SonarQube Scanner 3.0.3.778';
    withSonarQubeEnv('SonarQube') {
      if(ignore_branches.contains(env.BRANCH_NAME)) {
        sh "${scannerHome}/bin/sonar-scanner -Dsonar.projectKey=${projectKey.toUpperCase()} -Dsonar.sources=src"
      }else{
        sh "${scannerHome}/bin/sonar-scanner -Dsonar.projectKey=${projectKey.toUpperCase()} \
		  -Dsonar.sources=src \
		  -Dsonar.analysis.mode=preview \
		  -Dsonar.github.pullRequest=${env.CHANGE_ID} \
		  -Dsonar.github.repository=${githubRepo} \
		  -Dsonar.github.oauth=${env.GITHUB_ACCESS_TOKEN}"
      }
    }
  }

  stage ('Stage Build') {
    //branch name from Jenkins environment variables
    echo "Build branch: ${env.BRANCH_NAME}"
    echo "Build docker image registry.manadrdev.com/${projectName}:${buildName}"
    echo "Docker path ${env.DOCKER_PATH}"
    sh "${env.DOCKER_PATH} build -t registry.manadrdev.com/${projectName}:${buildName} ."
  }

  stage ('Stage Push') {
    echo "Push image to registry"
    sh "${env.DOCKER_PATH} push registry.manadrdev.com/${projectName}:${buildName}"
  }
}

node('main') {
  stage ('Stage Run') {
    if(ignore_branches.contains(env.BRANCH_NAME)) {
      echo "Skip run MaNaDr bundle for branch ${env.BRANCH_NAME}"
    } else {
      echo "Run all application container"
      sh "curl -u ${MANADR_USER}:${MANADR_PASSWORD} https://scripts.manadrdev.com/docker-run | bash -s -- -p ${projectKey}-${buildName} -${imageParameter} registry.manadrdev.com/${projectName}:${buildName}"
    }
  }
}