package hook

import "sort"

type Hook struct {
	pre  []NamedHook
	post []NamedHook
}

type NamedHook struct {
	Name string
	Fn   func()
}

func (h *Hook) AddPre(name string, fn func()) {
	h.pre = append(h.pre, NamedHook{name, fn})
}

func (h *Hook) AddPost(name string, fn func()) {
	h.post = append(h.post, NamedHook{name, fn})
}

func (h *Hook) RunPre() {
	sort.SliceStable(h.pre, func(i, j int) bool {
		return h.pre[i].Name < h.pre[j].Name
	})
	for _, hook := range h.pre {
		hook.Fn()
	}
}

func (h *Hook) RunPost() {
	sort.SliceStable(h.post, func(i, j int) bool {
		return h.post[i].Name < h.post[j].Name
	})
	for _, hook := range h.post {
		hook.Fn()
	}
}
