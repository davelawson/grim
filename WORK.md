# Current Work Slice

This file holds the design for one active slice of work. Before implementation
begins, fill in the goal and design below. Keep unresolved choices here while
the slice is in progress. `TODO.md` remains the backlog.

As each part is implemented and verified, move its lasting design details to
[ARCHITECTURE.md](ARCHITECTURE.md) for application and technical decisions or
[GAME.md](GAME.md) for game rules. Remove those completed details from this file
so it continues to describe only work in progress. When the slice is complete,
reset this file for the next slice.

For a design-only slice, move its accepted provisional rules to the appropriate
permanent document when the design is complete, even if implementation follows
in a later slice. Label those rules as provisional until they are tested.

## Active slice

None. The backend architecture interview and plan review are complete and
confirmed. The agreed design is recorded in
[ARCHITECTURE.md](ARCHITECTURE.md#planned-match-backend), with game rules in
[GAME.md](GAME.md), vocabulary in [CONTEXT.md](CONTEXT.md), and decisions in
[docs/adr/](docs/adr/).

## Goal

To be defined when the next implementation slice starts. The first planned
slice is backend foundations: incremental migrations, administrator storage and
authorization, lobby readiness/departure/ownership transfer, saved match setup,
and transactional one-time launch.

## Design

Use the confirmed architecture's stored-state model, delivery sequence, and
acceptance scenarios to define the next slice before implementation begins.
The complete-game milestone also requires the separate card catalogue and
unfinished rules tracked in [TODO.md](TODO.md).

## Open decisions

None in the completed backend design interview. Separate game-content design
prerequisites remain in the backlog.

## Progress

Design review complete. Application code is unchanged. Ready to define the
backend foundations implementation slice.
