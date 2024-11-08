package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strings"
	"time"
)

// gere constantes para o tamanho da população e o número de gerações, probabilidades de cruzamento e mutação e numero de selecionados
const (
	tamanhoPopulacao = 20
	nGeracoes        = 1000
	probMutacao      = 0.01
	selecionados     = 5
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

func gerarPopulacaoInicial(tamanhoPopulacao, nItens int, aleatorio *rand.Rand) [][]int {
	populacao := make([][]int, tamanhoPopulacao)

	for i := 0; i < tamanhoPopulacao; i++ {
		populacao[i] = make([]int, nItens)
		for j := 0; j < nItens; j++ {
			populacao[i][j] = aleatorio.Intn(2)
		}
	}
	return populacao
}

func fitness(arranjo []int, valores []int, tamanhos []int, tamanhoMaximo int) int {
	valor, _ := valorTotalTamanho(arranjo, valores, tamanhos, tamanhoMaximo)
	return valor
}

// crossover realiza o cruzamento entre dois pais para gerar dois filhos.
func crossover(pai1, pai2 []int, aleatorio *rand.Rand) ([]int, []int) {
	n := len(pai1)
	pontoCorte := aleatorio.Intn(n-1) + 1

	// Realizar o crossover em um ponto de corte
	//  [XXXYYYYY]
	//  [YYYXXXXX]

	filho1 := append([]int(nil), pai1[:pontoCorte]...)
	filho1 = append(filho1, pai2[pontoCorte:]...)

	filho2 := append([]int(nil), pai2[:pontoCorte]...)
	filho2 = append(filho2, pai1[pontoCorte:]...)

	return filho1, filho2
}

// mutacao realiza mutação em uma solução com base em uma probabilidade.
func mutacao(arranjo []int, aleatorio *rand.Rand, probMutacao float64) []int {
	for i := range arranjo {
		if aleatorio.Float64() < probMutacao {
			arranjo[i] = 1 - arranjo[i]
		}
	}
	return arranjo
}

func selecionarMelhores(populacao [][]int, valores, tamanhos []int, tamanhoMaximo, nMelhores int) [][]int {
	type individuo struct {
		arranjo []int
		valor   int
	}

	// calcula o valor da funcao fitness para cada individuo
	individuos := make([]individuo, len(populacao))
	for i, arranjo := range populacao {
		val := fitness(arranjo, valores, tamanhos, tamanhoMaximo)
		individuos[i] = individuo{arranjo, val}
	}

	// ordenar descendentemente
	sort.Slice(individuos, func(i, j int) bool {
		return individuos[i].valor > individuos[j].valor
	})

	// seleciona os nMelhores indivíduos
	melhores := make([][]int, nMelhores)
	for i := 0; i < nMelhores; i++ {
		melhores[i] = individuos[i].arranjo
	}

	return melhores
}

// algGenetico executa o algoritmo genético para o problema da mochila.
func algGenetico(nItens int, aleatorio *rand.Rand, valores []int, tamanhos []int, tamanhoMaximo, tamanhoPopulacao, maxGeracoes int) []int {
	probabilidadeMutacao := probMutacao
	var melhorAtual []int

	populacao := gerarPopulacaoInicial(tamanhoPopulacao, nItens, aleatorio)

	for geracao := 0; geracao < maxGeracoes; geracao++ {
		nMelhores := selecionados
		melhores := selecionarMelhores(populacao, valores, tamanhos, tamanhoMaximo, nMelhores)
		melhorAtual = melhores[0]

		var novaPopulacao [][]int

		// ETAPA de cruzamento e mutação
		for i := 0; i < len(melhores)-1; i += 2 {

			// Pega um par de melhores indivíduos e realiza o crossover
			filho1, filho2 := crossover(melhores[i], melhores[i+1], aleatorio)

			// Aplica a mutação nos filhos
			novaPopulacao = append(novaPopulacao, mutacao(filho1, aleatorio, probabilidadeMutacao))
			novaPopulacao = append(novaPopulacao, mutacao(filho2, aleatorio, probabilidadeMutacao))
		}

		// preenche o restante da população com indivíduos mutados
		for len(novaPopulacao) < tamanhoPopulacao {
			individuoMutado := mutacao(append([]int(nil), melhorAtual...), aleatorio, probabilidadeMutacao) // Aplica mutação
			novaPopulacao = append(novaPopulacao, individuoMutado)
		}

		populacao = novaPopulacao
	}

	melhores := selecionarMelhores(populacao, valores, tamanhos, tamanhoMaximo, 1)

	if fitness(melhores[0], valores, tamanhos, tamanhoMaximo) > fitness(melhorAtual, valores, tamanhos, tamanhoMaximo) {
		return melhores[0]
	}
	return melhorAtual
}

type Resultado struct {
	ValorMelhor   int
	TamanhoMelhor int
	MelhorSolucao []int
}

func main() {
	maxGeracoes := nGeracoes

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

	aleatorio := rand.New(rand.NewSource(time.Now().UnixNano()))

	fmt.Printf("Tamanho máximo da mochila = %d\n", tamanhoMaximo)

	fmt.Println("\nInício da demonstração de algoritmo genético com mochila")

	resultados := []Resultado{}

	// Dentro do loop, adicione os resultados em vez de sobrescrever elementos:
	for i := 0; i < 1000; i++ {
		solucao := algGenetico(len(valores), aleatorio, valores, tamanhos, tamanhoMaximo, tamanhoPopulacao, maxGeracoes)
		valor, tamanho := valorTotalTamanho(solucao, valores, tamanhos, tamanhoMaximo)

		resultado := Resultado{
			ValorMelhor:   valor,
			TamanhoMelhor: tamanho,
			MelhorSolucao: solucao,
		}
		// Adicione o resultado ao slice
		resultados = append(resultados, resultado)
	}

	// Serializa os resultados para JSON
	file, err := os.Create("resultados.json")
	if err != nil {
		fmt.Println("Erro ao criar o arquivo:", err)
		return
	}
	defer file.Close()

	// serialize to json
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	err = encoder.Encode(resultados)
	if err != nil {
		fmt.Println("Erro ao serializar para JSON:", err)
		return
	}
}
