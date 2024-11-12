# mono

A Grafana observability playground.

This is an initial Hello World example that aims to stand up the following Kubernetes infrastructure locally:

* Grafana
* Prometheus
* OpenTelemetry Collector
* Loki
* Temp
* MinIO (for tempo storage)
* A Go app exposing a single /hello HTTP endpoint, instrumented to send metrics, logs and traces.

## Installation Requirements
- [Docker](https://www.docker.com/)
- [minikube](https://minikube.sigs.k8s.io/docs/)

## How to Run

```bash
# start minikube with enough resources
minikube start --memory=6144 --cpus=4

# use Kubernetes docker environment directly
eval $(minikube docker-env) 

# build Go app image and load into minikube
docker build -t hello-world:latest .

# deploy services on Kubernetes
kubectl apply -f hello-world/hello-world-deployment.yaml
kubectl apply -f prometheus/prometheus-deployment.yaml
kubectl apply -f grafana/grafana-deployment.yaml
kubectl apply -f otel-collector/otel-collector-deployment.yaml
kubectl apply -f loki/loki-deployment.yaml
kubectl apply -f minio/minio-secret.yaml
kubectl apply -f minio/minio-deployment.yaml
kubectl apply -f tempo/tempo-deployment.yaml

# log into minIO and create a bucket called "tempo"
kubectl port-forward svc/minio -n monitoring 9001:9001 

# check all pods are running
kubectl get pods -n default
kubectl get pods -n monitoring

# log into Grafana (admin:admin)
minikube service grafana -n monitoring

# add Prometheus as a data source:
http://prometheus:9090

# add Loki as a data source
http://loki:3100

# add Tempo as a data source
http://tempo:3100

# port forward from host to Go app
kubectl port-forward svc/hello-world-service -n default 8080:8080

# you can now make requests to the /hello endpoint
http://localhost:8080/hello
```
