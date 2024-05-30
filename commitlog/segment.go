package commitlog

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

const (
	logSuffix   = ".log"
	indexSuffix = ".index"
	fileFormat  = "%05d%s"
)

type segment struct {
	writer         io.Writer
	reader         io.Reader
	log            *os.File
	index          *index
	path           string
	position       int
	maxBytes       int
	startingOffset int
	nextOffset     int
	file           string
	mu             sync.Mutex
}

/* Create a new segment */
func newSegment(directory string, offset int) (*segment, error) {
	seg := &segment{
		maxBytes:       SEGMENT_MAX_BYTES,
		position:       0,
		startingOffset: offset,
		nextOffset:     offset,
		file:           directory,
	}
	seg.path = seg.logPath()
	loggly, err := os.OpenFile(seg.path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return seg, err
	}
	seg.log = loggly
	seg.reader = loggly
	seg.writer = loggly

	ind, err := newIndex(seg.indexPath())
	if err != nil {
		return seg, err
	}
	//Add index pointer to our segment
	seg.index = ind

	return seg, nil
}

/* Load in a segment from disk, using the path to the logs and path to the index */
func loadSegment(logPath string, indexPath string) (*segment, error) {
	logBase := filepath.Base(logPath)
	offsetStr := strings.TrimSuffix(logBase, logSuffix)
	baseOffset, err := strconv.Atoi(offsetStr)
	if err != nil {
		return nil, err
	}

	loggly, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	fi, err := loggly.Stat()
	if err != nil {
		return nil, err
	}

	seg := &segment{
		path:           logPath,
		maxBytes:       SEGMENT_MAX_BYTES,
		position:       int(fi.Size()),
		startingOffset: baseOffset,
		log:            loggly,
		reader:         loggly,
		writer:         loggly,
	}

	ind, err := newIndex(indexPath)
	if err != nil {
		return seg, err
	}

	seg.index = ind

	if ind.indexFile == nil {
		return seg, errors.New("pointer to segment file is nil")
	}

	totalEntries, err := seg.index.loadIndex()
	if err != nil {
		slog.Error(err.Error())
		return seg, err
	}
	seg.nextOffset = totalEntries + seg.startingOffset
	return seg, nil
}

/* Write a new log to the segment. This will add a log to the log file with the data, and the metadata information to the index */
func (seg *segment) write(message []byte) (int, error) {
	seg.mu.Lock()
	defer seg.mu.Unlock()
	//Check the byte status..
	fileInfo, err := seg.log.Stat()
	if err != nil {
		return 0, err
	}
	computedSize := fileInfo.Size() + int64(len(message))
	//Strickly greater than
	if computedSize > int64(seg.maxBytes) {
		return 0, errors.New("max segment length")
	}

	numOfBytes, err := seg.writer.Write(message)
	if err != nil {
		return 0, err
	}

	err = seg.index.addEntry(seg.position, numOfBytes)
	if err != nil {
		return 0, err
	}

	seg.position += numOfBytes
	seg.nextOffset++

	return numOfBytes, nil
}

/*Given an arbitrary offset, read the data stored on the segment*/
func (seg *segment) readAt(offset int) (returnBuff []byte, err error) {
	seg.mu.Lock()
	defer seg.mu.Unlock()
	var buff []byte
	//Do we want to do this with the base offset??
	if offset >= seg.nextOffset-seg.startingOffset {
		slog.Error(fmt.Sprintf("offset given: %d, max offset: %d", offset, seg.nextOffset-seg.startingOffset))
		return nil, errors.New("offset out of bounds")
	} else {
		ent := seg.index.entries[offset]
		buff = make([]byte, ent.Total)
		seg.log.ReadAt(buff, int64(ent.Start))
	}
	return buff, nil
}

func (s *segment) logPath() string {
	return filepath.Join(s.file, fmt.Sprintf(fileFormat, s.startingOffset, logSuffix))
}

func (s *segment) indexPath() string {
	return filepath.Join(s.file, fmt.Sprintf(fileFormat, s.startingOffset, indexSuffix))
}

func (seg *segment) delete() error {
	if err := seg.close(); err != nil {
		return err
	}
	seg.mu.Lock()
	defer seg.mu.Unlock()
	if err := os.Remove(seg.path); err != nil {
		return err
	}
	if err := os.Remove(seg.index.path); err != nil {
		return err
	}
	return nil
}

func (seg *segment) close() error {
	seg.mu.Lock()
	defer seg.mu.Unlock()
	if err := seg.log.Close(); err != nil {
		return err
	}
	if err := seg.index.close(); err != nil {
		return err
	}
	return nil
}

/* Deprecated functions */

// func (seg *segment) read(offset int64, total int32) (string, error) {
// 	seg.mu.Lock()
// 	defer seg.mu.Unlock()
// 	_, err := seg.log.Seek(offset, 0)
// 	if err != nil {
// 		logger.Error.Println(err)
// 		return "", err
// 	}
// 	b2 := make([]byte, total)
// 	n2, err := seg.reader.Read(b2)
// 	if err != nil {
// 		logger.Error.Println(err)
// 		return "", err
// 	}
// 	return string(b2[:n2]), nil
// }
