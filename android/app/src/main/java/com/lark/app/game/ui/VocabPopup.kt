package com.lark.app.game.ui

import com.badlogic.gdx.graphics.Color
import com.badlogic.gdx.graphics.Pixmap
import com.badlogic.gdx.graphics.Texture
import com.badlogic.gdx.graphics.g2d.Batch
import com.badlogic.gdx.graphics.g2d.GlyphLayout
import com.badlogic.gdx.scenes.scene2d.Actor
import com.lark.app.data.api.VocabItem
import com.lark.app.game.util.PixelFont

/**
 * Vocabulary word list popup styled as pixel-art box.
 */
class VocabPopup(
    private val pixelFont: PixelFont,
    private val vocabulary: List<VocabItem>
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
        bgPix.setColor(0.05f, 0.1f, 0.05f, 0.94f)
        bgPix.fill()
        bgTexture = Texture(bgPix)
        bgPix.dispose()

        val borderPix = Pixmap(1, 1, Pixmap.Format.RGBA8888)
        borderPix.setColor(0.3f, 0.9f, 0.3f, 1f)
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
        var drawY = y + height - BORDER - PADDING

        // Title
        font.color = Color(0.3f, 1f, 0.3f, 1f)
        font.draw(batch, "Vocabulary", x + BORDER + PADDING, drawY)
        drawY -= LINE_HEIGHT * 1.2f

        // Words
        for (vocab in vocabulary) {
            font.color = Color.WHITE
            font.draw(batch, vocab.word, x + BORDER + PADDING, drawY)

            font.color = Color.LIGHT_GRAY
            font.draw(batch, vocab.translation, x + width / 2, drawY)

            drawY -= LINE_HEIGHT

            if (vocab.usage != null) {
                fontSmall.color = Color(0.7f, 0.7f, 0.7f, 1f)
                fontSmall.draw(batch, vocab.usage, x + BORDER + PADDING + 10, drawY)
                drawY -= LINE_HEIGHT * 0.8f
            }
        }

        // "Press A to continue"
        fontSmall.color = Color.GRAY
        fontSmall.draw(batch, "Press A to continue", x + BORDER + PADDING, y + BORDER + PADDING + 10)

        batch.color = c
    }
}
