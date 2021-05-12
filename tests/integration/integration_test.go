package integration

import (
	"crypto/sha1"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"syscall"
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

var ConfigLogStdout = os.Getenv("LOG_PROCESS_STDOUT") == "true"
var ConfigKeepRunning = os.Getenv("KEEP_RUNNING") == "true"

var runningPids []int

func init() {
	fmt.Printf("Checking environment...")
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
	fmt.Printf("\nWaiting for Nightscout to start")
	waitPortAlive(NSSCPort)
	fmt.Printf("\nWaiting for Nightscout to initialize")
	waitHTTPAlive(NSSCPort)
	fmt.Printf("Checking Nightscout configuration: ")
	configureNightscout()

	fmt.Printf("Nightscout token: %s (%s)\n", getNightscoutToken(), DefaultAPISecret)
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
	// Binding stdin ensures the process terminates properly
	cmd.Stdin = os.Stdin
	// Binding stdout/stderr prints status info to the terminal
	if ConfigLogStdout {
		cmd.Stdout = os.Stderr
	}
	cmd.Stderr = os.Stderr
	return cmd
}

// ExecCmdDir executes a command with the given current working directory
func ExecCmdDir(dir, arg string, args ...string) *exec.Cmd {
	cmd := ExecCmd(arg, args...)
	cmd.Dir = dir
	return cmd
}

func RunUntilStopped(cmd *exec.Cmd) error {
	glog.FatalIf(cmd.Start())
	runningPids = append(runningPids, cmd.Process.Pid)
	return cmd.Wait()
}

func cloneNightscout() {
	glog.FatalIf(ExecCmd("git", "clone", "https://github.com/jwoglom/cgm-remote-monitor").Run())
}

func deleteNightscoutData() {
	glog.FatalIf(ExecCmd("rm", "-rf", "cgm-remote-monitor/data").Run())
}

func startNightscout() {
	glog.FatalIf(RunUntilStopped(ExecCmdDir("cgm-remote-monitor", "docker-compose", "up")))
}

func startGoscout() {
	//glog.FatalIf(ExecCmdDir("../..", "make", "build").Run())
	glog.FatalIf(RunUntilStopped(ExecCmdDir("../..", "make", "run")))

}

func assert(s string, b bool, t *testing.T) {
	t.Helper()
	if !b {
		t.Errorf("Assertion failed: %s\n", s)
	} else {
		fmt.Printf("OK: %s\n", s)
	}
}

func assertEq(s string, l, r interface{}, t *testing.T) {
	t.Helper()
	if l != r {
		t.Errorf("Assertion failed: %s (%v != %v)\n", s, l, r)
	} else {
		fmt.Printf("OK: %s\n", s)
	}
}

func parseTimeNS(s string) int64 {
	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		glog.Errorf("%v", err)
	}
	return t.UnixNano()
}

func parseTimeGS(s string) int64 {
	t, err := time.Parse("2006-01-02T15:04:05.999999999-0700", s)
	if err != nil {
		glog.Errorf("%v", err)
	}
	return t.UnixNano()
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

func hasNightscoutProfile() bool {
	nsToken := getNightscoutToken()

	resp, err := req.Get(NSSC+"api/v1/profile.json", nsToken)
	glog.FatalIf(err)

	var profiles endpointsv1.Profiles
	resp.ToJSON(&profiles)

	return len(profiles) >= 1
}

func configureNightscout() {
	nsToken := getNightscoutToken()
	if hasNightscoutProfile() {
		fmt.Printf("profile(s) already exist\n")
		return
	}
	fmt.Printf("no Nightscout profile exists, adding default\n")

	profile := []req.Param{{
		"defaultProfile": "Default",
		"store": req.Param{
			"Default": req.Param{
				"dia": "3",
				"carbratio": []req.Param{{
					"time":          "00:00",
					"value":         "30",
					"timeAsSeconds": "0",
				}},
				"carbs_hr": "20",
				"delay":    "20",
				"sens": []req.Param{{
					"time":          "00:00",
					"value":         "100",
					"timeAsSeconds": "0",
				}},
				"timezone": "UTC",
				"basal": []req.Param{{
					"time":          "00:00",
					"value":         "0.1",
					"timeAsSeconds": "0",
				}},
				"target_low": []req.Param{{
					"time":          "00:00",
					"value":         "0",
					"timeAsSeconds": "0",
				}},
				"target_high": []req.Param{{
					"time":          "00:00",
					"value":         "0",
					"timeAsSeconds": "0",
				}},
				"units": "mg/dL",
			}},
		"startDate": "1970-01-01T00:00:00.000Z",
		"mills":     "0",
		"units":     "mg/dL",
	}}

	resp, err := req.Post(NSSC+"api/v1/profile.json", req.BodyJSON(&profile), nsToken)
	glog.FatalIf(err)

	if status := resp.Response().StatusCode; status != 200 {
		glog.Fatalf("Nightscout configuration failed: %d", status)
	}

}

func TestNsscApiAccessible(t *testing.T) {
	// gosc := r.Get(GOSC + "api/v1/status")
	nssc, err := req.Get(NSSC + "api/v1/status")
	reqOK(nssc, err, t)
	assertEq("body status ok", nssc.String(), "<h1>STATUS OK</h1>", t)
}

func TestGoscApiAccessible(t *testing.T) {
	// gosc := r.Get(GOSC + "api/v1/status")
	gosc, err := req.Get(GOSC + "api/v1/status")
	reqOK(gosc, err, t)
	assertEq("body status ok", gosc.String(), "<h1>STATUS OK</h1>", t)
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
	assert("name is nightscout", strings.ToLower(ns.Name) == "nightscout", t)
	assert("name is goscout", gs.Name == "Goscout", t)
}

func getNightscoutToken() req.Header {
	h := sha1.New()
	h.Write([]byte(DefaultAPISecret))
	tok := fmt.Sprintf("%x", h.Sum(nil))
	return req.Header{
		"Api-Secret":   tok,
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}
}

func getEntryIds(entries *[]map[string]interface{}) []string {
	var ids []string

	for _, entry := range *entries {
		fmt.Printf("entry: %+v\n", entry)
		if id, ok := entry["_id"]; ok {
			ids = append(ids, id.(string))
		}
	}

	return ids
}

func getEntryDates(entries *[]map[string]interface{}) []string {
	var dates []string

	for _, entry := range *entries {
		fmt.Printf("entry: %+v\n", entry)
		if id, ok := entry["dateString"]; ok {
			dates = append(dates, id.(string))
		}
	}

	return dates
}

func filterEntriesWithIds(entries endpointsv1.Entries, ids []string) endpointsv1.Entries {
	var ret []endpointsv1.Entry

	idmap := make(map[string]interface{})
	for _, id := range ids {
		idmap[id] = nil
	}

	for _, entry := range entries {
		if _, ok := idmap[entry.ID]; ok {
			ret = append(ret, entry)
		}
	}

	return ret
}

func filterEntriesWithDates(entries endpointsv1.Entries, dates []string) endpointsv1.Entries {
	var ret []endpointsv1.Entry

	idmap := make(map[string]interface{})
	for _, id := range dates {
		idmap[id] = nil
	}

	for _, entry := range entries {
		if _, ok := idmap[entry.DateString]; ok {
			ret = append(ret, entry)
		}
	}

	return ret
}

func TestFilterEntriesWithIds(t *testing.T) {
	entries := endpointsv1.Entries{{ID: "a"}, {ID: "b"}, {ID: "c"}}
	processed := filterEntriesWithIds(entries, []string{"b"})
	assertEq("filterEntries returns one entry", len(processed), 1, t)
	assertEq("filterEntries returns correct entry", processed[0].ID, "b", t)
}

func helperUploadEntry(entry []req.Param, t *testing.T) (endpointsv1.Entries, endpointsv1.Entries) {
	nsToken := getNightscoutToken()

	nssc, err := req.Post(NSSC+"api/v1/entries.json", req.BodyJSON(&entry), nsToken)
	reqOK(nssc, err, t)

	gosc, err := req.Post(GOSC+"api/v1/entries.json", req.BodyJSON(&entry))
	reqOK(gosc, err, t)

	var nsjs []map[string]interface{}
	assert("json parse nightscout", nssc.ToJSON(&nsjs) == nil, t)
	glog.Infof("nsjs: %+v\n", nsjs)

	nsDates := getEntryDates(&nsjs)
	assert("at least one returned nightscout date", len(nsDates) > 0, t)

	var gojs []map[string]interface{}
	assert("json parse goscout", gosc.ToJSON(&gojs) == nil, t)
	glog.Infof("gojs: %+v\n", gojs)

	gsDates := getEntryDates(&gojs)
	assert("at least one returned goscout date", len(gsDates) > 0, t)

	// Nightscout does not return old dates without a find query.
	nssce, err := req.Get(NSSC + "api/v1/entries.json?find[date][$gt]=0&count=999")
	reqOK(nssc, err, t)

	gosce, err := req.Get(GOSC + "api/v1/entries.json?count=999")
	reqOK(gosc, err, t)

	var nsentries endpointsv1.Entries
	nssce.ToJSON(&nsentries)

	var gsentries endpointsv1.Entries
	gosce.ToJSON(&gsentries)

	ns := filterEntriesWithDates(nsentries, nsDates)
	assert(fmt.Sprintf("at least one filtered ns entry found for '%s' %+v", strings.Join(nsDates, ","), nsentries), len(ns) > 0, t)

	gs := filterEntriesWithDates(gsentries, gsDates)
	assert(fmt.Sprintf("at least one filtered gs entry found: '%s' %+v", strings.Join(gsDates, ","), gsentries), len(gs) > 0, t)

	// var ns endpointsv1.Entries
	// nssc.ToJSON(&ns)

	// var gs endpointsv1.Entries
	// gosc.ToJSON(&gs)

	return ns, gs
}

func helperAssertEntry(a, b endpointsv1.Entry, t *testing.T) {
	t.Helper()
	t.Logf("helperAssertEntry:\n%+v\n%+v\n", a, b)
	assertEq("entry device", a.Device, b.Device, t)
	assertEq("entry date int", a.Date, b.Date, t)
	assertEq("entry datestring parsed", parseTimeNS(a.DateString), parseTimeGS(b.DateString), t)
	assertEq("entry sgv", a.Sgv, b.Sgv, t)
	assertEq("entry delta", a.Delta, b.Delta, t)
	assertEq("entry direction", a.Direction, b.Direction, t)
	assertEq("entry type", a.Type, b.Type, t)
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

	assertEq("one entry in nightscout", len(ns), 1, t)
	assertEq("one entry in goscout", len(gs), 1, t)
	if len(ns) != 1 || len(gs) != 1 {
		t.Fatal("can't compare nonexistent entry")
	}

	helperAssertEntry(ns[0], gs[0], t)

}

func TestUploadEntries(t *testing.T) {
	entry := []req.Param{{
		"device":     "xDrip-DexcomG5 G5 Native",
		"date":       1549340059737,
		"dateString": "2019-02-04T23:14:19.737-0500",
		"sgv":        171,
		"delta":      1.999,
		"direction":  "Flat",
		"type":       "sgv",
		"filtered":   182682,
		"unfiltered": 178976,
		"rssi":       100,
		"noise":      1,
		"sysTime":    "2019-02-04T23:14:19.737-0500",
	}, {
		"device":     "xDrip-DexcomG5 G5 Native",
		"date":       1549340359670,
		"dateString": "2019-02-04T23:19:19.670-0500",
		"sgv":        168,
		"delta":      1.999,
		"direction":  "Flat",
		"type":       "sgv",
		"filtered":   179690,
		"unfiltered": 175542,
		"rssi":       100,
		"noise":      1,
		"sysTime":    "2019-02-04T23:19:19.670-0500",
	}}

	ns, gs := helperUploadEntry(entry, t)

	assertEq("two entries in nightscout", len(ns), 2, t)
	assertEq("two entries in goscout", len(gs), 2, t)
	if len(ns) != 2 || len(gs) != 2 {
		t.Fatalf("can't compare nonexistent entry (ns=%d, gs=%d)", len(ns), len(gs))
	}

	helperAssertEntry(ns[0], gs[0], t)
	helperAssertEntry(ns[1], gs[1], t)
}


func TestZZZZ_Cleanup(t *testing.T) {
	if ConfigKeepRunning {
		glog.Infof("Skipping cleanup due to ConfigKeepRunning")
		return
	}

	glog.Infof("Cleaning up...")
	for _, pid := range runningPids {
		glog.Infof("Sending SIGTERM to %d", pid)
		glog.WarningIf(syscall.Kill(pid, syscall.SIGTERM))
	}
}