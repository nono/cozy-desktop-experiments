package client

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// Stack can be used to control a cozy-stack running locally for tests.
type Stack struct {
	Port  int
	FSDir string
	Cmd   *exec.Cmd
}

// NewStack builds a stack object, but you have to call Start on it to run the
// cozy-stack process.
func NewStack(port int, fsdir string) *Stack {
	return &Stack{Port: port, FSDir: fsdir}
}

// AdminPort returns the admin port as a string.
func (s *Stack) AdminPort() string {
	// keep the same interval as 6060 - 8080
	return fmt.Sprintf("%d", s.Port-2020)
}

// Start executes the cozy-stack serve command, and waits that the stack is
// listening to http requests.
func (s *Stack) Start() error {
	cmd := exec.Command("cozy-stack", "serve",
		"--port", fmt.Sprintf("%d", s.Port),
		"--admin-port", s.AdminPort(),
		"--fs-url", fmt.Sprintf("file://%s/", s.FSDir),
	)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	// Wait that the stack is ready
	buf := make([]byte, 16)
	if _, err := stdout.Read(buf); err != nil {
		return err
	}
	go func() { _, _ = io.Copy(os.Stdout, stdout) }() // TODO use a log file
	s.Cmd = cmd
	return nil
}

// Stop stops the cozy-stack process.
func (s *Stack) Stop() error {
	if s.Cmd == nil {
		return nil
	}
	if err := s.Cmd.Process.Kill(); err != nil {
		panic(err)
	}
	err := s.Cmd.Wait()
	s.Cmd = nil
	return err
}

// CreateInstance creates an instance with a name like "alice".
func (s *Stack) CreateInstance(name string) (*Instance, error) {
	inst := &Instance{Name: name, Stack: s}
	cmd := exec.Command("cozy-stack", "instance", "add", inst.Domain(),
		"--passphrase", "cozy",
		"--public-name", inst.Name,
		"--email", fmt.Sprintf("%s@cozy.tools", name),
		"--admin-port", inst.Stack.AdminPort(),
		"--locale", "en")
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	return inst, nil
}

// Instance describes a cozy instance on a stack.
type Instance struct {
	Name  string
	Stack *Stack
}

// Domain returns the main domain of the instance.
func (inst *Instance) Domain() string {
	return fmt.Sprintf("%s.localhost:%d", inst.Name, inst.Stack.Port)
}

// Address return the http address that can be used to access this instance.
func (inst *Instance) Address() string {
	return fmt.Sprintf("http://%s", inst.Domain())
}

// CreateAccessToken runs a cozy-stack command to create an access token for
// the given client (admin endpoint, no user interaction).
func (inst *Instance) CreateAccessToken(client *Client) error {
	scope := "io.cozy.files"
	cmd := exec.Command("cozy-stack", "instance", "token-oauth",
		"--admin-port", inst.Stack.AdminPort(),
		inst.Domain(), client.ClientID, scope)
	out, err := cmd.Output()
	if err != nil {
		return err
	}
	client.AccessToken = strings.TrimSpace(string(out))
	return nil
}

// Remove runs a cozy-stack command to destroy the instance.
func (inst *Instance) Remove() error {
	cmd := exec.Command("cozy-stack", "instance", "rm", inst.Domain(),
		"--force", "--admin-port", inst.Stack.AdminPort())
	return cmd.Run()
}
