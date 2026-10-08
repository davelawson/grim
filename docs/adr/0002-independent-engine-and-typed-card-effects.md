# Isolate the game engine and implement typed card effects

Game rules belong in a Go state machine independent of HTTP and SQLite, with
the application service handling authentication and persistence. The engine
advances through automatic steps until it needs player input or the match ends,
saving continuations as explicit state; embedding resolution in
request handlers or suspended functions would tie resumability to a running
request or process.

Card definitions contain structured properties and references to typed Go
effect handlers, while individual card instances hold ownership, location,
and changing state. This supports the initial small catalogue without the
language, interpreter, and compatibility burden of arbitrary card scripting;
new behavior can require an engine release.
