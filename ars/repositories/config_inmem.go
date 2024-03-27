package repositories

import "projekat/model"

type ConfigInMemRepository struct {
	configs map[string]model.Config
}

// todo: dodaj implementaciju metoda iz interfejsa ConfigRepository

func NewConfigInMemRepository() model.ConfigRepository {
	return ConfigInMemRepository{
		configs: make(map[string]model.Config),
	}
}
