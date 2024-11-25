# Grimoire Game Mechanics

Grimoire is a turn-based deck building game.  Each player takes on the role of a wizard striving to overtake their rivals, and be the first to ascend.  They develop spells and  artefacts, recruit minions to overcome challenges, in their pursuit of power.

## Basics of Play

## Victory

There are many paths to victory in Grimoire.  Every wizard's progress on each of these victory paths is visible to all.  The fist wizard to advance 5 steps along any victory path wins the game.

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

### Wizard

The Wizard card is picked at the beginning, and never drawn from the deck.  It remains in play for the entirety of the game, and provides a baseline magical aptitude and resource generation.

Parameters:

- Magical affinities
- Wis income
- Wealth income

### Assets

Cards that once played, tend to remain in play.  Assets tend to provide affinities and generate resources.

#### Minions

Minions are assets that represent loyal followers that are able to undertake challenges on behalf of the wizard.

Parameters:

- Cost: resources that are expended, and affinities that are satisfied, before the minion can be recruited
- Wage: resources that need to be expected at the end of each turn to prevent the minion from abandoning the wizard
- Attributes
  - Might: capacity for physical violence
  - Charm: ability to exert social influence
  - Acuity: aptitude for noticing things and unravelling mysteries
  - Magic: facility wielding and understanding magical phenomena

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

Activity cards represent basic actions that can be performed by the wizard.  Activity cards are played alongside another card, and are one of the primary methods to tailor the wizard's deck.

#### Research Activity

The research card allows the wizard to generate designs.  This adds cards to the wizard's library.

#### Deconstruct Activity

The deconstruct card allows the wizard to destroy a card, and recover resources in the process.

#### Specialize Activity

The specialize card allows the wizard to duplicate another card.

#### Collaborate Activity

The collaborate activity allows the wizard to generate rumours and opportunities for other wizards to interact with them.

### Obstacles

Obstacle cards are automatically played at the start of the turn, when found in the wizard's hand.  Obstacle cards are cards that serve to disrupt the wizard.  Some obstacles simply occupy handspace, limiting a wizards options.  Other obstacles remain in play and impede the wizard over a longer period of time.  Many obstacles include a challenge, that when overcome, results in the obstacle's destruction.

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
