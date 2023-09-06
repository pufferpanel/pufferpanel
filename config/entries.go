package config

import (
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"path/filepath"
)

var LogsFolder = asString("logs", "logs")
var WebHost = asString("web.host", "0.0.0.0:8080")

var PanelEnabled = asBool("panel.enable", true)
var DatabaseDialect = asString("panel.database.dialect", "sqlite3")
var DatabaseUrl = asString("panel.database.url", "")
var DatabaseLoggingEnabled = asBool("panel.database.log", false)
var WebRoot = asString("panel.web.files", "www")

var EmailTemplateFolder = asString("panel.email.templateFolder", "")
var EmailProvider = asString("panel.email.provider", "")
var EmailFrom = asString("panel.email.from", "")
var EmailDomain = asString("panel.email.domain", "")
var EmailHost = asString("panel.email.host", "")
var EmailKey = asString("panel.email.key", "")
var EmailUsername = asString("panel.email.username", "")
var EmailPassword = asString("panel.email.password", "")
var CompanyName = asString("panel.settings.companyName", "PufferPanel")
var DefaultTheme = asString("panel.settings.defaultTheme", "PufferPanel")
var ThemeSettings = asString("panel.settings.themeSettings", "{}")
var MasterUrl = asString("panel.settings.masterUrl", "http://localhost:8080")
var SessionKey = asString("panel.sessionKey", "")
var RegistrationEnabled = asBool("panel.registrationEnabled", true)

var DaemonEnabled = asBool("daemon.enable", true)
var ConsoleBuffer = asInt("daemon.console.buffer", 50)
var ConsoleForward = asBool("daemon.console.forward", false)
var SftpHost = asString("daemon.sftp.host", "0.0.0.0:5657")
var SftpKey = asString("daemon.sftp.key", "sftp.key")
var AuthUrl = asString("daemon.auth.url", "http://localhost:8080")
var ClientId = asString("daemon.auth.clientId", "")
var ClientSecret = asString("daemon.auth.clientSecret", "")
var CacheFolder = asDataFolder("daemon.data.cache", "cache")
var ServersFolder = asDataFolder("daemon.data.servers", "servers")
var BinariesFolder = asDataFolder("daemon.data.binaries", "binaries")
var CrashLimit = asInt("daemon.data.crashLimit", 3)
var WebSocketFileLimit = asInt64("daemon.data.maxWSDownloadSize", 1024*1024*20)
var CurseForgeKey = asString("daemon.curseforge.key", _hiddenCurseforgeKey)
var DataRootFolder = asString("daemon.data.root", "")

var _hiddenCurseforgeKey = ""

type entry[T ValueType] struct {
	key string
}

type StringEntry struct {
	entry[string]
}
type BoolEntry struct {
	entry[bool]
}
type IntEntry struct {
	entry[int]
}
type Int64Entry struct {
	entry[int64]
}
type DataFolder struct {
	StringEntry
}

type ValueType interface {
	int | int64 | bool | string
}

func (se StringEntry) Value() string {
	return cast.ToString(se.get())
}

func (se DataFolder) Value() string {
	p := cast.ToString(se.get())
	if filepath.IsAbs(p) {
		return p
	}

	return filepath.Join(DataRootFolder.Value(), p)
}

func (se BoolEntry) Value() bool {
	return cast.ToBool(se.get())
}
func (se IntEntry) Value() int {
	return cast.ToInt(se.get())
}
func (se Int64Entry) Value() int64 {
	return cast.ToInt64(se.get())
}

func (se entry[T]) Key() string {
	return se.key
}

func (se entry[T]) Set(value T, save bool) error {
	viper.Set(se.Key(), value)

	if save {
		return viper.WriteConfig()
	}
	return nil
}

func asDataFolder(key string, def string) DataFolder {
	return DataFolder{StringEntry{entry: as[string](key, def)}}
}
func asString(key string, def string) StringEntry {
	return StringEntry{entry: as[string](key, def)}
}
func asBool(key string, def bool) BoolEntry {
	return BoolEntry{entry: as[bool](key, def)}
}
func asInt(key string, def int) IntEntry {
	return IntEntry{entry: as[int](key, def)}
}
func asInt64(key string, def int64) Int64Entry {
	return Int64Entry{entry: as[int64](key, def)}
}

func as[T ValueType](key string, def T) entry[T] {
	//We are not using viper for this, because it writes a giant config with the defaults,
	//and we cannot do that to allow for changes in the future

	//viper.SetDefault(key, def)
	defaults[key] = def
	return entry[T]{key: key}
}

func (se entry[T]) get() interface{} {
	val := viper.Get(se.Key())
	if val == nil {
		val = defaults[se.Key()]
	}

	return val
}

var defaults = make(map[string]interface{})
