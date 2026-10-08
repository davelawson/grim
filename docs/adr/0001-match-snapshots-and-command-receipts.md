# Persist match snapshots with command receipts

The first playable backend must resume at any response window or required
choice, including choices inside an unfinished resolution. Store authoritative
state as a versioned match document in SQLite, with membership and searchable
metadata in ordinary tables; this keeps the Stack, choices, and continuation
together while avoiding the complexity of reconstructing state from events or
coordinating individual rows for every game object.

Commands carry a unique request ID and expected match revision. Save each
accepted command's resulting state and outcome in one transaction, return the
recorded outcome for an identical retry, and reject stale conflicting commands.
This prevents lost HTTP responses from repeating paid actions and competing
requests from overwriting a newer state.

The saved document requires deliberate schema compatibility as the engine
evolves. Command receipts support retries; they are not an event-sourced model
or a guarantee of a complete replay history.
