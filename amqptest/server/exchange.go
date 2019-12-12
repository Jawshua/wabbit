package server

import "fmt"

type Exchange interface {
	route(route string, d *Delivery) error
	addBinding(route string, q *Queue)
	delBinding(route string)
}

type TopicExchange struct {
	name     string
	bindings map[string]*Queue
}

func NewTopicExchange(name string) *TopicExchange {
	return &TopicExchange{
		name:     name,
		bindings: make(map[string]*Queue),
	}
}

func (t *TopicExchange) addBinding(route string, q *Queue) {
	t.bindings[route] = q
}

func (t *TopicExchange) delBinding(route string) {
	delete(t.bindings, route)
}

func (t *TopicExchange) route(route string, d *Delivery) error {
	for bname, q := range t.bindings {
		if topicMatch(bname, route) {
			q.data <- d
			return nil
		}
	}

	// The route doesnt match any binding, then will be discarded
	return nil
}

type DirectExchange struct {
	name     string
	bindings map[string]*Queue
}

func NewDirectExchange(name string) *DirectExchange {
	return &DirectExchange{
		name:     name,
		bindings: make(map[string]*Queue),
	}
}

func (d *DirectExchange) addBinding(route string, q *Queue) {
	if d.bindings == nil {
		d.bindings = make(map[string]*Queue)
	}

	d.bindings[route] = q
}

func (d *DirectExchange) delBinding(route string) {
	delete(d.bindings, route)
}

func (d *DirectExchange) route(route string, delivery *Delivery) error {
	if q, ok := d.bindings[route]; ok {
		q.data <- delivery
		return nil
	}

	return fmt.Errorf("No bindings to route: %s", route)

}

// FanoutExchange routes messages to all of the queues that are
// bound to it and the routing key is ignored.
type FanoutExchange struct {
	name     string
	bindings map[string]*Queue
}

func NewFanoutExchange(name string) *FanoutExchange {
	return &FanoutExchange{
		name:     name,
		bindings: make(map[string]*Queue),
	}
}

func (t *FanoutExchange) addBinding(route string, q *Queue) {
	t.bindings[q.name] = q
}

func (t *FanoutExchange) delBinding(route string) {
	delete(t.bindings, route)
}

func (t *FanoutExchange) route(route string, d *Delivery) error {
	for _, q := range t.bindings {
		q.data <- d
	}

	return nil
}
