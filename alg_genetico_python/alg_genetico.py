import random

# Função para calcular o valor total e o tamanho total de um arranjo
def valor_total_tamanho(arranjo, valores, tamanhos, tamanho_maximo):
    valor = 0
    tamanho = 0
    for i, v in enumerate(arranjo):
        if v == 1:
            valor += valores[i]
            tamanho += tamanhos[i]
    if tamanho > tamanho_maximo:
        valor = 0
    return valor, tamanho

# Função para imprimir a tabela de itens e a solução
def print_tabela(valores, tamanhos, arranjo):
    print("\n{:<10} {:<10} {:<10} {:<10}".format("Item", "Valor", "Tamanho", "Incluído"))
    print("-" * 40)
    for i in range(len(valores)):
        incluido = "Sim" if arranjo[i] == 1 else "Não"
        print("{:<10} {:<10} {:<10} {:<10}".format(i + 1, valores[i], tamanhos[i], incluido))

# Função para gerar a população inicial
def gerar_populacao_inicial(tamanho_populacao, n_itens):
    populacao = set()
    while len(populacao) < tamanho_populacao:
        arranjo = tuple(random.randint(0, 1) for _ in range(n_itens))
        populacao.add(arranjo)
    return list(populacao)

# Função de fitness para calcular o valor de um arranjo
def fitness(arranjo, valores, tamanhos, tamanho_maximo):
    valor, _ = valor_total_tamanho(arranjo, valores, tamanhos, tamanho_maximo)
    return valor

# Função de crossover uniforme para realizar cruzamento entre dois pais
def crossover_uniforme(pai1, pai2):
    return [(gene1 if random.random() < 0.5 else gene2) for gene1, gene2 in zip(pai1, pai2)]

# Função de mutação para realizar mutação em um arranjo
def mutacao(arranjo, prob_mutacao, tamanho_populacao, n_itens):
    novo_arranjo = [1 - gene if random.random() < prob_mutacao else gene for gene in arranjo]

    # Adiciona aleatoriamente um novo indivíduo à população com probabilidade
    if random.random() < 0.1:  # 10% de chance de gerar um novo indivíduo
        return [random.randint(0, 1) for _ in range(n_itens)]

    return novo_arranjo

# Função para selecionar os melhores indivíduos da população
def selecionar_melhores(populacao, valores, tamanhos, tamanho_maximo, n_melhores):
    populacao_avaliada = [(arranjo, fitness(arranjo, valores, tamanhos, tamanho_maximo)) for arranjo in populacao]
    populacao_avaliada.sort(key=lambda x: x[1], reverse=True)

    melhores = [individuo[0] for individuo in populacao_avaliada[:n_melhores]]

    # Adiciona alguns indivíduos aleatórios
    while len(melhores) < n_melhores + 2:  # Adicionando mais 2 indivíduos aleatórios
        aleatorio = random.choice(populacao)
        if aleatorio not in melhores:
            melhores.append(aleatorio)

    return melhores

# Função para calcular a diversidade da população
def calcular_diversidade(populacao):
    return len(set(tuple(individuo) for individuo in populacao)) / len(populacao)

# Função principal do algoritmo genético
def alg_genetico(n_itens, valores, tamanhos, tamanho_maximo, tamanho_populacao, max_geracoes):
    probabilidade_mutacao = 0.2  # Maior taxa no início
    populacao = gerar_populacao_inicial(tamanho_populacao, n_itens)

    for geracao in range(max_geracoes):
        if geracao > 20:  # Reduz a taxa de mutação após 20 gerações
            probabilidade_mutacao = 0.05

        n_melhores = tamanho_populacao // 2
        melhores = selecionar_melhores(populacao, valores, tamanhos, tamanho_maximo, n_melhores)

        nova_populacao = []

        # Etapa de cruzamento e mutação
        for i in range(0, len(melhores) - 1, 2):
            filho1, filho2 = crossover_uniforme(melhores[i], melhores[i + 1])
            nova_populacao.append(mutacao(filho1, probabilidade_mutacao, tamanho_populacao, n_itens))
            nova_populacao.append(mutacao(filho2, probabilidade_mutacao, tamanho_populacao, n_itens))

        if len(melhores) % 2 != 0:
            nova_populacao.append(melhores[-1])

        while len(nova_populacao) < tamanho_populacao:
            melhor_individuo = random.choice(melhores)
            nova_populacao.append(mutacao(melhor_individuo[:], probabilidade_mutacao, tamanho_populacao, n_itens))

        populacao = nova_populacao
        diversidade = calcular_diversidade(populacao)
        print(f"Geração {geracao + 1} concluída. Diversidade da população: {diversidade:.2f}")

    melhores = selecionar_melhores(populacao, valores, tamanhos, tamanho_maximo, 1)
    return melhores[0]

# Função principal para executar o programa
def main():
    valores = [95, 75, 60, 85, 40, 120, 30, 65, 50, 90]
    tamanhos = [50, 40, 30, 55, 25, 60, 35, 45, 40, 50]
    tamanho_maximo = 300
    tamanho_populacao = 15
    max_geracoes = 100

    print(f"Tamanho máximo da mochila = {tamanho_maximo}")
    print("\nInício da demonstração de algoritmo genético com mochila")
    print("\nValores e tamanhos dos itens:")
    print_tabela(valores, tamanhos, [0] * len(valores))

    solucao = alg_genetico(len(valores), valores, tamanhos, tamanho_maximo, tamanho_populacao, max_geracoes)

    print("\n\nSolução encontrada:")
    print_tabela(valores, tamanhos, solucao)

    valor, tamanho = valor_total_tamanho(solucao, valores, tamanhos, tamanho_maximo)
    print(f"\nValor total da solução = {valor}")
    print(f"Tamanho total da solução = {tamanho}")

if __name__ == "__main__":
    main()
 