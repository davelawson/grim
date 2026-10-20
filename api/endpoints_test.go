package api_test

import (
	"github.com/davelawson/grim/api"
	"net/url"
	"strings"
	"testing"
)

func TestEndpointPaths(t *testing.T) {
	tests := []struct {
		name string
		path func(string) string
		want string
	}{
		{"leave lobby", api.LeaveLobbyPath, "/lobby/ID/leave"},
		{"join lobby", api.JoinLobbyPath, "/lobby/ID/join"},
		{"delete lobby", api.DeleteLobbyPath, "/lobby/ID"},
		{"get lobby", api.GetLobbyPath, "/lobby/ID"},
		{"update lobby", api.UpdateLobbyPath, "/lobby/ID"},
		{"add user", api.AddUserToLobbyPath, "/lobby/ID/user"},
		{"remove user lobby ID", func(id string) string { return api.RemoveUserFromLobbyPath(id, "user") }, "/lobby/ID/user/user"},
		{"remove user user ID", func(id string) string { return api.RemoveUserFromLobbyPath("lobby", id) }, "/lobby/lobby/user/ID"},
		{"ready", api.SetReadyPath, "/lobby/ID/ready"},
		{"launch", api.LaunchMatchPath, "/lobby/ID/launch"},
		{"PlayCard", api.PlayCardPath, "/match/ID/actions/play-card"},
		{"Attack", api.AttackPath, "/match/ID/actions/attack"},
		{"AttemptChallenge", api.AttemptChallengePath, "/match/ID/actions/attempt-challenge"},
		{"ActivateAbility", api.ActivateAbilityPath, "/match/ID/actions/activate-ability"},
		{"ResearchElement", api.ResearchElementPath, "/match/ID/actions/research-element"},
		{"ResearchSpell", api.ResearchSpellPath, "/match/ID/actions/research-spell"},
		{"HealWithFocus", api.HealWithFocusPath, "/match/ID/actions/heal-with-focus"},
		{"EndTurn", api.EndTurnPath, "/match/ID/actions/end-turn"},
		{"Maintenance", func(id string) string { return api.MaintenancePath(id, "choice") }, "/match/ID/choices/choice/maintenance"},
		{"Kneel", func(id string) string { return api.KneelPath(id, "choice") }, "/match/ID/choices/choice/kneel"},
		{"ResolveEffect", func(id string) string { return api.ResolveEffectPath(id, "choice") }, "/match/ID/choices/choice/resolve-effect"},
		{"SkipEffect", func(id string) string { return api.SkipEffectPath(id, "choice") }, "/match/ID/choices/choice/skip-effect"},
		{"UndertakeHunt", func(id string) string { return api.UndertakeHuntPath(id, "choice") }, "/match/ID/choices/choice/undertake-hunt"},
		{"ChooseResearch", func(id string) string { return api.ChooseResearchPath(id, "choice") }, "/match/ID/choices/choice/research"},
		{"EndDiscard", func(id string) string { return api.EndDiscardPath(id, "choice") }, "/match/ID/choices/choice/end-discard"},
		{"HandDiscard", func(id string) string { return api.HandDiscardPath(id, "choice") }, "/match/ID/choices/choice/hand-discard"},
		{"AssetLoss", func(id string) string { return api.AssetLossPath(id, "choice") }, "/match/ID/choices/choice/asset-loss"},
		{"Scry", func(id string) string { return api.ScryPath(id, "choice") }, "/match/ID/choices/choice/scry"},
		{"React", api.ReactPath, "/match/ID/reactions/play-card"},
		{"Defend", api.DefendPath, "/match/ID/reactions/defend"},
		{"Pass", api.PassPath, "/match/ID/reactions/pass"},
		{"get match", api.GetMatchPath, "/match/ID"},
		{"get draft", api.GetDraftPath, "/match/ID/draft"},
		{"submit draft", api.SubmitDraftPath, "/match/ID/draft"},
		{"catalogue", api.GetCataloguePath, "/match/ID/catalogue"},
		{"end match", api.EndMatchPath, "/admin/match/ID/end"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parameters resembling route placeholders must remain literal values.
			for _, id := range []string{"abc-123", "a/b c?d#e%f", "雪", ":user_id", ":id"} {
				want := strings.Replace(tt.want, "ID", url.PathEscape(id), 1)
				got := tt.path(id)
				if got != want {
					t.Errorf("path(%q) = %q, want %q", id, got, want)
				}
				parsed, err := url.Parse("https://localhost:8080" + got)
				if err != nil {
					t.Fatal(err)
				}
				if parsed.RawQuery != "" || parsed.Fragment != "" || parsed.EscapedPath() != got {
					t.Errorf("parameter changed URL structure: %s", parsed)
				}
			}
		})
	}
}

func TestRemoveUserPathEscapesBothParameters(t *testing.T) {
	got := api.RemoveUserFromLobbyPath("lobby/one", "user/two?#")
	if want := "/lobby/lobby%2Fone/user/user%2Ftwo%3F%23"; got != want {
		t.Fatalf("path = %q, want %q", got, want)
	}
}

func TestChoicePathsEscapeBothParameters(t *testing.T) {
	if got := api.MaintenancePath(":choiceid/a", ":id/b ?#%"); got != "/match/"+url.PathEscape(":choiceid/a")+"/choices/"+url.PathEscape(":id/b ?#%")+"/maintenance" {
		t.Fatal(got)
	}
	if got := api.KneelPath(":choiceid/a", ":id/b ?#%"); got != "/match/"+url.PathEscape(":choiceid/a")+"/choices/"+url.PathEscape(":id/b ?#%")+"/kneel" {
		t.Fatal(got)
	}
	if got := api.ResolveEffectPath(":choiceid/a", ":id/b ?#%"); got != "/match/"+url.PathEscape(":choiceid/a")+"/choices/"+url.PathEscape(":id/b ?#%")+"/resolve-effect" {
		t.Fatal(got)
	}
	if got := api.SkipEffectPath(":choiceid/a", ":id/b ?#%"); got != "/match/"+url.PathEscape(":choiceid/a")+"/choices/"+url.PathEscape(":id/b ?#%")+"/skip-effect" {
		t.Fatal(got)
	}
	if got := api.UndertakeHuntPath(":choiceid/a", ":id/b ?#%"); got != "/match/"+url.PathEscape(":choiceid/a")+"/choices/"+url.PathEscape(":id/b ?#%")+"/undertake-hunt" {
		t.Fatal(got)
	}
	if got := api.ChooseResearchPath(":choiceid/a", ":id/b ?#%"); got != "/match/"+url.PathEscape(":choiceid/a")+"/choices/"+url.PathEscape(":id/b ?#%")+"/research" {
		t.Fatal(got)
	}
	if got := api.EndDiscardPath(":choiceid/a", ":id/b ?#%"); got != "/match/"+url.PathEscape(":choiceid/a")+"/choices/"+url.PathEscape(":id/b ?#%")+"/end-discard" {
		t.Fatal(got)
	}
	if got := api.HandDiscardPath(":choiceid/a", ":id/b ?#%"); got != "/match/"+url.PathEscape(":choiceid/a")+"/choices/"+url.PathEscape(":id/b ?#%")+"/hand-discard" {
		t.Fatal(got)
	}
	if got := api.AssetLossPath(":choiceid/a", ":id/b ?#%"); got != "/match/"+url.PathEscape(":choiceid/a")+"/choices/"+url.PathEscape(":id/b ?#%")+"/asset-loss" {
		t.Fatal(got)
	}
	if got := api.ScryPath(":choiceid/a", ":id/b ?#%"); got != "/match/"+url.PathEscape(":choiceid/a")+"/choices/"+url.PathEscape(":id/b ?#%")+"/scry" {
		t.Fatal(got)
	}
}
