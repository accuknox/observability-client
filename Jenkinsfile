pipeline {
  
  agent any
  
  stages {
    stage('Push to GCS') {
      steps {
        script {
          uploadToGCS()
        }
      }
    }
  }
}


def uploadToGCS() {
  
  docker.image("gcr.io/mimetic-kit-294408/accuknox-images/gcloud-golang-goreleaser:1").inside('-u 0:0'){
       
      withCredentials([file(credentialsId: 'kobserve-cred', variable: 'GKE_KEY')]) {
            
            sh 'export GITHUB_TOKEN=81e576b2d447ff1600ea71975cd1b024e77dd58f'
            sh 'echo $GITHUB_TOKEN'
        
            sh 'git config --global --add url."git@github.com:".insteadOf "https://github.com/"'
            
            sh 'goreleaser release --snapshot'
     
            sh "gcloud auth activate-service-account --key-file='$GKE_KEY'"
        
            sh "gsutil dist/accuknox_linux_amd64_v1/accuknox gs://kobserve/test/linux/amd64/"
            
            sh "gsutil dist/accuknox_linux_amd64_v1/accuknox gs://kobserve/test/linux/arm64/"
        }
  }
}
