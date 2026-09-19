# The EinBexiii/dragonfly fork

## How `next` is built

`next` is never edited directly. It is `upstream/master` (df-mc/dragonfly) plus
every branch below merged in with a merge commit, in the order listed. Each
branch is based on `upstream/master` unless it says otherwise, so any of them
can be dropped or sent upstream on its own. `next` is rebuilt from that list
whenever it changes; a branch that is not on the list is not in `next`.

`upstream/master` is at the version named by the most recent `dragonfly:
Updated to ...` commit below the merges.

## Features

| Branch | Adds |
|---|---|
| `pr/entity-bugfixes` | Collision reported for every entity, orb metadata only for experience orbs, fire-immune end crystals, lightning fire spread and disk-loaded strikes. Frozen. |
| `pr/ui-facade` | The player's display methods moved onto a UI facade. Frozen. |
| `feature/push-out-of-blocks` | An item stuck inside a block is pushed out of it. |
| `feature/death-animation` | A dead entity plays its death animation. Carries the entity refactor it needs: a shared living entity base (`LivingEnt`), the player embedded on it, shared damage, fall, fire, air, explosion, knock back and attack immunity logic, exported behaviour extension points and behaviour-contributed metadata. |
| `feature/integration` | Wires the entity features into each other: riding (seats, showing and steering a ride), entity inventories, entity actions, sounds and particles, `HandleEntityHurt`, and movement broadcast by what changed. Built on the death-animation refactor. |
| `feature/entity-target` | A mob's target is reported to viewers. |
| `feature/tack-items` | Saddle, horse armour, an entity's armour inventory, and armour rendered on an entity's body. |
| `feature/entity-trading` | Trading with entities, on top of the inventory an entity carries. |

## Fixes not yet in upstream

Each of these tracks an open upstream pull request and should be dropped from
the list once upstream merges it.

| Branch | Fixes | Upstream |
|---|---|---|
| `fix/chunk-callback-reentrancy` | Chunk callbacks run at the top level of a transaction, never nested inside another callback. Closes the world freeze when a player rejoins at a chunk border. | #1381 |
| `fix/loader-viewer-reentrancy` | The loader calls its viewer without holding its own lock, so a chunk request completed from inside `ViewChunk` cannot deadlock the world goroutine. | #1419 |
| `fix/loader-change-world` | Removing a loader's viewers when it changes world is scheduled through `World.Do` instead of queued directly, so a full transaction queue cannot freeze both worlds. Based on `fix/loader-viewer-reentrancy`. | #1407 |
| `fix/entity-chunk-unload-leak` | An entity whose `Close` does not remove it from the world is removed when its chunk unloads instead of leaking. | #1412 |
| `fix/unsaved-chunk-duplication` | A chunk stays dirty while it holds entities, and block entity changes are noticed by hashing their NBT, so items and entities are not duplicated by a chunk that was never saved. | #1400 |
| `fix/nil-block-entity-nbt` | A block entity that encodes to nil NBT no longer panics the chunk send path. | #1275 |
| `fix/spectator-game-mode` | Spectator is reproduced as measured on BDS 1.26.45: a game mode change is one player game type update addressed to the player's own unique ID (spectator is 6); a spectator's abilities carry the spectator layer ahead of the base layer, for the player and in the AddPlayer other clients get; the client hides a spectator by game type, so no invisibility is forced and no teleport is sent; a spectator cannot use items, and its request to stop flying is ignored. | #1285 |
| `perf/chunk-height-maps` | Height-map columns are cached and invalidated per column, and a chunk's surface is prepared once per sub-chunk response. | #1449 |

## Documentation

| Branch | Adds |
|---|---|
| `docs/fork` | This file. |
