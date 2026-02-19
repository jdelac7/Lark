package com.lark.app.game.ui

import com.badlogic.gdx.graphics.Color
import com.badlogic.gdx.graphics.Pixmap
import com.badlogic.gdx.graphics.Texture
import com.badlogic.gdx.graphics.g2d.Batch
import com.badlogic.gdx.graphics.g2d.GlyphLayout
import com.badlogic.gdx.scenes.scene2d.Actor
import com.badlogic.gdx.scenes.scene2d.InputEvent
import com.badlogic.gdx.scenes.scene2d.InputListener
import com.lark.app.data.api.Choice
import com.lark.app.game.util.PixelFont

/**
 * RPG-style choice menu with arrow cursor indicator.
 * Supports D-pad navigation and touch selection.
 */
class ChoiceMenu(
    private val pixelFont: PixelFont,
    private val choices: List<Choice>,
    var selectedIndex: Int = 0,
    private val onSelect: (Int) -> Unit
) : Actor() {

    private val bgTexture: Texture
    private val borderTexture: Texture
    private val highlightTexture: Texture
    private val layout = GlyphLayout()

    companion object {
        const val PADDING = 20f
        const val ITEM_HEIGHT = 60f
        const val BORDER = 5f
        const val ARROW_WIDTH = 24f
    }

    // Total items = choices + "Type your own..."
    private val totalItems get() = choices.size + 1

    init {
        val bgPix = Pixmap(1, 1, Pixmap.Format.RGBA8888)
        bgPix.setColor(0.05f, 0.05f, 0.15f, 0.94f)
        bgPix.fill()
        bgTexture = Texture(bgPix)
        bgPix.dispose()

        val borderPix = Pixmap(1, 1, Pixmap.Format.RGBA8888)
        borderPix.setColor(0.9f, 0.9f, 0.95f, 1f)
        borderPix.fill()
        borderTexture = Texture(borderPix)
        borderPix.dispose()

        val hlPix = Pixmap(1, 1, Pixmap.Format.RGBA8888)
        hlPix.setColor(0.2f, 0.3f, 0.6f, 0.5f)
        hlPix.fill()
        highlightTexture = Texture(hlPix)
        hlPix.dispose()

        addListener(object : InputListener() {
            override fun touchDown(event: InputEvent?, x: Float, y: Float, pointer: Int, button: Int): Boolean {
                // Determine which item was tapped
                val localY = height - BORDER - PADDING - y
                val idx = (localY / ITEM_HEIGHT).toInt()
                if (idx in 0 until totalItems) {
                    selectedIndex = idx
                    onSelect(idx)
                }
                return true
            }
        })
    }

    override fun draw(batch: Batch, parentAlpha: Float) {
        val c = batch.color.cpy()
        batch.setColor(1f, 1f, 1f, parentAlpha)

        // Border
        batch.draw(borderTexture, x, y, width, height)
        // Background
        batch.draw(bgTexture, x + BORDER, y + BORDER, width - BORDER * 2, height - BORDER * 2)

        val font = pixelFont.font
        val fontSmall = pixelFont.fontSmall
        val startY = y + height - BORDER - PADDING

        for (i in 0 until totalItems) {
            val itemY = startY - i * ITEM_HEIGHT
            val itemX = x + BORDER + PADDING

            // Highlight selected
            if (i == selectedIndex) {
                batch.draw(
                    highlightTexture,
                    x + BORDER, itemY - ITEM_HEIGHT + 5,
                    width - BORDER * 2, ITEM_HEIGHT
                )
            }

            // Arrow cursor
            if (i == selectedIndex) {
                font.color = Color.YELLOW
                font.draw(batch, ">", itemX, itemY)
            }

            val textX = itemX + ARROW_WIDTH

            if (i < choices.size) {
                // Choice text
                font.color = Color.WHITE
                layout.setText(font, choices[i].text)
                font.draw(batch, choices[i].text, textX, itemY)

                // Translation (smaller, below)
                fontSmall.color = Color.LIGHT_GRAY
                fontSmall.draw(batch, choices[i].translation, textX, itemY - layout.height - 2)
            } else {
                // "Type your own..."
                font.color = Color(0.6f, 0.8f, 1f, 1f)
                font.draw(batch, "Type your own...", textX, itemY)
            }
        }

        batch.color = c
    }
}
