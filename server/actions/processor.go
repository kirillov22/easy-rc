package actions

func Process(p Processable, robot Robot) (R any, err error) {
	debugLog("[DEBUG] Processing action: %s", p)
	result, err := p.Process(robot)
	if err != nil {
		debugLog("[DEBUG] Action %s failed: %v", p, err)
	}
	return result, err
}
