package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hybridgroup/gobot"
)

func initTestAPI() *api {
	log.SetOutput(gobot.NullReadWriteCloser{})
	g := gobot.NewGobot()
	a := NewAPI(g)
	a.start = func(m *api) {}
	a.Start()

	g.AddRobot(gobot.NewTestRobot("Robot 1"))
	g.AddRobot(gobot.NewTestRobot("Robot 2"))
	g.AddRobot(gobot.NewTestRobot("Robot 3"))

	return a
}

func TestRobots(t *testing.T) {
	a := initTestAPI()
	request, _ := http.NewRequest("GET", "/robots", nil)
	response := httptest.NewRecorder()
	a.server.ServeHTTP(response, request)

	body, _ := ioutil.ReadAll(response.Body)
	var i []map[string]interface{}
	json.Unmarshal(body, &i)
	gobot.Expect(t, len(i), 3)
}

func TestRobot(t *testing.T) {
	a := initTestAPI()
	request, _ := http.NewRequest("GET", "/robots/Robot%201", nil)
	response := httptest.NewRecorder()
	a.server.ServeHTTP(response, request)

	body, _ := ioutil.ReadAll(response.Body)
	var i map[string]interface{}
	json.Unmarshal(body, &i)
	gobot.Expect(t, i["name"].(string), "Robot 1")
}

func TestRobotDevices(t *testing.T) {
	a := initTestAPI()
	request, _ := http.NewRequest("GET", "/robots/Robot%201/devices", nil)
	response := httptest.NewRecorder()
	a.server.ServeHTTP(response, request)

	body, _ := ioutil.ReadAll(response.Body)
	var i []map[string]interface{}
	json.Unmarshal(body, &i)
	gobot.Expect(t, len(i), 3)
}

func TestRobotCommands(t *testing.T) {
	a := initTestAPI()
	request, _ := http.NewRequest("GET", "/robots/Robot%201/commands", nil)
	response := httptest.NewRecorder()
	a.server.ServeHTTP(response, request)

	body, _ := ioutil.ReadAll(response.Body)
	var i []string
	json.Unmarshal(body, &i)
	gobot.Expect(t, i, []string{"robotTestFunction"})
}

func TestExecuteRobotCommand(t *testing.T) {
	a := initTestAPI()
	request, _ := http.NewRequest("GET", "/robots/Robot%201/commands/robotTestFunction", bytes.NewBufferString(`{"message":"Beep Boop"}`))
	request.Header.Add("Content-Type", "application/json")
	response := httptest.NewRecorder()
	a.server.ServeHTTP(response, request)

	body, _ := ioutil.ReadAll(response.Body)
	var i interface{}
	json.Unmarshal(body, &i)
	gobot.Expect(t, i, "hey Robot 1, Beep Boop")
}

func TestUnknownRobotCommand(t *testing.T) {
	a := initTestAPI()
	request, _ := http.NewRequest("GET", "/robots/Robot%201/commands/robotTestFuntion1", bytes.NewBufferString(`{"message":"Beep Boop"}`))
	request.Header.Add("Content-Type", "application/json")
	response := httptest.NewRecorder()
	a.server.ServeHTTP(response, request)

	body, _ := ioutil.ReadAll(response.Body)
	var i interface{}
	json.Unmarshal(body, &i)
	gobot.Expect(t, i, "Unknown Command")
}

func TestRobotDevice(t *testing.T) {
	a := initTestAPI()
	request, _ := http.NewRequest("GET", "/robots/Robot%201/devices/Device%201", nil)
	response := httptest.NewRecorder()
	a.server.ServeHTTP(response, request)

	body, _ := ioutil.ReadAll(response.Body)
	var i map[string]interface{}
	json.Unmarshal(body, &i)
	gobot.Expect(t, i["name"].(string), "Device 1")
}

func TestRobotDeviceCommands(t *testing.T) {
	a := initTestAPI()
	request, _ := http.NewRequest("GET", "/robots/Robot%201/devices/Device%201/commands", nil)
	response := httptest.NewRecorder()
	a.server.ServeHTTP(response, request)

	body, _ := ioutil.ReadAll(response.Body)
	var i []string
	json.Unmarshal(body, &i)
	gobot.Expect(t, i, []string{"TestDriverCommand", "DriverCommand"})
}

func TestExecuteRobotDeviceCommand(t *testing.T) {
	a := initTestAPI()
	request, _ := http.NewRequest("GET", "/robots/Robot%201/devices/Device%201/commands/TestDriverCommand", bytes.NewBufferString(`{"name":"human"}`))
	request.Header.Add("Content-Type", "application/json")
	response := httptest.NewRecorder()
	a.server.ServeHTTP(response, request)

	body, _ := ioutil.ReadAll(response.Body)
	var i interface{}
	json.Unmarshal(body, &i)
	gobot.Expect(t, i, "hello human")
}

func TestUnknownRobotDeviceCommand(t *testing.T) {
	a := initTestAPI()
	request, _ := http.NewRequest("GET", "/robots/Robot%201/devices/Device%201/commands/DriverCommand1", bytes.NewBufferString(`{"name":"human"}`))
	request.Header.Add("Content-Type", "application/json")
	response := httptest.NewRecorder()
	a.server.ServeHTTP(response, request)

	body, _ := ioutil.ReadAll(response.Body)
	var i interface{}
	json.Unmarshal(body, &i)
	gobot.Expect(t, i, "Unknown Command")
}
