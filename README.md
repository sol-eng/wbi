# (W)ork(b)ench (I)nstaller 

## Functionality
- Verify Workbench is installed & output version
- Ask which languages will be used
    - R is required
        - Scan and output of found R installations
        - If no /opt/R locations found, tell user about the Posit Installation recommendations
    - Python is optional
        - Scan and output of found Python installations
        - If no /opt/python locations found, tell user about the Posit Installation recommendations
- Ask if Jupyter should be installed
    - Ask which Python location Jupyter should be installed into
        - Install jupyter, jupyterlab, rsp_jupyter, rsconnect_jupyter and workbench_jupyterlab
        - Install and enable Jupyter Notebook extensions
- Ask if SSL should be setup
    - Ask for cert location
    - Ask for cert key location
- Ask for desired authentication method
    - Current choices are:
        - SAML
            - Ask for IdP metadata URL
            - Ask for IdP username attribute (default provided)
            - Link to IdP setup in Admin guide provided
        - OIDC
            - Link to IdP setup in Admin guide provided
            - Ask for IdP client-id
            - Ask for IdP client-secret
            - Ask for IdP issuer URL
            - Ask for IdP username claim (default provided)
        - AD/LDAP
            - Provide links to support articles for integrating Active Directory for the operating systems below (detected automatically)
                - Ubuntu
                - RHEL
        - PAM
            - Provide link for PAM customization
        - Other
            - Provide link for other authentication methods
- Ask for Workbench license key
    - Activate Workbench


## Assumptions
- Single server
- SQLite database
- Workbench has already been installed
- R has already been installed
- Python has already been installed
- Internet access (online installation)

## TODO
- Present user at the end with all known configuration info
- Write out configuration files
- Verify SSL certs
- Install Workbench from WBI
- Present possible R versions and allow user to install from WBI
- Present possible Python versions and allow user to install from WBI
- Provide a branch for HA setup
    - PostgreSQL details
    - NFS details