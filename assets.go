package main

import (
	"os"
)

func (cfg apiConfig) ensureAssetsDir() error {
	if _, err := os.Stat(cfg.assetsRoot); os.IsNotExist(err) {
		return os.Mkdir(cfg.assetsRoot, 0755)
	}
	return nil
}

/*
func generatePresignedURL(s3Client *s3.Client, bucket, key string, expireTime time.Duration) (string, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	client := s3.NewPresignClient(s3Client)
	req, err := client.PresignGetObject(context.Background(), input)
	if err != nil {
		return "", err
	}

	return req.URL, nil
}

func (cfg *apiConfig) dbVideoToSignedVideo(video database.Video) (database.Video, error) {
	tokens := strings.Split(*video.VideoURL, ",")
	bucket := tokens[0]
	key := tokens[1]

	url, err := generatePresignedURL(cfg.s3Client, bucket, key, time.Duration(2)+time.Minute)
	if err != nil {
		return database.Video{}, err
	}
	video.VideoURL = &url

	return video, nil
}
*/
