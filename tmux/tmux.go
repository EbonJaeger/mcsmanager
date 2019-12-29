package tmux

import (
	"os"
	"os/exec"
	"strings"
	"syscall"
)

// SessionName is the name of our tmux session.
var SessionName = "MC Server Manager"

// Attach attempts to attach to a currently active tmux window.
func Attach(name string) error {
	// Replace current context with tmux attach session
	tmux, err := exec.LookPath("tmux")
	if err != nil {
		return err
	}
	args := []string{"tmux"}

	// Attach to the session if we're not already in tmux.
	// Otherwise, switch from our current session to the new one
	if os.Getenv("TMUX") == "" {
		args = append(args, "-u", "attach-session", "-t", getWindow(name))
	} else {
		args = append(args, "-u", "switch-client", "-t", getWindow(name))
	}

	// Replace our program context with tmux
	if sysErr := syscall.Exec(tmux, args, os.Environ()); sysErr != nil {
		return err
	}
	return nil
}

// CreateSession starts a named tmux session that runs a single command.
// If a session is already active, a new window for the server will be created.
func CreateSession(command, name string) ([]byte, error) {
	var cmd *exec.Cmd
	if IsSessionRunning() {
		cmd = exec.Command("tmux", "new-window", "-d", "-t", SessionName, "-n", name, command)
	} else {
		cmd = exec.Command("tmux", "new-session", "-d", "-s", SessionName, "-n", name, command)
	}

	return cmd.Output()
}

// Exec creates a command to send keys to the tmux session.
func Exec(command, name string) error {
	cmd := exec.Command("tmux", "send-keys", "-t", getWindow(name), strings.Replace(command, "\"", "\\\"", -1), "C-m")

	return cmd.Run()
}

// IsSessionRunning checks if our Minecraft tmux session is running.
func IsSessionRunning() bool {
	// List all sessions
	sessions, err := ListSessions()
	if err != nil {
		return false
	}

	// Check if our session is listed
	if strings.Contains(sessions, SessionName) {
		return true
	}

	return false
}

// IsServerRunning checks if a window with the given name exists.
func IsServerRunning(name string) bool {
	// List all windows
	w, err := ListWindows()
	if err != nil {
		return false
	}

	// Check if a window with the given name exists
	if strings.Contains(w, name) {
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

// ListWindows gets the list of active tmux windows.
func ListWindows() (string, error) {
	cmd := exec.Command("tmux", "list-windows", "-t", SessionName)
	out, err := cmd.Output()

	return string(out), err
}

// KillWindow closes an active tmux window.
func KillWindow(name string) error {
	cmd := exec.Command("tmux", "kill-window", "-t", getWindow(name))

	return cmd.Run()
}

func getWindow(name string) string {
	return SessionName + ":" + name
}
