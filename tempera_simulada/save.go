package main

import (
	"encoding/json"
	"os"
)

// Estrutura para a configuração da mochila
type MochilaConfig struct {
	Valores       []int `json:"valores"`
	Tamanhos      []int `json:"tamanhos"`
	TamanhoMaximo int   `json:"tamanho_maximo"`
}

const configFilename = "config_mochila.json"

// Função para carregar configuração do arquivo
func loadConfig() (MochilaConfig, error) {
	file, err := os.Open(configFilename)
	if err != nil {
		return MochilaConfig{}, err
	}
	defer file.Close()

	var config MochilaConfig
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return MochilaConfig{}, err
	}
	return config, nil
}

// Função para salvar configuração no arquivo
func saveConfig(config MochilaConfig) error {
	file, err := os.Create(configFilename)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(config)
}
