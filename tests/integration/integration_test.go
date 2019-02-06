package integration

import (
	"crypto/sha1"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/imroc/req"
	"github.com/jwoglom/goscout/app/endpointsv1"
	"github.com/ttacon/glog"
)

const GOSCPort = 3000
const GOSC = "http://localhost:3000/"
const NSSCPort = 1337
const NSSC = "http://localhost:1337/"
const DefaultAPISecret = "defaultapisecret"

func init() {
	if _, err := os.Stat("cgm-remote-monitor"); os.IsNotExist(err) {
		fmt.Printf("Cloning Nightscout...\n")
		cloneNightscout()
	}

	if _, err := os.Stat("cgm-remote-monitor/data"); os.IsExist(err) {
		fmt.Printf("Deleting Nightscout data...\n")
		deleteNightscoutData()
	}

	fmt.Printf("Starting Nightscout and Goscout\n")
	go startNightscout()
	go startGoscout()
	fmt.Printf("Waiting for Goscout to start")
	waitPortAlive(GOSCPort)
	fmt.Printf("Waiting for Goscout to initialize")
	waitHTTPAlive(GOSCPort)
	fmt.Printf("Waiting for Nightscout to start")
	waitPortAlive(NSSCPort)
	fmt.Printf("Waiting for Nightscout to initialize")
	waitHTTPAlive(NSSCPort)

	fmt.Printf("Starting tests...\n")
}

func waitPortAlive(port int) {
	for {
		alive := checkPortAlive(port)
		if alive {
			fmt.Printf(" Port %d is alive.\n", port)
			break
		}
		fmt.Printf(".")
		time.Sleep(2 * time.Second)
	}
}

func waitHTTPAlive(port int) {
	for {
		alive := checkHTTPAlive(port)
		if alive {
			fmt.Printf(" Port %d (http) is alive.\n", port)
			break
		}
		fmt.Printf(".")
		time.Sleep(2 * time.Second)
	}
}

func checkPortAlive(port int) bool {
	timeout := 2 * time.Second
	conn, err := net.DialTimeout("tcp", fmt.Sprintf(":%d", port), timeout)
	if err == nil {
		conn.Close()
		return true
	}
	return false
}

func checkHTTPAlive(port int) bool {
	_, err := req.Get(fmt.Sprintf("http://localhost:%d/", port))
	return err == nil
}

// ExecCmd executes a command and returns the Cmd object which can be .Run()
func ExecCmd(arg string, args ...string) *exec.Cmd {
	cmd := exec.Command(arg, args...)
	fmt.Printf("ExecCmd: %s %s\n", arg, strings.Join(args, " "))
	//cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}

// ExecCmdDir executes a command with the given current working directory
func ExecCmdDir(dir, arg string, args ...string) *exec.Cmd {
	cmd := ExecCmd(arg, args...)
	cmd.Dir = dir
	return cmd
}

func cloneNightscout() {
	glog.FatalIf(ExecCmd("git", "clone", "https://github.com/jwoglom/cgm-remote-monitor").Run())
}

func deleteNightscoutData() {
	glog.FatalIf(ExecCmd("rm", "-rf", "cgm-remote-monitor/data").Run())
}

func startNightscout() {
	glog.FatalIf(ExecCmdDir("cgm-remote-monitor", "docker-compose", "up").Run())
}

func startGoscout() {
	//glog.FatalIf(ExecCmdDir("../..", "make", "build").Run())
	glog.FatalIf(ExecCmdDir("../..", "make", "run").Run())
}

func assert(s string, b bool, t *testing.T) {
	t.Helper()
	if !b {
		t.Errorf("Assertion failed: %s", s)
	} else {
		fmt.Printf("OK: %s\n", s)
	}
}

func reqOK(resp *req.Resp, err error, t *testing.T) {
	t.Helper()
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	if status := resp.Response().StatusCode; status != 200 {
		t.Errorf("Unexpected status code: %d", status)
	}
}

func TestNsscApiAccessible(t *testing.T) {
	// gosc := r.Get(GOSC + "api/v1/status")
	nssc, err := req.Get(NSSC + "api/v1/status")
	reqOK(nssc, err, t)
	assert("body status ok", nssc.String() == "<h1>STATUS OK</h1>", t)
}

func TestGoscApiAccessible(t *testing.T) {
	// gosc := r.Get(GOSC + "api/v1/status")
	gosc, err := req.Get(GOSC + "api/v1/status")
	reqOK(gosc, err, t)
	assert("body status ok", gosc.String() == "<h1>STATUS OK</h1>", t)
}

func TestCompareStatusJSON(t *testing.T) {
	nssc, err := req.Get(NSSC + "api/v1/status.json")
	reqOK(nssc, err, t)
	gosc, err := req.Get(GOSC + "api/v1/status.json")
	reqOK(gosc, err, t)

	var ns endpointsv1.Status
	nssc.ToJSON(&ns)
	var gs endpointsv1.Status
	gosc.ToJSON(&gs)
	assert("status ok", ns.Status == gs.Status, t)
	assert("name is nightscout", ns.Name == "Nightscout", t)
	assert("name is goscout", gs.Name == "Goscout", t)
}

func getNightscoutToken() req.Header {
	h := sha1.New()
	h.Write([]byte(DefaultAPISecret))
	tok := fmt.Sprintf("%x", h.Sum(nil))
	fmt.Printf("token: %s\n", tok)
	return req.Header{
		"Api-Secret": tok,
	}
}

func helperUploadEntry(entry []req.Param, t *testing.T) (endpointsv1.Entries, endpointsv1.Entries) {
	nsToken := getNightscoutToken()

	nssc, err := req.Post(NSSC+"api/v1/entries.json", req.BodyJSON(&entry), nsToken)
	reqOK(nssc, err, t)

	gosc, err := req.Post(GOSC+"api/v1/entries.json", req.BodyJSON(&entry))
	reqOK(gosc, err, t)

	count := fmt.Sprintf("%d", len(entry))
	nssc, err = req.Get(NSSC + "api/v1/entries.json?count=" + count)
	reqOK(nssc, err, t)

	gosc, err = req.Get(GOSC + "api/v1/entries.json?count=" + count)
	reqOK(gosc, err, t)

	var ns endpointsv1.Entries
	nssc.ToJSON(&ns)

	var gs endpointsv1.Entries
	gosc.ToJSON(&gs)

	return ns, gs
}

func helperAssertEntry(a, b endpointsv1.Entry, t *testing.T) {
	t.Logf("helperAssertEntry: %+v %+v\n", a, b)
	assert("entry device", a.Device == b.Device, t)
	assert("entry date int", a.Date == b.Date, t)
	assert("entry datestring", a.DateString == b.DateString, t)
	assert("entry sgv", a.Sgv == b.Sgv, t)
	assert("entry delta", a.Delta == b.Delta, t)
	assert("entry direction", a.Direction == b.Direction, t)
	assert("entry type", a.Type == b.Type, t)
}

func TestUploadEntry(t *testing.T) {
	entry := []req.Param{{
		"device":     "xDrip-DexcomG5 G5 Native",
		"date":       1549339760137,
		"dateString": "2019-02-04T23:09:20.137-0500",
		"sgv":        172,
		"delta":      1.999,
		"direction":  "Flat",
		"type":       "sgv",
		"filtered":   186184,
		"unfiltered": 180982,
		"rssi":       100,
		"noise":      1,
		"sysTime":    "2019-02-04T23:09:20.137-0500",
	}}

	ns, gs := helperUploadEntry(entry, t)

	assert("one entry in nightscout", len(ns) == 1, t)
	assert("one entry in goscout", len(gs) == 1, t)
	if len(ns) != 1 || len(gs) != 1 {
		t.Fatal("can't compare nonexistent entry")
	}

	helperAssertEntry(ns[0], gs[0], t)

}
