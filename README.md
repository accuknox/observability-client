# observability-client

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

### Darwin

##### amd64

```bash
sudo curl https://storage.googleapis.com/kobserve/latest/darwin/amd64/accuknox -o accuknox && sudo chmod a+x accuknox | sudo mv accuknox /usr/bin/
```

##### arm64

```bash
sudo curl https://storage.googleapis.com/kobserve/latest/darwin/arm64/accuknox -o accuknox && sudo chmod a+x accuknox | sudo mv accuknox /usr/bin/
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
