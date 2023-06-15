package grades

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func RegisterHandlers() {
	handler := new(studentsHandler)
	http.Handle("/students", handler)
	http.Handle("/students/", handler)

}

type studentsHandler struct{}

//func (sh studentsHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
//	//TODO implement me
//	panic("implement me")
//}

// /students
// /students/{id}
// /students/{id}/grades
// ServeHTTP 是 Handler接口的一个方法，接口studentHandler实现该方法才能使用 http.Handle
func (sh studentsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 通过分割 / 的方式来匹配请求
	pathSegments := strings.Split(r.URL.Path, "/")
	switch len(pathSegments) {
	case 2:
		sh.getAll(w, r)
	case 3:
		id, err := strconv.Atoi(pathSegments[2])
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		sh.getOne(w, r, id)
	case 4:
		id, err := strconv.Atoi(pathSegments[2])
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		sh.addGrade(w, r, id)
	default:
		w.WriteHeader(http.StatusNotFound)
	}

}

func (sh studentsHandler) getAll(w http.ResponseWriter, r *http.Request) {
	studentsMutex.Lock()
	defer studentsMutex.Unlock()

	data, err := sh.toJSON(students)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(data)
}

func (sh studentsHandler) toJSON(obj interface{}) ([]byte, error) {
	var b bytes.Buffer
	// 新建编码器
	enc := json.NewEncoder(&b)
	// 对数据进行编码
	err := enc.Encode(obj)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize students: %q", err)
	}
	return b.Bytes(), nil
}

func (sh studentsHandler) getOne(w http.ResponseWriter, r *http.Request, id int) {
	studentsMutex.Lock()
	defer studentsMutex.Unlock()
	// 获取学生
	student, err := students.GetByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Println(err)
		return
	}
	// 通过toJSON将student转换成byteSlice
	data, err := sh.toJSON(student)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("failed to serialize student: %q\n", err)
		return
	}
	// 添加头
	w.Header().Add("Content-Type", "application/json")
	// 写数据
	w.Write(data)
}

func (sh studentsHandler) addGrade(w http.ResponseWriter, r *http.Request, id int) {
	studentsMutex.Lock()
	defer studentsMutex.Unlock()
	// 获取学生
	student, err := students.GetByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Println(err)
		return
	}
	var g Grade
	dec := json.NewDecoder(r.Body)
	// 获取到json数据中的Grade
	err = dec.Decode(&g)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}
	// 添加课程
	student.Grades = append(student.Grades, g)
	w.WriteHeader(http.StatusCreated)
	data, err := sh.toJSON(g)
	if err != nil {
		log.Println(err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	// 返回写入的数据
	w.Write(data)

}
