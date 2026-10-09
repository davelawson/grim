package api

// MatchSummary identifies a match without exposing saved game state.
type MatchSummary struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

// LobbySummary identifies a lobby without exposing its details.
type LobbySummary struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

type ListMatchesResponse struct {
	Matches []MatchSummary `json:"matches"`
}
type ListLobbiesResponse struct {
	Lobbies []LobbySummary `json:"lobbies"`
}
