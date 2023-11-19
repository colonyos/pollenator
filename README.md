# Introduction
> A pollinator is an organism that transfers pollen between flowers, aiding in plant reproduction and biodiversity. -- ChatGPT

## Usage
```console
mkdir lumi
cd lumi
```

## Usage
First generate a new Pollinator project.

```console
pollinator new -e lumi-standard-hpcexecutor
```

```console
INFO[0000] Creating directory                            Dir=./cfs/src
INFO[0000] Creating directory                            Dir=./cfs/data
INFO[0000] Creating directory                            Dir=./cfs/result
INFO[0000] Generating                                    Filename=./project.yaml
INFO[0000] Generating                                    Filename=./cfs/data/hello.txt
INFO[0000] Generating                                    Filename=./cfs/src/main.py
```

Modify the main.py to print the hostname of the compute node.
```console
cat ./cfs/main.py
```

```python
import socket

print("hostname:", hostname)
```

Change nodes to 10. 
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

Now execute the code at LUMI. The command below will:
1. Generate and submit a ColonyOS function spec.
2. Once assigned to the a remote exectuor, the remote exector will.
  3. Convert the specified Docker container to a Singularity container, and upload the container to a LUMI login node.
  4. Synchronize the src, data and result directories so that the data becomes available at the filesystem at LUMI.
  5. Generate a Slurm script to run the uploaded Singularity container, also bind the srs, data, and result dirs to the Singularity container.
  6. Run and monitor the Slurm script, such as uploading all standard output/error the Colonies server.
  7. Close the process
   
```console
pollinator run
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
