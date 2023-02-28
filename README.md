# (W)ork(b)ench (I)nstaller 

## Getting started

wbi is a CLI tool aimed at streamlining the installation and configuration of Posit Workbench. Please read through the assumptions below to ensure that your target architecture matches the current capabilities of wbi. Also note wbi requires to be run as root.

To get started run the setup command as root and follow the prompts:
```
sudo ./wbi setup
```

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
- Record, validate, and specify config for Posit Package Manger URL and R repo
- Automatically generate, validate and specify config for Posit Public Package Manager

### Posit Connect integration
- Record, validate and specify config for Posit Connect URL

### Configuration
- Inform of exact configuration changes needed
- Write changes to the correct config files