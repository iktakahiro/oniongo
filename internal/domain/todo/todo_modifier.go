package todo

import "time"

type TodoModifier struct {
	todo *Todo
}

func (t *Todo) Modify() *TodoModifier {
	return &TodoModifier{todo: t}
}

func (m *TodoModifier) SetTitle(title string) *TodoModifier {
	m.todo.title = title
	m.todo.updatedAt = time.Now()
	return m
}

func (m *TodoModifier) SetBody(body string) *TodoModifier {
	m.todo.body = body
	m.todo.updatedAt = time.Now()
	return m
}
