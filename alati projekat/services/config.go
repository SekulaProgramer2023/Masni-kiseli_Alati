package services

import (
	"alati_projekat/model"
	"fmt"
)

type ConfigService struct {
	repo model.ConfigRepository
}

func NewConfigService(repo model.ConfigRepository) ConfigService {
	return ConfigService{
		repo: repo,
	}
}

func (s ConfigService) Hello() {
	fmt.Println("hello from config service")
}

func (s ConfigService) Add(config model.Config) {
	s.repo.Add(config)
}

func (s ConfigService) Get(name string, version int) (model.Config, error) {
	return s.repo.Get(name, version)
}

func (s ConfigService) Delete(name string, version int) error {
	s.repo.Delete(name, version)

	err := s.repo.Delete(name, version)
	if err != nil {

		return err
	}

	return nil
}

// todo: implementiraj metode za dodavanje, brisanje, dobavljanje itd.
