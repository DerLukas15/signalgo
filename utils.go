package signalgo

import (
	"errors"
	"os"
	"syscall"
)

func commandValid(cmd string) (bool, error) {
	infos, err := os.Stat(cmd)
	if err != nil {
		return false, err
	}
	thisUID := os.Getuid()
	thisGID := os.Getgid()
	var fileUID int
	var fileGID int
	perm := infos.Mode().Perm()
	if stat, ok := infos.Sys().(*syscall.Stat_t); ok {
		fileUID = int(stat.Uid)
		fileGID = int(stat.Gid)
	} else {
		return false, errors.New("cannot determine permissions for '" + cmd + "'")
	}
	if fileUID == thisUID {
		valid := perm&0100 != 0
		if valid {
			return true, nil
		}
	}
	if fileGID == thisGID {
		valid := perm&0010 != 0
		if valid {
			return true, nil
		}
	}
	return perm&0001 != 0, nil
}
