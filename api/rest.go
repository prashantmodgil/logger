package api
import (
	"fmt"
	"encoding/json"
    bucket "getLogs/bucket/minio"
	"os"
	"io/ioutil"
	"log"
	"net/http"
	"math/rand"
	"time"
	"github.com/minio/minio-go/v6"
	"strings"
)

func SearchLogs(w http.ResponseWriter, r *http.Request) {
    searchKeyword := r.URL.Query().Get("searchKeyword")
    fromStr := r.URL.Query().Get("from")
    toStr := r.URL.Query().Get("to")

    fromTime, err := time.Parse(time.RFC3339, fromStr)
    if err != nil {
        http.Error(w, "Invalid 'from' timestamp"+err.Error(), http.StatusBadRequest)
        return
    }
    toTime, err := time.Parse(time.RFC3339, toStr)
    if err != nil {
        http.Error(w, "Invalid 'to' timestamp"+err.Error(), http.StatusBadRequest)
        return
    }

	doneCh := make(chan struct{})
	defer close(doneCh)
    objectsCh := bucket.MinioClient.ListObjects("mylogs", "",true,doneCh)
        var matchingLogLines []LogLine
        for obj := range objectsCh {
            if obj.Err != nil {
                http.Error(w, "Failed to list objects"+obj.Err.Error(), http.StatusInternalServerError)
                return
            }
            if !strings.HasSuffix(obj.Key, ".txt") {
                continue
            }
			fileFolder := strings.Split(obj.Key,"/")
			if len(fileFolder)!=2{
				continue
			}
			fileName := strings.Split(strings.TrimSuffix(fileFolder[1], ".txt"),"_")
			input := fileFolder[0]+"-"+fileName[1]+":00"
			layout := "2006-January-2-15:04"

			fileTimestamp, err := time.Parse(layout, input)
            if err != nil {
                log.Printf("Failed to parse timestamp from log file name: %v", err)
                continue
            }
			
            if fileTimestamp.After(fromTime) && (fileTimestamp.Before(toTime)|| fileTimestamp.Equal(toTime)){
                reader, err := bucket.MinioClient.GetObject("mylogs", obj.Key, minio.GetObjectOptions{})
                if err != nil {
                    log.Printf("Failed to download log file %s: %v", obj.Key, err)
                    continue
                }
                defer reader.Close()
                fileContents, err := ioutil.ReadAll(reader)
                if err != nil {
                    log.Printf("Failed to read log file %s: %v", obj.Key, err)
                    continue
                }
                lines := strings.Split(string(fileContents), "\n")
                for _, line := range lines {
                    if strings.Contains(line, searchKeyword) {
                        matchingLogLines = append(matchingLogLines, LogLine{
                            Timestamp: fileTimestamp,
                            Message:   line,
                        })
                    }
                }
            }
        }

        responseJSON, err := json.Marshal(matchingLogLines)
        if err != nil {
            http.Error(w, "Failed to marshal JSON response", http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        w.Write(responseJSON)
}

func SeedHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
	sl :=seedLogs[rand.Intn(len(seedLogs))]
	log.Printf("[%s] %s %s %s", time.Now().Format(LOGTIMEFORMATE), r.Method, r.URL.Path,sl)
    fmt.Fprintf(w, sl)
}

func LogRequest(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        hour := time.Now().Hour()
		year, month, day := time.Now().Date()
		date := strings.ToLower(fmt.Sprintf("%v-%v-%v",year, month, day))
		
        logFileName := fmt.Sprintf("log_%d.txt", hour)
        logFile, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
        if err != nil {
            log.Fatalf("Error opening log file: %v", err)
        }
        defer logFile.Close()

        log.SetOutput(logFile)

        log.Printf("[%s] %s %s", time.Now().Format(LOGTIMEFORMATE), r.Method, r.URL.Path)
		next(w, r)

        if err := bucket.UploadLogToMinio(logFileName, date,hour); err != nil {
            log.Printf("Error uploading log to MinIO: %v", err)
            http.Error(w, "Internal Server Error"+err.Error(), http.StatusInternalServerError)
            return
        }

    }
}



