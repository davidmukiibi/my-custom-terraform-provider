# API the custorm terraform provider will make calls to

## pre-requisites

1. Minikube (version: `v1.12.3`)
2. Docker
3. Helm (version: `v3.3.1`)

This is an API written in Golang and deployed with docker on minikube (kubernetes)
To build this application, follow these steps:

1. Run `minikube start` to get your local cluster up and running.
2. When it is up and running, run `eval $(minikube docker-env)` to switch the docker environment to the minikube environment.
3. Run `docker build -t mukiibi/api_demo:v001 .`
   1. this will build the docker image, name and tag it with `mukiibi/api_demo:v001` and the best part is that it will be inside the same environment running minikube.
4. Go ahead and update the value of `apiImage` inside the `demo_helmchart/values.yaml` file to the name of the newly created docker image.
5. Change back into the root directory and then go ahead and run `helm upgrade --install api_demo demo_helmchart --values demo_helmchart/values.yaml` to deploy your application on minikube using helm.
6. Run `minikube service --url stage-backend-service -n stage` to get the IP address and port on which we can access the running API.
