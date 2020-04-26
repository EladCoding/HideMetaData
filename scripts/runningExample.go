package scripts


func RunningExample() {
	CreateUsersMap()
	go RunInfrastructure()
	runClient()
}
