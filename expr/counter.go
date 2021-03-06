// Copyright 2018 Google LLC. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package expr

import (
	"github.com/google/nftables/binaryutil"
	"github.com/mdlayher/netlink"
	"golang.org/x/sys/unix"
)

type Counter struct {
	Bytes   uint64
	Packets uint64
}

func (e *Counter) marshal() ([]byte, error) {
	data, err := netlink.MarshalAttributes([]netlink.Attribute{
		{Type: unix.NFTA_COUNTER_BYTES, Data: binaryutil.BigEndian.PutUint64(e.Bytes)},
		{Type: unix.NFTA_COUNTER_PACKETS, Data: binaryutil.BigEndian.PutUint64(e.Packets)},
	})
	if err != nil {
		return nil, err
	}

	return netlink.MarshalAttributes([]netlink.Attribute{
		{Type: unix.NFTA_EXPR_NAME, Data: []byte("counter\x00")},
		{Type: unix.NLA_F_NESTED | unix.NFTA_EXPR_DATA, Data: data},
	})
}

func (e *Counter) unmarshal(data []byte) error {
	attrs, err := netlink.UnmarshalAttributes(data)
	if err != nil {
		return err
	}
	for _, attr := range attrs {
		switch attr.Type {
		case unix.NFTA_COUNTER_BYTES:
			e.Bytes = binaryutil.BigEndian.Uint64(attr.Data)
		case unix.NFTA_COUNTER_PACKETS:
			e.Packets = binaryutil.BigEndian.Uint64(attr.Data)
		}
	}

	return nil
}
