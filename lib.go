package goupdater

import (
	"time"

	"github.com/goupdate/goupdater/lib"
)

type updater struct {
	lib.ClientInfo
	timeout time.Duration
}

func New(serverIpPort, key, project, branch string) *updater {
	return &updater{
		ClientInfo: lib.ClientInfo{
			Server:  serverIpPort,
			Key:     key,
			Project: project,
			Branch:  branch,
		},
		timeout: 2*time.Minute,
	}
}

func (u *updater) Check() (bool, error) {
	return u.ClientInfo.Check()
}

func (u *updater) DownloadAndReplaceMe() error {
	return u.ClientInfo.DownloadAndReplaceMe(u.timeout)
}

func (u *updater) Upload(branch string) error {
	return u.ClientInfo.Upload(branch)
}

//set verbose to b
func (u *updater) Verbose(b bool) {
	u.ClientInfo.Verbose = b
}

func (u *updater) SetDownloadTimeout(t time.Duration) {
	u.timeout = t
}
