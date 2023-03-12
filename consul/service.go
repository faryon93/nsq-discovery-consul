package consul

// nsq-discovery-consul
// Copyright (C) 2023 Maximilian Pachl

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

// ---------------------------------------------------------------------------------------
//  imports
// ---------------------------------------------------------------------------------------

import (
	capi "github.com/hashicorp/consul/api"
	"github.com/sirupsen/logrus"
)

// ---------------------------------------------------------------------------------------
//  public functions
// ---------------------------------------------------------------------------------------

// WatchService calls onChange when the given consul service is changed.
// The onChange handler is called right after registering the service watch.
// From there on the handler is called whenever the service has changed.
// The call to this function blocks indefinitely.
func WatchService(client *capi.Client, service string, onChange func([]*capi.CatalogService) error) {
	log := logrus.WithField("service", service)

	if client == nil {
		log.Error("cannot watch service without a client: exiting")
		return
	}

	opts := capi.QueryOptions{
		WaitIndex: 0,
	}
	for {
		services, meta, err := client.Catalog().Service(service, "", &opts)
		if err != nil {
			log.Errorln("failed to query consul service:", err.Error())
			continue
		}

		err = onChange(services)
		if err != nil {
			continue
		}

		opts.WaitIndex = meta.LastIndex
	}
}
