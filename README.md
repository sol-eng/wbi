# (W)ork(b)ench (I)nstaller 

## Getting started

wbi is a CLI tool aimed at streamlining the installation and configuration of Posit Workbench. Please read through the assumptions below to ensure that your target architecture matches the current capabilities of wbi. Also note wbi requires to be run as root.

### Interactive Setup

To get started run the setup command as root and follow the prompts:
```
sudo ./wbi setup
```

You can also pass the `--step` flag to begin at a certain spot in the interactive flow. For example to start at Workbench installation:
```
sudo ./wbi setup --step workbench
```

The following steps are valid options: start, prereqs, firewall, security, languages, r, python, workbench, license, jupyter, prodrivers, ssl, auth, packagemanager, connect, restart, status, verify.

## Assumptions
- Single server
- SQLite database
- Internet access (online installation)

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

### Authentication
- Provide information about PAM and AD/LDAP setup steps
- Record and specify where to put values for SAML and OIDC SSO setups

### Posit Package Manager integration
- Record, validate, and specify config for Posit Package Manger URL and R/Python repos
- Automatically generate, validate and specify config for Posit Public Package Manager

### Posit Connect integration
- Record, validate and specify config for Posit Connect URL

### Configuration
- Inform of exact configuration changes needed
- Write changes to the correct config files