package commitlog

import (
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"sync/atomic"

	"integrand/utils"
)

var (
	//Max bytes of the log, used by segment
	SEGMENT_MAX_BYTES = utils.GetEnvrionmentVariableInt("SEGMENT_MAX_BYTES", 1000000000)
	//Max retention byte of the commitlog, used by cleaner
	RETENTION_BYTES = utils.GetEnvrionmentVariableInt("RETENTION_BYTES", 10000000000)
)

type Commitlog struct {
	path            string     //Path to the partition directory
	segments        []*segment //Individual segment files
	mu              sync.RWMutex
	vCurrentSegment atomic.Value
}

/*
New will open a commitlog using the provided path.
If no partion already exists, then one will be created.
*/
func New(path string) (*Commitlog, error) {
	cl := Commitlog{
		path: path,
	}
	_, err := os.Stat(path)
	if err != nil {
		//Partitions dont exist, create a new directory and return
		err := os.Mkdir(path, 0755)
		if err != nil {
			slog.Error(err.Error())
			return nil, err
		}
		return &cl, nil
	}
	//Since our partition already exists we need to load them in
	cl.loadSegments()
	return &cl, nil
}

func (cl *Commitlog) Delete() error {
	err := os.RemoveAll(cl.path)
	return err
}

func (cl *Commitlog) Clear() error {
	d, err := os.Open(cl.path)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(cl.path, name))
		if err != nil {
			return err
		}
	}
	//set the oldest offset and newest offset
	return nil
}

func (cl *Commitlog) GetPath() string {
	return cl.path
}

/*
Append appends a new entry to the commitlog
*/
func (cl *Commitlog) Append(message []byte) error {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	if len(cl.segments) == 0 {
		slog.Info("Creating a new segment")
		segment, err := newSegment(cl.path, 0)
		if err != nil {
			slog.Error("Unable to create new segment:")
			slog.Error(err.Error())
			cl.mu.Unlock()
			return err
		}
		cl.segments = append(cl.segments, segment)
		cl.vCurrentSegment.Store(segment)
	}
	curSegment, _ := cl.getCurrentSegment()
	_, err := curSegment.write(message)
	if err != nil {
		if err.Error() == "max segment length" {
			cl.mu.Unlock()
			// Clean the segments first....
			err = cl.clean()
			if err != nil {
				return err
			}
			//Check for error if too many bytes in the segment -> then split
			err = cl.split()
			if err != nil {
				return err
			}
			//Append again, this time on the new segment...
			cl.Append(message)
			cl.mu.Lock() //So defer unlock works..
			return nil
		}
		return err
	}
	return nil
}

/*
loadSegments loads all of the segments from disk into memory to read
*/
func (cl *Commitlog) loadSegments() error {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	files, err := os.ReadDir(cl.path)
	if err != nil {
		slog.Error("Unable to read directory")
		slog.Error(err.Error())
	}

	filesStrings := []string{}
	filesToRemove := []string{}
	for _, file := range files {
		filesStrings = append(filesStrings, file.Name())
	}
	sort.Strings(filesStrings)

	for _, file := range files {
		if strings.HasSuffix(file.Name(), logSuffix) {
			//Log file
			corespondingIndexFile := strings.Replace(file.Name(), logSuffix, indexSuffix, 1)
			isPresentValue := sort.SearchStrings(filesStrings, corespondingIndexFile)
			if isPresentValue == len(filesStrings) {
				filesToRemove = append(filesToRemove, file.Name())
			}
			// else {
			//This will not work if the path is passed in as "path/to/"
			// logfileToLoad := cl.path + "/" + file.Name()
			// indexFileToLoad := cl.path + "/" + corespondingIndexFile
			// seg, err := loadSegment(indexFileToLoad, logfileToLoad)
			// if err != nil {
			// 	logger.Error.Fatal(err)
			// }
			// cl.segments = append(cl.segments, seg)
			// }

		} else if strings.HasSuffix(file.Name(), indexSuffix) {
			corespondingLogFile := strings.Replace(file.Name(), indexSuffix, logSuffix, 1)
			isPresentValue := sort.SearchStrings(filesStrings, corespondingLogFile)
			if isPresentValue == len(filesStrings) {
				if err := os.Remove(file.Name()); err != nil {
					return err
				}
			} else {
				logfileToLoad := cl.path + "/" + corespondingLogFile
				indexFileToLoad := cl.path + "/" + file.Name()
				seg, err := loadSegment(logfileToLoad, indexFileToLoad)
				if err != nil {
					slog.Error(err.Error())
					return err
				}
				cl.segments = append(cl.segments, seg)
			}
		}
	}
	//Delete all of the corrupted files...
	//Should this be in a goroutine??
	for _, file := range filesToRemove {
		fullFilePath := cl.path + "/" + file
		if err := os.Remove(fullFilePath); err != nil {
			return err
		}
	}
	//If there are no segments, then we cannot set an active segment
	if len(cl.segments) > 0 {
		//Set our active segment to the latest segment
		latestSegment := cl.segments[len(cl.segments)-1]
		cl.vCurrentSegment.Store(latestSegment)
	}
	return nil
}

/*
Split will create a new segment based on the current segment within the commitlog
*/
func (cl *Commitlog) split() error {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	seg, _ := cl.getCurrentSegment()
	segment, err := newSegment(cl.path, seg.nextOffset)
	if err != nil {
		return err
	}
	cl.segments = append(cl.segments, segment)
	//Set our current segment to the new segment created
	cl.vCurrentSegment.Store(segment)

	return nil
}

func (cl *Commitlog) getCurrentSegment() (*segment, error) {
	if len(cl.segments) == 0 {
		return nil, errors.New("commitlog is empty")
	}
	return cl.vCurrentSegment.Load().(*segment), nil
}

/*
Read will read a specific 'offset' within the commitlog,
if the offset is out of bounds returns error
*/
func (cl *Commitlog) Read(offset int) ([]byte, error) {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	//This should be optomized. Should probably check the current segment first?
	//Then do this??
	var segmentContainsOffset *segment
	// Get the segment that holds the offset
	for _, seg := range cl.segments {
		if seg.startingOffset <= offset {
			segmentContainsOffset = seg
		} else {
			break
		}
	}

	if segmentContainsOffset == nil {
		return nil, errors.New("offset does not exist")
	}

	searchOffset := offset - segmentContainsOffset.startingOffset
	buff, err := segmentContainsOffset.readAt(searchOffset)
	if err != nil {
		return nil, err
	}
	return buff, err
}

type CommitlogDetails struct {
	OldestOffset   int
	NextOffset     int
	TotalSegments  int
	RetentionBytes int
}

func (cl *Commitlog) GetCommitlogDetails() CommitlogDetails {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	oldestOffset := 0
	nextOffset := 0
	currentSeg, err := cl.getCurrentSegment()
	if err == nil {
		nextOffset = currentSeg.nextOffset
		oldestOffset = cl.segments[0].startingOffset
	}
	return CommitlogDetails{
		OldestOffset:   oldestOffset,
		NextOffset:     nextOffset,
		TotalSegments:  len(cl.segments),
		RetentionBytes: RETENTION_BYTES,
	}
}
