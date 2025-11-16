# Weather API – `/api/weather`

API HTTP exposée par le backend Go pour récupérer la météo détaillée d'une ville ou d'un village.

## Endpoint

- Méthode : `GET`
- URL : `/api/weather`

## Paramètres

### Query

- `city` (string, obligatoire)  
  Nom de la ville ou du village.  
  Exemples : `Paris`, `Lyon`, `Marseille`, `New York`, `Tokyo`.

## Réponses

### 200 OK

Corps JSON (structure simplifiée) :

```json
{
  "city": "Paris",
  "country": "France",
  "region": "Ile-de-France",
  "lat": 48.86,
  "lon": 2.35,

  "temperature": 18.2,
  "feels_like": 17.5,
  "condition": "Ciel dégagé",
  "condition_icon_url": "https://cdn.weatherapi.com/weather/64x64/day/113.png",
  "humidity": 52,
  "wind_kph": 12.6,
  "wind_degree": 250,
  "wind_dir": "WSW",
  "pressure_mb": 1013,
  "visibility_km": 10.0,
  "uv": 4.0,
  "air_quality_index": 2,
  "cloud": 10,

  "forecast_days": [
    {
      "date": "2025-01-01",
      "min_temp": 12.0,
      "max_temp": 20.5,
      "avg_temp": 16.0,
      "condition": "Partiellement nuageux",
      "condition_icon_url": "https://cdn.weatherapi.com/weather/64x64/day/116.png",
      "chance_of_rain": 40,
      "chance_of_snow": 0,
      "risk_thunder": false,
      "wind_max_kph": 25.0,
      "gust_max_kph": 0,
      "sunrise": "08:30 AM",
      "sunset": "05:20 PM",
      "moon_phase": "Waxing Crescent"
    }
  ],

  "hourly": [
    {
      "time": "2025-01-01 09:00",
      "temp": 14.0,
      "condition": "Ciel dégagé",
      "chance_of_rain": 0,
      "wind_kph": 10.0,
      "gust_kph": 15.0,
      "pressure_mb": 1015,
      "uv": 3.0
    }
  ],

  "alerts": [
    {
      "type": "chaleur",
      "severity": "élevé",
      "message": "Épisode de chaleur (max 31.0°C). Pense à bien t’hydrater."
    }
  ]
}
```

> Remarque : certains champs (AQI, prévisions, alertes) dépendent des options de WeatherAPI et peuvent être absents si non disponibles.

### 400 Bad Request

Cas typiques :

- Paramètre `city` manquant.
- Ville invalide / requête vers WeatherAPI en erreur 400.

Exemple :

```json
{
  "error": "La ville saisie est invalide ou non supportée par l’API météo."
}
```

### 404 Not Found

- Aucune donnée trouvée pour la ville.

```json
{
  "error": "Aucune donnée météo trouvée pour cette ville."
}
```

### 500 Internal Server Error

- Problème de configuration serveur (clé API manquante, URL invalide).
- Erreur de décodage JSON interne.

```json
{
  "error": "Erreur interne du serveur météo. Réessaie dans quelques instants."
}
```

### 502 Bad Gateway

- WeatherAPI est en erreur ou ne répond pas correctement.

```json
{
  "error": "L’API météo externe ne répond pas correctement. Réessaie plus tard."
}
```

## Exemples de requêtes

### Exemple avec curl

```bash
curl "http://localhost:8080/api/weather?city=Paris"
```

### Exemple avec fetch (frontend)

```js
fetch("http://localhost:8080/api/weather?city=Paris")
  .then((res) => res.json())
  .then((data) => {
    console.log(data.city, data.temperature, data.condition);
  });
```
