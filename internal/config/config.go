package config

type OperatingSystem int

const (
	Unknown OperatingSystem = iota
	Ubuntu18
	Ubuntu20
	Ubuntu22
	Redhat7
	Redhat8
	Redhat9
)
