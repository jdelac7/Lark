package com.lark.app.game.ui

import com.badlogic.gdx.graphics.Color
import com.badlogic.gdx.graphics.Pixmap
import com.badlogic.gdx.graphics.Texture
import com.badlogic.gdx.graphics.g2d.Batch
import com.badlogic.gdx.graphics.g2d.GlyphLayout
import com.badlogic.gdx.scenes.scene2d.Actor
import com.lark.app.data.api.Correction
import com.lark.app.game.util.PixelFont

/**
 * Shows original vs corrected text with explanation.
 */
class CorrectionPopup(
    private val pixelFont: PixelFont,
    private val correction: Correction
) : Actor() {

    private val bgTexture: Texture
    private val borderTexture: Texture
    private val layout = GlyphLayout()

    companion object {
        const val PADDING = 24f
        const val LINE_HEIGHT = 50f
        const val BORDER = 5f
    }

    init {
        val bgPix = Pixmap(1, 1, Pixmap.Format.RGBA8888)
        bgPix.setColor(0.15f, 0.05f, 0.05f, 0.94f)
        bgPix.fill()
        bgTexture = Texture(bgPix)
        bgPix.dispose()

        val borderPix = Pixmap(1, 1, Pixmap.Format.RGBA8888)
        borderPix.setColor(0.9f, 0.4f, 0.3f, 1f)
        borderPix.fill()
        borderTexture = Texture(borderPix)
        borderPix.dispose()
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
        val textX = x + BORDER + PADDING
        var drawY = y + height - BORDER - PADDING

        // Title
        font.color = Color(1f, 0.4f, 0.3f, 1f)
        font.draw(batch, "Grammar Note", textX, drawY)
        drawY -= LINE_HEIGHT * 1.3f

        // Original (with strikethrough effect - draw in red)
        font.color = Color(0.9f, 0.3f, 0.3f, 1f)
        font.draw(batch, correction.original, textX, drawY)
        drawY -= LINE_HEIGHT

        // Corrected
        font.color = Color(0.3f, 1f, 0.3f, 1f)
        font.draw(batch, correction.corrected, textX, drawY)
        drawY -= LINE_HEIGHT * 1.2f

        // Explanation
        val wrapW = width - PADDING * 2 - BORDER * 2
        fontSmall.color = Color(0.85f, 0.85f, 0.85f, 1f)
        layout.setText(fontSmall, correction.explanation, fontSmall.color, wrapW,
            com.badlogic.gdx.utils.Align.left, true)
        fontSmall.draw(batch, layout, textX, drawY)

        // "Press A to continue"
        fontSmall.color = Color.GRAY
        fontSmall.draw(batch, "Press A to continue", textX, y + BORDER + PADDING + 10)

        batch.color = c
    }
}
