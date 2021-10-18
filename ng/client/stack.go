package client

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

type Stack struct {
	Port  int
	FSDir string
	Cmd   *exec.Cmd
}

func NewStack(port int, fsdir string) *Stack {
	return &Stack{Port: port}
}

func (s *Stack) AdminPort() string {
	// keep the same interval as 6060 - 8080
	return fmt.Sprintf("%d", s.Port-2020)
}

func (s *Stack) Start() error {
	cmd := exec.Command("cozy-stack", "serve",
		"--port", fmt.Sprintf("%d", s.Port),
		"--admin-port", s.AdminPort(),
		"--fsurl", fmt.Sprintf("file://%s/", s.FSDir),
	)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	// Wait that the stack is ready
	buf := make([]byte, 1024)
	if _, err := stdout.Read(buf); err != nil {
		return err
	}
	io.Copy(os.Stdout, stdout) // TODO use a log file
	s.Cmd = cmd
	return nil
}

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

type Instance struct {
	Name  string
	Stack *Stack
}

func (inst *Instance) Domain() string {
	return fmt.Sprintf("%s.localhost:%d", inst.Name, inst.Stack.Port)
}

func (inst *Instance) Remove() error {
	cmd := exec.Command("cozy-stack", "instance", "rm", inst.Domain(),
		"--force", "--admin-port", inst.Stack.AdminPort())
	return cmd.Run()
}
