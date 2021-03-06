// Code generated by go-bindata.
// sources:
// sql/20180124_init.sql
// sql/20180126_add_accounts.sql
// sql/20180127_sensitive_data.sql
// sql/20180127_z_1_null_and_default.sql
// sql/20180127_z_2_recurring_groups.sql
// sql/20180128_1_groups_finished.sql
// sql/20180303_1_transactions_edited_name.sql
// sql/20180729_1_recurring_to_groups.sql
// sql/20180729_2_trend_id.sql
// sql/20180729_3_groups_constraint.sql
// sql/20180729_4_group_trend_id.sql
// sql/20180729_5_trends.sql
// DO NOT EDIT!

package migrations

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

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _sql20180124_initSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x94\x94\xc1\x6f\xd3\x30\x14\xc6\xcf\xf1\x5f\xe1\x63\x2b\x32\x89\x01\xb7\x9d\xb2\x2d\x82\x08\x3a\x4a\x48\x91\xc6\xc5\x7a\xb3\x5f\x57\xab\x8d\x1d\xbd\x38\x65\xf9\xef\x51\x16\x95\x39\x76\xca\xb4\x53\x94\xdf\xf7\xd9\xcf\xcf\xef\x93\x2f\x2e\xf8\xbb\x5a\x3f\x12\x38\xe4\x9b\x86\xdd\x94\x79\x56\xe5\xbc\xca\xae\xbf\xe5\xdc\x6e\x9f\x84\x23\x30\x2d\x48\xa7\xad\x69\xf9\x82\x25\x5a\xf1\x5f\x59\x79\xf3\x25\x2b\x53\x96\xb4\xb6\x23\x89\x11\x10\x20\xa5\xed\x8c\x13\x33\xe6\x7f\x9a\xeb\x1b\x7f\xa5\x57\x27\x94\x14\x38\x14\x8d\x6d\x1d\x2a\x5e\x15\xab\xfc\x67\x95\xad\xd6\xd5\xef\x93\xd2\xb5\x48\x73\x1c\x8e\xa0\x0f\xf0\x70\xc0\x40\x84\x7a\xa8\xcf\xaf\x8b\xcf\xc5\x5d\x95\xb2\x64\xab\x85\x9b\x1e\x14\xe9\x88\x14\x40\xb9\x43\xb9\x17\xa6\xab\x3d\x46\xb8\x0d\x48\xab\xe5\xcb\xc6\x0d\xf4\x88\xd3\x3b\x30\x50\xfb\x9d\xe1\x93\x43\xa3\x50\x89\x80\xd7\x58\x5b\xef\x57\x9b\xa3\xf8\xf4\xfe\x72\x2f\xa2\xfb\x96\x1d\x11\x1a\xd9\x7b\xc8\x92\x7e\xd4\x06\x0e\x62\x46\x93\x96\x08\xa5\x13\x51\xcb\x27\x61\x1c\x80\x27\x8c\x3d\x04\xc7\x1b\x21\x28\x45\x97\xb3\xf4\xc3\x2c\xfd\x18\x51\xa9\x5d\x1f\xc1\xd6\x0d\x51\x0c\xe9\x30\xfe\xa1\x27\xab\x62\xed\x39\x4f\x14\xef\xd4\xec\xac\xf1\xdd\x0f\x60\xf6\x2f\xf1\xb3\xe2\xf9\x7f\x72\x0d\x91\x83\xc0\xc8\xdd\xff\x3d\xb3\x59\x3f\x67\x0a\x92\x7d\xce\xb6\xc7\xc9\xd4\x08\x95\x76\x42\x02\xa9\x57\x0b\xbf\xe2\x9d\x6c\xcc\x92\x75\x59\xac\xb2\xf2\x9e\x7f\xcd\xef\x17\x5a\xa5\x7c\xcc\xd7\xe9\xeb\x55\x88\xd0\xd0\xc8\x32\x65\xc9\xe6\xae\xf8\xb1\xc9\x17\x6f\x5b\x98\xf2\x31\x80\x4b\xb6\xbc\x62\xcc\x7f\x81\x6e\xed\x1f\xc3\x6e\xcb\xef\xeb\x33\x2f\xd0\x15\xfb\x1b\x00\x00\xff\xff\x65\x4e\x6e\xc3\xb0\x04\x00\x00")

func sql20180124_initSqlBytes() ([]byte, error) {
	return bindataRead(
		_sql20180124_initSql,
		"sql/20180124_init.sql",
	)
}

func sql20180124_initSql() (*asset, error) {
	bytes, err := sql20180124_initSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "sql/20180124_init.sql", size: 1200, mode: os.FileMode(436), modTime: time.Unix(1516946318, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _sql20180126_add_accountsSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xb4\xd1\xc1\x6e\xc3\x20\x0c\x06\xe0\x73\xfc\x14\x3e\xb6\x5a\xfa\x04\x39\xb1\x04\x69\xd5\xd2\xa4\x62\x74\x52\x4f\x88\x25\xb4\x42\x5d\x21\x03\xa2\xad\x6f\x3f\x65\xa3\x1a\x4d\xcf\x3b\x21\x7d\xb2\xec\xdf\x66\xb5\xc2\x87\xb3\x3e\x3a\x19\x14\xee\x06\x28\x19\x25\x9c\x22\x27\x8f\x35\x45\xd9\x75\x76\x34\xc1\xe3\x02\x32\xdd\x63\x96\xbd\x12\x56\x3e\x11\x96\x43\x66\xe4\x59\xdd\x40\xac\x15\xe1\x32\x28\x4c\xdc\xa9\x8f\x51\xf9\x5f\x4f\xf8\x4d\x9a\x93\xb0\xee\x38\xa7\x83\xee\xe7\x34\xba\xf7\x84\xb6\x6c\xbd\x21\x6c\x8f\xcf\x74\xbf\xd0\xfd\x12\x96\x05\x00\xa9\x39\x65\x31\xb3\x3d\x7c\x89\xe0\xa4\xf1\xb2\x0b\xda\x1a\x8f\x15\x6b\xb7\x58\xb6\xcd\x0b\x67\x64\xdd\xf0\xbb\x02\x31\x9c\xd4\x25\x87\x8c\x54\x15\x26\xcd\x71\xea\x1e\xb9\x6c\xeb\xdd\xa6\xb9\x9e\x43\xe8\x1e\x63\x9e\x02\x20\xbd\x5f\x65\x3f\x0d\xfc\xcc\xbb\xbd\x5f\xf1\x6f\x09\x73\xf4\x76\x74\x9d\xba\xbe\xe2\x2f\xe3\x1d\x4d\x3f\x30\x6d\x14\xc7\xcd\x56\x2a\xe0\x3b\x00\x00\xff\xff\x67\xde\xa9\x16\x0a\x02\x00\x00")

func sql20180126_add_accountsSqlBytes() ([]byte, error) {
	return bindataRead(
		_sql20180126_add_accountsSql,
		"sql/20180126_add_accounts.sql",
	)
}

func sql20180126_add_accountsSql() (*asset, error) {
	bytes, err := sql20180126_add_accountsSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "sql/20180126_add_accounts.sql", size: 522, mode: os.FileMode(436), modTime: time.Unix(1517040375, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _sql20180127_sensitive_dataSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x94\x92\x41\x4e\xeb\x30\x10\x86\xd7\xcf\xa7\x98\xe5\x43\xd0\x13\x64\xe5\x36\x91\xa8\x14\x1c\x30\x09\x5b\xcb\x72\xdc\xca\x02\xec\xc8\x9e\x00\xbe\x3d\xaa\x9a\x54\x0e\x49\x5b\x58\x45\x9a\x7c\x1a\xcf\x7c\xff\xac\x56\x70\xfb\x6e\xf6\x5e\xa2\x86\xa6\x23\xb4\xac\x0b\x0e\x35\x5d\x97\x05\xb8\xdd\x97\x40\x2f\x6d\x90\x0a\x8d\xb3\x01\x72\x5e\x3d\xc2\xa6\x62\xcf\x35\xa7\x5b\x56\xcf\x00\x11\x5c\xef\x95\x1e\x3f\x52\x29\xd7\x5b\x14\xa6\xfd\x59\xc1\x28\x5e\x75\xbc\x23\xff\x86\x8e\x65\xf3\xc0\xe0\xc8\x2c\x16\x93\x56\x97\xff\x63\xec\x0e\x1d\x68\x9e\x8f\x80\xd7\xaa\xf7\xde\xd8\xbd\x30\x2d\xbc\x50\xbe\xb9\xa7\xfc\x44\x9c\xdf\xa4\xf3\xfa\x43\x5b\x14\x6d\xdf\xbd\x19\x25\x51\x07\x68\xd8\xf6\xa9\x29\xe0\x7f\x32\x0b\xec\x8c\x40\xd3\xde\x64\x64\x22\x6e\x20\x02\x24\x73\x84\x68\x15\xac\xab\xaa\x2c\x28\xcb\x08\x49\xb5\xe7\xee\xd3\x5e\x16\x7f\x6d\xda\xbf\x79\x3f\x6d\x32\x18\x5f\x90\xbc\xe8\x75\xdc\x76\x2a\xf8\x48\xce\xd4\x2e\xa7\xf7\x2b\xec\xf0\x58\x02\x5e\xbd\xba\x79\x56\xe7\xf2\x98\x5c\x4e\xb4\x2a\x23\xdf\x01\x00\x00\xff\xff\x8a\x60\x67\xee\xff\x02\x00\x00")

func sql20180127_sensitive_dataSqlBytes() ([]byte, error) {
	return bindataRead(
		_sql20180127_sensitive_dataSql,
		"sql/20180127_sensitive_data.sql",
	)
}

func sql20180127_sensitive_dataSql() (*asset, error) {
	bytes, err := sql20180127_sensitive_dataSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "sql/20180127_sensitive_data.sql", size: 767, mode: os.FileMode(436), modTime: time.Unix(1517123581, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _sql20180127_z_1_null_and_defaultSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xa4\x97\xcf\x72\xe2\x38\x10\x87\xcf\x9b\xa7\xf0\x2d\xbb\xb5\x4b\x95\x81\xcd\x64\x8a\x39\xe5\x0f\x73\x62\x92\xa9\x0c\x9c\x55\x8a\x2c\x13\x17\x20\x31\xb2\x48\xc2\xdb\x4f\x19\xb0\x91\x6c\xab\xbb\x6d\x6e\x49\xf1\xb5\xe4\x6a\xcb\xfa\x7e\x3d\x18\x44\xff\x6e\xb2\xa5\xe1\x56\x46\x8b\xed\xd5\xdd\x6c\x3e\x7d\x89\xe6\x77\xf7\xb3\x69\xa4\xd3\x4f\x66\x0d\x57\x39\x17\x36\xd3\x2a\x8f\x8e\x3f\x3e\x3c\xcf\x16\x3f\x9e\xa2\x2c\x89\x7e\x4d\xe7\xd1\xe3\xf4\xfb\xdd\x62\x36\x8f\xae\xaf\xff\xbb\xfa\xcb\x03\x9c\x52\x66\xf7\x5b\x89\xe0\x09\xb7\x92\x6d\x75\x6e\xa5\xbf\xb0\xd2\x1f\x7f\xff\xd3\x0a\xef\x72\x69\xfc\x45\xe3\x38\x1e\x0e\x86\xa3\xc1\x78\x18\x0d\xbf\x4c\xe2\xdb\x49\x3c\x1a\xc4\xb7\x93\x9b\xd1\xe4\xe6\x6b\x74\xff\xd0\xbe\x27\x7f\xe7\xd9\x9a\xbf\xae\xe5\x65\x6b\xf1\x8d\xde\x29\xeb\xad\x11\xd7\x99\x34\x63\x16\xed\x5b\x2e\xcd\xbb\x34\x04\x50\xbc\x49\xb1\x62\x6a\xb7\x41\x38\x23\x53\x02\x95\x67\x02\x7e\xf8\x2d\xdf\x4b\xc9\xd0\xa7\x52\x7c\x83\xbd\x6a\xf9\x69\xa5\x4a\x64\xc2\x08\xec\x46\x6e\x34\x82\x64\xea\x9d\xfd\x1f\x0f\x57\x2c\xd7\x3b\x23\xb0\x05\xc5\xce\x18\xa9\xc4\x1e\xc1\xb4\xc9\x96\x99\xe2\x6b\x46\xe4\x85\x36\x46\x0a\xcb\x48\xaf\xb8\x84\x8f\x5f\x07\x02\x1f\xfb\x4e\x68\xd5\x11\xe4\x49\x62\x86\x64\x72\x44\x26\xc7\x24\x52\x64\x16\xeb\xd4\x11\xcc\x6d\x71\xe3\x50\xc8\xe2\x4e\x28\xde\x83\x4e\x68\xbc\x28\xbe\x43\x43\x7b\x8a\xed\x9b\x56\xd8\xaa\xaf\x5c\xad\x18\x17\x87\x65\x99\xd5\xec\xf0\x3f\xfa\x8a\x1b\x55\x86\x2b\xf1\xd6\xbd\xae\xfc\xb3\x77\x21\xe1\xea\x0d\x95\xae\x24\x7a\xea\x8d\x4c\x32\xcb\x04\x37\x49\xaf\x87\x46\xea\xf1\x07\x20\xef\x64\x64\xf1\x21\x67\x6a\xd9\x44\xbf\x5d\x75\xd7\xde\xd3\xf3\x3c\x7a\x5a\xcc\x66\x24\xe9\x85\xe0\xba\xf2\x40\xae\xb2\x1d\x48\xf9\x32\x0b\xa1\x8e\xab\x42\x88\x73\x8f\x85\x90\x9a\xa8\x42\x98\xaf\xa9\x10\xe5\x4a\x2a\xb8\xe1\x49\x51\xa1\xdf\x3d\x43\x85\xa0\xea\x22\x0d\x01\x4d\x3b\x85\xc8\xca\x4d\x21\xa0\xcd\x4c\xc1\x36\xb9\x9e\x09\x41\xed\x56\x0a\x2e\xd9\x74\x12\x86\x3a\x46\x82\xbb\x8c\x36\xa7\x6e\x23\x9c\x1b\x11\xb9\x31\x81\xab\x3c\x04\x63\x67\x0b\xc1\x5c\xdd\x41\xc8\xe6\x8e\x81\x90\x75\x2b\xff\x84\x38\xc8\x3e\xe4\x1a\xcf\x3d\xd4\xaa\xda\xd5\xda\xb5\x0c\xbd\xfd\x30\xeb\x04\x8f\x2a\xc9\x39\x3d\xab\xb1\xcd\x89\xbb\x34\x6c\x53\x82\x35\xd7\x9c\x96\xeb\x3a\x5a\x11\xc2\x60\x07\xfd\x1b\xf9\x7b\x27\x73\x7a\x52\xd0\x66\x49\xc1\x52\x5a\x64\xd9\x99\x35\x36\x99\xec\x95\x3f\x9a\xa4\x7c\x9d\xcb\x2e\x9d\xec\xed\x02\xf2\x61\x6e\xf4\x10\x3c\xf5\x65\x07\x41\x28\xa5\x7c\x78\x65\xf7\x82\xca\x2c\x7b\xe7\x1c\x40\x77\xe4\x7f\xd4\x1f\xaa\x5b\xfa\x79\x7c\x79\xfe\x59\xbe\x09\x34\xfd\x40\xb0\x9b\x7e\x50\xee\x90\x7e\x50\xea\x9c\x7e\x20\xf4\x94\x7e\x20\xe4\x64\x4c\x08\x71\xd2\x0f\x84\x9d\xd3\x0f\x44\x95\xe9\x07\xdc\x30\x13\xe0\xef\x55\xfa\x81\xa0\xc3\x89\x87\x00\x3f\xfd\x40\xe4\x21\xfd\x40\x40\x3d\xfd\x80\x6d\x2a\xf3\x0c\x04\x35\xd3\x0f\xb8\xa4\x9f\x7e\x28\xe8\x29\xfd\xe0\x5d\x46\x9b\xe3\xa6\x1f\x1a\x37\x22\x72\x63\x02\x77\x48\x3f\x38\x76\x4c\x3f\x38\xe7\xa6\x1f\xc2\xe6\xa7\xf4\x43\x58\xf7\x90\x7e\x20\x2e\x94\x7e\x3a\xd5\x54\xe9\xa7\x4b\x95\x23\xfa\x3e\x65\xe8\xed\x07\xa5\x1f\xf0\xa8\xa2\xe9\xe7\x82\x6a\x6c\x73\xe2\x2e\x5e\xfa\x71\xc1\x4e\x93\x76\xbb\x4c\x28\xb3\x34\x0e\x9e\x75\x42\x9c\xa6\x41\xd6\x15\x0a\x32\x4f\x83\x4c\x5d\x29\xf8\x44\x0d\x62\x9e\x54\xa0\x99\x1a\x04\x7c\xad\x80\x49\x0a\x24\x5a\xc4\x02\xce\xd5\x20\xd1\xaa\x16\x74\xb2\x06\xa9\x80\x5c\x88\xb3\x35\x89\x75\xf5\x42\x98\xae\x09\x9c\x23\x18\xca\x7c\x4d\x04\xc7\x14\xf0\xec\x18\xca\x88\x4d\x00\x1b\x96\xa1\x0d\xd9\x94\x95\xcf\x9e\xe9\x3a\x66\x77\x2b\xf2\x4d\xd3\x63\xd0\xee\x55\x87\x5f\x8e\xa8\x6c\xfa\xcf\xda\x97\x94\xa3\xdb\x53\xf7\x69\x0a\xa7\xeb\xbc\x7d\x51\x66\x26\x4b\xdf\x9b\x12\xd1\x74\x50\x4c\x89\x28\x94\x52\x02\x4a\x31\x25\x82\xa3\x45\x31\x25\x02\xa2\x6e\x6f\x1b\x7e\x49\xd1\x4f\x68\xb3\x33\xe8\x00\x8d\x53\x29\xe9\x8b\xaa\x9a\x03\xce\xd0\xb5\x53\xf5\x27\x00\x00\xff\xff\xf7\x87\x40\x4a\x37\x1f\x00\x00")

func sql20180127_z_1_null_and_defaultSqlBytes() ([]byte, error) {
	return bindataRead(
		_sql20180127_z_1_null_and_defaultSql,
		"sql/20180127_z_1_null_and_default.sql",
	)
}

func sql20180127_z_1_null_and_defaultSql() (*asset, error) {
	bytes, err := sql20180127_z_1_null_and_defaultSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "sql/20180127_z_1_null_and_default.sql", size: 7991, mode: os.FileMode(436), modTime: time.Unix(1517121326, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _sql20180127_z_2_recurring_groupsSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x84\xce\xc1\x0a\x82\x40\x14\x85\xe1\xf5\xdc\xa7\xb8\x3b\x95\xf2\x09\x5c\x4d\x3a\x51\x34\xa9\x5c\xc6\x40\x22\x44\x54\x64\x16\x8e\x32\x3a\xb4\x88\xde\x3d\x84\x16\x2d\x8a\xce\xfa\x3b\xf0\x87\x21\x6e\x06\xdd\xdb\x7a\xe9\xb0\x98\x20\x26\xc1\x95\x40\xc5\x77\x52\xa0\xed\x1a\x67\xad\x36\x7d\xd5\xdb\xd1\x4d\x33\xfa\xc0\x74\x8b\x8c\x5d\x38\xc5\x07\x4e\x98\x66\x0a\xd3\x42\x4a\x4c\xc4\x9e\x17\x52\xa1\xe7\x6d\x81\x99\x7a\xe8\xfe\xa2\xba\x69\x46\x67\x96\x4a\xb7\x33\xae\x7b\xeb\xeb\xed\x8b\x7f\x3c\xd7\x47\x4e\xc7\x33\xa7\x12\x4f\xa2\xf4\x75\x1b\x40\x10\x01\x7c\xd6\x27\xe3\xdd\x40\x42\x59\xfe\xa3\x3e\x82\x57\x00\x00\x00\xff\xff\xba\x36\x44\xe3\xec\x00\x00\x00")

func sql20180127_z_2_recurring_groupsSqlBytes() ([]byte, error) {
	return bindataRead(
		_sql20180127_z_2_recurring_groupsSql,
		"sql/20180127_z_2_recurring_groups.sql",
	)
}

func sql20180127_z_2_recurring_groupsSql() (*asset, error) {
	bytes, err := sql20180127_z_2_recurring_groupsSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "sql/20180127_z_2_recurring_groups.sql", size: 236, mode: os.FileMode(436), modTime: time.Unix(1517121557, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _sql20180128_1_groups_finishedSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x7c\xcd\x31\x0e\xc2\x20\x14\x06\xe0\x9d\x53\xfc\xbb\xe9\x09\x3a\x51\x1f\x4e\x4f\x9e\x69\x60\x36\x8d\x52\x24\x51\xda\x3c\x6c\xbc\xbe\xab\x71\xe8\x09\xbe\xae\xc3\xe1\x55\xb2\x4e\xef\x84\xb8\x1a\xcb\xc1\x8d\x08\x76\x60\x07\x4d\xb7\x4d\xb5\xd4\x7c\xcd\xba\x6c\x6b\x83\x25\xc2\x51\x38\x9e\x3d\xe6\x52\x4b\x7b\xa4\x3b\x06\x11\x76\xd6\xc3\x4b\x80\x8f\xcc\x20\x77\xb2\x91\x03\xe6\xe9\xd9\x52\x6f\xcc\x2f\x40\xcb\xa7\xee\x13\x34\xca\xe5\xdf\xe8\xcd\x37\x00\x00\xff\xff\x76\xac\xcb\x7e\xa5\x00\x00\x00")

func sql20180128_1_groups_finishedSqlBytes() ([]byte, error) {
	return bindataRead(
		_sql20180128_1_groups_finishedSql,
		"sql/20180128_1_groups_finished.sql",
	)
}

func sql20180128_1_groups_finishedSql() (*asset, error) {
	bytes, err := sql20180128_1_groups_finishedSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "sql/20180128_1_groups_finished.sql", size: 165, mode: os.FileMode(436), modTime: time.Unix(1517135413, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _sql20180303_1_transactions_edited_nameSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x7c\xcd\x31\x0e\xc2\x20\x14\x06\xe0\x9d\x53\xfc\x5b\x07\xd3\x13\x74\x7a\xf6\x61\x1c\x9e\x60\x08\xb8\x36\xc4\xa2\x61\x28\x98\x96\x44\x8f\xef\xea\x60\x7a\x82\xaf\xef\x71\x58\xf2\x73\x8d\x2d\x21\xbc\x14\x89\xd7\x0e\x9e\x8e\xa2\x51\x1f\x9f\xa9\xad\xb1\x6c\xf1\xde\x72\x2d\x1b\x88\x19\xa3\x95\x70\x31\x48\x73\x6e\x69\x9e\x4a\x5c\x12\x6e\xe4\xc6\x33\x39\x18\xeb\x61\x82\x08\x58\x9f\x28\x88\x47\xd7\x0d\x4a\xfd\x02\x5c\xdf\x65\x9f\x60\x67\xaf\x7f\x8c\x41\x7d\x03\x00\x00\xff\xff\x74\x3c\xc1\x45\xa8\x00\x00\x00")

func sql20180303_1_transactions_edited_nameSqlBytes() ([]byte, error) {
	return bindataRead(
		_sql20180303_1_transactions_edited_nameSql,
		"sql/20180303_1_transactions_edited_name.sql",
	)
}

func sql20180303_1_transactions_edited_nameSql() (*asset, error) {
	bytes, err := sql20180303_1_transactions_edited_nameSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "sql/20180303_1_transactions_edited_name.sql", size: 168, mode: os.FileMode(436), modTime: time.Unix(1520069813, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _sql20180729_1_recurring_to_groupsSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xd2\xd5\x55\xd0\xce\xcd\x4c\x2f\x4a\x2c\x49\x55\x08\x2d\xe0\x72\xf4\x09\x71\x0d\x52\x08\x71\x74\xf2\x71\x55\xc8\x4f\xab\x88\x2f\x29\x4a\xcc\x2b\x4e\x4c\x2e\xc9\xcc\xcf\x2b\x56\x08\x72\xf5\x73\xf4\x75\x55\x28\x4a\x4d\x2e\x2d\x2a\xca\xcc\x4b\x8f\xcf\x4c\x51\x08\xf1\x57\x48\x2f\xca\x2f\x2d\x88\xcf\x4c\xb1\x46\xd1\x8d\x50\x05\x96\x87\xeb\x0e\xf1\x57\x40\x32\x14\x2a\x69\xcd\xc5\x85\xec\x10\x97\xfc\xf2\x3c\xa2\x9c\x02\xb3\x1a\x64\x2a\xb2\xb3\x50\x9d\x82\x69\x1f\x92\x63\xd0\xdd\x69\xcd\x05\x08\x00\x00\xff\xff\xd1\xca\xcf\x4b\x13\x01\x00\x00")

func sql20180729_1_recurring_to_groupsSqlBytes() ([]byte, error) {
	return bindataRead(
		_sql20180729_1_recurring_to_groupsSql,
		"sql/20180729_1_recurring_to_groups.sql",
	)
}

func sql20180729_1_recurring_to_groupsSql() (*asset, error) {
	bytes, err := sql20180729_1_recurring_to_groupsSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "sql/20180729_1_recurring_to_groups.sql", size: 275, mode: os.FileMode(436), modTime: time.Unix(1532875795, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _sql20180729_2_trend_idSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x7c\xcd\x31\x0e\xc2\x20\x14\x06\xe0\x9d\x53\xfc\x5b\x07\xd3\x13\x74\x7a\xf6\x61\x1c\x9e\x60\x08\xb8\x36\xc4\x56\xc3\x20\x18\x20\xd1\xe3\xbb\x1a\x07\x2f\xf0\x7d\xe3\x88\xdd\x23\xdd\x6b\xec\x1b\xc2\x53\x91\x78\xed\xe0\x69\x2f\x1a\xe5\xf6\x5e\x7a\x8d\xb9\xc5\x6b\x4f\x25\x37\x10\x33\x66\x2b\xe1\x64\xd0\xeb\x96\xd7\x25\xad\xb8\x90\x9b\x8f\xe4\x60\xac\x87\x09\x22\x60\x7d\xa0\x20\x1e\xc3\x30\x29\xf5\xad\x73\x79\xe5\xff\x3e\x3b\x7b\xfe\x0d\x26\xf5\x09\x00\x00\xff\xff\x91\x31\xb2\x98\xa2\x00\x00\x00")

func sql20180729_2_trend_idSqlBytes() ([]byte, error) {
	return bindataRead(
		_sql20180729_2_trend_idSql,
		"sql/20180729_2_trend_id.sql",
	)
}

func sql20180729_2_trend_idSql() (*asset, error) {
	bytes, err := sql20180729_2_trend_idSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "sql/20180729_2_trend_id.sql", size: 162, mode: os.FileMode(436), modTime: time.Unix(1532876820, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _sql20180729_3_groups_constraintSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xd2\xd5\x55\xd0\xce\xcd\x4c\x2f\x4a\x2c\x49\x55\x08\x2d\xe0\x72\xf4\x09\x71\x0d\x52\xf0\xf4\x73\x71\x8d\x50\x28\x4a\x4d\x2e\x2d\x2a\xca\xcc\x4b\x8f\x4f\x2f\xca\x2f\x2d\x28\x8e\x2f\xc8\x4e\xad\x54\x08\x72\xf5\x73\xf4\x75\x55\x08\xf1\x57\x28\x29\x4a\xcc\x2b\x4e\x4c\x2e\xc9\xcc\xcf\x43\x56\x61\xcd\xc5\x85\x6c\xa8\x4b\x7e\x79\x1e\x8a\xb1\x38\xb4\x21\x19\x8c\xd5\x62\x6b\x2e\x40\x00\x00\x00\xff\xff\x53\xa2\xa7\x94\xab\x00\x00\x00")

func sql20180729_3_groups_constraintSqlBytes() ([]byte, error) {
	return bindataRead(
		_sql20180729_3_groups_constraintSql,
		"sql/20180729_3_groups_constraint.sql",
	)
}

func sql20180729_3_groups_constraintSql() (*asset, error) {
	bytes, err := sql20180729_3_groups_constraintSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "sql/20180729_3_groups_constraint.sql", size: 171, mode: os.FileMode(436), modTime: time.Unix(1532877277, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _sql20180729_4_group_trend_idSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x84\xcd\x31\x0e\xc2\x20\x14\x06\xe0\x9d\x53\xfc\x5b\x07\xd3\x13\x74\x7a\xf6\x61\x1c\x9e\x60\x08\xb8\x36\xc4\x36\x0d\x83\xd0\x50\x8c\xd7\x77\x35\x2e\xbd\xc0\xf7\xf5\x3d\x4e\xaf\xb4\xd6\xd8\x16\x84\x4d\x91\x78\xed\xe0\xe9\x2c\x1a\xad\xc6\xbc\xc7\x67\x4b\x25\x4f\x6b\x2d\xef\x6d\x07\x31\x63\xb4\x12\x6e\x06\xad\x2e\x79\x9e\xd2\x8c\x07\xb9\xf1\x4a\x0e\xc6\x7a\x98\x20\x02\xd6\x17\x0a\xe2\xd1\x75\x83\x52\xbf\x3e\x97\x4f\x3e\x1a\xd8\xd9\xfb\x7f\x31\xa8\x6f\x00\x00\x00\xff\xff\x91\x38\xff\x6c\xa6\x00\x00\x00")

func sql20180729_4_group_trend_idSqlBytes() ([]byte, error) {
	return bindataRead(
		_sql20180729_4_group_trend_idSql,
		"sql/20180729_4_group_trend_id.sql",
	)
}

func sql20180729_4_group_trend_idSql() (*asset, error) {
	bytes, err := sql20180729_4_group_trend_idSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "sql/20180729_4_group_trend_id.sql", size: 166, mode: os.FileMode(436), modTime: time.Unix(1532878356, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _sql20180729_5_trendsSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\xce\xb1\x6a\x85\x30\x18\xc5\xf1\x39\xdf\x53\x9c\xed\xde\x4b\x15\xba\x3b\xa5\x9a\x52\x69\xaa\x12\x62\xc1\x49\xc4\xa4\x92\xc1\x44\x3e\x95\xd2\xb7\x2f\x85\x0e\x85\x2e\xf7\xcc\xe7\x0f\xbf\x3c\xc7\xc3\x1a\x16\x9e\x0e\x8f\x7e\xa3\xd2\x28\x69\x15\xac\x7c\xd2\x0a\x07\xfb\xe8\x76\x5c\x49\x04\x07\x21\x00\xbc\x4b\x53\xbe\x48\x83\xa6\xb5\x68\x7a\xad\x51\xa9\x67\xd9\x6b\x8b\xcb\x25\x23\x11\xa7\xd5\xdf\xf3\xfb\x48\x3c\x7b\x37\x6e\x9e\x43\x72\xa3\x9b\xbe\x76\xd4\x8d\xc5\xcf\xfe\x05\x8f\x19\x89\xce\xd4\x6f\xd2\x0c\x78\x55\xc3\x35\xb8\x1b\xdd\x0a\xa2\xbf\xee\x2a\x7d\x46\xaa\x4c\xdb\xfd\xba\xd9\xcf\x27\x73\x88\xcb\xb8\x70\x3a\xb7\xbd\xa0\xef\x00\x00\x00\xff\xff\x9b\xbb\x66\x41\xe6\x00\x00\x00")

func sql20180729_5_trendsSqlBytes() ([]byte, error) {
	return bindataRead(
		_sql20180729_5_trendsSql,
		"sql/20180729_5_trends.sql",
	)
}

func sql20180729_5_trendsSql() (*asset, error) {
	bytes, err := sql20180729_5_trendsSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "sql/20180729_5_trends.sql", size: 230, mode: os.FileMode(436), modTime: time.Unix(1532878552, 0)}
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
	"sql/20180124_init.sql": sql20180124_initSql,
	"sql/20180126_add_accounts.sql": sql20180126_add_accountsSql,
	"sql/20180127_sensitive_data.sql": sql20180127_sensitive_dataSql,
	"sql/20180127_z_1_null_and_default.sql": sql20180127_z_1_null_and_defaultSql,
	"sql/20180127_z_2_recurring_groups.sql": sql20180127_z_2_recurring_groupsSql,
	"sql/20180128_1_groups_finished.sql": sql20180128_1_groups_finishedSql,
	"sql/20180303_1_transactions_edited_name.sql": sql20180303_1_transactions_edited_nameSql,
	"sql/20180729_1_recurring_to_groups.sql": sql20180729_1_recurring_to_groupsSql,
	"sql/20180729_2_trend_id.sql": sql20180729_2_trend_idSql,
	"sql/20180729_3_groups_constraint.sql": sql20180729_3_groups_constraintSql,
	"sql/20180729_4_group_trend_id.sql": sql20180729_4_group_trend_idSql,
	"sql/20180729_5_trends.sql": sql20180729_5_trendsSql,
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
	"sql": &bintree{nil, map[string]*bintree{
		"20180124_init.sql": &bintree{sql20180124_initSql, map[string]*bintree{}},
		"20180126_add_accounts.sql": &bintree{sql20180126_add_accountsSql, map[string]*bintree{}},
		"20180127_sensitive_data.sql": &bintree{sql20180127_sensitive_dataSql, map[string]*bintree{}},
		"20180127_z_1_null_and_default.sql": &bintree{sql20180127_z_1_null_and_defaultSql, map[string]*bintree{}},
		"20180127_z_2_recurring_groups.sql": &bintree{sql20180127_z_2_recurring_groupsSql, map[string]*bintree{}},
		"20180128_1_groups_finished.sql": &bintree{sql20180128_1_groups_finishedSql, map[string]*bintree{}},
		"20180303_1_transactions_edited_name.sql": &bintree{sql20180303_1_transactions_edited_nameSql, map[string]*bintree{}},
		"20180729_1_recurring_to_groups.sql": &bintree{sql20180729_1_recurring_to_groupsSql, map[string]*bintree{}},
		"20180729_2_trend_id.sql": &bintree{sql20180729_2_trend_idSql, map[string]*bintree{}},
		"20180729_3_groups_constraint.sql": &bintree{sql20180729_3_groups_constraintSql, map[string]*bintree{}},
		"20180729_4_group_trend_id.sql": &bintree{sql20180729_4_group_trend_idSql, map[string]*bintree{}},
		"20180729_5_trends.sql": &bintree{sql20180729_5_trendsSql, map[string]*bintree{}},
	}},
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

