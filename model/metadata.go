package model

import (
	"sync"

	"github.com/sitename/sitename/modules/json"
)

type ModelMetadata struct {
	Id              string    `json:"string"`
	Metadata        StringMap `json:"metadata"`
	PrivateMetadata StringMap `json:"private_metadata"`

	// mutex is used for safe access concurrenly
	mutex sync.RWMutex `db:"-"`
}

func (m *ModelMetadata) PreSave() {
	if m.Id == "" {
		m.Id = NewId()
	}
	if m.PrivateMetadata == nil {
		m.PrivateMetadata = make(map[string]string)
	}

	if m.Metadata == nil {
		m.Metadata = make(map[string]string)
	}
}

func (p *ModelMetadata) ToJson() string {
	b, _ := json.JSON.Marshal(p)
	return string(b)
}

type WhichMeta string

const (
	PrivateMetadata WhichMeta = "private"
	Metadata        WhichMeta = "metadata"
)

func (p *ModelMetadata) GetValueFromMeta(key string, defaultValue string, which WhichMeta) string {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	if which == PrivateMetadata { // get from private metadata
		if p.PrivateMetadata == nil {
			return defaultValue
		}

		if vl, ok := p.PrivateMetadata[key]; ok {
			return vl
		}
	} else if which == Metadata { // get from metadata
		if p.Metadata == nil {
			return defaultValue
		}

		if vl, ok := p.Metadata[key]; ok {
			return vl
		}
	}

	return defaultValue
}

func (p *ModelMetadata) StoreValueInMeta(items map[string]string, which WhichMeta) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if which == PrivateMetadata {
		if p.PrivateMetadata == nil {
			p.PrivateMetadata = make(map[string]string)
		}

		for k, vl := range items {
			p.PrivateMetadata[k] = vl
		}
	} else if which == Metadata {
		if p.Metadata == nil {
			p.Metadata = make(map[string]string)
		}

		for k, vl := range items {
			p.Metadata[k] = vl
		}
	}
}

func (p *ModelMetadata) ClearMeta(which WhichMeta) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if which == PrivateMetadata {
		for k := range p.PrivateMetadata {
			delete(p.PrivateMetadata, k)
		}
	} else if which == Metadata {
		for k := range p.Metadata {
			delete(p.Metadata, k)
		}
	}
}

func (p *ModelMetadata) DeleteValueFromMeta(key string, which WhichMeta) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if which == PrivateMetadata {
		delete(p.PrivateMetadata, key)
	} else if which == Metadata {
		delete(p.Metadata, key)
	}
}
