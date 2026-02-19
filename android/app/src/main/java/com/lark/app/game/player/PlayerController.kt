package com.lark.app.game.player

import com.lark.app.game.sprite.Direction

class PlayerController(
    startX: Int,
    startY: Int,
    private val collisionMap: CollisionMap
) {
    companion object {
        const val TILE_SIZE = 16f
        const val MOVE_DURATION = 0.2f // seconds per tile
    }

    var tileX: Int = startX
        private set
    var tileY: Int = startY
        private set
    var facing: Direction = Direction.DOWN
        private set

    // Smooth movement interpolation
    var worldX: Float = startX * TILE_SIZE
        private set
    var worldY: Float = startY * TILE_SIZE
        private set

    var isMoving: Boolean = false
        private set
    var walkFrame: Boolean = false // Alternates for walk animation
        private set

    var locked: Boolean = false // True when in dialog

    private var moveTimer: Float = 0f
    private var startWorldX: Float = 0f
    private var startWorldY: Float = 0f
    private var targetWorldX: Float = 0f
    private var targetWorldY: Float = 0f

    fun tryMove(direction: Direction) {
        if (isMoving || locked) return

        facing = direction
        val newX = tileX + direction.dx
        val newY = tileY + direction.dy

        if (!collisionMap.isWalkable(newX, newY)) return

        // Start movement
        isMoving = true
        moveTimer = 0f
        startWorldX = worldX
        startWorldY = worldY
        tileX = newX
        tileY = newY
        targetWorldX = newX * TILE_SIZE
        targetWorldY = newY * TILE_SIZE
        walkFrame = !walkFrame
    }

    fun update(delta: Float) {
        if (!isMoving) return

        moveTimer += delta
        val progress = (moveTimer / MOVE_DURATION).coerceIn(0f, 1f)

        worldX = startWorldX + (targetWorldX - startWorldX) * progress
        worldY = startWorldY + (targetWorldY - startWorldY) * progress

        if (progress >= 1f) {
            worldX = targetWorldX
            worldY = targetWorldY
            isMoving = false
        }
    }
}
