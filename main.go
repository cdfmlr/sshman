package main

func main() {
	InitLogger()
	InitDB()
	InitRouters()
	logger.Info("sshman is ready")
	logger.Infof("listening on %s", GlobalConfig.HTTP.Addr)
	logger.Fatal(
		Router.Run(GlobalConfig.HTTP.Addr),
	)
}
