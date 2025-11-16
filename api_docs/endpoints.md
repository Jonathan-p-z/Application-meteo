# Endpoints

## GET /api/health
Retourne l'état du service.

## GET /api/weather?city={name}
Retourne la météo simulée pour la ville donnée.
Réponse (200):
```json
{
  "city": "Paris",
  "temperature": 21.5,
  "condition": "Sunny in Paris"
}
```
