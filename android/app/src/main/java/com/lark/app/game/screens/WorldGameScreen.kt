package com.lark.app.game.screens

import com.badlogic.gdx.Gdx
import com.badlogic.gdx.ScreenAdapter
import com.badlogic.gdx.graphics.Color
import com.badlogic.gdx.graphics.GL20
import com.badlogic.gdx.graphics.OrthographicCamera
import com.badlogic.gdx.graphics.Pixmap
import com.badlogic.gdx.graphics.Texture
import com.badlogic.gdx.scenes.scene2d.Stage
import com.badlogic.gdx.utils.viewport.FitViewport
import com.badlogic.gdx.utils.viewport.ScreenViewport
import com.lark.app.data.api.RetrofitClient
import com.lark.app.data.repository.GameRepository
import com.lark.app.game.BeachTownGame
import com.lark.app.game.dialog.DialogManager
import com.lark.app.game.dialog.DialogState
import com.lark.app.game.map.BeachTownMap
import com.lark.app.game.map.BuildingRenderer
import com.lark.app.game.map.TileRenderer
import com.lark.app.game.map.TileType
import com.lark.app.game.npc.Npc
import com.lark.app.game.npc.NpcRegistry
import com.lark.app.game.npc.NpcRenderer
import com.lark.app.game.player.CollisionMap
import com.lark.app.game.player.PlayerController
import com.lark.app.game.sprite.NpcSprites
import com.lark.app.game.sprite.PlayerSprite
import com.lark.app.game.ui.ActionButtonActor
import com.lark.app.game.ui.DPadActor
import com.lark.app.game.util.PixelFont

class WorldGameScreen(private val game: BeachTownGame) : ScreenAdapter() {

    companion object {
        const val TILE_SIZE = 16f
        const val MAP_WIDTH = 40
        const val MAP_HEIGHT = 30
        // Wide viewport - 18 tiles wide, ~22 tiles tall (matches reference density)
        const val VIEWPORT_WIDTH = 288f
        const val VIEWPORT_HEIGHT = 352f
    }

    // Camera and rendering
    private val camera = OrthographicCamera()
    private val viewport = FitViewport(VIEWPORT_WIDTH, VIEWPORT_HEIGHT, camera)

    // Map
    private val mapData = BeachTownMap.create()
    private lateinit var tileRenderer: TileRenderer
    private lateinit var buildingRenderer: BuildingRenderer

    // Player
    private lateinit var playerSprite: PlayerSprite
    private lateinit var playerController: PlayerController
    private lateinit var collisionMap: CollisionMap

    // NPCs
    private lateinit var npcs: List<Npc>
    private lateinit var npcSprites: NpcSprites
    private lateinit var npcRenderer: NpcRenderer

    // UI
    private lateinit var uiStage: Stage
    private lateinit var dPad: DPadActor
    private lateinit var actionButton: ActionButtonActor
    private lateinit var pixelFont: PixelFont

    // Dialog
    private lateinit var dialogManager: DialogManager
    private val repository = GameRepository(RetrofitClient.api)

    // Control panel
    private var controlPanelHeight = 0
    private lateinit var separatorTex: Texture

    // State
    private var facingNpc: Npc? = null
    private var dpadChoiceHeld = false

    override fun show() {
        // Initialize rendering
        tileRenderer = TileRenderer()
        buildingRenderer = BuildingRenderer()
        playerSprite = PlayerSprite()
        npcSprites = NpcSprites()
        pixelFont = PixelFont()

        // Panel separator line texture
        val sp = Pixmap(1, 1, Pixmap.Format.RGBA8888)
        sp.setColor(0.25f, 0.26f, 0.35f, 1f)
        sp.fill()
        separatorTex = Texture(sp)
        sp.dispose()

        // Initialize NPCs
        npcs = NpcRegistry.getNpcs(game.languageCode)
        npcRenderer = NpcRenderer(npcSprites)

        // Initialize collision map
        collisionMap = CollisionMap(mapData, MAP_WIDTH, MAP_HEIGHT, npcs)

        // Initialize player at boardwalk spawn
        playerController = PlayerController(20, 20, collisionMap)

        // Camera initial position
        camera.position.set(
            playerController.worldX + TILE_SIZE / 2,
            playerController.worldY + TILE_SIZE / 2,
            0f
        )
        camera.update()

        // Initialize UI stage (full screen, screen-pixel coordinates)
        uiStage = Stage(ScreenViewport())

        // D-pad (size set in resize)
        dPad = DPadActor()
        uiStage.addActor(dPad)

        // Action button (size set in resize)
        actionButton = ActionButtonActor(pixelFont) {
            onActionPressed()
        }
        uiStage.addActor(actionButton)

        // Dialog manager
        dialogManager = DialogManager(uiStage, pixelFont, repository, game.languageCode) {
            playerController.locked = false
            facingNpc = null
        }

        Gdx.input.inputProcessor = uiStage
    }

    override fun render(delta: Float) {
        val screenW = Gdx.graphics.width
        val screenH = Gdx.graphics.height

        // 1. Clear entire screen with dark panel color
        Gdx.gl.glViewport(0, 0, screenW, screenH)
        Gdx.gl.glClearColor(0.10f, 0.11f, 0.16f, 1f)
        Gdx.gl.glClear(GL20.GL_COLOR_BUFFER_BIT)

        // 2. Update game state
        tileRenderer.updateAnimation(delta)

        val dir = dPad.currentDirection
        val dialogState = dialogManager.currentState
        if (dialogState is DialogState.Choices) {
            if (dir == com.lark.app.game.sprite.Direction.UP) {
                if (!dpadChoiceHeld) { dialogManager.moveChoiceUp(); dpadChoiceHeld = true }
            } else if (dir == com.lark.app.game.sprite.Direction.DOWN) {
                if (!dpadChoiceHeld) { dialogManager.moveChoiceDown(); dpadChoiceHeld = true }
            } else {
                dpadChoiceHeld = false
            }
        } else if (!playerController.locked) {
            dpadChoiceHeld = false
            if (dir != null) {
                playerController.tryMove(dir)
            }
        }
        playerController.update(delta)

        // Camera follow
        val targetX = playerController.worldX + TILE_SIZE / 2
        val targetY = playerController.worldY + TILE_SIZE / 2
        camera.position.x += (targetX - camera.position.x) * 5f * delta
        camera.position.y += (targetY - camera.position.y) * 5f * delta

        // Clamp camera to map bounds
        val halfW = viewport.worldWidth / 2
        val halfH = viewport.worldHeight / 2
        camera.position.x = camera.position.x.coerceIn(halfW, MAP_WIDTH * TILE_SIZE - halfW)
        camera.position.y = camera.position.y.coerceIn(halfH, MAP_HEIGHT * TILE_SIZE - halfH)
        camera.update()

        // Check for facing NPC
        facingNpc = if (!playerController.locked) findFacingNpc() else facingNpc

        // Update action button label
        val currentDialogState = dialogManager.currentState
        actionButton.label = when {
            currentDialogState is DialogState.Choices -> "OK"
            currentDialogState !is DialogState.Hidden -> ">"
            facingNpc != null -> "Talk"
            else -> "A"
        }

        // 3. Apply game viewport (top portion of screen, above control panel)
        viewport.apply()

        // 4. Render game world
        game.batch.projectionMatrix = camera.combined
        game.batch.begin()

        // Ground tiles
        tileRenderer.renderGround(game.batch, mapData, MAP_WIDTH, MAP_HEIGHT)

        // Y-sorted pass: buildings + objects + characters
        for (y in MAP_HEIGHT - 1 downTo 0) {
            // Composite buildings (rendered first so decorations/NPCs draw on top)
            buildingRenderer.renderRow(game.batch, y)
            // Non-building object tiles (fences, lamps, palms, etc.)
            for (x in 0 until MAP_WIDTH) {
                val tileType = TileType.fromInt(mapData[y * MAP_WIDTH + x])
                if (tileRenderer.isObjectTile(tileType) && !tileRenderer.isBuildingTile(tileType)) {
                    tileRenderer.renderObjectTile(game.batch, x, y, tileType, mapData, MAP_WIDTH, MAP_HEIGHT)
                }
            }
            for (npc in npcs) {
                if (npc.tileY == y) {
                    npcRenderer.renderSingle(game.batch, npc, facingNpc)
                }
            }
            if (playerController.tileY == y) {
                playerSprite.render(game.batch, playerController)
            }
        }

        game.batch.end()

        // 5. Draw control panel separator line (screen-space)
        Gdx.gl.glViewport(0, 0, screenW, screenH)
        game.batch.projectionMatrix = uiStage.viewport.camera.combined
        game.batch.begin()
        // Separator line at top of control panel
        game.batch.draw(separatorTex, 0f, controlPanelHeight.toFloat() - 2f, screenW.toFloat(), 2f)
        // Subtle inner highlight below separator
        val prevColor = game.batch.color.cpy()
        game.batch.setColor(0.16f, 0.17f, 0.24f, 1f)
        game.batch.draw(separatorTex, 0f, controlPanelHeight.toFloat() - 4f, screenW.toFloat(), 2f)
        game.batch.color = prevColor
        game.batch.end()

        // 6. Dialog
        dialogManager.update(delta)

        // 7. UI stage (d-pad, action button)
        uiStage.act(delta)
        uiStage.draw()
    }

    private fun findFacingNpc(): Npc? {
        val dir = playerController.facing
        val checkX = playerController.tileX + dir.dx
        val checkY = playerController.tileY + dir.dy
        return npcs.find { it.tileX == checkX && it.tileY == checkY }
    }

    private fun onActionPressed() {
        if (dialogManager.currentState !is DialogState.Hidden) {
            dialogManager.advance()
            return
        }
        val npc = facingNpc ?: return
        playerController.locked = true
        dialogManager.startInteraction(npc)
    }

    override fun resize(width: Int, height: Int) {
        // Split screen: bottom 42% = control panel, top 58% = game
        controlPanelHeight = (height * 0.42f).toInt()
        val gameH = height - controlPanelHeight

        // Update game viewport to fit in the top portion only
        viewport.update(width, gameH)
        // Shift viewport up above the control panel
        viewport.setScreenY(viewport.screenY + controlPanelHeight)

        // UI stage covers full screen
        uiStage.viewport.update(width, height, true)

        // Position controls in the control panel area
        positionControls(width.toFloat(), controlPanelHeight.toFloat())
    }

    private fun positionControls(screenW: Float, panelH: Float) {
        // D-pad: left side, sized proportionally
        val padSize = minOf(screenW * 0.38f, panelH * 0.72f)
        dPad.setSize(padSize, padSize)
        dPad.setPosition(screenW * 0.06f, (panelH - padSize) / 2f)

        // Action button: right side
        val btnSize = padSize * 0.50f
        actionButton.setSize(btnSize, btnSize)
        actionButton.setPosition(
            screenW - btnSize - screenW * 0.10f,
            (panelH - btnSize) / 2f
        )
    }

    override fun dispose() {
        dialogManager.dispose()
        tileRenderer.dispose()
        buildingRenderer.dispose()
        playerSprite.dispose()
        npcSprites.dispose()
        npcRenderer.dispose()
        uiStage.dispose()
        pixelFont.dispose()
        separatorTex.dispose()
    }
}
