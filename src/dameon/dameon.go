package dameon

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"path/filepath"
	"time"

	"../utils"

	"github.com/nightlyone/lockfile"
)

// Dameon is the go representation of the single-instance backend to gop.
type Dameon struct {
	pm *PM

	lockfile lockfile.Lockfile
	listener net.Listener

	loggingToFile bool
	logFileName   string
}

// Start launches the rpc server for the dameon and makes sure it is single instance
func (d *Dameon) Start(netinterface string) error {

	d.pm = NewPM()
	go d.monitorProcs()

	var err error
	var absPath string

	// Ensure single instance
	if absPath, err = filepath.Abs(utils.Options.LockFile); err != nil {
		return err
	}
	if d.lockfile, err = lockfile.New(absPath); err != nil {
		return err
	}

	if err = d.lockfile.TryLock(); err != nil {
		return err
	}

	// Start rpc
	gop := NewGOP(d)
	rpc.Register(gop)
	rpc.HandleHTTP()

	if d.listener, err = net.Listen("tcp", netinterface); err != nil {
		return err
	}

	return http.Serve(d.listener, nil)
}

func (d *Dameon) monitorProcs() {
	for procClose := range d.pm.ProcessExits {
		log.Println("Process ended [", procClose.Proc.InstanceName, "]:", procClose.Error)
	}
}

func (d *Dameon) SetLogging() error {

	t := time.Now()
	d.logFileName = filepath.Join("logs", fmt.Sprintf(
		"%d-%02d-%02d %02d-%02d-%02d.log",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second()))

	if err := os.MkdirAll(utils.Options.LogDirectory, 0777); err != nil {
		return err
	}

	logFile, err := os.Create(d.logFileName)
	d.loggingToFile = err == nil
	if d.loggingToFile {
		log.SetOutput(io.MultiWriter(os.Stdout, logFile))
	}

	return err
}

// Stop shuts down the rpc server and releases the lock file.
func (d *Dameon) Stop(closeProcs bool) {
	d.lockfile.Unlock()
	d.listener.Close()
	d.pm.Shutdown(closeProcs)
}
