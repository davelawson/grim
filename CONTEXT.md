# Grimoire

Terms used to describe the game's cards, players, and play.

## Language

**Ring**:
The arrangement of wizards that defines their Aggressor and Target
relationships.

**Aggressor**:
A wizard's neighbour on their left in the Ring, whose Target is that wizard.

**Target**:
A wizard's neighbour on their right in the Ring, against whom most of that
wizard's offensive actions are directed. Some offensive actions may instead
be directed against their Aggressor.

**Integrity**:
A measure of the strength of a wizard's personal protections, reduced by
successful raids or attacks and recoverable through Arcane Focus. A wizard
at zero Integrity remains in play; any further Integrity reduction eliminates
them.

**Immunity**:
A temporary protection from attacks during the turn immediately following
a loss of Integrity, including that turn's attack resolution.

**Bend the knee**:
A wizard's voluntary submission at the first Bend the Knee step after a
successful Attack, granting their Aggressor one Domination point and swapping
their adjacent Ring positions. Submission forfeits the remainder of the
current turn and hands play to their original Target; declining closes that
opportunity.

**Elimination**:
A wizard's removal from active play following a further Integrity reduction
while already at zero Integrity. They immediately lose response eligibility
and take no further turns, but their Ring position is removed only after the
current exchange finishes.

**Domination**:
A Victory path advanced when another wizard bends the knee to the wizard
pursuing it, or when that wizard eliminates their Target.

**Attack**:
A direct offensive action against a wizard's Target, available to a Minion
with Might Expertise and undertaken alone or with other such Minions. It
compares Might deterministically and reduces the Target's Integrity by 1
on success.

**Defence**:
A wizard's opposition to their Aggressor's Attack, including standing Might,
committed Minions and Arcane Focus, and defensive Spell effects.

**Reaction**:
A wizard's permitted response to a pending action, Reaction, or mandatory
card effect during an interaction.

**Response window**:
An opportunity for eligible wizards to declare Reactions before a pending
action, Reaction, or mandatory card effect resolves.

**Exchange**:
An action or card effect together with all its Reactions and resulting
mandatory effects. It finishes when no effects remain pending on the Stack,
unless five-step victory ends the game sooner.

**Pending effect**:
A declared action, Reaction, or mandatory card effect that has not yet resolved.

**Resolution**:
An uninterrupted execution of one pending effect's legal instructions and
applicable success or failure consequences.

**Stack**:
The ordered collection of pending effects, with the most recently declared
effect resolving first.

**Public information**:
The current table state and revealed declarations, commitments, targets,
choices, and payments visible to all players.

**Hidden information**:
The contents of each wizard's hand, deck, and discard pile.

**Priority**:
An eligible wizard's opportunity to declare a Reaction or pass during a
response window.

**Pass**:
A decision to decline the current opportunity to react; it does not waive
later opportunities after another Reaction is declared or an effect resolves.

**Minion**:
A durable asset representing a wizard's follower. A Minion in play can be
assigned to a challenge when ready. Minions are the only cards that possess
Expertise, and may also provide Affinity independently of that Expertise.

**Starting Minion**:
A named Minion card drafted from the three offered with a chosen Chantry into
the wizard's starting deck. Its name does not make it unique across wizards;
more than one wizard may have a copy of the same named card.

**Pupil**:
A newly recruited Minion from Promising Pupil, with a fixed random Expertise
domain different from the Element of its wizard's drafted Chantry.

**Exhausted**:
A card state preventing another assignment or activation; committing a Minion
exhausts it until its owner's next turn starts, even if its effect is countered.
Passive effects and its surviving pending contribution remain available.

**Ready**:
A card state in which the card is available for use.

**Spent**:
An asset removed from play and put into its owner's discard pile, for any
reason. "Spent" and "discarded" describe the same outcome for an asset; it
must be drawn and played again before it can be used.

**Rumour**:
A mandatory opportunity or disruption drawn from the rumour mill. It remains in play for the
duration printed on its card, alongside any other active rumours, then returns
to the rumour mill when discarded.

**Rumour mill**:
Each wizard's personal pool of rumour cards from which they draw each turn.
Its contents depend on the wizard's affinities and the cards' affinity
requirements; all listed affinities are required unless a card states an
alternative. Eligibility is checked only when a rumour is drawn. Investigate's
additional draw uses the final domain totals of its Challenge, including
assigned Minion Expertise and committed card contributions, and is set aside
until the next turn. A later loss of Affinity does not remove a rumour already
in play. The mill always has a rumour available to draw.

**Challenge**:
An action offered by a card for a Challenger to attempt. Rumours, victory-path
cards, and Activities may offer Challenges. Requirements, when present, can
combine cards, Affinities, and Minion Expertise. Investigate has no additional
requirements beyond a surviving assigned Challenger.

**Succeeded Challenge**:
A Challenge overcome before its opportunity window closes, typically earning
its printed reward for the wizard.

**Failed Challenge**:
A Challenge not overcome by its attempted resolution, or left unovercome when
its opportunity expires without an alternative succeeding. Its card determines
the failure consequences; an unused reusable victory-path card does not Fail.

**Challenger**:
The actor undertaking a Challenge: a ready Minion, or the wizard intervening
directly by playing Arcane Focus; several Minions may undertake one Challenge
together. A Spell may explicitly permit an attempt without a Challenger.

**Capacity**:
The combined cards, Affinities, Expertise, and payable resources available to
a wizard for meeting a Challenge's requirements.

**Affinity**:
A wizard's capacity to influence a particular domain, measured by the combined
contributions of the wizard's sources in that domain. Multiple sources of the
same Affinity stack. Affinity is neither spent nor reserved when used, and its
total cannot fall below zero.
The Affinities are Nobility, Church, Fire, Earth, Air, Water, Infernal, Might,
and Popularity.

**Expertise**:
A ranked measure of a Minion's capability within a particular domain. It adds
to the wizard's corresponding Affinity when that Minion undertakes a Challenge,
or contributes Might to an Attack or Defence, even when the wizard does not
otherwise have that Affinity.

**Elemental Affinity**:
One of the Fire, Earth, Air, or Water Affinities. Elemental Affinities are a
named group within the otherwise flat list of Affinities.

**Fire Affinity**:
A wizard's influence over flame, heat, energy, passion, destruction, and
purification.

**Earth Affinity**:
A wizard's influence over soil, stone, minerals, growth, stability, endurance,
and material wealth.

**Air Affinity**:
A wizard's influence over wind, atmosphere, weather, movement, thought, and
communication.

**Water Affinity**:
A wizard's influence over water, ice, healing, emotion, adaptation, and change.

**Church Affinity**:
A wizard's influence within the religious institution, including its clergy,
doctrine, temples, and political reach. It does not imply divine power.

**Nobility Affinity**:
A wizard's influence over the court and the government it controls.

**Infernal Affinity**:
A wizard's capacity to wield demonic and witchcraft practices through pacts,
communion, and related relationships with infernal powers. Witchcraft draws on
those powers even when its practitioner does not understand its ultimate source;
it includes potions, transformations, communion with animals, and moon rituals.

**Might Affinity**:
A wizard's influence over physical force, combat, weapons, courage,
intimidation, endurance, and martial leadership.

**Popularity Affinity**:
The favour a wizard holds among common people.

**Element card**:
A card that supplies one temporary point of its Elemental Affinity to a paired
action without changing the wizard's general Affinities. Each play supports
one action, and the card must be regained before reuse.

**Cost**:
The resource payment required to play a card or use an ability, distinct from
Affinity requirements. Paid resources remain spent if the effect is countered.

**Play requirement**:
A condition that must be met to play a card, such as a numeric Affinity threshold,
distinct from its resource Cost and any Challenge requirements.

**Ability**:
A card's passive, triggered, or activated behaviour, with its own applicable
timing, requirements, payments, targets, and effect.

**Persistence**:
A card property describing whether it remains in play after resolution and
what ends its stay.

**Disposal**:
The destination of a card when it would leave play or be discarded, subject to
its category and printed properties.

**Duration**:
The printed limit on how long a card or effect lasts.

**Arcane Focus**:
An Activity that lets the wizard undertake a Challenge directly and contributes
one Wis plus one point of any single Affinity to its paired action. It also
offers Integrity recovery as an alternative use of the card itself, or can
contribute one Might to an Attack or Defence.

**Spell**:
A card representing a magical working. Its Affinity requirements, if any, are
limited to Fire, Earth, Air, and Water; an explicit card property may permit
a Challenge attempt without a Minion or Arcane Focus.

**Association**:
The elemental inputs used when creating a Spell or Artifact. The finished card
prints an appropriate Affinity requirement rather than an association field.

**Artifact**:
A magical Asset that may stand alone or also belong to another card category,
such as an Artifact Minion.

**Consumable**:
A card that can remain in play until an ability explicitly consumes it;
passive benefits do not automatically consume it.

**Chantry**:
A wizard's starting Location, already in play when the game begins, that
continuously contributes one point of its associated Elemental Affinity.
Fire, Earth, Air, and Water each have a corresponding Chantry.

**Resource card**:
A card that provides Wealth or Wis when played. Recurring rewards typically
generate ephemeral resource cards into the wizard's hand during Start Phase
Income, separately from End Phase hand refill.

**Ephemeral card**:
A card that exists only for the current turn and is destroyed at its end,
whether played or unused, even when the End Phase is skipped.

**Retained card**:
A card kept in the hand through the End Phase Discard step, reducing the cards
needed to refill that hand. Multiple cards may be retained within maximum hand size.

**Maximum hand size**:
The number of ordinary cards the wizard may hold after End Phase Discard and
the size to which Draw refills that hand.

**Victory path**:
One of the ways a wizard can Ascend, tracked by visible steps of progress.
Reaching five steps wins the game as soon as the effect or submission awarding
that progress has finished.

**Turn order**:
The continuous sequence of wizards' turns following the Ring, with no rounds
or other grouping of turns. Normal handoff follows the active wizard's current
Target, kneeling uses their original Target, and elimination of the active
wizard passes play to the first survivor to their right in the former Ring.

**Turn**:
The span in which one active wizard takes ordinary actions and eligible other
wizards may respond to those actions.

**Active wizard**:
The wizard whose turn it is and who may take ordinary actions during that turn.

**Start Phase**:
The opening phase of a turn, comprising Recovery, Income, Maintenance,
Bend the Knee, Rumour, and Effects in that order.

**Bend the Knee step**:
The optional submission checkpoint after Maintenance and before Rumour.

**Rumour step**:
The Start Phase step after Bend the Knee in which personal rumours enter play,
with their effects deferred to Start.Effects.

**Maintenance step**:
The Start Phase step in which the active wizard pays recurring upkeep costs
after Income and before the Bend the Knee step.

**Start.Effects**:
The final Start Phase step for mandatory and optional card effects, including
rumour effects and deferred expiry effects, ordered by the active wizard
subject to printed timing.

**Deferred expiry effect**:
A mandatory effect caused by automatic expiry when kneeling or elimination
skips the End Phase, awaiting the surviving wizard's next Start.Effects even
though its source card has expired.

**Actions Phase**:
The phase for the active wizard's ordinary card plays, Challenges, and Attacks,
with response exchanges completed between ordinary actions.

**End Phase**:
The final phase of a normal turn, comprising mandatory End Phase effects,
Cleanup, hand Discard down to maximum size or below, and Draw in that order.
Kneeling or elimination skips its remaining steps.
