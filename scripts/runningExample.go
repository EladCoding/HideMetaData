package scripts

// Run a simulation of the mixnet architecture, with one client for the user to play with.
func PlayingExample() {
	go RunInfrastructure()
	runClient()
}
