package cache

import (
	"minik8s/pkg/apiobj"
	"sync"
)

type PodCache struct {
	mu   sync.RWMutex
	pods map[string]apiobj.Pod
}

type PodCacheInterface interface {
	GetAllPodFromCache() (map[string]*apiobj.Pod, error)
	DeletePodFromCache(podId string) error
	UpdatePodFromCache(*apiobj.Pod) error
}

func (p *PodCache) GetAllPodFromCache() (map[string]*apiobj.Pod, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	res := make(map[string]*apiobj.Pod)
	for podId, pod := range p.pods {
		// 创建新的对象以避免并发修改的问题
		newPod := pod
		res[podId] = &newPod
	}
	return res, nil
}

func (p *PodCache) DeletePodFromCache(podId string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.pods, podId)
	return nil
}

func (p *PodCache) UpdatePodFromCache(pod *apiobj.Pod) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	podId := pod.MetaData.UID
	p.pods[podId] = *pod
	return nil
}

func NewPodCache() *PodCache {
	return &PodCache{
		pods: make(map[string]apiobj.Pod),
	}
}
