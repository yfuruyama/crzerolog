package crzerolog

import "os"

func isCloudRun() bool {
	// There is no obvious way to detect whether the app is running on Clodu Run,
	// so we speculate from env var which is automatically added by Cloud Run.
	// ref. https://cloud.google.com/run/docs/reference/container-contract#env-vars
	// Note: we can't use K_SERVICE or K_REVISION since both are also used in Cloud Functions.
	return os.Getenv("K_CONFIGURATION") != ""
}

func isAppEngineSecond() bool {
	// ref. https://cloud.google.com/appengine/docs/standard/go/runtime#environment_variables
	return os.Getenv("GAE_ENV") == "standard"
}
