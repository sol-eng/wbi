package config

// Prints the WBConfig configuration struct information to the console
func (WBConfig *WBConfig) DetectConfigChange() bool {
	anyChange := false
	if WBConfig.PythonConfig.JupyterPath != "" {
		anyChange = true
	}
	if WBConfig.SSLConfig.UseSSL {
		anyChange = true
	}
	if WBConfig.AuthConfig.Using {
		anyChange = true
	}
	if WBConfig.PackageManagerConfig.Using {
		anyChange = true
	}
	if WBConfig.ConnectConfig.Using {
		anyChange = true
	}
	return anyChange
}
