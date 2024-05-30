package commitlog

import (
	"fmt"
	"log/slog"
)

func (cl *Commitlog) clean() error {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	cleanedSegments, err := cl.deleteCleaner()
	if err != nil {
		return err
	}
	cl.segments = cleanedSegments
	return nil
}

// Clean by size
func (cl *Commitlog) deleteCleaner() ([]*segment, error) {
	var cleanedSegments []*segment
	if len(cl.segments) == 0 {
		//return custom error...
		return cl.segments, nil
	}

	//If the retention bytes are set to a negative number,
	// We want to store these values indefinetly
	if RETENTION_BYTES < 0 {
		return cl.segments, nil
	}

	bytesSize := 0
	var i int
	for i = len(cl.segments) - 1; i >= 0; i-- {
		seg := cl.segments[i]
		if bytesSize > RETENTION_BYTES {
			break
		}
		slog.Info(fmt.Sprintf("Keeping: %s", seg.path))
		cleanedSegments = append(cleanedSegments, seg)
		bytesSize += seg.position
	}

	for j := 0; j <= i; j++ {
		seg := cl.segments[j]
		slog.Info(fmt.Sprintf("Deleting: %s", seg.path))
		if err := seg.delete(); err != nil {
			slog.Error("unable to delete segment")
			return cleanedSegments, err
		}
	}
	return cleanedSegments, nil
}
