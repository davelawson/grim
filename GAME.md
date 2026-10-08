# Grimoire Game Mechanics

Grimoire is a turn-based deck building game.  Each player takes on the role of a wizard striving to overtake their rivals, and be the first to ascend.  They develop spells and  artifacts, recruit minions to overcome challenges, in their pursuit of power.

## Basics of Play

### Provisional first playable draft

The rules in this section are a starting point for testing the flow of a game.
Card costs, hand size, and challenge thresholds are illustrative balance values,
not final numbers. The intended loop is to improve the deck and establish
assets, then use those cards to pursue rumours and visible victory progress.
Play is a continuous series of turns following Ring order, with no rounds or
other higher-level grouping and one active wizard taking ordinary actions at
a time. Their Aggressor and Target
can generally respond to those actions; global actions may require reactions
from other players. Player interactions allow back-and-forth responses,
including defensive Spellcasting. Every ordinary action opens a response
window before resolving; pending effects resolve with the latest response
first. This turn-flow design is agreed and remains provisional until playtested.

The working phase structure below is provisional. Each wizard receives an
initial five-card hand during setup. Thereafter, they keep the cards they
choose not to discard and refill their hand during their own End Phase,
providing cards for Reactions between turns and for the next Actions Phase.
Kneeling skips that End Phase and therefore provides no hand refresh.

#### Start Phase

Complete these steps in order:

1. **Recovery:** Recover the active wizard's exhausted cards. Newly played
   Minions become eligible to Attack and undertake Challenges at this owner's
   turn start.
2. **Income:** Recurring income generates its printed resources or resource
   cards. Generated ephemeral resource cards do not count toward the ordinary
   hand refill target. An asset first played during an earlier turn begins
   recurring income at its owner's next Start Phase.
3. **Maintenance:** Pay the active wizard's recurring upkeep costs. Resource
   cards, including cards generated during Income, may be played as part of a
   payment without separate response windows. An Asset whose upkeep is not
   paid is discarded. Grants from later Start Phase card effects are not yet
   available to fund this step.
4. **Bend the Knee:** A Target eligible to submit after a successful Attack
   may do so now. Their turn ends immediately, their Ring position updates,
   and their original Target begins the next turn. A wizard who kneels does
   not proceed to Rumour, Effects, Actions, or any End Phase step. There are
   no End Effects, Cleanup, Discard, or Draw. This forfeiture is the entire
   turn-loss cost; no additional future turn is skipped. Turn-based expiry
   still occurs automatically at the handoff, without a response window.
   Submission is available only at the first such step after that Attack;
   declining closes the opportunity.
5. **Rumour:** Draw and play the active wizard's normal personal rumour from
   their currently eligible pool. Under the existing Investigate rule, also
   play any extra rumour set aside for this turn, without rechecking its
   eligibility. If kneeling skipped an earlier Rumour step, the extra rumour
   stays set aside until this step is next reached. Rumour effects wait until
   Start.Effects.
6. **Effects:** Resolve mandatory and optional Start Phase card effects,
   including rumour effects and mandatory expiry effects deferred by kneeling,
   in an order chosen by the active wizard, subject to printed card-specific timing.
   Card effects use the Stack and response protocol; finish each exchange
   before proceeding to the next effect.

A resource card can be played as part of a required card-effect payment.
Its resource becomes available immediately. Rumours offer their Challenges
while in play; playing the rumour is separate from overcoming a Challenge.
Full Moon's optional hunt occurs during Start.Effects, subject to its printed
requirement to precede other optional actions.

Recovery, recurring income generation, and maintenance payment do not
independently open response windows. Mandatory card effects do.

#### Actions Phase

The active wizard plays new cards, undertakes Challenges, and carries out
Attacks, choosing one ordinary action at a time. Every declaration opens its
response exchange, and the entire exchange finishes before another ordinary
action is chosen. Each action resolves against the state surviving earlier
actions and Reactions.

Used hand cards must be regained before reuse. Minions exhaust on commitment
and recover during their owner's next Recovery step. Cards supply actions,
with no separate action-slot limit. One action can combine a Challenger and
committed support cards.

#### End Phase

Entering this phase is final: the wizard cannot return to ordinary actions.
Complete these steps in order:

1. **Effects:** Resolve mandatory End Phase card effects, including printed
   expiry-only Failure effects that are now due. Effects use the Stack and
   allow permitted Reactions. Complete their exchanges before Cleanup removes
   expired cards or effects.
2. **Cleanup:** Remove expired effects and destroy ephemeral cards, whether
   played or unused. All players' unspent available resources generated during the
   current turn expire, including resources generated by responders.
3. **Discard:** The active wizard may choose cards to discard, but must finish
   at or below their maximum hand size, provisionally five ordinary cards.
   They may keep multiple remaining cards, replacing the earlier
   one-retained-card limit. Played cards already follow their destinations when resolved or
   countered; gained cards enter the discard pile when gained.
4. **Draw:** Draw enough ordinary cards to refill the remaining hand to the
   wizard's maximum hand size. If the draw pile runs out, shuffle
   the discard pile to continue. The hand is available for Reactions before
   the owner's next ordinary turn, so cards used to respond reduce that turn's
   options.

Cleanup, voluntary discard, and hand refill do not independently open response
windows. Kneeling skips the entire End Phase, including these steps.

Turn order follows the Ring. After being successfully attacked, a Target may
Bend the Knee in the step immediately after Maintenance and before Rumour: their turn ends,
their Ring position updates, and their original Target begins the next turn.
Their existing hand is retained without discarding or drawing; the forfeited
turn is the only turn skipped. Even without Cleanup, unspent resources expire,
ephemeral cards are destroyed, and temporary effects end when their printed
turn limits are reached. This automatic expiry does not run the End Phase or
open a response window. A mandatory effect caused by expiry, such as Burn the
Witch's Failure, is deferred to the wizard's next Start.Effects and uses the
normal Stack and response rules. The expired rumour returns to the mill; its
pending Failure remains due even though the rumour has left play. If another
kneeling turn skips Start.Effects, keep that effect pending until the step is
actually reached, without creating another copy of the penalty.

Generated Wealth and Wis can take the form of resource cards placed in the
discard pile. Recurring rewards typically generate ephemeral resource cards
for the current turn instead. Playing a Wealth or Wis card makes its resource
available immediately. The wizard can spend that resource on a later cost in
the same turn; the played card follows its normal discard or destruction rule
as part of that play. Some effects, such as Solar Eclipse, grant resources directly.
An asset first played during a turn begins generating recurring cards at the
start of its owner's next turn.
Any unspent amount expires at the end of the current active wizard's turn,
including amounts generated as part of a Reaction. An Affinity is a
reusable numeric capacity, not a resource that is spent. Every active source of
the same Affinity contributes to the wizard's total, and those contributions
stack. Using an Affinity neither spends nor reserves it, so the same total can
meet multiple requirements across actions and turns. Effects may increase or
reduce a total for a stated duration, but no Affinity can fall below zero. Durable assets
remain in play and can generate future income. Available cards and active
opportunities, rather than a fixed action count, determine how much a wizard
can do.

### Attacks, Defence, and Integrity (provisional)

Each wizard starts at 3 Integrity and can recover up to a maximum of 3.
Reaching zero Integrity leaves the wizard in play; any further Integrity
reduction eliminates them.

Chosen offensive targets are restricted to the wizard's Ring neighbours.
Attacks always target their Target; other offensive abilities may target their
Aggressor where permitted by their printed rules. Global effects use their
printed affected group.

A ready Minion with Might Expertise can initiate an Attack without a separate
Attack card. Several such Minions may combine in the wizard's single Attack
for the turn. A newly played Minion cannot Attack until its owner's next turn
starts, but can defend immediately. Attacking and defending Minions must be
ready and exhaust when committed, retaining their surviving assigned
contributions even though exhausted.

Each side totals its wizard's current Might Affinity and the Might Expertise
of its committed Minions, plus 1 Might if Arcane Focus supports that side.
Defensive Spell effects may modify the comparison or remove participants at
their printed timing. Removed Minions do not contribute at resolution.
The Attack succeeds only when attacking Might exceeds defending Might by at
least 2. Success reduces the Target's Integrity by exactly 1; an Attack does
not instead select an Asset or resource as its objective.

On an unsuccessful Attack, the attacker chooses and discards one surviving
committed attacking Minion, if any, as part of the Attack's uninterrupted
resolution. Uncommitted passive Might sources and Arcane Focus cannot be
chosen for this loss. Interrupted or countered Attacks retain this failure
penalty, alongside the normal expenditure of cards and exhaustion of surviving
participants. Discarded Minions must be drawn and played again, rather than
being destroyed.

A loss of Integrity grants Immunity, during which the wizard cannot be
attacked. The exact duration under continuous turns and changing neighbours
remains an open design decision in [TODO.md](TODO.md).

Defence combines standing Might, deliberate commitments, and permitted
defensive Spell Reactions. Different strategies should have distinct defensive
tools so every wizard need not develop the same military engine. Whether
separate defensive preparation exists alongside responses remains open.

### Ring movement (provisional)

Play follows the Ring continuously, without a round boundary or a rule reserving
one turn per wizard in a cycle. When a wizard Bends the Knee, they grant their
current Aggressor one Domination point and immediately swap positions with
that Aggressor, without paying additional Integrity. Their original Target begins the
next turn, as specified by the kneeling handoff rule above. Kneeling remains
unavailable when only two wizards survive.

For example, in A → B → C → D → A, each arrow points toward the wizard's
Target. If B kneels to A, the Ring becomes B → A → C → D → B, and C takes
the next turn. Each subsequent kneel applies to the Ring as it then exists;
there is no batch of kneelings resolved together.

An eliminated wizard immediately loses response eligibility and takes no
further turns, but retains their Ring position until the entire current
exchange finishes, including all Reactions and resulting mandatory effects.
Then remove eliminated positions and reconnect their surviving neighbours
without changing the survivors' relative Ring order.

If the active wizard survives and the game continues, they finish their turn
normally using the updated Ring. Eliminating their Target does not grant a
second Attack that turn. A normal turn hands play to the active wizard's
current Target; kneeling instead uses the original Target as described above.

If the active wizard is eliminated, finish the exchange, remove eliminated
positions, and give the next turn to the first surviving wizard to their right
in the Ring as it stood before removal, unless the game has ended.

Before that handoff, skip the eliminated wizard's remaining phases and hand
refresh, but perform automatic turn-end expiry without a response window.
All players' unspent resources and ephemeral cards expire, as do temporary
effects whose printed turn limits have been reached. Any mandatory effect
created by that expiry for a surviving wizard is deferred to their next
Start.Effects, following the same non-duplication rule as kneeling.

With two survivors, both Aggressor and Target refer to the same opponent;
response rotation includes each wizard once, and kneeling is unavailable.

### Response protocol (provisional)

An exchange comprises the original action or card effect, its Reactions, and
all resulting mandatory effects; it finishes when the Stack is empty. Five-step
victory may end the game before the exchange finishes, as described under
Victory below.

1. Declare an ordinary action, commit its cards, and pay its costs. Resource
   cards may be played as part of a payment without separate response windows.
   Add the effect to the Stack. Countering it does not refund its cards or
   resources. Minions exhaust on commitment and retain their surviving assigned
   contributions to the pending effect. A counter leaves surviving Minions
   exhausted in play. Once resolved or countered, played cards go to their owner's discard pile unless
   they remain in play, such as a successfully played Asset; destruction
   properties replace discard as usual.
2. Fix the eligible participants for this exchange: the active wizard and
   their Aggressor and Target, or the wider group specified by a global action.
   Responding does not automatically invite the responder's other neighbours.
3. Rotate Priority through eligible players in turn order, beginning after
   whoever declared the latest effect. A player may pass or declare an
   expressly permitted Reaction or defensive commitment, paying and committing
   at declaration. Ordinary development actions are unavailable as Reactions.
4. A new Reaction goes on top of the Stack and resets the pass count. After
   every eligible participant passes consecutively, resolve the newest pending
   effect. A pass declines the current opportunity, not later opportunities.
5. After each uninterrupted resolution, check for five-step victory. If the
   game ends, stop without adding or resolving further triggers or effects.
   Otherwise, add mandatory effects triggered by resolution before checking
   whether the Stack is empty; the active wizard chooses and reveals the order of
   simultaneous triggers. If the active wizard has been eliminated, the next
   living wizard in Ring order chooses that order instead. If effects remain
   pending, open another window with the active wizard receiving Priority first.
   Eliminated players leave response rotation immediately, but their pending
   effects remain subject to normal resolution checks. Priority passes only
   among surviving eligible players.
6. When the Stack is empty, the exchange ends. Remove eliminated Ring positions
   and check for last-survivor victory or a zero-survivor draw. If the match
   continues and the active wizard survives, return to the phase that opened
   the exchange; another ordinary action is available only in Actions Phase.
   If the active wizard was eliminated, perform automatic turn-end expiry and
   hand play to their surviving successor as described under Ring movement.

Elimination of the active wizard is a contingency for possible future effects;
the currently defined cards and Attack rules do not cause it during their own turn.

Costs and play requirements are checked at declaration. Resolve against the
surviving state: recheck Challenge capacity, participants, and targets, and
exclude removed Minions' contributions. Every action attempts to resolve
regardless; if it fails, it counts as a failure and normal failure penalties
apply, without refunding costs, at their printed timing. For example, a failed
Burn the Witch attempt does not apply its expiry penalty early; succeeding
before its deadline still avoids that penalty. Each effect resolves as one
uninterrupted operation: perform portions that remain legal and apply normal
failure consequences if its success requirements are unmet. Challenge success
rewards require every requirement to be satisfied. On an unsuccessful Attack,
the attacker chooses and discards a surviving committed attacking Minion
during resolution, before another response window opens.
Neighbours may react to ordinary actions; only an original
global action admits wider participation, and a global Reaction never expands
the eligible participants.

All players see the current table state and declared cards, commitments,
targets, options, and payments as they happen. Each wizard can inspect their
own hand and the contents of their own draw pile, but not the draw pile's order.
Every wizard can inspect every discard pile, including its order. Opponents'
hands and draw-pile contents are hidden. Response eligibility still follows
the fixed participant rules above.

### Challenges and progress

A rumour can have an on-play effect without a Challenge, as Solar Eclipse does.
When a rumour offers actions, each is a Challenge; a victory-path card or an
Activity such as Investigate can also offer a Challenge. A Challenge Succeeds
when overcome and typically grants its printed reward. Most Challenges have a
limited opportunity window. If none of a rumour's offered Challenges is
overcome before that window expires, they Fail, and any printed Failure effect
applies once. A rumour enters play when
played and offers its Challenges while it remains in play. Investigate reserves
its draw until the next turn, without starting its duration. On Success of one of
them, the rumour is discarded immediately unless its card says it remains.
Playing a victory-path card starts its attempt. Every attempt requires a
Challenger: either a ready Minion already in play or the
wizard acting directly by playing Arcane Focus, unless a Spell explicitly
permits an attempt without either. The player may assign ready
Minions and commit eligible support cards from the hand. A Challenge's
requirements may involve resources provided by cards, Affinity thresholds, or
any combination.
Arcane Focus and an Element may both support the same Challenge: each pairs
with that Challenge, contributes only to that attempt, and is used once.
The wizard cannot begin an attempt without enough capacity from the chosen
Challenger, support, Affinities, Expertise, and payable resources to meet every
printed requirement.
For each printed Affinity threshold, add the wizard's current total in that
domain, every matching Expertise value on the assigned Minions, and any
action-specific contributions from committed cards. All matching Expertise on
an assigned Minion applies at once. A Minion needs no matching Expertise to
undertake a Challenge. Multiple Minions may participate together, combining
their Expertise. Arcane Focus allows the wizard to act as Challenger and adds
one Wis plus one point of any single Affinity to that attempt. There is no
random modifier yet. A newly played Minion cannot Attack or undertake a
Challenge until its owner's next turn starts, but can defend and provide
passive Affinity immediately. Committed support cards follow their normal discard or
destruction rule after an attempt. An assigned Minion exhausts when committed
and remains exhausted until the start of its wizard's next turn; its passive Affinities
remain available.
An unused, reusable victory-path card returns to the discard pile unchanged
unless retained; it does not Fail.
A rumour remains in play for its printed duration and leaves play as its card
directs.

At the opening draft, each wizard chooses a victory-path card for a primary
path. It returns through the deck, providing a relatively reliable chance to
advance. A successful attempt immediately adds one visible step, removes the
attempted card, and puts its costlier next rank into the discard pile. Other
cards may offer alternate paths or change the wizard's primary
path; those effects are not designed here. Reaching the fifth step wins the
game immediately, including during a turn, instead of creating another rank.

### Starting deck and card interactions

Every wizard begins with exactly eleven cards: a Chantry already in play and
a ten-card starting deck. Seven are fixed; four come from draft choices or
are linked to a choice:

| Starting card | Count | How it enters the deck |
| --- | ---: | --- |
| +Chantry | 1 | Draft one of four Chantries and place it in play during setup. |
| *Element | 1 | Add the Element linked to the chosen Chantry. |
| *Research Spell | 1 | Fixed starting card. |
| *Research Artifact | 1 | Fixed starting card. |
| **Basic Wis | 2 | Fixed starting cards. |
| *Victory challenge | 1 | Draft a card for a primary victory path. |
| *Basic Wealth | 1 | Fixed starting card. |
| *Minion | 1 | Draft one of the three Minions offered by the chosen Chantry. |
| *Arcane Focus | 1 | Fixed starting card. |
| *Investigate | 1 | Fixed starting card. |
| **Total** | **11** | |

The four Chantry cards have not yet been named or fully balanced. There is one
Chantry for each of Fire, Earth, Air, and Water. A Chantry is a Location that
provides one point of its linked Elemental Affinity and generates an ephemeral
Wis card during each of its owner's Income steps, including the first. Each
starts in play; if discarded, it costs 1 Wis and no Affinity to replay from
the hand. These are provisional values. Any additional effects depend on the
chosen card.
The Chantries give their wizards a similar opening rhythm. Their linked
Elements distinguish the spells and artifacts those wizards create through
Research. The twelve named starting Minions and their draft offers are listed
under Drafting below. Their Expertise changes which Challenges the wizard can
overcome early and which rumours Investigate can find.

**Arcane Inquiry I** remains an illustrative choice for the victory challenge
slot, with **Arcane Inquiry II** as its replacement after a success.

A basic Wis or Wealth card provides one unit of its named resource immediately
when played and costs no resources or Affinity. Generated ephemeral Wis and
Wealth cards have the same zero cost and one-unit value unless their source
states otherwise. Research Spell creates a new spell card in the discard pile;
Research Artifact creates a new artifact card there. When an Element is
committed with either Research card, it supplies one point of its Affinity and
the created card comes from that Element's creation pool and prints an
appropriate Affinity requirement. Association describes the elemental inputs
used during creation; it is not printed on the finished card. Without an
Element, the creation is unaligned. An Element may instead be committed with one other card to provide
one point of its Affinity for that action. It does not change the wizard's
general Affinity total or ordinary rumour eligibility. Each play supports one
action; the card follows its normal disposal and must be regained before reuse,
without any additional usage counter.

Arcane Focus is an Activity card that can be committed with one other card. It
can serve as the Challenger when the wizard intervenes directly. It
contributes one Wis and one point of any single Affinity, chosen when the play
is made, to that action only. Neither contribution enters the wizard's general
resource or Affinity pool. The paired cards form one action and are each used
once.

For example, pairing a Fire Element with Research Spell puts a spell with an
appropriate Fire requirement into the discard pile. Pairing Arcane Focus with an Arcane Inquiry I that
requires 1 Wis and Might 1 satisfies both requirements if Might is the Affinity
chosen for Arcane Focus. A starting Minion assigned to a Challenge contributes
all of its matching Expertise in addition to the wizard's reusable Affinity
totals.

## Victory

There are many paths to victory in Grimoire. Every wizard's progress on each
path is visible to all. A surviving wizard wins immediately by advancing 5 steps along
any one path, even during a turn; there is no wait for a turn handoff or a
round boundary. Victory Challenge progress and replacement occur immediately
on Success.

Finish the single uninterrupted effect or kneeling decision that awards the
winning progress, including its other instructions and rewards, then end the
game. Do not resolve any further pending effects, mandatory triggers, End Phase
steps, or turn handoff. No further expiry or Ring-removal bookkeeping is run
after the match has ended.

Last-survivor victory uses a different checkpoint: finish the entire exchange
before declaring the sole surviving wizard the winner. If no wizards survive
that exchange, the game is a draw. A five-step victory reached during the
exchange still ends the game at its earlier checkpoint.

If one uninterrupted effect leaves multiple surviving wizards with five or
more steps on a path, compare their other paths at that same checkpoint.
For each contender, set aside one winning path and sort the remaining path
progress from highest to lowest. Compare those values in order; the wizard
with more progress at the first difference wins. If every comparison ties,
the tied contenders share victory. No additional turn or response is granted.

Current progress awards are single steps. Future content that permits multiple
winning paths with unequal scores must specify which winning path is excluded
from this comparison before that content is introduced.

### Ending by agreement (provisional)

An unfinished match may be abandoned by a strict majority of its remaining,
non-eliminated players, or of all participants during setup. Starting a vote
freezes gameplay: no action, Reaction, or other gameplay progression occurs
until every eligible player has submitted a ballot, even if a majority is
already apparent. Ballots cannot be changed and voting has no deadline.
Eliminated players' ballots are invalidated and they do not participate further.

A completed vote with a majority ends the match without awarding victory or
resolving pending effects; a tie fails. A failed vote resumes the suspended
gameplay state. A new round may begin only after the pending round concludes,
with fresh ballots. Elimination cannot make a previously failed vote pass.

Voting can begin at a point where gameplay waits for input, regardless of
Priority. A server administrator may end the match with the same abandoned
outcome, including during a vote. Disconnecting alone never changes gameplay
or ends the match.

Opening a round casts no ballot: every eligible player, including its
initiator, submits a separate ballot. Ballots are open but voter identities
are not disclosed; players can see aggregate yes/no totals and how many
ballots remain outstanding, without identifying voters or non-voters.

### Domination

Gain one Domination point whenever another wizard Bends the Knee to you or
you eliminate your Target. A successful Attack that neither eliminates the
Target nor leads to submission awards no Domination point by itself.

Eliminating a Target also transfers that wizard's actual starting Chantry
card, wherever it currently resides, into the Aggressor's ownership and
discard pile. Additional Chantries acquired by the eliminated wizard do not
transfer. Apply these elimination rewards during the Attack's resolution,
before the victory checkpoint or any further response window.

### Arcane

By plumbing the depths of magical knowledge, the wizard progresses along the path to an arcane victory.

### Influence

The world is still largely controlled by mundane humans, and gaining sufficient sway in their society allows a wizard to progress towards an influence victory.

### Council

By cultivating favour with their peers, a wizard can eventually find themselves progressing towards becoming the recognized leader of the wizards, and a council victory.

### Fame

By completing epic quests that appear during play, a wizard can become eternally famous, and progress towards a fame victory.

## Cards

### Common Card Properties

#### Cost

Cost is the resource payment required to play a card or use an ability.
Paying Wis or Wealth depletes that resource. It is printed separately from
Affinity requirements.

#### Play Requirements

Affinity requirements are numeric thresholds met to play a card; meeting one
neither spends nor reserves the Affinity. A bare Affinity symbol is shorthand
for a requirement of 1. These requirements are distinct from requirements to
undertake a Challenge printed on that card.

#### Abilities

A card may have passive, triggered, or activated abilities. Each ability
specifies its own timing, requirements and payments, targets, and effect where
applicable. Reaction permission belongs to the relevant ability, including
any restriction on the actions or effects it can respond to. A passive benefit
and a separately activated ability may coexist on one card.

#### Offered Challenges

Cards that offer Challenges use a common block: the Challenge's name,
requirements, Success effect, Failure effect and its timing, and availability
(including attempt timing, repeatability, and expiry where applicable).
Requirements distinguish resource payments from Affinity thresholds and other
conditions. Rumours, Victory Challenges, and Activities may use this block;
a card can offer alternative Challenges with different requirements and outcomes.

#### Persistence

Persistence describes whether a card remains in play after resolving and what
ends its stay. It is separate from disposal, duration, and mandatory timing.

##### Instants

An instant card supplies an effect without remaining in play. Its effect
follows the applicable Stack rules, then the card enters its owner's discard
pile unless a destruction property applies.

##### Consumables

A Consumable can remain in play until an ability explicitly consumes it.
That ability identifies when the card is discarded, subject to any disposal
property. A passive benefit does not automatically consume the card.

##### Durables

Durable cards remain in play and active until something specifically removes them, or an upkeep is not paid.

#### Disposal

A card normally enters its owner's discard pile when it leaves play unless
its category or printed rules specify another destination; rumours return
to the rumour mill. Disposal properties can coexist with persistence properties.

##### Destructibles

Destructible cards are destroyed when they would otherwise go to the discard pile.

#### Duration

A printed duration limits how long a card or effect lasts. Duration is separate
from persistence and disposal.

##### Ephemeral

Ephemeral cards exist only for the current turn. At the end of that turn, they
are destroyed whether played or unused, including when kneeling skips Cleanup.

#### Play Timing

Play timing specifies when a card must or may be played; it is separate from
how long the card stays in play or where it goes afterward.

##### Mandatory

Mandatory cards must be played at their required timing. Rumours are played during Start.Rumour,
but their effects wait until Start.Effects. Required effects resolve through
the Stack with response windows, but a Challenge on the
card can remain available while the card stays in play. A rumour drawn by
Investigate is instead set aside and becomes mandatory to play during
Start.Rumour of the following turn.

#### Upkeep

A recurring cost paid during the owner's Start Phase Maintenance step; if
unpaid, the Asset is discarded.

## Card Categories

Cards have a category determining their applicable parameters and general
template. Some categories have subcategories. An Artifact may stand alone as
an Artifact or also belong to another category, such as an Artifact Minion.
In addition to parameters, any card can include text that doesn't necessarily conform to any category rules.

### Card Anatomy by Category (provisional)

This card anatomy is agreed and remains provisional until playtested.

Cards identify their name and category. Where applicable, they print resource
Costs, Affinity play requirements, abilities, persistence, disposal properties,
duration, and play timing. Ability payments and Challenge requirements are
distinct from requirements to play the card. An optional field does not have
to appear on a card that does not use it.

| Category | Category-specific information |
| --- | --- |
| Asset — Minion | Passive Affinities; Expertise by domain; recurring income; Upkeep; abilities |
| Asset — Location | Passive Affinities; recurring income; Upkeep; abilities |
| Asset — Artifact | Passive Affinities; recurring income; Upkeep; abilities; fields of another category if present |
| Spell | Effect-specific numbers; ability timing and Reaction permission; duration; explicit Challenger exception if present |
| Rumour | Eligibility requirements; duration; immediate or ongoing abilities, beneficial or disruptive; offered Challenges and their outcomes |
| Resource | Resource kind and yield |
| Element | Element represented and paired-action contribution |
| Victory Challenge | Victory path; rank; offered Challenge; advancement and replacement on Success |
| Activity | Abilities; targets or paired cards; offered Challenge if present |

Asset Affinities, income, Upkeep, and extra abilities are optional. Only Minions
possess Expertise. A card's readiness, exhaustion, remaining duration, and
temporary modifications are changing play state rather than extra printed
base stats. Creation association is not a field on the finished card.

### Assets

Cards that, once played, tend to remain in play, providing benefits to the wizard.  Assets tend to provide affinities and generate resources.

Every Asset may have passive Affinities, recurring income, Upkeep, and abilities;
these fields are optional. Only a Minion may possess Expertise. A standalone
Artifact uses these common Asset fields without requiring another category.

#### Minions

Minions are assets that represent loyal followers that are able to undertake challenges on behalf of the wizard.
Minions can become exhausted. When exhausted, any printed upkeep must still be
paid and they cannot attempt Challenges, but they continue to provide their
passive Affinities.
Assigning a Minion to a challenge exhausts it. Spending or discarding a Minion
or another asset means the same thing: move it from play to its owner's discard
pile. Its owner must draw and play it again before it can be used. This applies
whenever an asset is sent to the discard pile, including for unpaid upkeep or
a forced discard.
In the provisional first playable draft, a ready Minion in play may be assigned
to a Challenge and contribute all of its matching Expertise alongside support
cards committed from the hand. An assigned Minion becomes exhausted when
committed, even if its pending effect is countered. Exhausted cards cannot be assigned or activated again
until they recover at the start of their wizard's next turn. A newly played
Minion can defend immediately, but cannot Attack or undertake a Challenge
until its owner's next turn starts.

A Minion's ordinary Affinities become available as soon as it enters play,
whether ready or exhausted and whether or not it is assigned to an action.
Assignment does not reserve those Affinities: they may support the Minion's
Challenge and every other action for as long as the Minion remains in play.
Expertise is different: it is a ranked Minion capability added to a
Challenge that Minion undertakes, with Might Expertise also contributing to
an Attack or Defence the Minion joins. A Minion may print both Affinity and Expertise
in the same or different domains.

Parameters:

- Upkeep, if printed: resources paid during Start Phase Maintenance, or the Minion is discarded
- Affinities: numeric influence the Minion provides continuously while in play
- Expertise: ranked domain capability added to assigned Challenges, or Might to committed Attacks and Defence
- Recurring income, if printed
- Passive, triggered, or activated abilities, if printed

#### Artifacts

Artifacts are magical assets that are typically created by the wizard, although
they can sometimes be retrieved from Challenges or events, or taken from rivals.
An Artifact may be a standalone Artifact or also have another category, such as
a Minion golem. Consumable and Durable describe persistence, rather than
Artifact subcategories. A Consumable Artifact is discarded when an ability
explicitly consumes it; passive benefits do not automatically consume it.
A Durable Artifact remains in play until removed.

Parameters:

- Its printed abilities
- Optional passive Affinities, recurring income, and Upkeep
- Parameters appropriate to another category, if it has one

#### Locations

Locations are assets that represent structures or settings that fall under the dominion of the wizard.  They are placed in the play area, but not equipped to any individual minion or the wizard.  Locations typically afford affinities.

Parameters:

- Affinities: numeric influence provided while the Location remains in play
- Recurring income, Upkeep, and abilities, if printed

### Spells

Spells are probably the most common cards in your deck. A Spell's Affinity
requirements, if any, use only Fire, Earth, Air, or Water. Some Spells have no
Affinity requirement. Spells may also require resources such as Wis. Other card
types may require any Affinity that fits their fiction.
Spells have no universal Power stat: amounts are printed in individual effects.
Duration appears when applicable, and Reaction permission and any Challenger
exception are part of the relevant ability.

### Rumours

Rumours are opportunities or disruptions drawn separately from the wizard's own deck. The
rumour mill is each wizard's personal pool of rumour cards drawn from each
turn. Each rumour card may have one or more numeric Affinity requirements for joining
that mill. All listed affinities are required unless the card explicitly states
an alternative. The initial catalogue includes eight rumours with no Affinity
requirement, available even after a wizard loses their Chantry. The mill
expands and contracts as the wizard's affinities change. Eligibility is
recomputed from current affinities before each normal draw; a temporary
affinity can unlock a rumour only while it is still active at that draw.
Investigate makes an additional draw when its Challenge Succeeds, using the
domain totals available during that attempt, including assigned Minion
Expertise and committed card contributions. Losing an affinity does not remove
a rumour already in play; after it is discarded,
the rumour is not drawable until its affinity requirements are met again. A
rumour active or reserved by Investigate for a wizard cannot be drawn again
by that wizard until it is discarded. Discarded rumours return to the rumour
mill. The mill always has at least one drawable rumour at each draw.
The card set and printed durations must leave enough eligible rumours that
active and reserved ones leave enough drawable rumours at every draw;
there is no special draw rule for an empty mill.
Until weighting is designed, each drawable rumour has an equal chance of being
drawn. Wizard and rival actions may affect which types are eligible as their
effects on affinities are defined.

Normal rumours are drawn and played during Start.Rumour; an Investigate draw
is also mandatory to play in that step on the following turn. Rumour effects
resolve during Start.Effects, after Maintenance and the Bend the Knee step.
Each rumour card states
how long it remains in play. A rumour without a Challenge returns to the rumour mill
when its printed duration ends without a Success or Failure outcome. If it
offers Challenges, its wizard has an opportunity to overcome one of them on
every turn the rumour remains in play. If none is overcome before the window
closes, apply the rumour's printed Failure effect once and return it to the
rumour mill. The turn a rumour enters play counts as the first turn of its
duration. A wizard still draws a new rumour
each turn while earlier rumours remain in play, so several can be active at
once.

The [initial rumour catalogue](docs/rumours.md) defines twenty playable rumours
and every reward card they introduce: eight affinity-free rumours, nine with
one requirement covering all nine Affinities, and three requiring Fire with
Infernal, Earth with Nobility, or Water with Church. None awards victory
progress. The values are provisional until playtested.

The original five rumours remain in the catalogue. Promising Pupil creates
a new Minion on each Success, with Expertise 1 chosen uniformly from the eight
domains other than the drafted Chantry's Element. Its Expertise remains fixed
for the entire game and does not supply passive Affinity. Solar Eclipse grants
Wis and temporary Affinities during Start.Effects. Burn the Witch offers arrest or
recruitment and discards one chosen Asset only if unresolved after three turns.
New Leyline grants recurring Wis through a gained Location. Full Moon offers
an early Might hunt for Wealth and Affinities lasting through the next turn.

Most catalogue Challenges require domain totals of 1–2 and resource payments
of 1–2. Most unresolved opportunities expire without an additional penalty;
Burn the Witch is the only harmful Failure in this set. Some gained assets
provide lasting Affinity and require no Affinity to play, allowing new domains
first reached through Investigate to enter the normal rumour mill.

The catalogue's availability check covers the initial deck's one Investigate
per turn and durations of at most three turns. Future content that increases
draw frequency or duration must revisit that check before it is introduced.

### Resources

Resource Cards are cards that generate resources when played. Typical resources generated would be Wealth and Wis.

### Elements

An Element card represents Fire, Earth, Air, or Water. Commit it with a single
other card to supply one point of that Affinity for the paired action. This
action-specific point does not change general Affinity totals or ordinary
rumour eligibility; Investigate includes it for its additional draw. If the
paired action is Research Spell or Research Artifact, the
created card comes from that Element's creation pool and prints an appropriate
Affinity requirement. The creation association is not a field on the finished
card. Each play of an Element supports only one paired action. Used Elements
follow their normal disposal and must be regained before reuse; there is no
additional once-per-turn or once-per-round restriction.

### Victory Challenges

A drafted victory challenge provides a repeatable opportunity to advance one path. On success, it scores a step and is replaced by a more advanced version, as described above.

### Activities

Activity cards represent basic actions that can be performed by the wizard. In the provisional first playable draft, each Activity is played for its printed effect and may name another card as a target. An Activity does not always require a paired card. Activities are one of the primary methods to tailor the wizard's deck.

#### Research Activity

The two starting Research cards create new cards directly in the discard pile.
Pairing an Element with Research determines the creation association and the
appropriate Affinity requirement on the finished card; the association itself
is not printed. Without an Element, the creation is unaligned. How the specific
new Artifact is selected remains to be designed; Research Spell selection is
recorded in the parked Spell design notes.

##### Research Spell

Creates a new spell card and adds it to the wizard's discard pile.

##### Research Artifact

Creates a new artifact card and adds it to the wizard's discard pile.

#### Arcane Focus

Commit Arcane Focus with one other card to contribute one Wis and one point of
any single Affinity to that paired action. It also allows the wizard to serve as
the Challenger. Neither contribution enters the wizard's general pool.

Alternatively, play Arcane Focus itself to recover 1 Integrity, up to the
maximum of 3. Recovery costs 3 Wis total: Focus supplies 1 Wis toward that
action, leaving 2 Wis to pay from available resources. Focus must be drawn
normally and played for this use; it is not a permanently available recovery
ability.

Focus may instead support an Attack or Defence with 1 Might. Each play uses
one of its alternatives; it must be regained before reuse.

#### Investigate

Investigate is a starting Activity that offers a Challenge with no additional
Failure penalty. It requires a Challenger, which may be the wizard acting
through Arcane Focus. It has no additional success requirements beyond a
surviving eligible Challenger at resolution.
At that Success, draw one random eligible rumour from the wizard's normal
rumour mill using the final domain totals available for the attempt: the
wizard's current Affinities, the Expertise of Minions directly assigned to
Investigate, and any action-specific Affinity contributions from committed
cards. Set the generated rumour aside until the wizard's next Start.Rumour
step, when it is played in addition to the normal rumour. Do not check its
eligibility again when it is played. A rumour set aside this way is unavailable
for other draws until it returns to the mill. If kneeling skips Start.Rumour,
keep the extra rumour set aside until the wizard next reaches that step.

#### Deconstruct Activity

The deconstruct card allows the wizard to destroy a card, and recover resources in the process.

#### Specialize Activity

The specialize card allows the wizard to duplicate another card.

#### Collaborate Activity

The collaborate activity allows the wizard to generate rumours and opportunities for other wizards to interact with them.

## Deck Construction

The ten-card starting deck combines fixed cards with a victory challenge and
a Minion chosen during the opening draft. Choosing a Chantry adds its linked
Element to that deck; the Chantry itself starts in play. The full draft
procedure and the identities of the victory challenge offers remain to be
designed.

### Drafting

The opening draft offers all four Chantries, one each for Fire, Earth, Air, and
Water. Each wizard chooses one Chantry, one victory challenge, and one of three
Minions offered for that Chantry. The Chantry choice also determines the matching
Element card. Place the chosen Chantry in play before the first turn. Each
starting Minion costs 1 Wis and requires 1 Affinity in its
Chantry's Element to play from the hand. These are play costs, not payments made
during the draft. Each has exactly 1 Expertise in the domain shown below, no
ambient Affinity, no upkeep, and no additional ability. The Minion cannot
undertake a Challenge until its owner's next turn starts. Its Expertise
contributes only when it is directly assigned to a Challenge.

| Chantry | Starting Minion | Role | Expertise |
| --- | --- | --- | --- |
| Fire | Vera Coalheart | Furnace keeper | Fire 1 |
| Fire | Bastian Redhand | Enforcer | Might 1 |
| Fire | Orla Nine-Cinders | Pactmaker | Infernal 1 |
| Earth | Hedda Stonewake | Stonewright | Earth 1 |
| Earth | Leontine Grange | Estate agent | Nobility 1 |
| Earth | Rurik Ironvale | Quarry guard | Might 1 |
| Air | Mina Cloudscript | Weather reader | Air 1 |
| Air | Cassian Bellweather | Court messenger | Nobility 1 |
| Air | Elsie Many-Tongues | Storyteller | Popularity 1 |
| Water | Nadia Rivermark | Canal keeper | Water 1 |
| Water | Sister Agnes Reed | Abbey healer | Church 1 |
| Water | Samantha Twice-Born | Marsh witch | Infernal 1 |

Names identify card designs, not globally unique game pieces: different
wizards may draft copies of the same named Minion. The victory challenge offers
and draft order remain open.

### Static

Each starting deck also receives Research Spell, Research Artifact, two basic Wis,
one basic Wealth, Arcane Focus, and Investigate.

### Meta

The two starting Research cards create new spells or artifacts, allowing players to develop their deck over successive turns. The pace of this growth remains to be tested.

## Affinities and Expertise

The nine Affinities are Nobility, Church, Fire, Earth, Air, Water, Infernal,
Might, and Popularity. Fire, Earth, Air, and Water are collectively the
Elemental Affinities; otherwise the list is flat. Their domains are:

- **Nobility:** the court and the government it controls
- **Church:** the religious institution, including clergy, doctrine, temples,
  and political reach, but not divine power
- **Fire:** flame, heat, energy, passion, destruction, and purification
- **Earth:** soil, stone, minerals, growth, stability, endurance, and material
  wealth
- **Air:** wind, atmosphere, weather, movement, thought, and communication
- **Water:** water, ice, healing, emotion, adaptation, and change
- **Infernal:** demonic and witchcraft practices, including pacts, potions,
  transformations, communion with animals, and moon rituals; all witchcraft
  ultimately draws on infernal power, even unknowingly
- **Might:** physical force, combat, weapons, courage, intimidation, endurance,
  and martial leadership
- **Popularity:** favour among common people

Domains may overlap. A card may require more than one when each is essential.
Affinities never oppose or cancel one another automatically, and a wizard may
possess every Affinity. A specific card may restrict which Affinities or
Expertise can support its action; merely possessing an Affinity never
disqualifies the wizard.

Affinity requirements are numeric. All active sources in a domain stack, and
the same total remains available for every action. Any card or ongoing effect
may provide Affinity. Effects may modify the combined total for a stated
duration, to a minimum of zero. Ordinary sources and Expertise should usually
provide 1, exceptional ones may provide 2, and first-playable requirements
should generally range from 1 to 5. Affinity totals have no hard upper cap.

Only Minions possess Expertise. When one or more Minions undertake a Challenge,
add all of their matching Expertise to the wizard's Affinity totals for that
attempt. Action-specific contributions from cards such as Elements and Arcane
Focus are added after that. Expertise and action contributions do not affect
the wizard's general Affinity totals or ordinary rumour eligibility.
Investigate's additional draw explicitly uses both. Expertise applies to
assigned Challenges, with Might Expertise also contributing to committed
Attacks and Defence. An action-specific contribution can meet an Affinity
requirement of its paired action, including a card's play requirements.

## Notes To Incorporate

- When performing many types of actions that involve magic, playing additional spell cards allows you to inform the results of the action.
- When attempting a challenge, how do we add randomness to the result without screwing the player?
