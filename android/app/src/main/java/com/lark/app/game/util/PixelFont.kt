package com.lark.app.game.util

import com.badlogic.gdx.graphics.Color
import com.badlogic.gdx.graphics.Pixmap
import com.badlogic.gdx.graphics.Texture
import com.badlogic.gdx.graphics.g2d.BitmapFont
import com.badlogic.gdx.graphics.g2d.BitmapFont.BitmapFontData
import com.badlogic.gdx.graphics.g2d.TextureRegion
import com.badlogic.gdx.utils.Disposable

/**
 * Creates a simple pixel-art style BitmapFont using libGDX's built-in default font
 * with nearest-neighbor filtering for a crisp retro look.
 */
class PixelFont : Disposable {

    val font: BitmapFont
    val fontSmall: BitmapFont

    init {
        // Use libGDX's built-in Arial 15 bitmap font with nearest filtering
        font = BitmapFont().apply {
            region.texture.setFilter(
                Texture.TextureFilter.Nearest,
                Texture.TextureFilter.Nearest
            )
            data.setScale(2.5f)
            color = Color.WHITE
        }

        fontSmall = BitmapFont().apply {
            region.texture.setFilter(
                Texture.TextureFilter.Nearest,
                Texture.TextureFilter.Nearest
            )
            data.setScale(2.0f)
            color = Color.WHITE
        }
    }

    override fun dispose() {
        font.dispose()
        fontSmall.dispose()
    }
}
