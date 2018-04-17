package jobs

type QueueDef struct {
	Name string
}

type QueueSystem struct {
	Queues map[string]*PgJobQueue
	Defs   map[string]*QueueDef
}

func OpenQueueSystem(url string, defs []*QueueDef) (qs *QueueSystem, err error) {
	queues := make(map[string]*PgJobQueue)
	defsMap := make(map[string]*QueueDef)

	for _, def := range defs {
		jq, err := NewPqJobQueue(url, def.Name)
		if err != nil {
			return nil, err
		}
		queues[def.Name] = jq
		defsMap[def.Name] = def
	}

	qs = &QueueSystem{
		Queues: queues,
		Defs:   defsMap,
	}

	return
}
