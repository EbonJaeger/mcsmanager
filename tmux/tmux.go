package tmux

import (
	"os"
	"os/exec"
	"strings"
	"syscall"
)

// SessionName is the name of our tmux session.
var SessionName = "MC Server Manager"

// Attach attempts to attach to a currently active tmux session.
func Attach() error {
	// Replace current context with tmux attach session
	tmux, err := exec.LookPath("tmux")
	if err != nil {
		return err
	}
	args := []string{"tmux"}

	// Attach to the session if we're not already in tmux.
	// Otherwise, switch from our current session to the new one
	if os.Getenv("TMUX") == "" {
		args = append(args, "-u", "attach-session", "-t", SessionName)
	} else {
		args = append(args, "-u", "switch-client", "-t", SessionName)
	}

	// Replace our program context with tmux
	if sysErr := syscall.Exec(tmux, args, os.Environ()); sysErr != nil {
		return err
	}
	return nil
}

// CreateSession starts a named tmux session that runs a single command.
//
// The session will automatically end when the command ends.
func CreateSession(command string) ([]byte, error) {
	cmd := exec.Command("tmux", "new-session", "-d", "-s", SessionName, command)

	return cmd.Output()
}

// Exec creates a command to send keys to the tmux session.
func Exec(command string) error {
	cmd := exec.Command("tmux", "send-keys", "-t", getWindowName(), strings.Replace(command, "\"", "\\\"", -1), "C-m")

	return cmd.Run()
}

// IsSessionRunning checks if our Minecraft tmux session is running.
func IsSessionRunning() bool {
	sessions, err := ListSessions()
	if err != nil {
		return false
	}

	if strings.Contains(sessions, SessionName) {
		return true
	}

	return false
}

// ListSessions gets the list of active tmux sessions.
func ListSessions() (string, error) {
	cmd := exec.Command("tmux", "list-sessions")
	out, err := cmd.Output()

	return string(out), err
}

// KillSession creates a command to kill the tmux session.
func KillSession() error {
	cmd := exec.Command("tmux", "kill-session", "-t", SessionName)

	return cmd.Run()
}

func getWindowName() string {
	return SessionName + ".0"
}
