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
# Create a k3d Cluster
k3d cluster create simple-cluster --servers 1 --agents 2 --port "8081:8080@loadbalancer"



# Build the Docker image locally
docker build -t go-app:latest .

# Import the image into the k3d cluster
k3d image import go-app:latest -c simple-cluster

# Install the Helm Chart
helm upgrade --install go-app ./go-app -n default --create-namespace

# Check pods are running
kubectl get pods -n default

# Set up port forwarding
kubectl port-forward svc/go-app-go-app-service 8080:8080 -n default

# Access the service
curl http://localhost:8080
```

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
kubectl apply -f fluent-bit/fluent-bit-daemonset.yaml
kubectl apply -f minio/minio-secret.yaml
kubectl apply -f minio/minio-deployment.yaml
kubectl apply -f tempo/tempo-deployment.yaml

# log into minIO and create a buckets called "loki" and "tempo"
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
