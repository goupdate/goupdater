package goupdater

import (
	"goupdater/lib"
)

type updater struct {
	lib.ClientInfo
}

func New(serverIpPort, key, project, branch string) *updater {
	return &updater{
		ClientInfo: lib.ClientInfo{
			Server:  serverIpPort,
			Key:     key,
			Project: project,
			Branch:  branch,
		},
	}
}

func (u *updater) Check() (bool, error) {
	return u.ClientInfo.Check()
}

func (u *updater) DownloadAndReplaceMe() error {
	return u.ClientInfo.DownloadAndReplaceMe()
}

func (u *updater) Upload(branch string) error {
	return u.ClientInfo.Upload(branch)
}

//set verbose to b
func (u *updater) Verbose(b bool) {
	u.ClientInfo.Verbose = b
}
