package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Cache struct {
	cache map[Page]Teable
}

func NewCache() *Cache {
	return &Cache{
		cache: make(map[Page]Teable),
	}
}

func (c *Cache) Get(page Page) Teable {
	return c.cache[page]
}

func (c *Cache) Put(page Page, model Teable) {
	c.cache[page] = model
}

func (c *Cache) UpdateAll(msg tea.Msg) []tea.Cmd {
	cmds := make([]tea.Cmd, len(c.cache))
	var i int
	for k := range c.cache {
		cmds[i] = c.Update(k, msg)
		i++
	}
	return cmds
}

func (c *Cache) Update(key Page, msg tea.Msg) tea.Cmd {
	return c.cache[key].Update(msg)
}
