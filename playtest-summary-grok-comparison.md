# Playtest Summary: gpt-4o-mini vs grok-4.1-fast

## Test Setup
- Language: Spanish (es), 25 everyday scenarios, --max-turns 8, --seed 42
- gpt-4o-mini tested with original prompts (v1) and improved prompts (v3)
- grok-4.1-fast tested with improved prompts + reasoning disabled + max_tokens 4096

## Head-to-Head Comparison

| Dimension                | gpt-4o-mini (v1) | gpt-4o-mini (v3) | grok-4.1-fast |
|--------------------------|-------------------|-------------------|---------------|
| **Avg unique vocab**     | 8.2               | 14.0              | **17.4**      |
| **Turns with vocab**     | 36%               | 60%               | **97%**       |
| **Exactly 3 vocab/turn** | ~50%              | 100%              | **97%**       |
| **Avg dupes/scenario**   | 1.3               | 1.6               | 5.0           |
| **Usage field format**   | 0% correct        | ~95% correct      | **~70% correct** |
| **Choice diversity**     | Poor (yes/no)     | Improved           | **Excellent** |
| **NPC personality**      | Generic           | 60% improved       | **80% distinct** |
| **Complications**        | 0%                | ~40%               | **100%**      |
| **Movement**             | None              | 20%                | **83%**       |
| **Advanced grammar**     | A2 level          | B2+                | **B2+**       |
| **Sensory narration**    | Drops after turn 1| Partial             | **Persistent** |
| **Errors**               | 2/25              | 0/5                | 3/25          |

## Grok 4.1 Fast Wins

1. **Vocabulary every turn**: 97% of turns have exactly 3 vocab items (vs 60% for gpt-4o-mini v3). This was the #1 problem — gpt-4o-mini drops vocab after turn 1-2. Grok follows the instruction consistently.

2. **Complications in every scenario**: All 12 sampled scenarios had genuine complications (out-of-stock items, equipment failures, interruptions, bureaucratic obstacles). Some had double complications. gpt-4o-mini v3 only achieved ~40%.

3. **Spatial movement**: Players move through 3-6 distinct locations per scenario. gpt-4o-mini kept players at one counter for entire games.

4. **Sensory narration persists**: Rich multi-sensory detail (smell, sound, texture, body language) continues through turns 4-8. gpt-4o-mini dropped sensory detail after turn 0-1.

5. **NPC personality**: Characters have quirks, speech patterns, physical descriptions. "Your hair is a disaster, what are we fixing today?" vs generic "How can I help you?"

6. **Choice diversity**: 4 genuine options per turn with narrative-branching outcomes. Physical actions mixed with dialog. Advanced choices use subjunctive/conditional. No yes/no/ask patterns.

## Grok 4.1 Fast Weaknesses

1. **Higher duplicate vocabulary**: 5 dupes/scenario vs 1.6 for gpt-4o-mini v3. The model repeats core topic words (e.g. "permanencia" 6x in phone_plan). This is the main area where gpt-4o-mini is better.

2. **Occasional JSON errors**: 3/25 scenarios errored (2 truncation, 1 schema deviation). After bumping max_tokens to 4096, this dropped to ~2/25. gpt-4o-mini had 2/25 errors.

3. **Usage field content errors**: ~70% compliance vs ~95% for gpt-4o-mini v3. Format is correct (gender + sentence) but occasional content errors: invented words ("salmado" for "salado"), French words ("reposer" for "reposar"), split words ("plan cha" for "plancha").

## Recommendation

**Switch to grok-4.1-fast**. The improvements in vocab coverage, complications, movement, NPC personality, and narration quality are transformative for the game experience. The duplicate vocabulary issue is real but less impactful than missing vocabulary entirely (gpt-4o-mini's main failure).

## Changes Made
- `.env`: OPENROUTER_MODEL=x-ai/grok-4.1-fast
- `server/ai/openrouter.go`: Added `reasoning` field (disabled by default), increased max_tokens to 4096
- `server/ai/prompts.go`: Rewrote everyday + adventure prompts with stronger vocab/choice/movement/complication rules + difficulty grammar tiers
- `cli/playtest.go`: New playtest command
- `cli/main.go`: Added playtest route
