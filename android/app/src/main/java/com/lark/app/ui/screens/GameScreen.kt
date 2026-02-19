package com.lark.app.ui.screens

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.lazy.rememberLazyListState
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.ArrowBack
import androidx.compose.material.icons.filled.Send
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import com.lark.app.ui.components.*
import com.lark.app.viewmodel.ChatEntry
import com.lark.app.viewmodel.GameUiState
import kotlinx.coroutines.launch

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun GameScreen(
    uiState: GameUiState,
    onBack: () -> Unit,
    onChoiceSelected: (Int) -> Unit,
    onFreeTextSent: (String) -> Unit,
    modifier: Modifier = Modifier
) {
    var freeTextInput by remember { mutableStateOf("") }
    val listState = rememberLazyListState()
    val scope = rememberCoroutineScope()

    // Auto-scroll to bottom when new messages arrive
    LaunchedEffect(uiState.chatHistory.size) {
        if (uiState.chatHistory.isNotEmpty()) {
            scope.launch {
                listState.animateScrollToItem(uiState.chatHistory.size - 1)
            }
        }
    }

    Scaffold(
        topBar = {
            TopAppBar(
                title = {
                    Column {
                        Text(
                            uiState.scenarioName,
                            style = MaterialTheme.typography.titleMedium,
                            fontWeight = FontWeight.Bold
                        )
                        Text(
                            uiState.languageName,
                            style = MaterialTheme.typography.labelSmall,
                            color = MaterialTheme.colorScheme.onSurface.copy(alpha = 0.6f)
                        )
                    }
                },
                navigationIcon = {
                    IconButton(onClick = onBack) {
                        Icon(Icons.Filled.ArrowBack, contentDescription = "Back")
                    }
                }
            )
        },
        bottomBar = {
            if (!uiState.isFinished) {
                GameInputBar(
                    choices = uiState.currentMessage?.choices ?: emptyList(),
                    isLoading = uiState.isLoading,
                    freeTextInput = freeTextInput,
                    onFreeTextChanged = { freeTextInput = it },
                    onChoiceSelected = onChoiceSelected,
                    onFreeTextSent = {
                        onFreeTextSent(freeTextInput)
                        freeTextInput = ""
                    }
                )
            }
        }
    ) { padding ->
        LazyColumn(
            state = listState,
            modifier = modifier
                .fillMaxSize()
                .padding(padding),
            contentPadding = PaddingValues(vertical = 8.dp)
        ) {
            items(uiState.chatHistory) { entry ->
                ChatEntryView(entry)
            }

            if (uiState.isLoading) {
                item {
                    Box(
                        modifier = Modifier
                            .fillMaxWidth()
                            .padding(16.dp),
                        contentAlignment = Alignment.Center
                    ) {
                        CircularProgressIndicator(modifier = Modifier.size(24.dp))
                    }
                }
            }

            if (uiState.isFinished) {
                item {
                    Card(
                        modifier = Modifier
                            .fillMaxWidth()
                            .padding(16.dp),
                        colors = CardDefaults.cardColors(
                            containerColor = MaterialTheme.colorScheme.primaryContainer
                        )
                    ) {
                        Text(
                            text = "Scenario Complete!",
                            modifier = Modifier.padding(16.dp),
                            style = MaterialTheme.typography.headlineSmall,
                            fontWeight = FontWeight.Bold,
                            color = MaterialTheme.colorScheme.onPrimaryContainer
                        )
                    }
                }
            }

            if (uiState.error != null) {
                item {
                    Card(
                        modifier = Modifier
                            .fillMaxWidth()
                            .padding(16.dp),
                        colors = CardDefaults.cardColors(
                            containerColor = MaterialTheme.colorScheme.errorContainer
                        )
                    ) {
                        Text(
                            text = uiState.error,
                            modifier = Modifier.padding(16.dp),
                            color = MaterialTheme.colorScheme.onErrorContainer
                        )
                    }
                }
            }
        }
    }
}

@Composable
private fun ChatEntryView(entry: ChatEntry) {
    if (entry.isPlayerTurn && entry.playerInput != null) {
        PlayerBubble(text = entry.playerInput)
        return
    }

    val msg = entry.message

    if (msg.narrative.isNotBlank()) {
        NarrativeBubble(
            narrative = msg.narrative,
            translation = msg.translation
        )
    }

    if (!msg.npcDialog.isNullOrBlank()) {
        NpcDialogBubble(
            dialog = msg.npcDialog,
            translation = msg.npcDialogTranslation ?: ""
        )
    }

    if (entry.correction != null) {
        CorrectionCard(correction = entry.correction)
    }

    VocabCard(vocabulary = msg.vocabulary)
}

@Composable
private fun GameInputBar(
    choices: List<com.lark.app.data.api.Choice>,
    isLoading: Boolean,
    freeTextInput: String,
    onFreeTextChanged: (String) -> Unit,
    onChoiceSelected: (Int) -> Unit,
    onFreeTextSent: () -> Unit
) {
    Surface(tonalElevation = 3.dp) {
        Column(modifier = Modifier.padding(8.dp)) {
            // Choice buttons
            if (choices.isNotEmpty()) {
                ChoiceButtons(
                    choices = choices,
                    enabled = !isLoading,
                    onChoiceSelected = onChoiceSelected,
                    modifier = Modifier.padding(bottom = 8.dp)
                )
            }

            // Free text input
            Row(
                modifier = Modifier.fillMaxWidth(),
                verticalAlignment = Alignment.CenterVertically
            ) {
                OutlinedTextField(
                    value = freeTextInput,
                    onValueChange = onFreeTextChanged,
                    modifier = Modifier.weight(1f),
                    placeholder = { Text("Type your own response...") },
                    enabled = !isLoading,
                    singleLine = true
                )
                Spacer(modifier = Modifier.width(8.dp))
                IconButton(
                    onClick = onFreeTextSent,
                    enabled = !isLoading && freeTextInput.isNotBlank()
                ) {
                    Icon(Icons.Filled.Send, contentDescription = "Send")
                }
            }
        }
    }
}
