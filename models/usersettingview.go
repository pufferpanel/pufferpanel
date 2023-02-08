package models

type UserSettingView struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type UserSettingsView []*UserSettingView

func FromUserSetting(setting *UserSetting) *UserSettingView {
	return &UserSettingView{
		Key:   setting.Key,
		Value: setting.Value,
	}
}

func FromUserSettings(settings []*UserSetting) UserSettingsView {
	result := make(UserSettingsView, len(settings))

	for k, v := range settings {
		result[k] = FromUserSetting(v)
	}

	return result
}
