package com.lark.app.viewmodel

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.lark.app.data.api.CompletedScenario
import com.lark.app.data.api.VocabItem
import com.lark.app.data.repository.GameRepository
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.launch

data class ProgressUiState(
    val completedScenarios: List<CompletedScenario> = emptyList(),
    val vocabBank: List<VocabItem> = emptyList(),
    val isLoading: Boolean = true,
    val error: String? = null
)

class ProgressViewModel(
    private val repository: GameRepository = GameRepository(),
    private val playerId: String = ""
) : ViewModel() {

    private val _uiState = MutableStateFlow(ProgressUiState())
    val uiState: StateFlow<ProgressUiState> = _uiState

    fun loadProgress(playerId: String) {
        viewModelScope.launch {
            _uiState.value = _uiState.value.copy(isLoading = true, error = null)
            try {
                val progress = repository.getProgress(playerId)
                _uiState.value = ProgressUiState(
                    completedScenarios = progress.completedScenarios,
                    vocabBank = progress.vocabBank,
                    isLoading = false
                )
            } catch (e: Exception) {
                _uiState.value = _uiState.value.copy(
                    isLoading = false,
                    error = "Failed to load progress: ${e.message}"
                )
            }
        }
    }
}
