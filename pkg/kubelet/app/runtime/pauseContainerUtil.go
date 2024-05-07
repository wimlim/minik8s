package runtime

const (
	PauseContainerImageRef = "registry.cn-hangzhou.aliyuncs.com/google_containers/pause:3.6"
)

/*
 * createPauseContainer
 * 参数：
 * 返回：pauseContainer的Id，error
 */
 
func createPauseContainer(pod *apiobj.Pod) (string error) {
	image.PullImage(PauseContainerImageRef)
	
}