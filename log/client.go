package log

import (
	"bytes"
	"fmt"
	stlog "log"
	"net/http"
	"simple_distributed/registry"
)

func SetClientLogger(serviceURL string, clientService registry.ServiceName) {
	stlog.SetPrefix(fmt.Sprintf("[%v] - ", clientService))
	stlog.SetFlags(0)
	stlog.SetOutput(&clientLogger{url: serviceURL})
}

type clientLogger struct {
	url string
}

func (cl clientLogger) Write(data []byte) (n int, err error) {
	//b := bytes.NewBuffer([]byte(data[:len(data)-1]))
	b := bytes.NewBuffer([]byte(data))
	resp, err := http.Post(cl.url+"/log", "text/plain", b)
	if err != nil {
		return 0, err
	}
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to send log message.service responeded error")
	}

	return len(data), nil
}
