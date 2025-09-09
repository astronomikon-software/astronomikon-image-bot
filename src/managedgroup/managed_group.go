package managedgroup

import (
	"errors"
	"io"
	"ortemios/imgbot/types"
	"ortemios/imgbot/util"
	"os"
)

const unset = ""

var groupID = unset

var managedGroupFile string

func init() {
	managedGroupFile = util.MustGetEnv("MANAGED_GROUP_FILE")
}

var ErrGroupUnset = errors.New("group is not set")

func Set(id string) error {
	groupID = id
	return writeToFile(id)
}

func Get() (string, error) {
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

func writeToFile(id string) error {
	return os.WriteFile(managedGroupFile, []byte(id), 0644)
}

func readFromFile() (string, error) {
	file, err := os.Open(managedGroupFile)
	if err != nil {
		return unset, err
	}

	text, err := io.ReadAll(file)
	if err != nil {
		return unset, err
	}

	return types.GroupID(text), nil
}
