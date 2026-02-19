package com.lark.app

import android.content.Intent
import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.compose.runtime.*
import androidx.lifecycle.viewmodel.compose.viewModel
import com.lark.app.data.api.RetrofitClient
import com.lark.app.ui.screens.HomeScreen
import com.lark.app.ui.screens.ProgressScreen
import com.lark.app.ui.theme.LarkTheme
import com.lark.app.viewmodel.HomeViewModel
import com.lark.app.viewmodel.ProgressViewModel
import java.util.UUID

class MainActivity : ComponentActivity() {

    private val playerId by lazy {
        val prefs = getSharedPreferences("lark", MODE_PRIVATE)
        prefs.getString("playerId", null) ?: UUID.randomUUID().toString().also {
            prefs.edit().putString("playerId", it).apply()
        }
    }

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)

        RetrofitClient.setPlayerId(playerId)

        setContent {
            LarkTheme {
                LarkApp(
                    playerId = playerId,
                    onEnterTown = { languageCode, languageName ->
                        val intent = Intent(this, WorldActivity::class.java).apply {
                            putExtra("languageCode", languageCode)
                            putExtra("languageName", languageName)
                            putExtra("playerId", playerId)
                        }
                        startActivity(intent)
                    }
                )
            }
        }
    }
}

sealed class Screen {
    data object Home : Screen()
    data object Progress : Screen()
}

@Composable
fun LarkApp(
    playerId: String,
    onEnterTown: (languageCode: String, languageName: String) -> Unit
) {
    var currentScreen by remember { mutableStateOf<Screen>(Screen.Home) }

    when (currentScreen) {
        is Screen.Home -> {
            val homeVm: HomeViewModel = viewModel()
            val uiState by homeVm.uiState.collectAsState()

            HomeScreen(
                uiState = uiState,
                onLanguageSelected = { homeVm.selectLanguage(it) },
                onEnterTown = {
                    val lang = uiState.selectedLanguage ?: return@HomeScreen
                    onEnterTown(lang.code, lang.name)
                },
                onProgressClicked = { currentScreen = Screen.Progress }
            )
        }

        is Screen.Progress -> {
            val progressVm: ProgressViewModel = viewModel()
            val uiState by progressVm.uiState.collectAsState()

            LaunchedEffect(Unit) {
                progressVm.loadProgress(playerId)
            }

            ProgressScreen(
                uiState = uiState,
                onBack = { currentScreen = Screen.Home }
            )
        }
    }
}
