package nsq

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
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// ---------------------------------------------------------------------------------------
//  public functions
// ---------------------------------------------------------------------------------------

// SetLookupdAddrs sets the given lookupdAddrs in the given nsqd/nsqadmin instance
// via the specified nsqHttpAddr.
func SetLookupdAddrs(nsqHttpAddr string, lookupdAddrs []string) error {
	body, err := json.Marshal(lookupdAddrs)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", nsqHttpAddr, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%d / %s", resp.StatusCode, string(respBody))
	}

	return nil
}
