package match

import (
	"database/sql"
	"encoding/json"
	"errors"
	"strings"

	"github.com/davelawson/grim/game"
	"github.com/davelawson/grim/util"
	"github.com/google/uuid"
)

type gameplayOperation string

type commandDiagnostics struct {
	Replayed bool
	Random   []game.RandomDraw
}

// The operation is part of receipt identity while request IDs remain shared
// across all typed routes. The normalized engine input includes path-bound choice IDs.
type gameplayReceiptInput struct {
	Operation gameplayOperation `json:"operation"`
	Command   game.Command      `json:"command"`
}

// submitGameplay saves an accepted transition and receipt in the caller's
// transaction. Retry lookup precedes live timing checks, which can have changed
// since the original input was accepted.
func (s *Service) submitGameplay(tx *sql.Tx, id, actorID string, c game.Command, operation gameplayOperation, diagnostics *commandDiagnostics) (*game.View, error) {
	parsed, err := uuid.Parse(c.RequestID)
	if err != nil || !strings.EqualFold(c.RequestID, parsed.String()) || c.ExpectedRevision == nil || *c.ExpectedRevision < 0 {
		return nil, util.ErrInvalidRequest
	}
	c.RequestID = parsed.String()
	stored, err := s.repo.load(tx, id)
	if err != nil {
		return nil, err
	}
	member, err := s.repo.isParticipant(tx, id, actorID)
	if err != nil {
		return nil, err
	}
	if !member {
		return nil, util.ErrForbidden
	}
	input, err := json.Marshal(gameplayReceiptInput{Operation: operation, Command: c})
	if err != nil {
		return nil, err
	}
	if view, err := s.repo.gameplayReceipt(tx, id, actorID, c, string(input)); err != nil || view != nil {
		diagnostics.Replayed = view != nil
		return view, err
	}
	state, err := s.decode(tx, stored)
	if err != nil {
		return nil, err
	}
	if state.Revision != *c.ExpectedRevision {
		return nil, util.ErrStaleRevision
	}
	if state.Status == "playing" {
		if strings.HasPrefix(string(operation), "reactions/") {
			if state.Exchange == nil || state.Choice != nil || state.Resolving != nil {
				return nil, gameplayError(game.ErrTiming)
			}
		} else if c.Type != "resolve_choice" && state.Exchange != nil {
			return nil, gameplayError(game.ErrTiming)
		}
	}
	if state.Status == "finished" {
		return nil, gameplayError(game.ErrFinished)
	}
	if err := state.ValidateOperation(actorID, c, string(operation)); err != nil {
		return nil, gameplayError(err)
	}
	next, err := state.Apply(actorID, c)
	if err != nil {
		return nil, gameplayError(err)
	}
	next, err = advanceBots(next, next.Definitions())
	if err != nil {
		return nil, err
	}
	if err := s.repo.save(tx, next, state.Revision); err != nil {
		return nil, err
	}
	view := next.View(actorID)
	if err := s.repo.saveCommandReceipt(tx, id, actorID, c.RequestID, string(input), view); err != nil {
		return nil, err
	}
	diagnostics.Random = next.RandomDraws
	return view, nil
}

func (r *Repo) gameplayReceipt(tx *sql.Tx, id, actorID string, c game.Command, input string) (*game.View, error) {
	receipt, err := r.loadCommandReceipt(tx, id, actorID, c.RequestID)
	if err != nil || receipt == nil {
		return nil, err
	}
	if receipt.input != input {
		return nil, requestIDConflict()
	}
	var view game.View
	if err := json.Unmarshal(receipt.response, &view); err != nil {
		return nil, err
	}
	return &view, nil
}

func gameplayError(err error) error {
	switch {
	case errors.Is(err, game.ErrNotParticipant):
		return util.ErrForbidden
	case errors.Is(err, game.ErrUnsupportedVersion):
		return util.ErrUnsupportedState
	case errors.Is(err, game.ErrInvalidCommand):
		return util.ErrInvalidRequest
	case errors.Is(err, game.ErrInvalidChoice):
		return &util.APIError{Status: 400, Code: "invalid_choice", Message: "The selected cards, targets, or payments are not legal"}
	case errors.Is(err, game.ErrTiming):
		return &util.APIError{Status: 409, Code: "invalid_timing", Message: "This command is unavailable at the current input boundary"}
	case errors.Is(err, game.ErrPriority):
		return &util.APIError{Status: 409, Code: "wrong_priority", Message: "Another wizard must act"}
	case errors.Is(err, game.ErrStaleChoice):
		return &util.APIError{Status: 409, Code: "stale_choice", Message: "This choice is no longer pending"}
	case errors.Is(err, game.ErrFinished):
		return &util.APIError{Status: 409, Code: "match_finished", Message: "The match has finished"}
	default:
		return err
	}
}
