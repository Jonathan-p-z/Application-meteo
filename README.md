# Weather App

Application météo simple avec un backend en Go et un frontend HTML/CSS/JS.

## Structure du projet

- `backend/`
  - `main.go` : point d'entrée du serveur HTTP (API + static).
  - `config/` : chargement de `.env` (`config.go`, `.env`).
  - `handlers/` : handlers HTTP (par ex. `WeatherHandler`).
  - `services/` : logique métier (appel à WeatherAPI).
  - `models/` : structures de données (JSON).
  - `utils/` : helpers (client HTTP, etc.).
- `frontend/`
  - `index.html` : page d'accueil.
  - `weather.html` : page de consultation météo par ville.
  - `assets/css/` : styles séparés par rôle (`base`, `layout`, `components`, `pages`).
  - `assets/js/app.js` : logique front (recherche de ville, affichage des résultats).

## Conventions & hygiène

1. **Nommage clair et cohérent**
   - Langue du code : **anglais** (Go, JS, CSS).
   - Texte utilisateur : **français** (HTML, messages d'erreur front).
   - Noms descriptifs : ex. `GetWeatherForCity`, `WeatherHandler`, `.weather-card`.

2. **Structure lisible**
   - Backend / frontend séparés.
   - Dossiers pour `config`, `services`, `handlers`, `models`.
   - CSS découpé (`base`, `layout`, `components`, `pages`).

3. **Formatter & linter**
   - Go : `gofmt` + `go vet` sur `backend/`.
   - Recommandé : `golangci-lint` (non obligatoire).
   - Front : utiliser un formatter dans l’éditeur (Prettier) sur HTML/CSS/JS.

4. **Workflow Git**
   - Petites branches par fonctionnalité, par ex. :
     - `feature/weather-details`
     - `feature/error-handling`
   - Commits structurés :
     - `feat: add detailed weather info`
     - `fix: handle 400 error from WeatherAPI`
   - PR concises : une PR = une fonctionnalité / un fix.

## Principes de code

1. **KISS (Keep It Simple, Stupid)**
   - Préférer des fonctions courtes et simples.
   - Exemple Go :
     - `GetWeatherForCity` délègue à des helpers (`handleUpstreamStatus`, `buildWeatherModel`).

2. **DRY (Don’t Repeat Yourself)**
   - Facteur commun au même endroit :
     - `config.GetWeatherAPIKey()` : une seule logique pour la clé API.
     - CSS commun → `components.css`, `layout.css`, etc.
     - Messages d’erreur backend centralisés dans `WeatherHandler`.

3. **YAGNI (You Ain’t Gonna Need It)**
   - Pas de features non utilisées :
     - page `alerts`, carte du monde Leaflet, etc. retirées tant qu’elles ne sont pas nécessaires.
   - Ajouter seulement ce qui est consommé par le frontend.

4. **Séparation des responsabilités**
   - **Backend** :
     - `handlers` : gèrent les requêtes/réponses HTTP.
     - `services` : appellent WeatherAPI, font la logique métier.
     - `models` : décrivent la forme des données.
     - `config` : charge la configuration/env.
   - **Frontend** :
     - `index.html` : présentation + navigation.
     - `weather.html` : interaction + affichage détaillé.
     - CSS → apparence uniquement, JS → comportement uniquement.

## Lancement

```bash
# Backend
cd backend
go mod tidy
go run ./...

# Frontend
# Ouvrir frontend/index.html dans le navigateur
# ou servir le dossier frontend avec un petit serveur HTTP
```
