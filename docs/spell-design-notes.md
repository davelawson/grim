# Completed provisional Spell design

The Spell design interview is complete and the shared design is confirmed.
The canonical rules and card texts are recorded in [GAME.md](../GAME.md):

- [Spell support and timing](../GAME.md#spell-support-and-timing-provisional)
  covers declaration, support requirements, countering, stacking, and upkeep.
- [Delayed recovery and returning Minions](../GAME.md#delayed-recovery-and-returning-minions-provisional)
  defines overlapping delays, leaving play, and normal replay restrictions.
- [The first-playable Spell catalogue](../GAME.md#first-playable-spell-catalogue-provisional)
  contains 32 designs across 11 pools: four per Element, four unaligned, and
  two for each of the six elemental pairs.
- [Research Spell](../GAME.md#research-spell) costs 1 Wis and either creates a
  chosen Element card or discovers a Spell from the pool selected by zero,
  one, or two distinct Element inputs. Owning a design does not exclude it
  from future offers.
- [Spell interaction checks](../GAME.md#spell-interaction-checks) records the
  intended outcomes for support, countered costs, Attack failure, recovery,
  stacked assets, and private scrying.

These provisional playtest rules and all 32 Spell effects are implemented in
the playable-v1 engine, including Research, support, countering, and restart-safe
private scrying. Balance
values remain to be tested, particularly stacked lasting Spells, free support
and defensive Reactions, healing, and removal causing an additional discard
after a failed Attack. Research Artifact rules and three- or four-Element designs
remain outside this completed slice.
