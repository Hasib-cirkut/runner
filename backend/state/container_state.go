package state

import (
	"strconv"
	"sync"
)

type ContainerState struct {
	mu                  sync.Mutex
	availableContainers []string
	inUseContainers     map[string]struct{}
	poolSize            int
}

func NewContainerState() *ContainerState {
	return &ContainerState{
		availableContainers: []string{},
		inUseContainers:     make(map[string]struct{}),
		poolSize:            10,
	}
}

func (s *ContainerState) CreateContainers() ([]string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := 0; i < s.poolSize; i++ {
		containerID := "container_" + strconv.Itoa(i)
		s.availableContainers = append(s.availableContainers, containerID)
	}

	for i := 0; i < s.poolSize; i++ {
		print("container created: " + s.availableContainers[i] + "\n")
	}

	return s.availableContainers, true
}

func (s *ContainerState) GetContainer() (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.availableContainers) == 0 {
		return "", false
	}

	value, rest := s.availableContainers[0], s.availableContainers[1:]

	s.availableContainers = rest
	s.inUseContainers[value] = struct{}{}

	return value, true
}

func (s *ContainerState) ReleaseContainer(containerID string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.inUseContainers[containerID]

	if !exists {
		return "", false
	}

	s.availableContainers = append(s.availableContainers, containerID)

	return containerID, true
}
