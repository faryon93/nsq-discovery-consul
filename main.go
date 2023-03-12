package main

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
	"flag"
	"fmt"
	"github.com/faryon93/nsq-discovery-consul/consul"
	"github.com/faryon93/nsq-discovery-consul/nsq"
	"github.com/faryon93/util"
	capi "github.com/hashicorp/consul/api"
	"github.com/sirupsen/logrus"
	"os"
	"syscall"
	"time"
)

// ---------------------------------------------------------------------------------------
//  constants
// ---------------------------------------------------------------------------------------

const (
	nsqConfigRetryTimeout = 10 * time.Second
)

// ---------------------------------------------------------------------------------------
//  global variables
// ---------------------------------------------------------------------------------------

var (
	consulLookupdService string
	nsqConfigUrl         string
)

// ---------------------------------------------------------------------------------------
//  application entry
// ---------------------------------------------------------------------------------------

func main() {
	flag.StringVar(&consulLookupdService, "consul-lookupd-service", "", "consul service name of nsq lookupd tcp")
	flag.StringVar(&nsqConfigUrl, "nsq-conf-url", "", "http url of nsq/nsqadmin which should be managed")
	flag.Parse()

	if consulLookupdService == "" || nsqConfigUrl == "" {
		flag.Usage()
		os.Exit(-1)
	}

	logrus.Infoln("starting", GetAppVersion())

	cconsul, err := capi.NewClient(capi.DefaultConfig())
	if err != nil {
		logrus.Errorln("failed to connect to consul agent:", err)
		os.Exit(-1)
	}

	go consul.WatchService(cconsul, consulLookupdService, func(services []*capi.CatalogService) error {
		log := logrus.WithField("service", consulLookupdService)

		lookupdAddrs := make([]string, len(services))
		for i, svc := range services {
			lookupdAddrs[i] = fmt.Sprintf("%s:%d", svc.ServiceAddress, svc.ServicePort)
		}
		log.Infoln("lookupd endpoints changed:", lookupdAddrs)

		err := nsq.SetLookupdAddrs(nsqConfigUrl, lookupdAddrs)
		if err != nil {
			log.Errorf("failed to configure lookupd addresses (%s): retrying in %s", err, nsqConfigRetryTimeout)
			time.Sleep(nsqConfigRetryTimeout)
			return err
		}
		log.Infoln("nsq http configuration successful")

		return nil
	})

	util.WaitSignal(os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	logrus.Println("received SIGINT / SIGTERM going to shutdown")
}
