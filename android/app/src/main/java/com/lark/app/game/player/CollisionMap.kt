package com.lark.app.game.player

import com.lark.app.game.map.TileType
import com.lark.app.game.npc.Npc

class CollisionMap(
    mapData: IntArray,
    private val width: Int,
    private val height: Int,
    npcs: List<Npc>
) {
    private val walkable = BooleanArray(width * height)

    init {
        // Initialize from tile data
        for (i in mapData.indices) {
            walkable[i] = TileType.fromInt(mapData[i]).walkable
        }
        // Block NPC positions
        for (npc in npcs) {
            val idx = npc.tileY * width + npc.tileX
            if (idx in walkable.indices) {
                walkable[idx] = false
            }
        }
    }

    fun isWalkable(x: Int, y: Int): Boolean {
        if (x < 0 || x >= width || y < 0 || y >= height) return false
        return walkable[y * width + x]
    }
}
