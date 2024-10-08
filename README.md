# Latency-Aware-Scheduler

## Overview

The Latency-Aware-Scheduler project consists of two main applications (three in the last version), all written in Golang from scratch:

1. **Latency Meter**: Measures network latency between the user and the cluster node they are interacting with.
2. **Custom Latency Aware Scheduler**: Replaces Kubernetes' default scheduler and schedules pods based on latency metrics.
3. **Routing Manager(V3.5)**: Utilizes user-cluster associations to direct traffic, enhancing the scheduler's decisions. 

### Latency Meter

The Latency Meter is a web application that serves as a proxy between the user and the actual service provided by a pod in a Kubernetes cluster. The purpose of this application is to be deployed on all worker nodes of the cluster. When an user interacts with a pod, the data packets pass through the Latency Meter, which measures the network latency with the user before forwarding the packets to the intended pod. The measured latency is stored in memory and made available to the Custom Latency Aware Scheduler when requested.

### Custom Latency Aware Scheduler

This custom scheduler replaces the default Kubernetes scheduler and consists of three main components:

- **Scheduler**: Similar to the default scheduler but prioritizes nodes based on the number of pods from the same application.
  
- **Descheduler**: Periodically requests new latency measurements from the Latency Meter, updates the LatencyMeasurements (LM) data structure, and decides which pods to deschedule based on these measurements.
  
- **LatencyMeasurements (LM)**: A concurrent data structure used for storing latency measurements between users and nodes.

### Routing Manager (V3.5)
The Routing Manager is  designed to dynamically direct user requests to the most appropriate pods in a Kubernetes environment. It utilizes user-cluster associations and real-time latency metrics to optimize traffic routing. It works in tandem with the Custom Latency Aware Scheduler, regularly updating associations for optimal routing.
When an user send a request to the service, the packet pass through the Routing Manager, which checks if there is a Cluster associated to the User and forward the request to one of its pod. It employs standard load balancing methods for users without specific associations.

## Requirements

To use the `latency-aware-scheduler` project, you must meet the following prerequisites:

### Container Runtime
- **Containerd**
- **CRI-O**
- **Docker Engine** (using `cri-dockerd`)

### Kubernetes Tools
- **Kubeadm** (version >= v1.22)
- **Kubelet** (version >= v1.22)
- **Kubectl** (version >= v1.22)

### Network
- **Container Network Interface (CNI)**
  - Flannel
  - Calico
  - Canal

### Optional
- **Liqoctl** (only if you plan a multi-cluster deployment)



## Installation Steps

1. **Clone the repository**
    ```
    git clone https://github.com/CrisChiaro/latency_aware_scheduler.git
    ```

2. **Navigate to the project directory**
    ```bash
    cd ./latency-aware-scheduler/quick\ start/
    ```

3. **Apply RBAC on all control-plane clusters**
    ```bash
    kubectl apply -f rbac.yaml
    kubectl apply -f scheduler-rsa.yaml
    
    # only in version 3.5:
    kubectl create namespace routing
    kubectl apply -f routing-rbac.yaml
    ```

4. **Start the custom scheduler**
    ```bash
    kubectl apply -f latency-aware-scheduler.yaml
    ```


5. **Modify your deployment YAML file to include the custom scheduler and latency meter container** 
    Here's an example:

    ```yaml
    spec:
      serviceAccountName: latency-meter
      schedulerName: latency-aware-scheduler
      containers:
      - name: <container_name>
        image: <container_img>
        ports:
        - containerPort: 80
        env:
        - name: POD_NAME
          valueFrom:
            fieldRef:
              fieldPath: metadata.name
        - name: POD_NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace
      - name: latency-meter
        image: crischiaro/latency-meter:latest
        ports:
        - containerPort: 8080
        env:
        - name: POD_NAME
          valueFrom:
            fieldRef:
              fieldPath: metadata.name
        - name: POD_NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace
    ```


## Testing Steps

Navigate to `./tests/`:

- To start the test: `./start.sh`
- To stop the test: `./stop.sh`

### Simulating Latency with `tc`

The latency is artificially simulated using the `tc` command, which needs to be manually applied to the relevant nodes. For example:

```bash
tc qdisc show
```

```bash
tc qdisc add dev eth0 root netem delay 100ms
```
### Useful Commands
For testing network latency with tc:

```bash
tc qdisc show dev eth0
tc qdisc add dev eth0 root netem delay 100ms
tc qdisc delete dev eth0 root netem
```

For watching logs:

```bash
kubectl logs -f <pod-id>
```

For making requestes:

```bash
curl -H "X-Timestamp: $(date +%s%3N)" "http://<service_IP>/?id=123"
```
## Troubleshooting

### Multi-Cluster Issues
Ensure all clusters have the same podCIDRs. If you encounter issues in a multi-cluster setup, consult the [Liqo](https://docs.liqo.io/en/v0.9.4/index.html) support.

### Liqo Service Type
When you are installing Liqo on your cluster, please ensure the LoadBalancer has an external IP. If not, use the option `--service-type` and choose NodePort. For example:
```bash
 liqoctl install kubeadm --service-type NodePort --cluster-name paris
```

### CNI Configuration
Properly configure your CNIs based on the podCIDRs. If this is not done, the default will be set to `10.244.0.0/16`. Using CNIs other than those mentioned in the prerequisites may result in unpredictable behavior.

### Unexpected Latency Values
If you're experiencing latency values that don't align with your expectations (e.g., you set a 100ms delay through `tc` but observe 230ms), it could be due to the network path within the cluster. Depending on the network topology and CNI configuration, requests may go through other nodes in the cluster before reaching their destination. If these intermediary nodes introduce additional latency, it could be summed both at ingress and egress, resulting in a higher overall latency.

For more detailed troubleshooting guidelines, consult our [documentation](#).

