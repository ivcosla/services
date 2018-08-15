# Prerequisites

- Docker on each skyminer node.
- A [Docker Swarm][1] Cluster.

[1]: https://docs.docker.com/get-started/part4/#understanding-swarm-clusters

## Install Docker

For every skyminer node you need to ssh into:

    Instructions for ssh

Then install docker on it, if running armbian:

``` bash
# update and install
sudo apt-get update

sudo apt-get install \
    apt-transport-https \
    ca-certificates \
    curl \
    software-properties-common


# add docker official GPG key
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -

# verify
sudo apt-key fingerprint 0EBFCD88

# it should return:
pub   4096R/0EBFCD88 2017-02-22
      Key fingerprint = 9DC8 5822 9FC7 DD38 854A  E2D8 8D81 803C 0EBF CD88
uid                  Docker Release (CE deb) <docker@docker.com>
sub   4096R/F273FCD8 2017-02-22

# add the repository
sudo add-apt-repository \
   "deb [arch=armhf] https://download.docker.com/linux/ubuntu \
   $(lsb_release -cs) \
   stable"


# update and install docker-ce
sudo apt-get update
sudo apt-get install docker-ce
```

Now you may want to add your user to docker group so you can run docker without sudo and run it on startup

```bash
# Create docker group
sudo groupadd docker

# Add your user to the docker group
sudo usermod -aG docker $USER
# Now you should log out and in again for changes to take effect

# Make it available on startup
sudo usermod -aG docker $USER
sudo systemctl disable docker
```

## Install swarm

Docker swarm comes along with Docker-CE, and it allows to create a cluster of docker daemons and deploy services
on them.

For swarm to work there are some ports that need to be open, though on some systems they may be open
by default:

```bash
    TCP port 2377 for cluster management communications
    TCP and UDP port 7946 for communication among nodes
    UDP port 4789 for overlay network traffic
```

If your system doesn't open them by default, and assuming you are using armbian you can try to follow
the steps [here](https://wiki.debian.org/Uncomplicated%20Firewall%20%28ufw%29)

In order to get swarm ready we need to setup a cluster master node on one of the miner nodes, and then make every other
miner node to join the cluster as a worker.

First, on the node you want to be the master:

    docker swarm init

It will ask you which interface ip to advertise, choose the one you want. It will then prompt the command you must
run on the other nodes for them to join the cluster as workers:

```bash
  docker swarm join \
  --token <token> \
  <my ip>:<port>
```

Now log into every other miner node and join the cluster as a worker using the previous command.

## Deploying the docker stack

When every node has joined the cluster log into the master node and clone this repository somewhere. Now
you will launch the swarm stack from it. On the repository root run the following command:

    docker stack deploy --compose-file node-stack.yml nodes

Now you have commanded docker to deploy a stack of services called "nodes", but you can use another name
if you please.

Now you can check the status of the services (bear in mind that the first time you deploy the stack every docker
daemon has to pull the docker image that is mean to run, so it may take a while for them to be available):

    docker stack services nodes

You can also look at a service logs by calling:

    docker service logs <name of the service>

Or look where has been every replica of a service deployed:

    docker service ps <name of the service>

Finally if you want to remove every service deployed by the stack you can call:

    docker stack rm nodes
