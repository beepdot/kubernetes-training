### Agenda
In this session we will discuss briefly on the following topics:
- Virtual Box, VM, Vagrant
- Linux Namespaces and Cgroups
- Quick intro on Docker 
- Building docker images
- Running containers on docker
- Quick intro into docker swarm
- Setup and running workloads on docker swarm

#### Installing Virtual Box
```
sudo apt-get update
sudo apt-get install -y software-properties-common
wget -q https://www.virtualbox.org/download/oracle_vbox_2016.asc -O- | sudo apt-key add -
wget -q https://www.virtualbox.org/download/oracle_vbox.asc -O- | sudo apt-key add -
echo "deb [arch=amd64] https://download.virtualbox.org/virtualbox/debian  $(lsb_release -cs) contrib" | sudo tee /etc/apt/sources.list.d/virtualbox.list
sudo apt-get update
sudo apt-get install -y virtualbox-6.1
```

#### Installing Vagrant
```
curl -fsSL https://apt.releases.hashicorp.com/gpg | sudo apt-key add -
sudo apt-add-repository "deb [arch=amd64] https://apt.releases.hashicorp.com $(lsb_release -cs) main"
sudo apt-get update && sudo apt-get install -y vagrant
```

#### Starting an Ubuntu VM on Vagrant
```
mkdir -p ~/training/ubuntu-vagrant && cd ~/training/ubuntu-vagrant
vi Vagrantfile
# Copy and paste the contents from Part1/files/Vagrantfile
vagrant up
vagrant ssh
```

#### Linux Namespaces
```
ps -ef
unshare --user --pid --mount-proc --fork bash
ps -ef
lsns --output-all
```

#### Linux Cgroups
```
sudo apt install -y stress
systemd-run --scope --user -p MemoryLimit=100M stress --vm 1 --vm-bytes 1024M
systemd-run --scope --user -p CPUQuota=10% stress --cpu 1
# Use top command to view the process
# Press c and then press e after running top command
```

#### Installing Docker
```
sudo apt-get update
sudo apt-get install -y ca-certificates curl gnupg lsb-release
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt-get update
sudo apt-get install -y docker-ce docker-ce-cli containerd.io
sudo usermod -aG docker $(whoami)
sudo systemctl status docker.service --no-pager
sudo systemctl start docker.service
sudo systemctl enable docker.service
sudo docker run hello-world
```

#### Building and Pushing Docker Images
```
# Create a docker hub account first
docker build -t keshavprasad/python-server:0.0.1 -f ./DockerfilePython .
docker build -t keshavprasad/go-server:0.0.1 -f ./DockerfileGo .
docker images
docker logout
docker login
docker push keshavprasad/python-server:0.0.1
docker push keshavprasad/go-server:0.0.1
```

#### Running and Accessing Docker Containers
```
docker run --name pyserver -it -p 8000:8000 keshavprasad/python-server:0.0.1
docker run --name goserver -it -p 8090:8090 keshavprasad/go-server:0.0.1
```

#### Viewing Docker Containers Information
```
docker ps
docker ps -a
docker inspect CONTAINER_ID
```

#### Checking Docker Logs
```
docker logs CONTAINER_ID
docker logs CONTAINER_ID --tail 1 -f
```

#### Login to the containers
```
docker exec -it CONTAINER_ID sh
```

#### Copying Data to / from Containers
```
docker cp CONTAINER_ID:PATH LOCALPATH
docker cp LOCALPATH CONTAINER_ID:path
```

#### Stop / Restart / Kill / Remove Containers
```
# Use -d while running the containers to run in detached mode
docker stop CONTAINER_ID
docker start CONTAINER_ID
docker restart CONTAINER_ID
docker kill CONTAINER_ID
docker container rm CONTAINER_ID
```

#### Inter Container Communication
```
docker run -d --name pyserver -it -p 8000:8000 keshavprasad/python-server:0.0.1
docker run -d --name goserver -it -p 8090:8090 keshavprasad/go-server:0.0.1
docker inspect pyserver
# OR you can use jq to capture only IP address
# sudo apt install -y jq
# docker inspect pyserver | jq .[0].NetworkSettings.Networks.bridge.IPAddress
docker exec -it goserver sh
apk add curl
# Use the IP from the docker inspect pyserver command
curl 172.17.0.2:8000
curl pyserver:8000
docker network create shared-network
docker network connect shared-network goserver
docker network connect shared-network pyserver
curl pyserver:8000
# You can even use container id OR add a --hostname parameter while running the container
# To disconnect network use
# docker network disconnect shared-network
```

#### Restricting Resources and Metrics
```
docker run -it --rm --name ubuntu --cpus 0.5 --memory 100MB ubuntu:20.04 bash
apt update && apt install -y stress
docker stats ubuntu  # Run on another terminal
stress --vm 1 --vm-bytes 100M
stress --vm 1 --vm-bytes 80M
stress --cpu 1
# Also try with below
docker run --oom-kill-disable -it --rm --name ubuntu --cpus 1 --memory 100MB ubuntu:20.04 bash
```

#### Sharing namespaces (pid)
```
docker run -d --name pyserver -it -p 8000:8000 keshavprasad/python-server:0.0.1
docker run -d --name goserver -it -p 8090:8090 keshavprasad/go-server:0.0.1
docker exec -it goserver ps -ef
docker container rm goserver pyserver --force
docker run -d --name goserver -it -p 8090:8090 keshavprasad/go-server:0.0.1
docker run -d --pid container:pyserver --name goserver -it -p 8090:8090 keshavprasad/go-server:0.0.1
docker exec -it goserver ps -ef
docker container rm goserver pyserver --force
```

#### Persisting Container Data
```
mkdir /tmp/go-data
docker run -d --name goserver -v /tmp/go-data:/var/log -it -p 8090:8090 keshavprasad/go-server:0.0.1
curl localhost:8090/headers
cat /tmp/go-data/go-server.log
```

#### Environment Variables and Files
```
docker build -t keshavprasad/shell-app:0.0.1 -f ./DockerfileShell .
docker run -it --rm --env FIRSTNAME=Jhonny --env LASTNAME=Bravo keshavprasad/shell-app:0.0.1
docker run -it --rm --env-file envfile keshavprasad/shell-app:0.0.1
```

#### Creating Docker Swarm
```
docker swarm init
docker node ls
```

#### Running and Accessing Swarm Services
```
docker service create --name pyserver --replicas 1 --publish 8000:8000 keshavprasad/python-server:0.0.1
docker service create --name goserver --replicas 1 --publish 8090:8090 keshavprasad/go-server:0.0.1
docker service ls
curl localhost:8090/headers
```

#### Viewing Docker Service Information
```
docker service inspect goserver
docker service ps goserver --no-trunc
```

#### Inspecting Docker Service Logs
```
docker service logs goserver
docker service logs goserver --tail 1 -f
```

#### Scaling Docker Services
```
docker service scale goserver=2
docker service logs goserver --tail 1 -f
curl localhost:8090/headers # Run multiple times
docker service scale goserver=0
curl localhost:8090/headers
```

#### Inter Service Communication
```
docker service rm goserver pyserver
docker network create --driver swarm shared-network-swarm
docker service create --name pyserver --network shared-network-swarm --replicas 1 --publish 8000:8000 keshavprasad/python-server:0.0.1
docker service create --name goserver --network shared-network-swarm --replicas 1 --publish 8090:8090 keshavprasad/go-server:0.0.1
docker ps
docker exec -it CONTAINER_ID
apk add curl --no-cache
curl pyserver:8000
```

#### Persisting Docker Service Data
```
docker service rm goserver
rm -rf /tmp/go-data/*
ls /tmp/go-data/
docker service create --name goserver --network shared-network-swarm --replicas 1 --publish 8090:8090 --mount 'type=bind,src=/tmp/go-data,dst=/var/log' keshavprasad/go-server:0.0.1
curl localhost:8090/headers
cat /tmp/go-data/go-server.log
```

#### Restricting Resources and Metrics for Docker Services
```
docker ps
docker container update CONTAINER_ID --cpus 0.1 --memory 100M --memory-swap 200M
docker inspect CONTAINER_ID | grep -A1 "\"Memory\""
docker kill CONTAINER_ID
docker ps
docker inspect CONTAINER_ID | grep -A1 "\"Memory\""
docker service update goserver --limit-cpu 0.1 --limit-memory 100M
docker ps
docker inspect CONTAINER_ID | grep -A1 "\"Memory\""
docker stats
```

#### Rollback Service
```
docker service rm goserver
docker service create --name goserver --network shared-network-swarm --replicas 1 --publish 8090:8090 --mount 'type=bind,src=/tmp/go-data,dst=/var/log' keshavprasad/go-server:0.0.1
docker service update goserver --image keshavprasad/go-server:0.0.2
docker service ps goserver --no-trunc
docker service inspect goserver
docker service rollback goserver
docker service ps goserver --no-trunc
```

#### Docker Config, Environment Variables and Files
```
docker service create --name shellapp --network shared-network-swarm --replicas 1 keshavprasad/shell-app:0.0.1
docker service logs shellapp --tail 1 -f
docker service update shellapp --env-add FIRSTNAME=Jhonny
docker service logs shellapp --tail 1 -f
docker config create envconfig envfile
docker config ls
docker config inspect envconfig
docker config inspect envconfig | jq -r .[0].Spec.Data | base64 -d
docker service update shellapp --config-add envconfig
docker ps
docker exec -it CONTAINER_ID ls /
docker service logs shellapp --tail 1 -f
```

#### Docker Compose File
```
docker service rm goserver shellapp pyserver
sudo curl -L "https://github.com/docker/compose/releases/download/1.29.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
sudo ln -s /usr/local/bin/docker-compose /usr/bin/docker-compose
docker-compose --version
docker-compose -p example up
docker-compose -p example up -d
docker container rm files_pyserver_1 files_goserver_1 files_shellapp_1 --force
docker-compose scale shellapp=2
docker-compose rm -s
```

#### Docker Stack File
```
rm -rf /tmp/go-data/go-server.log
cat /tmp/go-data/go-server.log
docker stack deploy -c docker-stack.yml stack
docker service ls
docker service logs stack_shellapp  --tail 1 -f
curl localhost:8080/headers
cat /tmp/go-data/go-server.log
```