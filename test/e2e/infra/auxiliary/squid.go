package auxiliary

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	squidImage         = "docker.io/ubuntu/squid:latest"
	squidContainerName = "e2e-squid-proxy"
	squidPort          = "3128/tcp"
	squidPortNum       = "3128"
)

// Squid holds connection info and the container for the aux Squid forward proxy.
type Squid struct {
	URL       string
	Host      string
	Port      string
	container testcontainers.Container
}

// Start starts the Squid proxy container and sets URL, Host, Port.
func (s *Squid) Start(ctx context.Context, network string, reuse bool) error {
	logrus.Infof("Starting Squid proxy container (reuse=%v)", reuse)
	configPath, err := createSquidConfig()
	if err != nil {
		return fmt.Errorf("failed to create squid config: %w", err)
	}
	req := testcontainers.ContainerRequest{
		Image:        squidImage,
		Name:         squidContainerName,
		ExposedPorts: []string{squidPort},
		Files: []testcontainers.ContainerFile{
			{HostFilePath: configPath, ContainerFilePath: "/etc/squid/squid.conf", FileMode: 0644},
		},
		WaitingFor: wait.ForListeningPort(squidPortNum),
		SkipReaper: reuse,
	}
	container, err := CreateContainer(ctx, req, reuse, WithNetwork(network), WithHostAccess())
	if err != nil {
		return fmt.Errorf("failed to start squid proxy container: %w", err)
	}
	s.container = container
	s.Host = GetHostIP()
	mappedPort, err := container.MappedPort(ctx, squidPort)
	if err != nil {
		return fmt.Errorf("get mapped port for %s: %w", squidPort, err)
	}
	s.Port = mappedPort.Port()
	s.URL = fmt.Sprintf("http://%s", net.JoinHostPort(s.Host, s.Port))
	logrus.Infof("Squid proxy container started: %s", s.URL)
	return nil
}

// AccessLogs returns the Squid container's stdout which contains access log entries.
func (s *Squid) AccessLogs() (string, error) {
	//nolint:gosec // G204: squidContainerName is a package constant, not user input.
	cmd := exec.Command(containerRuntimeCLIName(), "logs", squidContainerName)
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to read squid logs: %w", err)
	}
	return string(out), nil
}

const squidConfigContent = `http_port 3128

# Allow all source networks (E2E test only)
acl localnet src all

# Allow CONNECT to any port — default squid.conf restricts CONNECT
# to port 443 via SSL_ports, but flightctl uses non-standard ports.
http_access allow CONNECT
http_access allow localnet
http_access allow localhost
http_access deny all

# Log to stdout so container logs capture access entries
access_log stdio:/dev/stdout
cache_log stdio:/dev/stderr
`

func createSquidConfig() (string, error) {
	tmpPath := filepath.Join(os.TempDir(), "e2e-squid-proxy.conf")
	return tmpPath, os.WriteFile(tmpPath, []byte(squidConfigContent), 0600)
}
