package tcpClient

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"log"
	"tcp_chat/types"

	"github.com/marcusolsson/tui-go"
)

func NewUI(name string, rw *bufio.ReadWriter, ep EntryPoint) *UI {
	return &UI{
		Username: name,
		rw:       rw,
		ep:       ep,
	}
}

type UI struct {
	Username string
	rw       *bufio.ReadWriter
	ep       EntryPoint
}

func (ui *UI) StartUI() {
	sidebar := tui.NewVBox(
		tui.NewLabel("USER"),
		tui.NewLabel(ui.Username),

		tui.NewSpacer(),
	)
	sidebar.SetBorder(true)

	history := tui.NewVBox()
	historyScroll := tui.NewScrollArea(history)
	historyScroll.SetAutoscrollToBottom(true)

	historyBox := tui.NewVBox(historyScroll)
	historyBox.SetBorder(true)

	input := tui.NewEntry()
	input.SetFocused(true)
	input.SetSizePolicy(tui.Expanding, tui.Maximum)

	inputBox := tui.NewHBox(input)
	inputBox.SetBorder(true)
	inputBox.SetSizePolicy(tui.Expanding, tui.Maximum)

	chat := tui.NewVBox(historyBox, inputBox)
	chat.SetSizePolicy(tui.Expanding, tui.Expanding)

	input.OnSubmit(func(e *tui.Entry) {
		messageCommand := types.MessageCommand{
			Message: e.Text(),
			Name:    ui.Username,
		}
		input.SetText("")

		ui.ep.SendMessage(messageCommand, ui.rw)

	})

	root := tui.NewHBox(sidebar, chat)

	uia, err := tui.New(root)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			var messageCommand types.MessageCommand
			enc := gob.NewDecoder(ui.ep.Conn)
			enc.Decode(&messageCommand)

			uia.Update(func() {
				history.Append(tui.NewHBox(
					// tui.NewLabel(m.time),
					tui.NewPadder(1, 0, tui.NewLabel(fmt.Sprintf("[%s]", messageCommand.Name))),
					tui.NewLabel(messageCommand.Message),
					tui.NewSpacer(),
				))
			})

		}

	}()

	uia.SetKeybinding("Esc", func() { uia.Quit() })

	if err := uia.Run(); err != nil {
		log.Fatal(err)
	}
}
