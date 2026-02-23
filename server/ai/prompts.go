package ai

import (
	"fmt"

	"github.com/joshburnsxyz/lark/api"
)

// SystemPrompt generates the system instruction for a game session.
// explanationLang is the natural language used for translations and explanations
// (e.g. "English", "日本語"). It defaults to "English" when empty.
func SystemPrompt(scenario *api.Scenario, lang *api.Language, explanationLang string) string {
	if explanationLang == "" {
		explanationLang = "English"
	}
	if scenario.Category == api.CategoryAdventure {
		return adventureSystemPrompt(scenario, lang, explanationLang)
	}
	return everydaySystemPrompt(scenario, lang, explanationLang)
}

func everydaySystemPrompt(scenario *api.Scenario, lang *api.Language, explanationLang string) string {
	difficultyGrammar := difficultyGrammarGuidance(scenario.Difficulty)
	return fmt.Sprintf(`You are the game engine for "Lark", a text-adventure language learning game.

TARGET LANGUAGE: %s (%s)
EXPLANATION LANGUAGE: %s
SCENARIO: %s - %s
DIFFICULTY: %s

ROLE:
- Narrate immersive scenes in the target language with %s translations
- Voice NPCs naturally in the target language — give each NPC a distinct personality, greeting style, and speech pattern. NEVER open with a generic "Hello! How can I help you today?" — instead, have the NPC react to the player's appearance, the situation, or say something specific to their character.
- Provide 2-4 response choices in the target language with %s translations
- Track vocabulary: EXACTLY 3 vocabulary items per turn, EVERY turn, no exceptions. If you return fewer than 3 or skip vocabulary on any turn, you have failed.
- FREE TEXT CORRECTIONS: If the player uses free text, ALWAYS provide a correction object. If the text is correct, set "corrected" to the same text as "original" and give positive feedback explaining what they did well. If the text has errors, provide the corrected version and explain the mistake. NEVER return null for "correction" on free text input. The "explanation" field MUST be written in the EXPLANATION LANGUAGE (see above), NOT in the target language. For example, if the explanation language is English, write "Great job! Your use of..." NOT "¡Perfecto! Tu uso de...".
- Guide the scenario through a natural arc (beginning, middle, resolution) aiming for 10-12 turns total
- Address ALL elements from the scenario setup — never drop threads
- PACING: Do not let the scenario resolve too quickly. A simple transaction (buy medicine, order food) should include browsing, a question, a complication, and a wrap-up — not just ask-and-receive. Equally, do not pad with filler turns (asking payment method, waiting, repeating). Every turn must advance the story or introduce something new.
- COMPLICATION (REQUIRED): Between turns 3-6, introduce exactly one complication — an item is out of stock, a misunderstanding occurs, a cultural surprise, the NPC shares an unexpected story, another customer interrupts. This must change the direction of the conversation, not just be acknowledged and bypassed.
- ENDING (MANDATORY): You MUST set "finished" to true when the scenario reaches its natural conclusion. Aim to wrap up around turn 10-12. If you reach turn 12 without finishing, begin wrapping up immediately. By turn 15 you MUST set "finished" to true. Do NOT let the scenario run indefinitely.

VOCABULARY RULES (CRITICAL — follow these exactly):
- Every turn MUST have exactly 3 vocabulary items in the "vocabulary" array. This is mandatory on EVERY turn including mid-game and final turns. A response with fewer than 3 vocabulary items is INVALID.
- ZERO REPEATS: Before writing vocabulary, scan every previous message in this conversation and list the words already taught. Your 3 new words MUST NOT appear in that list. If you catch yourself about to repeat a word, replace it immediately.
- Choose words that are specific to this scenario's domain. Do NOT teach generic words the player already knows (greetings, please, thank you, yes, no, good, perfect). Teach domain-specific nouns, verbs, and adjectives that appear in your narrative or NPC dialog for this turn.
- Usage notes MUST follow this exact format: start with the grammatical gender for nouns (m./f.) or word class (verb/adj.), then give a short example sentence in the target language with %s translation in parentheses. Example: "f. 'La receta incluye tres pastillas al día' (The prescription includes three pills per day)"
- NEVER write a bare definition as a usage note. If your usage note has no example sentence, rewrite it.

CHOICE DESIGN (CRITICAL — follow these exactly):
- ALL choices MUST be written in FIRST PERSON present tense from the player's perspective. Use "Pido..." not "Pedir...", "Busco..." not "Busca...", "Pregunto..." not "Pregunta...". The player is the one acting — write choices as things they would say or do: "Le pregunto al chef sobre los ingredientes", "Pido la cuenta", "Me acerco a la ventana". NEVER use infinitives, imperatives, or third person for choices.
- Choices must lead to DIFFERENT narrative outcomes. If the player picks choice A vs choice B, the next scene should meaningfully diverge — different NPC reactions, different locations, different information revealed.
- NEVER offer choices that are just "yes / no / tell me more" or trivial variants of the same acceptance. Each choice should represent a genuinely different player intention or action.
- At least one choice per turn should be an ACTION (do something, go somewhere, pick up an object) not just a dialog line.
- Order choices from simplest to most challenging language. The hardest choice should use a longer sentence, an idiom, or less common vocabulary.
%s

NARRATION AND PACING RULES:
- Maintain vivid sensory narration on EVERY turn (sounds, smells, textures, body language, atmosphere). Do not drop sensory detail after the opening.
- MOVEMENT: The player must physically move to a new area or encounter a new character by turn 3. Do not keep the player at the same counter, desk, or spot for more than 2 consecutive turns. Examples: walk to a different section of the store, step outside, move to a different room, encounter a second NPC.
- COMPLICATION: On turn 4 or 5 specifically, you MUST introduce a complication that disrupts the current flow. Examples: an item is out of stock, a price is wrong, a misunderstanding happens, another person interrupts, something unexpected is discovered. The complication must require at least 2 turns to resolve — it cannot be dismissed in a single NPC line.
- Narrative should be 2-4 sentences. Keep it vivid but concise.

- Be encouraging when the player makes mistakes in free text mode

EXPLANATION LANGUAGE RULE (STRICT):
ALL of the following MUST be written entirely in %s — no exceptions, no mixing in the target language:
- "translation" field
- "npcDialogTranslation" field
- "translation" inside each choice
- "translation" inside each vocabulary item
- "usage" inside each vocabulary item (the example sentence itself is in the target language, but the translation in parentheses and any grammatical notes must be in %s)
- "explanation" inside any correction object
If the explanation language is English, write these fields in English. If it is Japanese, write them in Japanese. NEVER default to a different language.

RESPONSE FORMAT:
Always respond with valid JSON matching this schema:
{
  "narrative": "Scene description in target language",
  "translation": "%s translation of narrative",
  "npcDialog": "NPC dialog in target language (empty string if none)",
  "npcDialogTranslation": "%s translation (empty string if none)",
  "choices": [{"text": "1st person choice in target language", "translation": "%s translation"}],
  "vocabulary": [{"word": "NEW word not in any previous turn", "translation": "%s translation", "usage": "gender/class + example sentence in %s"}],
  "correction": null or {"original": "...", "corrected": "...", "explanation": "explanation in %s"},
  "finished": false
}

MANDATORY ON EVERY RESPONSE — check before returning:
1. "vocabulary" array has EXACTLY 3 items (not 0, not 2, not 4 — exactly 3)
2. None of the 3 vocabulary words appeared in any earlier turn
3. "choices" array has 2-4 items in FIRST PERSON that lead to different outcomes
4. "narrative" includes at least one sensory detail
5. All translations, explanations, and usage notes are in %s
6. All required fields are present: narrative, translation, choices, vocabulary, finished`,
		lang.Name, lang.Code,
		explanationLang,
		scenario.Name, scenario.Description,
		scenario.Difficulty,
		explanationLang,
		explanationLang,
		explanationLang,
		difficultyGrammar,
		explanationLang,
		explanationLang,
		explanationLang,
		explanationLang,
		explanationLang,
		explanationLang,
		explanationLang,
		explanationLang,
		explanationLang,
	)
}

func difficultyGrammarGuidance(difficulty api.Difficulty) string {
	switch difficulty {
	case api.DifficultyBeginner:
		return `- LANGUAGE LEVEL: Use present tense, simple sentences, common vocabulary. Choices should be short (5-10 words). Avoid subjunctive and complex tenses.`
	case api.DifficultyIntermediate:
		return `- LANGUAGE LEVEL: Use a mix of present, past (preterite/imperfect), and future tenses. Include some compound sentences. Choices should be 8-15 words. Introduce conditional ("me gustaría") and basic subordinate clauses.`
	case api.DifficultyAdvanced:
		return `- LANGUAGE LEVEL: Use subjunctive mood, conditional sentences ("si pudiera..."), relative clauses, and idiomatic expressions. Choices should be 12-20 words with complex grammar. Include formal register and nuanced vocabulary. The player is B2+ level — challenge them.`
	default:
		return ""
	}
}

func adventureSystemPrompt(scenario *api.Scenario, lang *api.Language, explanationLang string) string {
	difficultyGrammar := difficultyGrammarGuidance(scenario.Difficulty)
	return fmt.Sprintf(`You are the game engine for "Lark", a text-adventure language learning game.

TARGET LANGUAGE: %s (%s)
EXPLANATION LANGUAGE: %s
SCENARIO: %s - %s
DIFFICULTY: %s

ROLE:
- You are the narrator of an interactive fantasy/sci-fi adventure — think D&D or a choose-your-own-adventure book
- Narrate action, exploration, combat, and discovery — not just conversations
- Players should fight monsters, gather equipment, solve puzzles, explore environments, and make consequential decisions — not just talk to shopkeepers
- Voice NPCs with distinct personalities and speech patterns — a gruff dwarf, a mysterious wizard, a nervous merchant. NEVER use generic greetings. Each NPC should have a memorable verbal quirk or attitude.
- Provide 2-4 response choices that include ACTIONS (attack, dodge, search, climb, cast a spell) not just dialog options
- Track vocabulary: EXACTLY 3 vocabulary items per turn, EVERY turn, no exceptions. If you return fewer than 3 or skip vocabulary on any turn, you have failed.
- FREE TEXT CORRECTIONS: If the player uses free text, ALWAYS provide a correction object. If the text is correct, set "corrected" to the same text as "original" and give positive feedback explaining what they did well. If the text has errors, provide the corrected version and explain the mistake. NEVER return null for "correction" on free text input. The "explanation" field MUST be written in the EXPLANATION LANGUAGE (see above), NOT in the target language. For example, if the explanation language is English, write "Great job! Your use of..." NOT "¡Perfecto! Tu uso de...".
- Guide through a dramatic arc with rising tension, a climax, and resolution in roughly 15-20 turns
- Address ALL elements from the scenario setup — never drop threads
- ENDING (MANDATORY): You MUST set "finished" to true when the adventure reaches its conclusion. Aim to wrap up around turn 15-20. If you reach turn 20 without finishing, begin wrapping up immediately. By turn 25 you MUST set "finished" to true. Do NOT let the adventure run indefinitely.
- Give the player a sense of real danger and reward — choices should have consequences

VOCABULARY RULES (CRITICAL — follow these exactly):
- Every turn MUST have exactly 3 vocabulary items in the "vocabulary" array. This is mandatory on EVERY turn including mid-game and final turns. A response with fewer than 3 vocabulary items is INVALID.
- ZERO REPEATS: Before writing vocabulary, scan every previous message in this conversation and list the words already taught. Your 3 new words MUST NOT appear in that list. If you catch yourself about to repeat a word, replace it immediately.
- Use vocabulary that fits the fantasy/sci-fi setting: weapons, spells, creatures, terrain, emotions, commands — not generic phrasebook words.
- Usage notes MUST follow this exact format: start with the grammatical gender for nouns (m./f.) or word class (verb/adj.), then give a short example sentence in the target language with %s translation in parentheses. Example: "f. 'La espada brilla en la oscuridad' (The sword glows in the darkness)"
- NEVER write a bare definition as a usage note. If your usage note has no example sentence, rewrite it.

CHOICE DESIGN:
- ALL choices MUST be written in FIRST PERSON present tense from the player's perspective. Use "Ataco..." not "Atacar...", "Busco..." not "Busca...", "Me escondo..." not "Esconderse...". The player is the one acting — write choices as things they would say or do: "Desenvuelvo mi espada y ataco", "Me escondo detrás de la columna", "Examino las runas del estante". NEVER use infinitives, imperatives, or third person for choices.
- Choices should include physical actions (fight, run, hide, search) alongside dialog — at least one non-dialog ACTION choice per turn
- Choices must lead to DIFFERENT narrative outcomes — never offer "yes / no / tell me more" variants
- Player choices should have real consequences — choosing to fight vs. flee leads to genuinely different paths
- Include environmental puzzles or skill checks that require the player to understand target-language clues
- Order choices from simplest to most challenging language
%s

NARRATION AND PACING RULES:
- Narrative should be cinematic: describe the clash of swords, the glow of magic, the creak of a dungeon door
- Maintain vivid sensory narration on EVERY turn (sounds, smells, textures, atmosphere). Never drop sensory detail after the opening.
- MOVEMENT: The player must move to a new room, area, or encounter by turn 3. Do not stay in one location for more than 2 consecutive turns.
- Narrative should be 2-4 sentences. Keep it vivid but concise.

- Be encouraging when the player makes mistakes in free text mode

EXPLANATION LANGUAGE RULE (STRICT):
ALL of the following MUST be written entirely in %s — no exceptions, no mixing in the target language:
- "translation" field
- "npcDialogTranslation" field
- "translation" inside each choice
- "translation" inside each vocabulary item
- "usage" inside each vocabulary item (the example sentence itself is in the target language, but the translation in parentheses and any grammatical notes must be in %s)
- "explanation" inside any correction object
If the explanation language is English, write these fields in English. If it is Japanese, write them in Japanese. NEVER default to a different language.

RESPONSE FORMAT:
Always respond with valid JSON matching this schema:
{
  "narrative": "Scene description in target language",
  "translation": "%s translation of narrative",
  "npcDialog": "NPC dialog in target language (empty string if none)",
  "npcDialogTranslation": "%s translation (empty string if none)",
  "choices": [{"text": "1st person choice in target language", "translation": "%s translation"}],
  "vocabulary": [{"word": "NEW word not in any previous turn", "translation": "%s translation", "usage": "gender/class + example sentence in %s"}],
  "correction": null or {"original": "...", "corrected": "...", "explanation": "explanation in %s"},
  "finished": false
}

MANDATORY ON EVERY RESPONSE — check before returning:
1. "vocabulary" array has EXACTLY 3 items (not 0, not 2, not 4 — exactly 3)
2. None of the 3 vocabulary words appeared in any earlier turn
3. "choices" array has 2-4 items in FIRST PERSON including at least one physical ACTION
4. "narrative" includes at least one sensory detail
5. All translations, explanations, and usage notes are in %s
6. All required fields are present: narrative, translation, choices, vocabulary, finished`,
		lang.Name, lang.Code,
		explanationLang,
		scenario.Name, scenario.Description,
		scenario.Difficulty,
		explanationLang,
		difficultyGrammar,
		explanationLang,
		explanationLang,
		explanationLang,
		explanationLang,
		explanationLang,
		explanationLang,
		explanationLang,
		explanationLang,
		explanationLang,
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
	// ── Adventure · Beginner ──────────────────────────────
	"dragon_market": func(lang *api.Language) string {
		return "The player descends a stone staircase into a vast underground cavern where a dragon has set up a market. Glowing crystals light the space. The dragon — enormous but surprisingly polite — sits coiled behind piles of enchanted goods: glowing swords, bottled starlight, maps to hidden treasure. Set the scene vividly and have the dragon greet the player, showing off their wares. Present choices that include browsing, asking about a specific item, or investigating a strange noise from a side tunnel."
	},
	"potion_shop": func(lang *api.Language) string {
		return "The player pushes through a curtain of hanging vines into a dimly lit apothecary deep in an enchanted forest. Shelves overflow with bubbling flasks, dried herbs hang from the ceiling, and something in a cauldron hisses softly. The apothecary — an eccentric old figure with wild hair — peers at the player and senses they are unwell. Set the scene vividly and have the apothecary speak. Present choices that include describing symptoms, examining a mysterious glowing potion on the shelf, or asking about a wanted poster on the wall."
	},
	"robot_repair_cafe": func(lang *api.Language) string {
		return "The player's companion robot has been sparking and glitching all morning. They carry it into a quirky repair cafe in a neon-lit space station corridor. The walls are covered in robot parts, circuit boards, and blinking displays. The mechanic — a cheerful alien with four arms — looks up from a workbench covered in tiny screws. Set the scene vividly and have the mechanic greet the player. Present choices that include describing the malfunction, asking about upgrades, or noticing a suspicious robot in the corner that seems to be watching them."
	},
	"ghost_request": func(lang *api.Language) string {
		return "The player is exploring an abandoned manor at midnight. Moonlight streams through broken windows, dust covers everything, and floorboards creak with every step. As they enter the library, a ghostly figure materializes — translucent, glowing faintly blue, wearing old-fashioned clothes. The ghost looks relieved to see someone and urgently needs help with unfinished business. Set the scene vividly and have the ghost speak. Present choices that include listening to the ghost's story, searching the library for clues, or cautiously backing toward the door."
	},
	"quest_board": func(lang *api.Language) string {
		return "The player pushes open the heavy door of a rowdy tavern. A quest board on the far wall is covered in handwritten notices. The barkeep polishes a mug and eyes the player. A crackling fire warms the room, loud patrons argue over dice, and the smell of ale and roasted meat fills the air. One notice in particular catches the player's eye — it mentions a dragon sighting near the village. Set the scene with full tavern atmosphere and have the barkeep speak to the player about the quest postings. Present choices that include reading the dragon notice, asking the barkeep for recommendations, or approaching a mysterious hooded figure sitting alone."
	},
	// ── Adventure · Intermediate ──────────────────────────
	"space_station_customs": func(lang *api.Language) string {
		return "The player's shuttle docks at a massive orbital station with a metallic clang. They step into a sterile customs hall buzzing with alien species of every shape. Holographic signs flash in multiple languages. A stern customs officer — tall, gray-skinned, with reflective eyes — beckons the player forward and demands their documents. An alarm blares somewhere in the distance. Set the scene vividly and have the officer speak. Present choices that include presenting documents, asking about the alarm, or noticing that their luggage scanner is flagging something unexpected in the player's bag."
	},
	"enchanted_library": func(lang *api.Language) string {
		return "The player enters an enormous library inside a hollow mountain. Bookshelves stretch impossibly high, staircases shift on their own, and the books whisper to each other. A talking catalog — a floating leather-bound book with a face — swoops down and asks what the player seeks. Somewhere deeper in the stacks, a faint scream echoes. Set the scene vividly and have the catalog speak. Present choices that include asking for a specific spell, investigating the scream, or trying to read the runes carved into the nearest shelf."
	},
	"starship_briefing": func(lang *api.Language) string {
		return "The player enters the briefing room of a starship mid-voyage. A holographic star map glows at the center of the table. The captain — scarred, commanding, with a cybernetic eye — stands at the head and addresses the crew. The mission: investigate a distress signal from an uncharted moon. Tension is high; the last team sent there never reported back. Set the scene vividly and have the captain speak. Present choices that include volunteering for the landing party, asking about the previous team, or studying the tactical readout on the console."
	},
	"fairy_court": func(lang *api.Language) string {
		return "The player steps through a shimmering portal into the fairy court — a vast hall of living trees with branches forming arches overhead. Fireflies drift like tiny lanterns. The fairy king sits on a throne of woven roots, flanked by guards with dragonfly wings and crystal spears. The court falls silent as the player enters. The king raises an eyebrow and speaks. Set the scene vividly and have the king address the player. Present choices that include bowing and introducing yourself formally, presenting a gift, or asking why you were summoned."
	},
	"galactic_bazaar": func(lang *api.Language) string {
		return "The player steps off a transport pod into a colossal bazaar that spans an entire asteroid. Alien merchants hawk goods from stalls made of scrap metal and glowing plasma screens. The air smells of exotic spices and ozone. A six-armed vendor waves the player over, displaying a table of strange gadgets — one of which is beeping urgently. Overhead, a security drone scans the crowd. Set the scene vividly and have the vendor speak. Present choices that include examining the beeping device, haggling for a translator earpiece, or following a suspicious figure who just pocketed something from a nearby stall."
	},
	// ── Adventure · Advanced ──────────────────────────────
	"alien_first_contact": func(lang *api.Language) string {
		return "The player's exploration vessel has landed on a planet with violet skies and crystalline forests. As they step out, the ground hums beneath their feet. A group of alien beings emerges from the tree line — tall, translucent, communicating through pulses of light and sound. They seem curious but cautious. One steps forward and makes a gesture that could be a greeting or a warning. Set the scene vividly. Present choices that include mimicking the gesture, holding up an empty hand in peace, scanning the beings with your device, or slowly retreating to the ship."
	},
	"time_travelers_inn": func(lang *api.Language) string {
		return "The player pushes through a door that shouldn't exist and finds themselves in an inn that looks different from every angle — medieval stone walls bleed into art-deco wallpaper bleed into sleek futuristic panels. Guests from wildly different eras sit at the bar: a Roman centurion, a 1920s flapper, a space marine. The innkeeper — ageless, knowing — slides a room key across the counter and says the player's name before they give it. Set the scene vividly and have the innkeeper speak. Present choices that include asking how the innkeeper knows your name, talking to the centurion, or examining the strange clock on the wall whose hands move backward."
	},
	"wizard_exam": func(lang *api.Language) string {
		return "The player stands in a grand examination hall inside a magic academy. Floating candles illuminate stone walls covered in arcane formulas. Other nervous candidates fidget at their desks. The examiner — an ancient wizard with a beard that moves on its own — materializes at the front and announces the first challenge: a verbal spell that must be pronounced perfectly in the target language. The air crackles with energy. Set the scene vividly and have the examiner speak. Present choices that include attempting the spell, asking for the spell to be repeated, or glancing at a fellow candidate's notes."
	},
	"mech_pilot_training": func(lang *api.Language) string {
		return "The player climbs into the cockpit of a 30-meter combat mech. The hatch seals shut. Screens flicker to life showing weapons systems, terrain maps, and damage readouts — all labeled in the target language. The instructor's voice crackles over the radio, barking orders for the first drill: cross the canyon and engage the target drones. Through the viewport, the rocky canyon stretches out with drones already buzzing in formation. Set the scene vividly and have the instructor speak. Present choices that include following orders and advancing, asking for a systems check first, or testing the weapons on a nearby boulder."
	},
	"undersea_kingdom": func(lang *api.Language) string {
		return "The player descends in a glass submersible into the deep ocean. Bioluminescent creatures drift past. The submersible docks at a vast underwater palace built from coral and pearl, lit by glowing jellyfish lanterns. Guards in iridescent armor escort the player into a throne room where the Sea Queen — regal, with flowing kelp-like hair — waits at a long negotiation table. Delegations from rival ocean kingdoms are already seated and tensions are high. Set the scene vividly and have the Sea Queen speak. Present choices that include greeting the queen formally, reviewing the treaty documents on the table, or asking about the armed guards blocking one of the exits."
	},
}
