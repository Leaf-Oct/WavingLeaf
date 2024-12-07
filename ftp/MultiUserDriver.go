package ftp

import (
	"errors"
	"io"
	"io/fs"

	"goftp.io/server/v2"
)

type MultiUserDriver struct {
	drivers map[string]server.Driver
}

var _ server.Driver = &MultiUserDriver{}

func NewMultiUserDriver(drivers map[string]server.Driver) server.Driver {
	return &MultiUserDriver{
		drivers: drivers,
	}
}

func getUser(ctx *server.Context) string {
	return ctx.Sess.LoginUser()
}

// DeleteDir implements server.Driver.
func (m *MultiUserDriver) DeleteDir(ctx *server.Context, path string) error {
	driver := m.drivers[getUser(ctx)]
	if driver == nil {
		return errors.New("No such a user " + getUser(ctx))
	}
	return driver.DeleteDir(ctx, path)
}

// DeleteFile implements server.Driver.
func (m *MultiUserDriver) DeleteFile(ctx *server.Context, path string) error {
	driver := m.drivers[getUser(ctx)]
	if driver == nil {
		return errors.New("No such a user " + getUser(ctx))
	}
	return driver.DeleteFile(ctx, path)
}

// GetFile implements server.Driver.
func (m *MultiUserDriver) GetFile(ctx *server.Context, path string, offset int64) (int64, io.ReadCloser, error) {
	driver := m.drivers[getUser(ctx)]
	if driver == nil {
		return 0, nil, errors.New("No such a user " + getUser(ctx))
	}
	return driver.GetFile(ctx, path, offset)
}

// ListDir implements server.Driver.
func (m *MultiUserDriver) ListDir(ctx *server.Context, path string, callback func(fs.FileInfo) error) error {
	driver := m.drivers[getUser(ctx)]
	if driver == nil {
		return errors.New("No such a user " + getUser(ctx))
	}
	return driver.ListDir(ctx, path, callback)
}

// MakeDir implements server.Driver.
func (m *MultiUserDriver) MakeDir(ctx *server.Context, path string) error {
	driver := m.drivers[getUser(ctx)]
	if driver == nil {
		return errors.New("No such a user " + getUser(ctx))
	}
	return driver.MakeDir(ctx, path)
}

// PutFile implements server.Driver.
func (m *MultiUserDriver) PutFile(ctx *server.Context, path string, data io.Reader, offset int64) (int64, error) {
	driver := m.drivers[getUser(ctx)]
	if driver == nil {
		return 0, errors.New("No such a user " + getUser(ctx))
	}
	return driver.PutFile(ctx, path, data, offset)
}

// Rename implements server.Driver.
func (m *MultiUserDriver) Rename(ctx *server.Context, fromPath string, toPath string) error {
	driver := m.drivers[getUser(ctx)]
	if driver == nil {
		return errors.New("No such a user " + getUser(ctx))
	}
	return driver.Rename(ctx, fromPath, toPath)
}

// Stat implements server.Driver.
func (m *MultiUserDriver) Stat(ctx *server.Context, path string) (fs.FileInfo, error) {
	driver := m.drivers[getUser(ctx)]
	if driver == nil {
		return nil, errors.New("No such a user " + getUser(ctx))
	}
	return driver.Stat(ctx, path)
}
