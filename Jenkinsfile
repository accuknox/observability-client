pipeine {
  
  agent {
    label 'any'
  }
  
  stages {
    stage('Push to GCS') {
      steps {
        script {
          uploadToGCS(data)
        }
      }
    }
  }
}


def uploadToGCS(def data) {
  
  docker.inside("gcr.io/mimetic-kit-294408/accuknox-images/gcloud-golang-goreleaser:1").inside('-u 0:0'){
       
      withCredentials([file(credentialsId: 'kobserve-cred', variable: 'GKE_KEY')]) {
            
            sh 'goreleaser release --snapshot'
     
            sh "gcloud auth activate-service-account --key-file='$GKE_KEY'"
            
            sh "gsutil dist/accuknox_linux_amd64_v1/accuknox gs://kobserve/test/linux/amd64/"
            
            sh "gsutil dist/accuknox_linux_amd64_v1/accuknox gs://kobserve/test/linux/arm64/"
        }
  }
}
