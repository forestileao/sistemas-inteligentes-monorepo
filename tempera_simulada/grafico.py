import json
import pandas as pd
import numpy as np
import plotly.express as px

# Carregar dados do arquivo JSON
with open("C:/Users/leand/OneDrive/Desktop/SI/sistemas-inteligentes-monorepo/tempera_simulada/resultados.json", "r") as file:
    dados = json.load(file)

# Extrair os valores atuais das iterações
valores = [item["valor_atual"] for item in dados]

# Criar um DataFrame com os valores
df = pd.DataFrame(valores, columns=["valor_atual"])

# Definir os intervalos específicos
bins = [0, 100, 200, 300, 400, 450, 480, 500, 515, 530, 545]
labels = ["0-100", "100-200", "200-300", "300-400", "400-450", "450-480", "480-500", "500-515", "515-530", "530-545"]

# Criar uma nova coluna de categorias no DataFrame
df['intervalo'] = pd.cut(df['valor_atual'], bins=bins, labels=labels, right=False)

# Calcular a frequência relativa para cada intervalo
contagens = df['intervalo'].value_counts(normalize=True) * 100
contagens = contagens.sort_index()  # Ordenar os índices para manter a ordem dos intervalos

# Filtrar para remover os intervalos sem ocorrências
contagens_filtradas = contagens[contagens > 0]

# Criar o gráfico com plotly
fig = px.bar(
    x=contagens_filtradas.index,
    y=contagens_filtradas.values,
    labels={'x': 'Intervalos de Valores Atingidos', 'y': 'Percentual (%)'},
    title='Percentual das Iterações por Intervalo de Valor',
    color=contagens_filtradas.values,
    color_continuous_scale='Viridis'  # Usando a paleta Viridis
)

# Ajustar a aparência do gráfico
fig.update_traces(marker=dict(line=dict(width=1, color='DarkSlateGrey')))  # Adicionar contorno aos barras
fig.update_layout(xaxis_title='Intervalos de Valores Atingidos', yaxis_title='Percentual (%)')

# Exibir o gráfico
fig.show()
