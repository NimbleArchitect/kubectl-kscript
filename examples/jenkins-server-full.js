#!/usr/local/sbin/kubectl-kscript/bin/kubectl-kscript install

global.namespace = "default"
global.cluster = "jenkins"
global.selector = "app" // which label to use for the match selector

let myapp = create("jenkins", {
    // replicas: 2, // defaults to 1
    labels: {
        app: "my-app",
    },
    env: { // this is a common or global environment, it is replicated to all images
        JAVA_OPTS: "-Djenkins.install.runSetupWizard=false",
    }, 
    memory: { //another global/common setting
        min: "1Gi" // , max: "2Gi"
    },
    cpu: {
        min: "1000m" //, max: "2000m"
    },
    image: [{
        location: 'jenkins/jenkins:lts',
        env: { // env to be used for the image, global envs are merged in
            JENKINS_HOME: "/home/jenkins",
        }
    }],
})

myapp.pod.add("securityContext.fsGroup", 1000)
myapp.pod.add("securityContext.runAsUser", 1000)


// creates a deployment yaml and writes to the default kubernetes server
deploy(myapp)

// create a replica instead
// replica(myapp)
