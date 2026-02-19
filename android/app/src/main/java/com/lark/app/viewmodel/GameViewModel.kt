package com.lark.app.viewmodel

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.lark.app.data.api.*
import com.lark.app.data.repository.GameRepository
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.launch

data class ChatEntry(
    val message: GameMessage,
    val correction: Correction? = null,
    val playerInput: String? = null,
    val isPlayerTurn: Boolean = false
)

data class GameUiState(
    val sessionId: String = "",
    val scenarioName: String = "",
    val languageName: String = "",
    val chatHistory: List<ChatEntry> = emptyList(),
    val currentMessage: GameMessage? = null,
    val isLoading: Boolean = false,
    val isFinished: Boolean = false,
    val error: String? = null
)

class GameViewModel(
    private val repository: GameRepository = GameRepository()
) : ViewModel() {

    private val _uiState = MutableStateFlow(GameUiState())
    val uiState: StateFlow<GameUiState> = _uiState

    fun startScenario(scenarioId: String, scenarioName: String, languageCode: String, languageName: String) {
        viewModelScope.launch {
            _uiState.value = GameUiState(
                scenarioName = scenarioName,
                languageName = languageName,
                isLoading = true
            )
            try {
                val response = repository.startScenario(scenarioId, languageCode)
                val entry = ChatEntry(message = response.message)
                _uiState.value = _uiState.value.copy(
                    sessionId = response.sessionId,
                    chatHistory = listOf(entry),
                    currentMessage = response.message,
                    isLoading = false,
                    isFinished = response.message.finished
                )
            } catch (e: Exception) {
                _uiState.value = _uiState.value.copy(
                    isLoading = false,
                    error = "Failed to start scenario: ${e.message}"
                )
            }
        }
    }

    fun sendChoice(choiceIndex: Int) {
        val state = _uiState.value
        if (state.sessionId.isEmpty() || state.isLoading) return

        val choiceText = state.currentMessage?.choices?.getOrNull(choiceIndex)?.text ?: return

        viewModelScope.launch {
            _uiState.value = state.copy(isLoading = true, error = null)

            // Add player's choice to history
            val playerEntry = ChatEntry(
                message = GameMessage(),
                playerInput = choiceText,
                isPlayerTurn = true
            )

            try {
                val response = repository.sendChoice(state.sessionId, choiceIndex)
                val responseEntry = ChatEntry(
                    message = response.message,
                    correction = response.correction
                )
                _uiState.value = state.copy(
                    chatHistory = state.chatHistory + playerEntry + responseEntry,
                    currentMessage = response.message,
                    isLoading = false,
                    isFinished = response.message.finished
                )
            } catch (e: Exception) {
                _uiState.value = state.copy(
                    isLoading = false,
                    error = "Failed to send input: ${e.message}"
                )
            }
        }
    }

    fun sendFreeText(text: String) {
        val state = _uiState.value
        if (state.sessionId.isEmpty() || state.isLoading || text.isBlank()) return

        viewModelScope.launch {
            _uiState.value = state.copy(isLoading = true, error = null)

            val playerEntry = ChatEntry(
                message = GameMessage(),
                playerInput = text,
                isPlayerTurn = true
            )

            try {
                val response = repository.sendFreeText(state.sessionId, text)
                val responseEntry = ChatEntry(
                    message = response.message,
                    correction = response.correction
                )
                _uiState.value = state.copy(
                    chatHistory = state.chatHistory + playerEntry + responseEntry,
                    currentMessage = response.message,
                    isLoading = false,
                    isFinished = response.message.finished
                )
            } catch (e: Exception) {
                _uiState.value = state.copy(
                    isLoading = false,
                    error = "Failed to send input: ${e.message}"
                )
            }
        }
    }

    fun clearError() {
        _uiState.value = _uiState.value.copy(error = null)
    }
}
