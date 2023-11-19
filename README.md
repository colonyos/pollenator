# Introduction
> A pollinator is an organism that transfers pollen between flowers, aiding in plant reproduction and biodiversity. -- ChatGPT

* **Pollinator** is a technology comparable to [Platform-as-a-Service](https://en.wikipedia.org/wiki/Platform_as_a_service) (PaaS), 
much like [Heroku](https://www.heroku.com), 
but is based on a distributed brokering orchestration, and is primarily designed to run batch jobs. 
While Heroku mainly targets deployment of web applications 
in the Cloud, Pollinator is designed to simplify job execution across platforms, e.g executing AI computations on 
HPC (High-Performance Computing), Edge, or Kubernetes platforms. **Pollinator** is also designed to ensure uniform workload execution across these diverse platforms.

* **Pollinator** is based on [ColonyOS](https://colonyos.io). Using Meta-Orchestration, it can create and run batch jobs over a network of 
loosely-connected and geographically disperse so-called **Executors**. These Executors, after receiving jobs assignment from a Colonies server, transform the job 
instructions (so-called *function specifications*) into a format that is compatible with the underlying system they're connected to, 
whether it's Kubernetes (K8s) or Slurm. On HPC systems, Docker containers are for example automatically converted to Singularity containers.

* **Pollinator** significantly simplifies interactions with HPC or Kubernetes systems. For instance, it completely eliminates the need to manually
login to HPC nodes to run Slurm jobs. It seamlessly synchronizes and transfers data from the 
user's local filesystem to remote Executors, *offering the convenience of a local development environment while harnessing powerful supercomputers 
and cloud platforms.*

*  With **Pollinator**, users are no longer required to have in-depth knowledge of Slurm or Kubernetes systems, speeding up development and making 
powerful HPC systems available to more people.

## How does it work? 
Pollinator assumes the existance of the directories in the table below.  

| Directory    | Purpose                                         | Synchronizion strategy                                                      |
|--------------|-------------------------------------------------|-----------------------------------------------------------------------------|
| ./cfs/src    | Contains source code or binaries to be executed | Will be synchronized from local computer to remote executor before execution                           |
| ./cfs/data   | Datasets or other data is stored here           | Will be synchronized from local computer before execution, but not removed after job completion        |
| ./cfs/result | Produced data can be stored here.               | Will be synchronized from remote executor to local computer after execution                            |

When running a job, Pollinator does the following:
1. Synchronize the source, data, and result directories to the ColonyOS meta-filesystem.
2. Automatically generate and submit a ColonyOS function specification to a Colonies server.
2. If the job is assigned to an HPC Executor:
    1. Pull the Docker container to the HPC environment, and convert it to a Singularity container.
    3. Synchronize the source, data, and result directories to make project file accessible on the remote HPC environment.
    4. Generate a Slurm script to execute the Singularity container, including binding the source, data, and result directories to the container.
    5. Execute and monitor the Slurm job, including uploading all standard outputs and error logs to a Colonies server.
    6. Close the process by making a request to the Colonies server.
2. If the job is assigned to a remote Kubernetes Executor:
    1. Synchronize the source, data, and result directories to a shared Persistent Volume.
    2. Generate and deploy a K8s batch job. 
    3. Monitor the execution of the batch job, including uploading logs to a Colonies server.
    4. Close the process by making a request to the Colonies server.

## Example
Let's run some Python code at the [LUMI](https://www.lumi-supercomputer.eu) supercomputer in Finland. First, we need to generate a new Pollinator project.
The example assumed that ColonyOS credentials (private keys and S3 keys) and configurations are available as 
environmental variables. It also assumes the existance of an HPC Executor named **lumi-standard-hpcexecutor**, connected
to the LUMI standard CPU partition.

```console
export COLONYOS_DASHBOARD_URL="..."
export COLONIES_TLS="true"
export COLONIES_SERVER_HOST="..."
export COLONIES_SERVER_PORT="443"
export COLONIES_SERVER_ID="..."
export COLONIES_SERVER_PRVKEY="..."
export COLONIES_COLONY_ID="..."
export COLONIES_EXECUTOR_ID="..."
export COLONIES_EXECUTOR_PRVKEY="..."
export AWS_S3_ENDPOINT="..."
export AWS_S3_ACCESSKEY="..."
export AWS_S3_SECRETKEY="..."
export AWS_S3_REGION_KEY=""
export AWS_S3_BUCKET="..."
export AWS_S3_TLS="true"
export AWS_S3_SKIPVERIFY="false"
```

```console
mkdir lumi
cd lumi
```

### Create an new Pollintor project
First, we need to generate a new Pollinator project.

```console
pollinator new -e lumi-standard-hpcexecutor
```

This will generate a **project.yaml** file and the **src**, **data**, **result** directories. It also generates some sample code in Python. 
Note that any language can be supported by specifing another Docker image and setting the *cmd* option in the **project.yaml**.

```console
INFO[0000] Creating directory     Dir=./cfs/src
INFO[0000] Creating directory     Dir=./cfs/data
INFO[0000] Creating directory     Dir=./cfs/result
INFO[0000] Generating             Filename=./project.yaml
INFO[0000] Generating             Filename=./cfs/data/hello.txt
INFO[0000] Generating             Filename=./cfs/src/main.py
```

### Edit source code
Modify the **main.py** source file to print the hostname of the compute node.
```console
cat ./cfs/main.py
```

```python
import socket

hostname = socket.gethostname()
print("hostname:", hostname)
```

### Update resource specifications
To run code on 10 nodes at LUMI, we need to update the **project.yaml** file.
```yaml
projectid: 4e3f0f068cdb08f78ba3992bf5ccb9f5eb321125fa696c477eb387d37ab5c15f
conditions:
  executorType: lumi-small-hpcexecutor
  nodes: 1
  processesPerNode: 1
  cpu: 1000m
  mem: 1000Mi
  walltime: 600
  gpu:
    count: 0
    name: ""
environment:
  docker: python:3.12-rc-bookworm
  rebuildImage: false
  cmd: python3
  source: main.py
```

### Run the code
Now, run the code.
```console
pollinator run --follow
```

We can now see the hostname of all the host that run the code.
```console
Uploading main.py 100% [===============] (581 kB/s)
INFO[0000] Process submitted                             ProcessID=14510690e58fadc8b326fd4b57586be8f197ec800071845100ea455b1edaed8a
INFO[0000] Follow process at https://dashboard.colonyos.io/process?processid=14510690e58fadc8b326fd4b57586be8f197ec800071845100ea455b1edaed8a
hostname: nid001054
hostname: nid001067
hostname: nid001082
hostname: nid001078
hostname: nid001066
hostname: nid001083
hostname: nid001055
hostname: nid001070
hostname: nid001079
hostname: nid001071
INFO[0440] Process finished successfully                 ProcessID=14510690e58fadc8b326fd4b57586be8f197ec800071845100ea455b1edaed8a
```
