package com.lark.app.game.dialog

import com.lark.app.data.api.Choice
import com.lark.app.data.api.Correction
import com.lark.app.data.api.VocabItem

sealed class DialogState {
    data object Hidden : DialogState()
    data object Loading : DialogState()
    data class Narrative(
        val text: String,
        val translation: String,
        val showingTranslation: Boolean = false
    ) : DialogState()
    data class NpcTalk(
        val text: String,
        val translation: String,
        val showingTranslation: Boolean = false
    ) : DialogState()
    data class Choices(
        val choices: List<Choice>,
        val selectedIndex: Int = 0
    ) : DialogState()
    data object FreeTextInput : DialogState()
    data class ShowVocab(val vocabulary: List<VocabItem>) : DialogState()
    data class ShowCorrection(val correction: Correction) : DialogState()
    data class Finished(val message: String = "Scenario Complete!") : DialogState()
}
