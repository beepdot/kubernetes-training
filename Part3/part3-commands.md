### Agenda
In this session we will discuss briefly on the following topics:
- Kubeconfig, Headless Services
- Labels, ReplicaSets, Deployments, DaemonSets, StatefulSets, Multi-Container Pods
- Probes, Resources Requests and Limits, Kubernetes Dashboard, Port Forwarding
- Jobs and Cronjobs
- Alternatives to docker (containerd, podman)


#### ReplicaSet
> Always prefer Deployment over ReplicaSet unless you want manual updates
```
k create -f replicaset.yaml
k get pods
k scale replicaset shell-rs --replicas=2
k get rs shell-rs -o jsonpath='{.spec.template.spec.containers[0].image}'
k patch rs shell-rs --type='json' -p='[{"op": "replace", "path": "/spec/template/spec/containers/0/image", "value":"keshavprasad/shell-app:0.0.2"}]'
k get rs shell-rs -o jsonpath='{.spec.template.spec.containers[0].image}'
k get pods
# Delete any one pod of the replica set and watch what happens
# You can path the pod now to fix the issue or patch the replica set and delete the pod
# k patch pods POD_ID --type='json' -p='[{"op": "replace", "path": "/spec/containers/0/image", "value":"keshavprasad/shell-app:0.0.1"}]'
k edit rs shell-rs
k delete pod POD_ID
```

#### Deployment
> Order of precedence

>Deployment -> ReplicaSet -> Pods -> Containers
```
k create deployment --image keshavprasad/shell-app:0.0.1 --replicas 2 shell-dp
k get pods
k get deployments.apps shell-dp -o yaml
k set image deployment shell-dp shell-app=keshavprasad/shell-app:0.0.1
k get pods
k rollout undo deployment shell-dp
k rollout status deployment shell-dp
k rollout history deployment shell-dp
# k set image deployment shell-dp shell-app=keshavprasad/shell-app:0.0.2 --record
# k set image deployment shell-dp shell-app=keshavprasad/shell-app:0.0.1 --record
# k patch deployment shell-dp --type='json' -p='[{"op": "replace", "path": "/spec/template/spec/containers/0/image", "value":"keshavprasad/shell-app:0.0.2"}]'
k get pods
k rollout restart deployment shell-dp
```

#### StatefulSet
```
k create -f statefulset.yaml
k get sts
k exec -it go-server-statefulset-0 -- sh -c 'apk add curl && curl localhost:8090/hello && curl localhost:8090/headers && cat /var/log/go-server.log'
k get pv,pvc
k scale statefulset go-server-statefulset --replicas=2
k get pv,pvc
k delete pod go-server-statefulset-0
k get pods
k exec -it go-server-statefulset-0 -- cat /var/log/go-server.log
```

#### Headless Service
```
k exec -it go-server-statefulset-0 -- sh -c 'nslookup pyserver.default.svc.cluster.local'
k exec -it go-server-statefulset-0 -- sh -c 'nslookup go-server-headless.default.svc.cluster.local'
```

#### DaemonSet
```
k create -f daemonset.yaml
k get pods -o wide
```

#### Multi-Container Pods
```
k create -f multi-pod-deployment.yaml
k get pods
k describe pod multi-pod-deployment-
k logs multi-pod-deployment-77bc7776fd-2lbhw -c shell-app --tail 10
k exec -it multi-pod-deployment-77bc7776fd-2lbhw -c python-server -- sh -c 'apk add curl && curl localhost:8000'
k logs multi-pod-deployment-77bc7776fd-2lbhw -c python-server --tail 10
k logs -l app=multi-pod-deployment -c shell-app --tail 1
```

#### Accessing the pods / applications
```
k get svc
k get svc pyserver -o yaml
k patch svc pyserver --type='json' -p='[{"op": "replace", "path": "/spec/type", "value":"NodePort"}]'
k get svc
curl $(minikube ip):$(k get svc pyserver -ojsonpath='{.spec.ports[0].nodePort}')
# Default NodePort range -30000-32767
k patch svc pyserver --type='json' -p='[{"op": "replace", "path": "/spec/type", "value":"LoadBalancer"}]'
# What are annotations? Check cloud provider annotations for private LB.
minikube addons enable metallb
k edit cm -n metallb-system config
# Add 192.168.49.200-192.168.49.21 under addresses as per below snipped

----
data:
  config: |
    address-pools:
    - name: default
      protocol: layer2
      addresses:
      - 192.168.49.200-192.168.49.210
kind: ConfigMap
----

k get svc
curl 192.168.49.200:8000
k patch svc pyserver --type='json' -p='[{"op": "replace", "path": "/spec/ports/0/port", "value":80}]'
k get svc
curl 192.168.49.200
```

#### Readiness and Liveness Probes
```
k create -f goserver-healthcheck.yaml
k expose deployment goserver-healthcheck --port 8090 --target-port 8090 --type LoadBalancer
curl 192.168.49.201:8090/hello
k get pods --watch
k describe pod goserver-healthcheck-
k expose pod pyserver --port 8000 --target-port 8000 --name readiness
k describe pod goserver-healthcheck-
curl 192.168.49.201:8090/hello
k get pods --watch
k expose pod pyserver --port 8000 --target-port 8000 --name liveness
k get pods --watch
k describe pod goserver-healthcheck-
curl 192.168.49.201:8090/hello
k logs -l app=goserver-healthcheck
```

#### Metrics Dashboard and Port Forward
```
k get pods -n kubernetes-dashboard
k get svc -n kubernetes-dashboard
minikube dashboard
k port-forward -n kubernetes-dashboard svc/kubernetes-dashboard 80
k port-forward -n kubernetes-dashboard svc/kubernetes-dashboard 8000:80
k port-forward -n kubernetes-dashboard kubernetes-dashboard-POD_ID 8000:9090
```

#### Jobs and Cronjobs
```
k create -f shellapp-job.yaml
k get jobs
k describe jobs.batch shellapp-job
k create cronjob shellapp-cronjob --image=keshavprasad/shell-app:0.0.2 --schedule="*/1 * * * *"
k get cronjobs
k describe cronjob shellapp-cronjob
```

#### Restricting Resources and Metrics
> Try it yourself activity
```
k run --image=ubuntu:20.04 ubuntu --command sleep 10000 --dry-run=client -o yaml
k create -f ubuntu-pod.yaml
k exec -it ubuntu -- bash
apt update && apt install -y stress
stress --cpu 1
# On another terminal run
k top pods
k top nodes
minikube dashboard
```

#### Deployment Rollout Strategies
> Try it yourself activity
```
Play around with deployment strategy types, maxSurge, maxUnavailable

  strategy:
    rollingUpdate:
      maxSurge: 25%
      maxUnavailable: 25%
    type: RollingUpdate

```

#### Podman
> Try it yourself activity
```
sudo apt-get -y install podman
sudo podman run -it --rm --env-file envfile keshavprasad/shell-app:0.0.1
sudo podman ps
```

#### Containerd and crictl
> Advanced use case (Try it yourself on vagrant to keep your system safe)
```
# Cannot have both dockerd and containerd engines at the same time
# Will try to provide a script to setup containerd in Vagrant
# More info https://kubernetes.io/docs/tasks/debug-application-cluster/crictl/
# In case of issue - unknown service runtime.v1alpha2.ImageService 
# Run sudo rm /etc/containerd/config.toml && sudo systemctl restart containerd
# Remove ipv6 ranges from /etc/cni/net.d/100-crio-bridge.conf if you get error failed to set bridge addr: could not add IP address to "cni0": permission denied

sudo crictl pull keshavprasad/shell-app:0.0.1
sudo crictl images
sudo crictl runp pod.config
sudo crictl create POD_ID_PREVIOUS_COMMAND container.config pod.config
sudo crictl ps
sudo crictl start CONTAINER_ID_FROM_PREVIOUS_COMMAND
sudo crictl ps
sudo crictl logs CONTAINER_ID_FROM_PREVIOUS_COMMAND
sudo crictl pods
sudo crictl stop CONTAINER_ID_FROM_PREVIOUS_COMMAND
sudo crictl stopp POD_ID
sudo crictl rm CONTAINER_ID_FROM_PREVIOUS_COMMAND
sudo crictl rmp POD_ID_PREVIOUS_COMMAND
```