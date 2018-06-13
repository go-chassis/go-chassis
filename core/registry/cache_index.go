package registry

import (
	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/hashicorp/go-version"
	cache "github.com/patrickmn/go-cache"
	"k8s.io/apimachinery/pkg/util/sets"
)

// noIndexCache return cache without index
type noIndexCache struct {
	latestV map[string]*version.Version
	cache   *cache.Cache
}

func newNoIndexCache() *noIndexCache {
	return &noIndexCache{
		cache:   cache.New(DefaultExpireTime, 0),
		latestV: map[string]*version.Version{},
	}
}

func (n *noIndexCache) SetIndexTags(tags sets.String) {}
func (n *noIndexCache) Items() map[string]cache.Item  { return n.cache.Items() }
func (n *noIndexCache) Delete(k string)               { n.cache.Delete(k); delete(n.latestV, k) }

func (n *noIndexCache) Set(k string, x interface{}) {
	latestV, _ := version.NewVersion("0.0.0")
	items, ok := x.([]*MicroServiceInstance)
	if !ok {
		return
	}
	for _, item := range items {
		v, _ := version.NewVersion(item.version())
		if v != nil && latestV.LessThan(v) {
			latestV = v
		}
	}
	// TODO: mutex should use
	n.latestV[k] = latestV
	n.cache.Set(k, x, 0)
}

func (n *noIndexCache) Get(k string, tags map[string]string) (interface{}, bool) {
	value, ok := n.cache.Get(k)
	if !ok {
		return nil, false
	}
	items, ok := value.([]*MicroServiceInstance)
	if !ok {
		return nil, false
	}
	n.setTagsBeforeQuery(k, tags)

	ret := make([]*MicroServiceInstance, 0, len(items))
	for _, item := range items {
		if item.has(tags) {
			ret = append(ret, item)
		}
	}
	if len(ret) == 0 {
		return nil, false
	}
	return ret, true
}

func (n *noIndexCache) setTagsBeforeQuery(k string, tags map[string]string) {
	if v, ok := tags[Version]; ok && v == common.LatestVersion {
		tags[Version] = n.latestV[k].String()
	}
}

// indexCache return cache with btree index
type indexCache struct {
	cache *noIndexCache

	// tags set according to archaius
	tags sets.String
	// indexers for each service name and tag key
	indexers map[string]map[string]Index
}

func newIndexCache() *indexCache {
	return &indexCache{
		indexers: make(map[string]map[string]Index, 0),
		tags:     sets.NewString(Version, App),
		cache:    newNoIndexCache(),
	}
}

// TODO: if tags rebuild, indexers should autoclear to remove
// index which is built from old tags
func (b *indexCache) SetIndexTags(tags sets.String) { b.tags = tags }
func (b *indexCache) Items() map[string]cache.Item  { return b.cache.Items() }
func (b *indexCache) Delete(k string)               { b.cache.Delete(k); delete(b.indexers, k) }

func (b *indexCache) Set(k string, x interface{}) {
	b.cache.Set(k, x)
	// no tags means no index need to be built
	if len(b.tags) == 0 {
		return
	}
	if _, ok := b.indexers[k]; !ok {
		b.indexers[k] = make(map[string]Index, 0)
	}
	for tag := range b.tags {
		if _, ok := b.indexers[k][tag]; !ok {
			b.indexers[k][tag] = newTreeIndex(tag)
		}
	}
	for tag := range b.indexers[k] {
		b.indexers[k][tag].Set(x.([]*MicroServiceInstance))
	}
}

func (b *indexCache) Get(k string, tags map[string]string) (interface{}, bool) {
	indexers, ok := b.indexers[k]
	if !ok || len(tags) == 0 {
		return b.cache.Get(k, tags)
	}

	b.cache.setTagsBeforeQuery(k, tags)
	ret := make([][]*MicroServiceInstance, 0, len(tags))
	for tag, value := range tags {
		ms := indexers[tag].Get(value)
		lager.Logger.Debugf("get instance(%v) from tag(%v) value(%v)", ms, tag, value)
		ret = append(ret, ms)
	}
	return multiInterSection(ret), true
}
