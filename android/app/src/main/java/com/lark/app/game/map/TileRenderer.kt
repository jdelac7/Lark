package com.lark.app.game.map

import com.badlogic.gdx.graphics.Color
import com.badlogic.gdx.graphics.Pixmap
import com.badlogic.gdx.graphics.Texture
import com.badlogic.gdx.graphics.g2d.SpriteBatch
import com.badlogic.gdx.utils.Disposable

/**
 * Mediterranean beach-town tile renderer — bold pixel art, high saturation.
 *
 * Ground: 16×16.  Objects extend upward for 3/4 perspective.
 * WALL/DOOR: 16×64 (south facade). ROOF: 16×56 (peeks above wall).
 * Decorations: BUSH_FLOWER 16×20, POT 16×18, CRATE 16×18,
 *              AWNING 16×12, LAMP 16×28.
 */
class TileRenderer : Disposable {

    private val textures = HashMap<TileType, Texture>()
    private var waterFrame = 0
    private var waterTimer = 0f
    private lateinit var waterAlt: Texture
    private lateinit var waterDeepAlt: Texture

    private lateinit var wallFrontTexture: Texture
    private lateinit var doorFrontTexture: Texture
    private lateinit var roofTallTexture: Texture
    private lateinit var palmTrunkTallTexture: Texture
    private lateinit var palmTopWideTexture: Texture
    private lateinit var fenceTallTexture: Texture
    private lateinit var signTallTexture: Texture
    private lateinit var bushFlowerTexture: Texture
    private lateinit var potTexture: Texture
    private lateinit var crateTexture: Texture
    private lateinit var awningTexture: Texture
    private lateinit var lampTexture: Texture

    companion object {
        val OBJECT_TILES = setOf(
            TileType.WALL, TileType.DOOR, TileType.ROOF,
            TileType.PALM_TRUNK, TileType.PALM_TOP,
            TileType.FENCE, TileType.SIGN,
            TileType.BUSH_FLOWER, TileType.POT, TileType.CRATE,
            TileType.AWNING, TileType.LAMP
        )
    }

    init { createAllTextures() }

    fun isObjectTile(type: TileType): Boolean = type in OBJECT_TILES

    fun isBuildingTile(type: TileType): Boolean =
        type == TileType.WALL || type == TileType.DOOR || type == TileType.ROOF

    private fun createAllTextures() {
        textures[TileType.SAND] = createSandTexture()
        textures[TileType.WATER] = createWaterTexture(false)
        textures[TileType.WATER_DEEP] = createDeepWaterTexture(false)
        textures[TileType.GRASS] = createGrassTexture()
        textures[TileType.PATH_STONE] = createStonePathTexture()
        textures[TileType.PATH_WOOD] = createWoodPathTexture()
        textures[TileType.FLOWERS] = createFlowersTexture()
        textures[TileType.DOCK] = createDockTexture()
        waterAlt = createWaterTexture(true)
        waterDeepAlt = createDeepWaterTexture(true)
        roofTallTexture = createRoofTexture()
        wallFrontTexture = createWallFrontTexture()
        doorFrontTexture = createDoorFrontTexture()
        palmTrunkTallTexture = createPalmTrunkTexture()
        palmTopWideTexture = createPalmTopTexture()
        fenceTallTexture = createFenceTexture()
        signTallTexture = createSignTexture()
        bushFlowerTexture = createBushFlowerTexture()
        potTexture = createPotTexture()
        crateTexture = createCrateTexture()
        awningTexture = createAwningTexture()
        lampTexture = createLampTexture()
    }

    fun renderGround(batch: SpriteBatch, mapData: IntArray, width: Int, height: Int) {
        val ts = 16f
        for (y in 0 until height) {
            for (x in 0 until width) {
                val type = TileType.fromInt(mapData[y * width + x])
                if (!isObjectTile(type)) {
                    batch.draw(getTexture(type), x * ts, y * ts, ts, ts)
                } else {
                    val ground = when (type) {
                        TileType.PALM_TRUNK, TileType.PALM_TOP -> textures[TileType.SAND]!!
                        TileType.CRATE -> textures[TileType.PATH_STONE]!!
                        else -> textures[TileType.GRASS]!!
                    }
                    batch.draw(ground, x * ts, y * ts, ts, ts)
                }
            }
        }
    }

    fun renderObjectTile(
        batch: SpriteBatch, x: Int, y: Int, type: TileType,
        mapData: IntArray, mapWidth: Int, mapHeight: Int
    ) {
        val ts = 16f
        when (type) {
            TileType.ROOF -> {
                batch.draw(roofTallTexture, x * ts, y * ts, ts, 56f)
            }
            TileType.WALL -> {
                if (isSouthFace(mapData, x, y, mapWidth)) {
                    batch.draw(wallFrontTexture, x * ts, y * ts, ts, 64f)
                }
            }
            TileType.DOOR -> {
                batch.draw(doorFrontTexture, x * ts, y * ts, ts, 64f)
            }
            TileType.PALM_TRUNK -> {
                batch.draw(palmTrunkTallTexture, x * ts, y * ts, ts, 28f)
            }
            TileType.PALM_TOP -> {
                batch.draw(palmTopWideTexture, x * ts - 8f, y * ts, 32f, 32f)
            }
            TileType.FENCE -> {
                batch.draw(fenceTallTexture, x * ts, y * ts, ts, 20f)
            }
            TileType.SIGN -> {
                batch.draw(signTallTexture, x * ts, y * ts, ts, 22f)
            }
            TileType.BUSH_FLOWER -> {
                batch.draw(bushFlowerTexture, x * ts, y * ts, ts, 20f)
            }
            TileType.POT -> {
                batch.draw(potTexture, x * ts, y * ts, ts, 18f)
            }
            TileType.CRATE -> {
                batch.draw(crateTexture, x * ts, y * ts, ts, 18f)
            }
            TileType.AWNING -> {
                batch.draw(awningTexture, x * ts, y * ts, ts, 12f)
            }
            TileType.LAMP -> {
                batch.draw(lampTexture, x * ts, y * ts, ts, 28f)
            }
            else -> {}
        }
    }

    private fun isSouthFace(mapData: IntArray, x: Int, y: Int, w: Int): Boolean {
        if (y == 0) return true
        val belowIdx = (y - 1) * w + x
        if (belowIdx < 0 || belowIdx >= mapData.size) return true
        val below = TileType.fromInt(mapData[belowIdx])
        return below != TileType.WALL && below != TileType.DOOR
    }

    fun updateAnimation(delta: Float) {
        waterTimer += delta
        if (waterTimer >= 0.8f) {
            waterTimer = 0f
            waterFrame = 1 - waterFrame
        }
    }

    private fun getTexture(type: TileType): Texture {
        if (waterFrame == 1) {
            if (type == TileType.WATER) return waterAlt
            if (type == TileType.WATER_DEEP) return waterDeepAlt
        }
        return textures[type] ?: textures[TileType.GRASS]!!
    }

    private fun finish(pixmap: Pixmap): Texture {
        val tex = Texture(pixmap)
        tex.setFilter(Texture.TextureFilter.Nearest, Texture.TextureFilter.Nearest)
        pixmap.dispose()
        return tex
    }

    private fun Pixmap.px(x: Int, y: Int, c: Color) { setColor(c); drawPixel(x, y) }
    private fun Pixmap.rect(x: Int, y: Int, w: Int, h: Int, c: Color) { setColor(c); fillRectangle(x, y, w, h) }

    // =====================================================================
    // GROUND TILES  (16×16)
    // =====================================================================

    private fun createGrassTexture(): Texture {
        val p = Pixmap(16, 16, Pixmap.Format.RGBA8888)
        val base = Color(0.30f, 0.60f, 0.22f, 1f)
        val dark = Color(0.20f, 0.46f, 0.14f, 1f)
        val vdark = Color(0.14f, 0.36f, 0.10f, 1f)
        val mid = Color(0.36f, 0.66f, 0.28f, 1f)
        val light = Color(0.46f, 0.76f, 0.34f, 1f)
        val bright = Color(0.56f, 0.84f, 0.42f, 1f)

        p.setColor(base); p.fill()

        p.setColor(dark)
        p.drawPixel(1, 2); p.drawPixel(2, 2); p.drawPixel(1, 3)
        p.drawPixel(9, 0); p.drawPixel(10, 0)
        p.drawPixel(6, 8); p.drawPixel(7, 8); p.drawPixel(7, 9)
        p.drawPixel(13, 5); p.drawPixel(14, 5)
        p.drawPixel(0, 12); p.drawPixel(1, 12)
        p.drawPixel(11, 11); p.drawPixel(12, 11)
        p.drawPixel(4, 14)

        p.setColor(vdark)
        p.drawPixel(2, 3); p.drawPixel(10, 1); p.drawPixel(12, 12)
        p.drawPixel(14, 6)

        p.setColor(mid)
        p.drawPixel(4, 1); p.drawPixel(8, 4); p.drawPixel(15, 3)
        p.drawPixel(3, 7); p.drawPixel(11, 6); p.drawPixel(0, 10)
        p.drawPixel(9, 13); p.drawPixel(14, 10); p.drawPixel(5, 5)
        p.drawPixel(12, 14); p.drawPixel(7, 2)

        p.setColor(light)
        p.drawPixel(3, 1); p.drawPixel(11, 0); p.drawPixel(8, 3)
        p.drawPixel(14, 4); p.drawPixel(0, 11); p.drawPixel(13, 10)
        p.drawPixel(5, 7); p.drawPixel(9, 12); p.drawPixel(2, 15)
        p.drawPixel(7, 14); p.drawPixel(15, 8)

        p.setColor(bright)
        p.drawPixel(3, 0); p.drawPixel(8, 7); p.drawPixel(14, 3)
        p.drawPixel(0, 10); p.drawPixel(12, 10); p.drawPixel(6, 13)

        return finish(p)
    }

    private fun createSandTexture(): Texture {
        val p = Pixmap(16, 16, Pixmap.Format.RGBA8888)
        val base = Color(0.94f, 0.84f, 0.56f, 1f)
        val warm = Color(0.96f, 0.88f, 0.64f, 1f)
        val dark = Color(0.84f, 0.72f, 0.44f, 1f)
        val light = Color(0.98f, 0.92f, 0.72f, 1f)
        val shell = Color(0.99f, 0.96f, 0.90f, 1f)

        p.setColor(base); p.fill()
        p.setColor(warm)
        p.drawPixel(2, 1); p.drawPixel(3, 1); p.drawPixel(8, 3); p.drawPixel(9, 3)
        p.drawPixel(13, 7); p.drawPixel(14, 7); p.drawPixel(1, 10); p.drawPixel(6, 13)
        p.setColor(dark)
        p.drawPixel(5, 2); p.drawPixel(12, 1); p.drawPixel(0, 6)
        p.drawPixel(7, 8); p.drawPixel(15, 5); p.drawPixel(3, 12)
        p.drawPixel(10, 14); p.drawPixel(14, 11)
        p.setColor(light)
        p.drawPixel(6, 0); p.drawPixel(13, 4); p.drawPixel(1, 8)
        p.drawPixel(10, 10); p.drawPixel(4, 15); p.drawPixel(15, 13)
        p.setColor(shell)
        p.drawPixel(11, 5); p.drawPixel(3, 11)

        return finish(p)
    }

    private fun createWaterTexture(alt: Boolean): Texture {
        val p = Pixmap(16, 16, Pixmap.Format.RGBA8888)
        val deep = Color(0.06f, 0.30f, 0.72f, 1f)
        val base = Color(0.10f, 0.40f, 0.82f, 1f)
        val mid = Color(0.18f, 0.50f, 0.88f, 1f)
        val light = Color(0.30f, 0.62f, 0.94f, 1f)
        val bright = Color(0.48f, 0.76f, 0.98f, 1f)
        val foam = Color(0.70f, 0.88f, 0.98f, 1f)
        val sparkle = Color(0.90f, 0.96f, 1.0f, 1f)

        p.setColor(base); p.fill()
        val off = if (alt) 5 else 0

        // Deep patches
        p.setColor(deep)
        p.drawPixel(2, 1); p.drawPixel(3, 1); p.drawPixel(3, 2)
        p.drawPixel(9, 7); p.drawPixel(10, 7); p.drawPixel(10, 8)
        p.drawPixel(0, 13); p.drawPixel(1, 13)

        // Mid variation
        p.setColor(mid)
        for (x in 0..15) if ((x + off) % 4 == 0) p.drawPixel(x, 6)
        p.drawPixel(5, 0); p.drawPixel(12, 2); p.drawPixel(7, 12); p.drawPixel(14, 14)

        // Wave crests
        p.setColor(light)
        for (i in 0..2) {
            val wx = ((i * 6) + off) % 16
            p.drawPixel(wx, 3); p.drawPixel((wx + 1) % 16, 3)
            p.drawPixel((wx + 2) % 16, 3); p.drawPixel((wx + 3) % 16, 4)
        }
        for (i in 0..2) {
            val wx = ((i * 6) + off + 3) % 16
            p.drawPixel(wx, 10); p.drawPixel((wx + 1) % 16, 10)
            p.drawPixel((wx + 2) % 16, 10); p.drawPixel((wx + 3) % 16, 11)
        }

        // Bright wave caps
        p.setColor(bright)
        p.drawPixel((1 + off) % 16, 3); p.drawPixel((2 + off) % 16, 3)
        p.drawPixel((7 + off) % 16, 3)
        p.drawPixel((4 + off) % 16, 10); p.drawPixel((5 + off) % 16, 10)

        // Foam
        p.setColor(foam)
        p.drawPixel((2 + off) % 16, 2); p.drawPixel((8 + off) % 16, 2)
        p.drawPixel((5 + off) % 16, 9); p.drawPixel((11 + off) % 16, 9)

        // Sparkle
        p.setColor(sparkle)
        p.drawPixel((3 + off) % 16, 2); p.drawPixel((14 + off) % 16, 9)
        p.drawPixel((9 + off) % 16, 1); p.drawPixel((6 + off) % 16, 14)

        return finish(p)
    }

    private fun createDeepWaterTexture(alt: Boolean): Texture {
        val p = Pixmap(16, 16, Pixmap.Format.RGBA8888)
        val abyss = Color(0.02f, 0.10f, 0.36f, 1f)
        val base = Color(0.04f, 0.16f, 0.48f, 1f)
        val mid = Color(0.08f, 0.24f, 0.56f, 1f)
        val light = Color(0.12f, 0.32f, 0.64f, 1f)
        val hl = Color(0.18f, 0.40f, 0.72f, 1f)

        p.setColor(base); p.fill()
        val off = if (alt) 4 else 0

        p.setColor(abyss)
        p.drawPixel(3, 2); p.drawPixel(4, 2); p.drawPixel(12, 8); p.drawPixel(13, 8)
        p.drawPixel(1, 11); p.drawPixel(8, 14)

        p.setColor(mid)
        for (i in 0..2) {
            val wx = ((i * 5) + off) % 16
            p.drawPixel(wx, 5); p.drawPixel((wx + 1) % 16, 5); p.drawPixel((wx + 2) % 16, 5)
            p.drawPixel((wx + 2) % 16, 12); p.drawPixel((wx + 3) % 16, 12)
        }

        p.setColor(light)
        p.drawPixel((3 + off) % 16, 5); p.drawPixel((10 + off) % 16, 12)

        p.setColor(hl)
        p.drawPixel((4 + off) % 16, 4); p.drawPixel((11 + off) % 16, 11)

        return finish(p)
    }

    private fun createStonePathTexture(): Texture {
        val p = Pixmap(16, 16, Pixmap.Format.RGBA8888)
        // Warm golden-brown cobblestone (matching reference)
        val grout = Color(0.50f, 0.40f, 0.28f, 1f)
        val stone1 = Color(0.76f, 0.64f, 0.46f, 1f)
        val stone2 = Color(0.70f, 0.58f, 0.40f, 1f)
        val stone3 = Color(0.80f, 0.68f, 0.50f, 1f)
        val hl = Color(0.88f, 0.78f, 0.60f, 1f)
        val shadow = Color(0.56f, 0.46f, 0.32f, 1f)

        p.setColor(grout); p.fill()

        // Large cobbles
        p.rect(0, 0, 7, 7, stone1)
        p.rect(8, 0, 8, 7, stone2)
        p.rect(0, 8, 5, 8, stone3)
        p.rect(6, 8, 5, 8, stone1)
        p.rect(12, 8, 4, 8, stone2)

        // Top-left highlights on each stone
        p.setColor(hl)
        p.drawPixel(1, 1); p.drawPixel(2, 1); p.drawPixel(1, 2)
        p.drawPixel(9, 1); p.drawPixel(10, 1); p.drawPixel(9, 2)
        p.drawPixel(1, 9); p.drawPixel(2, 9)
        p.drawPixel(7, 9); p.drawPixel(8, 9)
        p.drawPixel(13, 9)

        // Bottom-right shadows
        p.setColor(shadow)
        p.drawPixel(6, 6); p.drawPixel(5, 6); p.drawPixel(6, 5)
        p.drawPixel(15, 6); p.drawPixel(14, 6); p.drawPixel(15, 5)
        p.drawPixel(4, 15); p.drawPixel(3, 15)
        p.drawPixel(10, 15); p.drawPixel(9, 15)
        p.drawPixel(15, 15); p.drawPixel(14, 15)

        return finish(p)
    }

    private fun createWoodPathTexture(): Texture {
        val p = Pixmap(16, 16, Pixmap.Format.RGBA8888)
        val plank = Color(0.58f, 0.40f, 0.18f, 1f)
        val plankW = Color(0.66f, 0.48f, 0.24f, 1f)
        val gap = Color(0.28f, 0.16f, 0.06f, 1f)
        val grain = Color(0.50f, 0.34f, 0.14f, 1f)
        val light = Color(0.74f, 0.56f, 0.32f, 1f)
        val nail = Color(0.48f, 0.46f, 0.42f, 1f)

        p.setColor(plank); p.fill()
        p.setColor(plankW); p.fillRectangle(0, 4, 16, 3); p.fillRectangle(0, 12, 16, 3)
        p.setColor(gap)
        for (x in 0..15) { p.drawPixel(x, 3); p.drawPixel(x, 7); p.drawPixel(x, 11); p.drawPixel(x, 15) }
        p.setColor(grain)
        p.drawPixel(3, 1); p.drawPixel(4, 1); p.drawPixel(10, 0)
        p.drawPixel(7, 5); p.drawPixel(8, 5); p.drawPixel(13, 4)
        p.drawPixel(12, 9); p.drawPixel(13, 9); p.drawPixel(1, 8)
        p.drawPixel(9, 13); p.drawPixel(10, 13); p.drawPixel(14, 12)
        p.setColor(light)
        p.drawPixel(6, 0); p.drawPixel(14, 2); p.drawPixel(3, 4); p.drawPixel(11, 6)
        p.drawPixel(8, 8); p.drawPixel(1, 10); p.drawPixel(14, 12); p.drawPixel(6, 14)
        p.setColor(nail)
        p.drawPixel(2, 0); p.drawPixel(13, 0); p.drawPixel(2, 4); p.drawPixel(13, 4)
        p.drawPixel(2, 8); p.drawPixel(13, 8); p.drawPixel(2, 12); p.drawPixel(13, 12)

        return finish(p)
    }

    private fun createFlowersTexture(): Texture {
        val p = Pixmap(16, 16, Pixmap.Format.RGBA8888)
        val grass = Color(0.30f, 0.60f, 0.22f, 1f)
        val grassDk = Color(0.20f, 0.46f, 0.14f, 1f)
        val grassLt = Color(0.40f, 0.70f, 0.30f, 1f)
        p.setColor(grass); p.fill()
        p.setColor(grassDk)
        p.drawPixel(3, 4); p.drawPixel(9, 2); p.drawPixel(13, 10); p.drawPixel(1, 14)
        p.setColor(grassLt)
        p.drawPixel(5, 1); p.drawPixel(11, 7); p.drawPixel(0, 9); p.drawPixel(14, 14)

        // Pink bougainvillea
        val pink = Color(0.92f, 0.24f, 0.48f, 1f)
        val pinkLt = Color(0.98f, 0.48f, 0.64f, 1f)
        p.setColor(pink)
        p.drawPixel(3, 2); p.drawPixel(4, 1); p.drawPixel(4, 3); p.drawPixel(5, 2)
        p.drawPixel(12, 8); p.drawPixel(13, 7); p.drawPixel(13, 9); p.drawPixel(14, 8)
        p.setColor(pinkLt); p.drawPixel(4, 2); p.drawPixel(13, 8)

        // Yellow
        val yellow = Color(0.98f, 0.90f, 0.12f, 1f)
        val yelCtr = Color(0.90f, 0.60f, 0.10f, 1f)
        p.setColor(yellow)
        p.drawPixel(8, 5); p.drawPixel(9, 4); p.drawPixel(9, 6); p.drawPixel(10, 5)
        p.drawPixel(2, 11); p.drawPixel(3, 12); p.drawPixel(1, 12)
        p.setColor(yelCtr); p.drawPixel(9, 5); p.drawPixel(2, 12)

        // Purple
        val purple = Color(0.72f, 0.20f, 0.90f, 1f)
        p.setColor(purple)
        p.drawPixel(6, 10); p.drawPixel(7, 9); p.drawPixel(7, 11); p.drawPixel(8, 10)
        p.setColor(Color(0.96f, 0.80f, 0.96f, 1f)); p.drawPixel(7, 10)

        // White
        p.setColor(Color(0.98f, 0.98f, 1.0f, 1f))
        p.drawPixel(1, 6); p.drawPixel(2, 7); p.drawPixel(11, 13); p.drawPixel(15, 3)

        // Stems
        p.setColor(Color(0.16f, 0.40f, 0.12f, 1f))
        p.drawPixel(4, 4); p.drawPixel(9, 7); p.drawPixel(7, 12); p.drawPixel(13, 10)

        return finish(p)
    }

    private fun createDockTexture(): Texture {
        val p = Pixmap(16, 16, Pixmap.Format.RGBA8888)
        val water = Color(0.10f, 0.40f, 0.82f, 1f)
        val plank = Color(0.54f, 0.38f, 0.18f, 1f)
        val plankDk = Color(0.38f, 0.24f, 0.10f, 1f)
        val plankLt = Color(0.66f, 0.50f, 0.28f, 1f)
        val nail = Color(0.46f, 0.44f, 0.40f, 1f)

        p.setColor(water); p.fill()
        p.setColor(plank); p.fillRectangle(1, 1, 14, 14)
        p.setColor(plankDk)
        for (x in 1..14) { p.drawPixel(x, 4); p.drawPixel(x, 8); p.drawPixel(x, 12) }
        for (y in 1..14) { p.drawPixel(1, y); p.drawPixel(14, y) }
        p.setColor(plankLt)
        p.drawPixel(5, 2); p.drawPixel(6, 2); p.drawPixel(10, 3)
        p.drawPixel(4, 6); p.drawPixel(11, 7); p.drawPixel(7, 10); p.drawPixel(12, 11)
        p.setColor(nail)
        p.drawPixel(3, 1); p.drawPixel(12, 1); p.drawPixel(3, 5); p.drawPixel(12, 5)
        p.drawPixel(3, 9); p.drawPixel(12, 9); p.drawPixel(3, 13); p.drawPixel(12, 13)

        return finish(p)
    }

    // =====================================================================
    // 3/4 PERSPECTIVE — BUILDINGS
    // =====================================================================

    private fun createRoofTexture(): Texture {
        val p = Pixmap(16, 56, Pixmap.Format.RGBA8888)
        // Rich terracotta — peeks above wall (top 8px visible)
        val ridgeBr = Color(0.96f, 0.54f, 0.28f, 1f)
        val ridge = Color(0.90f, 0.44f, 0.20f, 1f)
        val roofLt = Color(0.86f, 0.36f, 0.14f, 1f)
        val roofMn = Color(0.76f, 0.28f, 0.10f, 1f)
        val roofDk = Color(0.60f, 0.20f, 0.06f, 1f)
        val eave = Color(0.50f, 0.16f, 0.04f, 1f)
        val shadow = Color(0f, 0f, 0f, 0.28f)

        // Ridge cap (rows 0-5)
        p.rect(0, 0, 16, 3, ridgeBr)
        p.rect(0, 3, 16, 3, ridge)

        // Upper slope (rows 6-18)
        p.rect(0, 6, 16, 13, roofLt)

        // Lower slope (rows 19-38)
        p.rect(0, 19, 16, 20, roofMn)

        // Tile row separators
        p.setColor(roofDk)
        for (x in 0..15) { p.drawPixel(x, 12); p.drawPixel(x, 18); p.drawPixel(x, 24); p.drawPixel(x, 30); p.drawPixel(x, 36) }
        // Staggered verticals
        for (y in 6..12) { p.drawPixel(4, y); p.drawPixel(8, y); p.drawPixel(12, y) }
        for (y in 13..18) { p.drawPixel(2, y); p.drawPixel(6, y); p.drawPixel(10, y); p.drawPixel(14, y) }
        for (y in 19..24) { p.drawPixel(4, y); p.drawPixel(8, y); p.drawPixel(12, y) }
        for (y in 25..30) { p.drawPixel(2, y); p.drawPixel(6, y); p.drawPixel(10, y); p.drawPixel(14, y) }

        // Tile highlights
        p.setColor(roofLt)
        p.drawPixel(2, 7); p.drawPixel(6, 8); p.drawPixel(10, 7); p.drawPixel(14, 8)
        p.drawPixel(4, 14); p.drawPixel(8, 15); p.drawPixel(12, 14)
        p.drawPixel(2, 21); p.drawPixel(10, 22)

        // Eave (rows 39-47)
        p.rect(0, 39, 16, 5, eave)
        p.rect(0, 44, 16, 4, Color(0.40f, 0.12f, 0.02f, 1f))

        // Shadow (rows 48-55)
        p.rect(0, 48, 16, 8, shadow)

        return finish(p)
    }

    /**
     * Build a bold Mediterranean window on a pixmap.
     * (x, y) = top-left of the window area, spanning full 16px width.
     * winH = height allocated for window area.
     */
    private fun drawWindow(p: Pixmap, wy: Int, winH: Int) {
        val white = Color(0.98f, 0.96f, 0.94f, 1f)
        val shutter = Color(0.16f, 0.32f, 0.72f, 1f)
        val shutDk = Color(0.10f, 0.22f, 0.54f, 1f)
        val shutLt = Color(0.24f, 0.42f, 0.82f, 1f)
        val frame = Color(0.86f, 0.82f, 0.76f, 1f)
        val glass = Color(0.28f, 0.50f, 0.76f, 1f)
        val glassLt = Color(0.50f, 0.72f, 0.92f, 1f)
        val glassDk = Color(0.16f, 0.34f, 0.56f, 1f)
        val sill = Color(0.90f, 0.86f, 0.80f, 1f)

        val cy = wy + 2 // leave 2px margin top
        val ch = winH - 4 // minus 2 top + 2 bottom

        // Shutters (bold blue blocks)
        p.rect(2, cy, 3, ch, shutter)
        p.rect(11, cy, 3, ch, shutter)
        // Shutter highlights
        p.setColor(shutLt)
        for (y in cy until cy + ch) p.drawPixel(2, y)
        for (y in cy until cy + ch) p.drawPixel(11, y)
        // Shutter slat lines
        p.setColor(shutDk)
        for (y in cy until cy + ch step 3) { p.drawPixel(3, y); p.drawPixel(4, y); p.drawPixel(12, y); p.drawPixel(13, y) }

        // Window frame
        p.rect(5, cy, 6, ch, frame)
        // Glass
        p.rect(6, cy + 1, 4, ch - 2, glass)
        // Glass reflection
        p.setColor(glassLt)
        p.drawPixel(6, cy + 1); p.drawPixel(7, cy + 1); p.drawPixel(6, cy + 2)
        // Glass shadow
        p.setColor(glassDk)
        p.drawPixel(9, cy + ch - 3); p.drawPixel(9, cy + ch - 2)
        // Cross bar
        p.setColor(frame)
        for (y in cy + 1 until cy + ch - 1) p.drawPixel(8, y) // vertical
        val midY = cy + ch / 2
        for (x in 6..9) p.drawPixel(x, midY) // horizontal

        // Sill below window
        p.setColor(sill)
        for (x in 4..11) p.drawPixel(x, wy + winH - 1)
    }

    /**
     * Draw a flower box under a window at row `y`, 16px wide.
     */
    private fun drawFlowerBox(p: Pixmap, y: Int) {
        val box = Color(0.52f, 0.34f, 0.14f, 1f)
        val pink = Color(0.92f, 0.24f, 0.48f, 1f)
        val mag = Color(0.84f, 0.16f, 0.38f, 1f)
        val leaf = Color(0.24f, 0.52f, 0.18f, 1f)
        val pinkLt = Color(0.98f, 0.50f, 0.66f, 1f)

        // Box
        p.rect(3, y + 1, 10, 2, box)
        // Flowers above box
        p.setColor(pink)
        p.drawPixel(4, y); p.drawPixel(6, y); p.drawPixel(8, y); p.drawPixel(10, y)
        p.setColor(mag)
        p.drawPixel(5, y); p.drawPixel(9, y)
        p.setColor(pinkLt)
        p.drawPixel(7, y); p.drawPixel(11, y)
        p.setColor(leaf)
        p.drawPixel(3, y); p.drawPixel(12, y)
    }

    private fun createWallFrontTexture(): Texture {
        val p = Pixmap(16, 64, Pixmap.Format.RGBA8888)

        // Color palette
        val white = Color(0.98f, 0.96f, 0.94f, 1f)
        val whiteCool = Color(0.94f, 0.92f, 0.88f, 1f)
        val foundation = Color(0.56f, 0.48f, 0.36f, 1f)
        val foundDk = Color(0.42f, 0.36f, 0.26f, 1f)
        val mortar = Color(0.86f, 0.82f, 0.76f, 1f)
        val edgeShadow = Color(0.88f, 0.84f, 0.78f, 1f)

        // Terracotta cornice
        val corrBr = Color(0.96f, 0.54f, 0.28f, 1f)
        val corr = Color(0.88f, 0.40f, 0.18f, 1f)
        val corrDk = Color(0.72f, 0.28f, 0.10f, 1f)
        val corrSh = Color(0.56f, 0.20f, 0.08f, 1f)

        // ── TERRACOTTA CORNICE (rows 0-11, 12px) ──
        p.rect(0, 0, 16, 2, corrBr)   // bright ridge
        p.rect(0, 2, 16, 4, corr)     // main terracotta
        p.rect(0, 6, 16, 3, corrDk)   // lower terracotta
        p.rect(0, 9, 16, 3, corrSh)   // overhang shadow
        // Tile pattern on cornice
        p.setColor(corrDk)
        for (x in 0..15 step 4) { p.drawPixel(x, 3); p.drawPixel(x, 4) }
        p.setColor(corrBr)
        p.drawPixel(2, 3); p.drawPixel(6, 4); p.drawPixel(10, 3); p.drawPixel(14, 4)

        // ── UPPER FLOOR (rows 12-33, 22px) ──
        p.rect(0, 12, 16, 22, white)
        // Plaster variation
        p.setColor(whiteCool)
        p.fillRectangle(1, 18, 5, 6)
        p.fillRectangle(10, 18, 5, 6)

        // Upper window (rows 14-27, 14px tall)
        drawWindow(p, 14, 14)

        // Flower box (row 28-30)
        drawFlowerBox(p, 28)

        // ── FLOOR DIVIDER (rows 31-32) ──
        p.setColor(mortar)
        for (x in 0..15) { p.drawPixel(x, 31); p.drawPixel(x, 32) }

        // ── LOWER FLOOR (rows 33-51, 19px) ──
        p.rect(0, 33, 16, 19, white)
        p.setColor(whiteCool)
        p.fillRectangle(1, 39, 5, 6)
        p.fillRectangle(10, 39, 5, 6)

        // Lower window (rows 35-47, 13px tall)
        drawWindow(p, 35, 13)

        // ── SIDE EDGE SHADOWS ──
        p.setColor(edgeShadow)
        for (y in 12..51) { p.drawPixel(0, y); p.drawPixel(15, y) }

        // ── FOUNDATION (rows 52-59, 8px) ──
        p.rect(0, 52, 16, 4, foundation)
        p.rect(0, 56, 16, 4, foundDk)
        // Stone detail
        p.setColor(Color(0.50f, 0.42f, 0.30f, 1f))
        p.drawPixel(3, 53); p.drawPixel(8, 54); p.drawPixel(12, 53)
        p.drawPixel(2, 57); p.drawPixel(6, 58); p.drawPixel(11, 57); p.drawPixel(14, 58)

        // ── GROUND SHADOW (rows 60-63) ──
        p.rect(0, 60, 16, 4, Color(0f, 0f, 0f, 0.22f))

        return finish(p)
    }

    private fun createDoorFrontTexture(): Texture {
        val p = Pixmap(16, 64, Pixmap.Format.RGBA8888)

        val white = Color(0.98f, 0.96f, 0.94f, 1f)
        val whiteCool = Color(0.94f, 0.92f, 0.88f, 1f)
        val foundation = Color(0.56f, 0.48f, 0.36f, 1f)
        val foundDk = Color(0.42f, 0.36f, 0.26f, 1f)
        val mortar = Color(0.86f, 0.82f, 0.76f, 1f)
        val edgeShadow = Color(0.88f, 0.84f, 0.78f, 1f)

        val corrBr = Color(0.96f, 0.54f, 0.28f, 1f)
        val corr = Color(0.88f, 0.40f, 0.18f, 1f)
        val corrDk = Color(0.72f, 0.28f, 0.10f, 1f)
        val corrSh = Color(0.56f, 0.20f, 0.08f, 1f)

        val doorWood = Color(0.48f, 0.30f, 0.12f, 1f)
        val doorLt = Color(0.60f, 0.42f, 0.20f, 1f)
        val doorDk = Color(0.32f, 0.18f, 0.06f, 1f)
        val doorFrame = Color(0.40f, 0.26f, 0.10f, 1f)
        val knob = Color(0.94f, 0.82f, 0.20f, 1f)

        val awRed = Color(0.86f, 0.18f, 0.12f, 1f)
        val awWhite = Color(0.98f, 0.96f, 0.92f, 1f)
        val awRedDk = Color(0.66f, 0.12f, 0.08f, 1f)

        val step = Color(0.74f, 0.68f, 0.58f, 1f)
        val stepLt = Color(0.82f, 0.76f, 0.66f, 1f)

        // ── TERRACOTTA CORNICE (rows 0-11) — same as wall ──
        p.rect(0, 0, 16, 2, corrBr)
        p.rect(0, 2, 16, 4, corr)
        p.rect(0, 6, 16, 3, corrDk)
        p.rect(0, 9, 16, 3, corrSh)
        p.setColor(corrDk)
        for (x in 0..15 step 4) { p.drawPixel(x, 3); p.drawPixel(x, 4) }
        p.setColor(corrBr)
        p.drawPixel(2, 3); p.drawPixel(6, 4); p.drawPixel(10, 3); p.drawPixel(14, 4)

        // ── UPPER FLOOR (rows 12-33) — window ──
        p.rect(0, 12, 16, 22, white)
        p.setColor(whiteCool)
        p.fillRectangle(1, 18, 5, 6)
        p.fillRectangle(10, 18, 5, 6)
        drawWindow(p, 14, 14)
        drawFlowerBox(p, 28)

        // ── DIVIDER (rows 31-32) ──
        p.setColor(mortar)
        for (x in 0..15) { p.drawPixel(x, 31); p.drawPixel(x, 32) }

        // ── LOWER FLOOR (rows 33-51) — DOOR with awning ──
        p.rect(0, 33, 16, 19, white)

        // Red/white striped awning (rows 33-38)
        for (row in 33..38) {
            val c = if ((row - 33) / 2 % 2 == 0) awRed else awWhite
            p.setColor(c)
            for (x in 2..13) p.drawPixel(x, row)
        }
        // Awning bottom edge
        p.setColor(awRedDk)
        for (x in 2..13) p.drawPixel(x, 38)
        // Scalloped fringe
        p.setColor(awRed)
        p.drawPixel(3, 39); p.drawPixel(4, 39); p.drawPixel(7, 39); p.drawPixel(8, 39); p.drawPixel(11, 39); p.drawPixel(12, 39)
        p.setColor(awWhite)
        p.drawPixel(5, 39); p.drawPixel(6, 39); p.drawPixel(9, 39); p.drawPixel(10, 39)

        // Door frame (rows 40-51)
        p.rect(4, 40, 8, 12, doorFrame)
        // Door wood
        p.rect(5, 41, 6, 10, doorWood)
        // Door panels
        p.setColor(doorLt)
        p.fillRectangle(6, 42, 2, 3); p.fillRectangle(9, 42, 2, 3)
        p.fillRectangle(6, 47, 2, 3); p.fillRectangle(9, 47, 2, 3)
        // Center split + mid rail
        p.setColor(doorDk)
        for (y in 41..50) p.drawPixel(8, y)
        for (x in 5..10) p.drawPixel(x, 46)
        // Knob
        p.setColor(knob); p.drawPixel(10, 47); p.drawPixel(10, 48)
        p.setColor(Color(0.98f, 0.94f, 0.50f, 1f)); p.drawPixel(10, 47)

        // Side edges
        p.setColor(edgeShadow)
        for (y in 12..51) { p.drawPixel(0, y); p.drawPixel(15, y) }

        // ── FOUNDATION (rows 52-59) ──
        p.rect(0, 52, 16, 4, foundation)
        // Door step
        p.rect(3, 52, 10, 3, step)
        p.setColor(stepLt); for (x in 3..12) p.drawPixel(x, 52)
        p.rect(0, 56, 16, 4, foundDk)
        p.rect(3, 56, 10, 1, step)

        // ── GROUND SHADOW ──
        p.rect(0, 60, 16, 4, Color(0f, 0f, 0f, 0.22f))

        return finish(p)
    }

    // =====================================================================
    // 3/4 PERSPECTIVE — NATURE & PROPS
    // =====================================================================

    private fun createPalmTrunkTexture(): Texture {
        val p = Pixmap(16, 28, Pixmap.Format.RGBA8888)
        p.setColor(Color.CLEAR); p.fill()

        val trunk = Color(0.50f, 0.34f, 0.16f, 1f)
        val trunkLt = Color(0.64f, 0.48f, 0.26f, 1f)
        val trunkDk = Color(0.36f, 0.22f, 0.08f, 1f)
        val bark = Color(0.42f, 0.28f, 0.12f, 1f)
        val barkLt = Color(0.56f, 0.40f, 0.22f, 1f)

        for (y in 0..26) {
            val w = when { y < 3 -> 4; y < 7 -> 5; y < 21 -> 6; else -> 5 }
            val sx = 8 - w / 2
            for (x in sx until sx + w) {
                val seg = y % 5
                val c = when {
                    x == sx -> trunkLt; x == sx + w - 1 -> trunkDk
                    seg == 0 -> bark; seg == 1 -> barkLt; else -> trunk
                }
                p.setColor(c); p.drawPixel(x, y)
            }
        }
        p.setColor(Color(0f, 0f, 0f, 0.22f))
        p.drawPixel(4, 27); p.drawPixel(5, 27); p.drawPixel(6, 27)
        p.drawPixel(9, 27); p.drawPixel(10, 27); p.drawPixel(11, 27)

        return finish(p)
    }

    private fun createPalmTopTexture(): Texture {
        val p = Pixmap(32, 32, Pixmap.Format.RGBA8888)
        p.setColor(Color.CLEAR); p.fill()

        val leaf = Color(0.10f, 0.54f, 0.12f, 1f)
        val leafLt = Color(0.24f, 0.70f, 0.20f, 1f)
        val leafDk = Color(0.06f, 0.38f, 0.06f, 1f)
        val leafBr = Color(0.38f, 0.80f, 0.30f, 1f)
        val coco = Color(0.46f, 0.32f, 0.14f, 1f)
        val cocoDk = Color(0.32f, 0.20f, 0.06f, 1f)

        val cx = 16; val cy = 16
        drawFrond(p, cx, cy, -9, -11, leaf, leafLt, leafDk)
        drawFrond(p, cx, cy, 0, -13, leaf, leafLt, leafDk)
        drawFrond(p, cx, cy, 9, -11, leaf, leafLt, leafDk)
        drawFrond(p, cx, cy, -12, -3, leaf, leafLt, leafDk)
        drawFrond(p, cx, cy, 12, -3, leaf, leafLt, leafDk)
        drawFrond(p, cx, cy, -8, 7, leaf, leafLt, leafDk)
        drawFrond(p, cx, cy, 0, 9, leaf, leafLt, leafDk)
        drawFrond(p, cx, cy, 8, 7, leaf, leafLt, leafDk)

        p.setColor(leafDk); p.fillRectangle(cx - 3, cy - 3, 6, 6)
        p.setColor(leaf); p.fillRectangle(cx - 2, cy - 2, 4, 4)
        p.setColor(leafBr); p.drawPixel(cx - 1, cy - 1); p.drawPixel(cx, cy - 1)

        p.setColor(coco)
        p.drawPixel(cx - 2, cy + 2); p.drawPixel(cx - 1, cy + 2); p.drawPixel(cx, cy + 2)
        p.drawPixel(cx + 1, cy + 1); p.drawPixel(cx + 2, cy + 1)
        p.setColor(cocoDk)
        p.drawPixel(cx - 1, cy + 3); p.drawPixel(cx, cy + 3); p.drawPixel(cx + 2, cy + 2)

        return finish(p)
    }

    private fun drawFrond(
        p: Pixmap, cx: Int, cy: Int, dx: Int, dy: Int,
        main: Color, light: Color, dark: Color
    ) {
        val steps = maxOf(kotlin.math.abs(dx), kotlin.math.abs(dy))
        for (i in 0..steps) {
            val t = i.toFloat() / steps
            val x = cx + (dx * t).toInt()
            val y = cy + (dy * t).toInt()
            val th = if (i < steps / 3) 3 else if (i < steps * 2 / 3) 2 else 1
            for (off in -th / 2..th / 2) {
                val px = x + if (kotlin.math.abs(dy) > kotlin.math.abs(dx)) off else 0
                val py = y + if (kotlin.math.abs(dx) >= kotlin.math.abs(dy)) off else 0
                if (px in 0..31 && py in 0..31) {
                    p.setColor(if (off == -th / 2) light else if (off == th / 2) dark else main)
                    p.drawPixel(px, py)
                }
            }
            if (i > steps / 3 && i % 2 == 0) {
                val lx1 = x + if (kotlin.math.abs(dy) > kotlin.math.abs(dx)) 2 else 0
                val ly1 = y + if (kotlin.math.abs(dx) >= kotlin.math.abs(dy)) 2 else 0
                val lx2 = x + if (kotlin.math.abs(dy) > kotlin.math.abs(dx)) -2 else 0
                val ly2 = y + if (kotlin.math.abs(dx) >= kotlin.math.abs(dy)) -2 else 0
                if (lx1 in 0..31 && ly1 in 0..31) { p.setColor(light); p.drawPixel(lx1, ly1) }
                if (lx2 in 0..31 && ly2 in 0..31) { p.setColor(dark); p.drawPixel(lx2, ly2) }
            }
        }
    }

    private fun createFenceTexture(): Texture {
        val p = Pixmap(16, 20, Pixmap.Format.RGBA8888)
        p.setColor(Color.CLEAR); p.fill()

        val wood = Color(0.88f, 0.82f, 0.70f, 1f)
        val woodLt = Color(0.94f, 0.90f, 0.80f, 1f)
        val woodDk = Color(0.70f, 0.64f, 0.52f, 1f)
        val rail = Color(0.76f, 0.70f, 0.58f, 1f)

        for (px in intArrayOf(1, 6, 11)) {
            p.px(px + 1, 1, woodLt)
            p.setColor(wood); p.fillRectangle(px, 2, 3, 15)
            p.setColor(woodLt); for (y in 2..16) p.drawPixel(px, y)
            p.setColor(woodDk); for (y in 2..16) p.drawPixel(px + 2, y)
        }
        p.setColor(rail)
        for (x in 0..15) { p.drawPixel(x, 6); p.drawPixel(x, 7); p.drawPixel(x, 13); p.drawPixel(x, 14) }
        p.setColor(woodDk)
        for (x in 0..15) { p.drawPixel(x, 7); p.drawPixel(x, 14) }
        p.setColor(Color(0f, 0f, 0f, 0.14f))
        for (x in 0..15) { p.drawPixel(x, 18); p.drawPixel(x, 19) }

        return finish(p)
    }

    private fun createSignTexture(): Texture {
        val p = Pixmap(16, 22, Pixmap.Format.RGBA8888)
        p.setColor(Color.CLEAR); p.fill()

        val post = Color(0.50f, 0.34f, 0.16f, 1f)
        val postDk = Color(0.38f, 0.24f, 0.10f, 1f)
        val board = Color(0.92f, 0.86f, 0.66f, 1f)
        val boardDk = Color(0.76f, 0.70f, 0.50f, 1f)
        val boardLt = Color(0.96f, 0.92f, 0.76f, 1f)
        val text = Color(0.28f, 0.22f, 0.14f, 1f)

        p.rect(7, 10, 2, 12, post)
        p.setColor(postDk); p.drawPixel(8, 12); p.drawPixel(8, 15); p.drawPixel(8, 18); p.drawPixel(8, 21)
        p.rect(1, 1, 14, 10, board)
        p.setColor(boardDk)
        for (x in 1..14) { p.drawPixel(x, 1); p.drawPixel(x, 10) }
        for (y in 1..10) { p.drawPixel(1, y); p.drawPixel(14, y) }
        p.setColor(boardLt); for (x in 2..13) p.drawPixel(x, 2)
        p.drawPixel(2, 3); p.drawPixel(2, 4)
        p.setColor(text)
        for (x in 3..12) p.drawPixel(x, 5)
        for (x in 4..11) p.drawPixel(x, 7)

        p.setColor(Color(0f, 0f, 0f, 0.10f))
        p.drawPixel(5, 21); p.drawPixel(6, 21); p.drawPixel(9, 21); p.drawPixel(10, 21)
        return finish(p)
    }

    // =====================================================================
    // DECORATIONS
    // =====================================================================

    private fun createBushFlowerTexture(): Texture {
        val p = Pixmap(16, 20, Pixmap.Format.RGBA8888)
        p.setColor(Color.CLEAR); p.fill()

        val bushDk = Color(0.14f, 0.38f, 0.10f, 1f)
        val bush = Color(0.22f, 0.52f, 0.16f, 1f)
        val bushLt = Color(0.34f, 0.66f, 0.26f, 1f)
        val bushBr = Color(0.44f, 0.76f, 0.34f, 1f)
        val pink = Color(0.92f, 0.24f, 0.48f, 1f)
        val mag = Color(0.84f, 0.16f, 0.38f, 1f)
        val pinkLt = Color(0.98f, 0.50f, 0.66f, 1f)

        // Bush body — rounded mass
        p.rect(4, 2, 8, 2, bush)
        p.rect(3, 4, 10, 2, bush)
        p.rect(2, 6, 12, 8, bush)
        p.rect(3, 14, 10, 2, bush)
        p.rect(4, 16, 8, 1, bushDk)

        // Dark foliage
        p.setColor(bushDk)
        p.drawPixel(3, 5); p.drawPixel(5, 8); p.drawPixel(8, 7); p.drawPixel(10, 9)
        p.drawPixel(6, 12); p.drawPixel(9, 13); p.drawPixel(3, 11); p.drawPixel(12, 6)
        // Light foliage
        p.setColor(bushLt)
        p.drawPixel(5, 3); p.drawPixel(8, 3); p.drawPixel(6, 5); p.drawPixel(10, 5)
        p.drawPixel(4, 7); p.drawPixel(11, 7); p.drawPixel(7, 10); p.drawPixel(12, 10)
        p.setColor(bushBr)
        p.drawPixel(6, 2); p.drawPixel(9, 4); p.drawPixel(5, 6)

        // Bougainvillea flowers (bold!)
        p.setColor(pink)
        p.drawPixel(5, 3); p.drawPixel(9, 2); p.drawPixel(7, 4); p.drawPixel(4, 5)
        p.drawPixel(11, 4); p.drawPixel(6, 6); p.drawPixel(10, 6); p.drawPixel(3, 7)
        p.drawPixel(8, 8); p.drawPixel(12, 7); p.drawPixel(5, 10); p.drawPixel(10, 11)
        p.setColor(mag)
        p.drawPixel(8, 3); p.drawPixel(5, 4); p.drawPixel(11, 5); p.drawPixel(7, 7)
        p.drawPixel(4, 9); p.drawPixel(9, 10)
        p.setColor(pinkLt)
        p.drawPixel(6, 3); p.drawPixel(10, 3); p.drawPixel(9, 6); p.drawPixel(3, 8)

        // Shadow
        p.setColor(Color(0f, 0f, 0f, 0.18f))
        for (x in 3..12) { p.drawPixel(x, 17); p.drawPixel(x, 18) }
        for (x in 5..10) p.drawPixel(x, 19)

        return finish(p)
    }

    private fun createPotTexture(): Texture {
        val p = Pixmap(16, 18, Pixmap.Format.RGBA8888)
        p.setColor(Color.CLEAR); p.fill()

        val pot = Color(0.80f, 0.40f, 0.14f, 1f)
        val potLt = Color(0.90f, 0.52f, 0.24f, 1f)
        val potDk = Color(0.62f, 0.28f, 0.08f, 1f)
        val potRim = Color(0.86f, 0.46f, 0.18f, 1f)
        val leaf = Color(0.22f, 0.52f, 0.16f, 1f)
        val leafDk = Color(0.14f, 0.38f, 0.10f, 1f)
        val red = Color(0.94f, 0.20f, 0.16f, 1f)
        val redLt = Color(0.98f, 0.40f, 0.32f, 1f)
        val redDk = Color(0.74f, 0.12f, 0.08f, 1f)

        // Flowers
        p.setColor(leaf)
        p.drawPixel(5, 4); p.drawPixel(6, 3); p.drawPixel(8, 3); p.drawPixel(9, 4); p.drawPixel(10, 5)
        p.drawPixel(4, 5); p.drawPixel(11, 4); p.drawPixel(7, 5)
        p.setColor(leafDk); p.drawPixel(6, 5); p.drawPixel(9, 5); p.drawPixel(5, 6); p.drawPixel(10, 6)
        p.setColor(red)
        p.drawPixel(6, 1); p.drawPixel(7, 0); p.drawPixel(8, 1); p.drawPixel(5, 2)
        p.drawPixel(9, 2); p.drawPixel(10, 1); p.drawPixel(7, 3); p.drawPixel(4, 3); p.drawPixel(11, 2)
        p.setColor(redLt); p.drawPixel(7, 1); p.drawPixel(5, 1); p.drawPixel(10, 2)
        p.setColor(redDk); p.drawPixel(6, 2); p.drawPixel(8, 2); p.drawPixel(9, 3)

        // Rim
        p.rect(4, 7, 8, 1, potRim)
        // Pot body (tapers)
        for (row in 8..14) {
            val ins = (row - 8) / 2
            val l = 4 + ins; val r = 11 - ins
            for (x in l..r) {
                val c = when { x == l -> potLt; x == r -> potDk; else -> pot }
                p.setColor(c); p.drawPixel(x, row)
            }
        }
        p.rect(6, 15, 4, 1, potDk)
        // Shadow
        p.setColor(Color(0f, 0f, 0f, 0.16f))
        for (x in 5..10) { p.drawPixel(x, 16); p.drawPixel(x, 17) }

        return finish(p)
    }

    private fun createCrateTexture(): Texture {
        val p = Pixmap(16, 18, Pixmap.Format.RGBA8888)
        p.setColor(Color.CLEAR); p.fill()

        val wood = Color(0.56f, 0.40f, 0.18f, 1f)
        val woodLt = Color(0.68f, 0.52f, 0.28f, 1f)
        val woodDk = Color(0.40f, 0.26f, 0.10f, 1f)
        val iron = Color(0.36f, 0.34f, 0.32f, 1f)
        val ironDk = Color(0.26f, 0.24f, 0.22f, 1f)

        // Top face
        p.rect(2, 0, 12, 4, woodLt)
        p.setColor(wood); for (x in 2..13) p.drawPixel(x, 3)
        p.setColor(woodDk); p.drawPixel(2, 1); p.drawPixel(2, 2); p.drawPixel(2, 3)
        // Front face
        p.rect(2, 4, 12, 10, wood)
        p.setColor(woodDk); for (x in 2..13) { p.drawPixel(x, 7); p.drawPixel(x, 10) }
        p.setColor(woodLt); for (y in 4..13) p.drawPixel(2, y)
        p.setColor(woodDk); for (y in 4..13) p.drawPixel(13, y)
        // Iron bands
        p.setColor(iron); for (x in 2..13) { p.drawPixel(x, 5); p.drawPixel(x, 12) }
        p.setColor(ironDk); for (x in 2..13) { p.drawPixel(x, 6); p.drawPixel(x, 13) }
        // Corner brackets
        p.setColor(iron)
        p.drawPixel(3, 4); p.drawPixel(3, 5); p.drawPixel(12, 4); p.drawPixel(12, 5)
        p.drawPixel(3, 12); p.drawPixel(3, 13); p.drawPixel(12, 12); p.drawPixel(12, 13)
        // Shadow
        p.setColor(Color(0f, 0f, 0f, 0.18f))
        for (x in 2..13) { p.drawPixel(x, 14); p.drawPixel(x, 15) }
        for (x in 4..11) p.drawPixel(x, 16)

        return finish(p)
    }

    private fun createAwningTexture(): Texture {
        val p = Pixmap(16, 12, Pixmap.Format.RGBA8888)
        p.setColor(Color.CLEAR); p.fill()

        val red = Color(0.86f, 0.18f, 0.12f, 1f)
        val redDk = Color(0.66f, 0.12f, 0.08f, 1f)
        val white = Color(0.98f, 0.96f, 0.92f, 1f)
        val whiteDk = Color(0.86f, 0.84f, 0.80f, 1f)
        val pole = Color(0.36f, 0.32f, 0.28f, 1f)

        for (row in 0..8) {
            val str = if ((row / 2) % 2 == 0) red else white
            val strDk = if ((row / 2) % 2 == 0) redDk else whiteDk
            for (x in 1..14) p.px(x, row, if (row == 8) strDk else str)
        }
        // Scallop
        p.setColor(red); p.drawPixel(2, 9); p.drawPixel(3, 9); p.drawPixel(6, 9); p.drawPixel(7, 9)
        p.drawPixel(10, 9); p.drawPixel(11, 9)
        p.setColor(white); p.drawPixel(4, 9); p.drawPixel(5, 9); p.drawPixel(8, 9); p.drawPixel(9, 9)
        p.drawPixel(12, 9); p.drawPixel(13, 9)
        // Poles
        p.setColor(pole); for (y in 0..11) { p.drawPixel(1, y); p.drawPixel(14, y) }

        return finish(p)
    }

    private fun createLampTexture(): Texture {
        val p = Pixmap(16, 28, Pixmap.Format.RGBA8888)
        p.setColor(Color.CLEAR); p.fill()

        val iron = Color(0.20f, 0.18f, 0.16f, 1f)
        val ironLt = Color(0.32f, 0.30f, 0.26f, 1f)
        val ironDk = Color(0.12f, 0.10f, 0.08f, 1f)
        val glassW = Color(0.98f, 0.92f, 0.60f, 1f)
        val glassG = Color(1.0f, 0.96f, 0.76f, 1f)
        val glow = Color(1.0f, 0.94f, 0.58f, 0.30f)

        // Cap
        p.px(7, 0, ironLt); p.px(8, 0, ironLt)
        p.setColor(iron)
        p.drawPixel(6, 1); p.drawPixel(7, 1); p.drawPixel(8, 1); p.drawPixel(9, 1)
        p.drawPixel(5, 2); p.drawPixel(6, 2); p.drawPixel(9, 2); p.drawPixel(10, 2)
        // Lantern
        p.rect(5, 3, 6, 5, glassW)
        p.setColor(glassG); p.drawPixel(7, 4); p.drawPixel(8, 4); p.drawPixel(7, 5)
        p.setColor(iron)
        for (y in 3..7) { p.drawPixel(5, y); p.drawPixel(10, y) }
        p.drawPixel(6, 3); p.drawPixel(9, 3)
        // Glow
        p.setColor(glow)
        for (y in 3..7) { p.drawPixel(4, y); p.drawPixel(11, y) }
        // Bracket
        p.setColor(iron)
        p.drawPixel(6, 8); p.drawPixel(7, 8); p.drawPixel(8, 8); p.drawPixel(9, 8)
        // Post
        for (y in 9..25) { p.px(7, y, ironLt); p.px(8, y, iron) }
        p.setColor(ironDk); p.drawPixel(7, 9); p.drawPixel(8, 9)
        p.drawPixel(7, 15); p.drawPixel(8, 15)
        p.drawPixel(7, 21); p.drawPixel(8, 21)
        // Base
        p.setColor(iron)
        p.drawPixel(5, 25); p.drawPixel(6, 25); p.drawPixel(9, 25); p.drawPixel(10, 25)
        for (x in 5..10) p.drawPixel(x, 26)
        // Shadow
        p.setColor(Color(0f, 0f, 0f, 0.14f)); for (x in 5..10) p.drawPixel(x, 27)

        return finish(p)
    }

    override fun dispose() {
        textures.values.forEach { it.dispose() }
        if (::waterAlt.isInitialized) waterAlt.dispose()
        if (::waterDeepAlt.isInitialized) waterDeepAlt.dispose()
        if (::roofTallTexture.isInitialized) roofTallTexture.dispose()
        if (::wallFrontTexture.isInitialized) wallFrontTexture.dispose()
        if (::doorFrontTexture.isInitialized) doorFrontTexture.dispose()
        if (::palmTrunkTallTexture.isInitialized) palmTrunkTallTexture.dispose()
        if (::palmTopWideTexture.isInitialized) palmTopWideTexture.dispose()
        if (::fenceTallTexture.isInitialized) fenceTallTexture.dispose()
        if (::signTallTexture.isInitialized) signTallTexture.dispose()
        if (::bushFlowerTexture.isInitialized) bushFlowerTexture.dispose()
        if (::potTexture.isInitialized) potTexture.dispose()
        if (::crateTexture.isInitialized) crateTexture.dispose()
        if (::awningTexture.isInitialized) awningTexture.dispose()
        if (::lampTexture.isInitialized) lampTexture.dispose()
    }
}
