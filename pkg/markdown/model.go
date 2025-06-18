package markdown

type DownloadImageWithUrlRespModel struct {
	ImageFilePathList            []string
	FailedToDownloadImageUrlList []string
}

type DownloadImageWithUrlReqModel struct {
	Url           string
	ImageFileName string
}
