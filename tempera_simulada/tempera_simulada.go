package main

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"
)

// Estrutura para salvar a solução em JSON
type Solucao struct {
	Iteracao      int     `json:"iteracao"`
	ValorAtual    int     `json:"valor_atual"`
	Temperatura   float64 `json:"temperatura"`
	MelhorSolucao []int   `json:"melhor_solucao"`
}

// valorTotalTamanho calcula o valor total e o tamanho total dos itens selecionados.
func valorTotalTamanho(arranjo []int, valores []int, tamanhos []int, tamanhoMaximo int) (int, int) {
	valor, tamanho := 0, 0
	for i, v := range arranjo {
		if v == 1 {
			valor += valores[i]
			tamanho += tamanhos[i]
		}
	}
	if tamanho > tamanhoMaximo {
		valor = 0
	}
	return valor, tamanho
}

// adjacente gera um arranjo adjacente alterando aleatoriamente um item.
func adjacente(arranjo []int, aleatorio *rand.Rand) []int {
	n := len(arranjo)
	resultado := make([]int, n)
	copy(resultado, arranjo)
	i := aleatorio.Intn(n)
	resultado[i] = 1 - resultado[i]
	return resultado
}

// resfriar calcula a nova temperatura.
func resfriar(temperaturaAtual, fatorResfriamento float64) float64 {
	novaTemperatura := temperaturaAtual * fatorResfriamento
	if novaTemperatura < 0.0001 {
		novaTemperatura = 0
	}
	return novaTemperatura
}

// Função principal de têmpera simulada
func temperaSimulada(nItens int, aleatorio *rand.Rand, valores, tamanhos []int, tamanhoMaximo int, temperaturaInicial, fatorResfriamento float64) ([]int, []Solucao) {
	temperaturaAtual := temperaturaInicial
	solucao := make([]int, nItens)
	melhorSolucao := solucao
	valorAtual, _ := valorTotalTamanho(solucao, valores, tamanhos, tamanhoMaximo)
	valorMelhorSolucao := valorAtual

	fmt.Printf("\n%-15s %-15s %-15s\n", "Iteração", "Valor Atual", "Temperatura")
	fmt.Println("-----------------------------------------")

	// Lista para armazenar dados para o JSON
	historico := []Solucao{}

	for iteracao := 0; temperaturaAtual > 0; iteracao++ {
		arranjoAdjacente := adjacente(solucao, aleatorio)
		valorAdjacente, _ := valorTotalTamanho(arranjoAdjacente, valores, tamanhos, tamanhoMaximo)

		delta_e := valorAdjacente - valorAtual

		if delta_e > 0 || math.Exp(float64(delta_e)/temperaturaAtual) > rand.Float64() {
			solucao = arranjoAdjacente
			valorAtual = valorAdjacente
			if valorAtual > valorMelhorSolucao {
				melhorSolucao = solucao
				valorMelhorSolucao = valorAtual
			}
		}

		if iteracao%10 == 0 {
			fmt.Printf("%-15d %-15d %-15.2f\n", iteracao, valorAtual, temperaturaAtual)
		}

		// Adicionar dados ao histórico
		historico = append(historico, Solucao{
			Iteracao:      iteracao,
			ValorAtual:    valorAtual,
			Temperatura:   temperaturaAtual,
			MelhorSolucao: melhorSolucao,
		})

		temperaturaAtual = resfriar(temperaturaAtual, fatorResfriamento)
	}

	return melhorSolucao, historico
}

// Função para salvar o histórico no arquivo JSON
func salvarJSON(nomeArquivo string, dados []Solucao) error {
	arquivo, err := os.Create(nomeArquivo)
	if err != nil {
		return err
	}
	defer arquivo.Close()

	encoder := json.NewEncoder(arquivo)
	encoder.SetIndent("", "  ") // Formatar o JSON com indentação
	return encoder.Encode(dados)
}

func printTabela(valores, tamanhos, arranjo []int) {
	fmt.Printf("\n%-10s %-10s %-10s\n", "Item", "Valor", "Tamanho")
	fmt.Println(strings.Repeat("-", 30))
	for i := range valores {
		incluido := "Não"
		if arranjo[i] == 1 {
			incluido = "Sim"
		}
		fmt.Printf("%-10d %-10d %-10d %-10s\n", i+1, valores[i], tamanhos[i], incluido)
	}
}

func main() {
	valores := []int{95, 75, 60, 85, 40, 120, 30, 65, 50, 90}
	tamanhos := []int{50, 40, 30, 55, 25, 60, 35, 45, 40, 50}
	tamanhoMaximo := 300

	aleatorio := rand.New(rand.NewSource(time.Now().UnixNano()))
	temperaturaInicial := 1000.0
	fatorResfriamento := 0.998

	fmt.Printf("Tamanho máximo da mochila = %d\n", tamanhoMaximo)
	fmt.Printf("Temperatura inicial = %.1f\n", temperaturaInicial)
	fmt.Printf("Fator de Resfriamento = %.2f\n", fatorResfriamento)

	fmt.Println("\nInício da demonstração de têmpera simulada com mochila")
	fmt.Println("\nValores e tamanhos dos itens:")
	printTabela(valores, tamanhos, make([]int, len(valores)))

	solucao, historico := temperaSimulada(10, aleatorio, valores, tamanhos, tamanhoMaximo, temperaturaInicial, fatorResfriamento)

	fmt.Println("\n\nSolução encontrada:")
	printTabela(valores, tamanhos, solucao)

	valor, tamanho := valorTotalTamanho(solucao, valores, tamanhos, tamanhoMaximo)
	fmt.Printf("\nValor total da solução = %d\n", valor)
	fmt.Printf("Tamanho total da solução = %d\n", tamanho)

	// Salvar histórico em JSON
	if err := salvarJSON("historico.json", historico); err != nil {
		fmt.Println("Erro ao salvar JSON:", err)
	} else {
		fmt.Println("Histórico salvo em 'historico.json'")
	}
}
