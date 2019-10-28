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
	"github.com/xmidt-org/kratos"
	"github.com/xmidt-org/webpa-common/logging"
	"go.uber.org/fx"
)

func StartUpstreamConnection(config Config, lc fx.Lifecycle, logger log.Logger) (kratos.Client, error) {
	queueConfig := kratos.QueueConfig{
		MaxWorkers: 5,
		Size:       100,
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
		Handlers:             []kratos.HandlerConfig{},
		HandlePingMiss: func() error {
			logging.Error(logger).Log(logging.MessageKey(), "Ping Miss")
			// TODO: handle reconnect
			return nil
		},
		ClientLogger: logger,
	})

	if err != nil {
		return nil, err
	}
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