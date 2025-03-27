package genesis

import _ "embed"

var (
	//go:embed dpos/proxy.bin
	ProxyContract string

	//go:embed dpos/hub.bin
	HubContract string

	//go:embed dpos/validatorManager.bin
	ValidatorManagerContract string

	//go:embed dpos/vrf.bin
	VRFContract string
)
