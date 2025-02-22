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
	"context"
	"github.com/go-kit/kit/log"
	"github.com/stepherg/kratos"
	"github.com/xmidt-org/webpa-common/logging"
	"go.uber.org/fx"
	"time"
)

func StartUpstreamConnection(config Config, lc fx.Lifecycle, logger log.Logger) (kratos.Client, error) {
	queueConfig := kratos.QueueConfig{
		MaxWorkers: 5,
		Size:       100,
	}

	waitForNetwork(config.URL)

	tlsConfig, err := getTLSConfig(config.MtlsClientCertPath, config.MtlsClientKeyPath)
	if err != nil {
		logging.Error(logger).Log(logging.MessageKey(), "Failed to load mTLS config.")
	}

	token, err := getThemisToken(config, tlsConfig, logger)
	if err != nil {
		logging.Error(logger).Log(logging.MessageKey(), "Failed to get themis token.")
	}

	client, err := kratos.NewClient(kratos.ClientConfig{
		DeviceName:           config.DeviceID,
		FirmwareName:         config.FirmwareName,
		ModelName:            config.HardwareModel,
		Manufacturer:         config.HardwareManufacturer,
		DestinationURL:       config.URL,
		OutboundQueue:        queueConfig,
		WRPEncoderQueue:      queueConfig,
		WRPDecoderQueue:      queueConfig,
		HandlerRegistryQueue: queueConfig,
		HandleMsgQueue:       queueConfig,
		TlsConfig:            tlsConfig,
		Token:                token,
		Handlers:             []kratos.HandlerConfig{},
		HandlePingMiss: func() error {
			logging.Error(logger).Log(logging.MessageKey(), "Ping Miss")
			// TODO: handle reconnect
			return nil
		},
		ClientLogger: logger,
		PingConfig: kratos.PingConfig{
			PingWait:    time.Second * time.Duration(config.PingTimeout),
			MaxPingMiss: 3,
		},
	})
	if err != nil {
		logging.Error(logger).Log(logging.MessageKey(), "failed to create client", logging.ErrorKey(), err)
		if client != nil {
			closeErr := client.Close()
			logging.Info(logger).Log(logging.MessageKey(), "failed to close bad client", logging.ErrorKey(), closeErr)
		}
	}

	logging.Info(logger).Log(logging.MessageKey(), "kratos client created")
	lc.Append(fx.Hook{
		OnStart: func(context context.Context) error {
			return nil
		},
		OnStop: func(context context.Context) error {
			return client.Close()
		},
	})
	return client, nil
}
