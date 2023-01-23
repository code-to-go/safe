package main

import (
	"fmt"
	"math"
	"regexp"
	"strings"

	"github.com/code-to-go/safepool.lib/api/chat"
	"github.com/code-to-go/safepool.lib/core"
	"github.com/code-to-go/safepool.lib/pool"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
)

func printChatHelp() {
	color.White("commands: ")
	color.White("  '' refresh chat content")
	color.White("  \\x exit chat")
	color.White("  \\c create a sub pool")
}

var isValidName = regexp.MustCompile(`^[a-zA-Z0-9#]+$`).MatchString

func createChat(c chat.Chat) {

	var name string
	for {
		prompt := promptui.Prompt{
			Label:       "Pool name (only alphanumeric and #): ",
			HideEntered: true,
		}

		name, _ = prompt.Run()
		if name == "" {
			return
		}
		if isValidName(name) {
			break
		}
		color.Red("Invalid name '%s'. Name can contain only alphanumeric letters and #.", name)
	}

	selfId := c.Pool.Self.Id()
	selected := map[string]bool{}
	for {
		items := []string{"Complete"}
		identities, _ := c.Pool.Identities()
		for idx, i := range identities {
			id := i.Id()
			if id == selfId {
				if idx < len(identities)-1 {
					identities[idx] = identities[len(identities)-1]
					identities[len(identities)-1] = i
				} else {
					continue
				}
			}
			if selected[id] {
				items = append(items, fmt.Sprintf("✓ %s [%s]", i.Nick, id))
			} else {
				items = append(items, fmt.Sprintf("  %s [%s]", i.Nick, id))
			}
		}

		sel := promptui.Select{
			Label: "Select users for the new pool",
			Items: items,
		}
		idx, _, err := sel.Run()
		if err != nil {
			return
		}

		if idx == 0 {
			break
		}
		id := identities[idx-1].Id()
		selected[id] = !selected[id]
	}

	var ids []string
	for id, ok := range selected {
		if ok {
			ids = append(ids, id)
		}
	}

	co, err := c.Pool.CreateBranch(name, ids)
	if core.IsErr(err, "cannot create branch in pool %v: %v", c.Pool) {
		color.Red("😱 something went wrong!")
	}

	token := pool.Token{
		Config: co,
		Host:   c.Pool.Self,
	}
	for _, id := range ids {
		tk, err := pool.EncodeToken(token, id)
		if err == nil {
			c.SendMessage(fmt.Sprintf("%s,%s", id, tk), "application/token", nil)
		}
	}
}

func Chat(p *pool.Pool) {
	var lastId uint64
	c := chat.Get(p, "chat")

	identities, err := p.Identities()
	if err != nil {
		color.Red("cannot retrieve identities for pool '%s': %v", p.Name)
		return
	}

	id2nick := map[string]string{}
	for _, i := range identities {
		id2nick[i.Id()] = i.Nick
	}

	selfId := p.Self.Id()
	color.Green("Enter \\? for list of commands")
	for {
		messages, err := c.GetMessages(lastId, math.MaxInt64, 32)
		if err != nil {
			color.Red("cannot retrieve chat messages from pool '%s': %v", p.Name)
			return
		}
		for _, m := range messages {
			if m.Author == selfId {
				color.Blue("%s: %s", id2nick[m.Author], m.Content)
			} else {
				color.Green("%s: %s", id2nick[m.Author], m.Content)
			}
			if m.Id > lastId {
				lastId = m.Id
			}
		}
		prompt := promptui.Prompt{
			Label:       "> ",
			HideEntered: true,
		}

		t, _ := prompt.Run()
		t = strings.Trim(t, " ")

		switch {
		case len(t) == 0:
		case strings.HasPrefix(t, "\\x"):
			return
		case strings.HasPrefix(t, "\\c"):
			createChat(c)
		case strings.HasPrefix(t, "\\?"):
			printChatHelp()
		case strings.HasPrefix(t, "\\"):
			printChatHelp()
		default:
			_, err := c.SendMessage(t, "text/html", nil)
			if err != nil {
				color.Red("cannot send message: %s")
			}
		}
	}
}
