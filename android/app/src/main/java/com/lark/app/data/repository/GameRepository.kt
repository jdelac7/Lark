package com.lark.app.data.repository

import com.lark.app.data.api.*

class GameRepository(private val api: LarkApi = RetrofitClient.api) {

    suspend fun getScenarios(): List<Scenario> = api.getScenarios()

    suspend fun getLanguages(): List<Language> = api.getLanguages()

    suspend fun startScenario(scenarioId: String, language: String): StartResponse =
        api.startScenario(StartRequest(scenarioId, language))

    suspend fun sendChoice(sessionId: String, choiceIndex: Int): PlayerInputResponse =
        api.sendInput(
            PlayerInputRequest(
                sessionId = sessionId,
                mode = "choice",
                choiceIndex = choiceIndex
            )
        )

    suspend fun sendFreeText(sessionId: String, text: String): PlayerInputResponse =
        api.sendInput(
            PlayerInputRequest(
                sessionId = sessionId,
                mode = "free_text",
                text = text
            )
        )

    suspend fun getGameState(sessionId: String): GameStateResponse =
        api.getGameState(sessionId)

    suspend fun getProgress(playerId: String): ProgressResponse =
        api.getProgress(playerId)
}
