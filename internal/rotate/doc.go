// Package rotate detects log file rotation by monitoring inode
// changes and file truncation, emitting notifications so the tail
// reader can reopen the file and resume streaming from the start.
//
// Usage:
//
//	mon, err := rotate.New("/var/log/app.log")
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer mon.Stop()
//
//	for range mon.Notify() {
//		// file has been rotated – reopen it
//	}
package rotate
