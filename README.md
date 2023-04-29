# (W)ork(b)ench (I)nstaller 

## Getting started

wbi is a CLI tool aimed at streamlining the installation and configuration of Posit Workbench. Please read through the assumptions below to ensure that your target architecture matches the current capabilities of wbi. Also note wbi requires to be run as root.

## Installation

### Apt Install

To install wbi on Ubuntu 18.04, 20.04 or 22.04 using apt install:
```
echo "deb [trusted=yes] https://apt.fury.io/wbi/ /" | sudo tee -a /etc/apt/sources.list.d/fury.list && sudo apt update && sudo apt install wbi

sudo wbi setup
```

### Yum Install

To install wbi on RHEL 7/CentOS 7, RHEL 8 or RHEL 9 using yum install:
```
sudo tee -a /etc/yum.repos.d/fury.repo > /dev/null <<EOT
[fury]
name=wbi
baseurl=https://yum.fury.io/wbi/
enabled=1
gpgcheck=0
EOT

sudo yum install -y wbi

sudo wbi setup
```

### Manual Install

Visit the [release page](https://github.com/sol-eng/wbi/releases) to find install instructions for the latest release.


## Usage

### Interactive Prompts

To get started run the setup command as root and follow the prompts:
```
sudo wbi setup
```

You can also pass the `--step` flag to begin at a certain spot in the interactive flow. For example to start at the Workbench installation step:
```
sudo wbi setup --step workbench
```

The following steps are valid options: start, prereqs, user, firewall, security, languages, r, python, workbench, license, jupyter, prodrivers, ssl, packagemanager, connect, restart, status, verify.

## Assumptions
- Single server
- SQLite database
- Internet access (online installation)

## Supported Operating Systems
- RHEL 9
- RHEL 8
- RHEL 7/CentOS 7
- Ubuntu 22.04
- Ubuntu 20.04
- Ubuntu 18.04

## Functionality

### R and Python installations
- Scan for existing R installations
- Install one or more R version from binary
- Symlinks R and Rscript
- Scan for existing Python installations
- Install one or more Python version from binary
- Adds a version of Python to PATH

### Posit Workbench installation
- Scan for an existing Workbench installation
- Install Workbench

### Licensing
- Detect if a license is already activated
- Activate a new license

### Jupyter
- Install Jupyter & extensions into a specified Python location
- Enable Jupyter Notebook extensions

### Posit Pro Drivers installation
- Scan for an existing Pro Drivers installation
- Install Posit Pro Drivers

### SSL
- Record and specify where to put cert and key paths

### Posit Package Manager integration
- Record, validate, and specify config for Posit Package Manger URL and R/Python repos
- Automatically generate, validate and specify config for Posit Public Package Manager

### Posit Connect integration
- Record, validate and specify config for Posit Connect URL

### Configuration
- Write changes to the correct config files