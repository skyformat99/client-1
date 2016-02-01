// Copyright (c) 2014 The gomqtt Authors. All rights reserved.
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

package client

import (
	"sync"

	"github.com/gomqtt/packet"
)

// A store is used to persists incoming or outgoing packets until they are
// successfully acknowledged by the other side.
type Store interface {
	// Open will open the store. Opening an already opened store must not
	// return an error.
	Open() error

	// Put will persist a packet to the store. An eventual existing packet with
	// the same id gets overwritten.
	Put(packet.Packet) error

	// Get will retrieve a packet from the store.
	Get(uint16) (packet.Packet, error)

	// Del will remove a packet from the store. Removing a nonexistent packet
	// must not return an error.
	Del(uint16) error

	// All will return all packets currently in the store.
	All() ([]packet.Packet, error)

	// Reset will completely reset the store.
	Reset() error

	// Close will close the store. Closing an already closed store must not
	// return an error.
	Close() error
}

type MemoryStore struct {
	store map[uint16]packet.Packet
	mutex sync.Mutex
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		store: make(map[uint16]packet.Packet),
	}
}

func (s *MemoryStore) Open() error {
	return nil
}

func (s *MemoryStore) Put(pkt packet.Packet) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	id, ok := packet.PacketID(pkt)
	if ok {
		s.store[id] = pkt
	}

	return nil
}

func (s *MemoryStore) Get(id uint16) (packet.Packet, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.store[id], nil
}

func (s *MemoryStore) Del(id uint16) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.store, id)

	return nil
}

func (s *MemoryStore) All() ([]packet.Packet, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	all := make([]packet.Packet, len(s.store))

	i := 0
	for _, pkt := range s.store {
		all[i] = pkt
		i++
	}

	return all, nil
}

func (s *MemoryStore) Reset() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.store = make(map[uint16]packet.Packet)

	return nil
}

func (s *MemoryStore) Close() error {
	return nil
}
