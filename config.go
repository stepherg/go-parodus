/*
 *  Copyright 2019 Comcast Cable Communications Management, LLC
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"go.uber.org/fx"
)

const (
	DEVICEID = "mac:%s"
)

const (
	HardwareModelKeyName            = "hw-model"
	HardwareSerialNumberKeyName     = "hw-serial-number"
	HardwareManufacturerKeyName     = "hw-manufacturer"
	HardwareMACKeyName              = "hw-mac"
	HardwareLastRebootReasonKeyName = "hw-last-reboot-reason"
	FirmwareNameKeyName             = "fw-name"
	BootTimeKeyName                 = "boot-time"
	PingTimeoutKeyName              = "webpa-ping-timeout"
	URLKeyName                      = "webpa-url"
	MaxBackoffKeyName               = "webpa-backoff-max"
	InterfaceKeyName                = "webpa-interface-used"
	ProtocolKeyName                 = ""
	UUIDKeyName                     = ""
	LocalURLKeyName                 = "parodus-local-url"
	PartnerIDKeyName                = "partner-id"
	CertPathKeyName                 = "ssl-cert-path"
	AuthTokenUrlKeyName             = "token-server-url"
	IPv4KeyName                     = "force-ipv4"
	IPv6KeyName                     = "force-ipv6"
	DebugKeyName                    = "debug"
	VersionKeyName                  = "version"
	ClientCertKeyName               = "client-cert-path"
	MtlsClientKeyKeyName            = "mtls-client-key-path"
	MtlsClientCertKeyName           = "mtls-client-cert-path"
)

const (
	XMIDTPathURL = "/api/v2/device"
)

func SetupFlagSet(fs *pflag.FlagSet) error {
	// Mark Device Info
	fs.StringP(HardwareModelKeyName, "m", "RPI-RDKB", "the hardware model name")
	fs.StringP(HardwareSerialNumberKeyName, "s", "0000000000000000", "the serial number")
	fs.StringP(HardwareManufacturerKeyName, "f", "Raspberry", "the device manufacturer")
	fs.StringP(HardwareMACKeyName, "d", "11:22:33:44:55:66", "the MAC address used to manage the device")
	fs.StringP(HardwareLastRebootReasonKeyName, "r", "unknown", "the last known reboot reason")
	fs.StringP(FirmwareNameKeyName, "n", "unknown", "firmware name and version currently running")
	fs.Int64P(BootTimeKeyName, "b", time.Now().Unix(), "the boot time in unix time")
	fs.StringP(URLKeyName, "u", "https://petasos:6400", "the petasos host")
	fs.IntP(MaxBackoffKeyName, "o", 60, "the maximum value in seconds for the backoff algorithm")
	fs.IntP(PingTimeoutKeyName, "t", 180, "the maximum time to wait between pings before assuming the upstream is broken")
	fs.StringP(InterfaceKeyName, "i", "eth0", "the device interface being used to connect to the cloud")
	fs.StringP(LocalURLKeyName, "l", "tcp://127.0.0.1:6666", "Parodus local server url")
	fs.StringP(PartnerIDKeyName, "p", "comcast", "partner ID of iot/gateway device")
	fs.StringP(CertPathKeyName, "c", "", "provide the certs for establishing secure upstream")
	fs.BoolP(IPv4KeyName, "4", false, "forcefully connect parodus to ipv4 address")
	fs.BoolP(IPv6KeyName, "6", false, "forcefully connect parodus to ipv6 address")
	fs.StringP(AuthTokenUrlKeyName, "U", "https://themis:6501/issue", "the themis host")
	fs.StringP(ClientCertKeyName, "P", "", "Client cert path")
	fs.StringP(MtlsClientKeyKeyName, "K", "", "mTLS client key path")
	fs.StringP(MtlsClientCertKeyName, "M", "", "mTLS client cert path")

	fs.BoolP(DebugKeyName, "", false, "enables debug logging")
	fs.BoolP(VersionKeyName, "v", false, "print version and exit")
	return nil
}

type Config struct {
	HardwareModel            string
	HardwareSerialNumber     string
	HardwareManufacturer     string
	HardwareMAC              string
	HardwareLastRebootReason string
	FirmwareName             string
	BootTime                 int64
	PingTimeout              int
	URL                      string
	MaxBackoff               int
	Interface                string
	Protocol                 string
	UUID                     string
	LocalURL                 string
	PartnerID                string
	CertPath                 string
	AuthTokenURL             string
	DeviceID                 string
	IPv4                     bool
	IPv6                     bool
	Debug                    bool
	PrintVersion             bool
	ClientCertPath           string
	MtlsClientKeyPath        string
	MtlsClientCertPath       string
}

type ConfigFlagIn struct {
	fx.In

	FlagSet *pflag.FlagSet

	PrintVersionFunc func()
}

func Provide(in ConfigFlagIn) (Config, error) {
	config := Config{}
	config.HardwareModel, _ = in.FlagSet.GetString(HardwareModelKeyName)
	config.HardwareSerialNumber, _ = in.FlagSet.GetString(HardwareSerialNumberKeyName)
	config.HardwareManufacturer, _ = in.FlagSet.GetString(HardwareManufacturerKeyName)
	config.HardwareMAC, _ = in.FlagSet.GetString(HardwareMACKeyName)
	config.HardwareMAC = strings.ToLower(config.HardwareMAC)

	config.HardwareLastRebootReason, _ = in.FlagSet.GetString(HardwareLastRebootReasonKeyName)
	config.FirmwareName, _ = in.FlagSet.GetString(FirmwareNameKeyName)
	config.BootTime, _ = in.FlagSet.GetInt64(BootTimeKeyName)
	config.URL, _ = in.FlagSet.GetString(URLKeyName)
	config.URL += XMIDTPathURL
	config.MaxBackoff, _ = in.FlagSet.GetInt(MaxBackoffKeyName)
	config.PingTimeout, _ = in.FlagSet.GetInt(PingTimeoutKeyName)
	config.Interface, _ = in.FlagSet.GetString(InterfaceKeyName)
	config.LocalURL, _ = in.FlagSet.GetString(LocalURLKeyName)
	config.PartnerID, _ = in.FlagSet.GetString(PartnerIDKeyName)
	config.CertPath, _ = in.FlagSet.GetString(CertPathKeyName)
	config.IPv4, _ = in.FlagSet.GetBool(IPv4KeyName)
	config.IPv6, _ = in.FlagSet.GetBool(IPv6KeyName)
	config.DeviceID = fmt.Sprintf(DEVICEID, strings.Replace(config.HardwareMAC, ":", "", -1))

	config.Debug, _ = in.FlagSet.GetBool(DebugKeyName)
	config.PrintVersion, _ = in.FlagSet.GetBool(VersionKeyName)

	config.AuthTokenURL, _ = in.FlagSet.GetString(AuthTokenUrlKeyName)
	config.ClientCertPath, _ = in.FlagSet.GetString(ClientCertKeyName)
	config.MtlsClientKeyPath, _ = in.FlagSet.GetString(MtlsClientKeyKeyName)
	config.MtlsClientCertPath, _ = in.FlagSet.GetString(MtlsClientCertKeyName)

	if config.PrintVersion {
		in.PrintVersionFunc()
	}

	return config, validateConfig(config)
}

func validateConfig(config Config) error {
	if config.HardwareModel == "" {
		return fmt.Errorf("%s must be set", HardwareModelKeyName)
	}
	if config.HardwareSerialNumber == "" {
		return fmt.Errorf("%s must be set", HardwareSerialNumberKeyName)
	}
	if config.HardwareManufacturer == "" {
		return fmt.Errorf("%s must be set", HardwareManufacturerKeyName)
	}
	if !validateMAC(config.HardwareMAC) {
		return fmt.Errorf("bad mac address: %s", config.HardwareMAC)
	}
	if config.URL == "" {
		return fmt.Errorf("%s must be set", URLKeyName)
	}
	return nil
}
