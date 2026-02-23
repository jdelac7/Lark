# Playtest Summary: Prompt Improvements (gpt-4o-mini)

## Test Setup
- Model: openai/gpt-4o-mini via OpenRouter
- Language: Spanish (es)
- Max turns: 8, seed: 42
- 25 everyday scenarios tested (v1), 5 retested across v2/v3

## V1 Baseline Issues (original prompts)

### Critical
1. **Vocabulary disappears mid-game**: 13/20 scenarios only had vocab on turns 0-1. `apartment` taught only 4 words in 8 turns.
2. **Duplicate vocabulary**: 13/20 scenarios had dupes despite explicit prohibition. `phone_plan` repeated "plan" 3x.
3. **Usage fields are dictionary definitions**: "A place where you can buy medicine" instead of example sentences with gender.
4. **Formulaic choices**: Nearly all turns followed "yes / no / ask question" pattern. No meaningful branching.
5. **Advanced scenarios use beginner grammar**: No subjunctive, conditionals, or complex structures in advanced tier.

### Moderate
6. **No mid-scenario complications**: Only `hair_salon` had a natural complication. Most scenarios were linear transactions.
7. **Generic NPC openers**: "Hola! En que puedo ayudarte?" in 10+ scenarios.
8. **Sensory narration drops after turn 0-1**.
9. **Scenarios finish too early** (pharmacy: 4 turns) **or drag** (cafe, grocery, beach stuck at turn 8 with no progress).
10. **2 server errors**: cooking_class failed to start, wine_tasting errored at turn 5 (JSON decode failures).

### Minor
11. Vocabulary typo: "lechuaga" instead of "lechuga" in market.
12. Grammar error: "una taxi" instead of "un taxi".
13. Translation inconsistency: "With gas" vs "Sparkling" for "con gas".

## Prompt Changes Made

### Everyday prompt (`everydaySystemPrompt`)
- **Vocabulary**: Changed from "2-4 items" to "EXACTLY 3 items per turn, EVERY turn, no exceptions. A response with fewer than 3 is INVALID."
- **No-repeat enforcement**: Added "ZERO REPEATS: scan every previous message" instruction with explicit recovery step.
- **Usage format**: Required "gender/class + example sentence (translation)" with concrete example.
- **Choice design**: Banned "yes/no/tell me more" explicitly. Required at least one ACTION choice per turn.
- **Difficulty grammar tiers**: Added `difficultyGrammarGuidance()` function:
  - Beginner: present tense, short sentences (5-10 words)
  - Intermediate: past/future tenses, compound sentences (8-15 words), conditional
  - Advanced: subjunctive, conditional sentences, relative clauses, idioms (12-20 words)
- **NPC personality**: Banned generic "How can I help you?" openers. Required NPCs to react to situation/player.
- **Movement**: Required player to move to new area by turn 3. No more than 2 consecutive turns in one spot.
- **Complication**: Required specifically on turn 4 or 5, must take 2+ turns to resolve.
- **Checklist**: Added 5-point mandatory checklist at end of prompt (recency bias).
- **Vocab in schema**: Changed `"word"` placeholder to `"NEW word not in any previous turn"`.

### Adventure prompt (`adventureSystemPrompt`)
- Same vocab, usage, movement, and checklist improvements applied.

## Results: V1 → V2 → V3 Comparison

### Unique vocabulary words (average across 5 test scenarios)
| Version | Avg Unique Words | Avg Turns w/ Vocab | Count=3 Compliance |
|---------|-----------------|--------------------|--------------------|
| V1      | 8.2             | 3.2 / 9            | ~50%               |
| V2      | 10.8            | 4.4 / 9            | ~95%               |
| V3      | 14.0            | 5.4 / 9            | 100%               |

### Per-scenario detail (V3)
| Scenario      | Unique | Dupes | Turns w/ Vocab | Finished |
|---------------|--------|-------|----------------|----------|
| apartment     | 14     | 1     | 5/9            | No       |
| cafe          | 16     | 3     | 7/9            | No       |
| pharmacy      | 16     | 2     | 6/9            | No       |
| job_interview | 9      | 0     | 3/9            | No       |
| grocery       | 15     | 3     | 6/9            | No       |

### Qualitative Improvements (V3)
- **Usage fields**: Near 100% correct format (gender + example sentence). Was 0% in V1.
- **Advanced grammar**: Job interview uses subjunctive ("si pudiera"), conditional ("diría"), relative clauses. Was A2-level in V1.
- **NPC personality**: Landlord shows pride, pharmacist reacts to player's appearance, barista is warm. 3/5 improved.
- **Movement**: Apartment tours rooms (living room → kitchen). Grocery moves through aisles. Pharmacy moves to checkout.
- **Complications**: Grocery has another customer blocking shelf (turn 5). Apartment has lease negotiation pushback.

### Remaining Issues
1. **Vocab still drops off** in later turns (turns 5-8). Model complies for 4-6 turns then stops. This appears to be a gpt-4o-mini context-following limitation.
2. **Duplicates persist** at low levels (1-3 per scenario). Model doesn't reliably scan history.
3. **Complications not universal** — only ~2/5 scenarios had a real one.
4. **job_interview** had weakest vocab retention (only 3/9 turns).

## Files Changed
- `server/ai/prompts.go` — rewrote `everydaySystemPrompt`, `adventureSystemPrompt`, added `difficultyGrammarGuidance()`
- `cli/playtest.go` — new file (playtest command)
- `cli/main.go` — added `case "playtest"` route
