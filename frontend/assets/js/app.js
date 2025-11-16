// Logique page d'accueil (index.html)
document.addEventListener("DOMContentLoaded", () => {
	// Pour l'instant, la page d'accueil n'a pas de JS spécifique.
	// Si tu ajoutes plus tard des interactions, mets-les ici
	// en vérifiant la présence d'éléments propres à index.html.
});

// Logique page météo (weather.html)
document.addEventListener("DOMContentLoaded", () => {
	const cityInput = document.getElementById("city-input");
	const searchBtn = document.getElementById("search-btn");
	const result = document.getElementById("result");

	if (!cityInput || !searchBtn || !result) {
		// On n'est pas sur la page météo, on ne fait rien ici.
		return;
	}

	// --- Sélecteurs page météo (weather.html) ---
	const inputPage = document.getElementById("city-input");
	const btnPage = document.getElementById("search-btn");
	const resultPage = document.getElementById("result");
	const worldMap = document.getElementById("world-weather-map");

	// Fonction générique de rendu pour une carte météo
	const makeWeatherCardHTML = (data) => {
		const temp = data.temperature ?? data.TempC ?? data.temp_c;
		const condition = data.condition ?? data.Condition ?? "";
		const iconUrl = data.condition_icon_url ?? data.icon ?? "";

		// Bloc d'infos supplémentaires sur la ville
		const locationInfo = `
			<div class="weather-location-extra">
				<p><strong>Pays :</strong> ${data.country || "Inconnu"}</p>
				<p><strong>Région :</strong> ${data.region || "Inconnue"}</p>
				<p><strong>Coordonnées :</strong> ${
					data.latitude !== undefined && data.longitude !== undefined
						? data.latitude.toFixed(2) + "° / " + data.longitude.toFixed(2) + "°"
						: "Non disponibles"
				}</p>
			</div>
		`;

		// Détails météo
		const details = `
			<div class="weather-details">
				<p><strong>Température ressentie :</strong> ${
					data.feels_like !== undefined ? Math.round(data.feels_like) + "°C" : "NC"
				}</p>
				<p><strong>Humidité :</strong> ${
					data.humidity !== undefined ? data.humidity + " %" : "NC"
				}</p>
				<p><strong>Vent :</strong> ${
					data.wind_kph !== undefined
						? data.wind_kph + " km/h" + (data.wind_dir ? " (" + data.wind_dir + ")" : "")
						: "NC"
				}</p>
				<p><strong>Pression :</strong> ${
					data.pressure_mb !== undefined ? data.pressure_mb + " hPa" : "NC"
				}</p>
				<p><strong>Visibilité :</strong> ${
					data.visibility_km !== undefined ? data.visibility_km + " km" : "NC"
				}</p>
				<p><strong>Indice UV :</strong> ${data.uv !== undefined ? data.uv : "NC"}</p>
			</div>
		`;

		return `
			<div class="weather-card">
				<h2>${data.city || data.name || "Ville"}</h2>
				<p class="weather-location">Conditions actuelles</p>
				<div class="weather-main">
					<span class="weather-temp">${Math.round(temp)}°C</span>
					<span class="weather-condition">${condition}</span>
					${
						iconUrl
							? `<div class="weather-icon-wrapper">
									<img src="${iconUrl}" alt="${condition}" />
							   </div>`
							: ""
					}
				</div>
				${locationInfo}
				${details}
			</div>
		`;
	};

	// ------ LOGIQUE POUR LA PAGE METEO (weather.html) ------

	if (inputPage && btnPage && resultPage) {
		const renderLoadingPage = (city) => {
			resultPage.className = "status-message";
			resultPage.innerHTML = `Chargement de la météo pour <strong>${city}</strong>...`;
		};

		const renderErrorPage = (userMessage) => {
			resultPage.className = "status-message error";
			resultPage.innerHTML = `
				<p>${userMessage}</p>
			`;
		};

		const applyWeatherAnimation = (card, conditionRaw) => {
			if (!card || !conditionRaw) return;
			const c = conditionRaw.toLowerCase();

			card.classList.remove("sunny", "rainy", "cloudy", "stormy");

			if (c.includes("soleil") || c.includes("sunny") || c.includes("clear")) {
				card.classList.add("sunny");
			} else if (c.includes("pluie") || c.includes("rain") || c.includes("drizzle")) {
				card.classList.add("rainy");
			} else if (c.includes("orage") || c.includes("thunder")) {
				card.classList.add("stormy");
			} else if (c.includes("nuage") || c.includes("cloud")) {
				card.classList.add("cloudy");
			}
		};

		const renderWeatherPage = (data) => {
			resultPage.className = "";
			resultPage.innerHTML = makeWeatherCardHTML(data);

			const card = resultPage.querySelector(".weather-card");
			const condition = data.condition ?? data.Condition ?? "";
			applyWeatherAnimation(card, condition);
		};

		btnPage.addEventListener("click", async () => {
			const city = inputPage.value.trim();
			if (!city) {
				renderErrorPage("Merci d’entrer une ville ou un village.");
				return;
			}

			renderLoadingPage(city);

			try {
				const res = await fetch(
					`http://localhost:8080/api/weather?city=${encodeURIComponent(city)}`
				);

				if (!res.ok) {
					// On essaie de lire un JSON { "error": "..." }
					let backendMsg = "";
					try {
						const data = await res.json();
						if (data && data.error) backendMsg = data.error;
					} catch {
						backendMsg = "";
					}

					// Message dédié en fonction du code HTTP
					if (res.status === 400) {
						renderErrorPage(
							"Requête invalide : vérifie le nom de la ville ou le format de la demande."
						);
					} else if (res.status === 404) {
						renderErrorPage("Ville introuvable dans l’API météo.");
					} else if (res.status === 500) {
						renderErrorPage(
							"Erreur interne du serveur météo. Réessaie dans quelques instants."
						);
					} else {
						renderErrorPage(
							`Erreur inattendue (${res.status})${
								backendMsg ? " : " + backendMsg : ""
							}`
						);
					}
					return;
				}

				const data = await res.json();
				renderWeatherPage(data);
			} catch (e) {
				renderErrorPage(
					"Erreur de connexion au serveur backend (Go). Vérifie qu’il est bien démarré."
				);
			}
		});
	}

	// --- "Carte du monde" météo simplifiée ---
	if (worldMap) {
		// Données de villes globales fictives (pour la carte)
		const globalCities = [
			{ name: "Paris", temp: 18, region: "europe" },
			{ name: "New York", temp: 22, region: "americas" },
			{ name: "Tokyo", temp: 25, region: "asia" },
			{ name: "Sydney", temp: 20, region: "oceania" },
			{ name: "Le Caire", temp: 30, region: "africa" }
		];

		const getTempColor = (t) => {
			// renvoie une classe selon la température
			if (t <= 0) return "temp-cold";
			if (t <= 10) return "temp-fresh";
			if (t <= 20) return "temp-mild";
			if (t <= 30) return "temp-warm";
			return "temp-hot";
		};

		worldMap.innerHTML = globalCities
			.map(
				(c) => `
					<div class="world-city world-city-${c.region} ${getTempColor(c.temp)}">
						<span class="world-city-name">${c.name}</span>
						<span class="world-city-temp">${c.temp}°C</span>
					</div>
				`
			)
			.join("");
	}

	// --- Initialisation Leaflet: carte du monde ---
	if (worldMap && typeof L !== "undefined") {
		// Créer la carte centrée sur [20, 0] avec un zoom global
		const map = L.map(worldMap).setView([20, 0], 2);

		// Tu peux changer de tiles provider si besoin (OpenStreetMap ici)
		L.tileLayer("https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png", {
			maxZoom: 19,
			attribution: "&copy; OpenStreetMap contributors",
		}).addTo(map);

		// Données de quelques grandes villes (coordonnées + temp estimée)
		const globalCities = [
			{ name: "Paris", temp: 18, lat: 48.8566, lon: 2.3522 },
			{ name: "New York", temp: 22, lat: 40.7128, lon: -74.006 },
			{ name: "Tokyo", temp: 25, lat: 35.6895, lon: 139.6917 },
			{ name: "Sydney", temp: 20, lat: -33.8688, lon: 151.2093 },
			{ name: "Le Caire", temp: 30, lat: 30.0444, lon: 31.2357 },
		];

		const getTempColor = (t) => {
			if (t <= 0) return "#0ea5e9";
			if (t <= 10) return "#22c55e";
			if (t <= 20) return "#eab308";
			if (t <= 30) return "#f97316";
			return "#ef4444";
		};

		globalCities.forEach((c) => {
			const marker = L.circleMarker([c.lat, c.lon], {
				radius: 8,
				fillColor: getTempColor(c.temp),
				color: "#0b1120",
				weight: 1,
				opacity: 1,
				fillOpacity: 0.9,
			}).addTo(map);

			marker.bindPopup(
				`<strong>${c.name}</strong><br/>Température approx. : ${c.temp}°C`
			);
		});
	}
});
