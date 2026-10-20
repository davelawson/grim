package match

import (
	"encoding/json"
	"github.com/davelawson/grim/api"
	"github.com/davelawson/grim/game"
	"github.com/davelawson/grim/model"
	"github.com/davelawson/grim/util"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type Controller struct{ service *ServiceFacade }

func NewController(service *ServiceFacade) *Controller { return &Controller{service: service} }

// Get returns the authenticated participant's private and public match view.
// @Summary Get a match
// @Security ApiKeyAuth
// @Tags match
// @Produce json
// @Param id path string true "Match ID"
// @Success 200 {object} api.MatchResponse
// @Failure 401,403,404,409,500 {object} api.ErrorResponse
// @Router /match/{id} [get]
func (c *Controller) Get(ctx *gin.Context) {
	actor := ctx.MustGet("reqUser").(*model.User)
	view, err := c.service.Get(ctx.Param("id"), actor.Id)
	if err != nil {
		util.WriteAPIError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, api.MatchResponse{Match: view})
}

// GetDraft returns the caller's current step and draftable card definitions.
// @Summary Get current draft options
// @Security ApiKeyAuth
// @Tags match
// @Produce json
// @Param id path string true "Match ID"
// @Success 200 {object} api.DraftResponse
// @Failure 401,403,404,409,500 {object} api.ErrorResponse
// @Router /match/{id}/draft [get]
func (c *Controller) GetDraft(ctx *gin.Context) {
	actor := ctx.MustGet("reqUser").(*model.User)
	response, err := c.service.GetDraft(ctx.Param("id"), actor.Id)
	if err != nil {
		util.WriteAPIError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, response)
}

// SubmitDraft locks one opening choice and returns the next draft options.
// @Summary Choose an opening card
// @Security ApiKeyAuth
// @Tags match
// @Accept json
// @Produce json
// @Param id path string true "Match ID"
// @Param command body api.DraftPickRequest true "Draft choice; expectedrevision is required"
// @Success 200 {object} api.DraftResponse
// @Failure 400,401,403,404,409,500 {object} api.ErrorResponse
// @Router /match/{id}/draft [post]
func (c *Controller) SubmitDraft(ctx *gin.Context) {
	var command api.DraftPickRequest
	decoder := json.NewDecoder(ctx.Request.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&command); err != nil {
		util.WriteAPIError(ctx, util.ErrInvalidRequest)
		return
	}
	if err := decoder.Decode(new(any)); err != io.EOF {
		util.WriteAPIError(ctx, util.ErrInvalidRequest)
		return
	}
	actor := ctx.MustGet("reqUser").(*model.User)
	response, err := c.service.SubmitDraft(ctx.Param("id"), actor.Id, command)
	if err != nil {
		util.WriteAPIError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, response)
}

// Catalogue returns the definitions pinned to this participant's match.
// @Summary Get the match card catalogue
// @Security ApiKeyAuth
// @Tags match
// @Produce json
// @Param id path string true "Match ID"
// @Success 200 {object} api.CatalogueResponse
// @Failure 401,403,404,409,500 {object} api.ErrorResponse
// @Router /match/{id}/catalogue [get]
func (c *Controller) Catalogue(ctx *gin.Context) {
	actor := ctx.MustGet("reqUser").(*model.User)
	catalogue, err := c.service.Catalogue(ctx.Param("id"), actor.Id)
	if err != nil {
		util.WriteAPIError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, api.CatalogueResponse{Catalogue: catalogue})
}

// End soft-deletes a match while preserving its saved state.
// @Summary Administratively end a match
// @Description Requires current stored administrator permission; already-deleted matches succeed unchanged.
// @Security ApiKeyAuth
// @Tags admin
// @Produce json
// @Param id path string true "Match ID"
// @Success 204 "No Content"
// @Failure 401,403,404,500 {object} api.ErrorResponse
// @Router /admin/match/{id}/end [post]
func (c *Controller) End(ctx *gin.Context) {
	actor := ctx.MustGet("reqUser").(*model.User)
	if err := c.service.End(ctx.Param("id"), actor.Id); err != nil {
		util.WriteAPIError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func submitTyped[T any](c *Controller, ctx *gin.Context, convert func(T) game.Command, operation gameplayOperation) {
	var request T
	decoder := json.NewDecoder(ctx.Request.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		util.WriteAPIError(ctx, util.ErrInvalidRequest)
		return
	}
	if err := decoder.Decode(new(any)); err != io.EOF {
		util.WriteAPIError(ctx, util.ErrInvalidRequest)
		return
	}
	if err := binding.Validator.ValidateStruct(request); err != nil {
		util.WriteAPIError(ctx, util.ErrInvalidRequest)
		return
	}
	command := convert(request)
	if strings.HasPrefix(string(operation), "choices/") {
		command.ChoiceID = ctx.Param("choiceid")
	}
	actor := ctx.MustGet("reqUser").(*model.User)
	view, err := c.service.submitGameplay(ctx.Param("id"), actor.Id, command, operation)
	if err != nil {
		util.WriteAPIError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, api.MatchResponse{Match: view})
}

// PlayCard submits play-card.
// @Summary play-card
// @Security ApiKeyAuth
// @Tags match
// @Accept json
// @Produce json
// @Param id path string true "Match ID"
// @Param input body api.PlayCardRequest true "play-card input; requestid and expectedrevision are required"
// @Success 200 {object} api.MatchResponse
// @Failure 400,401,403,404,409,500 {object} api.ErrorResponse
// @Router /match/{id}/actions/play-card [post]
func (c *Controller) PlayCard(ctx *gin.Context) {
	submitTyped(c, ctx, func(r api.PlayCardRequest) game.Command {
		command := game.Command{RequestID: r.RequestID, ExpectedRevision: r.ExpectedRevision}
		command.Type = "play_card"
		command.CardID = r.CardID
		command.Targets = r.Targets
		command.Support = r.Support
		command.Payment = r.Payment
		command.Discard = r.Discard
		return command
	}, "actions/play-card")
}

// Attack submits attack.
// @Summary attack
// @Security ApiKeyAuth
// @Tags match
// @Accept json
// @Produce json
// @Param id path string true "Match ID"
// @Param input body api.AttackRequest true "attack input; requestid and expectedrevision are required"
// @Success 200 {object} api.MatchResponse
// @Failure 400,401,403,404,409,500 {object} api.ErrorResponse
// @Router /match/{id}/actions/attack [post]
func (c *Controller) Attack(ctx *gin.Context) {
	submitTyped(c, ctx, func(r api.AttackRequest) game.Command {
		command := game.Command{RequestID: r.RequestID, ExpectedRevision: r.ExpectedRevision}
		command.Type = "play_card"
		command.CardID = r.CardID
		command.Minions = r.Minions
		command.Support = r.Support
		command.Payment = r.Payment
		return command
	}, "actions/attack")
}

// AttemptChallenge submits attempt-challenge.
// @Summary attempt-challenge
// @Security ApiKeyAuth
// @Tags match
// @Accept json
// @Produce json
// @Param id path string true "Match ID"
// @Param input body api.AttemptChallengeRequest true "attempt-challenge input; requestid and expectedrevision are required"
// @Success 200 {object} api.MatchResponse
// @Failure 400,401,403,404,409,500 {object} api.ErrorResponse
// @Router /match/{id}/actions/attempt-challenge [post]
func (c *Controller) AttemptChallenge(ctx *gin.Context) {
	submitTyped(c, ctx, func(r api.AttemptChallengeRequest) game.Command {
		command := game.Command{RequestID: r.RequestID, ExpectedRevision: r.ExpectedRevision}
		command.Type = "attempt_challenge"
		command.CardID = r.CardID
		command.ChallengeID = r.ChallengeID
		command.Minions = r.Minions
		command.Support = r.Support
		command.Payment = r.Payment
		return command
	}, "actions/attempt-challenge")
}

// ActivateAbility submits activate-ability.
// @Summary activate-ability
// @Security ApiKeyAuth
// @Tags match
// @Accept json
// @Produce json
// @Param id path string true "Match ID"
// @Param input body api.ActivateAbilityRequest true "activate-ability input; requestid and expectedrevision are required"
// @Success 200 {object} api.MatchResponse
// @Failure 400,401,403,404,409,500 {object} api.ErrorResponse
// @Router /match/{id}/actions/activate-ability [post]
func (c *Controller) ActivateAbility(ctx *gin.Context) {
	submitTyped(c, ctx, func(r api.ActivateAbilityRequest) game.Command {
		command := game.Command{RequestID: r.RequestID, ExpectedRevision: r.ExpectedRevision}
		command.Type = "activate_ability"
		command.CardID = r.CardID
		command.AbilityID = r.AbilityID
		command.Support = r.Support
		command.Payment = r.Payment
		return command
	}, "actions/activate-ability")
}

// ResearchElement submits research-element.
// @Summary research-element
// @Security ApiKeyAuth
// @Tags match
// @Accept json
// @Produce json
// @Param id path string true "Match ID"
// @Param input body api.ResearchElementRequest true "research-element input; requestid and expectedrevision are required"
// @Success 200 {object} api.MatchResponse
// @Failure 400,401,403,404,409,500 {object} api.ErrorResponse
// @Router /match/{id}/actions/research-element [post]
func (c *Controller) ResearchElement(ctx *gin.Context) {
	submitTyped(c, ctx, func(r api.ResearchElementRequest) game.Command {
		command := game.Command{RequestID: r.RequestID, ExpectedRevision: r.ExpectedRevision}
		command.Type = "play_card"
		command.CardID = r.CardID
		command.Element = r.Element
		command.Mode = "element"
		return command
	}, "actions/research-element")
}

// ResearchSpell submits research-spell.
// @Summary research-spell
// @Security ApiKeyAuth
// @Tags match
// @Accept json
// @Produce json
// @Param id path string true "Match ID"
// @Param input body api.ResearchSpellRequest true "research-spell input; requestid and expectedrevision are required"
// @Success 200 {object} api.MatchResponse
// @Failure 400,401,403,404,409,500 {object} api.ErrorResponse
// @Router /match/{id}/actions/research-spell [post]
func (c *Controller) ResearchSpell(ctx *gin.Context) {
	submitTyped(c, ctx, func(r api.ResearchSpellRequest) game.Command {
		command := game.Command{RequestID: r.RequestID, ExpectedRevision: r.ExpectedRevision}
		command.Type = "play_card"
		command.CardID = r.CardID
		command.Mode = "discover"
		for _, id := range r.Elements {
			command.Support = append(command.Support, api.CardSupport{CardID: id})
		}
		return command
	}, "actions/research-spell")
}

// HealWithFocus submits heal-with-focus.
// @Summary heal-with-focus
// @Security ApiKeyAuth
// @Tags match
// @Accept json
// @Produce json
// @Param id path string true "Match ID"
// @Param input body api.HealWithFocusRequest true "heal-with-focus input; requestid and expectedrevision are required"
// @Success 200 {object} api.MatchResponse
// @Failure 400,401,403,404,409,500 {object} api.ErrorResponse
// @Router /match/{id}/actions/heal-with-focus [post]
func (c *Controller) HealWithFocus(ctx *gin.Context) {
	submitTyped(c, ctx, func(r api.HealWithFocusRequest) game.Command {
		command := game.Command{RequestID: r.RequestID, ExpectedRevision: r.ExpectedRevision}
		command.Type = "play_card"
		command.CardID = r.CardID
		command.Support = r.Support
		command.Payment = r.Payment
		command.Mode = "heal"
		return command
	}, "actions/heal-with-focus")
}

// EndTurn submits end-turn.
// @Summary end-turn
// @Security ApiKeyAuth
// @Tags match
// @Accept json
// @Produce json
// @Param id path string true "Match ID"
// @Param input body api.EndTurnRequest true "end-turn input; requestid and expectedrevision are required"
// @Success 200 {object} api.MatchResponse
// @Failure 400,401,403,404,409,500 {object} api.ErrorResponse
// @Router /match/{id}/actions/end-turn [post]
func (c *Controller) EndTurn(ctx *gin.Context) {
	submitTyped(c, ctx, func(r api.EndTurnRequest) game.Command {
		command := game.Command{RequestID: r.RequestID, ExpectedRevision: r.ExpectedRevision}
		command.Type = "end_turn"
		return command
	}, "actions/end-turn")
}

// Maintenance submits maintenance.
// @Summary maintenance
// @Security ApiKeyAuth
// @Tags match
// @Accept json
// @Produce json
// @Param id path string true "Match ID"
// @Param choiceid path string true "Pending choice ID"
// @Param input body api.MaintenanceRequest true "maintenance input; requestid and expectedrevision are required"
// @Success 200 {object} api.MatchResponse
// @Failure 400,401,403,404,409,500 {object} api.ErrorResponse
// @Router /match/{id}/choices/{choiceid}/maintenance [post]
func (c *Controller) Maintenance(ctx *gin.Context) {
	submitTyped(c, ctx, func(r api.MaintenanceRequest) game.Command {
		command := game.Command{RequestID: r.RequestID, ExpectedRevision: r.ExpectedRevision}
		command.Type = "resolve_choice"
		command.Payment = r.Payment
		command.Selections = r.Keep
		return command
	}, "choices/maintenance")
}

// Kneel submits kneel.
// @Summary kneel
// @Security ApiKeyAuth
// @Tags match
// @Accept json
// @Produce json
// @Param id path string true "Match ID"
// @Param choiceid path string true "Pending choice ID"
// @Param input body api.KneelRequest true "kneel input; requestid and expectedrevision are required"
// @Success 200 {object} api.MatchResponse
// @Failure 400,401,403,404,409,500 {object} api.ErrorResponse
// @Router /match/{id}/choices/{choiceid}/kneel [post]
func (c *Controller) Kneel(ctx *gin.Context) {
	submitTyped(c, ctx, func(r api.KneelRequest) game.Command {
		command := game.Command{RequestID: r.RequestID, ExpectedRevision: r.ExpectedRevision}
		command.Type = "resolve_choice"
		command.Selections = []string{"decline"}
		if *r.Kneel {
			command.Selections[0] = "kneel"
		}
		return command
	}, "choices/kneel")
}

// ResolveEffect submits resolve-effect.
// @Summary resolve-effect
// @Security ApiKeyAuth
// @Tags match
// @Accept json
// @Produce json
// @Param id path string true "Match ID"
// @Param choiceid path string true "Pending choice ID"
// @Param input body api.ResolveEffectRequest true "resolve-effect input; requestid and expectedrevision are required"
// @Success 200 {object} api.MatchResponse
// @Failure 400,401,403,404,409,500 {object} api.ErrorResponse
// @Router /match/{id}/choices/{choiceid}/resolve-effect [post]
func (c *Controller) ResolveEffect(ctx *gin.Context) {
	submitTyped(c, ctx, func(r api.ResolveEffectRequest) game.Command {
		command := game.Command{RequestID: r.RequestID, ExpectedRevision: r.ExpectedRevision}
		command.Type = "resolve_choice"
		command.Selections = []string{r.EffectID}
		return command
	}, "choices/resolve-effect")
}

// SkipEffect submits skip-effect.
// @Summary skip-effect
// @Security ApiKeyAuth
// @Tags match
// @Accept json
// @Produce json
// @Param id path string true "Match ID"
// @Param choiceid path string true "Pending choice ID"
// @Param input body api.SkipEffectRequest true "skip-effect input; requestid and expectedrevision are required"
// @Success 200 {object} api.MatchResponse
// @Failure 400,401,403,404,409,500 {object} api.ErrorResponse
// @Router /match/{id}/choices/{choiceid}/skip-effect [post]
func (c *Controller) SkipEffect(ctx *gin.Context) {
	submitTyped(c, ctx, func(r api.SkipEffectRequest) game.Command {
		command := game.Command{RequestID: r.RequestID, ExpectedRevision: r.ExpectedRevision}
		command.Type = "resolve_choice"
		command.Selections = []string{"skip:" + r.EffectID}
		return command
	}, "choices/skip-effect")
}

// UndertakeHunt submits undertake-hunt.
// @Summary undertake-hunt
// @Security ApiKeyAuth
// @Tags match
// @Accept json
// @Produce json
// @Param id path string true "Match ID"
// @Param choiceid path string true "Pending choice ID"
// @Param input body api.UndertakeHuntRequest true "undertake-hunt input; requestid and expectedrevision are required"
// @Success 200 {object} api.MatchResponse
// @Failure 400,401,403,404,409,500 {object} api.ErrorResponse
// @Router /match/{id}/choices/{choiceid}/undertake-hunt [post]
func (c *Controller) UndertakeHunt(ctx *gin.Context) {
	submitTyped(c, ctx, func(r api.UndertakeHuntRequest) game.Command {
		command := game.Command{RequestID: r.RequestID, ExpectedRevision: r.ExpectedRevision}
		command.Type = "resolve_choice"
		command.Minions = r.Minions
		command.Support = r.Support
		command.Payment = r.Payment
		command.Selections = []string{r.EffectID}
		return command
	}, "choices/undertake-hunt")
}

// ChooseResearch submits research.
// @Summary research
// @Security ApiKeyAuth
// @Tags match
// @Accept json
// @Produce json
// @Param id path string true "Match ID"
// @Param choiceid path string true "Pending choice ID"
// @Param input body api.ChooseResearchRequest true "research input; requestid and expectedrevision are required"
// @Success 200 {object} api.MatchResponse
// @Failure 400,401,403,404,409,500 {object} api.ErrorResponse
// @Router /match/{id}/choices/{choiceid}/research [post]
func (c *Controller) ChooseResearch(ctx *gin.Context) {
	submitTyped(c, ctx, func(r api.ChooseResearchRequest) game.Command {
		command := game.Command{RequestID: r.RequestID, ExpectedRevision: r.ExpectedRevision}
		command.Type = "resolve_choice"
		command.Selections = []string{r.DesignID}
		return command
	}, "choices/research")
}

// EndDiscard submits end-discard.
// @Summary end-discard
// @Security ApiKeyAuth
// @Tags match
// @Accept json
// @Produce json
// @Param id path string true "Match ID"
// @Param choiceid path string true "Pending choice ID"
// @Param input body api.EndDiscardRequest true "end-discard input; requestid and expectedrevision are required"
// @Success 200 {object} api.MatchResponse
// @Failure 400,401,403,404,409,500 {object} api.ErrorResponse
// @Router /match/{id}/choices/{choiceid}/end-discard [post]
func (c *Controller) EndDiscard(ctx *gin.Context) {
	submitTyped(c, ctx, func(r api.EndDiscardRequest) game.Command {
		command := game.Command{RequestID: r.RequestID, ExpectedRevision: r.ExpectedRevision}
		command.Type = "resolve_choice"
		command.Selections = r.Cards
		return command
	}, "choices/end-discard")
}

// HandDiscard submits hand-discard.
// @Summary hand-discard
// @Security ApiKeyAuth
// @Tags match
// @Accept json
// @Produce json
// @Param id path string true "Match ID"
// @Param choiceid path string true "Pending choice ID"
// @Param input body api.HandDiscardRequest true "hand-discard input; requestid and expectedrevision are required"
// @Success 200 {object} api.MatchResponse
// @Failure 400,401,403,404,409,500 {object} api.ErrorResponse
// @Router /match/{id}/choices/{choiceid}/hand-discard [post]
func (c *Controller) HandDiscard(ctx *gin.Context) {
	submitTyped(c, ctx, func(r api.HandDiscardRequest) game.Command {
		command := game.Command{RequestID: r.RequestID, ExpectedRevision: r.ExpectedRevision}
		command.Type = "resolve_choice"
		command.CardID = r.CardID
		command.CardID = ""
		command.Selections = []string{r.CardID}
		command.Mode = r.Disposition
		return command
	}, "choices/hand-discard")
}

// AssetLoss submits asset-loss.
// @Summary asset-loss
// @Security ApiKeyAuth
// @Tags match
// @Accept json
// @Produce json
// @Param id path string true "Match ID"
// @Param choiceid path string true "Pending choice ID"
// @Param input body api.AssetLossRequest true "asset-loss input; requestid and expectedrevision are required"
// @Success 200 {object} api.MatchResponse
// @Failure 400,401,403,404,409,500 {object} api.ErrorResponse
// @Router /match/{id}/choices/{choiceid}/asset-loss [post]
func (c *Controller) AssetLoss(ctx *gin.Context) {
	submitTyped(c, ctx, func(r api.AssetLossRequest) game.Command {
		command := game.Command{RequestID: r.RequestID, ExpectedRevision: r.ExpectedRevision}
		command.Type = "resolve_choice"
		command.CardID = r.CardID
		command.CardID = ""
		command.Selections = []string{r.CardID}
		return command
	}, "choices/asset-loss")
}

// Scry submits scry.
// @Summary scry
// @Security ApiKeyAuth
// @Tags match
// @Accept json
// @Produce json
// @Param id path string true "Match ID"
// @Param choiceid path string true "Pending choice ID"
// @Param input body api.ScryRequest true "scry input; requestid and expectedrevision are required"
// @Success 200 {object} api.MatchResponse
// @Failure 400,401,403,404,409,500 {object} api.ErrorResponse
// @Router /match/{id}/choices/{choiceid}/scry [post]
func (c *Controller) Scry(ctx *gin.Context) {
	submitTyped(c, ctx, func(r api.ScryRequest) game.Command {
		command := game.Command{RequestID: r.RequestID, ExpectedRevision: r.ExpectedRevision}
		command.Type = "resolve_choice"
		command.Order = r.Order
		command.Selections = []string{r.Take}
		return command
	}, "choices/scry")
}

// React submits play-card.
// @Summary play-card
// @Security ApiKeyAuth
// @Tags match
// @Accept json
// @Produce json
// @Param id path string true "Match ID"
// @Param input body api.ReactRequest true "play-card input; requestid and expectedrevision are required"
// @Success 200 {object} api.MatchResponse
// @Failure 400,401,403,404,409,500 {object} api.ErrorResponse
// @Router /match/{id}/reactions/play-card [post]
func (c *Controller) React(ctx *gin.Context) {
	submitTyped(c, ctx, func(r api.ReactRequest) game.Command {
		command := game.Command{RequestID: r.RequestID, ExpectedRevision: r.ExpectedRevision}
		command.Type = "play_card"
		command.CardID = r.CardID
		command.Targets = r.Targets
		command.Support = r.Support
		command.Payment = r.Payment
		return command
	}, "reactions/play-card")
}

// Defend submits defend.
// @Summary defend
// @Security ApiKeyAuth
// @Tags match
// @Accept json
// @Produce json
// @Param id path string true "Match ID"
// @Param input body api.DefendRequest true "defend input; requestid and expectedrevision are required"
// @Success 200 {object} api.MatchResponse
// @Failure 400,401,403,404,409,500 {object} api.ErrorResponse
// @Router /match/{id}/reactions/defend [post]
func (c *Controller) Defend(ctx *gin.Context) {
	submitTyped(c, ctx, func(r api.DefendRequest) game.Command {
		command := game.Command{RequestID: r.RequestID, ExpectedRevision: r.ExpectedRevision}
		command.Type = "defend"
		command.Minions = r.Minions
		command.Support = r.Support
		return command
	}, "reactions/defend")
}

// Pass submits pass.
// @Summary pass
// @Security ApiKeyAuth
// @Tags match
// @Accept json
// @Produce json
// @Param id path string true "Match ID"
// @Param input body api.PassRequest true "pass input; requestid and expectedrevision are required"
// @Success 200 {object} api.MatchResponse
// @Failure 400,401,403,404,409,500 {object} api.ErrorResponse
// @Router /match/{id}/reactions/pass [post]
func (c *Controller) Pass(ctx *gin.Context) {
	submitTyped(c, ctx, func(r api.PassRequest) game.Command {
		command := game.Command{RequestID: r.RequestID, ExpectedRevision: r.ExpectedRevision}
		command.Type = "pass"
		return command
	}, "reactions/pass")
}
