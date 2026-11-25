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

// TestReplacingAlertDoesNotArmSecondChain: a new alert arriving while
// one is active must reuse the in-flight chain, not start another.
func TestReplacingAlertDoesNotArmSecondChain(t *testing.T) {
	m, _ := spawnActiveAlert(t, *NewAlertModel(50, false, time.Minute))

	out, cmd := m.Update(alertMsg{alertKey: WarnKey, msg: "second", dur: time.Minute})
	if cmd != nil {
		t.Fatal("replacing an active alert armed a second tick chain")
	}
	if !out.(AlertModel).HasActiveAlert() {
		t.Fatal("replacement alert should be active")
	}
}

// TestEscThenNewAlertReusesInFlightTick: after Esc, the dismissed
// alert's tick is still in flight, so an alert arriving before it
// lands must not arm a second chain.
func TestEscThenNewAlertReusesInFlightTick(t *testing.T) {
	m := *NewAlertModel(50, false, time.Minute)
	m.allowEscToClose = true
	m, _ = spawnActiveAlert(t, m)

	out, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyEscape})
	m = out.(AlertModel)
	if m.HasActiveAlert() {
		t.Fatal("esc should have closed the alert")
	}

	out, cmd := m.Update(alertMsg{alertKey: InfoKey, msg: "again", dur: time.Minute})
	if cmd != nil {
		t.Fatal("new alert armed a second chain while the old tick was in flight")
	}

	// The in-flight tick lands and adopts the new alert.
	if _, cmd = out.(AlertModel).Update(tickMsg(time.Now())); cmd == nil {
		t.Fatal("in-flight tick should keep ticking for the new alert")
	}
}

// TestChainRestartsAfterEnding: once a chain has ended, the next
// alert must arm a fresh one.
func TestChainRestartsAfterEnding(t *testing.T) {
	m, _ := spawnActiveAlert(t, *NewAlertModel(50, false, time.Minute))

	out, _ := m.Update(tickMsg(time.Now().Add(2 * time.Minute)))
	if _, cmd := out.(AlertModel).Update(alertMsg{alertKey: InfoKey, msg: "next", dur: time.Minute}); cmd == nil {
		t.Fatal("alert after chain ended should arm a new tick chain")
	}
}
