package com.lark.app.game.sprite

import com.badlogic.gdx.graphics.Color
import com.badlogic.gdx.graphics.Texture
import com.badlogic.gdx.utils.Disposable

/**
 * 7 distinct NPC character textures (32x48 pixels each, taller Pokemon-style proportions).
 * All NPCs face down (toward the player).
 * 0=Waiter (Carlos), 1=Barista (Maria), 2=Hotel Clerk (Eduardo),
 * 3=Market Vendor (Sofia), 4=Doctor (Dr. Rivera),
 * 5=Station Attendant (Miguel), 6=Friendly Local (Isabella)
 */
class NpcSprites : Disposable {

    val textures: Array<Texture>

    init {
        textures = arrayOf(
            createWaiter(),
            createBarista(),
            createHotelClerk(),
            createMarketVendor(),
            createDoctor(),
            createStationAttendant(),
            createFriendlyLocal()
        )
    }

    /** NPC 0: Waiter (Carlos) - black hair, white shirt, black vest, red bow tie, black pants */
    private fun createWaiter(): Texture {
        val palette = arrayOf(
            Color.CLEAR,                          // 0: transparent
            Color(0.12f, 0.08f, 0.05f, 1f),       // 1: black hair
            Color(0.85f, 0.68f, 0.52f, 1f),       // 2: skin
            Color(0.95f, 0.95f, 0.95f, 1f),       // 3: white shirt
            Color(0.1f, 0.1f, 0.1f, 1f),          // 4: black vest/pants
            Color(0.85f, 0.12f, 0.12f, 1f),       // 5: red bow tie
            Color(0.08f, 0.08f, 0.08f, 1f),       // 6: black shoes
            Color(0.05f, 0.05f, 0.05f, 1f),       // 7: eyes
        )
        // 16x24 - facing down
        val pixels = intArrayOf(
            // Row 0-2: top of hair
            0,0,0,0,0,0,1,1,1,1,0,0,0,0,0,0,
            0,0,0,0,0,1,1,1,1,1,1,0,0,0,0,0,
            0,0,0,0,1,1,1,1,1,1,1,1,0,0,0,0,
            // Row 3-5: hair
            0,0,0,0,1,1,1,1,1,1,1,1,0,0,0,0,
            0,0,0,0,1,1,1,1,1,1,1,1,0,0,0,0,
            0,0,0,0,1,1,1,1,1,1,1,1,0,0,0,0,
            // Row 6-8: face (eyes at row 7)
            0,0,0,0,2,2,2,2,2,2,2,2,0,0,0,0,
            0,0,0,0,2,7,2,2,2,2,7,2,0,0,0,0,
            0,0,0,0,2,2,2,2,2,2,2,2,0,0,0,0,
            // Row 9: neck with bow tie
            0,0,0,0,0,2,2,5,5,2,2,0,0,0,0,0,
            // Row 10-14: torso (vest over white shirt)
            0,0,0,0,4,4,3,3,3,3,4,4,0,0,0,0,
            0,0,0,4,4,4,3,3,3,3,4,4,4,0,0,0,
            0,0,0,4,4,3,3,3,3,3,3,4,4,0,0,0,
            0,0,0,2,4,3,3,3,3,3,3,4,2,0,0,0,
            0,0,0,0,4,4,3,3,3,3,4,4,0,0,0,0,
            // Row 15-17: pants
            0,0,0,0,4,4,4,4,4,4,4,4,0,0,0,0,
            0,0,0,0,4,4,4,0,0,4,4,4,0,0,0,0,
            0,0,0,0,4,4,4,0,0,4,4,4,0,0,0,0,
            // Row 18-20: legs
            0,0,0,0,0,4,4,0,0,4,4,0,0,0,0,0,
            0,0,0,0,0,4,4,0,0,4,4,0,0,0,0,0,
            0,0,0,0,0,2,2,0,0,2,2,0,0,0,0,0,
            // Row 21-23: feet
            0,0,0,0,0,6,6,0,0,6,6,0,0,0,0,0,
            0,0,0,0,6,6,6,0,0,6,6,6,0,0,0,0,
            0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,
        )
        return PixelArt.createTexture(32, 48, upscale2x(pixels, 16, 24), palette)
    }

    /** NPC 1: Barista (Maria) - brown ponytail, green apron, white tee, blue jeans, brown shoes */
    private fun createBarista(): Texture {
        val palette = arrayOf(
            Color.CLEAR,                          // 0: transparent
            Color(0.4f, 0.22f, 0.1f, 1f),        // 1: brown hair
            Color(0.82f, 0.65f, 0.5f, 1f),        // 2: skin
            Color(0.2f, 0.55f, 0.2f, 1f),         // 3: green apron
            Color(0.92f, 0.9f, 0.88f, 1f),        // 4: white t-shirt
            Color(0.25f, 0.35f, 0.6f, 1f),        // 5: blue jeans
            Color(0.4f, 0.25f, 0.12f, 1f),        // 6: brown shoes
            Color(0.05f, 0.05f, 0.05f, 1f),       // 7: eyes
        )
        val pixels = intArrayOf(
            // Row 0-2: top of hair (ponytail extends right)
            0,0,0,0,0,0,1,1,1,1,0,0,0,0,0,0,
            0,0,0,0,0,1,1,1,1,1,1,1,0,0,0,0,
            0,0,0,0,1,1,1,1,1,1,1,1,1,0,0,0,
            // Row 3-5: hair with ponytail on right
            0,0,0,0,1,1,1,1,1,1,1,1,1,0,0,0,
            0,0,0,0,1,1,1,1,1,1,1,1,1,1,0,0,
            0,0,0,0,1,1,1,1,1,1,1,1,0,1,0,0,
            // Row 6-8: face (eyes at row 7)
            0,0,0,0,2,2,2,2,2,2,2,2,0,1,0,0,
            0,0,0,0,2,7,2,2,2,2,7,2,0,1,0,0,
            0,0,0,0,2,2,2,2,2,2,2,2,0,1,0,0,
            // Row 9: neck
            0,0,0,0,0,2,2,2,2,2,2,0,0,0,0,0,
            // Row 10-14: torso (green apron over white tee)
            0,0,0,0,3,3,3,3,3,3,3,3,0,0,0,0,
            0,0,0,2,3,3,4,4,4,4,3,3,2,0,0,0,
            0,0,0,2,3,3,4,4,4,4,3,3,2,0,0,0,
            0,0,0,0,3,3,4,4,4,4,3,3,0,0,0,0,
            0,0,0,0,3,3,3,3,3,3,3,3,0,0,0,0,
            // Row 15-17: jeans
            0,0,0,0,5,5,5,5,5,5,5,5,0,0,0,0,
            0,0,0,0,5,5,5,0,0,5,5,5,0,0,0,0,
            0,0,0,0,5,5,5,0,0,5,5,5,0,0,0,0,
            // Row 18-20: legs
            0,0,0,0,0,5,5,0,0,5,5,0,0,0,0,0,
            0,0,0,0,0,5,5,0,0,5,5,0,0,0,0,0,
            0,0,0,0,0,2,2,0,0,2,2,0,0,0,0,0,
            // Row 21-23: feet (brown shoes)
            0,0,0,0,0,6,6,0,0,6,6,0,0,0,0,0,
            0,0,0,0,6,6,6,0,0,6,6,6,0,0,0,0,
            0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,
        )
        return PixelArt.createTexture(32, 48, upscale2x(pixels, 16, 24), palette)
    }

    /** NPC 2: Hotel Clerk (Eduardo) - slicked black hair, navy blazer, white shirt, red tie */
    private fun createHotelClerk(): Texture {
        val palette = arrayOf(
            Color.CLEAR,                          // 0: transparent
            Color(0.08f, 0.06f, 0.04f, 1f),       // 1: slicked black hair
            Color(0.88f, 0.72f, 0.58f, 1f),       // 2: skin
            Color(0.12f, 0.15f, 0.4f, 1f),        // 3: navy blue blazer
            Color(0.92f, 0.92f, 0.92f, 1f),       // 4: white shirt
            Color(0.8f, 0.15f, 0.15f, 1f),        // 5: red tie
            Color(0.22f, 0.22f, 0.25f, 1f),       // 6: dark gray pants
            Color(0.05f, 0.05f, 0.05f, 1f),       // 7: eyes/black shoes
        )
        val pixels = intArrayOf(
            // Row 0-2: top of slicked hair
            0,0,0,0,0,0,1,1,1,1,0,0,0,0,0,0,
            0,0,0,0,0,1,1,1,1,1,1,0,0,0,0,0,
            0,0,0,0,1,1,1,1,1,1,1,1,0,0,0,0,
            // Row 3-5: hair (slicked back, neat)
            0,0,0,0,1,1,1,1,1,1,1,1,0,0,0,0,
            0,0,0,0,1,1,1,1,1,1,1,1,0,0,0,0,
            0,0,0,0,1,1,1,1,1,1,1,1,0,0,0,0,
            // Row 6-8: face (eyes at row 7)
            0,0,0,0,2,2,2,2,2,2,2,2,0,0,0,0,
            0,0,0,0,2,7,2,2,2,2,7,2,0,0,0,0,
            0,0,0,0,2,2,2,2,2,2,2,2,0,0,0,0,
            // Row 9: neck with tie starting
            0,0,0,0,0,2,2,5,5,2,2,0,0,0,0,0,
            // Row 10-14: torso (navy blazer, white shirt, red tie)
            0,0,0,0,3,4,4,5,5,4,4,3,0,0,0,0,
            0,0,0,3,3,4,4,5,5,4,4,3,3,0,0,0,
            0,0,0,3,3,3,4,5,5,4,3,3,3,0,0,0,
            0,0,0,2,3,3,4,5,5,4,3,3,2,0,0,0,
            0,0,0,0,3,3,3,3,3,3,3,3,0,0,0,0,
            // Row 15-17: dark gray pants
            0,0,0,0,6,6,6,6,6,6,6,6,0,0,0,0,
            0,0,0,0,6,6,6,0,0,6,6,6,0,0,0,0,
            0,0,0,0,6,6,6,0,0,6,6,6,0,0,0,0,
            // Row 18-20: legs
            0,0,0,0,0,6,6,0,0,6,6,0,0,0,0,0,
            0,0,0,0,0,6,6,0,0,6,6,0,0,0,0,0,
            0,0,0,0,0,2,2,0,0,2,2,0,0,0,0,0,
            // Row 21-23: feet (black shoes)
            0,0,0,0,0,7,7,0,0,7,7,0,0,0,0,0,
            0,0,0,0,7,7,7,0,0,7,7,7,0,0,0,0,
            0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,
        )
        return PixelArt.createTexture(32, 48, upscale2x(pixels, 16, 24), palette)
    }

    /** NPC 3: Market Vendor (Sofia) - dark brown hair with headband, orange apron, cream shirt */
    private fun createMarketVendor(): Texture {
        val palette = arrayOf(
            Color.CLEAR,                          // 0: transparent
            Color(0.28f, 0.15f, 0.06f, 1f),       // 1: dark brown hair
            Color(0.75f, 0.58f, 0.42f, 1f),       // 2: skin
            Color(0.92f, 0.55f, 0.1f, 1f),        // 3: orange apron
            Color(0.92f, 0.88f, 0.78f, 1f),       // 4: cream shirt
            Color(0.48f, 0.35f, 0.22f, 1f),       // 5: brown pants
            Color(0.65f, 0.5f, 0.3f, 1f),         // 6: sandals (tan)
            Color(0.05f, 0.05f, 0.05f, 1f),       // 7: eyes/headband
        )
        val pixels = intArrayOf(
            // Row 0-2: top of hair
            0,0,0,0,0,0,1,1,1,1,0,0,0,0,0,0,
            0,0,0,0,0,1,1,1,1,1,1,0,0,0,0,0,
            0,0,0,0,1,1,1,1,1,1,1,1,0,0,0,0,
            // Row 3-5: hair with headband (row 3)
            0,0,0,0,7,7,7,7,7,7,7,7,0,0,0,0,
            0,0,0,0,1,1,1,1,1,1,1,1,0,0,0,0,
            0,0,0,0,1,1,1,1,1,1,1,1,0,0,0,0,
            // Row 6-8: face (eyes at row 7)
            0,0,0,0,2,2,2,2,2,2,2,2,0,0,0,0,
            0,0,0,0,2,7,2,2,2,2,7,2,0,0,0,0,
            0,0,0,0,2,2,2,2,2,2,2,2,0,0,0,0,
            // Row 9: neck
            0,0,0,0,0,2,2,2,2,2,2,0,0,0,0,0,
            // Row 10-14: torso (orange apron over cream shirt)
            0,0,0,0,3,3,4,4,4,4,3,3,0,0,0,0,
            0,0,0,2,3,3,4,4,4,4,3,3,2,0,0,0,
            0,0,0,2,3,3,4,4,4,4,3,3,2,0,0,0,
            0,0,0,0,3,3,4,4,4,4,3,3,0,0,0,0,
            0,0,0,0,3,3,3,3,3,3,3,3,0,0,0,0,
            // Row 15-17: brown pants
            0,0,0,0,5,5,5,5,5,5,5,5,0,0,0,0,
            0,0,0,0,5,5,5,0,0,5,5,5,0,0,0,0,
            0,0,0,0,5,5,5,0,0,5,5,5,0,0,0,0,
            // Row 18-20: legs
            0,0,0,0,0,5,5,0,0,5,5,0,0,0,0,0,
            0,0,0,0,0,5,5,0,0,5,5,0,0,0,0,0,
            0,0,0,0,0,2,2,0,0,2,2,0,0,0,0,0,
            // Row 21-23: feet (sandals)
            0,0,0,0,0,6,6,0,0,6,6,0,0,0,0,0,
            0,0,0,0,6,6,6,0,0,6,6,6,0,0,0,0,
            0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,
        )
        return PixelArt.createTexture(32, 48, upscale2x(pixels, 16, 24), palette)
    }

    /** NPC 4: Doctor (Dr. Rivera) - gray hair, white lab coat, light blue shirt, red cross */
    private fun createDoctor(): Texture {
        val palette = arrayOf(
            Color.CLEAR,                          // 0: transparent
            Color(0.6f, 0.58f, 0.55f, 1f),        // 1: gray/white hair
            Color(0.88f, 0.72f, 0.58f, 1f),       // 2: skin
            Color(0.95f, 0.95f, 0.95f, 1f),       // 3: white lab coat
            Color(0.55f, 0.72f, 0.88f, 1f),       // 4: light blue shirt
            Color(0.25f, 0.25f, 0.3f, 1f),        // 5: dark pants
            Color(0.88f, 0.15f, 0.15f, 1f),       // 6: red cross
            Color(0.05f, 0.05f, 0.05f, 1f),       // 7: eyes/shoes
        )
        val pixels = intArrayOf(
            // Row 0-2: top of gray hair
            0,0,0,0,0,0,1,1,1,1,0,0,0,0,0,0,
            0,0,0,0,0,1,1,1,1,1,1,0,0,0,0,0,
            0,0,0,0,1,1,1,1,1,1,1,1,0,0,0,0,
            // Row 3-5: hair
            0,0,0,0,1,1,1,1,1,1,1,1,0,0,0,0,
            0,0,0,0,1,1,1,1,1,1,1,1,0,0,0,0,
            0,0,0,0,1,1,1,1,1,1,1,1,0,0,0,0,
            // Row 6-8: face (eyes at row 7)
            0,0,0,0,2,2,2,2,2,2,2,2,0,0,0,0,
            0,0,0,0,2,7,2,2,2,2,7,2,0,0,0,0,
            0,0,0,0,2,2,2,2,2,2,2,2,0,0,0,0,
            // Row 9: neck
            0,0,0,0,0,2,2,2,2,2,2,0,0,0,0,0,
            // Row 10-14: torso (lab coat over blue shirt, red cross on left chest)
            0,0,0,3,3,4,4,4,4,4,4,3,3,0,0,0,
            0,0,3,3,3,4,6,4,4,4,4,3,3,3,0,0,
            0,0,3,3,3,6,6,6,4,4,4,3,3,3,0,0,
            0,0,3,2,3,4,6,4,4,4,4,3,2,3,0,0,
            0,0,3,3,3,3,3,3,3,3,3,3,3,3,0,0,
            // Row 15-17: coat continues + pants visible
            0,0,0,3,3,5,5,5,5,5,5,3,3,0,0,0,
            0,0,0,3,3,5,5,0,0,5,5,3,3,0,0,0,
            0,0,0,0,0,5,5,0,0,5,5,0,0,0,0,0,
            // Row 18-20: legs
            0,0,0,0,0,5,5,0,0,5,5,0,0,0,0,0,
            0,0,0,0,0,5,5,0,0,5,5,0,0,0,0,0,
            0,0,0,0,0,2,2,0,0,2,2,0,0,0,0,0,
            // Row 21-23: feet (black shoes)
            0,0,0,0,0,7,7,0,0,7,7,0,0,0,0,0,
            0,0,0,0,7,7,7,0,0,7,7,7,0,0,0,0,
            0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,
        )
        return PixelArt.createTexture(32, 48, upscale2x(pixels, 16, 24), palette)
    }

    /** NPC 5: Station Attendant (Miguel) - black hair, blue cap, blue uniform, gold buttons */
    private fun createStationAttendant(): Texture {
        val palette = arrayOf(
            Color.CLEAR,                          // 0: transparent
            Color(0.12f, 0.08f, 0.05f, 1f),       // 1: black hair
            Color(0.82f, 0.65f, 0.48f, 1f),       // 2: skin
            Color(0.15f, 0.28f, 0.55f, 1f),       // 3: blue uniform jacket
            Color(0.88f, 0.75f, 0.18f, 1f),       // 4: gold buttons/badge
            Color(0.1f, 0.18f, 0.38f, 1f),        // 5: dark blue pants
            Color(0.08f, 0.08f, 0.08f, 1f),       // 6: black shoes
            Color(0.05f, 0.05f, 0.05f, 1f),       // 7: eyes
        )
        val pixels = intArrayOf(
            // Row 0-2: top of cap (blue cap with brim)
            0,0,0,0,0,3,3,3,3,3,3,0,0,0,0,0,
            0,0,0,3,3,3,3,3,3,3,3,3,3,0,0,0,
            0,0,0,0,3,3,3,3,3,3,3,3,0,0,0,0,
            // Row 3-5: hair under cap
            0,0,0,0,1,1,1,1,1,1,1,1,0,0,0,0,
            0,0,0,0,1,1,1,1,1,1,1,1,0,0,0,0,
            0,0,0,0,1,1,1,1,1,1,1,1,0,0,0,0,
            // Row 6-8: face (eyes at row 7)
            0,0,0,0,2,2,2,2,2,2,2,2,0,0,0,0,
            0,0,0,0,2,7,2,2,2,2,7,2,0,0,0,0,
            0,0,0,0,2,2,2,2,2,2,2,2,0,0,0,0,
            // Row 9: neck
            0,0,0,0,0,2,2,2,2,2,2,0,0,0,0,0,
            // Row 10-14: torso (blue uniform jacket with gold buttons)
            0,0,0,0,3,3,3,4,3,3,3,3,0,0,0,0,
            0,0,0,2,3,3,3,4,3,3,3,3,2,0,0,0,
            0,0,0,2,3,3,3,4,3,3,3,3,2,0,0,0,
            0,0,0,0,3,3,3,4,3,3,3,3,0,0,0,0,
            0,0,0,0,3,3,3,3,3,3,3,3,0,0,0,0,
            // Row 15-17: dark blue pants
            0,0,0,0,5,5,5,5,5,5,5,5,0,0,0,0,
            0,0,0,0,5,5,5,0,0,5,5,5,0,0,0,0,
            0,0,0,0,5,5,5,0,0,5,5,5,0,0,0,0,
            // Row 18-20: legs
            0,0,0,0,0,5,5,0,0,5,5,0,0,0,0,0,
            0,0,0,0,0,5,5,0,0,5,5,0,0,0,0,0,
            0,0,0,0,0,2,2,0,0,2,2,0,0,0,0,0,
            // Row 21-23: feet (black shoes)
            0,0,0,0,0,6,6,0,0,6,6,0,0,0,0,0,
            0,0,0,0,6,6,6,0,0,6,6,6,0,0,0,0,
            0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,
        )
        return PixelArt.createTexture(32, 48, upscale2x(pixels, 16, 24), palette)
    }

    /** NPC 6: Friendly Local (Isabella) - long dark hair, sun hat, pink sundress, tan skin */
    private fun createFriendlyLocal(): Texture {
        val palette = arrayOf(
            Color.CLEAR,                          // 0: transparent
            Color(0.22f, 0.12f, 0.05f, 1f),       // 1: long dark hair
            Color(0.82f, 0.65f, 0.48f, 1f),       // 2: skin (tan)
            Color(0.95f, 0.45f, 0.55f, 1f),       // 3: pink sundress
            Color(0.98f, 0.6f, 0.68f, 1f),        // 4: lighter pink (dress accent)
            Color(0.92f, 0.85f, 0.35f, 1f),       // 5: yellow/straw sun hat
            Color(0.55f, 0.38f, 0.2f, 1f),        // 6: brown sandals
            Color(0.05f, 0.05f, 0.05f, 1f),       // 7: eyes
        )
        val pixels = intArrayOf(
            // Row 0-2: wide sun hat (straw/yellow)
            0,0,0,5,5,5,5,5,5,5,5,5,5,0,0,0,
            0,0,5,5,5,5,5,5,5,5,5,5,5,5,0,0,
            0,5,5,5,5,5,5,5,5,5,5,5,5,5,5,0,
            // Row 3-5: hair beneath hat, long and dark
            0,0,0,0,1,1,1,1,1,1,1,1,0,0,0,0,
            0,0,0,0,1,1,1,1,1,1,1,1,0,0,0,0,
            0,0,0,0,1,1,1,1,1,1,1,1,0,0,0,0,
            // Row 6-8: face (eyes at row 7), hair on sides
            0,0,0,1,2,2,2,2,2,2,2,2,1,0,0,0,
            0,0,0,1,2,7,2,2,2,2,7,2,1,0,0,0,
            0,0,0,1,2,2,2,2,2,2,2,2,1,0,0,0,
            // Row 9: neck, hair on sides
            0,0,0,1,0,2,2,2,2,2,2,0,1,0,0,0,
            // Row 10-14: torso (pink sundress)
            0,0,0,1,3,3,3,4,4,3,3,3,1,0,0,0,
            0,0,0,1,3,3,4,4,4,4,3,3,1,0,0,0,
            0,0,0,0,3,3,4,4,4,4,3,3,0,0,0,0,
            0,0,0,2,3,3,3,4,4,3,3,3,2,0,0,0,
            0,0,0,0,3,3,3,3,3,3,3,3,0,0,0,0,
            // Row 15-17: dress skirt (flares out)
            0,0,0,0,3,3,4,3,3,4,3,3,0,0,0,0,
            0,0,0,3,3,3,3,4,4,3,3,3,3,0,0,0,
            0,0,0,3,3,4,3,3,3,3,4,3,3,0,0,0,
            // Row 18-20: legs below dress
            0,0,0,0,0,2,2,0,0,2,2,0,0,0,0,0,
            0,0,0,0,0,2,2,0,0,2,2,0,0,0,0,0,
            0,0,0,0,0,2,2,0,0,2,2,0,0,0,0,0,
            // Row 21-23: feet (brown sandals)
            0,0,0,0,0,6,6,0,0,6,6,0,0,0,0,0,
            0,0,0,0,6,6,6,0,0,6,6,6,0,0,0,0,
            0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,
        )
        return PixelArt.createTexture(32, 48, upscale2x(pixels, 16, 24), palette)
    }

    private fun upscale2x(src: IntArray, srcW: Int, srcH: Int): IntArray {
        val dstW = srcW * 2
        val dstH = srcH * 2
        val dst = IntArray(dstW * dstH)
        for (y in 0 until srcH) {
            for (x in 0 until srcW) {
                val v = src[y * srcW + x]
                dst[(y * 2) * dstW + (x * 2)] = v
                dst[(y * 2) * dstW + (x * 2 + 1)] = v
                dst[(y * 2 + 1) * dstW + (x * 2)] = v
                dst[(y * 2 + 1) * dstW + (x * 2 + 1)] = v
            }
        }
        return dst
    }

    override fun dispose() {
        textures.forEach { it.dispose() }
    }
}
