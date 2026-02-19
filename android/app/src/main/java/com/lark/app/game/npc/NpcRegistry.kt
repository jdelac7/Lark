package com.lark.app.game.npc

import com.lark.app.game.sprite.Direction

object NpcRegistry {

    fun getNpcs(languageCode: String): List<Npc> = listOf(
        Npc(
            id = "waiter",
            scenarioId = "restaurant",
            displayName = "Carlos",
            tileX = 5,
            tileY = 15,
            facing = Direction.DOWN,
            spriteId = 0
        ),
        Npc(
            id = "barista",
            scenarioId = "cafe",
            displayName = "Maria",
            tileX = 14,
            tileY = 15,
            facing = Direction.DOWN,
            spriteId = 1
        ),
        Npc(
            id = "hotel_clerk",
            scenarioId = "hotel",
            displayName = "Eduardo",
            tileX = 24,
            tileY = 15,
            facing = Direction.DOWN,
            spriteId = 2
        ),
        Npc(
            id = "market_vendor",
            scenarioId = "market",
            displayName = "Sofia",
            tileX = 30,
            tileY = 15,
            facing = Direction.DOWN,
            spriteId = 3
        ),
        Npc(
            id = "doctor",
            scenarioId = "doctor",
            displayName = "Dr. Rivera",
            tileX = 6,
            tileY = 5,
            facing = Direction.DOWN,
            spriteId = 4
        ),
        Npc(
            id = "station_attendant",
            scenarioId = "train_station",
            displayName = "Miguel",
            tileX = 33,
            tileY = 5,
            facing = Direction.DOWN,
            spriteId = 5
        ),
        Npc(
            id = "friendly_local",
            scenarioId = "directions",
            displayName = "Isabella",
            tileX = 19,
            tileY = 7,
            facing = Direction.DOWN,
            spriteId = 6
        )
    )
}
