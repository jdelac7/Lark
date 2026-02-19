package com.lark.app.game.ui

import com.badlogic.gdx.graphics.Color
import com.badlogic.gdx.graphics.Pixmap
import com.badlogic.gdx.graphics.Texture
import com.badlogic.gdx.graphics.g2d.Batch
import com.badlogic.gdx.graphics.g2d.GlyphLayout
import com.badlogic.gdx.scenes.scene2d.Actor
import com.badlogic.gdx.scenes.scene2d.InputEvent
import com.badlogic.gdx.scenes.scene2d.InputListener
import com.lark.app.game.util.PixelFont

/**
 * Retro circular action button (SNES/GBA style).
 * Size is set externally via setSize().
 */
class ActionButtonActor(
    private val pixelFont: PixelFont,
    private val onClick: () -> Unit
) : Actor() {

    var label: String = "A"
    private val buttonTexture: Texture
    private val pressedTexture: Texture
    private var pressed = false
    private val layout = GlyphLayout()

    init {
        buttonTexture = createButtonTexture(false)
        pressedTexture = createButtonTexture(true)

        addListener(object : InputListener() {
            override fun touchDown(event: InputEvent?, x: Float, y: Float, pointer: Int, button: Int): Boolean {
                pressed = true
                return true
            }

            override fun touchUp(event: InputEvent?, x: Float, y: Float, pointer: Int, button: Int) {
                if (pressed) {
                    pressed = false
                    onClick()
                }
            }
        })
    }

    override fun draw(batch: Batch, parentAlpha: Float) {
        val c = batch.color.cpy()
        batch.setColor(1f, 1f, 1f, parentAlpha)

        val tex = if (pressed) pressedTexture else buttonTexture
        batch.draw(tex, x, y, width, height)

        // Draw label centered
        val font = pixelFont.font
        val prevColor = font.color.cpy()
        font.color = Color.WHITE
        layout.setText(font, label)
        font.draw(batch, label, x + (width - layout.width) / 2, y + (height + layout.height) / 2)
        font.color = prevColor

        batch.color = c
    }

    private fun createButtonTexture(isPressed: Boolean): Texture {
        val s = 64
        val p = Pixmap(s, s, Pixmap.Format.RGBA8888)
        p.setColor(Color.CLEAR); p.fill()

        val cx = s / 2
        val cy = s / 2
        val r = 28

        if (isPressed) {
            // Pressed: darker, flatter look (pushed in)
            val face = Color(0.62f, 0.12f, 0.14f, 1f)
            val rim = Color(0.45f, 0.08f, 0.10f, 1f)
            val inner = Color(0.68f, 0.16f, 0.18f, 1f)

            // Outer rim
            p.setColor(rim); p.fillCircle(cx, cy, r)
            // Face
            p.setColor(face); p.fillCircle(cx, cy, r - 4)
            // Subtle inner highlight (smaller, pressed look)
            p.setColor(inner)
            p.fillCircle(cx, cy + 2, r - 10)
        } else {
            // Normal: raised, beveled look
            val face = Color(0.78f, 0.18f, 0.20f, 1f)
            val highlight = Color(0.92f, 0.32f, 0.34f, 1f)
            val shadow = Color(0.52f, 0.10f, 0.12f, 1f)
            val rim = Color(0.40f, 0.06f, 0.08f, 1f)
            val shine = Color(0.98f, 0.42f, 0.44f, 1f)

            // Dark rim (shadow at bottom-right)
            p.setColor(rim); p.fillCircle(cx, cy, r)
            // Shadow crescent (bottom)
            p.setColor(shadow); p.fillCircle(cx, cy + 2, r - 2)
            // Main face
            p.setColor(face); p.fillCircle(cx, cy, r - 4)
            // Highlight crescent (top)
            p.setColor(highlight); p.fillCircle(cx, cy - 4, r - 8)
            // Face center
            p.setColor(face); p.fillCircle(cx, cy - 2, r - 10)
            // Specular shine (2x sized)
            p.setColor(shine)
            p.fillRectangle(cx - 8, cy - 13, 2, 2)
            p.fillRectangle(cx - 6, cy - 14, 2, 2)
            p.fillRectangle(cx - 6, cy - 12, 2, 2)
            p.fillRectangle(cx - 4, cy - 14, 2, 2)
        }

        val tex = Texture(p)
        tex.setFilter(Texture.TextureFilter.Nearest, Texture.TextureFilter.Nearest)
        p.dispose()
        return tex
    }
}
