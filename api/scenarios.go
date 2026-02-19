package api

// Scenarios is the full scenario catalog.
var Scenarios = []Scenario{
	// ── Beginner ────────────────────────────────────────────
	{
		ID:          "restaurant",
		Name:        "At the Restaurant",
		Description: "Order food and drinks at a local restaurant. Practice menu vocabulary, polite requests, and dining customs.",
		Difficulty:  DifficultyBeginner,
	},
	{
		ID:          "hotel",
		Name:        "Hotel Check-In",
		Description: "Check into a hotel. Handle reservations, room preferences, and common front-desk interactions.",
		Difficulty:  DifficultyBeginner,
	},
	{
		ID:          "market",
		Name:        "At the Market",
		Description: "Shop at an outdoor market. Negotiate prices, ask about products, and practice numbers and quantities.",
		Difficulty:  DifficultyBeginner,
	},
	{
		ID:          "cafe",
		Name:        "Morning at the Cafe",
		Description: "Order your morning coffee and pastry. Practice casual greetings and cafe vocabulary.",
		Difficulty:  DifficultyBeginner,
	},
	{
		ID:          "taxi",
		Name:        "Taking a Taxi",
		Description: "Hail a taxi and give directions to your destination. Practice addresses, landmarks, and small talk.",
		Difficulty:  DifficultyBeginner,
	},
	{
		ID:          "pharmacy",
		Name:        "At the Pharmacy",
		Description: "Buy medicine and basic supplies. Practice describing simple ailments and understanding dosage instructions.",
		Difficulty:  DifficultyBeginner,
	},
	{
		ID:          "bakery",
		Name:        "At the Bakery",
		Description: "Browse fresh bread and pastries. Practice food vocabulary, quantities, and polite small talk.",
		Difficulty:  DifficultyBeginner,
	},
	{
		ID:          "beach",
		Name:        "At the Beach",
		Description: "Rent an umbrella, order drinks at a beach bar, and chat with fellow travelers. Practice weather and leisure vocabulary.",
		Difficulty:  DifficultyBeginner,
	},
	{
		ID:          "grocery",
		Name:        "At the Grocery Store",
		Description: "Navigate aisles, find items, and check out. Practice food names, quantities, and asking for help.",
		Difficulty:  DifficultyBeginner,
	},
	// ── Intermediate ───────────────────────────────────────
	{
		ID:          "directions",
		Name:        "Asking for Directions",
		Description: "Find your way around town. Practice location words, landmarks, and understanding directions.",
		Difficulty:  DifficultyIntermediate,
	},
	{
		ID:          "doctor",
		Name:        "Visiting the Doctor",
		Description: "Describe symptoms and understand medical advice. Practice body vocabulary and health expressions.",
		Difficulty:  DifficultyIntermediate,
	},
	{
		ID:          "train_station",
		Name:        "At the Train Station",
		Description: "Buy tickets and navigate a train station. Practice travel vocabulary, schedules, and platform announcements.",
		Difficulty:  DifficultyIntermediate,
	},
	{
		ID:          "bank",
		Name:        "At the Bank",
		Description: "Open an account or exchange currency. Practice financial vocabulary, forms, and formal requests.",
		Difficulty:  DifficultyIntermediate,
	},
	{
		ID:          "post_office",
		Name:        "At the Post Office",
		Description: "Send a package or buy stamps. Practice weights, addresses, shipping options, and polite queuing.",
		Difficulty:  DifficultyIntermediate,
	},
	{
		ID:          "museum",
		Name:        "Visiting a Museum",
		Description: "Buy tickets, ask about exhibits, and discuss art. Practice descriptive language and cultural vocabulary.",
		Difficulty:  DifficultyIntermediate,
	},
	{
		ID:          "airport",
		Name:        "At the Airport",
		Description: "Navigate check-in, security, and boarding. Practice travel documents, luggage, and flight vocabulary.",
		Difficulty:  DifficultyIntermediate,
	},
	{
		ID:          "hair_salon",
		Name:        "At the Hair Salon",
		Description: "Describe the haircut you want. Practice appearance vocabulary, preferences, and small talk.",
		Difficulty:  DifficultyIntermediate,
	},
	{
		ID:          "clothes_shop",
		Name:        "Shopping for Clothes",
		Description: "Try on clothes, ask about sizes and colors. Practice fashion vocabulary, comparisons, and returns.",
		Difficulty:  DifficultyIntermediate,
	},
	// ── Advanced ───────────────────────────────────────────
	{
		ID:          "job_interview",
		Name:        "Job Interview",
		Description: "Attend a job interview in a foreign country. Practice professional vocabulary, answering questions, and discussing experience.",
		Difficulty:  DifficultyAdvanced,
	},
	{
		ID:          "apartment",
		Name:        "Renting an Apartment",
		Description: "Tour an apartment and negotiate a lease. Practice housing vocabulary, contract terms, and expressing requirements.",
		Difficulty:  DifficultyAdvanced,
	},
	{
		ID:          "car_rental",
		Name:        "Renting a Car",
		Description: "Choose a car, understand insurance options, and handle paperwork. Practice vehicle and contract vocabulary.",
		Difficulty:  DifficultyAdvanced,
	},
	{
		ID:          "police_report",
		Name:        "Reporting a Lost Item",
		Description: "File a report at a police station for a lost bag. Practice describing objects, recounting events, and formal language.",
		Difficulty:  DifficultyAdvanced,
	},
	{
		ID:          "cooking_class",
		Name:        "Cooking Class",
		Description: "Follow a chef's instructions to cook a local dish. Practice kitchen vocabulary, measurements, and imperative forms.",
		Difficulty:  DifficultyAdvanced,
	},
	{
		ID:          "phone_plan",
		Name:        "Buying a Phone Plan",
		Description: "Compare plans at a mobile shop and set up service. Practice tech vocabulary, plan comparisons, and contracts.",
		Difficulty:  DifficultyAdvanced,
	},
	{
		ID:          "wine_tasting",
		Name:        "Wine Tasting",
		Description: "Visit a vineyard and describe flavors. Practice sensory vocabulary, opinions, and cultural conversation.",
		Difficulty:  DifficultyAdvanced,
	},
}

// ScenarioByID returns a scenario by its ID, or nil if not found.
func ScenarioByID(id string) *Scenario {
	for i := range Scenarios {
		if Scenarios[i].ID == id {
			return &Scenarios[i]
		}
	}
	return nil
}
