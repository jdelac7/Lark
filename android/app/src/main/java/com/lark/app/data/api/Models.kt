package com.lark.app.data.api

data class Scenario(
    val id: String,
    val name: String,
    val description: String,
    val difficulty: String
)

data class Language(
    val code: String,
    val name: String
)

data class StartRequest(
    val scenarioId: String,
    val language: String
)

data class StartResponse(
    val sessionId: String,
    val message: GameMessage
)

data class PlayerInputRequest(
    val sessionId: String,
    val mode: String,
    val text: String? = null,
    val choiceIndex: Int? = null
)

data class PlayerInputResponse(
    val message: GameMessage,
    val correction: Correction? = null
)

data class GameMessage(
    val narrative: String = "",
    val translation: String = "",
    val npcDialog: String? = null,
    val npcDialogTranslation: String? = null,
    val choices: List<Choice> = emptyList(),
    val vocabulary: List<VocabItem> = emptyList(),
    val finished: Boolean = false
)

data class Choice(
    val text: String,
    val translation: String
)

data class VocabItem(
    val word: String,
    val translation: String,
    val usage: String? = null
)

data class Correction(
    val original: String,
    val corrected: String,
    val explanation: String
)

data class GameStateResponse(
    val sessionId: String,
    val scenarioId: String,
    val language: String,
    val turnCount: Int,
    val message: GameMessage
)

data class ProgressResponse(
    val playerId: String,
    val completedScenarios: List<CompletedScenario> = emptyList(),
    val vocabBank: List<VocabItem> = emptyList()
)

data class CompletedScenario(
    val scenarioId: String,
    val language: String,
    val turnCount: Int
)
