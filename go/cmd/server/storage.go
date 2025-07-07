package server



import (
	"os"
	"sync"
	"time"
	"fmt"
)

type LogStorage struct {
	file *os.File
	fileData  string
	directory string
	mutex 	sync.Mutex

}



func NewLogStorage(directory string) (*LogStorage, error) {

//create directory if it does not exist
	if err := os.MkdirAll(directory, 0700);
	err != nil {
		return nil, fmt.Errorf("failed to create directory : %w", err)

	}


	//Get current date for file name 
	currentDate := time.Now().Format("2005-01-02")
	filename := filepath.Join(directory + "security_" + currentDatea + ".log")


	// opne log file
	
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file : %w", err)
	}


	//Return initialized LogStorage 
	return &LogStorage{
		file: file,
		fileDate: currentDate,
		directory: directory,
	}, nil 
	}


		func (s *LogStorage)WriteLog(source, ip , message string) error {
		s.mutex.Lock()
		defer s.mutex.Unlock()

		//check if we need to rotate logs (new day)
		currentDate := time.Now().Format("2006-01-02")
		if currentDate != s.fileDate {
			if err := s.rotate(currentDate); err != nil {
				return fmt.Errorf("failed to rotate log file: %w", err)
			}
		}
		//write log entry
		logEntry := formatLogEntry(source, ip , message)
		 err := s.file.WriteString(logEntry)
		if err != nil {
			return fmt.Errorf("failed to write log entry: %w", err)
		}

		
		func (s *LogStorage) rotate(newDate string) error {
			//close current file
			if err := s.file.Close(); err != nil {
				return fmt.Errorf("failed to close current log file: %w", err)
			}

			//Create new file with new date
			filename := filepath.Join(s.directory, "security_" + newDate + ".log")

			// Open new file
			file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
			if err != nil {
				return fmt.Errorf("failed to open new log file: %w", err)
			}

			//Update storage with new file and date
			s.file = file
			s.fileDate = newDate
			return nil

		}




func formatLogEntry(source, ip, message string) string {
	return fmt.Sprintf("[%s] [%s] [%s] %s\n", 
		time.Now().UTC().Format(time.RFC3339),
		source,
		ip,
		message)
}


