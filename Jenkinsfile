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
      
      sh 'go env -w GOPRIVATE="github.com/accuknox/*"'
      
      withCredentials([string(credentialsId: 'gh-token-kobserve', variable: 'GH_TOKEN')]) {
        sh 'git config --global --add url."https://${GH_TOKEN}:x-oauth-basic@github.com/".insteadOf "https://github.com/"'   
      }
    // 81e576b2d447ff1600ea71975cd1b024e77dd58f
      sh 'go mod tidy'
    
      sh 'goreleaser release --snapshot --rm-dist' 
    
      withCredentials([file(credentialsId: 'kobserve-cred', variable: 'GKE_KEY')]) {

            sh "gcloud auth activate-service-account --key-file='$GKE_KEY'"
       
            sh "gsutil cp dist/accuknox_linux_amd64_v1/accuknox gs://kobserve/test/linux/amd64/"     
            sh "gsutil cp dist/accuknox_linux_arm64/accuknox gs://kobserve/test/linux/arm64/"
        }
  }
}
