package com.lark.app.data.api

import retrofit2.http.*

interface LarkApi {

    @GET("api/v1/health")
    suspend fun health(): Map<String, String>

    @GET("api/v1/languages")
    suspend fun getLanguages(): List<Language>

    @GET("api/v1/scenarios")
    suspend fun getScenarios(): List<Scenario>

    @POST("api/v1/scenarios/start")
    suspend fun startScenario(@Body request: StartRequest): StartResponse

    @POST("api/v1/game/input")
    suspend fun sendInput(@Body request: PlayerInputRequest): PlayerInputResponse

    @GET("api/v1/game/state")
    suspend fun getGameState(@Query("sessionId") sessionId: String): GameStateResponse

    @GET("api/v1/progress")
    suspend fun getProgress(@Query("playerId") playerId: String): ProgressResponse
}
