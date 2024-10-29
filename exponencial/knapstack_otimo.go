package main

import (
	"fmt"
	"math/rand"
	"time"
)

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

// knapsackRecursive é uma função recursiva que resolve o problema da mochila
func knapsackRecursive(tamanhos []int, valores []int, tamanhoMaximo int, n int) int {
	// Caso base: se não há itens ou a capacidade da mochila é 0
	if n == 0 || tamanhoMaximo == 0 {
		return 0
	}

	// Se o tamanho do item atual é maior que a capacidade da mochila,
	// não podemos incluí-lo
	if tamanhos[n-1] > tamanhoMaximo {
		return knapsackRecursive(tamanhos, valores, tamanhoMaximo, n-1)
	} else {
		// Retorna o máximo entre incluir ou não o item atual
		include := valores[n-1] + knapsackRecursive(tamanhos, valores, tamanhoMaximo-tamanhos[n-1], n-1)
		exclude := knapsackRecursive(tamanhos, valores, tamanhoMaximo, n-1)
		return max(include, exclude)
	}
}

// max é uma função auxiliar para retornar o maior valor
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	var valores, tamanhos []int
	var tamanhoMaximo int

	// Tenta carregar configuração existente, se não, gera nova
	config, err := loadConfig()
	if err == nil {
		valores = config.Valores
		tamanhos = config.Tamanhos
		tamanhoMaximo = config.TamanhoMaximo
	} else {
		valores, tamanhos, tamanhoMaximo = generateKnapsackConfig(10)
		saveConfig(MochilaConfig{valores, tamanhos, tamanhoMaximo})
	}

	fmt.Printf("Valores: %v\n", valores)
	fmt.Printf("Tamanhos: %v\n", tamanhos)
	fmt.Printf("Tamanho máximo da mochila: %d\n", tamanhoMaximo)

	// Usando o algoritmo recursivo
	valorMaximo := knapsackRecursive(tamanhos, valores, tamanhoMaximo, len(valores))
	fmt.Printf("Valor máximo encontrado pelo algoritmo exponencial = %d\n", valorMaximo)
}
