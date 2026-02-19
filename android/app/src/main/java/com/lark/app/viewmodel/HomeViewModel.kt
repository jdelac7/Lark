package com.lark.app.viewmodel

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.lark.app.data.api.Language
import com.lark.app.data.api.Scenario
import com.lark.app.data.repository.GameRepository
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.launch

data class HomeUiState(
    val scenarios: List<Scenario> = emptyList(),
    val languages: List<Language> = emptyList(),
    val selectedLanguage: Language? = null,
    val isLoading: Boolean = true,
    val error: String? = null
)

class HomeViewModel(
    private val repository: GameRepository = GameRepository()
) : ViewModel() {

    private val _uiState = MutableStateFlow(HomeUiState())
    val uiState: StateFlow<HomeUiState> = _uiState

    init {
        loadData()
    }

    fun loadData() {
        viewModelScope.launch {
            _uiState.value = _uiState.value.copy(isLoading = true, error = null)
            try {
                val scenarios = repository.getScenarios()
                val languages = repository.getLanguages()
                _uiState.value = _uiState.value.copy(
                    scenarios = scenarios,
                    languages = languages,
                    selectedLanguage = languages.firstOrNull(),
                    isLoading = false
                )
            } catch (e: Exception) {
                _uiState.value = _uiState.value.copy(
                    isLoading = false,
                    error = "Failed to connect to server: ${e.message}"
                )
            }
        }
    }

    fun selectLanguage(language: Language) {
        _uiState.value = _uiState.value.copy(selectedLanguage = language)
    }
}
