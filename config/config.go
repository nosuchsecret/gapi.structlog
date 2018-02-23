package config

import (
	"os"
	"fmt"
	//"time"
	"path/filepath"
	goconf "github.com/msbranco/goconfig"
	"github.com/nosuchsecret/gapi/errors"
	"github.com/nosuchsecret/gapi/variable"
)

// Config of server
type Config struct {
	HttpAddr    string  /* http server bind address */
	UdpAddr     string  /* udp server bind address */
	TcpAddr     string  /* tcp server bind address */
	UsocketAddr string  /* usocket server bind address */

	//UdpNFI     string  /* udp server multicast receive interface */

	Location    string  /* handler location */

	Log         string  /* log file */
	Level       string  /* log level */
	RotateSize  int     /* log rotate line */
	BackupSize  int     /* log backup size */

	File        string  /* config file */
	C           *goconf.ConfigFile /* goconfig struct */
}

func (conf *Config) SetConf(file string) {
	conf.File = filepath.Join(variable.DEFAULT_CONFIG_PATH, file)
}
// ReadConf reads conf from file
func (conf *Config) ReadConf(file string) error {
	if file == "" {
		if conf.File == "" {
			file = filepath.Join(variable.DEFAULT_CONFIG_PATH, variable.DEFAULT_CONFIG_FILE)
		} else {
			file = conf.File
		}
	}

	c, err := goconf.ReadConfigFile(file)
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Error] Read conf file %s failed", file)
		return err
	}
	conf.C = c
	return nil
}

// ParseConf parses config
func (conf *Config) ParseConf() error {
	var err error

	if conf.C == nil {
		fmt.Fprintln(os.Stderr, "[Error] Must read config first")
		return errors.BadConfigError
	}

	conf.HttpAddr, err = conf.C.GetString("default", "http_addr")
	if err != nil {
		//fmt.Fprintln(os.Stderr, "[Info] [Default] Read conf: No HttpAddr")
		conf.HttpAddr = ""
	} else {
		fmt.Fprintln(os.Stderr, "[Info] [Default] listen on http addr:", conf.HttpAddr)
	}

	conf.TcpAddr, err = conf.C.GetString("default", "tcp_addr")
	if err != nil {
		//fmt.Fprintln(os.Stderr, "[Info] [Default] Read conf: No TcpAddr")
		conf.UdpAddr = ""
	} else {
		fmt.Fprintln(os.Stderr, "[Info] [Default] listen on tcp addr:", conf.TcpAddr)
	}
	conf.UdpAddr, err = conf.C.GetString("default", "udp_addr")
	if err != nil {
		//fmt.Fprintln(os.Stderr, "[Info] [Default] Read conf: No UdpAddr")
		conf.UdpAddr = ""
	} else {
		fmt.Fprintln(os.Stderr, "[Info] [Default] listen on udp addr:", conf.UdpAddr)
	}

	conf.UsocketAddr, err = conf.C.GetString("default", "usocket_addr")
	if err != nil {
		conf.UsocketAddr = ""
	} else {
		fmt.Fprintln(os.Stderr, "[Info] [Default] listen on usocket addr:", conf.UsocketAddr)
	}
	//conf.UdpAddr, err = conf.C.GetString("default", "udp_interface")
	//if err != nil {
	//	conf.UdpNIF = ""
	//} else {
	//	fmt.Fprintln(os.Stderr, "[Info] [Default] use udp network interface:", conf.UdpNIF)
	//}

	conf.Log, err = conf.C.GetString("default", "log")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Info] [Default] log not found, use default log file")
		conf.Log = ""
	}
	conf.Level, err = conf.C.GetString("default", "level")
	if err != nil {
		conf.Level = "error"
		fmt.Fprintln(os.Stderr, "[Info] [Default] level not found, use default log level error")
	}
	rsize, err := conf.C.GetInt64("default", "rotate_size")
	if err != nil {
		rsize = variable.DEFAULT_ROTATE_SIZE
		fmt.Fprintln(os.Stderr, "[Info] [Default] rotate_size not found, use default", rsize)
	}
	conf.RotateSize = int(rsize)

	bksize, err := conf.C.GetInt64("default", "backup_size")
	if err != nil {
		bksize = variable.DEFAULT_BACKUP_SIZE
		fmt.Fprintln(os.Stderr, "[Info] [Default] backup_size not found, use default", bksize)
	}
	conf.BackupSize = int(bksize)

	return nil
}

