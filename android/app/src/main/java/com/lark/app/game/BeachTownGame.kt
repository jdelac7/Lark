package com.lark.app.game

import com.badlogic.gdx.Game
import com.badlogic.gdx.graphics.g2d.SpriteBatch
import com.lark.app.game.screens.WorldGameScreen

class BeachTownGame(
    val languageCode: String,
    val languageName: String,
    val playerId: String,
    val onExit: () -> Unit
) : Game() {

    lateinit var batch: SpriteBatch
        private set

    override fun create() {
        batch = SpriteBatch()
        setScreen(WorldGameScreen(this))
    }

    override fun dispose() {
        batch.dispose()
        screen?.dispose()
    }
}
