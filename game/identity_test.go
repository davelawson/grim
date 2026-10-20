package game

import (
	"encoding/json"
	"testing"

	"github.com/davelawson/grim/api"
)

func TestPrivateIdentityIndependentOfCardsAndTiming(t *testing.T) {
	s := table(t, 3)
	for _, actor := range s.Participants {
		w := s.Wizards[actor]
		w.Hand = nil
		w.DrawPile = nil
		w.Eliminated = true
		for _, status := range []string{"playing", "finished"} {
			s.Status = status
			v := s.View(actor)
			if v.You == nil || v.You.UserID != actor || len(v.You.Hand) != 0 {
				t.Fatalf("%s/%s: %+v", actor, status, v.You)
			}
			// Backward-compatible receipts without identity still decode.
			data := []byte(`{"draft":{"step":"complete"},"hand":[],"drawpile":[]}`)
			var old api.PrivateView
			if err := json.Unmarshal(data, &old); err != nil || old.UserID != "" {
				t.Fatalf("old receipt: %+v %v", old, err)
			}
		}
	}
	if s.View("outsider").You != nil {
		t.Fatal("outsider received private identity")
	}
	// Identity is a projection addition, never authoritative saved state.
	var snapshot map[string]any
	if err := json.Unmarshal(jsonBytes(t, s), &snapshot); err != nil {
		t.Fatal(err)
	}
	if _, ok := snapshot["you"]; ok {
		t.Fatal("private projection entered snapshot")
	}
}
