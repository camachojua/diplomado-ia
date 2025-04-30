import pandas as pd
import numpy as np
import matplotlib.pyplot as plt
import torch
import torch.nn as nn
import torch.nn.functional as F
import torch.optim as optim
from torch.utils.data import DataLoader, TensorDataset
from sklearn.preprocessing import StandardScaler
from sklearn.model_selection import train_test_split

# Verificar si CUDA está disponible
device = torch.device("cuda" if torch.cuda.is_available() else "cpu")
print(f"Usando dispositivo: {device}")

# Cargar los datos
df = pd.read_csv('/content/geria2023gro.csv', encoding='latin1')

# Verificar las dimensiones del dataset
print(f"Dimensiones del dataset: {df.shape}")

# Seleccionar columnas relevantes para el análisis de ECNT
# Basándonos en el análisis previo, seleccionamos columnas relacionadas con:
# - Diabetes (ADM, PDM, HBA)
# - Hipertensión (AHA)
# - Obesidad (AOB)

# Primero, identificamos todas las columnas con estos prefijos
diabetes_cols = [col for col in df.columns if col.startswith(('ADM', 'PDM', 'HBA'))]
hipertension_cols = [col for col in df.columns if col.startswith('AHA')]
obesidad_cols = [col for col in df.columns if col.startswith('AOB')]

# Combinar todas las columnas de interés
ecnt_cols = diabetes_cols + hipertension_cols + obesidad_cols

print(f"Total de columnas seleccionadas: {len(ecnt_cols)}")
print(f"Columnas de diabetes: {len(diabetes_cols)}")
print(f"Columnas de hipertensión: {len(hipertension_cols)}")
print(f"Columnas de obesidad: {len(obesidad_cols)}")

# Extraer solo las columnas relevantes
X = df[ecnt_cols].copy()

# Verificar valores faltantes
missing_percentage = X.isnull().mean() * 100
print("\nPorcentaje de valores faltantes por columna:")
print(missing_percentage.describe())

# Imputar valores faltantes con la media (esto es una simplificación)
X = X.fillna(X.mean())

# Escalar los datos
scaler = StandardScaler()
X_scaled = scaler.fit_transform(X)

# Dividir en conjunto de entrenamiento y prueba
X_train, X_test = train_test_split(X_scaled, test_size=0.2, random_state=42)

print(f"\nDimensiones de X_train: {X_train.shape}")
print(f"Dimensiones de X_test: {X_test.shape}")

# Convertir a tensores de PyTorch
X_train_tensor = torch.FloatTensor(X_train)
X_test_tensor = torch.FloatTensor(X_test)

# Crear datasets y dataloaders
train_dataset = TensorDataset(X_train_tensor, X_train_tensor)
test_dataset = TensorDataset(X_test_tensor, X_test_tensor)

batch_size = 128
train_loader = DataLoader(train_dataset, batch_size=batch_size, shuffle=True)
test_loader = DataLoader(test_dataset, batch_size=batch_size)

# Definir hiperparámetros del VAE
original_dim = X_train.shape[1]  # Número de características
latent_dim = 2                   # Dimensión del espacio latente (2 para visualización)
intermediate_dim = 64            # Dimensión de la capa intermedia

# Definir la clase VAE
class VAE(nn.Module):
    def __init__(self, input_dim, hidden_dim, latent_dim):
        super(VAE, self).__init__()
        
        # Encoder
        self.encoder = nn.Sequential(
            nn.Linear(input_dim, hidden_dim),
            nn.ReLU()
        )
        
        # Capas para la media y la varianza del espacio latente
        self.fc_mu = nn.Linear(hidden_dim, latent_dim)
        self.fc_var = nn.Linear(hidden_dim, latent_dim)
        
        # Decoder
        self.decoder = nn.Sequential(
            nn.Linear(latent_dim, hidden_dim),
            nn.ReLU(),
            nn.Linear(hidden_dim, input_dim)
        )
    
    def encode(self, x):
        h = self.encoder(x)
        mu = self.fc_mu(h)
        log_var = self.fc_var(h)
        return mu, log_var
    
    def reparameterize(self, mu, log_var):
        std = torch.exp(0.5 * log_var)
        eps = torch.randn_like(std)
        z = mu + eps * std
        return z
    
    def decode(self, z):
        return self.decoder(z)
    
    def forward(self, x):
        mu, log_var = self.encode(x)
        z = self.reparameterize(mu, log_var)
        x_recon = self.decode(z)
        return x_recon, mu, log_var

# Función de pérdida
def loss_function(recon_x, x, mu, log_var):
    # Pérdida de reconstrucción
    recon_loss = F.mse_loss(recon_x, x, reduction='sum')
    
    # Divergencia KL
    kl_loss = -0.5 * torch.sum(1 + log_var - mu.pow(2) - log_var.exp())
    
    return recon_loss + kl_loss

# Inicializar el modelo y el optimizador
model = VAE(original_dim, intermediate_dim, latent_dim).to(device)
optimizer = optim.Adam(model.parameters(), lr=1e-3)

# Función de entrenamiento
def train(model, train_loader, optimizer, epoch):
    model.train()
    train_loss = 0
    
    for batch_idx, (data, _) in enumerate(train_loader):
        data = data.to(device)
        optimizer.zero_grad()
        
        recon_batch, mu, log_var = model(data)
        loss = loss_function(recon_batch, data, mu, log_var)
        
        loss.backward()
        train_loss += loss.item()
        optimizer.step()
    
    return train_loss / len(train_loader.dataset)

# Función de evaluación
def test(model, test_loader):
    model.eval()
    test_loss = 0
    
    with torch.no_grad():
        for batch_idx, (data, _) in enumerate(test_loader):
            data = data.to(device)
            
            recon_batch, mu, log_var = model(data)
            loss = loss_function(recon_batch, data, mu, log_var)
            
            test_loss += loss.item()
    
    return test_loss / len(test_loader.dataset)

# Entrenamiento del modelo
epochs = 50
train_losses = []
test_losses = []

print("\nEntrenando el modelo VAE...")
for epoch in range(1, epochs + 1):
    train_loss = train(model, train_loader, optimizer, epoch)
    test_loss = test(model, test_loader)
    
    train_losses.append(train_loss)
    test_losses.append(test_loss)
    
    if epoch % 10 == 0:
        print(f'Época {epoch}, Pérdida de entrenamiento: {train_loss:.4f}, Pérdida de prueba: {test_loss:.4f}')

# Visualizar la pérdida durante el entrenamiento
plt.figure(figsize=(10, 6))
plt.plot(train_losses, label='Pérdida de entrenamiento')
plt.plot(test_losses, label='Pérdida de validación')
plt.title('Pérdida del VAE durante el entrenamiento')
plt.xlabel('Épocas')
plt.ylabel('Pérdida')
plt.legend()
plt.grid(True)
plt.show()

# Codificar los datos en el espacio latente
model.eval()
with torch.no_grad():
    # Codificar los datos de entrenamiento
    z_mean_train, _ = model.encode(X_train_tensor.to(device))
    z_mean_train = z_mean_train.cpu().numpy()
    
    # Codificar los datos de prueba
    z_mean_test, _ = model.encode(X_test_tensor.to(device))
    z_mean_test = z_mean_test.cpu().numpy()

# Visualizar el espacio latente
plt.figure(figsize=(12, 10))
plt.scatter(z_mean_train[:, 0], z_mean_train[:, 1], c='blue', alpha=0.5, label='Entrenamiento')
plt.scatter(z_mean_test[:, 0], z_mean_test[:, 1], c='red', alpha=0.5, label='Prueba')
plt.title('Visualización del espacio latente 2D')
plt.xlabel('z[0]')
plt.ylabel('z[1]')
plt.legend()
plt.grid(True)
plt.show()

# Generar nuevos datos a partir del espacio latente
n = 15  # Número de puntos en cada dimensión
grid_x = np.linspace(-3, 3, n)
grid_y = np.linspace(-3, 3, n)

# Visualizar algunos ejemplos generados
plt.figure(figsize=(12, 10))
for i, yi in enumerate(grid_y):
    for j, xi in enumerate(grid_x):
        z_sample = torch.tensor([[xi, yi]], dtype=torch.float).to(device)
        with torch.no_grad():
            x_decoded = model.decode(z_sample).cpu().numpy()
        
        # Si i=0 y j=0, mostrar la forma de los datos generados
        if i == 0 and j == 0:
            print(f"\nForma de los datos generados: {x_decoded.shape}")
        
        # Solo visualizamos algunos ejemplos (4x4)
        if i < 4 and j < 4:
            plt.subplot(4, 4, i * 4 + j + 1)
            plt.imshow(x_decoded.reshape(1, -1), aspect='auto', cmap='viridis')
            plt.title(f'z=({xi:.1f}, {yi:.1f})')
            plt.axis('off')

plt.tight_layout()
plt.suptitle('Ejemplos generados desde el espacio latente', y=1.02)
plt.show()

# Analizar la estructura del espacio latente
# Vamos a ver si hay agrupaciones naturales en el espacio latente
from sklearn.cluster import KMeans

# Aplicar K-means al espacio latente
kmeans = KMeans(n_clusters=3, random_state=42)
clusters = kmeans.fit_predict(z_mean_train)

# Visualizar los clusters
plt.figure(figsize=(12, 10))
plt.scatter(z_mean_train[:, 0], z_mean_train[:, 1], c=clusters, cmap='viridis', alpha=0.7)
plt.scatter(kmeans.cluster_centers_[:, 0], kmeans.cluster_centers_[:, 1], c='red', marker='x', s=200)
plt.title('Clusters en el espacio latente')
plt.xlabel('z[0]')
plt.ylabel('z[1]')
plt.colorbar(label='Cluster')
plt.grid(True)
plt.show()

# Analizar las características de cada cluster
# Primero, asignar los clusters a los datos originales
df_clusters = pd.DataFrame(X_train, columns=ecnt_cols)
df_clusters['cluster'] = clusters

# Calcular la media de cada característica por cluster
cluster_means = df_clusters.groupby('cluster').mean()

# Visualizar las características más distintivas por cluster
plt.figure(figsize=(15, 10))
cluster_means.T.plot(kind='bar', ax=plt.gca())
plt.title('Valores medios de características por cluster')
plt.xlabel('Características')
plt.ylabel('Valor medio escalado')
plt.legend(title='Cluster')
plt.xticks(rotation=90)
plt.tight_layout()
plt.show()

# Seleccionar las 10 características más distintivas
# Calculamos la varianza entre clusters para cada característica
feature_variance = cluster_means.var(axis=0).sort_values(ascending=False)
top_features = feature_variance.head(10).index

# Visualizar solo las características más distintivas
plt.figure(figsize=(15, 8))
cluster_means[top_features].T.plot(kind='bar', ax=plt.gca())
plt.title('Top 10 características más distintivas por cluster')
plt.xlabel('Características')
plt.ylabel('Valor medio escalado')
plt.legend(title='Cluster')
plt.xticks(rotation=45)
plt.tight_layout()
plt.show()

print("\nCaracterísticas más distintivas entre clusters:")
for feature in top_features:
    print(f"- {feature}")
