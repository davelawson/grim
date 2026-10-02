# Continuous turns and immediate kneeling swaps

Grimoire is a continuous series of wizards' turns following the Ring, without
rounds or a per-cycle turn entitlement; this supersedes the round-based
organization accompanying [the earlier turn-flow decision](0011-round-robin-turns.md).
A wizard who Bends the Knee immediately swaps positions with their current
Aggressor and hands play to their original Target, retaining the established
turn-forfeiture and two-survivor restrictions. Each kneel changes the current
Ring independently, replacing the historical batch chain-reversal model in
[the earlier kneeling decision](0005-change-neighbours-after-a-raid.md), so
movement may bring a wizard's next turn forward or delay it without round
bookkeeping.

An eliminated wizard immediately leaves response rotation and takes no further
turns, but their position remains until the entire current exchange finishes;
then their surviving neighbours reconnect. This keeps neighbours stable
through the Attack, its Reactions, and any resulting mandatory effects.

The submission opportunity exists only at the first Bend the Knee step after
the successful Attack; declining closes it, so a defeat does not grant a
permanent option to change neighbours later.

A surviving active wizard continues their turn after eliminating a Target,
using their new neighbours and retaining the already-consumed Attack limit.
Normal handoff follows their current Target; if the active wizard is
eliminated, finish the exchange and hand play to the first surviving wizard
to their right in the former Ring, unless the match has ended.

Elimination skips the active wizard's remaining phases and hand refresh, but
automatic turn-end expiry still removes all players' unspent resources,
ephemeral cards, and temporary effects whose limits have been reached. It
opens no response window; any resulting mandatory expiry effect for a
surviving wizard waits until their next Start.Effects, as with kneeling.
