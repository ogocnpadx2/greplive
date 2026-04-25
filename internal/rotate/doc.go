// Package rotate detects log file rotation by monitoring inode changes
// and file size resets, notifying consumers so they can reopen the file
// and continue reading from the beginning.
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
//		// file has been rotated — reopen it
//	}
package rotate
