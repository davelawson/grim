package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	randv2 "math/rand/v2"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/davelawson/grim/api"
	"github.com/davelawson/grim/game"
	"github.com/google/uuid"
)

// The CLI's module imports only api. These opt-in tests host production routes
// on the existing isolated SQLite fixture and execute a freshly built external
// binary through complete games. No production account or server is needed.
func TestCLICompleteGames(t *testing.T) {
	binary := os.Getenv("GRIMCLI_TEST_BINARY")
	if binary == "" {
		t.Skip("set GRIMCLI_TEST_BINARY to a freshly built grimcli binary")
	}
	binary, err := filepath.Abs(binary)
	if err != nil {
		t.Fatal(err)
	}
	if info, err := os.Stat(binary); err != nil || info.IsDir() {
		t.Fatalf("GRIMCLI_TEST_BINARY is not a file: %s (%v)", binary, err)
	}
	for _, size := range []int{2, 3, 4} {
		t.Run(fmt.Sprintf("%d", size), func(t *testing.T) {
			f := newLifecycleFixture(t)
			server := httptest.NewServer(f.router)
			defer server.Close()
			dir := t.TempDir()
			invoke := func(actor, command string, flags ...string) []byte {
				t.Helper()
				ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
				defer cancel()
				args := []string{command, "--server", server.URL, "--token", "token-" + actor}
				args = append(args, flags...)
				cmd := exec.CommandContext(ctx, binary, args...)
				cmd.Dir = dir
				var out, diagnostics bytes.Buffer
				cmd.Stdout = &out
				cmd.Stderr = &diagnostics
				if err := cmd.Run(); err != nil {
					t.Fatalf("CLI %q: %v\n%s", args, err, diagnostics.String())
				}
				return out.Bytes()
			}
			current := func(id, actor string) *game.View {
				t.Helper()
				var result api.MatchResponse
				if err := json.Unmarshal(invoke(actor, "match", "--id", id), &result); err != nil || result.Match == nil {
					t.Fatalf("decode match: %v", err)
				}
				if result.Match.You == nil || result.Match.You.UserID != f.ids[actor] {
					t.Fatal("caller identity missing")
				}
				return result.Match
			}
			var lobby api.CreateLobbyResponse
			if err := json.Unmarshal(invoke("owner", "lobby-create", "--name", "CLI complete game"), &lobby); err != nil {
				t.Fatal(err)
			}
			for _, actor := range roster[1:size] {
				invoke(actor, "lobby-join", "--id", lobby.Id)
			}
			for _, actor := range roster[:size] {
				invoke(actor, "lobby-ready", "--id", lobby.Id, "--ready=true")
			}
			var launched api.MatchResponse
			if err := json.Unmarshal(invoke("owner", "lobby-launch", "--id", lobby.Id), &launched); err != nil || launched.Match == nil {
				t.Fatalf("launch: %v", err)
			}
			id := launched.Match.ID
			// Reuse the reproducible initial deal used by the backend HTTP lifecycle
			// tests. Every choice and gameplay transition below goes through the CLI.
			f.pinCLISeed(id)
			var catalogue api.CatalogueResponse
			if err := json.Unmarshal(invoke("owner", "catalogue", "--id", id), &catalogue); err != nil || catalogue.Catalogue == nil {
				t.Fatalf("catalogue: %v", err)
			}
			for _, actor := range roster[:size] {
				for _, pick := range []struct{ typ, card string }{{"choose_chantry", "chantry_fire"}, {"choose_victory_path", "subjugate_rival"}, {"choose_minion", "bastian_redhand"}} {
					var draft api.DraftResponse
					if err := json.Unmarshal(invoke(actor, "draft", "--id", id), &draft); err != nil {
						t.Fatal(err)
					}
					invoke(actor, "draft-pick", "--id", id, "--type", pick.typ, "--cardid", pick.card, "--revision", fmt.Sprint(draft.Revision), "--requestid", uuid.NewString())
				}
			}
			count := 0
			submit := func(id, actor string, c game.Command) *game.View {
				t.Helper()
				view := current(id, actor)
				path, input := typedGameplay(t, view, c)
				parts := strings.Split(path, "/")
				command := parts[len(parts)-1]
				if command == "attempt-challenge" {
					command = "challenge"
				}
				if command == "activate-ability" {
					command = "activate"
				}
				if command == "research" {
					command = "choose-research"
				}
				if strings.Contains(path, "/reactions/play-card") {
					command = "react"
				}
				args := []string{"--id", id, "--revision", fmt.Sprint(view.Revision), "--requestid", uuid.NewString()}
				if c.ChoiceID != "" {
					args = append(args, "--choiceid", c.ChoiceID)
				}
				var fields map[string]any
				if err := json.Unmarshal(mustJSON(t, input), &fields); err != nil {
					t.Fatal(err)
				}
				for key, value := range fields {
					if key == "requestid" || key == "expectedrevision" {
						continue
					}
					flag := map[string]string{"targets": "target", "minions": "minion", "elements": "element-card", "cards": "discard"}[key]
					if flag == "" {
						flag = key
					}
					if key == "support" {
						for _, entry := range value.([]any) {
							support := entry.(map[string]any)
							card := support["cardid"].(string)
							args = append(args, "--support", card)
							if domains, ok := support["domains"].([]any); ok {
								for _, domain := range domains {
									args = append(args, "--support-domain", card+"="+domain.(string))
								}
							}
						}
					} else if values, ok := value.([]any); ok {
						for _, v := range values {
							args = append(args, "--"+flag, v.(string))
						}
					} else if value != nil {
						args = append(args, "--"+flag, fmt.Sprint(value))
					}
				}
				var response api.MatchResponse
				if err := json.Unmarshal(invoke(actor, command, args...), &response); err != nil || response.Match == nil {
					t.Fatalf("command response: %v", err)
				}
				authoritative := current(id, actor)
				if !reflect.DeepEqual(response.Match, authoritative) {
					t.Fatal("CLI command response differs from authoritative view")
				}
				count++
				return authoritative
			}
			finished := f.playToVictoryUsing(current(id, "owner"), size > 2, submit, false)
			if finished.Status != "finished" || len(finished.Winners) != 1 || count == 0 {
				t.Fatal("CLI did not reach victory")
			}
			// The terminal result remains visible to every participant.
			for _, actor := range roster[:size] {
				view := current(id, actor)
				if !reflect.DeepEqual(view.Winners, finished.Winners) {
					t.Fatal("winner not visible")
				}
			}
			t.Logf("completed %d-player game through %d CLI gameplay submissions", size, count)
		})
	}
}

func (f *lifecycleFixture) pinCLISeed(id string) {
	f.t.Helper()
	_, s := f.saved(id)
	seed := [32]byte{5, 8, 13}
	g := randv2.NewChaCha8(seed)
	state, err := g.MarshalBinary()
	if err != nil {
		f.t.Fatal(err)
	}
	s.Random.Seed = seed[:]
	s.Random.State = state
	data, err := json.Marshal(s)
	if err != nil {
		f.t.Fatal(err)
	}
	f.exec("update matches set snapshot = ? where id = ?", string(data), id)
}
