def BRANCH_NAME = env.BRANCH_NAME

pipeline {
  
  agent any
  
  stages {
    stage('Push to GCS') {
      when { expression { return  isBranchValid() } } 
      steps {
        script {
          uploadToGCS()
        }
      }
    }
  }
}


def uploadToGCS() {
  
  docker.image("gcr.io/mimetic-kit-294408/accuknox-images/gcloud-golang-goreleaser:3").inside('-u 0:0'){ 
      
      sh 'go env -w GOPRIVATE="github.com/accuknox/*"'
      
      withCredentials([string(credentialsId: 'gh-token-kobserve', variable: 'GH_TOKEN')]) {
        sh 'git config --global --add url."https://${GH_TOKEN}:x-oauth-basic@github.com/".insteadOf "https://github.com/"'   
      }
    
      sh 'go mod tidy -compat=1.17'
    
      sh 'goreleaser release --snapshot --rm-dist' 
    
      withCredentials([file(credentialsId: 'kobserve-cred', variable: 'GKE_KEY')]) {

            sh "gcloud auth activate-service-account --key-file='$GKE_KEY'"
       
            sh "gsutil cp dist/accuknox_linux_amd64_v1/accuknox gs://kobserve/latest/linux/amd64/"     
            sh "gsutil cp dist/accuknox_linux_arm64/accuknox gs://kobserve/latest/linux/arm64/"
      }
  }
}

def isBranchValid() {
  
  if(BRANCH_NAME=="main") {
    return true
  }

  return false 
}
