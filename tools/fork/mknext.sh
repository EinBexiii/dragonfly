#!/bin/sh
set -e
# Run from a docs/fork worktree; optionally rebuild a separate fork checkout.
if [ "$#" -gt 1 ]; then
  echo "usage: $0 [fork-checkout]" >&2
  exit 2
fi
here=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
repo=${1:-"$here/../.."}
repo=$(git -C "$repo" rev-parse --show-toplevel)
if [ -n "$(git -C "$repo" status --porcelain)" ]; then
  echo "refusing to rebuild a checkout with existing work: $repo" >&2
  exit 1
fi
branches='pr/entity-bugfixes pr/ui-facade feature/push-out-of-blocks feature/death-animation feature/integration feature/entity-target feature/tack-items feature/entity-trading fix/chunk-callback-reentrancy fix/loader-viewer-reentrancy fix/loader-change-world fix/entity-chunk-unload-leak fix/unsaved-chunk-duplication fix/nil-block-entity-nbt fix/spectator-game-mode perf/chunk-height-maps fix/immunity-excess-knockback fix/break-time-check fix/transfer-inventory-resync fix/persona-skin-data fix/live-skin-decoder-hardening fix/entity-handles-within feature/entity-view docs/fork'
# Fail before replacing next if a required local ref is missing.
for b in upstream/master $branches; do
  git -C "$repo" rev-parse --verify "$b^{commit}" >/dev/null
done
# upstream/master has no fork tools. Keep the resolver outside the checkout,
# and parse the whole rebuild function before checkout removes this script.
snapshot=$(mktemp -d "${TMPDIR:-/tmp}/dragonfly-rebuild.XXXXXX")
trap 'rm -rf -- "$snapshot"' EXIT
trap 'exit 130' INT
trap 'exit 143' TERM
cp "$here/resolve.py" "$snapshot/resolve.py"
rebuild() (
  cd "$repo"
  git checkout -q -B next upstream/master
  resolve() { python3 "$snapshot/resolve.py"; gofmt -w server/ 2>/dev/null; git add -A; git commit -q --no-edit; }
  for b in $branches; do
    git merge --no-ff --no-edit "$b" >/dev/null 2>&1 || resolve
  done
)
rebuild
