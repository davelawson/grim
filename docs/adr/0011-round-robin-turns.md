# Round-robin turns with player responses

The turn sequence has no round hierarchy, as clarified by
[Continuous turns and immediate kneeling swaps](0015-continuous-turns-and-kneeling-swaps.md).
The response and owner-turn recovery decisions below are retained.

Grimoire replaces simultaneous independent turns and a shared post-submission
interaction phase with round-robin turns, in which one active wizard takes
ordinary actions at a time and other eligible players may respond. Generally
the active wizard's Aggressor and Target can respond, while global actions may
require reactions from other players. This trades simultaneous ordinary play
for a more straightforward turn flow with back-and-forth interaction. Turn
order follows the Ring rather than an independent fixed sequence. After a
successful Attack, its Target may Bend the Knee in the Start Phase step
immediately after Maintenance and before Rumour,
ending their turn and updating their Ring position; their original Target then
begins their turn. Kneeling skips the remaining Start steps, Actions, and
the entire End Phase, with no additional future turn forfeited. There is no
hand refresh; an extra Investigate rumour stays set aside until the wizard
next reaches Start.Rumour. Turn-based expiry still happens automatically at
the handoff, without a response window: resources expire, ephemeral cards are
destroyed, and temporary effects end at their printed turn limits.
Mandatory effects caused by that expiry wait until the wizard's next
Start.Effects, using normal response rules; the expired card leaves play
without cancelling or duplicating its deferred effect.

Exhausted Minions recover at the start of their owner's turn without an
additional once-per-round restriction. This lets a Minion defend before its
owner's turn and act after recovering, while one used during the owner's turn
remains exhausted during subsequent turns. Used hand cards must be regained
before reuse.
