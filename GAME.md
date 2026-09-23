# Grimoire Game Mechanics

Grimoire is a turn-based deck building game.  Each player takes on the role of a wizard striving to overtake their rivals, and be the first to ascend.  They develop spells and  artifacts, recruit minions to overcome challenges, in their pursuit of power.

## Basics of Play

### Provisional first playable draft

The rules in this section are a starting point for testing the flow of a game.
Card costs, hand size, and challenge thresholds are illustrative balance values,
not final numbers. The intended loop is to improve the deck and establish
assets, then use those cards to pursue rumours and visible victory progress.
Each shared round contains every wizard's simultaneous turn. This draft
covers independent actions; direct confrontation and shared contested
opportunities need later rules.

A round works as follows:

1. Each wizard has the resources generated in the previous round available,
   draws five cards, and receives one personal rumour card. A new wizard begins
   with no available resources. If the draw pile runs out during the draw,
   shuffle the discard pile to continue. Obstacles drawn into the hand resolve
   automatically, as described below.
2. Each player privately submits a set of card plays. A card can be used only
   once that round. Cards supply the actions; there is no separate action-slot
   limit. A play can use one card or combine cards for one action, such as a
   challenge with support. The total costs of the submitted plays and upkeep
   must fit the resources already available at the start of the round.
3. Once all plans are submitted, resolve the plays concurrently. Each action
   uses the state and resources available at the start of the round. An action's
   result generally cannot fund or modify another action in the same round;
   cards explicitly combined for one action are the exception. In this draft,
   players' actions are independent, so no cross-player order is needed.
   Generated resources are set aside for the following round.
4. Pay upkeep from this round's available resources. Move played non-durable
   cards and unused hand cards to the discard pile; newly gained cards enter
   the discard pile too. Destroy ephemeral cards instead. Unspent available
   resources expire, and generated resources become available next round.

Wealth and Wis can be spent only in the round after they are generated. Any
unspent amount expires at the end of that spending round. An affinity is a
prerequisite, not a resource that is spent. Durable assets remain in play and
can generate future income. Cards in the hand, rather than a fixed number of
actions, determine how much a wizard can do.

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
advance. A successful attempt adds one visible step when the round resolves,
removes the attempted card, and puts its costlier next rank into the discard
pile. Other cards may offer alternate paths or change the wizard's primary
path; those effects are not designed here. Reaching the fifth step wins the
game instead of creating another rank.

### Starting deck and card interactions

Every wizard begins with exactly ten cards. Six are fixed; four come from draft
choices or are linked to a choice:

| Starting card | Count | How it enters the deck |
| --- | ---: | --- |
| Chantry | 1 | Draft one of three Chantries. |
| Element | 1 | Add the Element linked to the chosen Chantry. |
| Research Spell | 1 | Fixed starting card. |
| Research Artifact | 1 | Fixed starting card. |
| Basic Wis | 2 | Fixed starting cards. |
| Victory challenge | 1 | Draft a card for a primary victory path. |
| Basic Wealth | 1 | Fixed starting card. |
| Minion | 1 | Draft a Minion themed to match the chosen Chantry. |
| Arcane Focus | 1 | Fixed starting card. |
| **Total** | **10** | |

The three Chantry choices, their linked Elements, and the Minion offers have
not yet been named or balanced. A Chantry is a Location that generates Wis;
its other effects and exact cost depend on the chosen card. The Minion's theme
follows that choice, but its challenge contribution is still to be designed.
**Arcane Inquiry I** remains an illustrative choice for the victory challenge
slot, with **Arcane Inquiry II** as its replacement after a success.

A basic Wis or Wealth card generates one unit of its named resource for the
next round when played. Research Spell creates a new spell card in the discard pile;
Research Artifact creates a new artifact card there. When an Element is
committed with either Research card, it supplies its affinity and the created
card is associated with that Element. Without an Element, the result is
unaligned. An Element may instead be committed with one other card to provide
its affinity for that action. It cannot be used for another action in the
same round.

Arcane Focus is an Activity card that can be committed with one other card. It
reduces either that card's Wis cost or its Wealth cost by 1, chosen when the
play is submitted, to a minimum of zero. The paired cards form one action and
are each used once. Arcane Focus changes that action's cost; it does not make
resources generated by another action available in the same round.

For example, pairing an Element with Research Spell puts an Element-associated
spell into the discard pile. Pairing Arcane Focus with an Arcane Inquiry I that
costs 1 Wis reduces that attempt's Wis cost to zero. It does not add any
challenge capability, so the attempt still needs appropriate support to
succeed. The way a starting Minion supplies such support remains open.

## Victory

There are many paths to victory in Grimoire. Every wizard's progress on each path is visible to all. A wizard wins by advancing 5 steps along any one path. If multiple wizards reach a fifth step in the same round, the tie rule remains to be designed.

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

#### Artifacts

Artifacts are magical assets that are typically created by the wizard, although they can sometimes be retrieved from challenges or events, or taken from rivals.  Artifacts fall into two subcategories: Consumables, and Durables.  Consumables, upon use are discarded, while durables remain in play.  An artifact is used when the card has any impact on the game at all, other than paying upkeep.

Artifacts are always assets that belong to one of the other subcategories as well (ie. an artifact minion golem).

Parameters:

- Parameters appropriate to the secondary asset category

#### Locations

Locations are assets that represent structures or settings that fall under the dominion of the wizard.  They are placed in the play area, but not equipped to any individual minion or the wizard.  Locations typically afford affinities.

Parameters:

- Affinities: spheres of influence and knowledge that are impacted by controlling this location

### Spells

Spells are probably the most common cards in your deck.  The cost to play a spell always includes magical affinities, but often includes other resources, such as Wis.

### Designs

Design cards represent unhatched arcane plans. These could become spells or artifacts. The design gives some, but not all details of what will be created when the design is implemented. Designs, once discovered, are placed in the discoverer's library. The two starting Research cards create completed spells or artifacts directly, rather than Designs.

### Rumours

Rumours are cards that are added to your hand at the start of the turn.  Rumour cards are ephemeral.  Rumours cards are typically generated by your actions, and those of your rivals.  A rumour represents a momentary opportunity to embark upon a challenge.  If not quickly seized upon, the opportunity passes.

### Resources

Resource Cards are cards that generate resources when played. Typical resources generated would be Wealth and Wis.

### Elements

An Element card represents one elemental affinity. Commit it with a single other card to supply that affinity for the paired action. If that action is Research Spell or Research Artifact, the created card is associated with the Element. Each Element can support only one paired action in a round.

### Victory Challenges

A drafted victory challenge provides a repeatable opportunity to advance one path. On success, it scores a step and is replaced by a more advanced version, as described above.

### Activities

Activity cards represent basic actions that can be performed by the wizard. In the provisional first playable draft, each Activity is played for its printed effect and may name another card as a target. An Activity does not always require a paired card. Activities are one of the primary methods to tailor the wizard's deck.

#### Research Activity

The two starting Research cards create new cards directly in the discard pile. Pairing an Element with Research determines the new card's association; without an Element, the result is unaligned. How the specific new card is selected remains to be designed.

##### Research Spell

Creates a new spell card and adds it to the wizard's discard pile.

##### Research Artifact

Creates a new artifact card and adds it to the wizard's discard pile.

#### Arcane Focus

Commit Arcane Focus with one other card to reduce either that card's Wis or Wealth cost by 1, to a minimum of zero. It affects only that paired action.

#### Deconstruct Activity

The deconstruct card allows the wizard to destroy a card, and recover resources in the process.

#### Specialize Activity

The specialize card allows the wizard to duplicate another card.

#### Collaborate Activity

The collaborate activity allows the wizard to generate rumours and opportunities for other wizards to interact with them.

### Obstacles

Obstacle cards are automatically played at the start of the turn, when found in the wizard's hand.  Obstacle cards are cards that serve to disrupt the wizard.  Some obstacles simply occupy hand space, limiting a wizards options.  Other obstacles remain in play and impede the wizard over a longer period of time.  Many obstacles include a challenge, that when overcome, results in the obstacle's destruction.

## Deck Construction

The ten-card starting deck combines fixed cards with a Chantry, a victory challenge, and a Minion chosen during the opening draft. Choosing a Chantry also adds its linked Element. The full draft procedure and the identities of the offered cards remain to be designed; the card counts and relationships above define the current starting structure.

### Drafting

The opening draft includes one Chantry chosen from three offers, one victory challenge, and one Minion themed by the chosen Chantry. The Chantry choice also determines the Element card. Details of the victory and Minion offers and the draft order remain open.

### Static

Each starting deck also receives Research Spell, Research Artifact, two basic Wis, one basic Wealth, and Arcane Focus.

### Meta

The two starting Research cards create new spells or artifacts, allowing players to develop their deck over successive rounds. The pace of this growth remains to be tested.

## Aspects and Affinities

Does it make sense to have opposed affinities?  For instance City vs Remote, or Holy vs Demonic?  Sending a demon to attempt a negotiation challenge with a priest might be impossible.

## Notes To Incorporate

- When performing many types of actions that involve magic, playing additional spell cards allows you to inform the results of the action.
- When attempting a challenge, how do we add randomness to the result without screwing the player?
