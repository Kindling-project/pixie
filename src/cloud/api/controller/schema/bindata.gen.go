// Code generated for package schema by go-bindata DO NOT EDIT. (@generated)
// sources:
// schema.graphql
package schema

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

// Name return file name
func (fi bindataFileInfo) Name() string {
	return fi.name
}

// Size return file size
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}

// Mode return file mode
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}

// Mode return file modify time
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}

// IsDir return file whether a directory
func (fi bindataFileInfo) IsDir() bool {
	return fi.mode&os.ModeDir != 0
}

// Sys return file is sys mode
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _schemaGraphql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xb4\x58\x5f\x6f\xe3\xb8\x11\x7f\xf7\xa7\x98\x5c\x1e\x2e\x01\x92\x3c\x14\xed\xa1\xf0\x53\x55\xdb\x7b\xab\x26\x71\xdc\xd8\xd9\xed\x75\x11\x04\xb4\x38\xb6\x08\x4b\xa4\x96\xa4\x9c\xb8\x87\xfb\xee\xc5\x90\x94\x4c\xda\xce\x1e\xf6\x80\x7b\xb2\x48\xce\x9f\x1f\xe7\x1f\x67\x7c\x0e\x8b\x52\x18\x58\x89\x0a\x81\xa3\x29\xb4\x58\xa2\x01\x5b\x22\x98\xa2\xc4\x9a\xc1\x4a\xab\xda\xad\xb3\x59\x0e\x06\xf5\x56\x14\x78\x33\x38\x1f\x9c\x43\x6e\x7f\x34\x20\x95\x05\xc1\x91\x55\x57\xb0\x6c\x2d\xbc\x22\x48\x44\x0e\x56\x41\xcd\x64\xcb\xaa\x6a\x07\x6b\x94\xa8\x99\x45\xb0\xbb\x06\x0d\xac\x94\x76\xf2\x16\xbb\x06\xe7\x85\x16\x8d\x85\xa7\x7c\x70\x0e\xaf\x25\x4a\xb0\x3d\x18\x61\xa0\x6d\x38\xb3\xc8\x6f\x3c\xc4\x82\x49\x58\x22\x70\x25\x11\x96\x3b\xd0\xad\x94\x42\xae\x87\x83\x73\x80\xb5\x66\x4d\xf9\xb5\xba\xf6\x90\xaf\x9d\x1e\x2f\xb9\xd3\x7d\x6d\x4d\xb8\xd0\x4d\x20\x86\xeb\x6b\xd5\xda\xa6\xb5\xdd\x3e\xbf\xb1\xc6\xc1\x10\x45\x09\xaf\xa2\xaa\x22\xe0\x25\x42\x20\x26\xd9\x1e\xa0\x2d\x99\xf5\x74\x4b\x84\x46\x14\x1b\xe4\xd0\x36\x04\x8d\xc8\x9f\xf2\x9b\x41\xb0\x6d\x24\xdf\x71\x1a\x30\xa5\x6a\x2b\x0e\xf8\x26\x8c\x05\x21\xbd\xb9\x59\x8d\xc0\x85\xc6\xc2\x2a\xbd\x03\x16\x3b\xa1\xc7\x4c\xec\x37\x83\x41\x70\xcd\xaf\x03\x80\xaf\x2d\xea\xdd\x10\xfe\x4d\x3f\x03\x80\xba\xb5\xcc\x0a\x25\x87\x70\x1f\xbe\x06\xbf\x0d\x06\x0e\xf4\x93\x41\x9d\xcb\x95\x72\x6c\x82\x0f\x21\x1f\x9f\x0d\x00\x24\xab\x71\x08\x73\xab\x85\x5c\xd3\x1a\x6b\x26\xaa\x78\xa3\x11\x85\x6d\x75\x42\xa3\xf4\x7a\x9a\xb0\xc5\x3a\xe6\x68\xad\x90\x6b\xa7\x66\x83\xbb\x98\x6f\xcb\xaa\x36\xe5\x42\xd9\xd6\x90\x69\x2b\x56\xac\xb0\x14\x11\x8e\x0d\x20\x5b\xbc\x3c\x4d\x6f\xa7\x0f\x9f\xa7\xdd\xf2\x2e\x9f\x3e\xfd\xe7\x25\xbb\x1f\xff\xf4\xd7\x6e\x6b\x9c\x3d\x7e\xce\xa7\xe9\xde\xe8\x61\xba\xc8\xf2\xe9\xe4\xf1\x65\x3e\x59\xbc\xfc\x92\xdd\xdf\xcd\x4f\x1f\xc5\xf2\x7a\x20\xad\x55\x85\xaa\x9b\x0a\x2d\x4e\xa4\x15\x76\x37\xb7\xe4\x7f\xc2\x94\x4d\xe6\x11\x24\x5a\xcd\x26\xd3\x71\x3e\xfd\x39\xac\x1e\x9f\xa6\xd3\xfd\xea\x43\x96\xdf\x4d\xc6\x61\xb1\x98\x3c\xde\xe7\xd3\x6c\x31\x19\x9f\xd4\x94\x15\xe4\xa7\xfe\xf2\x59\x72\xf7\x73\xc8\x24\x20\x17\x16\x98\x23\x03\x55\x14\xad\x36\x20\x56\xc0\xa0\x35\xa8\xa1\x64\x06\x6a\xc5\xc5\x4a\x50\xde\x95\x08\x42\xba\x40\xc5\x37\x4b\xc1\x28\xa4\x41\xed\x1c\xa2\x34\x70\xac\xd0\x7d\x17\x25\xd3\xac\xb0\xa8\xcd\x8d\x53\xe2\x02\x55\xc8\xa2\x6a\x39\xa5\xff\xae\x71\x0c\x3e\x32\x37\xb8\x5b\x2a\xa6\x39\x30\xc9\xa1\x61\xc6\x0b\x50\x75\xcd\x24\x77\xec\x84\x78\x32\xce\x17\x1e\x2e\x18\xac\xb0\xd8\xe3\x95\xd5\xee\x34\xe8\xa2\x54\x06\x25\x30\x09\x2c\xb2\x06\x98\x76\xbd\x46\x43\xbc\x37\x1d\x2c\x2e\x0a\x66\x09\x97\x72\x2a\x08\x54\xc2\xe2\x52\x51\xd8\x2e\xaf\x6a\xb5\xf5\x39\x4b\xaa\x7e\x34\x40\xba\xa9\xe8\x28\xb7\x29\xc9\x30\xac\x69\xb4\x6a\xb4\x70\xd9\xcd\x96\xdd\x2d\xe6\x93\xbb\xc9\x68\xf1\x8d\x78\xb8\x15\x92\x87\x70\xb8\x4d\xc2\xe1\xf6\x65\xf6\x30\x0e\x5f\xf3\x4f\xa3\xee\x6b\xf4\x98\xcf\x16\x61\x31\xcd\xee\x27\xf3\x59\x36\x9a\xf4\xe9\x32\xc6\xa6\x52\xbb\x1a\xa5\xbd\xc5\xdd\x41\x5e\x1e\xa4\x4e\xa1\x91\x6a\x61\x66\xef\xcd\x10\x3e\x54\x8a\x59\xda\xa5\x8a\x7d\x9c\x85\xae\x18\x38\x71\x64\x80\x61\x9f\xf9\x67\x61\x27\x64\xa8\xb9\xd8\xe0\xce\x0c\xe1\x8b\xe7\x7f\x3e\xbb\x1c\xc2\x97\x28\x83\x9f\xcf\x7c\x68\x3c\x8c\x1f\x2e\x28\x88\xb4\x90\xea\x72\x08\xf7\x6c\x83\x90\x8f\x41\xe3\xd7\x56\x68\xe4\xa0\x64\x41\xf5\xce\x99\xdd\x80\xda\xa2\x33\x75\xdd\x56\x56\x5c\x17\x55\x6b\x2c\x6a\x30\x6d\xd3\x28\x6d\xc9\xce\x61\xeb\xc2\x5f\xf5\x72\x08\x23\xbf\xd1\x21\x0c\xe7\x04\x2c\x3e\xf9\x73\xd1\x8c\x94\x94\xe8\x02\xf6\x08\xd7\xfe\x68\x8f\x50\x74\x25\xeb\x82\x45\xb5\x6b\x98\x54\x32\x92\x70\x97\x77\x3b\xc4\xd7\xd1\x9a\x9e\x2b\xae\xa2\x97\x7b\x76\xd3\x69\x8a\xa3\xfc\xc2\xe5\x75\x47\x7d\x15\xa2\x7a\xa6\xcc\x10\x72\x69\xaf\x42\xbe\x0d\xdf\x29\x2d\x57\xdd\x4d\x9f\xf2\x71\xac\x31\x22\x7e\x44\xd3\x56\xf6\x50\xed\x07\x81\x15\x3f\xd4\xbd\xa2\xcd\x70\xe5\x93\x39\x72\xe5\x8a\x6e\xe7\x94\x4c\xaf\x89\x98\x5c\x7a\x9a\xfc\xf9\x34\xbc\x84\x7a\xde\xd7\x85\xe7\x81\x0b\x05\xdf\x3c\xd4\x6b\x0d\x28\x79\xa3\x84\xb4\xe6\x0a\x34\xae\xbc\xc7\xb9\x2a\xa8\x74\x40\x51\xa9\x96\xb3\x46\xdc\x34\x5a\xb9\xfa\x51\x89\x2d\x7e\x12\xf8\x4a\x68\xee\xc2\xf7\x3d\x5a\xc6\x99\x65\x3e\xca\x3a\x8a\x91\x92\x16\xa5\x35\x21\x24\x28\x3f\xee\x0e\x8e\x88\xdc\xb7\x1a\x2e\x91\xdc\x57\x2a\xcc\x9f\x9e\x10\x35\x4f\x0e\xce\xfc\x9d\x7c\x45\xa0\xe4\x37\x2e\xbd\xa3\xfa\x40\x0a\x92\x82\xe1\xe5\x27\x34\x91\xf8\x94\xb4\xaf\x0d\xc7\x0e\x77\x85\x82\x5e\x00\xa4\xd6\xac\x66\xd6\x22\x0f\x6f\x88\x30\xd1\x83\x62\x82\xef\x7d\x83\x44\x05\x7c\x89\x28\xa1\x61\xda\x20\xef\xda\x9e\xb4\x2c\xab\xbe\x76\xfb\xba\xcd\x96\x73\xab\x1a\x68\x94\x11\xe4\x47\xf7\x78\xf4\x3a\xf3\x38\xc4\x1c\xfd\xe7\x12\x6d\x89\xfa\x08\x03\xe1\x62\xd4\x4c\x08\x7e\x05\xf8\x86\x45\x6b\xd9\xb2\xc2\xee\x4d\x22\xa9\xc2\x4c\xfa\xfd\x21\xfc\x53\xa9\x0a\x99\xf4\xef\x53\x55\x45\x4f\x8c\x6f\x47\x91\x15\x25\xa8\x95\x53\x14\x40\x3a\x6c\xf4\xbd\x27\x1d\xc2\x97\x45\xbc\xf1\xdc\x1b\x35\xd9\x8e\xec\x29\x24\xc7\xb7\x48\xb0\x7f\xa8\x6c\x89\x06\x13\x0c\x4c\x3b\xdb\x07\x95\x39\x71\xb9\xa4\x4e\xac\xe0\x9f\x55\xba\x3e\x8b\x98\x43\x3b\x4d\x9e\x62\xcb\xa0\xd0\x35\xa5\x35\x15\x46\xd2\x1b\xac\x12\x19\x8a\xf4\xec\x57\xd9\xca\x52\xc9\x27\xe1\xb1\xa5\x4c\x72\xf1\xf7\x12\xf1\x54\x58\x1d\x98\x62\x23\x24\x7f\xaf\x4c\x1c\xf4\x9f\xe1\x45\xa3\xbc\x70\xa5\xac\xdf\xad\x99\x2d\x4a\x0a\x11\x8e\x6f\xae\x8c\xe4\xd2\x3e\x13\x48\x6a\xcf\x4e\x09\x77\x7d\x5b\xff\x8e\x87\x4a\x4e\x9b\xad\x89\xfc\xc3\x71\xc5\x28\x03\x9c\x18\xea\x4e\xa4\xb2\x65\x08\xb0\x8d\x54\xaf\x92\x2c\xf5\xe9\xbf\x2f\xf3\xb4\x23\x23\xd6\xc0\x62\xa0\x44\x56\xd9\x72\x47\xdc\x25\x32\x6d\x97\xc8\xac\xf7\xa8\xc6\x02\xc5\xd6\x3d\x49\xa0\x71\xdd\x56\x4c\x83\x90\x16\xf5\x96\x55\xc6\x35\x53\xb6\xf4\x89\xd1\xbd\x4b\xc2\x80\x46\xd3\x28\xc9\x09\x84\x55\xae\x80\xa2\xb1\x66\x8f\xe3\xe3\x24\xbb\x5b\x7c\xfc\xe5\x00\x87\x1f\x57\x94\xab\x7b\xc2\x14\xfe\xc5\xa2\x34\xf6\xa1\xf7\xf3\xe3\x6c\x04\x45\xff\x8e\xc1\x52\x23\xdb\x98\x1b\x27\xa0\x54\x0d\xfa\x44\x67\xb6\xef\xae\x3a\x40\x4e\x6e\xa1\x6a\x84\x25\x2b\x36\xd4\xcb\x09\x89\x0e\xba\x46\xd3\xd6\x14\xe1\x10\x10\x79\x24\x7b\xa0\xe3\x7c\x3e\x7a\x98\x4e\x27\xa3\x85\xeb\x83\x0f\xac\xe6\x46\x3b\xba\x64\x98\xfa\x30\xb6\x41\x98\x88\x1a\xad\x0a\x34\x86\xf2\xa7\x23\x8f\xfc\x31\x1b\x67\x0b\xdf\x6f\x7b\xd1\x5b\xf1\x3f\xd1\x35\x96\xdd\xfd\xfd\x4c\x4a\x5b\x34\xa6\x1a\x94\x16\x98\xdc\x81\x72\xf9\xb4\x6a\xb5\xcf\x2b\x1f\x15\x7e\xd8\x34\xc0\x96\xaa\xf5\x86\x78\x0d\x89\x27\x6c\xec\x67\xa5\x4f\xa0\x39\xbe\x69\x80\xf3\x4a\x73\x9c\xde\x05\x77\x7a\x1d\x1e\xd5\x8a\x89\x0a\x7d\x4f\x2d\x08\xdf\x2b\x5d\x9b\xc1\x92\xf1\x43\x4b\xba\xab\x4e\xba\x91\xa2\xcb\xb8\x4f\x4e\xc1\x48\xc9\x95\xf0\xb3\x56\xc3\x8c\xb1\xa5\x56\xed\xba\x9c\x48\x4a\x6c\xbe\x4f\xe7\x8e\x89\x5e\x1b\x26\x64\x92\x0a\x87\x03\xe0\xe9\x4e\x33\x64\xda\x9e\xac\x46\x63\xd8\x3a\x4e\x5d\x8d\xcc\x44\x59\xdb\xe9\xbc\xfd\xbb\x99\x6c\xc9\xf4\xbf\x9e\xe4\x5a\x09\x6d\xec\x42\xd4\x98\xa8\xab\xd8\xd1\x66\x27\x6f\xa6\xf8\x1f\x42\xdf\x9a\xef\x82\x0f\x14\x46\xde\x58\xae\x15\x4d\x2d\xe7\xdf\x5e\xa4\x7b\xd1\x69\x77\x47\xda\xee\x6d\xbd\x6f\x5e\x0f\x5a\xfb\x0e\x4d\x52\x96\xba\x5b\x7f\xec\x8a\x48\x72\x81\x6d\xe4\xed\x61\xe2\xfb\xfd\xe9\x27\xd4\x26\x2d\x9b\x21\xa7\xde\x3d\x98\xa6\xb5\xb7\xd1\x68\xed\x6e\x74\xf2\xec\xb8\x3d\x0b\x16\xd2\xaa\x9a\x55\x4c\x62\xef\x16\x57\x9f\xfb\x95\x37\x94\x6c\xeb\xa9\xe2\xe8\x5b\xd5\xb0\x91\x4b\x63\x75\x4b\x0d\x0a\xf2\xf8\xf0\xc0\x7e\x69\xfb\xed\x2d\xd9\x64\x9c\x6b\x34\x89\x43\xad\xda\xa0\x3c\x9e\x85\xba\xbf\x41\x1c\xe3\xc8\x05\x47\x10\x9c\x0c\x1e\xf0\x0f\x8e\x8d\x46\x6a\x55\xf8\x45\x17\x09\x3f\x04\x02\x5f\xce\x29\x45\x43\x74\xc1\x56\x30\x68\xde\x42\xef\xf5\xc3\xe5\x00\xe0\xc9\xa5\x76\xec\x98\x8b\x60\x32\xb2\x58\x3e\x3e\xbb\xfa\x56\x82\x5e\xf6\x5f\x67\x3d\xcc\xa4\x81\x3b\xea\xe7\x00\xc6\x34\xcb\xa7\x54\x51\xfb\x17\x89\xf3\xd0\x9e\xbe\x31\xf8\x5d\xf9\x7f\x66\x0e\x66\xc1\x5e\x44\xef\x91\xfd\x38\x13\xfe\xa6\x69\x75\xf2\x57\x11\x80\x29\xd9\x5f\xfe\xf6\xd3\xb1\x1b\x92\xc9\xc6\x3b\xd1\x62\xed\x3a\x8b\x70\xf2\x7c\x44\xeb\xc8\xb6\x69\xec\xba\x14\x2f\x99\x5c\x63\xa5\xd6\x89\xfb\x45\x8d\xc6\xb2\xba\x49\x2b\xc6\x39\x3c\xfe\xce\x40\xe0\x54\x1e\xce\x01\xbf\xf3\x1f\xd9\xd1\xd8\xfd\x9d\x6a\xba\xa6\x3f\x94\x44\xaf\x73\x78\x84\xc2\xfd\xfb\xf6\x56\x75\xd4\xc9\x3f\x69\xc2\xfc\x6b\xfe\x30\xfd\x23\x20\xd2\x21\xe5\xbb\x6e\x0a\xf4\x96\x76\x28\xd3\x00\xf9\x2e\xe5\xef\xdc\xff\x60\x7c\x0a\x05\x26\xbd\xfa\x6f\x83\xff\x07\x00\x00\xff\xff\x8c\x0e\x52\x3e\xb2\x16\x00\x00")

func schemaGraphqlBytes() ([]byte, error) {
	return bindataRead(
		_schemaGraphql,
		"schema.graphql",
	)
}

func schemaGraphql() (*asset, error) {
	bytes, err := schemaGraphqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "schema.graphql", size: 5810, mode: os.FileMode(436), modTime: time.Unix(1606245061, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"schema.graphql": schemaGraphql,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"schema.graphql": &bintree{schemaGraphql, map[string]*bintree{}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
