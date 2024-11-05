import json
import matplotlib.pyplot as plt
from collections import Counter

graficos = [
    "tam10-ger-100-prob0.1,selecionados5.json",
    "tam20-ger-1000-prob0.1,selecionados10.json",
    "tam10-ger-1000-prob0.1,selecionados5.json",
    "tam20-ger-1000-prob0.01,selecionados5.json",
    "resultados.json"
]

for grafico in graficos:
    # Carrega os dados do JSON
    with open(grafico) as f:
        resultados = json.load(f)

    # Função para calcular a distribuição dos valores em intervalos definidos
    def calcular_distribuicao(valores, intervalo):
        max_valor = max(valores)
        # Define os intervalos de valores
        bins = range(0, max_valor + intervalo, intervalo)
        # Conta quantos valores estão em cada intervalo
        contagem = Counter((v // intervalo) * intervalo for v in valores)
        # Calcula a porcentagem para cada intervalo
        total = len(valores)
        porcentagens = {b: (contagem[b] / total) * 100 for b in bins}
        return bins, porcentagens

    # Extrai os valores dos melhores e últimos de cada entrada
    melhores_valores = [res["ValorMelhor"] for res in resultados]

    # Define o intervalo de cada faixa (exemplo: 200 unidades)
    intervalo = 200

    # Calcula as distribuições para o gráfico
    bins, melhores_porcentagens = calcular_distribuicao(melhores_valores, intervalo)

    # Cria o gráfico
    plt.figure(figsize=(10, 6))
    plt.bar([str(b) + '-' + str(b + intervalo) for b in bins],
            [melhores_porcentagens[b] for b in bins],
            color='blue', alpha=0.7)
    plt.xlabel('Intervalo de Valores')
    plt.ylabel('Porcentagem')
    plt.title('Distribuição dos Melhores Valores')
    plt.xticks(rotation=45)

    # Salva a figura
    plt.tight_layout()
    plt.savefig(f"{grafico.replace('.json', '')}_distribuicao.png")
    plt.close()
