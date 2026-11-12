package bubbleup

import (
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
)

// spawnActiveAlert feeds a real alertMsg through Update so tests
// exercise the actual spawn path.
func spawnActiveAlert(t *testing.T, m AlertModel) (AlertModel, tea.Cmd) {
	t.Helper()
	out, cmd := m.Update(alertMsg{alertKey: InfoKey, msg: "hello", dur: time.Minute})
	spawned := out.(AlertModel)
	if !spawned.HasActiveAlert() {
		t.Fatal("alert should be active after alertMsg")
	}
	if cmd == nil {
		t.Fatal("spawning an alert should arm its tick chain")
	}
	return spawned, cmd
}

// TestUnrelatedMessagesDoNotArmTicks pins the no-amplification
// contract: while an alert is active, messages unrelated to the
// alert lifecycle must not spawn additional tick chains.
func TestUnrelatedMessagesDoNotArmTicks(t *testing.T) {
	m, _ := spawnActiveAlert(t, *NewAlertModel(50, false, time.Minute))

	for _, msg := range []tea.Msg{
		tea.MouseMotionMsg{},
		tea.WindowSizeMsg{Width: 80, Height: 24},
		struct{}{},
	} {
		out, cmd := m.Update(msg)
		if cmd != nil {
			t.Fatalf("%T armed a duplicate tick chain", msg)
		}
		if !out.(AlertModel).HasActiveAlert() {
			t.Fatalf("%T cleared the active alert", msg)
		}
	}
}

// TestTickSustainsWhileActive guards against over-fixing: the
// alert's own tick must keep the single chain running while alive.
func TestTickSustainsWhileActive(t *testing.T) {
	m, _ := spawnActiveAlert(t, *NewAlertModel(50, false, time.Minute))

	out, cmd := m.Update(tickMsg(time.Now()))
	if cmd == nil {
		t.Fatal("an active alert should re-arm its tick chain")
	}
	if !out.(AlertModel).HasActiveAlert() {
		t.Fatal("alert should still be active after its own tick")
	}
}

// TestExpiredTickStopsChain verifies the lifecycle end: once past
// deathTime, the tick clears the alert and stops ticking.
func TestExpiredTickStopsChain(t *testing.T) {
	m, _ := spawnActiveAlert(t, *NewAlertModel(50, false, time.Minute))

	out, cmd := m.Update(tickMsg(time.Now().Add(2 * time.Minute)))
	if out.(AlertModel).HasActiveAlert() {
		t.Fatal("alert should have expired")
	}
	if cmd != nil {
		t.Fatal("an expired alert should not re-arm the tick chain")
	}
}
