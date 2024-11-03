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

type IterationData struct {
	Iteracao      int     `json:"iteracao"`
	ValorAtual    int     `json:"valor_atual"`
	Temperatura   float64 `json:"temperatura"`
	MelhorSolucao []int   `json:"melhor_solucao"`
	Execucao      int     `json:"execucao"`
}

func generateKnapsackConfig(n int) ([]int, []int, int) {
	aleatorio := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Inicializa arrays de valores e tamanhos
	values := make([]int, n)
	sizes := make([]int, n)

	// Define os valores e tamanhos aleatórios para cada item
	for i := 0; i < n; i++ {
		values[i] = aleatorio.Intn(100) + 1 // Valor entre 1 e 100
		sizes[i] = aleatorio.Intn(50) + 1   // Tamanho entre 1 e 50
	}

	// Calcula a soma dos tamanhos para definir o limite da mochila
	totalSize := 0
	for _, size := range sizes {
		totalSize += size
	}

	// Define a capacidade máxima como uma fração (ex.: 50%) da soma total dos tamanhos
	maxWeight := int(float64(totalSize) * 0.5)

	return values, sizes, maxWeight
}

// valorTotalTamanho calcula o valor total e o tamanho total dos itens selecionados no arranjo.
// Retorna zero para o valor se o tamanho total exceder o tamanho máximo permitido.
func valorTotalTamanho(arranjo []int, valores []int, tamanhos []int, tamanhoMaximo int) (int, int) {
	valor := 0
	tamanho := 0
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

// adjacente gera um arranjo adjacente alterando aleatoriamente um item do arranjo atual.
func adjacente(arranjo []int, aleatorio *rand.Rand) []int {
	n := len(arranjo)
	resultado := make([]int, n)
	copy(resultado, arranjo)
	i := aleatorio.Intn(n)
	resultado[i] = 1 - resultado[i]
	return resultado
}

func temperaSimulada(nItens int, aleatorio *rand.Rand, valores []int, tamanhos []int, tamanhoMaximo int, execucao int) ([]int, []IterationData) {
	solucao := make([]int, nItens)
	for i := 0; i < nItens; i++ {
		solucao[i] = 0
	}
	melhorSolucao := solucao

	valorAtual, _ := valorTotalTamanho(solucao, valores, tamanhos, tamanhoMaximo)
	valorMelhorSolucao := valorAtual

	resultados := []IterationData{}

	fmt.Printf("\n%-15s %-15s %-15s\n", "Iteração", "Valor Atual", "Temperatura")
	fmt.Println("-----------------------------------------")

	// Itera até a temperatura atingir zero
	for iteracao := 0; temperatura(iteracao) > 0.0001; iteracao++ {
		arranjoAdjacente := adjacente(solucao, aleatorio)
		valorAdjacente, _ := valorTotalTamanho(arranjoAdjacente, valores, tamanhos, tamanhoMaximo)

		delta_e := valorAdjacente - valorAtual

		// Se a solução adjacente for melhor, será aceita
		if delta_e > 0 {
			solucao = arranjoAdjacente
			valorAtual = valorAdjacente

			if valorAtual > valorMelhorSolucao {
				melhorSolucao = solucao
				valorMelhorSolucao = valorAtual
			}
		} else {
			probAceitar := math.Exp(float64(delta_e) / temperatura(iteracao))

			if rand.Float64() < probAceitar {
				solucao = arranjoAdjacente
				valorAtual = valorAdjacente

				if valorAtual > valorMelhorSolucao {
					melhorSolucao = solucao
					valorMelhorSolucao = valorAtual
				}
			}
		}

		if iteracao%10 == 0 {
			fmt.Printf("%-15d %-15d %-15.2f\n", iteracao, valorAtual, temperatura(iteracao))
		}

		// Adiciona os dados da iteração no slice resultados
		resultados = append(resultados, IterationData{
			Iteracao:      iteracao,
			ValorAtual:    valorAtual,
			Temperatura:   temperatura(iteracao),
			MelhorSolucao: append([]int{}, melhorSolucao...),
			Execucao:      execucao,
		})
	}

	return melhorSolucao, resultados
}

func printTabela(valores, tamanhos []int, arranjo []int) {
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

func temperaturaLinear(x int) float64 {
	return 1000.0 - float64(x)
}

func temperatura1(x int) float64 {
	return 1000.0 * math.Pow(10, -3.0*float64(x)/1000.0)
}

func temperatura2(x int) float64 {
	return -0.001*math.Pow(float64(x), 2) + 750
}

func temperatura(x int) float64 {
	return temperatura1(x)
}

func main() {
	var valores, tamanhos []int
	var tamanhoMaximo int

	// Gera nova configuração
	valores, tamanhos, tamanhoMaximo = generateKnapsackConfig(50)
	aleatorio := rand.New(rand.NewSource(time.Now().UnixNano()))

	resultadosTotais := []IterationData{}

	fmt.Printf("Tamanho máximo da mochila = %d\n", tamanhoMaximo)
	fmt.Printf("Temperatura inicial = %.1f\n", temperatura(0))

	fmt.Println("\nInício da demonstração de têmpera simulada com mochila")

	for execucao := 1; execucao <= 2; execucao++ {
		fmt.Println("\nValores e tamanhos dos itens:")
		printTabela(valores, tamanhos, make([]int, len(valores)))

		solucao, resultados := temperaSimulada(len(valores), aleatorio, valores, tamanhos, tamanhoMaximo, execucao)
		resultadosTotais = append(resultadosTotais, resultados...)

		fmt.Println("\n\nSolução encontrada:")
		printTabela(valores, tamanhos, solucao)

		valor, tamanho := valorTotalTamanho(solucao, valores, tamanhos, tamanhoMaximo)

		fmt.Printf("\nValor total da solucao = %d\n", valor)
		fmt.Printf("Tamanho total da solucao = %d\n", tamanho)
	}

	// Serializa os resultados em JSON
	jsonData, err := json.MarshalIndent(resultadosTotais, "", "  ")
	if err != nil {
		fmt.Println("Erro ao gerar JSON:", err)
		return
	}

	// Salva o JSON em um arquivo
	err = os.WriteFile("resultados.json", jsonData, 0644)
	if err != nil {
		fmt.Println("Erro ao salvar o arquivo JSON:", err)
		return
	}
	fmt.Println(string(jsonData))
}
