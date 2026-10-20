# Store the shared card catalogue in one versioned table

Card designs are shared content, distinct from the individual cards in a
wizard's deck. Store the complete catalogue in one `card_designs` table, with
columns for ordinary properties and individual stats, and JSON only for lists
such as categories, effect identifiers, abilities, and Challenges. SQL scripts
manage content; there is no editing API or scripting language. Typed effect
implementations and draft selection rules remain Go code.

The application loads and validates each match's pinned catalogue inside its
transaction, then injects immutable definitions into the SQL-independent engine.
Definitions are runtime dependencies, not additional match snapshot data. The
initial seed preserves every existing design, value, and catalogue position;
positions matter to reproducible random choices. New matches select the version
configured by `GRIM_CATALOGUE_VERSION`, defaulting to `playable-v1`.

A version freezes when its first match is saved, including matches later ended
or soft-deleted. Triggers prohibit inserting, updating, or deleting its rows;
changes create a new catalogue version. A separate additive installation script
preserves existing matches and receipts. Startup validates installed data and
never seeds, resets, migrates, or falls back to Go definitions. This preserves
[version compatibility](0003-pin-match-rules-and-catalogue.md) while allowing
catalogue data changes without rebuilding the server.
