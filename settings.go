package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/op/go-logging"
	"golang.org/x/net/websocket"

	"github.com/microscaling/microscaling/config"
	"github.com/microscaling/microscaling/demand"
	"github.com/microscaling/microscaling/engine"
	"github.com/microscaling/microscaling/engine/localEngine"
	"github.com/microscaling/microscaling/engine/serverEngine"
	"github.com/microscaling/microscaling/monitor"
	"github.com/microscaling/microscaling/scheduler"
	"github.com/microscaling/microscaling/scheduler/docker"
	"github.com/microscaling/microscaling/scheduler/kubernetes"
	"github.com/microscaling/microscaling/scheduler/marathon"
	"github.com/microscaling/microscaling/scheduler/toy"
)

type settings struct {
	schedulerType   string
	sendMetrics     bool
	monitorTypes    string
	microscalingAPI string
	userID          string
	pullImages      bool
	dockerHost      string
	demandEngine    string
	marathonAPI     string
	config          string
	kubeConfig      string
	kubeNamespace   string
}

func initLogging() {
	// The MSS_LOG_DEBUG environment variable controls what logging is output
	// By default the log level is INFO for all components
	// Adding a component name to MSS_LOG_DEBUG makes its logging level DEBUG
	// In addition, if "detail" is included in the environment variable details of the process ID and file name / line number are included in the logs
	// MSS_LOG_DEBUG="all" - turn on DEBUG for all components
	// MSS_LOG_DEBUG="mssapi,detail" - turn on DEBUG for the api package, and use the detailed logging format
	basicLogFormat := logging.MustStringFormatter(`%{color}%{level:.4s} %{time:15:04:05.000}: %{color:reset} %{message}`)
	detailLogFormat := logging.MustStringFormatter(`%{color}%{level:.4s} %{time:15:04:05.000} %{pid} %{shortfile}: %{color:reset} %{message}`)

	logComponents := getEnvOrDefault("MSS_LOG_DEBUG", "none")
	if strings.Contains(logComponents, "detail") {
		logging.SetFormatter(detailLogFormat)
	} else {
		logging.SetFormatter(basicLogFormat)
	}

	logBackend := logging.NewLogBackend(os.Stdout, "", 0)
	logging.SetBackend(logBackend)

	var components = []string{"mssagent", "mssapi", "mssconfig", "mssdemand", "mssengine", "mssmetric", "mssscheduler", "msstarget", "mssutils"}

	for _, component := range components {
		if strings.Contains(logComponents, component) || strings.Contains(logComponents, "all") {
			logging.SetLevel(logging.DEBUG, component)
		} else {
			logging.SetLevel(logging.INFO, component)
		}
	}
}

func getSettings() settings {
	var st settings
	st.schedulerType = getEnvOrDefault("MSS_SCHEDULER", "DOCKER")
	st.microscalingAPI = getEnvOrDefault("MSS_API_ADDRESS", "app.microscaling.com")
	st.userID = getEnvOrDefault("MSS_USER_ID", "5k5gk")
	st.sendMetrics = (getEnvOrDefault("MSS_SEND_METRICS_TO_API", "true") == "true")
	st.monitorTypes = getEnvOrDefault("MSS_MONITOR", "SERVER")
	st.pullImages = (getEnvOrDefault("MSS_PULL_IMAGES", "true") == "true")
	st.dockerHost = getEnvOrDefault("DOCKER_HOST", "unix:///var/run/docker.sock")
	st.demandEngine = getEnvOrDefault("MSS_DEMAND_ENGINE", "LOCAL")
	st.marathonAPI = getEnvOrDefault("MSS_MARATHON_API", "http://localhost:8080")
	st.config = getEnvOrDefault("MSS_CONFIG", "SERVER")
	// To run locally set kube config location. Otherwise uses the built in cluster config.
	st.kubeConfig = getEnvOrDefault("MSS_KUBE_CONFIG", "")
	st.kubeNamespace = getEnvOrDefault("MSS_KUBE_NAMESPACE", "default")
	return st
}

func getScheduler(st settings, demandUpdate chan struct{}) (scheduler.Scheduler, error) {
	var s scheduler.Scheduler

	switch st.schedulerType {
	case "DOCKER":
		log.Info("Scheduling with Docker remote API")
		s = docker.NewScheduler(st.pullImages, st.dockerHost)
	case "MARATHON":
		log.Info("Scheduling with Mesos / Marathon")
		s = marathon.NewScheduler(st.marathonAPI, demandUpdate)
	case "ECS":
		return nil, fmt.Errorf("Scheduling with ECS not yet supported. Tweet with hashtag #MicroscaleECS if you'd like us to add this next!")
	case "KUBERNETES":
		log.Info("Scheduling with Kubernetes")
		s = kubernetes.NewScheduler(st.kubeConfig, st.kubeNamespace, demandUpdate)
	case "NOMAD":
		return nil, fmt.Errorf("Scheduling with Nomad not yet supported. Tweet with hashtag #MicroscaleNomad if you'd like us to add this next!")
	case "TOY":
		log.Info("Scheduling with toy scheduler")
		s = toy.NewScheduler()
	default:
		return nil, fmt.Errorf("Bad value for MSS_SCHEDULER: %s", st.schedulerType)
	}

	if s == nil {
		return nil, fmt.Errorf("No scheduler")
	}

	return s, nil
}

func getTasks(st settings) (tasks *demand.Tasks, err error) {
	var c config.Config

	tasks = new(demand.Tasks)

	// Get the tasks that have been configured by this user
	switch st.config {
	case "FILE":
		return nil, fmt.Errorf("Config file not yet supported")
	case "SERVER":
		c = config.NewServerConfig(st.microscalingAPI)
	case "HARDCODED":
		c = config.NewHardcodedConfig()
	case "LABEL":
		switch st.schedulerType {
		case "DOCKER":
			// Gets labels for the configured image from the MicroBadger API
			c = config.NewLabelConfig(st.microscalingAPI)
		case "KUBERNETES":
			// Gets the image from the k8s deployment and the labels from the MicroBadger API
			c = config.NewKubeLabelConfig(st.microscalingAPI, st.kubeConfig, st.kubeNamespace)
		default:
			return nil, fmt.Errorf("Label config not supported for scheduler: %s", st.config)
		}
	default:
		return nil, fmt.Errorf("Bad value for MSS_CONFIG: %s", st.config)
	}

	t, maxContainers, err := c.GetApps(st.userID)
	tasks.MaxContainers = maxContainers

	// For now pass the whole environment to all containers.
	globalEnv := os.Environ()

	for _, task := range t {
		task.Env = globalEnv
		log.Debugf("%+v", task)
	}

	tasks.Tasks = t

	if err != nil {
		log.Errorf("Error getting tasks: %v", err)
	}

	return tasks, err
}

func getDemandEngine(st settings, ws *websocket.Conn) (e engine.Engine, err error) {
	switch st.demandEngine {
	case "LOCAL":
		log.Info("Calculate demand locally")
		e = localEngine.NewEngine()
	case "SERVER":
		log.Info("Get demand from server")
		e = serverEngine.NewEngine(ws)
	default:
		return nil, fmt.Errorf("Bad value for MSS_DEMAND_ENGINE: %s", st.demandEngine)
	}
	return e, nil
}

func getMonitors(st settings, ws *websocket.Conn) (m []monitor.Monitor) {
	// Monitor is where we send results & output. There might be more than one so we return a list
	if strings.Contains(st.monitorTypes, "SERVER") {
		log.Info("Server is a monitor")
		ms := monitor.NewServerMonitor(ws, st.userID)
		m = append(m, ms)
	}

	return
}

func getEnvOrDefault(name string, defaultValue string) string {
	v := os.Getenv(name)
	if v == "" {
		v = defaultValue
	}

	return v
}
