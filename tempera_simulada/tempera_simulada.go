package main

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"
)

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

// resfriar calcula a nova temperatura com base na temperatura atual e no fator fatorResfriamento.
// Retorna zero se a nova temperatura for menor que 0.0001.
func resfriar(temperaturaAtual, fatorResfriamento float64) float64 {
	novaTemperatura := temperaturaAtual * fatorResfriamento

	if novaTemperatura < 0.0001 {
		novaTemperatura = 0
	}

	return novaTemperatura
}

func temperaSimulada(nItens int, aleatorio *rand.Rand, valores []int, tamanhos []int, tamanhoMaximo int, temperaturaInicial, fatorResfriamento float64) []int {
	temperaturaAtual := temperaturaInicial
	solucao := []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	for i := 0; i < nItens; i++ {
		solucao[i] = 0
	}
	valorAtual, _ := valorTotalTamanho(solucao, valores, tamanhos, tamanhoMaximo)

	fmt.Printf("\n%-15s %-15s %-15s\n", "Iteração", "Valor Atual", "Temperatura")
	fmt.Println("-----------------------------------------")

	// Itera até a temperatura atingir zero
	for iteracao := 0; temperaturaAtual > 0; iteracao++ {
		arranjoAdjacente := adjacente(solucao, aleatorio)
		valorAdjacente, _ := valorTotalTamanho(arranjoAdjacente, valores, tamanhos, tamanhoMaximo)

		delta_e := valorAdjacente - valorAtual

		// Se a solução adjacente for melhor, será aceita
		if delta_e > 0 {
			solucao = arranjoAdjacente
			valorAtual = valorAdjacente
		} else {
			probAceitar := math.Exp(float64(delta_e) / temperaturaAtual)

			if rand.Float64() < probAceitar {
				solucao = arranjoAdjacente
				valorAtual = valorAdjacente
			}
		}

		if iteracao%100 == 0 {
			fmt.Printf("%-15d %-15d %-15.2f\n", iteracao, valorAtual, temperaturaAtual)
		}

		temperaturaAtual = resfriar(temperaturaAtual, fatorResfriamento)
	}

	return solucao
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

func main() {
	valores := []int{95, 75, 60, 85, 40, 120, 30, 65, 50, 90}
	tamanhos := []int{50, 40, 30, 55, 25, 60, 35, 45, 40, 50}
	tamanhoMaximo := 250

	aleatorio := rand.New(rand.NewSource(time.Now().UnixNano()))
	var temperaturaInicial float64 = 12400
	fatorResfriamento := 0.98

	fmt.Printf("Tamanho máximo da mochila = %d\n", tamanhoMaximo)
	fmt.Printf("Temperatura inicial = %.1f\n", temperaturaInicial)
	fmt.Printf("Fator de Resfriamento = %.2f\n", fatorResfriamento)

	fmt.Println("\nInício da demonstração de têmpera simulada com mochila")

	fmt.Println("\nValores e tamanhos dos itens:")
	printTabela(valores, tamanhos, make([]int, len(valores)))

	solucao := temperaSimulada(10, aleatorio, valores, tamanhos, tamanhoMaximo, temperaturaInicial, fatorResfriamento)
	fmt.Println("\n\nArranjo encontrado:")
	printTabela(valores, tamanhos, solucao)

	valor, tamanho := valorTotalTamanho(solucao, valores, tamanhos, tamanhoMaximo)

	fmt.Printf("\nValor total do arranjo = %d\n", valor)
	fmt.Printf("Tamanho total do arranjo = %d\n", tamanho)
}
