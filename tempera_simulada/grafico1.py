import json
import matplotlib.pyplot as plt
import numpy as np

# Carregar dados do arquivo JSON
with open(
        "tempera_simulada\historico.json",
        "r") as file:
    dados = json.load(file)

# Extrair os valores atuais das iterações
valores = [item["valor_atual"] for item in dados]

# Definir os intervalos específicos
bins = [0, 100, 200, 300, 400, 450, 480, 500, 515, 530, 545]
labels = [
    "0-100", "100-200", "200-300", "300-400", "400-450", "450-480", "480-500",
    "500-515", "515-530", "530-545"
]

# Calcular a frequência relativa para cada intervalo
contagens, _ = np.histogram(valores, bins=bins)
percentual = (contagens / len(valores)) * 100

# Filtrar para remover os intervalos sem ocorrências
labels_filtrados = [label for label, pct in zip(labels, percentual) if pct > 0]
percentual_filtrado = [pct for pct in percentual if pct > 0]

# Plotar o gráfico com os intervalos filtrados
plt.figure(figsize=(12, 6))
plt.bar(labels_filtrados,
        percentual_filtrado,
        color="skyblue",
        edgecolor="black")
plt.xlabel("Intervalos de Valores Atingidos")
plt.ylabel("Percentual (%)")
plt.title("Percentual das Iterações por Intervalo de Valor")
plt.xticks(rotation=45)
plt.tight_layout()
plt.show()
