# CloudSweep

**Plateforme SaaS FinOps/GreenOps pour automatiser le nettoyage des ressources cloud inutilisees**

## Vision

CloudSweep aide les entreprises a:
- **Reduire les couts cloud** en identifiant et supprimant automatiquement les ressources inutilisees
- **Diminuer l'empreinte carbone** en eliminant le gaspillage de ressources
- **Automatiser la gouvernance cloud** avec des politiques de nettoyage configurables

### Clouds supportes
- Amazon Web Services (AWS)
- Microsoft Azure
- Google Cloud Platform (GCP)

### Ressources detectees
- Instances EC2/VM arretees
- Volumes EBS/Disques non attaches
- Snapshots obsoletes
- Adresses IP elastiques non utilisees
- Load balancers sans cibles
- Buckets S3 vides ou abandonnes

## Architecture

```
cloudsweep/
├── cmd/
│   ├── api/            # Point d'entree API REST
│   └── worker/         # Point d'entree Worker asynchrone
├── internal/
│   ├── domain/         # Entites et interfaces (ports)
│   ├── application/    # Cas d'usage (services)
│   ├── infrastructure/ # Adaptateurs (repositories, clients cloud)
│   └── interfaces/     # Controleurs HTTP, handlers
├── pkg/                # Code partage reutilisable
└── config/             # Fichiers de configuration
```

## Stack Technique

| Composant | Technologie |
|-----------|-------------|
| Langage | Go 1.21+ |
| API REST | Gin |
| Queue | Asynq (Redis) |
| ORM | GORM |
| Base de donnees | PostgreSQL |
| Configuration | Viper |
| Conteneurisation | Docker |

## Demarrage rapide

### Prerequis
- Go 1.21+
- Docker & Docker Compose
- Make

### Installation

```bash
# Cloner le repository
git clone https://github.com/cloudsweep/cloudsweep.git
cd cloudsweep

# Lancer l'infrastructure
make docker-up

# Installer les dependances
make deps

# Lancer l'API en mode developpement
make run-api

# Lancer le worker en mode developpement
make run-worker
```

### Commandes Make

```bash
make build       # Compile les binaires
make test        # Execute les tests
make lint        # Analyse statique du code
make docker-up   # Demarre les conteneurs
make docker-down # Arrete les conteneurs
make migrate     # Execute les migrations
```

## Configuration

Les variables d'environnement peuvent etre definies dans un fichier `.env`:

```env
# Server
SERVER_PORT=8080
SERVER_ENV=development

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=cloudsweep
DB_PASSWORD=secret
DB_NAME=cloudsweep

# Redis
REDIS_ADDR=localhost:6379

# Cloud Providers
AWS_REGION=eu-west-1
```

## API Endpoints

| Methode | Endpoint | Description |
|---------|----------|-------------|
| GET | /health | Health check |
| GET | /api/v1/resources | Liste des ressources |
| POST | /api/v1/scans | Lancer un scan |
| GET | /api/v1/scans/:id | Statut d'un scan |
| POST | /api/v1/cleanup | Executer un nettoyage |
| GET | /api/v1/policies | Liste des politiques |
| POST | /api/v1/policies | Creer une politique |

## Licence

MIT License - voir [LICENSE](LICENSE)
