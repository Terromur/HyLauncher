package app

import (
	"HyLauncher/internal/config"
	"HyLauncher/pkg/hyerrors"
)

func (a *App) SetNick(nick string) error {
	if nick == "" {
		err := hyerrors.Validation("nickname cannot be empty")
		hyerrors.Report(err)
		return err
	}

	a.cfg.Nick = nick

	if err := config.Save(a.cfg); err != nil {
		appErr := hyerrors.WrapConfig(err, "failed to save nickname").
			WithContext("nick", nick)
		hyerrors.Report(appErr)
		return appErr
	}

	return nil
}

func (a *App) GetNick() string {
	return a.cfg.Nick
}

func (a *App) GetLauncherVersion() string {
	return config.Default().Version
}

func (a *App) SetLocalGameVersion(version int) error {
	if err := config.SaveLocalGameVersion(version); err != nil {
		appErr := hyerrors.WrapConfig(err, "failed to save game version").
			WithContext("version", version)
		hyerrors.Report(appErr)
		return appErr
	}
	return nil
}

func (a *App) GetLocalGameVersion() int {
	return a.cfg.CurrentGameVersion
}

func (a *App) SetBranch(branch string) error {
	if err := config.SaveBranch(branch); err != nil {
		appErr := hyerrors.WrapConfig(err, "failed to save branch").
			WithContext("branch", branch)
		hyerrors.Report(appErr)
		return appErr
	}
	return nil
}

func (a *App) GetBranch() (string, error) {
	branch, err := config.GetBranch()
	if err != nil {
		appErr := hyerrors.WrapConfig(err, "failed to get branch")
		hyerrors.Report(appErr)
		return "", appErr
	}
	return branch, nil
}
