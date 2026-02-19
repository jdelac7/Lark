package com.lark.app.game.map

import com.badlogic.gdx.graphics.Color
import com.badlogic.gdx.graphics.Pixmap
import com.badlogic.gdx.graphics.Texture
import com.badlogic.gdx.graphics.g2d.SpriteBatch
import com.badlogic.gdx.utils.Disposable

/**
 * Renders buildings as composite multi-tile sprites with isometric side faces
 * for 3D depth. Each building is one wide texture instead of per-tile
 * wall/door/roof rendering.
 *
 * Texture layout (80px tall, top-to-bottom in Pixmap = top-to-bottom on screen):
 *   Rows 0-17:  Terracotta roof (18px)
 *   Rows 18-19: Eave shadow (2px)
 *   Rows 20-41: Upper floor — windows + flower boxes (22px)
 *   Rows 42-43: Floor divider (2px)
 *   Rows 44-67: Ground floor — door + windows + awning (24px)
 *   Rows 68-73: Stone foundation (6px)
 *   Rows 74-79: Ground shadow (6px)
 *
 * Width = building_tiles * 16 + SIDE_W (6px isometric side face on right).
 */
class BuildingRenderer : Disposable {

    private data class BuildingDef(
        val name: String,
        val startX: Int,
        val wallY: Int,
        val width: Int,
        val doorLocalX: Int,
        val accent: Color,
        val accentDk: Color
    )

    companion object {
        private const val SIDE_W = 6
        private const val HEIGHT = 80
        private const val SHADOW_H = 6
    }

    private val buildings = listOf(
        BuildingDef("Restaurant", 3, 16, 6, 2,
            Color(0.86f, 0.18f, 0.12f, 1f), Color(0.66f, 0.12f, 0.08f, 1f)),
        BuildingDef("Cafe", 12, 16, 5, 2,
            Color(0.18f, 0.62f, 0.22f, 1f), Color(0.12f, 0.44f, 0.14f, 1f)),
        BuildingDef("Hotel", 22, 16, 5, 2,
            Color(0.16f, 0.38f, 0.78f, 1f), Color(0.10f, 0.26f, 0.58f, 1f)),
        BuildingDef("Market", 28, 16, 6, 2,
            Color(0.92f, 0.58f, 0.12f, 1f), Color(0.72f, 0.42f, 0.08f, 1f)),
        BuildingDef("Doctor", 4, 6, 6, 2,
            Color(0.94f, 0.94f, 0.96f, 1f), Color(0.80f, 0.80f, 0.84f, 1f)),
        BuildingDef("TrainStation", 30, 6, 6, 3,
            Color(0.52f, 0.16f, 0.12f, 1f), Color(0.38f, 0.10f, 0.08f, 1f))
    )

    private val textures = HashMap<String, Texture>()

    init {
        for (b in buildings) textures[b.name] = createBuildingTexture(b)
    }

    /** Render all buildings whose wall row == y. Call per row in the y-sort loop. */
    fun renderRow(batch: SpriteBatch, y: Int) {
        for (b in buildings) {
            if (b.wallY == y) {
                val tex = textures[b.name] ?: continue
                batch.draw(
                    tex,
                    b.startX * 16f,
                    y * 16f - SHADOW_H,
                    tex.width.toFloat(),
                    tex.height.toFloat()
                )
            }
        }
    }

    // ═════════════════════════════════════════════════════════════════════════
    //  TEXTURE GENERATION
    // ═════════════════════════════════════════════════════════════════════════

    private fun createBuildingTexture(b: BuildingDef): Texture {
        val mainW = b.width * 16
        val totalW = mainW + SIDE_W
        val p = Pixmap(totalW, HEIGHT, Pixmap.Format.RGBA8888)
        p.setColor(Color.CLEAR); p.fill()

        drawRoof(p, mainW)
        drawEaveShadow(p, mainW)
        drawUpperFloor(p, mainW, b.width)
        drawDivider(p, mainW)
        drawGroundFloor(p, mainW, b)
        drawFoundation(p, mainW, b.doorLocalX)
        drawGroundShadow(p, totalW)
        drawSideFace(p, mainW)
        drawFrontEdgeShadows(p, mainW)

        val tex = Texture(p)
        tex.setFilter(Texture.TextureFilter.Nearest, Texture.TextureFilter.Nearest)
        p.dispose()
        return tex
    }

    // ── Roof (rows 0-17) ────────────────────────────────────────────────────

    private fun drawRoof(p: Pixmap, mw: Int) {
        val bright = Color(0.96f, 0.54f, 0.28f, 1f)
        val main   = Color(0.88f, 0.40f, 0.18f, 1f)
        val mid    = Color(0.76f, 0.30f, 0.12f, 1f)
        val dk     = Color(0.60f, 0.22f, 0.08f, 1f)
        val eave   = Color(0.50f, 0.16f, 0.04f, 1f)

        p.fill(0, 0, mw, 2, bright)    // ridge cap
        p.fill(0, 2, mw, 2, main)      // ridge body
        p.fill(0, 4, mw, 6, main)      // upper slope
        p.fill(0, 10, mw, 5, mid)      // lower slope
        p.fill(0, 15, mw, 3, eave)     // eave lip

        // Horizontal tile row lines
        p.setColor(dk)
        for (x in 0 until mw) { p.drawPixel(x, 6); p.drawPixel(x, 9); p.drawPixel(x, 12) }

        // Staggered vertical tile separators
        for (x in 0 until mw step 4) {
            for (y in 4..5)  safePx(p, x, y)
            for (y in 7..8)  safePx(p, x + 2, y)
            for (y in 10..11) safePx(p, x, y)
            for (y in 13..14) safePx(p, x + 2, y)
        }

        // Tile surface highlights
        p.setColor(bright)
        for (x in 1 until mw step 8) {
            safePx(p, x, 5); safePx(p, x + 4, 8); safePx(p, x, 11); safePx(p, x + 4, 14)
        }
    }

    // ── Eave shadow (rows 18-19) ────────────────────────────────────────────

    private fun drawEaveShadow(p: Pixmap, mw: Int) {
        p.fill(0, 18, mw, 2, Color(0f, 0f, 0f, 0.25f))
    }

    // ── Upper floor (rows 20-41) ────────────────────────────────────────────

    private fun drawUpperFloor(p: Pixmap, mw: Int, tileCount: Int) {
        val white = Color(0.98f, 0.96f, 0.92f, 1f)
        val cool  = Color(0.95f, 0.92f, 0.88f, 1f)

        p.fill(0, 20, mw, 22, white)

        // Subtle plaster variation every other tile
        for (t in 0 until tileCount) {
            if (t % 2 == 1) p.fill(t * 16 + 1, 22, 14, 16, cool)
        }

        // Window per tile
        for (t in 0 until tileCount) drawWindow(p, t * 16 + 2, 23, 12, 13)

        // Flower box under each window
        for (t in 0 until tileCount) drawFlowerBox(p, t * 16 + 3, 37, 10)
    }

    // ── Floor divider (rows 42-43) ──────────────────────────────────────────

    private fun drawDivider(p: Pixmap, mw: Int) {
        p.fill(0, 42, mw, 2, Color(0.84f, 0.80f, 0.74f, 1f))
    }

    // ── Ground floor (rows 44-67) ───────────────────────────────────────────

    private fun drawGroundFloor(p: Pixmap, mw: Int, b: BuildingDef) {
        val white = Color(0.98f, 0.96f, 0.92f, 1f)
        val cool  = Color(0.95f, 0.92f, 0.88f, 1f)

        p.fill(0, 44, mw, 24, white)

        for (t in 0 until b.width) {
            if (t % 2 == 0) p.fill(t * 16 + 1, 46, 14, 18, cool)
        }

        // Door tile — awning + door
        val dx = b.doorLocalX * 16
        drawAwning(p, dx + 1, 44, 14, 6, b.accent, b.accentDk)
        drawDoor(p, dx + 4, 52, 8, 15)

        // Windows on non-door tiles
        for (t in 0 until b.width) {
            if (t == b.doorLocalX) continue
            drawWindow(p, t * 16 + 2, 47, 12, 13)
        }
    }

    // ── Foundation (rows 68-73) ─────────────────────────────────────────────

    private fun drawFoundation(p: Pixmap, mw: Int, doorLocalX: Int) {
        val stone   = Color(0.58f, 0.50f, 0.38f, 1f)
        val stoneDk = Color(0.44f, 0.38f, 0.28f, 1f)
        val detail  = Color(0.52f, 0.44f, 0.32f, 1f)
        val step    = Color(0.76f, 0.70f, 0.60f, 1f)
        val stepLt  = Color(0.84f, 0.78f, 0.68f, 1f)

        p.fill(0, 68, mw, 3, stone)
        p.fill(0, 71, mw, 3, stoneDk)

        // Stone block texture
        p.setColor(detail)
        for (x in 0 until mw step 5) { safePx(p, x + 1, 69); safePx(p, x + 3, 72) }

        // Door step
        val dx = doorLocalX * 16
        p.fill(dx + 2, 68, 12, 3, step)
        p.setColor(stepLt)
        for (x in dx + 2 until dx + 14) safePx(p, x, 68)
    }

    // ── Ground shadow (rows 74-79) ──────────────────────────────────────────

    private fun drawGroundShadow(p: Pixmap, tw: Int) {
        p.fill(0, 74, tw, 3, Color(0f, 0f, 0f, 0.20f))
        p.fill(0, 77, tw, 3, Color(0f, 0f, 0f, 0.08f))
    }

    // ── Isometric side face (right edge, SIDE_W pixels) ─────────────────────

    private fun drawSideFace(p: Pixmap, mw: Int) {
        val sx = mw // side face starts at pixel mw

        val sRoof    = Color(0.72f, 0.28f, 0.10f, 1f)
        val sRoofDk  = Color(0.56f, 0.20f, 0.06f, 1f)
        val sEave    = Color(0.40f, 0.14f, 0.04f, 1f)
        val sWall    = Color(0.88f, 0.84f, 0.78f, 1f)
        val sWallDk  = Color(0.78f, 0.74f, 0.68f, 1f)
        val sLine    = Color(0.82f, 0.78f, 0.72f, 1f)
        val sFound   = Color(0.48f, 0.40f, 0.30f, 1f)
        val sFoundDk = Color(0.36f, 0.30f, 0.22f, 1f)

        // Side roof
        p.fill(sx, 0, SIDE_W, 15, sRoof)
        p.setColor(sRoofDk)
        for (y in 0..14) safePx(p, sx + SIDE_W - 1, y)
        // Tile row lines on side roof
        for (y in 0..14 step 3) {
            p.setColor(sRoofDk)
            for (x in sx until sx + SIDE_W) safePx(p, x, y)
        }
        // Side eave
        p.fill(sx, 15, SIDE_W, 3, sEave)

        // Side wall (rows 18-67)
        p.fill(sx, 18, SIDE_W, 50, sWall)

        // Depth gradient — darker toward right (further from viewer)
        p.setColor(sWallDk)
        for (y in 18..67) { safePx(p, sx + SIDE_W - 1, y); safePx(p, sx + SIDE_W - 2, y) }

        // Horizontal mortar lines
        p.setColor(sLine)
        for (y in 18..67 step 6) {
            for (x in sx until sx + SIDE_W) safePx(p, x, y)
        }

        // Small side windows (upper + lower floor)
        drawSideWindow(p, sx + 1, 28, 4, 8)
        drawSideWindow(p, sx + 1, 52, 4, 8)

        // Side foundation
        p.fill(sx, 68, SIDE_W, 3, sFound)
        p.fill(sx, 71, SIDE_W, 3, sFoundDk)

        // Right-edge darkening (depth cue)
        p.setColor(Color(0f, 0f, 0f, 0.15f))
        for (y in 0..73) safePx(p, sx + SIDE_W - 1, y)
    }

    private fun drawSideWindow(p: Pixmap, x: Int, y: Int, w: Int, h: Int) {
        val frame   = Color(0.80f, 0.76f, 0.70f, 1f)
        val glass   = Color(0.32f, 0.48f, 0.68f, 1f)
        val glassLt = Color(0.48f, 0.64f, 0.82f, 1f)

        p.fill(x, y, w, h, frame)
        p.fill(x + 1, y + 1, w - 2, h - 2, glass)
        p.setColor(glassLt); safePx(p, x + 1, y + 1)
        // Cross bar
        p.setColor(frame)
        for (sy in y + 1 until y + h - 1) safePx(p, x + w / 2, sy)
        for (ex in x + 1 until x + w - 1) safePx(p, ex, y + h / 2)
    }

    // ── Front face edge shadows ─────────────────────────────────────────────

    private fun drawFrontEdgeShadows(p: Pixmap, mw: Int) {
        p.setColor(Color(0.90f, 0.87f, 0.82f, 1f))
        for (y in 20..73) { safePx(p, 0, y); safePx(p, mw - 1, y) }
    }

    // ═════════════════════════════════════════════════════════════════════════
    //  COMPONENT DRAWING
    // ═════════════════════════════════════════════════════════════════════════

    private fun drawWindow(p: Pixmap, x: Int, y: Int, w: Int, h: Int) {
        val shutter   = Color(0.16f, 0.32f, 0.72f, 1f)
        val shutterLt = Color(0.24f, 0.42f, 0.82f, 1f)
        val shutterDk = Color(0.10f, 0.22f, 0.54f, 1f)
        val glass     = Color(0.28f, 0.50f, 0.76f, 1f)
        val glassLt   = Color(0.50f, 0.72f, 0.92f, 1f)
        val frame     = Color(0.88f, 0.84f, 0.78f, 1f)
        val sill      = Color(0.92f, 0.88f, 0.82f, 1f)

        val shutW = 3
        val gx = x + shutW
        val gw = w - shutW * 2

        // Left shutter
        p.fill(x, y + 1, shutW, h - 2, shutter)
        p.setColor(shutterLt)
        for (sy in y + 1 until y + h - 1) safePx(p, x, sy)
        p.setColor(shutterDk)
        for (sy in y + 1 until y + h - 1 step 3) {
            for (sx in x until x + shutW) safePx(p, sx, sy)
        }

        // Right shutter
        val rx = x + w - shutW
        p.fill(rx, y + 1, shutW, h - 2, shutter)
        p.setColor(shutterLt)
        for (sy in y + 1 until y + h - 1) safePx(p, rx, sy)
        p.setColor(shutterDk)
        for (sy in y + 1 until y + h - 1 step 3) {
            for (sx in rx until rx + shutW) safePx(p, sx, sy)
        }

        // Frame + glass
        p.fill(gx, y, gw, h, frame)
        p.fill(gx + 1, y + 1, gw - 2, h - 2, glass)

        // Glass reflection (top-left highlight)
        p.setColor(glassLt)
        safePx(p, gx + 1, y + 1); safePx(p, gx + 2, y + 1); safePx(p, gx + 1, y + 2)

        // Cross bars
        p.setColor(frame)
        val midX = gx + gw / 2
        for (sy in y + 1 until y + h - 1) safePx(p, midX, sy)
        val midY = y + h / 2
        for (sx in gx + 1 until gx + gw - 1) safePx(p, sx, midY)

        // Sill
        p.setColor(sill)
        for (sx in x until x + w) safePx(p, sx, y + h - 1)
    }

    private fun drawFlowerBox(p: Pixmap, x: Int, y: Int, w: Int) {
        val box    = Color(0.52f, 0.34f, 0.14f, 1f)
        val pink   = Color(0.92f, 0.24f, 0.48f, 1f)
        val mag    = Color(0.84f, 0.16f, 0.38f, 1f)
        val pinkLt = Color(0.98f, 0.50f, 0.66f, 1f)
        val leaf   = Color(0.22f, 0.52f, 0.16f, 1f)

        p.fill(x, y + 1, w, 2, box)
        for (fx in x until x + w) {
            val c = when ((fx - x) % 3) { 0 -> pink; 1 -> mag; else -> pinkLt }
            p.setColor(c); safePx(p, fx, y)
        }
        p.setColor(leaf); safePx(p, x, y); safePx(p, x + w - 1, y)
    }

    private fun drawAwning(
        p: Pixmap, x: Int, y: Int, w: Int, h: Int,
        accent: Color, accentDk: Color
    ) {
        val white = Color(0.98f, 0.96f, 0.92f, 1f)
        for (row in 0 until h) {
            val c = if ((row / 2) % 2 == 0) accent else white
            p.setColor(c)
            for (ax in x until x + w) safePx(p, ax, y + row)
        }
        p.setColor(accentDk)
        for (ax in x until x + w) safePx(p, ax, y + h - 1)
        // Scalloped fringe
        for (ax in x until x + w) {
            val c = if ((ax - x) % 4 < 2) accent else white
            p.setColor(c); safePx(p, ax, y + h)
        }
    }

    private fun drawDoor(p: Pixmap, x: Int, y: Int, w: Int, h: Int) {
        val wood   = Color(0.48f, 0.30f, 0.12f, 1f)
        val woodLt = Color(0.60f, 0.42f, 0.20f, 1f)
        val woodDk = Color(0.32f, 0.18f, 0.06f, 1f)
        val frame  = Color(0.40f, 0.26f, 0.10f, 1f)
        val knob   = Color(0.94f, 0.82f, 0.20f, 1f)
        val knobHi = Color(0.98f, 0.94f, 0.50f, 1f)

        // Arched top — semicircle of frame above door
        p.setColor(frame)
        val cx = x + w / 2
        for (ax in x until x + w) {
            val dx = ax - cx
            val archH = kotlin.math.sqrt((w / 2.0) * (w / 2.0) - dx.toDouble() * dx).toInt()
            for (ay in y - archH until y) safePx(p, ax, ay)
        }

        // Frame
        p.fill(x, y, w, h, frame)
        // Wood fill
        p.fill(x + 1, y + 1, w - 2, h - 1, wood)

        // Door panels (2×2)
        val panW = (w - 4) / 2
        val panH = (h - 4) / 2
        if (panW > 0 && panH > 0) {
            p.fill(x + 2, y + 2, panW, panH, woodLt)
            p.fill(x + 2 + panW + 1, y + 2, panW, panH, woodLt)
            p.fill(x + 2, y + 3 + panH, panW, panH, woodLt)
            p.fill(x + 2 + panW + 1, y + 3 + panH, panW, panH, woodLt)
        }

        // Center split
        p.setColor(woodDk)
        for (sy in y + 1 until y + h) safePx(p, cx, sy)
        // Mid rail
        val my = y + h / 2
        for (sx in x + 1 until x + w - 1) safePx(p, sx, my)

        // Knob
        p.setColor(knob); safePx(p, cx + 1, my + 1)
        p.setColor(knobHi); safePx(p, cx + 1, my + 2)
    }

    // ═════════════════════════════════════════════════════════════════════════
    //  HELPERS
    // ═════════════════════════════════════════════════════════════════════════

    /** Bounds-checked fillRectangle. */
    private fun Pixmap.fill(x: Int, y: Int, w: Int, h: Int, c: Color) {
        setColor(c); fillRectangle(x, y, w, h)
    }

    /** Bounds-checked single pixel draw (color must be set beforehand). */
    private fun safePx(p: Pixmap, x: Int, y: Int) {
        if (x in 0 until p.width && y in 0 until p.height) p.drawPixel(x, y)
    }

    override fun dispose() {
        textures.values.forEach { it.dispose() }
    }
}
