package service

import "awesomeProject/internal/store"

type Service struct {
	store *store.Store
}

func NewService(store *store.Store) *Service {
	return &Service{store: store}
}
