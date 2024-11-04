import json
import matplotlib.pyplot as plt
from collections import Counter

graficos = [
    'resultados_exponencial.json',
    'resultados_linear.json',
    'resultados_quadratica.json',
    'resultados_sigmoidal.json'
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
    ultimos_valores = [res["ValorUltimo"] for res in resultados]

    # Define o intervalo de cada faixa (exemplo: 200 unidades)
    intervalo = 200

    # Calcula as distribuições para os gráficos
    bins, melhores_porcentagens = calcular_distribuicao(melhores_valores, intervalo)
    _, ultimos_porcentagens = calcular_distribuicao(ultimos_valores, intervalo)

    # Cria uma figura com subplots para os três gráficos
    fig, axs = plt.subplots(3, 1, figsize=(10, 15))

    # Gráfico dos Melhores Valores
    axs[0].bar([str(b) + '-' + str(b + intervalo) for b in bins],
               [melhores_porcentagens[b] for b in bins],
               color='blue', alpha=0.7)
    axs[0].set_xlabel('Intervalo de Valores')
    axs[0].set_ylabel('Porcentagem')
    axs[0].set_title('Distribuição dos Melhores Valores')
    axs[0].tick_params(axis='x', rotation=45)

    # Gráfico dos Últimos Valores
    axs[1].bar([str(b) + '-' + str(b + intervalo) for b in bins],
               [ultimos_porcentagens[b] for b in bins],
               color='red', alpha=0.7)
    axs[1].set_xlabel('Intervalo de Valores')
    axs[1].set_ylabel('Porcentagem')
    axs[1].set_title('Distribuição dos Últimos Valores')
    axs[1].tick_params(axis='x', rotation=45)

    # Gráfico Comparativo
    axs[2].bar([str(b) + '-' + str(b + intervalo) for b in bins],
               [melhores_porcentagens[b] for b in bins],
               color='blue', alpha=0.5, label='Melhores Valores')
    axs[2].bar([str(b) + '-' + str(b + intervalo) for b in bins],
               [ultimos_porcentagens[b] for b in bins],
               color='red', alpha=0.5, label='Últimos Valores')
    axs[2].set_xlabel('Intervalo de Valores')
    axs[2].set_ylabel('Porcentagem')
    axs[2].set_title('Distribuição Comparativa dos Melhores e Últimos Valores')
    axs[2].tick_params(axis='x', rotation=45)
    axs[2].legend()

    # Ajusta o layout e salva a figura
    plt.tight_layout()
    plt.savefig(f"{grafico.replace('.json', '')}_distribuicao.png")
    plt.close()
