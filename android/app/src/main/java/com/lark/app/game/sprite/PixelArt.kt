package com.lark.app.game.sprite

import com.badlogic.gdx.graphics.Color
import com.badlogic.gdx.graphics.Pixmap
import com.badlogic.gdx.graphics.Texture

/**
 * Utility for creating textures from pixel grid arrays.
 * Each grid is a 2D array of palette indices, where -1 means transparent.
 */
object PixelArt {

    fun createTexture(width: Int, height: Int, pixels: IntArray, palette: Array<Color>): Texture {
        val pixmap = Pixmap(width, height, Pixmap.Format.RGBA8888)
        pixmap.setColor(Color.CLEAR)
        pixmap.fill()

        for (y in 0 until height) {
            for (x in 0 until width) {
                val idx = pixels[y * width + x]
                if (idx >= 0 && idx < palette.size) {
                    pixmap.setColor(palette[idx])
                    pixmap.drawPixel(x, y)
                }
            }
        }

        val texture = Texture(pixmap)
        texture.setFilter(Texture.TextureFilter.Nearest, Texture.TextureFilter.Nearest)
        pixmap.dispose()
        return texture
    }

    fun createSolidTexture(width: Int, height: Int, color: Color): Texture {
        val pixmap = Pixmap(width, height, Pixmap.Format.RGBA8888)
        pixmap.setColor(color)
        pixmap.fill()
        val texture = Texture(pixmap)
        texture.setFilter(Texture.TextureFilter.Nearest, Texture.TextureFilter.Nearest)
        pixmap.dispose()
        return texture
    }
}
