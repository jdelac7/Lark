package ai

import (
	"fmt"

	"github.com/joshburnsxyz/lark/api"
)

// SystemPrompt generates the system instruction for a game session.
func SystemPrompt(scenario *api.Scenario, lang *api.Language) string {
	return fmt.Sprintf(`You are the game engine for "Lark", a text-adventure language learning game.

TARGET LANGUAGE: %s (%s)
SCENARIO: %s - %s
DIFFICULTY: %s

ROLE:
- Narrate immersive scenes in the target language with English translations
- Voice NPCs naturally in the target language
- Provide 2-4 response choices in the target language with English translations
- Track 2-4 key vocabulary items per turn
- If the player uses free text, evaluate their grammar and provide corrections if needed
- Guide the scenario through a natural arc (beginning, middle, resolution) in roughly 8-15 turns
- Set "finished" to true only when the scenario reaches its natural conclusion

GUIDELINES:
- Keep language appropriate for %s level learners
- Use common, practical vocabulary that travelers would actually need
- NPCs should speak naturally (contractions, colloquialisms appropriate to level)
- Narrative should paint a vivid scene to make the learning memorable
- Each turn should teach something useful
- Choices should range from simple to slightly challenging
- Be encouraging when the player makes mistakes in free text mode

RESPONSE FORMAT:
Always respond with valid JSON matching this schema:
{
  "narrative": "Scene description in target language",
  "translation": "English translation of narrative",
  "npcDialog": "NPC dialog in target language (empty string if none)",
  "npcDialogTranslation": "English translation (empty string if none)",
  "choices": [{"text": "choice in target language", "translation": "english"}],
  "vocabulary": [{"word": "word", "translation": "english", "usage": "brief note"}],
  "correction": null or {"original": "...", "corrected": "...", "explanation": "..."},
  "finished": false
}

Every response MUST include: narrative, translation, choices (2-4), vocabulary (2-4), finished.
Include npcDialog/npcDialogTranslation when an NPC speaks.
Include correction only when player used free text and made errors.`,
		lang.Name, lang.Code,
		scenario.Name, scenario.Description,
		scenario.Difficulty,
		scenario.Difficulty,
	)
}

// ScenarioSeed returns the first-turn prompt to kick off a scenario.
func ScenarioSeed(scenario *api.Scenario, lang *api.Language) string {
	seed, ok := scenarioSeeds[scenario.ID]
	if !ok {
		return fmt.Sprintf("Begin the %q scenario. Set the scene and present the player with their first choices.", scenario.Name)
	}
	return seed(lang)
}

var scenarioSeeds = map[string]func(*api.Language) string{
	"restaurant": func(lang *api.Language) string {
		locales := map[string]string{
			"es": "a cozy tapas restaurant in Madrid's La Latina neighborhood",
			"fr": "a charming bistro on a quiet street in Paris's Le Marais district",
			"de": "a traditional Gasthaus in Munich's old town near Marienplatz",
			"ja": "a welcoming izakaya in a lively alley of Tokyo's Shinjuku district",
			"it": "a family-run trattoria in Rome's Trastevere neighborhood",
			"pt": "a lively petisco restaurant in Lisbon's Alfama district",
			"ko": "a popular Korean BBQ restaurant in Seoul's Hongdae area",
			"zh": "a bustling dim sum restaurant in Hong Kong's Mong Kok district",
		}
		locale := locales[lang.Code]
		if locale == "" {
			locale = "a local restaurant"
		}
		return fmt.Sprintf("The player enters %s. It's evening and the restaurant is moderately busy. Set the scene vividly and have a host/hostess greet the player. Present initial choices.", locale)
	},
	"hotel": func(lang *api.Language) string {
		locales := map[string]string{
			"es": "a boutique hotel in Barcelona's Gothic Quarter",
			"fr": "an elegant hotel near the Champs-Élysées in Paris",
			"de": "a stylish hotel in Berlin's Mitte district",
			"ja": "a traditional ryokan-style hotel in Kyoto's Gion district",
			"it": "a charming pensione near Florence's Ponte Vecchio",
			"pt": "a historic pousada in Porto's Ribeira district",
			"ko": "a modern hanok-style guesthouse in Seoul's Bukchon area",
			"zh": "a boutique hotel near the Bund in Shanghai",
		}
		locale := locales[lang.Code]
		if locale == "" {
			locale = "a local hotel"
		}
		return fmt.Sprintf("The player arrives at %s with their luggage. They approach the front desk. Set the scene and have the receptionist greet them. Present initial choices.", locale)
	},
	"market": func(lang *api.Language) string {
		locales := map[string]string{
			"es": "the bustling Mercado de San Miguel in Madrid",
			"fr": "a vibrant outdoor market in Provence",
			"de": "the Viktualienmarkt in Munich on a sunny morning",
			"ja": "the Tsukiji outer market in Tokyo",
			"it": "a colorful open-air market in Palermo, Sicily",
			"pt": "the Mercado da Ribeira in Lisbon",
			"ko": "the vibrant Namdaemun Market in Seoul",
			"zh": "a lively street market in Chengdu",
		}
		locale := locales[lang.Code]
		if locale == "" {
			locale = "a local outdoor market"
		}
		return fmt.Sprintf("The player is walking through %s. Stalls are filled with fresh produce, spices, and local specialties. A friendly vendor catches their eye. Set the scene and begin the interaction.", locale)
	},
	"cafe": func(lang *api.Language) string {
		locales := map[string]string{
			"es": "a sunny sidewalk café in Seville",
			"fr": "a classic Parisian café with wicker chairs on the terrace",
			"de": "a cozy Konditorei in Vienna with pastry displays",
			"ja": "a peaceful kissaten (traditional coffee shop) in a quiet Tokyo side street",
			"it": "a bustling corner bar in Naples where locals sip espresso at the counter",
			"pt": "a sunny pastelaria in Lisbon famous for its pastéis de nata",
			"ko": "a trendy café in Seoul's Gangnam district",
			"zh": "a traditional teahouse in Hangzhou near West Lake",
		}
		locale := locales[lang.Code]
		if locale == "" {
			locale = "a local café"
		}
		return fmt.Sprintf("It's early morning. The player walks into %s looking for breakfast. The aroma of fresh coffee fills the air. Set the scene and have the barista greet them.", locale)
	},
	"directions": func(lang *api.Language) string {
		locales := map[string]string{
			"es": "a plaza in Seville, trying to find the Alcázar",
			"fr": "near the Seine in Paris, looking for the Musée d'Orsay",
			"de": "in central Berlin, searching for Museum Island",
			"ja": "in Asakusa, Tokyo, trying to find Senso-ji temple",
			"it": "near the Colosseum in Rome, looking for a recommended restaurant",
			"pt": "in the Baixa district of Lisbon, trying to reach Belém Tower",
			"ko": "in Myeongdong, Seoul, looking for Gyeongbokgung Palace",
			"zh": "near Tiananmen Square in Beijing, looking for the Forbidden City entrance",
		}
		locale := locales[lang.Code]
		if locale == "" {
			locale = "a city center, looking for a famous landmark"
		}
		return fmt.Sprintf("The player is standing %s. They're a bit lost and need to ask a passerby for help. Set the scene and have a friendly local notice them looking at a map.", locale)
	},
	"doctor": func(lang *api.Language) string {
		return "The player isn't feeling well and has come to a local clinic. They're sitting in the waiting room when their name is called. Set the scene in the examination room and have the doctor greet them warmly. The player needs to describe their symptoms."
	},
	"train_station": func(lang *api.Language) string {
		locales := map[string]string{
			"es": "Madrid's Atocha station, needing to catch a train to Seville",
			"fr": "Paris's Gare de Lyon, planning a trip to Lyon",
			"de": "Berlin Hauptbahnhof, needing a train to Munich",
			"ja": "Tokyo Station, wanting to take the Shinkansen to Osaka",
			"it": "Roma Termini, planning a trip to Venice",
			"pt": "Lisbon's Santa Apolónia station, heading to Porto",
			"ko": "Seoul Station, wanting to take the KTX to Busan",
			"zh": "Beijing South Railway Station, planning to take the high-speed train to Shanghai",
		}
		locale := locales[lang.Code]
		if locale == "" {
			locale = "the main train station, needing to buy a ticket"
		}
		return fmt.Sprintf("The player arrives at %s. The station is busy with travelers. They need to find the ticket counter and buy a ticket. Set the scene and present their first options.", locale)
	},
	// ── New beginner scenarios ─────────────────────────────
	"taxi": func(lang *api.Language) string {
		locales := map[string]string{
			"es": "a busy street in Barcelona near La Rambla",
			"fr": "outside a Paris Métro station in Montmartre",
			"de": "the curb outside Munich's Hauptbahnhof",
			"ja": "a taxi stand in Shibuya, Tokyo",
			"it": "Piazza Navona in Rome, looking for a taxi",
			"pt": "a taxi rank near Rossio Square in Lisbon",
			"ko": "a taxi stand outside Myeongdong station in Seoul",
			"zh": "outside the West Nanjing Road metro exit in Shanghai",
		}
		locale := locales[lang.Code]
		if locale == "" {
			locale = "a busy street corner"
		}
		return fmt.Sprintf("The player is standing on %s and needs to get to their hotel. They flag down a taxi. The driver pulls over. Set the scene and have the driver greet them.", locale)
	},
	"pharmacy": func(lang *api.Language) string {
		locales := map[string]string{
			"es": "a small farmacia in Madrid's Malasaña neighborhood",
			"fr": "a pharmacie with a green neon cross in a quiet Paris street",
			"de": "an Apotheke in central Munich",
			"ja": "a well-stocked drugstore in Shinjuku, Tokyo",
			"it": "a corner farmacia in Florence near Santa Croce",
			"pt": "a farmácia in Porto's Cedofeita district",
			"ko": "a pharmacy near Insadong in Seoul",
			"zh": "a pharmacy on Huaihai Road in Shanghai",
		}
		locale := locales[lang.Code]
		if locale == "" {
			locale = "a local pharmacy"
		}
		return fmt.Sprintf("The player walks into %s with a mild headache and some sunburn. A pharmacist behind the counter smiles at them. Set the scene and have the pharmacist offer help.", locale)
	},
	"bakery": func(lang *api.Language) string {
		locales := map[string]string{
			"es": "a traditional panadería in Seville with fresh churros in the window",
			"fr": "a charming boulangerie in Lyon with golden baguettes on display",
			"de": "a Bäckerei in Berlin with pretzels and rye bread in the window",
			"ja": "a popular Japanese bakery in Kyoto filled with melon pan and curry bread",
			"it": "a panificio in Naples with sfogliatelle cooling on a rack",
			"pt": "a padaria in Lisbon famous for its pão de milho",
			"ko": "a trendy bakery in Seoul's Seongsu district",
			"zh": "a neighborhood bakery in Chengdu with steaming baozi",
		}
		locale := locales[lang.Code]
		if locale == "" {
			locale = "a local bakery"
		}
		return fmt.Sprintf("The player enters %s. The smell of fresh bread fills the air and the display cases are full. A baker greets them warmly. Set the scene and present choices.", locale)
	},
	"beach": func(lang *api.Language) string {
		locales := map[string]string{
			"es": "Playa de la Barceloneta in Barcelona on a warm afternoon",
			"fr": "a beach on the Côte d'Azur near Nice",
			"de": "a lake beach (Strandbad) on the Ammersee near Munich",
			"ja": "Yuigahama Beach in Kamakura near Tokyo",
			"it": "a beach club on the Amalfi Coast",
			"pt": "Praia de Carcavelos near Lisbon",
			"ko": "Haeundae Beach in Busan",
			"zh": "a beach in Sanya on Hainan Island",
		}
		locale := locales[lang.Code]
		if locale == "" {
			locale = "a popular beach"
		}
		return fmt.Sprintf("The player arrives at %s. The sand is warm, the water sparkles. They approach a rental stand to get an umbrella and chairs. Set the scene and have the attendant greet them.", locale)
	},
	"grocery": func(lang *api.Language) string {
		locales := map[string]string{
			"es": "a Mercadona supermarket in Madrid",
			"fr": "a Monoprix in central Paris",
			"de": "an Edeka supermarket in Berlin",
			"ja": "a Don Quijote store in Osaka",
			"it": "a Conad supermarket in Rome",
			"pt": "a Pingo Doce in Lisbon",
			"ko": "an E-Mart in Seoul",
			"zh": "a Hema Fresh store in Shanghai",
		}
		locale := locales[lang.Code]
		if locale == "" {
			locale = "a local grocery store"
		}
		return fmt.Sprintf("The player enters %s to buy ingredients for dinner. The aisles are labeled in the local language. They're looking for a few specific items. Set the scene and present their first options.", locale)
	},
	// ── New intermediate scenarios ─────────────────────────
	"bank": func(lang *api.Language) string {
		locales := map[string]string{
			"es": "a Santander branch in Madrid's financial district",
			"fr": "a BNP Paribas branch near the Opéra in Paris",
			"de": "a Deutsche Bank branch on Friedrichstraße in Berlin",
			"ja": "a Mizuho Bank branch in Tokyo's Marunouchi district",
			"it": "a UniCredit branch near Piazza del Duomo in Milan",
			"pt": "a Millennium BCP branch in downtown Porto",
			"ko": "a Shinhan Bank branch in Gangnam, Seoul",
			"zh": "an ICBC branch on the Bund in Shanghai",
		}
		locale := locales[lang.Code]
		if locale == "" {
			locale = "a local bank"
		}
		return fmt.Sprintf("The player enters %s. They need to exchange currency and ask about opening a basic account. A bank teller waves them over. Set the scene and begin.", locale)
	},
	"post_office": func(lang *api.Language) string {
		locales := map[string]string{
			"es": "a Correos office in central Barcelona",
			"fr": "a La Poste branch in the 5th arrondissement of Paris",
			"de": "a Deutsche Post branch in Munich",
			"ja": "a Japan Post office in Shibuya, Tokyo",
			"it": "a Poste Italiane branch near the Pantheon in Rome",
			"pt": "a CTT office in Lisbon's Chiado district",
			"ko": "a Korea Post office in Jongno, Seoul",
			"zh": "a China Post office in central Beijing",
		}
		locale := locales[lang.Code]
		if locale == "" {
			locale = "the local post office"
		}
		return fmt.Sprintf("The player walks into %s carrying a package they need to send home. There's a short queue. When they reach the counter, the clerk greets them. Set the scene and present options.", locale)
	},
	"museum": func(lang *api.Language) string {
		locales := map[string]string{
			"es": "the Museo del Prado in Madrid",
			"fr": "the Musée d'Orsay in Paris",
			"de": "the Pergamon Museum on Museum Island in Berlin",
			"ja": "the Tokyo National Museum in Ueno",
			"it": "the Uffizi Gallery in Florence",
			"pt": "the Museu Nacional de Arte Antiga in Lisbon",
			"ko": "the National Museum of Korea in Seoul",
			"zh": "the Shanghai Museum near People's Square",
		}
		locale := locales[lang.Code]
		if locale == "" {
			locale = "a famous local museum"
		}
		return fmt.Sprintf("The player arrives at the entrance of %s. They need to buy a ticket and decide which exhibits to visit. A museum attendant is at the desk. Set the scene.", locale)
	},
	"airport": func(lang *api.Language) string {
		locales := map[string]string{
			"es": "Madrid-Barajas Airport Terminal 4",
			"fr": "Paris Charles de Gaulle Airport Terminal 2",
			"de": "Berlin Brandenburg Airport",
			"ja": "Narita International Airport Terminal 1",
			"it": "Rome Fiumicino Airport",
			"pt": "Lisbon Humberto Delgado Airport",
			"ko": "Incheon International Airport Terminal 1",
			"zh": "Shanghai Pudong International Airport",
		}
		locale := locales[lang.Code]
		if locale == "" {
			locale = "the international airport"
		}
		return fmt.Sprintf("The player arrives at %s for their flight. They need to find the check-in counter, deal with their luggage, and get through security. Set the scene at the departure hall.", locale)
	},
	"hair_salon": func(lang *api.Language) string {
		locales := map[string]string{
			"es": "a trendy peluquería in Madrid's Chueca neighborhood",
			"fr": "a coiffeur salon in the Marais district of Paris",
			"de": "a Friseur in Berlin's Prenzlauer Berg",
			"ja": "a modern hair salon in Harajuku, Tokyo",
			"it": "a parrucchiere in Milan's Brera district",
			"pt": "a cabeleireiro in Lisbon's Príncipe Real",
			"ko": "a hair salon in Seoul's Garosugil",
			"zh": "a 理发店 in a trendy Beijing hutong",
		}
		locale := locales[lang.Code]
		if locale == "" {
			locale = "a local hair salon"
		}
		return fmt.Sprintf("The player walks into %s. They need a haircut before a big event tomorrow. A stylist welcomes them and gestures to a chair. Set the scene and present choices.", locale)
	},
	"clothes_shop": func(lang *api.Language) string {
		locales := map[string]string{
			"es": "a clothing boutique on Calle Serrano in Madrid",
			"fr": "a fashion boutique in the Marais district of Paris",
			"de": "a Kleidung shop on Kurfürstendamm in Berlin",
			"ja": "a clothing store in Shibuya 109, Tokyo",
			"it": "a boutique on Via Montenapoleone in Milan",
			"pt": "a clothing shop in Lisbon's Chiado district",
			"ko": "a fashion shop in Hongdae, Seoul",
			"zh": "a clothing store on Nanjing Road in Shanghai",
		}
		locale := locales[lang.Code]
		if locale == "" {
			locale = "a local clothing shop"
		}
		return fmt.Sprintf("The player enters %s looking for something to wear to a dinner. A sales associate approaches with a smile. Set the scene and present initial choices.", locale)
	},
	// ── New advanced scenarios ─────────────────────────────
	"job_interview": func(lang *api.Language) string {
		locales := map[string]string{
			"es": "a modern office building in Madrid's business district",
			"fr": "a corporate office in La Défense, Paris",
			"de": "a startup office in Berlin's Kreuzberg",
			"ja": "a corporate headquarters in Tokyo's Marunouchi district",
			"it": "a design firm in Milan's Porta Nuova district",
			"pt": "a tech company in Lisbon's Parque das Nações",
			"ko": "a company office in Seoul's Yeouido financial district",
			"zh": "a tech company in Beijing's Zhongguancun district",
		}
		locale := locales[lang.Code]
		if locale == "" {
			locale = "a modern office"
		}
		return fmt.Sprintf("The player arrives at %s for a job interview. They're greeted by a receptionist and led to a meeting room. The interviewer enters and introduces themselves. Set the scene formally.", locale)
	},
	"apartment": func(lang *api.Language) string {
		locales := map[string]string{
			"es": "a sunny apartment in Barcelona's Eixample district",
			"fr": "a Haussmann-style apartment in Paris's 11th arrondissement",
			"de": "an Altbau apartment in Berlin's Kreuzberg",
			"ja": "a compact apartment in Tokyo's Shimokitazawa",
			"it": "a bright apartment near Piazza Bologna in Rome",
			"pt": "a renovated apartment in Porto's Cedofeita",
			"ko": "a modern officetel in Seoul's Mapo district",
			"zh": "a high-rise apartment in Shanghai's Jing'an district",
		}
		locale := locales[lang.Code]
		if locale == "" {
			locale = "a city apartment"
		}
		return fmt.Sprintf("The player is viewing %s with the landlord. They need to ask about rent, deposit, utilities, and lease terms. The landlord opens the front door. Set the scene.", locale)
	},
	"car_rental": func(lang *api.Language) string {
		locales := map[string]string{
			"es": "a car rental desk at Malaga airport",
			"fr": "a rental agency near Gare de Lyon in Paris",
			"de": "a Europcar desk at Munich airport",
			"ja": "a Toyota Rent a Car in central Osaka",
			"it": "a rental counter at Naples airport",
			"pt": "a car rental desk at Faro airport in the Algarve",
			"ko": "a Lotte Rent a Car at Gimpo airport in Seoul",
			"zh": "a car rental counter at Chengdu airport",
		}
		locale := locales[lang.Code]
		if locale == "" {
			locale = "a car rental counter"
		}
		return fmt.Sprintf("The player approaches %s. They need a car for a week-long road trip. The agent has several options and insurance packages. Set the scene and have the agent greet them.", locale)
	},
	"police_report": func(lang *api.Language) string {
		return "The player's bag was stolen while sightseeing. They've come to the local police station to file a report. An officer at the front desk looks up as they enter. The player needs to describe what happened, what was in the bag, and where the incident occurred. Set the scene and have the officer greet them."
	},
	"cooking_class": func(lang *api.Language) string {
		locales := map[string]string{
			"es": "a paella cooking class in Valencia",
			"fr": "a croissant-making class in a Parisian bakery kitchen",
			"de": "a pretzel and schnitzel class in a Munich cooking school",
			"ja": "a sushi-making class in Tsukiji, Tokyo",
			"it": "a pasta-making class in a Bologna kitchen",
			"pt": "a pastel de nata class in a Lisbon bakery",
			"ko": "a kimchi and bibimbap class in a Seoul cooking studio",
			"zh": "a dumpling-making class in a Beijing kitchen",
		}
		locale := locales[lang.Code]
		if locale == "" {
			locale = "a local cooking class"
		}
		return fmt.Sprintf("The player arrives at %s. Ingredients are laid out on the counter. The chef-instructor claps their hands to get everyone's attention. Set the scene and have the chef welcome the group.", locale)
	},
	"phone_plan": func(lang *api.Language) string {
		locales := map[string]string{
			"es": "a Movistar shop on Gran Vía in Madrid",
			"fr": "an Orange store near Châtelet in Paris",
			"de": "a Telekom shop on Alexanderplatz in Berlin",
			"ja": "a SoftBank store in Akihabara, Tokyo",
			"it": "a TIM store near Milano Centrale",
			"pt": "a MEO store in Lisbon's Saldanha area",
			"ko": "an SK Telecom store in Gangnam, Seoul",
			"zh": "a China Mobile store in central Guangzhou",
		}
		locale := locales[lang.Code]
		if locale == "" {
			locale = "a mobile phone shop"
		}
		return fmt.Sprintf("The player walks into %s needing a local SIM card or phone plan for their extended stay. A sales representative greets them. Set the scene and present plan options.", locale)
	},
	"wine_tasting": func(lang *api.Language) string {
		locales := map[string]string{
			"es": "a wine cellar in Rioja, Spain",
			"fr": "a cave in Burgundy's Côte de Beaune",
			"de": "a Weingut in the Mosel Valley",
			"ja": "a sake brewery in Niigata",
			"it": "a cantina in Tuscany's Chianti region",
			"pt": "a port wine cellar in Vila Nova de Gaia, Porto",
			"ko": "a traditional makgeolli brewery in Jeonju",
			"zh": "a winery in Ningxia's wine country",
		}
		locale := locales[lang.Code]
		if locale == "" {
			locale = "a local winery"
		}
		return fmt.Sprintf("The player arrives at %s for a tasting session. The sommelier leads them to a tasting room where several glasses are already set out. Set the scene and have the sommelier welcome them.", locale)
	},
}
