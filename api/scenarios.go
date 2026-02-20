package api

// Scenarios is the full scenario catalog.
var Scenarios = []Scenario{
	// ── Everyday · Beginner ─────────────────────────────────
	{
		ID:          "restaurant",
		Name:        "At the Restaurant",
		Description: "Order food and drinks at a local restaurant. Practice menu vocabulary, polite requests, and dining customs.",
		Difficulty:  DifficultyBeginner,
		Category:    CategoryEveryday,
	},
	{
		ID:          "hotel",
		Name:        "Hotel Check-In",
		Description: "Check into a hotel. Handle reservations, room preferences, and common front-desk interactions.",
		Difficulty:  DifficultyBeginner,
		Category:    CategoryEveryday,
	},
	{
		ID:          "market",
		Name:        "At the Market",
		Description: "Shop at an outdoor market. Negotiate prices, ask about products, and practice numbers and quantities.",
		Difficulty:  DifficultyBeginner,
		Category:    CategoryEveryday,
	},
	{
		ID:          "cafe",
		Name:        "Morning at the Cafe",
		Description: "Order your morning coffee and pastry. Practice casual greetings and cafe vocabulary.",
		Difficulty:  DifficultyBeginner,
		Category:    CategoryEveryday,
	},
	{
		ID:          "taxi",
		Name:        "Taking a Taxi",
		Description: "Hail a taxi and give directions to your destination. Practice addresses, landmarks, and small talk.",
		Difficulty:  DifficultyBeginner,
		Category:    CategoryEveryday,
	},
	{
		ID:          "pharmacy",
		Name:        "At the Pharmacy",
		Description: "Buy medicine and basic supplies. Practice describing simple ailments and understanding dosage instructions.",
		Difficulty:  DifficultyBeginner,
		Category:    CategoryEveryday,
	},
	{
		ID:          "bakery",
		Name:        "At the Bakery",
		Description: "Browse fresh bread and pastries. Practice food vocabulary, quantities, and polite small talk.",
		Difficulty:  DifficultyBeginner,
		Category:    CategoryEveryday,
	},
	{
		ID:          "beach",
		Name:        "At the Beach",
		Description: "Rent an umbrella, order drinks at a beach bar, and chat with fellow travelers. Practice weather and leisure vocabulary.",
		Difficulty:  DifficultyBeginner,
		Category:    CategoryEveryday,
	},
	{
		ID:          "grocery",
		Name:        "At the Grocery Store",
		Description: "Navigate aisles, find items, and check out. Practice food names, quantities, and asking for help.",
		Difficulty:  DifficultyBeginner,
		Category:    CategoryEveryday,
	},
	// ── Everyday · Intermediate ─────────────────────────────
	{
		ID:          "directions",
		Name:        "Asking for Directions",
		Description: "Find your way around town. Practice location words, landmarks, and understanding directions.",
		Difficulty:  DifficultyIntermediate,
		Category:    CategoryEveryday,
	},
	{
		ID:          "doctor",
		Name:        "Visiting the Doctor",
		Description: "Describe symptoms and understand medical advice. Practice body vocabulary and health expressions.",
		Difficulty:  DifficultyIntermediate,
		Category:    CategoryEveryday,
	},
	{
		ID:          "train_station",
		Name:        "At the Train Station",
		Description: "Buy tickets and navigate a train station. Practice travel vocabulary, schedules, and platform announcements.",
		Difficulty:  DifficultyIntermediate,
		Category:    CategoryEveryday,
	},
	{
		ID:          "bank",
		Name:        "At the Bank",
		Description: "Open an account or exchange currency. Practice financial vocabulary, forms, and formal requests.",
		Difficulty:  DifficultyIntermediate,
		Category:    CategoryEveryday,
	},
	{
		ID:          "post_office",
		Name:        "At the Post Office",
		Description: "Send a package or buy stamps. Practice weights, addresses, shipping options, and polite queuing.",
		Difficulty:  DifficultyIntermediate,
		Category:    CategoryEveryday,
	},
	{
		ID:          "museum",
		Name:        "Visiting a Museum",
		Description: "Buy tickets, ask about exhibits, and discuss art. Practice descriptive language and cultural vocabulary.",
		Difficulty:  DifficultyIntermediate,
		Category:    CategoryEveryday,
	},
	{
		ID:          "airport",
		Name:        "At the Airport",
		Description: "Navigate check-in, security, and boarding. Practice travel documents, luggage, and flight vocabulary.",
		Difficulty:  DifficultyIntermediate,
		Category:    CategoryEveryday,
	},
	{
		ID:          "hair_salon",
		Name:        "At the Hair Salon",
		Description: "Describe the haircut you want. Practice appearance vocabulary, preferences, and small talk.",
		Difficulty:  DifficultyIntermediate,
		Category:    CategoryEveryday,
	},
	{
		ID:          "clothes_shop",
		Name:        "Shopping for Clothes",
		Description: "Try on clothes, ask about sizes and colors. Practice fashion vocabulary, comparisons, and returns.",
		Difficulty:  DifficultyIntermediate,
		Category:    CategoryEveryday,
	},
	// ── Everyday · Advanced ─────────────────────────────────
	{
		ID:          "job_interview",
		Name:        "Job Interview",
		Description: "Attend a job interview in a foreign country. Practice professional vocabulary, answering questions, and discussing experience.",
		Difficulty:  DifficultyAdvanced,
		Category:    CategoryEveryday,
	},
	{
		ID:          "apartment",
		Name:        "Renting an Apartment",
		Description: "Tour an apartment and negotiate a lease. Practice housing vocabulary, contract terms, and expressing requirements.",
		Difficulty:  DifficultyAdvanced,
		Category:    CategoryEveryday,
	},
	{
		ID:          "car_rental",
		Name:        "Renting a Car",
		Description: "Choose a car, understand insurance options, and handle paperwork. Practice vehicle and contract vocabulary.",
		Difficulty:  DifficultyAdvanced,
		Category:    CategoryEveryday,
	},
	{
		ID:          "police_report",
		Name:        "Reporting a Lost Item",
		Description: "File a report at a police station for a lost bag. Practice describing objects, recounting events, and formal language.",
		Difficulty:  DifficultyAdvanced,
		Category:    CategoryEveryday,
	},
	{
		ID:          "cooking_class",
		Name:        "Cooking Class",
		Description: "Follow a chef's instructions to cook a local dish. Practice kitchen vocabulary, measurements, and imperative forms.",
		Difficulty:  DifficultyAdvanced,
		Category:    CategoryEveryday,
	},
	{
		ID:          "phone_plan",
		Name:        "Buying a Phone Plan",
		Description: "Compare plans at a mobile shop and set up service. Practice tech vocabulary, plan comparisons, and contracts.",
		Difficulty:  DifficultyAdvanced,
		Category:    CategoryEveryday,
	},
	{
		ID:          "wine_tasting",
		Name:        "Wine Tasting",
		Description: "Visit a vineyard and describe flavors. Practice sensory vocabulary, opinions, and cultural conversation.",
		Difficulty:  DifficultyAdvanced,
		Category:    CategoryEveryday,
	},
	// ── Adventure · Beginner ────────────────────────────────
	{
		ID:          "dragon_market",
		Name:        "The Dragon's Market",
		Description: "Trade with a dragon merchant who hoards rare goods. Haggle for enchanted items using the target language.",
		Difficulty:  DifficultyBeginner,
		Category:    CategoryAdventure,
	},
	{
		ID:          "potion_shop",
		Name:        "The Potion Shop",
		Description: "Order potions from a mysterious apothecary. Describe your symptoms and learn magical remedy vocabulary.",
		Difficulty:  DifficultyBeginner,
		Category:    CategoryAdventure,
	},
	{
		ID:          "robot_repair_cafe",
		Name:        "Robot Repair Cafe",
		Description: "Get your companion robot fixed at a quirky repair shop. Describe malfunctions and understand tech jargon.",
		Difficulty:  DifficultyBeginner,
		Category:    CategoryAdventure,
	},
	{
		ID:          "ghost_request",
		Name:        "The Ghost's Request",
		Description: "Help a friendly ghost with unfinished business. Listen to their story and carry out their last wish.",
		Difficulty:  DifficultyBeginner,
		Category:    CategoryAdventure,
	},
	{
		ID:          "quest_board",
		Name:        "Quest Board at the Tavern",
		Description: "Accept a quest from a tavern notice board. Read postings, ask the barkeep for details, and gear up.",
		Difficulty:  DifficultyBeginner,
		Category:    CategoryAdventure,
	},
	// ── Adventure · Intermediate ────────────────────────────
	{
		ID:          "space_station_customs",
		Name:        "Space Station Customs",
		Description: "Clear immigration at an orbital station. Present documents, answer questions, and navigate bureaucracy in zero gravity.",
		Difficulty:  DifficultyIntermediate,
		Category:    CategoryAdventure,
	},
	{
		ID:          "enchanted_library",
		Name:        "The Enchanted Library",
		Description: "Find a spell in a magical library where books talk back. Navigate the catalog and decipher arcane text.",
		Difficulty:  DifficultyIntermediate,
		Category:    CategoryAdventure,
	},
	{
		ID:          "starship_briefing",
		Name:        "Starship Crew Briefing",
		Description: "Attend a mission briefing on a starship. Understand orders, ask tactical questions, and volunteer for roles.",
		Difficulty:  DifficultyIntermediate,
		Category:    CategoryAdventure,
	},
	{
		ID:          "fairy_court",
		Name:        "The Fairy Court",
		Description: "Navigate etiquette at a fairy king's court. Use formal language, make requests, and avoid giving offense.",
		Difficulty:  DifficultyIntermediate,
		Category:    CategoryAdventure,
	},
	{
		ID:          "galactic_bazaar",
		Name:        "Galactic Bazaar",
		Description: "Shop at a massive interstellar marketplace. Compare alien wares, convert currencies, and barter across species.",
		Difficulty:  DifficultyIntermediate,
		Category:    CategoryAdventure,
	},
	// ── Adventure · Advanced ────────────────────────────────
	{
		ID:          "alien_first_contact",
		Name:        "Alien First Contact",
		Description: "Communicate with a newly discovered alien species. Establish basic vocabulary and build trust through language.",
		Difficulty:  DifficultyAdvanced,
		Category:    CategoryAdventure,
	},
	{
		ID:          "time_travelers_inn",
		Name:        "Time Traveler's Inn",
		Description: "Check into an inn that exists across timelines. Speak with guests from different eras and unravel a temporal mystery.",
		Difficulty:  DifficultyAdvanced,
		Category:    CategoryAdventure,
	},
	{
		ID:          "wizard_exam",
		Name:        "The Wizard's Exam",
		Description: "Take an entrance exam at a magic academy. Answer riddles, cast verbal spells, and debate magical theory.",
		Difficulty:  DifficultyAdvanced,
		Category:    CategoryAdventure,
	},
	{
		ID:          "mech_pilot_training",
		Name:        "Mech Pilot Training",
		Description: "Learn to operate a giant mech suit. Follow instructor commands, read control labels, and pass the field test.",
		Difficulty:  DifficultyAdvanced,
		Category:    CategoryAdventure,
	},
	{
		ID:          "undersea_kingdom",
		Name:        "The Undersea Kingdom",
		Description: "Attend a diplomatic meeting in an underwater palace. Use formal speech, negotiate treaties, and respect aquatic customs.",
		Difficulty:  DifficultyAdvanced,
		Category:    CategoryAdventure,
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

// ScenariosByCategory returns all scenarios matching the given category.
func ScenariosByCategory(cat Category) []Scenario {
	var out []Scenario
	for _, s := range Scenarios {
		if s.Category == cat {
			out = append(out, s)
		}
	}
	return out
}
