// Package rotate detects log file rotation by monitoring inode changes
// and notifies consumers so they can reopen the file from the beginning.
//
// Usage:
//
//	r, err := rotate.New("/var/log/app.log", rotate.DefaultInterval)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer r.Stop()
//
//	for range r.Notify() {
//		// file has been rotated – reopen it
//	}
package rotate
