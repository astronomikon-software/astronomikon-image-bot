package managedgroup

import (
	"errors"
	"io"
	"ortemios/imgbot/types"
	"ortemios/imgbot/util"
	"os"
	"strconv"
)

const unset = types.ChatID(-1)

var groupID = unset

var managedGroupFile string

func init() {
	managedGroupFile = util.MustGetEnv("MANAGED_GROUP_FILE")
}

var ErrGroupUnset = errors.New("group is not set")

func Set(id types.ChatID) error {
	groupID = id
	return writeToFile(id)
}

func Get() (types.ChatID, error) {
	if groupID != unset {
		return groupID, nil
	} else {
		gid, err := readFromFile()
		if err == nil {
			groupID = gid
			return groupID, nil
		} else {
			return unset, ErrGroupUnset
		}
	}
}

func writeToFile(id types.ChatID) error {
	return os.WriteFile(managedGroupFile, []byte(strconv.Itoa(int(id))), 0644)
}

func readFromFile() (types.ChatID, error) {
	file, err := os.Open(managedGroupFile)
	if err != nil {
		return unset, err
	}

	text, err := io.ReadAll(file)
	if err != nil {
		return unset, err
	}

	gid, err := strconv.Atoi(string(text))
	if err != nil {
		return unset, err
	}

	return types.ChatID(gid), nil
}
