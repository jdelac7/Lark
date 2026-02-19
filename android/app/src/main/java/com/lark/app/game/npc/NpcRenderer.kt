package com.lark.app.game.npc

import com.badlogic.gdx.graphics.Color
import com.badlogic.gdx.graphics.Pixmap
import com.badlogic.gdx.graphics.Texture
import com.badlogic.gdx.graphics.g2d.SpriteBatch
import com.badlogic.gdx.utils.Disposable
import com.lark.app.game.sprite.NpcSprites

class NpcRenderer(private val npcSprites: NpcSprites) : Disposable {

    private val exclamationTexture: Texture
    private val checkTexture: Texture
    private var indicatorBob = 0f

    init {
        exclamationTexture = createExclamationTexture()
        checkTexture = createCheckTexture()
    }

    /** Render a single NPC (used during Y-sorted pass). */
    fun renderSingle(batch: SpriteBatch, npc: Npc, facingNpc: Npc?) {
        indicatorBob += 0.02f

        val tex = npcSprites.textures.getOrNull(npc.spriteId) ?: return
        val x = npc.tileX * 16f
        val y = npc.tileY * 16f

        // Draw NPC at 16x24 (extends 8px above tile)
        batch.draw(tex, x, y, 16f, 24f)

        // Show "!" if player is facing this NPC
        if (npc == facingNpc && !npc.completed) {
            val bobOffset = kotlin.math.sin(indicatorBob.toDouble()).toFloat() * 2f
            batch.draw(exclamationTexture, x + 3f, y + 26f + bobOffset, 10f, 10f)
        }

        // Show checkmark if completed
        if (npc.completed) {
            batch.draw(checkTexture, x + 3f, y + 26f, 10f, 10f)
        }
    }

    fun render(batch: SpriteBatch, npcs: List<Npc>, facingNpc: Npc?) {
        for (npc in npcs) {
            renderSingle(batch, npc, facingNpc)
        }
    }

    private fun createExclamationTexture(): Texture {
        val p = Pixmap(16, 16, Pixmap.Format.RGBA8888)
        p.setColor(Color.CLEAR)
        p.fill()
        val yellow = Color(1f, 0.9f, 0.1f, 1f)
        val outline = Color(0.7f, 0.6f, 0f, 1f)
        val bright = Color(1f, 0.95f, 0.4f, 1f)
        // Exclamation body (2x scale)
        p.setColor(yellow)
        p.fillRectangle(6, 0, 4, 10)
        p.fillRectangle(6, 12, 4, 4)
        // Bright center highlight
        p.setColor(bright)
        p.fillRectangle(7, 1, 2, 8)
        p.fillRectangle(7, 13, 2, 2)
        // Outline accents
        p.setColor(outline)
        p.drawPixel(5, 0); p.drawPixel(5, 1); p.drawPixel(10, 0); p.drawPixel(10, 1)
        p.drawPixel(5, 12); p.drawPixel(5, 13); p.drawPixel(10, 14); p.drawPixel(10, 15)
        val tex = Texture(p)
        tex.setFilter(Texture.TextureFilter.Nearest, Texture.TextureFilter.Nearest)
        p.dispose()
        return tex
    }

    private fun createCheckTexture(): Texture {
        val p = Pixmap(16, 16, Pixmap.Format.RGBA8888)
        p.setColor(Color.CLEAR)
        p.fill()
        val green = Color(0.2f, 0.85f, 0.2f, 1f)
        val dark = Color(0.1f, 0.6f, 0.1f, 1f)
        val bright = Color(0.4f, 0.95f, 0.4f, 1f)
        // Checkmark (2x scale, thicker strokes)
        p.setColor(green)
        p.fillRectangle(2, 8, 2, 2); p.fillRectangle(4, 10, 2, 2); p.fillRectangle(6, 12, 2, 2)
        p.fillRectangle(8, 10, 2, 2); p.fillRectangle(10, 8, 2, 2); p.fillRectangle(12, 6, 2, 2); p.fillRectangle(14, 4, 2, 2)
        // Shadow line
        p.setColor(dark)
        p.fillRectangle(4, 8, 2, 2); p.fillRectangle(6, 10, 2, 2); p.fillRectangle(8, 8, 2, 2)
        p.fillRectangle(10, 6, 2, 2); p.fillRectangle(12, 4, 2, 2)
        // Bright highlight
        p.setColor(bright)
        p.drawPixel(3, 8); p.drawPixel(5, 10); p.drawPixel(7, 12)
        p.drawPixel(9, 10); p.drawPixel(11, 8); p.drawPixel(13, 6); p.drawPixel(15, 4)
        val tex = Texture(p)
        tex.setFilter(Texture.TextureFilter.Nearest, Texture.TextureFilter.Nearest)
        p.dispose()
        return tex
    }

    override fun dispose() {
        exclamationTexture.dispose()
        checkTexture.dispose()
    }
}
