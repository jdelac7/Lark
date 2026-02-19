package com.lark.app.game.dialog

import com.badlogic.gdx.Gdx
import com.badlogic.gdx.scenes.scene2d.Stage
import com.lark.app.data.api.GameMessage
import com.lark.app.data.api.Correction
import com.lark.app.data.repository.GameRepository
import com.lark.app.game.npc.Npc
import com.lark.app.game.ui.ChoiceMenu
import com.lark.app.game.ui.CorrectionPopup
import com.lark.app.game.ui.DialogBox
import com.lark.app.game.ui.VocabPopup
import com.lark.app.game.util.PixelFont
import kotlinx.coroutines.*

class DialogManager(
    private val stage: Stage,
    private val pixelFont: PixelFont,
    private val repository: GameRepository,
    private val languageCode: String,
    private val onDialogFinished: () -> Unit
) {
    var currentState: DialogState = DialogState.Hidden
        private set

    private var sessionId: String = ""
    private var currentNpc: Npc? = null
    private var pendingStates = mutableListOf<DialogState>()
    private var currentCorrection: Correction? = null
    private var currentMessage: GameMessage? = null

    // UI actors
    private var dialogBox: DialogBox? = null
    private var choiceMenu: ChoiceMenu? = null
    private var vocabPopup: VocabPopup? = null
    private var correctionPopup: CorrectionPopup? = null

    // Coroutine scope for API calls
    private val scope = CoroutineScope(Dispatchers.IO + SupervisorJob())

    // Typewriter state
    private var typewriterText = ""
    private var typewriterTarget = ""
    private var typewriterTimer = 0f
    private var typewriterSpeed = 0.03f // seconds per character
    private var typewriterComplete = false

    fun startInteraction(npc: Npc) {
        currentNpc = npc
        setState(DialogState.Loading)

        scope.launch {
            try {
                val response = repository.startScenario(npc.scenarioId, languageCode)
                sessionId = response.sessionId
                Gdx.app.postRunnable {
                    processResponse(response.message, null)
                }
            } catch (e: Exception) {
                Gdx.app.postRunnable {
                    showDialog("Error: ${e.message ?: "Connection failed"}")
                    setState(DialogState.Hidden)
                    onDialogFinished()
                }
            }
        }
    }

    fun processResponse(message: GameMessage, correction: Correction?) {
        currentMessage = message
        currentCorrection = correction
        pendingStates.clear()

        // Build the sequence of dialog states for this response
        if (correction != null) {
            pendingStates.add(DialogState.ShowCorrection(correction))
        }

        if (message.narrative.isNotBlank()) {
            pendingStates.add(DialogState.Narrative(message.narrative, message.translation))
        }

        if (!message.npcDialog.isNullOrBlank()) {
            pendingStates.add(
                DialogState.NpcTalk(
                    message.npcDialog,
                    message.npcDialogTranslation ?: ""
                )
            )
        }

        if (message.vocabulary.isNotEmpty()) {
            pendingStates.add(DialogState.ShowVocab(message.vocabulary))
        }

        if (message.finished) {
            pendingStates.add(DialogState.Finished())
        } else if (message.choices.isNotEmpty()) {
            pendingStates.add(DialogState.Choices(message.choices))
        }

        advanceToNextPending()
    }

    fun advance() {
        when (val state = currentState) {
            is DialogState.Loading -> { /* Wait */ }

            is DialogState.Narrative -> {
                if (!typewriterComplete) {
                    // Skip typewriter, show full text
                    typewriterComplete = true
                    typewriterText = typewriterTarget
                    updateDialogBoxText(typewriterTarget)
                } else if (!state.showingTranslation && state.translation.isNotBlank()) {
                    // Show translation
                    val newState = state.copy(showingTranslation = true)
                    currentState = newState
                    updateDialogBoxText("${state.text}\n\n${state.translation}")
                } else {
                    advanceToNextPending()
                }
            }

            is DialogState.NpcTalk -> {
                if (!typewriterComplete) {
                    typewriterComplete = true
                    typewriterText = typewriterTarget
                    updateDialogBoxText(typewriterTarget)
                } else if (!state.showingTranslation && state.translation.isNotBlank()) {
                    val newState = state.copy(showingTranslation = true)
                    currentState = newState
                    val npcName = currentNpc?.displayName ?: "NPC"
                    updateDialogBoxText("$npcName:\n${state.text}\n\n${state.translation}")
                } else {
                    advanceToNextPending()
                }
            }

            is DialogState.Choices -> {
                // A press on choices = confirm selection
                val choice = state.choices.getOrNull(state.selectedIndex) ?: return
                hideAllUI()
                setState(DialogState.Loading)
                scope.launch {
                    try {
                        val response = repository.sendChoice(sessionId, state.selectedIndex)
                        Gdx.app.postRunnable {
                            processResponse(response.message, response.correction)
                        }
                    } catch (e: Exception) {
                        Gdx.app.postRunnable {
                            showDialog("Error: ${e.message}")
                            setState(DialogState.Hidden)
                            onDialogFinished()
                        }
                    }
                }
            }

            is DialogState.FreeTextInput -> { /* Handled by text input callback */ }

            is DialogState.ShowVocab -> advanceToNextPending()

            is DialogState.ShowCorrection -> advanceToNextPending()

            is DialogState.Finished -> {
                currentNpc?.completed = true
                hideAllUI()
                setState(DialogState.Hidden)
                onDialogFinished()
            }

            is DialogState.Hidden -> { /* Nothing */ }
        }
    }

    fun moveChoiceUp() {
        val state = currentState
        if (state is DialogState.Choices && state.selectedIndex > 0) {
            val newState = state.copy(selectedIndex = state.selectedIndex - 1)
            currentState = newState
            choiceMenu?.selectedIndex = newState.selectedIndex
        }
    }

    fun moveChoiceDown() {
        val state = currentState
        if (state is DialogState.Choices && state.selectedIndex < state.choices.size) {
            // +1 for "Type your own..." option
            val newState = state.copy(selectedIndex = state.selectedIndex + 1)
            currentState = newState
            choiceMenu?.selectedIndex = newState.selectedIndex
        }
    }

    fun submitFreeText(text: String) {
        if (text.isBlank()) return
        hideAllUI()
        setState(DialogState.Loading)

        scope.launch {
            try {
                val response = repository.sendFreeText(sessionId, text)
                Gdx.app.postRunnable {
                    processResponse(response.message, response.correction)
                }
            } catch (e: Exception) {
                Gdx.app.postRunnable {
                    showDialog("Error: ${e.message}")
                    setState(DialogState.Hidden)
                    onDialogFinished()
                }
            }
        }
    }

    fun update(delta: Float) {
        // Typewriter effect
        if (!typewriterComplete && typewriterTarget.isNotEmpty()) {
            typewriterTimer += delta
            val charsToShow = (typewriterTimer / typewriterSpeed).toInt()
            if (charsToShow >= typewriterTarget.length) {
                typewriterText = typewriterTarget
                typewriterComplete = true
            } else {
                typewriterText = typewriterTarget.substring(0, charsToShow)
            }
            updateDialogBoxText(typewriterText)
        }
    }

    private fun advanceToNextPending() {
        hideAllUI()

        if (pendingStates.isEmpty()) {
            // No more pending — back to hidden
            setState(DialogState.Hidden)
            onDialogFinished()
            return
        }

        val next = pendingStates.removeAt(0)
        setState(next)

        when (next) {
            is DialogState.Loading -> showDialog("...")

            is DialogState.Narrative -> {
                startTypewriter(next.text)
                showDialogBox()
            }

            is DialogState.NpcTalk -> {
                val npcName = currentNpc?.displayName ?: "NPC"
                startTypewriter("$npcName:\n${next.text}")
                showDialogBox()
            }

            is DialogState.Choices -> {
                showChoiceMenu(next)
            }

            is DialogState.FreeTextInput -> {
                // Show text input - handled via Android keyboard
                Gdx.input.getTextInput(
                    object : com.badlogic.gdx.Input.TextInputListener {
                        override fun input(text: String?) {
                            if (text != null) {
                                Gdx.app.postRunnable { submitFreeText(text) }
                            }
                        }
                        override fun canceled() {
                            Gdx.app.postRunnable {
                                // Go back to choices if available
                                currentMessage?.let { msg ->
                                    if (msg.choices.isNotEmpty()) {
                                        val choiceState = DialogState.Choices(msg.choices)
                                        pendingStates.add(0, choiceState)
                                        advanceToNextPending()
                                    }
                                }
                            }
                        }
                    },
                    "Type your response",
                    "",
                    ""
                )
            }

            is DialogState.ShowVocab -> showVocabPopup(next)

            is DialogState.ShowCorrection -> showCorrectionPopup(next)

            is DialogState.Finished -> showDialog(next.message)

            is DialogState.Hidden -> { /* Should not happen */ }
        }
    }

    private fun setState(state: DialogState) {
        currentState = state
    }

    private fun startTypewriter(text: String) {
        typewriterTarget = text
        typewriterText = ""
        typewriterTimer = 0f
        typewriterComplete = false
    }

    private fun showDialog(text: String) {
        hideAllUI()
        showDialogBox()
        updateDialogBoxText(text)
        typewriterComplete = true
    }

    private fun showDialogBox() {
        if (dialogBox == null) {
            dialogBox = DialogBox(pixelFont)
            stage.addActor(dialogBox)
        }
        dialogBox?.isVisible = true
        positionDialogBox()
    }

    private fun updateDialogBoxText(text: String) {
        dialogBox?.setText(text)
    }

    private fun showChoiceMenu(state: DialogState.Choices) {
        choiceMenu = ChoiceMenu(pixelFont, state.choices, state.selectedIndex) { index ->
            if (index == state.choices.size) {
                // "Type your own..." selected
                hideAllUI()
                setState(DialogState.FreeTextInput)
                advanceToNextPending()
                // Re-trigger free text since we just set the state
                Gdx.input.getTextInput(
                    object : com.badlogic.gdx.Input.TextInputListener {
                        override fun input(text: String?) {
                            if (text != null) {
                                Gdx.app.postRunnable { submitFreeText(text) }
                            }
                        }
                        override fun canceled() {
                            Gdx.app.postRunnable {
                                pendingStates.add(0, state)
                                advanceToNextPending()
                            }
                        }
                    },
                    "Type your response",
                    "",
                    ""
                )
            } else {
                // Choice selected via tap
                val s = currentState
                if (s is DialogState.Choices) {
                    currentState = s.copy(selectedIndex = index)
                    advance()
                }
            }
        }
        stage.addActor(choiceMenu)
        positionChoiceMenu()
    }

    private fun showVocabPopup(state: DialogState.ShowVocab) {
        vocabPopup = VocabPopup(pixelFont, state.vocabulary)
        stage.addActor(vocabPopup)
        positionVocabPopup()
    }

    private fun showCorrectionPopup(state: DialogState.ShowCorrection) {
        correctionPopup = CorrectionPopup(pixelFont, state.correction)
        stage.addActor(correctionPopup)
        positionCorrectionPopup()
    }

    private fun hideAllUI() {
        dialogBox?.remove()
        dialogBox = null
        choiceMenu?.remove()
        choiceMenu = null
        vocabPopup?.remove()
        vocabPopup = null
        correctionPopup?.remove()
        correctionPopup = null
    }

    private fun getGameAreaBottom(): Float {
        // Control panel takes bottom 42% of screen
        return Gdx.graphics.height * 0.42f
    }

    private fun positionDialogBox() {
        val sw = Gdx.graphics.width.toFloat()
        val panelTop = getGameAreaBottom()
        val boxW = sw * 0.92f
        val boxH = Gdx.graphics.height * 0.22f
        dialogBox?.setBounds((sw - boxW) / 2, panelTop + 12f, boxW, boxH)
    }

    private fun positionChoiceMenu() {
        val sw = Gdx.graphics.width.toFloat()
        val sh = Gdx.graphics.height.toFloat()
        val panelTop = getGameAreaBottom()
        val gameH = sh - panelTop
        val boxW = sw * 0.85f
        val boxH = gameH * 0.75f
        choiceMenu?.setBounds((sw - boxW) / 2, panelTop + (gameH - boxH) / 2, boxW, boxH)
    }

    private fun positionVocabPopup() {
        val sw = Gdx.graphics.width.toFloat()
        val sh = Gdx.graphics.height.toFloat()
        val panelTop = getGameAreaBottom()
        val gameH = sh - panelTop
        val boxW = sw * 0.85f
        val boxH = gameH * 0.70f
        vocabPopup?.setBounds((sw - boxW) / 2, panelTop + (gameH - boxH) / 2, boxW, boxH)
    }

    private fun positionCorrectionPopup() {
        val sw = Gdx.graphics.width.toFloat()
        val sh = Gdx.graphics.height.toFloat()
        val panelTop = getGameAreaBottom()
        val gameH = sh - panelTop
        val boxW = sw * 0.85f
        val boxH = gameH * 0.60f
        correctionPopup?.setBounds((sw - boxW) / 2, panelTop + (gameH - boxH) / 2, boxW, boxH)
    }

    fun dispose() {
        scope.cancel()
        hideAllUI()
    }
}
