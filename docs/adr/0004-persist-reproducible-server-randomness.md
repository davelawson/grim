# Persist reproducible server randomness

The server supplies every random outcome and stores private random-generator
state alongside each match snapshot, committing generated outcomes with the
accepted command. This lets a match continue consistently after restart and
makes its randomness reproducible, unlike unrecorded process-local or
client-supplied random outcomes.

Generator state is excluded from player views. Reproduction requires the
random behavior associated with the match's pinned rules version. Retain the
initial seed and use diagnostic logs to reconstruct game flow when necessary;
a complete replay store and replay API are outside the first milestone.

Operator-accessible diagnostic logs include accepted command inputs and
choices, match/actor/revision identifiers, pinned versions, and random outcomes.
They remain separate from player views and exclude credentials and tokens;
no stored public player history or player-history endpoint is required.
