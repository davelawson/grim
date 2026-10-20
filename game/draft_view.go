package game

import (
	"errors"
	"fmt"

	"github.com/davelawson/grim/api"
)

// DraftView uses only a participant view, including a recorded receipt's view,
// so retries reveal the options for the originally returned step and revision.
func DraftView(view *View, definitions *Definitions) (*api.DraftResponse, error) {
	if view == nil || view.You == nil {
		return nil, errors.New("missing private draft view")
	}
	if definitions == nil || view.RulesVersion != definitions.catalogue.RulesVersion || view.CatalogueVersion != definitions.catalogue.CatalogueVersion {
		return nil, ErrUnsupportedVersion
	}
	response := &api.DraftResponse{
		Draft: view.You.Draft, Cards: []api.CardDesign{},
		Revision: view.Revision, Status: view.Status,
	}
	if view.Status != "setup" {
		return response, nil
	}
	offers := OpeningOffers()
	var ids []string
	switch response.Draft.Step {
	case "chantry":
		ids = offers.Chantries
	case "victory_path":
		ids = offers.VictoryPaths
	case "minion":
		ids = offers.Minions[response.Draft.Chantry]
		if len(ids) == 0 {
			return nil, errors.New("missing draft Chantry offers")
		}
	case "complete":
		return response, nil
	default:
		return nil, errors.New("unknown draft step")
	}
	catalogue := &definitions.catalogue
	for _, id := range ids {
		card, ok := catalogue.Card(id)
		if !ok {
			return nil, fmt.Errorf("missing draft card %s", id)
		}
		response.Cards = append(response.Cards, card)
	}
	return response, nil
}
