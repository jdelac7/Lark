package com.lark.app.game.ui

import com.badlogic.gdx.graphics.Color
import com.badlogic.gdx.graphics.Pixmap
import com.badlogic.gdx.graphics.Texture
import com.badlogic.gdx.graphics.g2d.Batch
import com.badlogic.gdx.scenes.scene2d.Actor
import com.badlogic.gdx.scenes.scene2d.InputEvent
import com.badlogic.gdx.scenes.scene2d.InputListener
import com.lark.app.game.sprite.Direction

/**
 * Retro cross-shaped D-pad with beveled pixel art look.
 * Size is set externally via setSize().
 */
class DPadActor : Actor() {

    var currentDirection: Direction? = null
        private set

    private val crossTexture: Texture
    private val highlightTex: Texture

    init {
        crossTexture = createCrossTexture()

        val hp = Pixmap(1, 1, Pixmap.Format.RGBA8888)
        hp.setColor(1f, 1f, 1f, 0.18f)
        hp.fill()
        highlightTex = Texture(hp)
        hp.dispose()

        addListener(object : InputListener() {
            override fun touchDown(event: InputEvent?, x: Float, y: Float, pointer: Int, button: Int): Boolean {
                updateDirection(x, y)
                return true
            }

            override fun touchDragged(event: InputEvent?, x: Float, y: Float, pointer: Int) {
                updateDirection(x, y)
            }

            override fun touchUp(event: InputEvent?, x: Float, y: Float, pointer: Int, button: Int) {
                currentDirection = null
            }
        })
    }

    private fun updateDirection(touchX: Float, touchY: Float) {
        val cx = width / 2
        val cy = height / 2
        val dx = touchX - cx
        val dy = touchY - cy
        val deadZone = width * 0.10f

        if (dx * dx + dy * dy < deadZone * deadZone) {
            currentDirection = null
            return
        }

        currentDirection = if (Math.abs(dx) > Math.abs(dy)) {
            if (dx > 0) Direction.RIGHT else Direction.LEFT
        } else {
            if (dy > 0) Direction.UP else Direction.DOWN
        }
    }

    override fun draw(batch: Batch, parentAlpha: Float) {
        val c = batch.color.cpy()
        batch.setColor(1f, 1f, 1f, parentAlpha)

        // Draw cross-shaped d-pad
        batch.draw(crossTexture, x, y, width, height)

        // Highlight pressed arm
        val armFrac = 32f / 96f
        val armW = width * armFrac
        val armH = height * armFrac
        val midX = x + (width - armW) / 2
        val midY = y + (height - armH) / 2

        when (currentDirection) {
            Direction.UP -> batch.draw(highlightTex, midX, y + height - height * 32f / 96f, armW, armH)
            Direction.DOWN -> batch.draw(highlightTex, midX, y, armW, armH)
            Direction.LEFT -> batch.draw(highlightTex, x, midY, armW, armH)
            Direction.RIGHT -> batch.draw(highlightTex, x + width - width * 32f / 96f, midY, armW, armH)
            null -> {}
        }

        batch.color = c
    }

    private fun createCrossTexture(): Texture {
        val s = 96
        val p = Pixmap(s, s, Pixmap.Format.RGBA8888)
        p.setColor(Color.CLEAR); p.fill()

        val armW = 32
        val half = armW / 2
        val c = s / 2

        // Colors - dark gunmetal grey with subtle warmth
        val face = Color(0.28f, 0.28f, 0.32f, 1f)
        val highlight = Color(0.42f, 0.42f, 0.48f, 1f)
        val shadow = Color(0.14f, 0.14f, 0.17f, 1f)
        val outline = Color(0.08f, 0.08f, 0.10f, 1f)
        val arrow = Color(0.20f, 0.20f, 0.24f, 1f)
        val arrowLt = Color(0.30f, 0.30f, 0.35f, 1f)
        val centerDot = Color(0.22f, 0.22f, 0.26f, 1f)

        // Draw cross shape - vertical bar
        p.setColor(face)
        p.fillRectangle(c - half, 0, armW, s)
        // Horizontal bar
        p.fillRectangle(0, c - half, s, armW)

        // === Outline (2px dark border around the cross) ===
        p.setColor(outline)
        // Vertical bar outline
        for (y2 in 0 until s) {
            if (y2 < c - half || y2 >= c + half) {
                p.drawPixel(c - half, y2); p.drawPixel(c - half + 1, y2)
                p.drawPixel(c + half - 1, y2); p.drawPixel(c + half - 2, y2)
            }
        }
        for (x2 in c - half until c + half) { p.drawPixel(x2, 0); p.drawPixel(x2, 1); p.drawPixel(x2, s - 1); p.drawPixel(x2, s - 2) }
        // Horizontal bar outline
        for (x2 in 0 until s) {
            if (x2 < c - half || x2 >= c + half) {
                p.drawPixel(x2, c - half); p.drawPixel(x2, c - half + 1)
                p.drawPixel(x2, c + half - 1); p.drawPixel(x2, c + half - 2)
            }
        }
        for (y2 in c - half until c + half) { p.drawPixel(0, y2); p.drawPixel(1, y2); p.drawPixel(s - 1, y2); p.drawPixel(s - 2, y2) }
        // Corner transitions
        for (dy in 0..1) for (dx in 0..1) {
            p.drawPixel(c - half + dx, c - half + dy); p.drawPixel(c + half - 1 - dx, c - half + dy)
            p.drawPixel(c - half + dx, c + half - 1 - dy); p.drawPixel(c + half - 1 - dx, c + half - 1 - dy)
        }

        // === Beveled highlight (top/left inner edges, 2px wide) ===
        p.setColor(highlight)
        for (x2 in c - half + 2 until c + half - 2) { p.drawPixel(x2, 2); p.drawPixel(x2, 3) }
        for (y2 in c - half + 2 until c + half - 2) { p.drawPixel(2, y2); p.drawPixel(3, y2) }
        for (x2 in 2 until c - half) { p.drawPixel(x2, c - half + 2); p.drawPixel(x2, c - half + 3) }
        for (x2 in c + half until s - 2) { p.drawPixel(x2, c - half + 2); p.drawPixel(x2, c - half + 3) }
        for (y2 in 2 until c - half) { p.drawPixel(c - half + 2, y2); p.drawPixel(c - half + 3, y2) }
        for (y2 in c + half until s - 2) { p.drawPixel(c - half + 2, y2); p.drawPixel(c - half + 3, y2) }

        // === Beveled shadow (bottom/right inner edges, 2px wide) ===
        p.setColor(shadow)
        for (x2 in c - half + 2 until c + half - 2) { p.drawPixel(x2, s - 3); p.drawPixel(x2, s - 4) }
        for (y2 in c - half + 2 until c + half - 2) { p.drawPixel(s - 3, y2); p.drawPixel(s - 4, y2) }
        for (x2 in 2 until c - half) { p.drawPixel(x2, c + half - 3); p.drawPixel(x2, c + half - 4) }
        for (x2 in c + half until s - 2) { p.drawPixel(x2, c + half - 3); p.drawPixel(x2, c + half - 4) }
        for (y2 in 2 until c - half) { p.drawPixel(c + half - 3, y2); p.drawPixel(c + half - 4, y2) }
        for (y2 in c + half until s - 2) { p.drawPixel(c + half - 3, y2); p.drawPixel(c + half - 4, y2) }

        // === Arrow indicators (doubled size, bolder) ===
        p.setColor(arrow)
        // Up arrow
        p.fillRectangle(c - 2, 10, 4, 2)
        p.fillRectangle(c - 4, 12, 8, 2)
        p.fillRectangle(c - 6, 14, 12, 2)
        p.setColor(arrowLt)
        p.fillRectangle(c - 1, 10, 2, 2)

        // Down arrow
        p.setColor(arrow)
        p.fillRectangle(c - 6, 80, 12, 2)
        p.fillRectangle(c - 4, 82, 8, 2)
        p.fillRectangle(c - 2, 84, 4, 2)
        p.setColor(arrowLt)
        p.fillRectangle(c - 1, 84, 2, 2)

        // Left arrow
        p.setColor(arrow)
        p.fillRectangle(10, c - 2, 2, 4)
        p.fillRectangle(12, c - 4, 2, 8)
        p.fillRectangle(14, c - 6, 2, 12)
        p.setColor(arrowLt)
        p.fillRectangle(10, c - 1, 2, 2)

        // Right arrow
        p.setColor(arrow)
        p.fillRectangle(84, c - 2, 2, 4)
        p.fillRectangle(82, c - 4, 2, 8)
        p.fillRectangle(80, c - 6, 2, 12)
        p.setColor(arrowLt)
        p.fillRectangle(84, c - 1, 2, 2)

        // Center circle indent
        p.setColor(centerDot)
        p.fillRectangle(c - 5, c - 5, 10, 10)
        p.setColor(shadow)
        for (i in 0..1) {
            p.drawPixel(c - 5 + i, c - 5); p.drawPixel(c + 4 - i, c - 5)
            p.drawPixel(c - 5, c - 5 + i); p.drawPixel(c - 5, c + 4 - i)
            p.drawPixel(c - 5 + i, c + 4); p.drawPixel(c + 4 - i, c + 4)
            p.drawPixel(c + 4, c - 5 + i); p.drawPixel(c + 4, c + 4 - i)
        }

        val tex = Texture(p)
        tex.setFilter(Texture.TextureFilter.Nearest, Texture.TextureFilter.Nearest)
        p.dispose()
        return tex
    }
}
