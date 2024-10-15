package main

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"time"
)

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
	probabilidadeMutacao := 0.05

	populacao := gerarPopulacaoInicial(tamanhoPopulacao, nItens, aleatorio)

	for geracao := 0; geracao < maxGeracoes; geracao++ {
		nMelhores := tamanhoPopulacao / 2
		melhores := selecionarMelhores(populacao, valores, tamanhos, tamanhoMaximo, nMelhores)

		var novaPopulacao [][]int

		// ETAPA de cruzamento e mutação
		for i := 0; i < len(melhores)-1; i += 2 {

			// Pega um par de melhores indivíduos e realiza o crossover
			filho1, filho2 := crossover(melhores[i], melhores[i+1], aleatorio)

			// Aplica a mutação nos filhos
			novaPopulacao = append(novaPopulacao, mutacao(filho1, aleatorio, probabilidadeMutacao))
			novaPopulacao = append(novaPopulacao, mutacao(filho2, aleatorio, probabilidadeMutacao))
		}

		// Se sobrar um melhor indivíduo, ele é adicionado sem crossover
		if len(melhores)%2 != 0 {
			novaPopulacao = append(novaPopulacao, melhores[len(melhores)-1])
		}

		// Se a nova população não atingir o tamanho original, aplicar mutações nos melhores
		for len(novaPopulacao) < tamanhoPopulacao {
			// Seleciona um dos melhores para aplicar mutação
			melhorIndividuo := melhores[aleatorio.Intn(len(melhores))]
			individuoMutado := mutacao(append([]int(nil), melhorIndividuo...), aleatorio, probabilidadeMutacao) // Aplica mutação
			novaPopulacao = append(novaPopulacao, individuoMutado)
		}

		populacao = novaPopulacao

		fmt.Printf("Geração %d concluida.\n", geracao+1)
	}

	melhores := selecionarMelhores(populacao, valores, tamanhos, tamanhoMaximo, 1)
	return melhores[0]
}

func main() {
	valores := []int{95, 75, 60, 85, 40, 120, 30, 65, 50, 90}
	tamanhos := []int{50, 40, 30, 55, 25, 60, 35, 45, 40, 50}
	tamanhoMaximo := 300
	tamanhoPopulacao := 15
	maxGeracoes := 100

	aleatorio := rand.New(rand.NewSource(time.Now().UnixNano()))

	fmt.Printf("Tamanho máximo da mochila = %d\n", tamanhoMaximo)

	fmt.Println("\nInício da demonstração de algoritmo genético com mochila")

	fmt.Println("\nValores e tamanhos dos itens:")
	printTabela(valores, tamanhos, make([]int, len(valores)))

	solucao := algGenetico(len(valores), aleatorio, valores, tamanhos, tamanhoMaximo, tamanhoPopulacao, maxGeracoes)

	fmt.Println("\n\nSolução encontrada:")
	printTabela(valores, tamanhos, solucao)

	valor, tamanho := valorTotalTamanho(solucao, valores, tamanhos, tamanhoMaximo)

	fmt.Printf("\nValor total da solução = %d\n", valor)
	fmt.Printf("Tamanho total da solução = %d\n", tamanho)
}
