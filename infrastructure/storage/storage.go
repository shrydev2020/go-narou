package storage

import "narou/config"

type storageManager struct {
	dist    string
	subDist string
}

type Manager interface {
	GetDist() string
	GetSubDist() string
}

func NewManager(configure config.IConfigure) Manager {
	dist, subDist := configure.GetStorageConfig()
	return &storageManager{dist: dist, subDist: subDist}
}

func (s *storageManager) GetDist() string {
	return s.dist
}

func (s *storageManager) GetSubDist() string {
	return s.subDist
}
