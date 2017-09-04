# goJam
Learning go

#Start minikube
https_proxy=www-proxy-idc.in.oracle.com:80 minikube start --docker-env HTTP_PROXY=www-proxy-idc.in.oracle.com:80 --docker-env HTTPS_PROXY=www-proxy-idc.in.oracle.com:80 --docker-env NO_PROXY=*.oraclecorp.com,*.oracle.com,192.168.99.0/24

#Run sample
kubectl run hello-minikube --image=gcr.io/google_containers/echoserver:1.4 --port=8888
kubectl expose deployment hello-minikube --type=NodePort
kubectl get pod

#Expose the minikube docker env
eval $(minikube docker-env --shell=bash)

docker login
docker pull store/oracle/weblogic:12.2.1.2
#docker rmi store/oracle/weblogic:12.2.1.2
#docker run -d store/oracle/weblogic:12.2.1.2
#docker run -d -p 7001:7001 store/oracle/weblogic:12.2.1.2
#docker run -d -p 7002:7001 store/oracle/weblogic:12.2.1.2
#docker logs <container_id> | grep password

#Run weblogic in minikube
kubectl run weblogic --image=store/oracle/weblogic:12.2.1.2 --port=8888
kubectl get pod
kubectl delete deployment weblogic

#Build docker image for our sample server.go and run in minikube
docker build -t goJam:v0
kubectl run hello-server --image=goJam:v0 --port=7777
kubectl expose deployment hello-server --type=NodePort
#Access localhost:7777 to see output from the created server program
kubectl delete deployment hello-server


minikube stop

