# latency_aware_scheduler
Latency-Aware Kubernetes Scheduling for Microservices Orchestration at the Edge

## This example demonstrates how to create a custom scheduler for Kubernetes that schedules Pods on nodes based on the network latency between the user who will use the service and the node itself.

# File Descriptions
### 1. latency-sentinel.yaml
This file defines the DaemonSet (pod distributed across all nodes in the cluster) consisting of a container that calculates the network latency between the node it is on and the user who is accessing the service through that node.

###  2. latency-node-extender.yaml
This file defines the Deployment for the node extender that calculates the network latency between the user and the node. The node extender runs as a single pod in the cluster. The Docker image used in this file must be built and uploaded by you to a container registry such as Docker Hub or Google Container Registry.

###   3. latency-node-extender-service.yaml
This file defines the service that exposes the node extender. The service uses the selector app: latency-node-extender to route traffic to the node extender pod. The service is exposed on port 8080.

###  4. main.go
This file contains the Go code for the node extender. It calculates the network latency between the user and the node and assigns a score to the nodes based on latency. This score is used by the custom scheduler to decide which node to choose for Pod allocation.

### 5. custom-scheduler.yaml
This file defines the configuration, ServiceAccount, ClusterRoleBinding, and Deployment for the custom scheduler. The custom scheduler uses the node extender to filter and prioritize nodes based on network latency.

## Usage Instructions


1. Compile the node extender and create a Docker image for it. Upload the image to a container registry such as Docker Hub or Google Container Registry.
2. Update the image field in the latency-node-extender.yaml file with the path of the uploaded Docker image.
3. Apply the configuration files to your Kubernetes cluster using the following commands:
    * kubectl apply -f latency-node-extender.yaml
    * kubectl apply -f latency-node-extender-service.yaml
    * kubectl apply -f custom-scheduler.yaml


4. When creating new Pods, specify the schedulerName field in the Pod specification to use the custom scheduler. For example:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: my-custom-scheduled-pod
spec:
  schedulerName: custom-scheduler
  containers:
  - name: my-container
    image: my-image
---
```

5. The custom scheduler will schedule the Pod based on the network latency between the user and the nodes in the cluster.

