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
	"crypto/tls"
	"regexp"
	"net/http"
	"fmt"
	"github.com/xmidt-org/webpa-common/logging"
	"github.com/rs/xid"
	"io/ioutil"
	"github.com/go-kit/kit/log"
)

func validateMAC(mac string) bool {
	var (
		macAddrRE    *regexp.Regexp
		macAddrSlice []string
	)

	macAddrRE = regexp.MustCompile(`(([0-9A-Fa-f]{2}(?:[:-]?)){5}[0-9A-Fa-f]{2})|(([0-9A-Fa-f]{4}\.){2}[0-9A-Fa-f]{4})`)
	macAddrSlice = macAddrRE.FindAllString(mac, -1)

	if len(macAddrSlice) >= 1 {
		return true
	}

	return false
}

func getTLSConfig(mtlsCert string, mtlsKey string) (*tls.Config, error) {

	cert, err := tls.LoadX509KeyPair(mtlsCert, mtlsKey)
	if err != nil {
		return nil, err
	}
	return &tls.Config{
		Certificates: []tls.Certificate{cert},
	}, nil
}

func getThemisToken(config Config, tlsConfig *tls.Config, logger log.Logger) (token string, err error) {

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, _ := http.NewRequest("GET", config.AuthTokenURL, nil)
	req.Header.Set("X-Midt-Mac-Address", fmt.Sprintf("mac:%s", config.HardwareMAC))
	req.Header.Set("X-Midt-Serial-Number", config.HardwareSerialNumber)
	uuid := xid.New()
	req.Header.Set("X-Midt-Uuid", uuid.String())

	req.Header.Set("X-Midt-Partner-Id", config.PartnerID)
	req.Header.Set("X-Midt-Hardware-Model", config.HardwareModel)
	req.Header.Set("X-Midt-Hardware-Manufacturer", config.HardwareManufacturer)
	req.Header.Set("X-Midt-Firmware-Name", config.FirmwareName)
	req.Header.Set("X-Midt-Protocol", config.Protocol)
	req.Header.Set("X-Midt-Interface-Used", config.Interface)
	req.Header.Set("X-Midt-Last-Reboot-Reason", config.HardwareLastRebootReason)

	resp, err := client.Do(req)
	if err != nil {
		logging.Error(logger).Log(logging.MessageKey(), "Failed to get themis token.")
		return "", err
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logging.Error(logger).Log(logging.MessageKey(), "Failed to read themis http response.")
	}
	token = string(bodyBytes)
	return token, err
}

