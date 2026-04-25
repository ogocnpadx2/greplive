// Package rotate detects log file rotation by monitoring inode changes
// and file size resets, notifying consumers so they can reopen the file
// and continue tailing from the beginning.
package rotate
