// Go Substrate RPC Client (GSRPC) provides APIs and types around Polkadot and any Substrate-based chain RPC calls
//
// Copyright 2019 Centrifuge GmbH
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

package types

import (
	"fmt"
	"hash"
	"strings"

	"github.com/centrifuge/go-substrate-rpc-client/scale"
)

// Modelled after packages/types/src/Metadata/v8/Metadata.ts
type MetadataV8 struct {
	Modules []ModuleMetadataV8
}

func (m *MetadataV8) Decode(decoder scale.Decoder) error {
	err := decoder.Decode(&m.Modules)
	if err != nil {
		return err
	}
	return nil
}

func (m MetadataV8) Encode(encoder scale.Encoder) error {
	err := encoder.Encode(m.Modules)
	if err != nil {
		return err
	}
	return nil
}

func (m *MetadataV8) FindCallIndex(call string) (CallIndex, error) {
	s := strings.Split(call, ".")
	mi := uint8(0)
	for _, mod := range m.Modules {
		if !mod.HasEvents {
			continue
		}
		if string(mod.Name) != s[0] {
			mi++
			continue
		}
		for ci, f := range mod.Calls {
			if string(f.Name) == s[1] {
				return CallIndex{mi, uint8(ci)}, nil
			}
		}
		return CallIndex{}, fmt.Errorf("method %v not found within module %v for call %v", s[1], mod.Name, call)
	}
	return CallIndex{}, fmt.Errorf("module %v not found in metadata for call %v", s[0], call)
}

func (m *MetadataV8) FindEventNamesForEventID(eventID EventID) (Text, Text, error) {
	mi := uint8(0)
	for _, mod := range m.Modules {
		if !mod.HasEvents {
			continue
		}
		if mi != eventID[0] {
			mi++
			continue
		}
		if int(eventID[1]) >= len(mod.Events) {
			return "", "", fmt.Errorf("event index %v for module %v out of range", eventID[1], mod.Name)
		}
		return mod.Name, mod.Events[eventID[1]].Name, nil
	}
	return "", "", fmt.Errorf("module index %v out of range", eventID[0])
}

func (m *MetadataV8) FindStorageKeyHasher(module string, fn string) (hash.Hash, error) {
	for _, mod := range m.Modules {
		if !mod.HasStorage {
			continue
		}
		if string(mod.Storage.Prefix) != module {
			continue
		}
		for _, s := range mod.Storage.Items {
			if string(s.Name) != fn {
				continue
			}
			if s.Type.IsMap {
				return s.Type.AsMap.Hasher.HashFunc()
			}
			if s.Type.IsDoubleMap {
				return s.Type.AsDoubleMap.Hasher.HashFunc()
			}
			return nil, nil
		}
		return nil, fmt.Errorf("storage %v not found within module %v", fn, module)
	}
	return nil, fmt.Errorf("module %v not found in metadata", module)
}

type ModuleMetadataV8 struct {
	Name       Text
	HasStorage bool
	Storage    StorageMetadata
	HasCalls   bool
	Calls      []FunctionMetadataV4
	HasEvents  bool
	Events     []EventMetadataV4
	Constants  []ModuleConstantMetadataV6
	Errors     []ErrorMetadataV8
}

func (m *ModuleMetadataV8) Decode(decoder scale.Decoder) error {
	err := decoder.Decode(&m.Name)
	if err != nil {
		return err
	}

	err = decoder.Decode(&m.HasStorage)
	if err != nil {
		return err
	}

	if m.HasStorage {
		err = decoder.Decode(&m.Storage)
		if err != nil {
			return err
		}
	}

	err = decoder.Decode(&m.HasCalls)
	if err != nil {
		return err
	}

	if m.HasCalls {
		err = decoder.Decode(&m.Calls)
		if err != nil {
			return err
		}
	}

	err = decoder.Decode(&m.HasEvents)
	if err != nil {
		return err
	}

	if m.HasEvents {
		err = decoder.Decode(&m.Events)
		if err != nil {
			return err
		}
	}

	err = decoder.Decode(&m.Constants)
	if err != nil {
		return err
	}

	return decoder.Decode(&m.Errors)
}

func (m ModuleMetadataV8) Encode(encoder scale.Encoder) error {
	err := encoder.Encode(m.Name)
	if err != nil {
		return err
	}

	err = encoder.Encode(m.HasStorage)
	if err != nil {
		return err
	}

	if m.HasStorage {
		err = encoder.Encode(m.Storage)
		if err != nil {
			return err
		}
	}

	err = encoder.Encode(m.HasCalls)
	if err != nil {
		return err
	}

	if m.HasCalls {
		err = encoder.Encode(m.Calls)
		if err != nil {
			return err
		}
	}

	err = encoder.Encode(m.HasEvents)
	if err != nil {
		return err
	}

	if m.HasEvents {
		err = encoder.Encode(m.Events)
		if err != nil {
			return err
		}
	}

	err = encoder.Encode(m.Constants)
	if err != nil {
		return err
	}

	return encoder.Encode(m.Errors)
}

type ErrorMetadataV8 struct {
	Name          Text
	Documentation []Text
}