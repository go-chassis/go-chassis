package registry

import (
	"sync"

	"github.com/ServiceComb/go-chassis/core/common"
	cache "github.com/patrickmn/go-cache"
	"k8s.io/apimachinery/pkg/util/sets"
)

// Index interface provide set and get for microservice instance
type Index interface {
	Set(k string, x interface{})
	Get(k string, tags map[string]string) (interface{}, bool)
	Delete(k string)
	Items() map[string]*cache.Cache

	//TODO: tags set to defaultTag with version and app now
	// it can be set according to config
	GetTags() []string
	SetTags(tags sets.String)
}

// TODO: index can be rewrite in other ways
// now its benchmark shows it costs 300ns/op
type hashIndex struct {
	cache  map[string]*cache.Cache
	mux    sync.RWMutex
	labels []string
}

func newHashIndex(tags sets.String) Index {
	hi := &hashIndex{cache: make(map[string]*cache.Cache)}
	hi.SetTags(tags)
	return hi
}

func (hi *hashIndex) Items() map[string]*cache.Cache { return hi.cache }
func (hi *hashIndex) GetTags() []string              { return hi.labels }

func (hi *hashIndex) SetTags(tags sets.String) {
	labels := make([]string, 0, len(tags))
	for tag := range tags {
		labels = append(labels, tag)
	}
	hi.mux.Lock()
	hi.labels = labels
	hi.mux.Unlock()
}

func (hi *hashIndex) Set(k string, x interface{}) {
	microServiceInstances := x.([]*MicroServiceInstance)
	exist := make(map[string][]*MicroServiceInstance)

	for _, m := range microServiceInstances {
		ss := hi.getLabelkey(m.Metadata)
		if _, ok := exist[ss]; !ok {
			exist[ss] = make([]*MicroServiceInstance, 0)
		}
		exist[ss] = append(exist[ss], m)
	}

	hi.mux.Lock()
	defer hi.mux.Unlock()

	if _, ok := hi.cache[k]; !ok {
		hi.cache[k] = cache.New(0, 0)
	}
	for label, m := range exist {
		hi.cache[k].Set(label, m, 0)
	}
	// TODO: how to clear cache auto clear does not work
	// hi.autoClearCache(k, exist)
}

func (hi *hashIndex) autoClearCache(k string, exist map[string][]*MicroServiceInstance) {
	old := hi.cache[k].Items()
	delsets := sets.NewString()
	for key := range old {
		if _, ok := exist[key]; !ok {
			delsets.Insert(key)
		}
	}

	for in := range delsets {
		hi.cache[k].Delete(in)
	}
}

func (hi *hashIndex) Get(k string, tags map[string]string) (interface{}, bool) {
	s := hi.getLabelkey(tags)
	if len(hi.labels) != len(tags) {
		return nil, false
	}

	hi.mux.RLock()
	defer hi.mux.RUnlock()

	if ca, ok := hi.cache[k]; ok {
		return ca.Get(s)
	}
	return nil, false
}

func (hi *hashIndex) Delete(k string) {
	hi.mux.Lock()
	defer hi.mux.Unlock()
	delete(hi.cache, k)
}

func (hi *hashIndex) getLabelkey(tags map[string]string) (ss string) {
	for _, label := range hi.labels {
		if ss != "" {
			ss += ":" + tags[label]
			continue
		}
		ss += tags[label]
	}
	return
}

func (m *MicroServiceInstance) appID() string   { return m.Metadata[common.BuildinTagApp] }
func (m *MicroServiceInstance) version() string { return m.Metadata[common.BuildinTagVersion] }

// Has return whether microservice has tags
func (m *MicroServiceInstance) Has(tags map[string]string) bool {
	for k, v := range tags {
		if mt, ok := m.Metadata[k]; !ok || mt != v {
			return false
		}
	}
	return true
}

// WithAppID add app tag for microservice instance
func (m *MicroServiceInstance) WithAppID(v string) *MicroServiceInstance {
	m.Metadata[common.BuildinTagApp] = v
	return m
}
