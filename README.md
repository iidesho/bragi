# Bragi

Bragi is a simple log lib that is designed to be a dropin repleasement for go's std log lib. But with log output that reflect the output given with java logback. This lib is also going to rotate log files and remove old logs.

## How to use

This is a dropin replace ment, so just add the dependencie and use it as you would do go log.
To log to file add the following

```go
log.SetPrefix("vili")
cloaser := log.SetOutputFolder(logDir)
if cloaser == nil {
	log.Fatal("Unable to sett logdir")
}
defer cloaser()
```

SetPrefix sets the prefix name of the logs. When they rotate the one without any time data in the name is the current and the rest will have the date appended.

SetOutputFolder creates and sets the output folder for bouth human and json logs and returns a function to cloase the files if the process was successfull. 

To use the rotating feature add a call to ` log.StartRotate(done chan func()) ` after the code example above. Full example below

```go
log.SetPrefix("vili")
cloaser := log.SetOutputFolder(logDir)
if cloaser == nil {
	log.Fatal("Unable to sett logdir")
}
defer cloaser()
log.StartRoute(nil)
```

## Extra functions

If you want to add an error to debug, info, notice, error, crit, fatal with the following pattern

```go
 log.AddError(err).Debug("Information about the error") 
```

With that pattern you will get bouth the error and the information text in the same log object