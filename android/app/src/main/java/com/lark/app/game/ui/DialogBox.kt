package com.lark.app.game.ui

import com.badlogic.gdx.graphics.Color
import com.badlogic.gdx.graphics.Pixmap
import com.badlogic.gdx.graphics.Texture
import com.badlogic.gdx.graphics.g2d.Batch
import com.badlogic.gdx.graphics.g2d.GlyphLayout
import com.badlogic.gdx.scenes.scene2d.Actor
import com.lark.app.game.util.PixelFont

/**
 * RPG-style dialog box: dark background with white pixel border,
 * typewriter text rendering.
 */
class DialogBox(private val pixelFont: PixelFont) : Actor() {

    private val bgTexture: Texture
    private val borderTexture: Texture
    private var displayText: String = ""
    private val layout = GlyphLayout()

    companion object {
        const val PADDING = 24f
        const val BORDER = 5f
    }

    init {
        // Dark semi-transparent background
        val bgPix = Pixmap(1, 1, Pixmap.Format.RGBA8888)
        bgPix.setColor(0.05f, 0.05f, 0.15f, 0.92f)
        bgPix.fill()
        bgTexture = Texture(bgPix)
        bgPix.dispose()

        // White border
        val borderPix = Pixmap(1, 1, Pixmap.Format.RGBA8888)
        borderPix.setColor(0.9f, 0.9f, 0.95f, 1f)
        borderPix.fill()
        borderTexture = Texture(borderPix)
        borderPix.dispose()
    }

    fun setText(text: String) {
        displayText = text
    }

    override fun draw(batch: Batch, parentAlpha: Float) {
        val c = batch.color.cpy()
        batch.setColor(1f, 1f, 1f, parentAlpha)

        // Border
        batch.draw(borderTexture, x, y, width, height)

        // Background (inset by border width)
        batch.draw(
            bgTexture,
            x + BORDER, y + BORDER,
            width - BORDER * 2, height - BORDER * 2
        )

        // Text
        if (displayText.isNotEmpty()) {
            val font = pixelFont.font
            font.color = Color.WHITE
            val textX = x + PADDING
            val textY = y + height - PADDING
            val textW = width - PADDING * 2

            layout.setText(font, displayText, Color.WHITE, textW, com.badlogic.gdx.utils.Align.left, true)
            font.draw(batch, layout, textX, textY)
        }

        batch.color = c
    }
}
