package spongedl

import (
	"fmt"
	"github.com/pufferpanel/pufferpanel/v3"
	"os"
	"time"
)

type TestEnvironment struct {
	pufferpanel.Environment
}

func (te *TestEnvironment) Execute(steps pufferpanel.ExecutionData) error {
	return nil
}

func (te *TestEnvironment) ExecuteAsync(steps pufferpanel.ExecutionData) error {
	return nil
}

func (te *TestEnvironment) ExecuteInMainProcess(cmd string) error {
	return nil
}

func (te *TestEnvironment) Kill() error {
	return nil
}

func (te *TestEnvironment) Create() error {
	return nil
}

func (te *TestEnvironment) Delete() error {
	return nil
}

func (te *TestEnvironment) Update() error {
	return nil
}

func (te *TestEnvironment) IsRunning() (isRunning bool, err error) {
	return false, nil
}

func (te *TestEnvironment) WaitForMainProcess() error {
	return nil
}

func (te *TestEnvironment) WaitForMainProcessFor(timeout time.Duration) error {
	return nil
}

func (te *TestEnvironment) GetRootDirectory() string {
	dir, _ := os.Getwd()
	return dir
}

func (te *TestEnvironment) GetConsole() (console []string, epoch int64) {
	return []string{}, 0
}

func (te *TestEnvironment) GetConsoleFrom(time int64) (console []string, epoch int64) {
	return []string{}, 0
}

func (te *TestEnvironment) AddListener(ws *pufferpanel.Socket) {
}

func (te *TestEnvironment) GetStats() (*pufferpanel.ServerStats, error) {
	return nil, nil
}

func (te *TestEnvironment) DisplayToConsole(prefix bool, msg string, data ...interface{}) {
	fmt.Printf(msg, data...)
}

func (te *TestEnvironment) SendCode(code int) error {
	return nil
}

func (te *TestEnvironment) GetBase() *pufferpanel.BaseEnvironment {
	return nil
}
