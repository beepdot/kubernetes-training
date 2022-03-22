### Agenda
In this session we will discuss briefly on the following topics:
- Quick Intro to Kubernetes
- Setting up Kubernetes cluster
- Running containers on Kubernetes
- Differences between docker swarm and kubernetes
- Kubernetes Components and Architecture, Kubeconfig
- ConfigMaps, Secrets, Environment Variables, Services, Headless Services
- Namespaces, Labels, Pods, ReplicaSets, Deployments, DaemonSets, StatefulSets, Multi-Container Pods
- Pod Logs, Probes, Resources Requests and Limits, Kubernetes Dashboard, Port Forwarding
- Jobs and Cronjobs
- Tools and Shortcuts for Kubernetes
- Alternatives to docker (containerd, podman)


#### Installing Kubernetes
```
curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
sudo install minikube-linux-amd64 /usr/local/bin/minikube
curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl
kubectl version --client
sudo sysctl net/netfilter/nf_conntrack_max=393216
sudo usermod -aG docker $USER && newgrp docker
sudo swapoff -a
minikube start --nodes 2 --driver=docker
minikube addons enable metrics-server
minikube addons enable dashboard
kubectl get po -A
source <(kubectl completion bash) # setup autocomplete in bash into the current shell, bash-completion package should be installed first.
echo "source <(kubectl completion bash)" >> ~/.bashrc # add autocomplete permanently to your bash shell.
alias k=kubectl
complete -F __start_kubectl k
k get po -A
minikube stop
minikube start
```

#### Running Kubernetes Pods
```
k run --image=keshavprasad/shell-app:0.0.1 shellapp
```

#### Viewing Pod Information
```
k get pods
```

#### Exposing Pod and Accessing the Applications
```
k describe pod shellapp
```

#### Inspecting Pod Logs
```
k logs shellapp
k logs shellapp --tail 1 -f
```

#### Login to the Pods
```
k exec -it shellapp -- sh
ls
hostname
exit
```

#### Stop / Restart / Kill Pods
```
k delete pod shellapp
# Stop and Restart we will see in subsequent commands
```

#### Running Pods with Replicas
```
k create deployment --image=keshavprasad/shell-app:0.0.1 shellapp
k get deployments
k get pods
k delete pod POD_ID --grace-period 0 --force
k get pods
k scale deployment shellapp --replicas=2
k get pods
k scale deployment shellapp --replicas=0
k get pods
```

#### Inter Pod Communication
```
k run --image=keshavprasad/python-server:0.0.1 pyserver
k run --image=keshavprasad/go-server:0.0.1 goserver
k get pods -o wide
k exec -it goserver -- apk add curl --no-cache
k exec -it goserver -- curl PYSERVER_POD_IP:8000
k expose pod pyserver --port 8000 --target-port 8000
k get svc
k exec -it goserver -- curl pyserver:8000
```

#### Namespaces
```
k create namespace myns
k get ns
k run -n myns --image=keshavprasad/go-server:0.0.1 goserver-myns
k get pods -n myns -o wide
k exec -n myns -it goserver-myns -- apk add curl --no-cache
k exec -n myns -it goserver-myns -- curl --connect-timeout 1 pyserver:8000
k exec -n myns -it goserver-myns -- curl pyserver.default:8000
```

#### Environment Variables and Configmaps
```
k create -f shellapp-envs.yaml
k logs shellappenv --tail 10
k create configmap --from-file envfile envfile-cm
k get cm envfile-cm
k get cm envfile-cm -o yaml
k create -f shellapp-configmap-as-mount.yaml
k logs shellappcm --tail 10
k create configmap envfromliteral --from-literal=FIRSTNAME=Popeye --from-literal=LASTNAME=Sailor
k get cm envfromliteral -o yaml
k create -f shellapp-configmap-as-env.yaml
k logs shellappcmenv --tail 10
k delete -f shellapp-configmap-as-env.yaml
k delete cm envfromliteral
k create configmap --from-env-file envfile envfromliteral
k get cm envfromliteral -o yaml
k create -f shellapp-configmap-as-env.yaml
k logs shellappcmenv --tail 10
```

#### Secrets
```
docker build -t keshavprasad/shell-app:0.0.2 -f ./DockerfileShell .
minikube image load keshavprasad/shell-app:0.0.2
k create secret generic literalsecret --from-literal=SECRETNAME=SHHH!
k create secret generic secretenvfile --from-env-file=secret-env-file
k create secret generic secretenvfile --from-env-file=secret-env-file
k create secret generic secretfile --from-file secrets=secret-env-file
k run secretpod --image=keshavprasad/shell-app:0.0.2
k get pods -o wide --watch
k run secretpod --image=keshavprasad/shell-app:0.0.2 --restart OnFailure
k delete pod secretpod --force --grace-period 0
k create -f secretpod.yaml
k get secrets secretenvfile -o yaml
k get secrets literalsecret -o jsonpath='{.data.SECRETNAME}' | base64 -d && echo
```

#### Manifest File
```
ubuntu-pod.yaml shellapp-envs.yaml shellapp-configmap-as-mount.yaml are all called manifests files
```
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

#### Few good tools
```
krew - https://krew.sigs.k8s.io/
stern - https://github.com/wercker/stern
k9s - https://github.com/derailed/k9s
kubens and kubectx - https://github.com/ahmetb/kubectx#installation

(
  set -x; cd "$(mktemp -d)" &&
  OS="$(uname | tr '[:upper:]' '[:lower:]')" &&
  ARCH="$(uname -m | sed -e 's/x86_64/amd64/' -e 's/\(arm\)\(64\)\?.*/\1\2/' -e 's/aarch64$/arm64/')" &&
  KREW="krew-${OS}_${ARCH}" &&
  curl -fsSLO "https://github.com/kubernetes-sigs/krew/releases/latest/download/${KREW}.tar.gz" &&
  tar zxvf "${KREW}.tar.gz" &&
  ./"${KREW}" install krew
)

export PATH="${KREW_ROOT:-$HOME/.krew}/bin:$PATH"
kubectl krew install ctx
kubectl krew install ns

git clone https://github.com/ahmetb/kubectx.git ~/.kubectx
COMPDIR=$(pkg-config --variable=completionsdir bash-completion)
ln -sf ~/.kubectx/completion/kubens.bash $COMPDIR/kubens
ln -sf ~/.kubectx/completion/kubectx.bash $COMPDIR/kubectx
# Fuzzy Search
sudo apt-get install -y fzf
cat << EOF >> ~/.bashrc

# Krew
export PATH="${KREW_ROOT:-$HOME/.krew}/bin:$PATH"

#kubectx and kubens
export PATH=~/.kubectx:\$PATH
alias kns=kubens
alias kctx=kubectx
EOF

git clone https://github.com/jonmosco/kube-ps1.git
# Add the below to your bashrc

#kube-ps1
source /path/to/kube-ps1/kube-ps1.sh
PS1='[\u@\h \w $(kube_ps1)]\$ '

# Default Ubuntu Style with Kube PS1
PS1='\[\033[01;32m\]\u@\h:\[\033[01;34m\]\w\[\033[00m\] $(kube_ps1)]\$ '

kubeon
kubeoff

wget https://github.com/wercker/stern/releases/download/1.11.0/stern_linux_amd64
wget https://github.com/derailed/k9s/releases/download/v0.25.18/k9s_Linux_x86_64.tar.gz

mv stern_linux_amd64 stern
chmod +x stern
sudo mv stern /usr/local/bin/
stern shell --tail 1

tar -xf k9s_Linux_x86_64.tar.gz
sudo mv k9s /usr/local/bin/
k9s
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

#### Restricting Resources and Metrics
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
> Advanced use case (Try it on vagrant to be safe)
```
# Cannot have both dockerd and containerd engines at the same time
# A script will be given later to setup containerd in vagrant
# Try to install in Vagrant
# https://kubernetes.io/docs/tasks/debug-application-cluster/crictl/
sudo crictl pull keshavprasad/shell-app:0.0.1
sudo crictl images
sudo crictl runp pod.config
sudo crictl create 4c2fd02afacb6715ba58cad33cd2322d848deb214782089ca59fc1915717d0ba container.config pod.config
sudo crictl ps
sudo crictl start CONTAINER_ID
sudo crictl ps
sudo crictl logs CONTAINER_ID
sudo crictl pods
sudo crictl stop CONTAINER_ID
sudo crictl stopp POD_ID
sudo crictl rm CONTAINER_ID
sudo crictl rmp POD_ID
```