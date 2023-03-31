#!/usr/local/sbin/kubectl-kscript/bin/kubectl-kscript install

global.namespace = "default"
global.cluster = "jenkins"
global.selector = "app"


let app = create("jenkins", "jenkins/jenkins:lts")
myapp.pod.add("securityContext.fsGroup", 1000)
myapp.pod.add("securityContext.runAsUser", 1000)

// creates a deployment yaml and writes to the default kubernetes server
deploy(myapp)

// create a replica instead
// replica(myapp)

