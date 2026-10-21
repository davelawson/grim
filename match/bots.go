package match

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/davelawson/grim/bot"
	"github.com/davelawson/grim/game"
)

// advanceBots runs inside the initiating human's transaction. Only the final
// state is saved; internal bot moves need no externally replayable receipts.
func advanceBots(state *game.State, definitions *game.Definitions) (result *game.State, err error) {
	defer func() {
		if err != nil {
			slog.Error("automated match transition failed", "match", state.ID, "revision", state.Revision, "error", err)
		}
	}()
	if len(state.Bots) == 0 {
		return state, nil
	}
	catalogue := definitions.Catalogue()
	random := append([]game.RandomDraw{}, state.RandomDraws...)
	turn, actions := 0, 0
	for steps := 0; steps < 20000; steps++ {
		if state.Status == "finished" {
			state.RandomDraws = random
			return state, nil
		}
		actor := ""
		if state.Status == "setup" {
			for _, id := range state.Participants {
				if _, ok := state.Bots[id]; ok && state.Wizards[id].Draft.Step != "complete" {
					actor = id
					break
				}
			}
			if actor == "" {
				state.RandomDraws = random
				return state, nil
			}
			offer, err := game.DraftView(state.View(actor), definitions)
			if err != nil {
				return nil, err
			}
			kind, id := bot.Draft(offer, catalogue)
			next, err := state.Choose(actor, kind, id)
			if err != nil {
				return nil, fmt.Errorf("automated draft: %w", err)
			}
			state = next
			continue
		}
		switch {
		case state.Choice != nil:
			actor = state.Choice.Actor
		case state.Exchange != nil:
			actor = state.Exchange.Priority
		case state.Turn != nil:
			actor = state.Turn.ActiveWizard
		}
		if _, ok := state.Bots[actor]; !ok {
			state.RandomDraws = random
			return state, nil
		}
		if state.Turn.Number != turn {
			turn, actions = state.Turn.Number, 0
		}
		var accepted *game.State
		for _, command := range bot.Candidates(state.View(actor), catalogue, actions) {
			next, err := state.Apply(actor, command)
			if err == nil {
				accepted = next
				break
			}
			if !errors.Is(err, game.ErrInvalidChoice) && !errors.Is(err, game.ErrTiming) {
				return nil, fmt.Errorf("automated command: %w", err)
			}
		}
		if accepted == nil {
			return nil, fmt.Errorf("no legal automated decision for %s at revision %d", actor, state.Revision)
		}
		random = append(random, accepted.RandomDraws...)
		state = accepted
		actions++
	}
	return nil, errors.New("automated play exceeded transition limit")
}
