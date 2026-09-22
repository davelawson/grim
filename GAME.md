# Grimoire Game Mechanics

Grimoire is a turn-based deck building game.  Each player takes on the role of a wizard striving to overtake their rivals, and be the first to ascend.  They develop spells and  artefacts, recruit minions to overcome challenges, in their pursuit of power.

## Basics of Play

### Provisional first playable draft

The rules in this section are a starting point for testing the flow of a game.
Card costs, hand size, and challenge thresholds are illustrative balance values,
not final numbers. The intended loop is to improve the deck and establish
assets, then use those cards to pursue rumours and visible victory progress.
This draft covers independent player turns in a simultaneous round; direct
confrontation and shared contested opportunities need later rules.

A round works as follows:

1. Each wizard receives income from assets already in play, draws five cards,
   and receives one personal rumour card. If the draw pile runs out during the
   draw, shuffle the discard pile to continue. Obstacles drawn into the hand
   resolve automatically, as described below.
2. Each player privately submits an ordered sequence of card plays. A card can
   be used only once that round. Cards supply the actions; there is no separate
   action-slot limit. A play can use one card or combine an opportunity with
   supporting cards. A player may stop before using every card.
3. Once all plans are submitted, resolve each player's plays in the order that
   player specified. Each player's plays are independent in this draft, so
   there is no cross-player resolution order. A resource gained earlier in a
   sequence can pay for a later play. If a planned play cannot meet its cost
   when reached, it has no effect; its cards remain unused until cleanup.
4. Pay upkeep for cards that require it, then move played non-durable cards
   and unused hand cards to the discard pile. A newly gained card also enters
   the discard pile. Ephemeral cards are destroyed instead. Clear unspent
   Wealth and Wis, then begin the next round.

Wealth and Wis are turn-only resources in this draft. A resource card produces
its printed amount when played; the resource itself is spent on later plays.
An affinity is a prerequisite, not a resource that is spent. Durable assets
remain in play after being played and can provide future income. Cards in the
hand, rather than a fixed number of actions, determine how much a wizard can do.

### Challenges and progress

A rumour or victory-path card can offer a challenge. Playing that card starts
one attempt. The player commits any eligible support cards from the hand and
pays the costs of the opportunity and its support. Only those committed cards
supply challenge capability in this first draft; the wizard has no automatic
base score. Add their printed value for the challenge's attribute and compare
it with the threshold. Meeting or exceeding the threshold succeeds; there is
no random modifier yet. Committed cards are spent whether the attempt succeeds
or fails. An uncompleted, reusable victory-path card returns to the discard
pile unchanged. A rumour is ephemeral and is destroyed after its attempt or at
round end.

At the opening draft, each wizard chooses a victory-path card for a primary
path. It returns through the deck, providing a relatively reliable chance to
advance. A successful attempt immediately adds one visible step on that path,
removes the attempted card, and puts its costlier next rank into the discard
pile. Other cards may offer alternate paths or change the wizard's primary
path; those effects are not designed here. Reaching the fifth step wins the
game instead of creating another rank.

### Example starting deck and cards

For a paper playthrough, start with three **Wis**, two **Wealth**, two
**Focus**, one **Specialize**, and one **Chantry**. The opening draft adds one
**Arcane Inquiry I** as the tenth card. This example wizard has the Arcane
affinity. Give the wizard one **Veiled Archive** rumour each round; how rumours
are generated in the full game remains to be designed.

| Card | Category and cost | Effect and destination |
| --- | --- | --- |
| Wis | Resource; none | Gain 1 Wis, then discard. |
| Wealth | Resource; none | Gain 1 Wealth, then discard. |
| Focus | Spell; Arcane affinity | Commit with a challenge for +2 Essence, then discard. |
| Specialize | Activity; 1 Wis | Copy one resource or spell played earlier this round into the discard pile, then discard Specialize. Victory-path cards cannot be copied by this example effect. |
| Chantry | Location; 2 Wealth | Remain in play and gain 1 Wis at the start of each later round. |
| Veiled Archive | Personal rumour; none | Attempt an Essence 2 challenge. On success, gain a Focus into the discard pile. Destroy this rumour after the attempt or at round end. |
| Arcane Inquiry I | Victory-path card; 1 Wis | Attempt an Essence 2 challenge. On success, gain 1 Arcane step and replace this card with Arcane Inquiry II in the discard pile; otherwise discard this card unchanged. |
| Arcane Inquiry II | Victory-path card; 2 Wis | Attempt an Essence 3 challenge. On success, gain 1 Arcane step and replace this card with the next rank in the discard pile; otherwise discard this card unchanged. Later ranks and their numbers remain open. |

A hand containing **Wis**, **Specialize**, and other cards can play Wis, then
spend it on Specialize to copy a previously played resource into the discard
pile. A hand containing **Wis**, **Focus**, and **Arcane Inquiry I** can gain
1 Wis, spend it to start the Inquiry, and commit Focus for 2 Essence. That
meets the threshold, scores an Arcane step, and replaces the Inquiry. If the
player instead attempts the Inquiry without Focus, the challenge fails: the
Wis is spent, the Inquiry is discarded unchanged, and no step is gained.
Committing Focus with Veiled Archive similarly meets its Essence 2 threshold
and gains another Focus into the discard pile. Two Wealth cards and Chantry can
be played in that order to create persistent income for later rounds.

## Victory

There are many paths to victory in Grimoire.  Every wizard's progress on each of these victory paths is visible to all.  The first wizard to advance 5 steps along any victory path wins the game.

### Domination

Every time that the wizard confront and defeats another wizard, they progress along path to domination.

### Arcane

By plumbing the depths of magical knowledge, the wizard progresses along the path to an arcane victory.

### Influence

The world is still largely controlled by mundane humans, and gaining sufficient sway in their society allows a wizard to progress towards an influence victory.

### Council

By cultivating favour with their peers, a wizard can eventually find themselves progressing towards becoming the recognized leader of the wizards, and a council victory.

### Fame

By completing epic quests that appear during play, a wizard can become eternally famous, and progress towards a fame victory.

## Cards

### Common Card Attributes

***Should the concept of durable and consumable apply more generally across the board?***
***Should cards have a generic 'Noise' attribute that informs rumour generation in rivals?***

#### Cost

Cost indicates what resources and affinities need to be present to play the card from your hand.  Inability to pay the cost means that the card can not be played.

#### Longevity

All cards have one or more longevity attributes.  These inform how the card remains in play, is placed in the discard pile, and perhaps destroyed.

##### Instants

An instant card is a played and triggered immediately, then placed in the discard pile.

##### Consumables

A consumable card isn't necessarily triggered when played, but once triggered, will be discarded.

##### Durables

Durable cards remain in play and active until something specifically removes them, or an upkeep is not paid.

##### Destructibles

Destructible cards are destroyed when they would otherwise go to the discard pile.

##### Ephemeral

Ephemeral cards are cards that are destroyed at the end of your turn.

#### Visibility

How greatly the card impacts the rumour mill.

#### Upkeep

A cost that must be played by the end every turn, or the card will be discarded.

## Card Categories

All cards fall into one of the below categories.  Some of those categories are then further divided into subcategories.  The category of a card determines the parameters that can be present on the card, and the general template.
In addition to parameters, any card can include text that doesn't necessarily conform to any category rules.

### Assets

Cards that, once played, tend to remain in play, providing benefits to the wizard.  Assets tend to provide affinities and generate resources.

#### Minions

Minions are assets that represent loyal followers that are able to undertake challenges on behalf of the wizard.
Minions can become exhausted.  When exhausted, their upkeep must still be paid, but they are unable to attempt challenges, or provide affinities to the wizard.
A common way for minions to become exhausted is to undertake a challenge.  Another way is to suffer injury.
The provisional first playable draft does not yet define how minions contribute to challenges; its examples use cards committed from the hand.

Parameters:

- Upkeep: resources that must be paid at the end of each turn, or the minion is discarded
- Attributes (typical range 0-3)
  - Might: capacity for physical violence
  - Charm: ability to exert social influence
  - Acuity: aptitude for noticing things and unravelling mysteries
  - Essence: facility wielding and understanding magical phenomena
- Affinities: a set of domains the minion is familiar with

#### Artefacts

Artefacts are magical assets that are typically created by the wizard, although they can sometimes be retrieved from challenges or events, or taken from rivals.  Artefacts fall into two subcategories: Consumables, and Durables.  Consumables, upon use are discarded, while durables remain in play.  An artefact is used when the card has any impact on the game at all, other than paying upkeep.

Artefacts are always assets that belong to one of the other subcategories as well (ie. an artefact minion golem).

Parameters:

- Parameters appropriate to the secondary asset category

#### Locations

Locations are assets that represent structures or settings that fall under the dominion of the wizard.  They are placed in the play area, but not equipped to any individual minion or the wizard.  Locations typically afford affinities.

Parameters:

- Affinities: spheres of influence and knowledge that are impacted by controlling this location

### Spells

Spells are probably the most common cards in your deck.  The cost to play a spell always includes magical affinities, but often includes other resources, such as Wis.

### Designs

Design cards represent unhatched arcane plans.  These could become spells or artefacts.  The design gives some, but not all details of what will be created when the design is implemented.  Designs, once discovered, are placed in the discoverer's library.

### Rumours

Rumours are cards that are added to your hand at the start of the turn.  Rumour cards are ephemeral.  Rumours cards are typically generated by your actions, and those of your rivals.  A rumour represents a momentary opportunity to embark upon a challenge.  If not quickly seized upon, the opportunity passes.

### Resources

Resource Cards are cards that generate resources when played.  Typical resources generated would be Wealth and Wis.

### Activities

Activity cards represent basic actions that can be performed by the wizard. In the provisional first playable draft, each Activity is played for its printed effect and may name another card as a target. An Activity does not always require a paired card. Activities are one of the primary methods to tailor the wizard's deck.

#### Research Activity

The research card allows the wizard to generate designs.  This adds cards to the wizard's library.

#### Deconstruct Activity

The deconstruct card allows the wizard to destroy a card, and recover resources in the process.

#### Specialize Activity

The specialize card allows the wizard to duplicate another card.

#### Collaborate Activity

The collaborate activity allows the wizard to generate rumours and opportunities for other wizards to interact with them.

### Obstacles

Obstacle cards are automatically played at the start of the turn, when found in the wizard's hand.  Obstacle cards are cards that serve to disrupt the wizard.  Some obstacles simply occupy hand space, limiting a wizards options.  Other obstacles remain in play and impede the wizard over a longer period of time.  Many obstacles include a challenge, that when overcome, results in the obstacle's destruction.

## Deck Construction

A player's initial deck is constructed by performing a series of drafts, and adding a set of prescribed basic cards.

Spells, and wis offered in the drafting of the initial deck will necessarily applicable for affinities possessed by the wizard.

During some of the drafting phases, the wizard will be able to choose between different categories of cards.  For instance they might be offered a choice between a spell, a wis resource, or a minion card.  Each of these would be tailored to their affinities, but the wizard would be able to bias their deck towards one of those categories.

### Drafting

Pick one of several cards several times to fill your deck.

### Static

You simply start with a preset collection of cards.

### Meta

Many of the initial cards create other cards, giving the player an opportunity to tailor their deck.  This concept doesn't work super well if we are playing only a turn every x amount of time.

## Aspects and Affinities

Does it make sense to have opposed affinities?  For instance City vs Remote, or Holy vs Demonic?  Sending a demon to attempt a negotiation challenge with a priest might be impossible.

## Notes To Incorporate

- When drafting at the start, one of the selections will be for a chantry location.  A chantry is a location that generates wis income.
- When performing many types of actions that involve magic, playing additional spell cards allows you to inform the results of the action.
- When attempting a challenge, how do we add randomness to the result without screwing the player?
