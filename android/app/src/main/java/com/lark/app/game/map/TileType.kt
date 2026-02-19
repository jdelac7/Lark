package com.lark.app.game.map

enum class TileType(val walkable: Boolean) {
    SAND(true),
    WATER(false),
    WATER_DEEP(false),
    GRASS(true),
    PATH_STONE(true),
    PATH_WOOD(true),
    WALL(false),
    DOOR(true),
    ROOF(false),
    PALM_TRUNK(false),
    PALM_TOP(false),  // Rendered above, walkable underneath in some cases
    FLOWERS(true),
    FENCE(false),
    SIGN(false),
    DOCK(true),
    BUSH_FLOWER(false),  // 15 - Lush flowering bush
    POT(false),          // 16 - Terracotta flower pot
    CRATE(false),        // 17 - Wooden crate/barrel
    AWNING(false),       // 18 - Striped canopy
    LAMP(false),         // 19 - Iron street lamp
    CLIFF(false),        // 20 - Rocky cliff face
    WATER_LIGHT(false);  // 21 - Shallow turquoise water

    companion object {
        fun fromInt(value: Int): TileType = entries.getOrElse(value) { GRASS }
    }
}
