package lib

import (
	"sync"

	"github.com/Psiphon-Labs/psiphon-tunnel-core/MobileLibrary/psi"

	xcore "github.com/xtls/xray-core/core"
	"github.com/xtls/xray-core/infra/conf"

	_ "github.com/xtls/xray-core/main/distro/all"
)

//
// ====================
// PSIPHON
// ====================
//

type PsiphonProvider interface {
	Notice(noticeJSON string)
	BindToDevice(fileDescriptor int) (string, error)
	HasNetworkConnectivity() int
	GetDNSServersAsString() string
	IPv6Synthesize(IPv4Addr string) string
	HasIPv6Route() int
	GetNetworkID() string
}

func StartPsiphon(
	config string,
	embeddedServerEntries string,
	encodedAuthorizations string,
	provider PsiphonProvider,
	isVpn bool,
	useIPv6Synthesizer bool,
	useHasIPv6Route bool,
) error {

	return psi.Start(
		config,
		embeddedServerEntries,
		encodedAuthorizations,
		provider,
		isVpn,
		useIPv6Synthesizer,
		useHasIPv6Route,
	)
}

func StopPsiphon() {
	psi.Stop()
}

//
// ====================
// XRAY
// ====================
//

var (
	xrayInstance *xcore.Instance
	xrayLock     sync.Mutex
)

func StartXray(configJSON string) string {

	xrayLock.Lock()
	defer xrayLock.Unlock()

	if xrayInstance != nil {
		return "xray already running"
	}

	config, err := conf.ParseJSONConfig(
		"json",
		[]byte(configJSON),
	)

	if err != nil {
		return err.Error()
	}

	coreConfig, err := config.Build()
	if err != nil {
		return err.Error()
	}

	instance, err := xcore.New(coreConfig)
	if err != nil {
		return err.Error()
	}

	err = instance.Start()
	if err != nil {
		return err.Error()
	}

	xrayInstance = instance

	return ""
}

func StopXray() string {

	xrayLock.Lock()
	defer xrayLock.Unlock()

	if xrayInstance == nil {
		return ""
	}

	err := xrayInstance.Close()

	xrayInstance = nil

	if err != nil {
		return err.Error()
	}

	return ""
}

func XrayVersion() string {
	return xcore.Version()
}
