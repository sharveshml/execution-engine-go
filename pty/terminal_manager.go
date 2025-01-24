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

func (tm *TerminalManager) CreatePty(id string, replId string, onData func(data string, id int)) (*os.File, error) {
	cmd := exec.Command(SHELL)

	ptmx, err := pty.Start(cmd)
	if err != nil {
		return nil, err
	}

	tm.mu.Lock()
	tm.sessions[id] = &Session{
		Terminal: ptmx,
		ReplID:   replId,
	}
	tm.mu.Unlock()

	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := ptmx.Read(buf)
			if err != nil {
				break
			}
			onData(string(buf[:n]), cmd.Process.Pid)
		}
	}()

	go func() {
		cmd.Wait()
		tm.mu.Lock()
		delete(tm.sessions, id)
		tm.mu.Unlock()
	}()

	return ptmx, nil
}

func (tm *TerminalManager) Write(terminalId string, data string) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	session, exists := tm.sessions[terminalId]
	if !exists {
		return fmt.Errorf("session not found")
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
