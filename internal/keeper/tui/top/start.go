package top

// import (
// 	"context"
// 	"gophkeeper/internal/keeper/config"
// 	"sync"

// 	tea "github.com/charmbracelet/bubbletea"
// )

// const logFilePath = "debug.log"

// // Start starts the TUI and blocks until the user exits.
// func Start(cfg *config.Config) error {
// 	m, err := NewModel(cfg)
// 	if err != nil {
// 		return err
// 	}

// 	_, err = tea.LogToFile(logFilePath, "debug")

// 	p := tea.NewProgram(m, tea.WithAltScreen())

// 	ch, unsub := setupSubscriptions()
// 	defer unsub()

// 	// Relay events to model in background
// 	go func() {
// 		for msg := range ch {
// 			p.Send(msg)
// 		}
// 	}()

// 	// Blocks until user quits
// 	_, err = p.Run()
// 	return err
// }

// // // StartTest starts the TUI and returns a test model for testing purposes.
// // func StartTest(t *testing.T, cfg app.Config, width, height int) *teatest.TestModel {
// // 	app, err := app.New(cfg)
// // 	if err != nil {
// // 		return nil
// // 	}
// // 	t.Cleanup(app.Cleanup)

// // 	m, err := newModel(cfg, app)
// // 	require.NoError(t, err)

// // 	ch, unsub := setupSubscriptions(app, cfg)
// // 	t.Cleanup(unsub)

// // 	tm := teatest.NewTestModel(t, m, teatest.WithInitialTermSize(width, height))

// // 	// Relay events to model in background
// // 	go func() {
// // 		for msg := range ch {
// // 			tm.Send(msg)
// // 		}
// // 	}()

// // 	t.Cleanup(func() {
// // 		tm.Quit()
// // 	})
// // 	return tm
// // }

// func setupSubscriptions() (chan tea.Msg, func()) {
// 	ch := make(chan tea.Msg)
// 	wg := sync.WaitGroup{}

// 	_, cancel := context.WithCancel(context.Background())

// 	// {
// 	// 	sub := app.Logger.Subscribe(ctx)
// 	// 	wg.Add(1)
// 	// 	go func() {
// 	// 		for ev := range sub {
// 	// 			ch <- ev
// 	// 		}
// 	// 		wg.Done()
// 	// 	}()
// 	// }

// 	// cleanup function to be invoked when program is terminated.
// 	return ch, func() {
// 		cancel()
// 		// Wait for relays to finish before closing channel, to avoid sends
// 		// to a closed channel, which would result in a panic.
// 		wg.Wait()
// 		close(ch)
// 	}
// }
