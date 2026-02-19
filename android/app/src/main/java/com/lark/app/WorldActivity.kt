package com.lark.app

import android.os.Bundle
import com.badlogic.gdx.backends.android.AndroidApplication
import com.badlogic.gdx.backends.android.AndroidApplicationConfiguration
import com.lark.app.game.BeachTownGame

class WorldActivity : AndroidApplication() {

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)

        val languageCode = intent.getStringExtra("languageCode") ?: "es"
        val languageName = intent.getStringExtra("languageName") ?: "Spanish"
        val playerId = intent.getStringExtra("playerId") ?: ""

        val config = AndroidApplicationConfiguration().apply {
            useAccelerometer = false
            useCompass = false
            useImmersiveMode = true
        }

        val game = BeachTownGame(languageCode, languageName, playerId) {
            finish()
        }

        initialize(game, config)
    }
}
