package com.lark.app.game.map

/**
 * 40x30 tile map of the beach town.
 * Buildings are 2 rows (roof + wall/door) for seamless 3/4 perspective.
 */
object BeachTownMap {

    private const val W = 40
    private const val H = 30

    private const val SD = 0  // SAND
    private const val WA = 1  // WATER
    private const val WD = 2  // WATER_DEEP
    private const val GR = 3  // GRASS
    private const val PS = 4  // PATH_STONE
    private const val PW = 5  // PATH_WOOD
    private const val WL = 6  // WALL
    private const val DR = 7  // DOOR
    private const val RF = 8  // ROOF
    private const val PT = 9  // PALM_TRUNK
    private const val PO = 10 // PALM_TOP
    private const val FL = 11 // FLOWERS
    private const val FE = 12 // FENCE
    private const val SI = 13 // SIGN
    private const val DK = 14 // DOCK
    private const val BF = 15 // BUSH_FLOWER
    private const val CP = 16 // POT
    private const val CR = 17 // CRATE
    private const val AW = 18 // AWNING
    private const val LA = 19 // LAMP

    fun create(): IntArray {
        val map = IntArray(W * H) { GR }

        // === OCEAN (rows 25-29) ===
        for (y in 25..29) for (x in 0 until W) set(map, x, y, WD)

        // === SHALLOW WATER (rows 23-24) ===
        for (y in 23..24) for (x in 0 until W) set(map, x, y, WA)

        // Dock (x=33..37, y=23..24)
        for (y in 23..24) for (x in 33..37) set(map, x, y, DK)

        // === BEACH (rows 21-22) ===
        for (y in 21..22) for (x in 0 until W) set(map, x, y, SD)

        // Palm trees
        placePalm(map, 3, 21)
        placePalm(map, 10, 22)
        placePalm(map, 18, 21)
        placePalm(map, 27, 22)
        placePalm(map, 35, 21)

        // === BOARDWALK (row 20) ===
        for (x in 0 until W) set(map, x, 20, PW)

        // === TOWN (rows 0-19) ===

        // Main path (row 19)
        for (x in 0 until W) set(map, x, 19, PS)

        // Vertical paths
        for (y in 1..19) set(map, 10, y, PS)
        for (y in 1..19) set(map, 20, y, PS)
        for (y in 1..19) set(map, 32, y, PS)

        // Cross path (row 11)
        for (x in 10..32) set(map, x, 11, PS)

        // === BUILDINGS (2 rows: roof on top, wall/door on bottom) ===

        // Upper row (wall at y=16, roof at y=17) — NPCs at y=15
        placeBuilding(map, 3, 16, 6, 5)    // Restaurant, door at x=5
        placeBuilding(map, 12, 16, 5, 14)   // Cafe, door at x=14
        placeBuilding(map, 22, 16, 5, 24)   // Hotel, door at x=24
        placeBuilding(map, 28, 16, 6, 30)   // Market, door at x=30

        // Lower row (wall at y=6, roof at y=7) — NPCs at y=5
        placeBuilding(map, 4, 6, 6, 6)     // Doctor, door at x=6
        placeBuilding(map, 30, 6, 6, 33)    // Train Station, door at x=33

        // Town Square (cols 15-24, rows 5-9)
        for (y in 5..9) for (x in 15..24) set(map, x, y, PS)
        set(map, 17, 6, FL); set(map, 19, 7, FL); set(map, 22, 6, FL)
        set(map, 18, 8, FL); set(map, 21, 7, FL)
        for (y in 9..11) set(map, 20, y, PS)

        // === DECORATIONS ===

        // Side fences (only sides + top, not bottom where NPCs are)
        placeSideFences(map, 3, 16, 6, 2)
        placeSideFences(map, 12, 16, 5, 2)
        placeSideFences(map, 22, 16, 5, 2)
        placeSideFences(map, 28, 16, 6, 2)
        placeSideFences(map, 4, 6, 6, 2)
        placeSideFences(map, 30, 6, 6, 2)

        // Signs (south of buildings, near doors)
        set(map, 3, 15, SI);  set(map, 12, 15, SI)
        set(map, 22, 15, SI); set(map, 28, 15, SI)
        set(map, 4, 5, SI);   set(map, 31, 5, SI)
        set(map, 16, 4, SI)   // Town square sign

        // Flower patches
        set(map, 1, 18, FL); set(map, 38, 18, FL)
        set(map, 15, 2, FL); set(map, 25, 2, FL)
        set(map, 1, 9, FL);  set(map, 38, 9, FL)
        set(map, 9, 14, FL); set(map, 34, 14, FL)

        // === MEDITERRANEAN DECORATIONS ===

        // Flowering bushes (lush bougainvillea at building corners/edges)
        set(map, 1, 14, BF);  set(map, 1, 13, BF)    // Left of town
        set(map, 11, 14, BF); set(map, 17, 14, BF)    // Between upper buildings
        set(map, 18, 14, BF); set(map, 21, 14, BF)
        set(map, 27, 14, BF); set(map, 36, 13, BF)    // Right side
        set(map, 37, 14, BF)
        // Replace 4 edge flowers with flowering bushes
        set(map, 1, 9, BF);   set(map, 38, 9, BF)
        set(map, 9, 14, BF);  set(map, 34, 14, BF)

        // Flower pots (flanking every building entrance)
        set(map, 4, 15, CP);  set(map, 6, 15, CP)     // Restaurant
        set(map, 13, 15, CP); set(map, 15, 15, CP)     // Cafe
        set(map, 23, 15, CP); set(map, 25, 15, CP)     // Hotel
        set(map, 29, 15, CP); set(map, 31, 15, CP)     // Market
        set(map, 5, 5, CP);   set(map, 7, 5, CP)       // Doctor
        set(map, 34, 5, CP);  set(map, 35, 5, CP)      // Train Station

        // Crates (near market/dock area)
        set(map, 35, 15, CR); set(map, 36, 15, CR)
        set(map, 37, 18, CR); set(map, 38, 14, CR)

        // Street lamps (along main path + vertical paths)
        set(map, 0, 18, LA);  set(map, 18, 18, LA)     // Main road
        set(map, 35, 18, LA); set(map, 39, 18, LA)
        set(map, 9, 3, LA);   set(map, 11, 13, LA)     // Near x=10 path
        set(map, 19, 13, LA); set(map, 21, 3, LA)      // Near x=20 path
        set(map, 31, 13, LA)                            // Near x=32 path

        return map
    }

    private fun set(map: IntArray, x: Int, y: Int, tile: Int) {
        if (x in 0 until W && y in 0 until H) map[y * W + x] = tile
    }

    private fun get(map: IntArray, x: Int, y: Int): Int {
        if (x !in 0 until W || y !in 0 until H) return -1
        return map[y * W + x]
    }

    private fun placePalm(map: IntArray, x: Int, y: Int) {
        set(map, x, y, PT)
        set(map, x, y + 1, PO)
    }

    /** Place a 2-row building: wall/door at startY, roof at startY+1. */
    private fun placeBuilding(map: IntArray, startX: Int, startY: Int, width: Int, doorX: Int) {
        // Roof on top row
        for (x in startX until startX + width) set(map, x, startY + 1, RF)
        // Walls on bottom row
        for (x in startX until startX + width) set(map, x, startY, WL)
        // Door
        set(map, doorX, startY, DR)
    }

    /** Place fences on sides and top of building (not bottom, NPCs walk there). */
    private fun placeSideFences(map: IntArray, startX: Int, startY: Int, width: Int, height: Int) {
        val left = startX - 1
        val right = startX + width
        val top = startY + height
        // Top fence
        for (x in left..right) {
            if (get(map, x, top) == GR) set(map, x, top, FE)
        }
        // Side fences
        for (y in startY - 1..top) {
            if (get(map, left, y) == GR) set(map, left, y, FE)
            if (get(map, right, y) == GR) set(map, right, y, FE)
        }
    }
}
