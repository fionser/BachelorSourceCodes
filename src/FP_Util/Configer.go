package FP_Util

import (
	"bufio"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	COMMENT = "#"
	SPLITOR = "="
)

var (
	instance IConfiger
)

func init() {
	instance = NewConfiger()
	instance.Load("../configer.properties")
}

type IConfiger interface {
	Load(file string) bool
	Close()
	GetString(key, defaultVal string) string
	GetInt(key string, defaultVal int) int
}

func GetInstance() IConfiger {
	return instance
}

func NewConfiger() IConfiger {
	c := new(configer)
	c.constructor()
	return c
}

type configer struct {
	pairs map[string]string
	fd    *os.File
}

func (this *configer) constructor() {
	this.pairs = make(map[string]string)
	this.fd = nil
}

func (this *configer) parseLine(line string) {
	line = strings.TrimSpace(line)
	index := strings.Index(line, COMMENT)

	if index != -1 {
		if index != 0 {
			index -= 1
		}
		line = line[:index]
	}

	if len(line) > 0 {
		fields := strings.Split(line, SPLITOR)
		if len(fields) == 2 {
			this.pairs[fields[0]] = fields[1]
		}
	}
}

func (this *configer) doLoad() {
	reader := bufio.NewReader(this.fd)
	var line string
	var err error
	for {
		line, err = reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				this.parseLine(line)
			}
			break
		}
		this.parseLine(line)
	}
}

func (this *configer) Load(file string) (ok bool) {
	var err error
	this.fd, err = os.Open(file)
	if err != nil {
		ok = false
		log.Fatal("Can not load configuration file %s. Error %s\n", file, err.Error())
	} else {
		ok = true
		this.doLoad()
	}
	return
}

func (this *configer) Close() {
	if this.fd != nil {
		this.fd.Close()
	}
}

func (this *configer) GetString(key, defaultVal string) string {
	var returnVal string
	var ok bool
	if returnVal, ok = this.pairs[key]; !ok {
		returnVal = defaultVal
	}
	return returnVal
}

func (this *configer) GetInt(key string, defaultVal int) int {
	var returnVal string
	var returnInt int
	var ok bool
	var err error

	if returnVal, ok = this.pairs[key]; !ok {
		returnInt = defaultVal
	} else {
		if returnInt, err = strconv.Atoi(returnVal); err != nil {
			returnInt = defaultVal
			log.Fatalf("Configuration Error. Key %s has Value %s\n", key, returnVal)
		}
	}
	return returnInt
}
