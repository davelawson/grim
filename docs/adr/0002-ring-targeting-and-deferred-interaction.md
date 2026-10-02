# Ring targeting and deferred interaction

The prohibition on live responses is superseded by
[Allow player responses to directed actions](0010-allow-player-responses.md).
The shared phase after all turn submissions is superseded by
[Round-robin turns with player responses](0011-round-robin-turns.md).
Ring targeting is retained.

To limit coordinated attacks on one rival while supporting play without live
responses, Grimoire arranges wizards in a Ring: each wizard is the Aggressor
of the wizard on their right and the Target of the wizard on their left;
chosen offensive recipients are restricted to those two neighbours, with most
offensive actions directed at the Target and some at the Aggressor.
Player-to-player impacts resolve in a shared phase after every wizard has
completed their independent actions and submitted their turn for the round,
with no intervening negotiation or response decisions.
This trades unrestricted targeting and reactive negotiation for bounded
rivalries and interactions that resolve from submitted choices.
