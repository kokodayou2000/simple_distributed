package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

const ServerPort = ":3000"
const ServicesURL = "http://localhost" + ServerPort + "/services"

type registry struct {
	registrations []Registration
	mutex         *sync.RWMutex
}

// 注册
func (r *registry) add(reg Registration) error {
	r.mutex.Lock()
	r.registrations = append(r.registrations, reg)
	r.mutex.Unlock()
	// 添加依赖请求
	err := r.sendRequiredServices(reg)
	return err
}

func (r *registry) sendRequiredServices(reg Registration) error {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	var p patch
	// foreach 已经注册的服务
	for _, serviceReg := range r.registrations {
		// 需要的服务
		for _, reqService := range reg.RequiredServices {
			if serviceReg.ServiceName == reqService {
				// 添加服务信息
				p.Added = append(p.Added, patchEntry{
					Name: serviceReg.ServiceName,
					URL:  serviceReg.ServiceURL,
				})
			}
		}
	}
	// send patch
	err := r.sendPatch(p, reg.ServiceUpdateURL)
	if err != nil {
		return err
	}
	return nil
}

// sendPatch 向指定的url发送patch
func (r *registry) sendPatch(p patch, url string) error {
	d, err := json.Marshal(p)
	if err != nil {
		return err
	}
	_, err = http.Post(url, "application/json", bytes.NewBuffer(d))
	if err != nil {
		return err
	}
	return nil
}

func (r *registry) remove(url string) error {

	for i := range reg.registrations {
		if reg.registrations[i].ServiceURL == url {
			r.mutex.Lock()
			reg.registrations = append(reg.registrations[:i], r.registrations[i+1:]...)
			r.mutex.Unlock()
			return nil
		}
	}
	return fmt.Errorf("service at URL %s not found", url)
}

var reg = registry{
	registrations: make([]Registration, 0),
	mutex:         new(sync.RWMutex),
}

// RegistryService 注册服务
type RegistryService struct {
}

func (s RegistryService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Request received")
	switch r.Method {
	case http.MethodPost:
		dec := json.NewDecoder(r.Body)
		var r Registration
		err := dec.Decode(&r)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Printf("Adding service: %v  with URL: %s\n", r.ServiceName, r.ServiceURL)
		err = reg.add(r)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
