// Copyright (C) 2019 Cisco Systems Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package types

import "github.com/jwhited/wgsd/vpplink/binapi/vppapi/interface_types"

type RxMode uint32

const (
	InvalidInterface = interface_types.InterfaceIndex(^uint32(0))
)

const (
	UnknownRxMode RxMode = 0
	Polling       RxMode = 1
	Interrupt     RxMode = 2
	Adaptative    RxMode = 3
	DefaultRxMode RxMode = 4

	AllQueues = ^uint32(0)
)

type VppInterfaceDetails struct {
	SwIfIndex uint32
	IsUp      bool
	Name      string
	Tag       string
	Type      string
}

type VppXDPInterface struct {
	SwIfIndex         uint32
	Name              string
	HostInterfaceName string
	RxQueueSize       int
	TxQueueSize       int
	NumRxQueues       int
}

func UnformatRxMode(str string) RxMode {
	switch str {
	case "interrupt":
		return Interrupt
	case "polling":
		return Polling
	case "adaptive":
		return Adaptative
	default:
		return UnknownRxMode
	}
}

func FormatRxMode(rxMode RxMode) string {
	switch rxMode {
	case Interrupt:
		return "interrupt"
	case Polling:
		return "polling"
	case Adaptative:
		return "adaptive"
	default:
		return "default"
	}
}
