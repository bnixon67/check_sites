# check_sites

check_sites will attempt to connect to each URL provided on stdin

each check will be done via a goroutine to allow multi-tasking since the connection may take some time to process or timeout
