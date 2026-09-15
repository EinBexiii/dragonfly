# The EinBexiii/dragonfly fork

## How `next` is built

`next` is never edited directly. It is `upstream/master` (df-mc/dragonfly) plus
every branch in the Features through Documentation sections below merged in
with a merge commit, in the order listed. Each
branch is based on `upstream/master` unless it says otherwise, so any of them
can be dropped or sent upstream on its own. `next` is rebuilt from that list
whenever it changes; a branch that is not on the list is not in `next`.

`upstream/master` is at the version named by the most recent `dragonfly:
Updated to ...` commit below the merges.

`next` remains the Minecraft 1.26.45 source line for backends. The separate
1.26.50 preview below is excluded from its constituent list and rebuild tool.
Check the upstream protocol version before a rebuild too: moving the base
to a newer protocol is a coordinated upgrade, not routine fork maintenance.

## Rebuild tooling

Maintain [mknext.sh](tools/fork/mknext.sh), its [resolver](tools/fork/resolve.py),
and the [fork maintainer agent](.claude/agents/fork-maintainer.md) on `docs/fork`.

Run `sh tools/fork/mknext.sh` from a clean `docs/fork` worktree to rebuild that
checkout, or `sh tools/fork/mknext.sh /path/to/fork-checkout` to rebuild a
separate clean checkout. The target needs `upstream/master` and every branch
in the script as local refs; the tool does not fetch or push. It refuses a
dirty target or missing refs before replacing `next`. The resolver is copied
to temporary storage because checking out upstream removes these tools until
`docs/fork` is merged again. The shell loads the rebuild function first.

A rebuild replaces `next` and replays the full list onto `upstream/master`.
Inspect the target checkout and branch topology before using it. To land one
constituent branch, maintain the script's list on `docs/fork` and merge that
branch into `next`; a full rebuild is a separate operation.

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
| `fix/wall-collision-height` | A wall's collision box is the one the Bedrock server computes: a block and a half tall like a fence's whatever the drawn post and arm heights, one enclosing box, a straight postless run as narrow as its arms. Entity movement searches half a block further down for a fence or wall that reaches up into the box. | not yet sent |
| `fix/low-block-collision` | Rails, pressure plates and buttons have empty collision in every vanilla state, including all wood materials in the bundled palette. Snow layers collide at (layers−1)/8 height with the first layer empty; mud and soul sand collide at 7/8 height. Registers the missing states and their hashes without adding rail or redstone mechanics. Based on `a36ed0ed`; geometry follows the 2026-09-10 Rook block collision audit and state coverage follows the vanilla palette. | not yet sent |
| `fix/lily-pad-collision` | Lily pads collide 3/32 high, inset a sixteenth, as Bedrock has them (BDS constructor 0xc6840da); the old 1/64 top left a landing player unsupported for the anticheat. Based on `a36ed0ed`; row 6 of the 2026-09-10 Rook block collision audit. | not yet sent |
| `fix/trapdoor-collision` | Trapdoors collide 0.1825 thick in every orientation, as Bedrock has them (BDS shape helper 0x148a233d0, constants 0.1825 and 0.8175 at 0x14a8262e4 and 0x14a8262e0); Java's 3/16 left a landing player standing 0.005 inside the box, which the anticheat read as an unexplained clipped fall. Based on `a36ed0ed`. | not yet sent |
| `fix/flower-pot-collision` | Flower pots are registered with Bedrock's centred 3/8 by 3/8 collision box, collision only (BDS constructor 0x1489deba0 stores (0.3125, 0, 0.3125) to (0.6875, 0.375, 0.6875); the box getter 0x1414e4a20 ignores the pot's contents); the palette had the states but no block, so the fallback full cube left a standing player 5/8 inside it. Based on `a36ed0ed`; row 22 of the 2026-09-10 Rook block collision audit. | not yet sent |
| `fix/partial-block-collision` | Hoppers, grindstones, composters, brewing stands, iron trapdoors and cauldrons collide as the Bedrock server has them, read from its binary (Rook's reviews/2026-09-12-partial-block-boxes-bds-astra.md): hopper floor plate 10/16 to 11/16 with walls, body and facing outlet; grindstone inset an eighth so a standing one reaches the top; empty composter floor an eighth; brewing stand stem 0 to 7/8; iron trapdoor on the trapdoor model; cauldron floor 5/16 with eighth walls for every state. Based on `a36ed0ed`; audit rows 7, 11, 12, 13, 16 and 17. | not yet sent |
| `fix/fence-gate-connection` | A fence joins a gate only at the ends of the gate's bar, open or closed, as the Bedrock server connects them (Rook's reviews/2026-09-13-fence-gate-connection-bds-astra.md: gate connection mask at 0x1489d7ba0, fence helper 0x14158ee50); the faces a gate swings through carry no arm. Connecting to every gate put a 1.5 tall box where the client has none. | not yet sent |
| `fix/door-collision` | Door leaves are 0.1825 thick and both halves read direction and open state from the lower half and the hinge from the upper, as the Bedrock server does (Rook's reviews/2026-09-13-door-collision-bds-astra.md: shape routine 0x148e95ca0, pair resolver 0x148e95a20); iron doors exist, sharing the model and placement. A saved world's upper half carries default state, so the fork drew its leaf on the other side of the block. | not yet sent |
| `fix/hub-blocks` | Every block state the lobby's saved world holds decodes: rotated quartz, redstone lamps, huge and small mushrooms, dispensers, frosted ice, hardened stained glass panes, tripwire hooks, daylight detectors, pistons with their arms, chorus plants and flowers, with collision as the Bedrock server has it (Rook's reviews/2026-09-13-hub-blocks-bds-astra.md): the piston base a full cube in every state, the arm three boxes with the rear connector a quarter into the base, the chorus plant one box grown toward its neighbours, the detector 0.375 high, hooks and small mushrooms without collision. A one-sided pane or bar arm ends at the middle of the block, not the post edge; `cmd/blockhash` learned `MushroomBlockType`. No placement or redstone logic. | not yet sent |
| `fix/block-states` | Blocks the fork had but whose saved states it refused decode: exploding tnt, ladders and wall signs, banners and glazed terracotta with any facing_direction, occupied beds, the deprecated property on bone and hay blocks, item frame map and photo bits, the multiple grindstone attachment, powered lecterns, a portal without an axis, purpur on every axis; the poplar wood type with its three leaf colours. No box moves for a state the fork already accepted; directions hash in three bits. | not yet sent |
| `fix/blocks-shaped` | Nineteen block kinds with their own collision, read from the Bedrock server (Rook's reviews/2026-09-14-shaped-blocks-bds-astra.md): hanging signs, shelves, bell, dripleaves, pointed dripstone, amethyst buds and cluster, coral fans, copper bulbs, lightning rods, respawn anchor, turtle and sniffer eggs, chiseled bookshelf, sculk vein, glow lichen, resin clump, pale moss carpet, leaf litter, mangrove and hanging roots, frog spawn, pitcher plant and crop. Decode, collision, creative items and placement: each block is placed the way it was clicked, turning its front to the player or growing out of the face it was put on, and breaks when what carried it is gone; a copper bulb toggles on the rising edge of a redstone signal and lights the room, dimmed by its oxidation. Saved block-entity data still rides through untouched, and there is still no break info and no oxidation tick. A coral wall fan is decoded but never placed: the corpus does not resolve which cardinal each of its four coral_direction values is. Scaffolding stays unknown: Bedrock collides it only for an actor standing on its top, which the collision interface cannot see. | not yet sent |

## Fixes on the fork's own features

Built on `feature/death-animation` and `feature/integration`, so they cannot
go upstream on their own; they follow those branches.

| Branch | Fixes |
|---|---|
| `fix/immunity-excess-knockback` | A hit inside the attack immunity window deals its excess and counts as landed, but the window remembers that it did and `KnockBack` refuses it, on players and living entities, so a crit after a plain hit no longer sends the victim flying twice; the hurt animation and sound stay silent for it. |
| `fix/break-time-check` | A survival break is credited one mining frame per admitted client input frame, admitted by a budget that only takes the frames real time has passed, that a new block may carry at most two of (naming another block earns nothing), and that an episode already underway may bank so delayed inputs are all credited; a finish must name the block the episode started on, in reach, with its whole break time earned, an early one keeps the progress; a block the held tool breaks within a frame needs no episode but spends a frame. `StopBreak` is an abort. |

## Additions not yet sent upstream

Independent of the fork's features, candidates for upstream pull requests.

| Branch | Adds |
|---|---|
| `fix/transfer-inventory-resync` | The inventories and the held slot are sent again on the first input after a spawn, since a client arriving by transfer discards what reached it before its own spawn completed. |
| `feature/named-particle` | `particle.Custom` shows a resource pack's own particle effect by name, as `sound.Custom` plays its sounds. |
| `fix/entity-handles-within` | `Tx.EntityHandlesWithin`, `Tx.EntityHandles` and `Tx.EntityPosition`: entities as handles, without opening them, so a mob ranking hundreds of candidates opens only the ones it keeps. A handle's position is read through its world's transaction, which reports false for a handle that world no longer holds. |
| `feature/entity-view` | `Tx.EntityHandlesOf`: the World's entity handles grouped by identifier, rebuilt only when an entity is added or removed, so a crowd of mobs looking for the few players costs the players rather than the crowd. |

## Documentation

| Branch | Adds |
|---|---|
| `docs/fork` | This file, the rebuild tools, and the fork maintainer agent. |

## Protocol preview outside `next`

These branches are not part of the `next` list. Do not add either to
`tools/fork/mknext.sh`: the proxy moves to 1.26.50 before the backends, which
must keep their 1.26.45 source line. A later backend upgrade is a separate,
coordinated change.

| Branch | Adds |
|---|---|
| `fix/26.50-compat` | The 1.26.50 source adaptations, based on `a36ed0edb548298ab482939e1c653e39f9683719`, the upstream base of this `next` snapshot. Keeps the renamed block interaction action ignored and expresses dimension bounds as minimum Y and highest-Y distance (383 for -64..319). Pins the gophertunnel preview for standalone builds. A constituent of the preview only. |
| `next-26.50` | The `next` snapshot at `c7811091fcfd08b785fdb9e0acb5728d3c112aa4` with `fix/26.50-compat` merged in, for Minecraft 1.26.50 / protocol 2193. Composed separately without rebuilding or moving `next`. |

The preview's root `go.mod` selects
`github.com/EinBexiii/gophertunnel v0.0.0-20260915193041-cb7562a22ded`,
the 26.50 integration retaining v1.61.0's corrections. With Go 1.26.5,
`go build ./...` and `go vet ./...` verify the preview standalone. Replacements
in a library are not inherited: consuming roots selecting this Dragonfly
line must also replace `github.com/sandertv/gophertunnel` with that version.

`DefaultBiome` remains explicitly empty. Gophertunnel defines it as a biome
identifier and serialises the string unchanged, but Dragonfly's dimension
interface supplies no default biome and void generation identifies none.
The client's empty-name fallback is unverified; custom-dimension client
validation is still needed before treating the preview as rollout-ready.
