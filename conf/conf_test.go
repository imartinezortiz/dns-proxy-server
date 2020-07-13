package conf

import (
	"testing"
	"os"
	"github.com/stretchr/testify/assert"
	"github.com/mageddo/dns-proxy-server/flags"
	"github.com/mageddo/dns-proxy-server/events/local"
	"github.com/mageddo/dns-proxy-server/utils"
	"flag"
	"github.com/mageddo/dns-proxy-server/utils/env"
	"github.com/mageddo/go-logging"
)

func TestDefaultFlagValues(t *testing.T) {
	assert.Equal(t, *flags.WebServerPort, WebServerPort())
	assert.Equal(t, *flags.DnsServerPort, DnsServerPort())
	assert.Equal(t, *flags.SetupResolvconf, SetupResolvConf())
}

func TestFlagValuesFromArgs(t *testing.T) {
	os.Args = []string{"cmd", "-web-server-port=8282", "--server-port=61", "-default-dns=false"}
	flag.Parse()
	assert.Equal(t, 8282, WebServerPort())
	assert.Equal(t, 61, DnsServerPort())
	assert.Equal(t, false, SetupResolvConf())
}

func TestFlagValuesFromConf(t *testing.T) {

	// arrange
	local.ResetConf()

	// act
	err := utils.WriteToFile(`{ "webServerPort": 8080, "dnsServerPort": 62, "defaultDns": false }`, utils.SolveRelativePath(*flags.ConfPath))

	// assert
	assert.Nil(t, err)
	assert.Equal(t, 8080, WebServerPort())
	assert.Equal(t, 62, DnsServerPort())
	assert.Equal(t, false, SetupResolvConf())
}


func TestLogLevel_DefaultValue(t *testing.T) {
	assert.Equal(t, logging.INFO, LogLevel())
}

func TestLogLevel_ReadFromConfig(t *testing.T) {

	// arrange
	local.ResetConf()
	c, err := local.LoadConfiguration()
	assert.Nil(t, err)
	c.LogLevel = "INFO"
	local.SaveConfiguration(c)

	// act
	level := LogLevel()

	// assert
	assert.Equal(t, logging.INFO, level)

	os.Remove(local.GetConfPath())
}

func TestLogLevel_ReadFromEnv(t *testing.T) {

	// arrange
	os.Setenv(env.MG_LOG_LEVEL, "WARNING")

	// act
	level := LogLevel()

	// assert
	assert.Equal(t, logging.WARNING, level)

}

func TestLogFile_DefaultValue(t *testing.T) {

	// arrange

	// act
	level := LogFile()

	// assert
	assert.Equal(t,"console", level)
}

func TestLogFile_ReadFromConfig(t *testing.T) {

	// arrange
	c, err := local.LoadConfiguration()
	assert.Nil(t, err)
	c.LogFile = "false"
	local.SaveConfiguration(c)

	// act
	level := LogFile()

	// assert
	assert.Equal(t,"", level)

	os.Remove(local.GetConfPath())
}

func TestLogFile_ReadFromEnv(t *testing.T) {

	// arrange
	os.Setenv(env.MG_LOG_FILE, "true")

	// act
	level := LogFile()

	// assert
	assert.Equal(t,"/var/log/dns-proxy-server.log", level)

	os.Remove(local.GetConfPath())
}

func TestLogFile_CustomPath(t *testing.T) {

	// arrange
	os.Setenv(env.MG_LOG_FILE, "custom-file.log")

	// act
	level := LogFile()

	// assert
	assert.Equal(t,"custom-file.log", level)

	os.Remove(local.GetConfPath())
}


const TEST_ENV_VAR_NAME = "DPS_TEST_ENV_VAR"

func TestParseBooleanEnv_TrueValues(t *testing.T) {
	trueValues := []string{"1", "t", "T", "TRUE", "true", "True"}
	for _, s := range trueValues {
		// arrange
		os.Setenv(TEST_ENV_VAR_NAME, s)

		// act
		declared, value, err := parseBooleanEnv(TEST_ENV_VAR_NAME)

		// assert
		assert.True(t, declared)
		assert.True(t, value, "Environment variable value: "+s)
		assert.NoError(t, err)
	}
}

func TestParseBooleanEnv_FalseValues(t *testing.T) {

	falseValues := []string{"0", "f", "F", "FALSE", "false", "False"}
	for _, s := range falseValues {
		// arrange
		os.Setenv(TEST_ENV_VAR_NAME, s)

		// act
		declared, value, err := parseBooleanEnv(TEST_ENV_VAR_NAME)

		// assert
		assert.True(t, declared)
		assert.False(t, value)
		assert.NoError(t, err)
	}
}

func TestParseBooleanEnv_Error(t *testing.T) {

	falseValues := []string{"x", "y"}
	for _, s := range falseValues {
		// arrange
		os.Setenv(TEST_ENV_VAR_NAME, s)

		// act
		declared, _, err := parseBooleanEnv(TEST_ENV_VAR_NAME)

		// assert
		assert.True(t, declared)
		assert.Error(t, err)
	}
}


func TestParseBooleanEnv_NotDeclared(t *testing.T) {

	// arrange

	// act
	declared, _, err := parseBooleanEnv("DPS_NOT_EXISTING_ENV_VAR")

	// assert
	assert.False(t, declared)
	assert.NoError(t, err)
}