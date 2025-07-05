package main 



import (
	"os"
	"sync"
	"time"
)

type LogStorage struct {
	file *os.File
	fileData  string
	directory string
	mutex 	sync.Mutex

}