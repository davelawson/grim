# Pin rules and catalogue versions for each match

Matches may wait indefinitely for a response or required choice, so a deployment
can occur during an unfinished effect. Pin each match to its rules and catalogue
versions rather than applying the latest behavior at load time: automatic
rebalance could otherwise change paid costs, pending effects, or the options
available to a player midway through play.

Deployments must preserve the behavior needed by unfinished matches or provide
an explicit migration. This deliberately carries a compatibility cost so that
resuming a match does not silently reinterpret its saved state.
