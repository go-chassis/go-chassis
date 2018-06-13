package registry

import (
	"sort"
	"sync"

	"github.com/google/btree"
)

// Index interface provide put and get for microservice instance
type Index interface {
	Set(m []*MicroServiceInstance)
	Get(key string) []*MicroServiceInstance
}

// treeIndex build index for select micro service instance metadata
type treeIndex struct {
	sync.RWMutex
	tree  *btree.BTree
	label string
}

func newTreeIndex(label string) Index {
	return &treeIndex{label: label, tree: btree.New(32)}
}

func (ti *treeIndex) Set(microServiceInstances []*MicroServiceInstance) {
	// TODO: this separation according to label could move to keyindex
	ret := make(map[string][]*MicroServiceInstance)
	for _, m := range microServiceInstances {
		iKey, ok := m.Metadata[ti.label]
		if !ok {
			ret[iKey] = make([]*MicroServiceInstance, 0)
		}
		ret[iKey] = append(ret[iKey], m)
	}

	ti.RLock()
	defer ti.RUnlock()

	for labelKey, m := range ret {
		keyi := newKeyIndex(ti.label, labelKey)

		item := ti.tree.Get(keyi)
		if item == nil {
			keyi.Set(m)
			ti.tree.ReplaceOrInsert(keyi)
			continue
		}
		okeyi := item.(*keyIndex)
		okeyi.Set(m)
	}
}

func (ti *treeIndex) Get(key string) []*MicroServiceInstance {
	keyi := newKeyIndex(ti.label, key)
	ti.RLock()
	defer ti.RUnlock()
	var item btree.Item

	item = ti.tree.Get(keyi)
	if item == nil {
		return nil
	}
	keyi = item.(*keyIndex)
	return keyi.Get()
}

// keyIndex is the node type of btree
type keyIndex struct {
	label string
	key   string

	microIns []*MicroServiceInstance
}

func newKeyIndex(label, key string) *keyIndex {
	return &keyIndex{label: label, key: key, microIns: make([]*MicroServiceInstance, 0)}
}

func (k *keyIndex) Less(b btree.Item) bool { return k.key < b.(*keyIndex).key }

func (k *keyIndex) Get() []*MicroServiceInstance { return k.microIns }
func (k *keyIndex) Set(m []*MicroServiceInstance) {
	sort.Sort(MicroServiceInstances(m))
	k.microIns = m
}

// TODO: Put can replace set
func (k *keyIndex) Put(m *MicroServiceInstance) {
	cur := sort.Search(len(k.microIns), func(i int) bool { return !k.microIns[i].less(m) })
	if cur != len(k.microIns) && k.microIns[cur].equal(m) {
		return
	}
	ss := append(k.microIns[:cur], m)
	k.microIns = append(ss, k.microIns[cur:]...)
}

// MicroServiceInstances defines slice for point of MicroServiceInstance
type MicroServiceInstances []*MicroServiceInstance

func (m MicroServiceInstances) Len() int           { return len(m) }
func (m MicroServiceInstances) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }
func (m MicroServiceInstances) Less(i, j int) bool { return m[i].InstanceID < m[j].InstanceID }

func (m *MicroServiceInstance) less(b *MicroServiceInstance) bool { return m.InstanceID < b.InstanceID }
func (m *MicroServiceInstance) equal(b *MicroServiceInstance) bool {
	return m.InstanceID == b.InstanceID
}

func (m *MicroServiceInstance) appID() string   { return m.Metadata[App] }
func (m *MicroServiceInstance) version() string { return m.Metadata[Version] }
func (m *MicroServiceInstance) has(tags map[string]string) bool {
	for k, v := range tags {
		if mt, ok := m.Metadata[k]; !ok || mt != v {
			return false
		}
	}
	return true
}

// WithAppID add app tag for microservice instance
func (m *MicroServiceInstance) WithAppID(v string) *MicroServiceInstance {
	m.Metadata[App] = v
	return m
}

func multiInterSection(d [][]*MicroServiceInstance) []*MicroServiceInstance {
	switch len(d) {
	case 0:
		return nil
	case 1:
		return d[0]
	case 2:
		return twoInterSection(d[0], d[1])
	}
	s1 := multiInterSection(d[0 : len(d)/2])
	s2 := multiInterSection(d[len(d)/2:])
	return twoInterSection(s1, s2)
}

func twoInterSection(s1, s2 []*MicroServiceInstance) []*MicroServiceInstance {
	cap := len(s1)
	switch {
	case s2[len(s2)-1].less(s1[0]) || s1[len(s1)-1].less(s2[0]):
		return nil
	case len(s1) == 0:
		return s2
	case len(s2) == 0:
		return s1
	case len(s2) < len(s1):
		cap = len(s2)
	}
	ret := make([]*MicroServiceInstance, 0, cap)
	for i, j := 0, 0; i < len(s1) && j < len(s2); i++ {
		for j < len(s2) && s2[j].less(s1[i]) {
			j++
		}
		if j == len(s2) {
			break
		}
		if s1[i].equal(s2[j]) {
			ret = append(ret, s2[j])
			j++
		}
		for i < len(s1)-1 && s1[i].equal(s1[i+1]) {
			i++
		}
	}
	return ret
}
