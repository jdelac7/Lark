package com.lark.app.game.npc

import com.lark.app.game.sprite.Direction

data class Npc(
    val id: String,
    val scenarioId: String,
    val displayName: String,
    val tileX: Int,
    val tileY: Int,
    val facing: Direction,
    val spriteId: Int, // Index into NpcSprites
    var completed: Boolean = false
)
