package pty

import (
	"fmt"
	"os"
	"os/exec"
	"sync"

	"github.com/creack/pty"
)

const SHELL = "bash"

type Session struct {
	Terminal *os.File
	ReplID   string
}

type TerminalManager struct {
	sessions map[string]*Session
	mu       sync.Mutex
}

func NewTerminalManager() *TerminalManager {
	return &TerminalManager{
		sessions: make(map[string]*Session),
	}
}

func (tm *TerminalManager) CreatePty(terminalID string, replID string, onData func(data string, terminalID string)) (*os.File, error) {
	cmd := exec.Command(SHELL)
	ptmx, err := pty.Start(cmd)
	if err != nil {
		return nil, err
	}

	tm.mu.Lock()
	tm.sessions[terminalID] = &Session{
		Terminal: ptmx,
		ReplID:   replID,
	}
	tm.mu.Unlock()

	// Read from PTY and send data back via onData
	go func() {
		defer ptmx.Close()
		buf := make([]byte, 1024)
		for {
			n, err := ptmx.Read(buf)
			if err != nil {
				break
			}
			onData(string(buf[:n]), terminalID)
		}
	}()

	return ptmx, nil
}

func (tm *TerminalManager) Write(terminalID string, data string) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	session, exists := tm.sessions[terminalID]
	if !exists {
		return fmt.Errorf("terminal session not found")
	}

	_, err := session.Terminal.Write([]byte(data))
	return err
}

func (tm *TerminalManager) Clear(terminalId string) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	session, exists := tm.sessions[terminalId]
	if !exists {
		return fmt.Errorf("session not found")
	}

	err := session.Terminal.Close()
	delete(tm.sessions, terminalId)
	return err
}
