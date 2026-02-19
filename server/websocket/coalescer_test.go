package websocket

import (
	"easy-rc-server/actions"
	"easy-rc-server/actions/model"
	"testing"
)

func TestCoalesceActionsEmpty(t *testing.T) {
	result := CoalesceActions(nil)
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d items", len(result))
	}
}

func TestCoalesceActionsSingleMove(t *testing.T) {
	batch := []actions.Processable{model.NewMove(10, 20)}
	result := CoalesceActions(batch)

	if len(result) != 1 {
		t.Fatalf("expected 1 action, got %d", len(result))
	}
	m, ok := result[0].(model.Move)
	if !ok {
		t.Fatalf("expected Move, got %T", result[0])
	}
	if m.MoveX() != 10 || m.MoveY() != 20 {
		t.Errorf("expected (10, 20), got (%d, %d)", m.MoveX(), m.MoveY())
	}
}

func TestCoalesceActionsConsecutiveMoves(t *testing.T) {
	batch := []actions.Processable{
		model.NewMove(5, 10),
		model.NewMove(3, -2),
		model.NewMove(7, 8),
	}
	result := CoalesceActions(batch)

	if len(result) != 1 {
		t.Fatalf("expected 1 action, got %d", len(result))
	}
	m := result[0].(model.Move)
	if m.MoveX() != 15 || m.MoveY() != 16 {
		t.Errorf("expected (15, 16), got (%d, %d)", m.MoveX(), m.MoveY())
	}
}

func TestCoalesceActionsMoveClickMove(t *testing.T) {
	click := model.NewClick(actions.LeftButton)
	batch := []actions.Processable{
		model.NewMove(5, 10),
		model.NewMove(3, 2),
		click,
		model.NewMove(1, 1),
		model.NewMove(2, 3),
	}
	result := CoalesceActions(batch)

	if len(result) != 3 {
		t.Fatalf("expected 3 actions, got %d", len(result))
	}

	m1 := result[0].(model.Move)
	if m1.MoveX() != 8 || m1.MoveY() != 12 {
		t.Errorf("first move: expected (8, 12), got (%d, %d)", m1.MoveX(), m1.MoveY())
	}

	if _, ok := result[1].(model.Click); !ok {
		t.Errorf("expected Click at index 1, got %T", result[1])
	}

	m2 := result[2].(model.Move)
	if m2.MoveX() != 3 || m2.MoveY() != 4 {
		t.Errorf("second move: expected (3, 4), got (%d, %d)", m2.MoveX(), m2.MoveY())
	}
}

func TestCoalesceActionsNonMoveOnly(t *testing.T) {
	ping := model.NewPing()
	batch := []actions.Processable{ping}
	result := CoalesceActions(batch)

	if len(result) != 1 {
		t.Fatalf("expected 1 action, got %d", len(result))
	}
	if _, ok := result[0].(model.Ping); !ok {
		t.Errorf("expected Ping, got %T", result[0])
	}
}

func TestCoalesceActionsMixedSequence(t *testing.T) {
	batch := []actions.Processable{
		model.NewMove(1, 2),
		model.NewMove(3, 4),
		model.NewClick(actions.LeftButton),
		model.NewMove(5, 6),
		model.NewMove(7, 8),
		model.NewPing(),
	}
	result := CoalesceActions(batch)

	if len(result) != 4 {
		t.Fatalf("expected 4 actions, got %d", len(result))
	}

	m1 := result[0].(model.Move)
	if m1.MoveX() != 4 || m1.MoveY() != 6 {
		t.Errorf("first move: expected (4, 6), got (%d, %d)", m1.MoveX(), m1.MoveY())
	}

	if _, ok := result[1].(model.Click); !ok {
		t.Errorf("expected Click at index 1, got %T", result[1])
	}

	m2 := result[2].(model.Move)
	if m2.MoveX() != 12 || m2.MoveY() != 14 {
		t.Errorf("second move: expected (12, 14), got (%d, %d)", m2.MoveX(), m2.MoveY())
	}

	if _, ok := result[3].(model.Ping); !ok {
		t.Errorf("expected Ping at index 3, got %T", result[3])
	}
}
