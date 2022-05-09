# Observability Client

This provides observation of workloads (k8s pods, processes in VMs) at runtime. This delivers deep visibility into the workloads and their behavior with respect to the host environment and other services.

AccuKnox insights is part of a CLI tool that provides runtime visibility in an aggregated form (by pods, processes, workloads)

- Network: The L3, L4 and L7 connections with eBPF based observability. Ingress and egress..

- System: The files accessed, processes forked and network connections requested

## Installation

The following sections show how to install the accuknox CLI.

### Linux

##### amd64

```bash
sudo curl https://storage.googleapis.com/kobserve/latest/linux/amd64/accuknox -o accuknox && sudo chmod a+x accuknox | sudo mv accuknox /usr/bin/
```

##### arm64

```bash
sudo curl https://storage.googleapis.com/kobserve/latest/linux/arm64/accuknox -o accuknox && sudo chmod a+x accuknox | sudo mv accuknox /usr/bin/
```

## Usage

```
CLI Utility to help manage Observability
	
Observability is based on Cilium and Kubearmor Logs. 
Using this we can identify behavior of container, vm, pod, node at the network and system level.

Usage:
  accuknox [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  logs        To Get the logs based on Network or System

Flags:
  -h, --help   help for accuknox

Use "accuknox [command] --help" for more information about a command.

```
