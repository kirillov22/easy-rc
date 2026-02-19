package websocket

import (
	"easy-rc-server/actions"
	"easy-rc-server/actions/model"
)

func CoalesceActions(batch []actions.Processable) []actions.Processable {
	if len(batch) == 0 {
		return nil
	}

	var result []actions.Processable
	var accX, accY int32
	var hasPendingMove bool

	for _, action := range batch {
		if m, ok := action.(model.Move); ok {
			accX += m.MoveX()
			accY += m.MoveY()
			hasPendingMove = true
			continue
		}

		if hasPendingMove {
			result = append(result, model.NewMove(accX, accY))
			accX, accY = 0, 0
			hasPendingMove = false
		}
		result = append(result, action)
	}

	if hasPendingMove {
		result = append(result, model.NewMove(accX, accY))
	}

	return result
}

func drainChannel(first actions.Processable, ch <-chan actions.Processable) []actions.Processable {
	batch := []actions.Processable{first}
	for {
		select {
		case item, ok := <-ch:
			if !ok {
				return batch
			}
			batch = append(batch, item)
		default:
			return batch
		}
	}
}
